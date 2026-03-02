package teacherleavelogs

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
	dateTimeLayout = "2006-01-02T15:04:05Z07:00"
)

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
	Type              string  `json:"type" binding:"omitempty,oneof=sick business vacation other"`
	StartDate         *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate           *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Reason            *string `json:"reason"`
	Status            string  `json:"status" binding:"omitempty,oneof=pending approved rejected"`
	ApprovedByStaffID *string `json:"approved_by_staff_id" binding:"omitempty,uuid"`
}

type updateRequest struct {
	Type              string  `json:"type" binding:"required,oneof=sick business vacation other"`
	StartDate         *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate           *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Reason            *string `json:"reason"`
	Status            string  `json:"status" binding:"required,oneof=pending approved rejected"`
	ApprovedByStaffID *string `json:"approved_by_staff_id" binding:"omitempty,uuid"`
}

type response struct {
	ID                string  `json:"id"`
	TeacherID         string  `json:"teacher_id"`
	Type              string  `json:"type"`
	StartDate         *string `json:"start_date"`
	EndDate           *string `json:"end_date"`
	Reason            *string `json:"reason"`
	Status            string  `json:"status"`
	ApprovedByStaffID *string `json:"approved_by_staff_id"`
	CreatedAt         string  `json:"created_at"`
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
	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate)
	if !ok {
		return
	}
	approvedByStaffID, err := parseUUIDPtr(req.ApprovedByStaffID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	leaveType := ent.TeacherLeaveTypeOther
	if req.Type != "" {
		leaveType = ent.ToTeacherLeaveType(req.Type)
	}
	leaveStatus := ent.TeacherLeaveStatusPending
	if req.Status != "" {
		leaveStatus = ent.ToTeacherLeaveStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{TeacherID: teacherID, Type: leaveType, StartDate: startDate, EndDate: endDate, Reason: req.Reason, Status: leaveStatus, ApprovedByStaffID: approvedByStaffID})
	if err != nil {
		log.Errf("teacher-leave-logs.create.error: %v", err)
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
		log.Errf("teacher-leave-logs.list.error: %v", err)
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
	teacherID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate)
	if !ok {
		return
	}
	approvedByStaffID, err := parseUUIDPtr(req.ApprovedByStaffID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), teacherID, childID, &UpdateInput{Type: ent.ToTeacherLeaveType(req.Type), StartDate: startDate, EndDate: endDate, Reason: req.Reason, Status: ent.ToTeacherLeaveStatus(req.Status), ApprovedByStaffID: approvedByStaffID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherLeaveLogNotFound, nil)
			return
		}
		log.Errf("teacher-leave-logs.update.error: %v", err)
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

func parseUUIDPtr(raw *string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func parseDatePtr(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateLayoutOnly, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseDateRange(ctx *gin.Context, startRaw *string, endRaw *string) (*time.Time, *time.Time, bool) {
	startDate, err := parseDatePtr(startRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return nil, nil, false
	}
	endDate, err := parseDatePtr(endRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return nil, nil, false
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		base.ValidateFailed(ctx, ci18n.TeacherInvalidDateRange, nil)
		return nil, nil, false
	}
	return startDate, endDate, true
}

func toResponse(item *ent.TeacherLeaveLog) response {
	return response{ID: item.ID.String(), TeacherID: item.TeacherID.String(), Type: string(item.Type), StartDate: dateToStringPtr(item.StartDate), EndDate: dateToStringPtr(item.EndDate), Reason: item.Reason, Status: string(item.Status), ApprovedByStaffID: uuidToStringPtr(item.ApprovedByStaffID), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func uuidToStringPtr(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	parsed := value.String()
	return &parsed
}

func dateToStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	parsed := value.Format(dateLayoutOnly)
	return &parsed
}
