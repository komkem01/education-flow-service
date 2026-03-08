package schoolcalendarevents

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"education-flow/app/modules/auth"
	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const dateLayout = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	SchoolID          string  `json:"school_id" binding:"required,uuid"`
	CreatedByMemberID *string `json:"created_by_member_id" binding:"omitempty,uuid"`
	Title             string  `json:"title" binding:"required,min=1,max=255"`
	Description       *string `json:"description"`
	EventType         string  `json:"event_type" binding:"required,oneof=holiday exam activity meeting other"`
	StartDate         string  `json:"start_date" binding:"required"`
	EndDate           *string `json:"end_date"`
	IsActive          *bool   `json:"is_active"`
}

type updateRequest = createRequest

type response struct {
	ID                string  `json:"id"`
	SchoolID          string  `json:"school_id"`
	CreatedByMemberID *string `json:"created_by_member_id"`
	Title             string  `json:"title"`
	Description       *string `json:"description"`
	EventType         string  `json:"event_type"`
	StartDate         string  `json:"start_date"`
	EndDate           *string `json:"end_date"`
	IsActive          bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, creatorID, startDate, endDate, ok := parseCreateUpdateFields(ctx, req)
	if !ok {
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{
		SchoolID:          schoolID,
		CreatedByMemberID: creatorID,
		Title:             req.Title,
		Description:       req.Description,
		EventType:         ent.ToSchoolCalendarEventType(req.EventType),
		StartDate:         startDate,
		EndDate:           endDate,
		IsActive:          isActive,
	})
	if err != nil {
		log.Errf("school-calendar-events.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok {
		if schoolID != nil && *schoolID != claims.SchoolID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		schoolID = &claims.SchoolID
	}

	var eventType *ent.SchoolCalendarEventType
	if raw := strings.TrimSpace(ctx.Query("event_type")); raw != "" {
		if raw != "holiday" && raw != "exam" && raw != "activity" && raw != "meeting" && raw != "other" {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		v := ent.ToSchoolCalendarEventType(raw)
		eventType = &v
	}

	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "true"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), schoolID, eventType, onlyActive)
	if err != nil {
		log.Errf("school-calendar-events.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	out := make([]response, 0, len(items))
	for _, item := range items {
		out = append(out, toResponse(item))
	}
	base.Success(ctx, out)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "school-calendar-event-not-found", nil)
			return
		}
		log.Errf("school-calendar-events.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && item.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, creatorID, startDate, endDate, ok := parseCreateUpdateFields(ctx, req)
	if !ok {
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{
		SchoolID:          schoolID,
		CreatedByMemberID: creatorID,
		Title:             req.Title,
		Description:       req.Description,
		EventType:         ent.ToSchoolCalendarEventType(req.EventType),
		StartDate:         startDate,
		EndDate:           endDate,
		IsActive:          isActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "school-calendar-event-not-found", nil)
			return
		}
		log.Errf("school-calendar-events.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("school-calendar-events.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": id.String()})
}

func parseID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func parseCreateUpdateFields(ctx *gin.Context, req createRequest) (uuid.UUID, *uuid.UUID, time.Time, *time.Time, bool) {
	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, time.Time{}, nil, false
	}

	var creatorID *uuid.UUID
	if req.CreatedByMemberID != nil {
		v, parseErr := uuid.Parse(*req.CreatedByMemberID)
		if parseErr != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return uuid.Nil, nil, time.Time{}, nil, false
		}
		creatorID = &v
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, nil, time.Time{}, nil, false
	}

	var endDate *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		parsed, parseErr := time.Parse(dateLayout, strings.TrimSpace(*req.EndDate))
		if parseErr != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return uuid.Nil, nil, time.Time{}, nil, false
		}
		endDate = &parsed
	}

	if endDate != nil && endDate.Before(startDate) {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, nil, time.Time{}, nil, false
	}

	return schoolID, creatorID, startDate, endDate, true
}

func toResponse(item *ent.SchoolCalendarEvent) response {
	var creatorID *string
	if item.CreatedByMemberID != nil {
		v := item.CreatedByMemberID.String()
		creatorID = &v
	}
	var endDate *string
	if item.EndDate != nil {
		v := item.EndDate.UTC().Format(dateLayout)
		endDate = &v
	}

	return response{
		ID:                item.ID.String(),
		SchoolID:          item.SchoolID.String(),
		CreatedByMemberID: creatorID,
		Title:             item.Title,
		Description:       item.Description,
		EventType:         string(item.EventType),
		StartDate:         item.StartDate.UTC().Format(dateLayout),
		EndDate:           endDate,
		IsActive:          item.IsActive,
	}
}
