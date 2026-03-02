package studentgraderecords

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

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	EnrollmentID string   `json:"enrollment_id" binding:"required,uuid"`
	GradeItemID  string   `json:"grade_item_id" binding:"required,uuid"`
	Score        *float64 `json:"score"`
	TeacherNote  *string  `json:"teacher_note"`
}

type updateRequest = createRequest

type response struct {
	ID           string   `json:"id"`
	EnrollmentID string   `json:"enrollment_id"`
	GradeItemID  string   `json:"grade_item_id"`
	Score        *float64 `json:"score"`
	TeacherNote  *string  `json:"teacher_note"`
	UpdatedAt    string   `json:"updated_at"`
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
	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	gradeItemID, err := uuid.Parse(req.GradeItemID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, EnrollmentID: enrollmentID, GradeItemID: gradeItemID, Score: req.Score, TeacherNote: req.TeacherNote})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentNotFound, nil)
			return
		}
		log.Errf("student-grade-records.create.error: %v", err)
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
		log.Errf("student-grade-records.list.error: %v", err)
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
	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	gradeItemID, err := uuid.Parse(req.GradeItemID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{StudentID: studentID, EnrollmentID: enrollmentID, GradeItemID: gradeItemID, Score: req.Score, TeacherNote: req.TeacherNote})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.GradeRecordNotFound, nil)
			return
		}
		log.Errf("student-grade-records.update.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.GradeRecordNotFound, nil)
			return
		}
		log.Errf("student-grade-records.delete.error: %v", err)
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

func toResponse(item *ent.GradeRecord) response {
	return response{ID: item.ID.String(), EnrollmentID: item.EnrollmentID.String(), GradeItemID: item.GradeItemID.String(), Score: item.Score, TeacherNote: item.TeacherNote, UpdatedAt: item.UpdatedAt.UTC().Format(time.RFC3339)}
}
