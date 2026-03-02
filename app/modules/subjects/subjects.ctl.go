package subjects

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

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createSubjectRequest struct {
	SchoolID    string   `json:"school_id" binding:"required,uuid"`
	SubjectCode *string  `json:"subject_code" binding:"omitempty,max=50"`
	Name        string   `json:"name" binding:"required,min=1,max=255"`
	Credits     *float64 `json:"credits"`
	Type        string   `json:"type" binding:"omitempty,oneof=core elective activity"`
}

type updateSubjectRequest struct {
	SchoolID    string   `json:"school_id" binding:"required,uuid"`
	SubjectCode *string  `json:"subject_code" binding:"omitempty,max=50"`
	Name        string   `json:"name" binding:"required,min=1,max=255"`
	Credits     *float64 `json:"credits"`
	Type        string   `json:"type" binding:"required,oneof=core elective activity"`
}

type subjectResponse struct {
	ID          string   `json:"id"`
	SchoolID    string   `json:"school_id"`
	SubjectCode *string  `json:"subject_code"`
	Name        string   `json:"name"`
	Credits     *float64 `json:"credits"`
	Type        string   `json:"type"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createSubjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	subjectType := ent.SubjectTypeCore
	if req.Type != "" {
		subjectType = ent.ToSubjectType(req.Type)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectInput{SchoolID: schoolID, SubjectCode: req.SubjectCode, Name: req.Name, Credits: req.Credits, Type: subjectType})
	if err != nil {
		log.Errf("subjects.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListSubjectsInput{SchoolID: schoolID})
	if err != nil {
		log.Errf("subjects.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]subjectResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toSubjectResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectNotFound, nil)
			return
		}
		log.Errf("subjects.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectID(ctx)
	if !ok {
		return
	}

	var req updateSubjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateSubjectInput{SchoolID: schoolID, SubjectCode: req.SubjectCode, Name: req.Name, Credits: req.Credits, Type: ent.ToSubjectType(req.Type)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectNotFound, nil)
			return
		}
		log.Errf("subjects.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("subjects.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseSubjectID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func toSubjectResponse(item *ent.Subject) subjectResponse {
	return subjectResponse{ID: item.ID.String(), SchoolID: item.SchoolID.String(), SubjectCode: item.SubjectCode, Name: item.Name, Credits: item.Credits, Type: string(item.Type)}
}
