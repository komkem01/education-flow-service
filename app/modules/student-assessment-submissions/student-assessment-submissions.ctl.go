package studentassessmentsubmissions

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

type createRequest struct {
	AssessmentSetID string   `json:"assessment_set_id" binding:"required,uuid"`
	SubmitTime      *string  `json:"submit_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TotalScore      *float64 `json:"total_score"`
	Status          string   `json:"status" binding:"omitempty,oneof=in_progress submitted graded"`
}

type updateRequest struct {
	AssessmentSetID string   `json:"assessment_set_id" binding:"required,uuid"`
	SubmitTime      *string  `json:"submit_time" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	TotalScore      *float64 `json:"total_score"`
	Status          string   `json:"status" binding:"required,oneof=in_progress submitted graded"`
}

type response struct {
	ID              string   `json:"id"`
	AssessmentSetID string   `json:"assessment_set_id"`
	StudentID       string   `json:"student_id"`
	SubmitTime      *string  `json:"submit_time"`
	TotalScore      *float64 `json:"total_score"`
	Status          string   `json:"status"`
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
	assessmentSetID, err := uuid.Parse(req.AssessmentSetID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	submitTime, err := utils.ParseTimePtrWithLayout(req.SubmitTime, dateTimeLayout)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	status := ent.StudentAssessmentSubmissionStatusInProgress
	if req.Status != "" {
		status = ent.ToStudentAssessmentSubmissionStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, AssessmentSetID: assessmentSetID, SubmitTime: submitTime, TotalScore: req.TotalScore, Status: status})
	if err != nil {
		log.Errf("student-assessment-submissions.create.error: %v", err)
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
		log.Errf("student-assessment-submissions.list.error: %v", err)
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
	assessmentSetID, err := uuid.Parse(req.AssessmentSetID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	submitTime, err := utils.ParseTimePtrWithLayout(req.SubmitTime, dateTimeLayout)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{AssessmentSetID: assessmentSetID, SubmitTime: submitTime, TotalScore: req.TotalScore, Status: ent.ToStudentAssessmentSubmissionStatus(req.Status)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentAssessmentSubmissionNotFound, nil)
			return
		}
		log.Errf("student-assessment-submissions.update.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.StudentAssessmentSubmissionNotFound, nil)
			return
		}
		log.Errf("student-assessment-submissions.delete.error: %v", err)
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

func toResponse(item *ent.StudentAssessmentSubmission) response {
	return response{ID: item.ID.String(), AssessmentSetID: item.AssessmentSetID.String(), StudentID: item.StudentID.String(), SubmitTime: utils.TimeToStringPtrWithLayout(item.SubmitTime, dateTimeLayout), TotalScore: item.TotalScore, Status: string(item.Status)}
}
