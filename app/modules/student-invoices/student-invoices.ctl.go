package studentinvoices

import (
	"database/sql"
	"errors"
	"time"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const (
	dateLayoutOnly = "2006-01-02"
)

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	FeeCategoryID  string   `json:"fee_category_id" binding:"required,uuid"`
	AcademicYearID string   `json:"academic_year_id" binding:"required,uuid"`
	Amount         *float64 `json:"amount"`
	DueDate        *string  `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Status         string   `json:"status" binding:"omitempty,oneof=unpaid paid partial cancelled"`
}

type updateRequest struct {
	FeeCategoryID  string   `json:"fee_category_id" binding:"required,uuid"`
	AcademicYearID string   `json:"academic_year_id" binding:"required,uuid"`
	Amount         *float64 `json:"amount"`
	DueDate        *string  `json:"due_date" binding:"omitempty,datetime=2006-01-02"`
	Status         string   `json:"status" binding:"required,oneof=unpaid paid partial cancelled"`
}

type response struct {
	ID             string   `json:"id"`
	StudentID      string   `json:"student_id"`
	FeeCategoryID  string   `json:"fee_category_id"`
	AcademicYearID string   `json:"academic_year_id"`
	Amount         *float64 `json:"amount"`
	DueDate        *string  `json:"due_date"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	feeCategoryID, err := uuid.Parse(req.FeeCategoryID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	dueDate, err := utils.ParseTimePtrWithLayout(req.DueDate, dateLayoutOnly)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	status := ent.StudentInvoiceStatusUnpaid
	if req.Status != "" {
		status = ent.ToStudentInvoiceStatus(req.Status)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, FeeCategoryID: feeCategoryID, AcademicYearID: academicYearID, Amount: req.Amount, DueDate: dueDate, Status: status})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.FeeCategoryNotFound, nil)
			return
		}
		log.Errf("student-invoices.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	items, err := c.svc.ListByStudentID(ctx.Request.Context(), studentID)
	if err != nil {
		log.Errf("student-invoices.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}
	base.Success(ctx, responseList)
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	feeCategoryID, err := uuid.Parse(req.FeeCategoryID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	dueDate, err := utils.ParseTimePtrWithLayout(req.DueDate, dateLayoutOnly)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{FeeCategoryID: feeCategoryID, AcademicYearID: academicYearID, Amount: req.Amount, DueDate: dueDate, Status: ent.ToStudentInvoiceStatus(req.Status)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentInvoiceNotFound, nil)
			return
		}
		log.Errf("student-invoices.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), studentID, childID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentInvoiceNotFound, nil)
			return
		}
		log.Errf("student-invoices.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": childID.String()})
}

func parseIDs(ctx *gin.Context, childRequired bool) (uuid.UUID, uuid.UUID, bool) {
	studentID, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	if !childRequired {
		return studentID, uuid.Nil, true
	}
	childID, err := utils.ParsePathUUID(ctx, "child_id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	return studentID, childID, true
}

func toResponse(item *ent.StudentInvoice) response {
	return response{ID: item.ID.String(), StudentID: item.StudentID.String(), FeeCategoryID: item.FeeCategoryID.String(), AcademicYearID: item.AcademicYearID.String(), Amount: item.Amount, DueDate: utils.TimeToStringPtrWithLayout(item.DueDate, dateLayoutOnly), Status: string(item.Status), CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339)}
}
