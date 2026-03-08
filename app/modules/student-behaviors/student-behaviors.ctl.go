package studentbehaviors

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
	SchoolID     string  `json:"school_id" binding:"required,uuid"`
	StudentID    string  `json:"student_id" binding:"required,uuid"`
	BehaviorType string  `json:"behavior_type" binding:"required,oneof=good bad"`
	Category     *string `json:"category" binding:"omitempty,max=255"`
	Description  *string `json:"description"`
	Points       int     `json:"points"`
	RecordedOn   string  `json:"recorded_on" binding:"required"`
	IsActive     *bool   `json:"is_active"`
}

type updateRequest = createRequest

type response struct {
	ID                 string  `json:"id"`
	SchoolID           string  `json:"school_id"`
	StudentID          string  `json:"student_id"`
	RecordedByMemberID string  `json:"recorded_by_member_id"`
	BehaviorType       string  `json:"behavior_type"`
	Category           *string `json:"category"`
	Description        *string `json:"description"`
	Points             int     `json:"points"`
	RecordedOn         string  `json:"recorded_on"`
	IsActive           bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	claims, hasClaims := auth.GetClaimsFromGin(ctx)
	if !hasClaims {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	schoolID, studentID, recordedOn, ok := parseCreateUpdateFields(ctx, req)
	if !ok {
		return
	}
	if claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{
		SchoolID:           schoolID,
		StudentID:          studentID,
		RecordedByMemberID: claims.MemberID,
		BehaviorType:       ent.ToStudentBehaviorType(req.BehaviorType),
		Category:           req.Category,
		Description:        req.Description,
		Points:             req.Points,
		RecordedOn:         recordedOn,
		IsActive:           isActive,
	})
	if err != nil {
		log.Errf("student-behaviors.create.error: %v", err)
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
	studentID, err := utils.ParseQueryUUID(ctx.Query("student_id"))
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

	var behaviorType *ent.StudentBehaviorType
	if raw := strings.TrimSpace(ctx.Query("behavior_type")); raw != "" {
		if raw != "good" && raw != "bad" {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		v := ent.ToStudentBehaviorType(raw)
		behaviorType = &v
	}

	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "true"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), schoolID, studentID, behaviorType, onlyActive)
	if err != nil {
		log.Errf("student-behaviors.list.error: %v", err)
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
			base.ValidateFailed(ctx, "student-behavior-not-found", nil)
			return
		}
		log.Errf("student-behaviors.get.error: %v", err)
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
	claims, hasClaims := auth.GetClaimsFromGin(ctx)
	if !hasClaims {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	schoolID, studentID, recordedOn, ok := parseCreateUpdateFields(ctx, req)
	if !ok {
		return
	}
	if claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{
		SchoolID:           schoolID,
		StudentID:          studentID,
		RecordedByMemberID: claims.MemberID,
		BehaviorType:       ent.ToStudentBehaviorType(req.BehaviorType),
		Category:           req.Category,
		Description:        req.Description,
		Points:             req.Points,
		RecordedOn:         recordedOn,
		IsActive:           isActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "student-behavior-not-found", nil)
			return
		}
		log.Errf("student-behaviors.update.error: %v", err)
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
		log.Errf("student-behaviors.delete.error: %v", err)
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

func parseCreateUpdateFields(ctx *gin.Context, req createRequest) (uuid.UUID, uuid.UUID, time.Time, bool) {
	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, time.Time{}, false
	}
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, time.Time{}, false
	}
	recordedOn, err := time.Parse(dateLayout, req.RecordedOn)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, uuid.Nil, time.Time{}, false
	}
	return schoolID, studentID, recordedOn, true
}

func toResponse(item *ent.StudentBehavior) response {
	return response{
		ID:                 item.ID.String(),
		SchoolID:           item.SchoolID.String(),
		StudentID:          item.StudentID.String(),
		RecordedByMemberID: item.RecordedByMemberID.String(),
		BehaviorType:       string(item.BehaviorType),
		Category:           item.Category,
		Description:        item.Description,
		Points:             item.Points,
		RecordedOn:         item.RecordedOn.UTC().Format(dateLayout),
		IsActive:           item.IsActive,
	}
}
