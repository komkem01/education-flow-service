package teachereducations

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

type uriRequest struct {
	ID      string `uri:"id" binding:"required"`
	ChildID string `uri:"child_id"`
}

type createRequest struct {
	DegreeLevel    *string `json:"degree_level" binding:"omitempty,max=100"`
	DegreeName     *string `json:"degree_name" binding:"omitempty,max=255"`
	Major          *string `json:"major" binding:"omitempty,max=255"`
	University     *string `json:"university" binding:"omitempty,max=255"`
	GraduationYear *string `json:"graduation_year" binding:"omitempty,max=10"`
}

type updateRequest = createRequest

type response struct {
	ID             string  `json:"id"`
	TeacherID      string  `json:"teacher_id"`
	DegreeLevel    *string `json:"degree_level"`
	DegreeName     *string `json:"degree_name"`
	Major          *string `json:"major"`
	University     *string `json:"university"`
	GraduationYear *string `json:"graduation_year"`
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

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{TeacherID: teacherID, DegreeLevel: req.DegreeLevel, DegreeName: req.DegreeName, Major: req.Major, University: req.University, GraduationYear: req.GraduationYear})
	if err != nil {
		log.Errf("teacher-educations.create.error: %v", err)
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
		log.Errf("teacher-educations.list.error: %v", err)
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

	item, err := c.svc.UpdateByID(ctx.Request.Context(), childID, &UpdateInput{DegreeLevel: req.DegreeLevel, DegreeName: req.DegreeName, Major: req.Major, University: req.University, GraduationYear: req.GraduationYear})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherEducationNotFound, nil)
			return
		}
		log.Errf("teacher-educations.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	_, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), childID); err != nil {
		log.Errf("teacher-educations.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": childID.String()})
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

func toResponse(item *ent.TeacherEducation) response {
	return response{ID: item.ID.String(), TeacherID: item.TeacherID.String(), DegreeLevel: item.DegreeLevel, DegreeName: item.DegreeName, Major: item.Major, University: item.University, GraduationYear: item.GraduationYear}
}
