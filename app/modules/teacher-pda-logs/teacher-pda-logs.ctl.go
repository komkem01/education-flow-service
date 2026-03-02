package teacherpdalogs

import (
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
	CourseName     *string `json:"course_name" binding:"omitempty,max=255"`
	Hours          *int    `json:"hours" binding:"omitempty,min=0"`
	CertificateURL *string `json:"certificate_url" binding:"omitempty,url,max=2048"`
}

type response struct {
	ID             string  `json:"id"`
	TeacherID      string  `json:"teacher_id"`
	CourseName     *string `json:"course_name"`
	Hours          *int    `json:"hours"`
	CertificateURL *string `json:"certificate_url"`
	CreatedAt      string  `json:"created_at"`
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
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{TeacherID: teacherID, CourseName: req.CourseName, Hours: req.Hours, CertificateURL: req.CertificateURL})
	if err != nil {
		log.Errf("teacher-pda-logs.create.error: %v", err)
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
		log.Errf("teacher-pda-logs.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}
	base.Success(ctx, responseList)
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	_, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), childID); err != nil {
		log.Errf("teacher-pda-logs.delete.error: %v", err)
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

func toResponse(item *ent.TeacherPDALog) response {
	return response{ID: item.ID.String(), TeacherID: item.TeacherID.String(), CourseName: item.CourseName, Hours: item.Hours, CertificateURL: item.CertificateURL, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}
