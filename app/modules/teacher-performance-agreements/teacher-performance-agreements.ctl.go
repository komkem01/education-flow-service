package teacherperformanceagreements

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const dateTimeLayout = "2006-01-02T15:04:05Z07:00"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type uriRequest struct {
	ID      string `uri:"id" binding:"required"`
	ChildID string `uri:"child_id"`
}

type createRequest struct {
	AcademicYearID   string  `json:"academic_year_id" binding:"required,uuid"`
	AgreementDetail  *string `json:"agreement_detail"`
	ExpectedOutcomes *string `json:"expected_outcomes"`
	Status           string  `json:"status" binding:"omitempty,oneof=draft active completed"`
}

type updateRequest struct {
	AcademicYearID   string  `json:"academic_year_id" binding:"required,uuid"`
	AgreementDetail  *string `json:"agreement_detail"`
	ExpectedOutcomes *string `json:"expected_outcomes"`
	Status           string  `json:"status" binding:"required,oneof=draft active completed"`
}

type response struct {
	ID               string  `json:"id"`
	TeacherID        string  `json:"teacher_id"`
	AcademicYearID   string  `json:"academic_year_id"`
	AgreementDetail  *string `json:"agreement_detail"`
	ExpectedOutcomes *string `json:"expected_outcomes"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	status := ent.TeacherPerformanceAgreementStatusDraft
	if req.Status != "" {
		status = ent.ToTeacherPerformanceAgreementStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{TeacherID: teacherID, AcademicYearID: academicYearID, AgreementDetail: req.AgreementDetail, ExpectedOutcomes: req.ExpectedOutcomes, Status: status})
	if err != nil {
		log.Errf("teacher-performance-agreements.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	items, err := c.svc.ListByTeacherID(ctx.Request.Context(), teacherID)
	if err != nil {
		log.Errf("teacher-performance-agreements.list.error: %v", err)
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
	_, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	academicYearID, err := uuid.Parse(req.AcademicYearID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), childID, &UpdateInput{AcademicYearID: academicYearID, AgreementDetail: req.AgreementDetail, ExpectedOutcomes: req.ExpectedOutcomes, Status: ent.ToTeacherPerformanceAgreementStatus(req.Status)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherPerformanceAgreementNotFound, nil)
			return
		}
		log.Errf("teacher-performance-agreements.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func parseIDs(ctx *gin.Context, childRequired bool) (uuid.UUID, uuid.UUID, bool) {
	var req uriRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, false
	}
	teacherID, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	if !childRequired {
		return teacherID, uuid.Nil, true
	}
	childID, err := uuid.Parse(req.ChildID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	return teacherID, childID, true
}

func toResponse(item *ent.TeacherPerformanceAgreement) response {
	return response{ID: item.ID.String(), TeacherID: item.TeacherID.String(), AcademicYearID: item.AcademicYearID.String(), AgreementDetail: item.AgreementDetail, ExpectedOutcomes: item.ExpectedOutcomes, Status: string(item.Status), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout), UpdatedAt: item.UpdatedAt.UTC().Format(dateTimeLayout)}
}
