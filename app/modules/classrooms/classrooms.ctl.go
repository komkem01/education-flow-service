package classrooms

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

type createClassroomRequest struct {
	SchoolID         string  `json:"school_id" binding:"required,uuid"`
	Name             string  `json:"name" binding:"required,min=1,max=100"`
	GradeLevel       *string `json:"grade_level" binding:"omitempty,max=50"`
	RoomNo           *string `json:"room_no" binding:"omitempty,max=50"`
	AdvisorTeacherID *string `json:"advisor_teacher_id" binding:"omitempty,uuid"`
}

type updateClassroomRequest = createClassroomRequest

type classroomResponse struct {
	ID               string  `json:"id"`
	SchoolID         string  `json:"school_id"`
	Name             string  `json:"name"`
	GradeLevel       *string `json:"grade_level"`
	RoomNo           *string `json:"room_no"`
	AdvisorTeacherID *string `json:"advisor_teacher_id"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, advisorTeacherID, ok := parseClassroomCreateUpdateFields(ctx, req.SchoolID, req.AdvisorTeacherID)
	if !ok {
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateClassroomInput{SchoolID: schoolID, Name: req.Name, GradeLevel: req.GradeLevel, RoomNo: req.RoomNo, AdvisorTeacherID: advisorTeacherID})
	if err != nil {
		log.Errf("classrooms.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toClassroomResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListClassroomsInput{SchoolID: schoolID})
	if err != nil {
		log.Errf("classrooms.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]classroomResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toClassroomResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseClassroomID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ClassroomNotFound, nil)
			return
		}
		log.Errf("classrooms.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toClassroomResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseClassroomID(ctx)
	if !ok {
		return
	}

	var req updateClassroomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, advisorTeacherID, ok := parseClassroomCreateUpdateFields(ctx, req.SchoolID, req.AdvisorTeacherID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateClassroomInput{SchoolID: schoolID, Name: req.Name, GradeLevel: req.GradeLevel, RoomNo: req.RoomNo, AdvisorTeacherID: advisorTeacherID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ClassroomNotFound, nil)
			return
		}
		log.Errf("classrooms.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toClassroomResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseClassroomID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("classrooms.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseClassroomID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseClassroomCreateUpdateFields(ctx *gin.Context, schoolIDRaw string, advisorTeacherIDRaw *string) (uuid.UUID, *uuid.UUID, bool) {
	schoolID, err := uuid.Parse(schoolIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, false
	}

	advisorTeacherID, err := utils.ParseUUIDPtr(advisorTeacherIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, false
	}

	return schoolID, advisorTeacherID, true
}

func toClassroomResponse(item *ent.Classroom) classroomResponse {
	return classroomResponse{ID: item.ID.String(), SchoolID: item.SchoolID.String(), Name: item.Name, GradeLevel: item.GradeLevel, RoomNo: item.RoomNo, AdvisorTeacherID: utils.UUIDToStringPtr(item.AdvisorTeacherID)}
}
