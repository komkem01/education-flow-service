package teacherprofilerequests

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
	RequestedData map[string]any `json:"requested_data"`
	Reason        *string        `json:"reason"`
	Status        string         `json:"status" binding:"omitempty,oneof=pending approved rejected"`
	Comment       *string        `json:"comment"`
}

type updateRequest struct {
	RequestedData      map[string]any `json:"requested_data"`
	Reason             *string        `json:"reason"`
	Status             string         `json:"status" binding:"required,oneof=pending approved rejected"`
	Comment            *string        `json:"comment"`
	ProcessedByStaffID *string        `json:"processed_by_staff_id" binding:"omitempty,uuid"`
	ProcessedAt        *string        `json:"processed_at" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
}

type response struct {
	ID                 string         `json:"id"`
	TeacherID          string         `json:"teacher_id"`
	RequestedData      map[string]any `json:"requested_data"`
	Reason             *string        `json:"reason"`
	Status             string         `json:"status"`
	Comment            *string        `json:"comment"`
	ProcessedByStaffID *string        `json:"processed_by_staff_id"`
	ProcessedAt        *string        `json:"processed_at"`
	CreatedAt          string         `json:"created_at"`
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
	status := ent.TeacherProfileRequestStatusPending
	if req.Status != "" {
		status = ent.ToTeacherProfileRequestStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{TeacherID: teacherID, RequestedData: req.RequestedData, Reason: req.Reason, Status: status, Comment: req.Comment})
	if err != nil {
		log.Errf("teacher-profile-requests.create.error: %v", err)
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
		log.Errf("teacher-profile-requests.list.error: %v", err)
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
	processedByStaffID, err := parseUUIDPtr(req.ProcessedByStaffID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	processedAt, err := parseDateTimePtr(req.ProcessedAt)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), teacherID, childID, &UpdateInput{RequestedData: req.RequestedData, Reason: req.Reason, Status: ent.ToTeacherProfileRequestStatus(req.Status), Comment: req.Comment, ProcessedByStaffID: processedByStaffID, ProcessedAt: processedAt})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherProfileRequestNotFound, nil)
			return
		}
		log.Errf("teacher-profile-requests.update.error: %v", err)
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

func parseDateTimePtr(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateTimeLayout, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func toResponse(item *ent.TeacherProfileRequest) response {
	return response{ID: item.ID.String(), TeacherID: item.TeacherID.String(), RequestedData: item.RequestedData, Reason: item.Reason, Status: string(item.Status), Comment: item.Comment, ProcessedByStaffID: uuidToStringPtr(item.ProcessedByStaffID), ProcessedAt: dateTimeToStringPtr(item.ProcessedAt), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func uuidToStringPtr(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	parsed := value.String()
	return &parsed
}

func dateTimeToStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	parsed := value.UTC().Format(dateTimeLayout)
	return &parsed
}
