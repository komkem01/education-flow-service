package assessmentsets

import (
	"database/sql"
	"errors"
	"strconv"

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

type createAssessmentSetRequest struct {
	SubjectAssignmentID string   `json:"subject_assignment_id" binding:"required,uuid"`
	Title               *string  `json:"title" binding:"omitempty,max=255"`
	DurationMinutes     *int     `json:"duration_minutes"`
	TotalScore          *float64 `json:"total_score"`
	IsPublished         *bool    `json:"is_published"`
}

type updateAssessmentSetRequest = createAssessmentSetRequest

type assessmentSetResponse struct {
	ID                  string   `json:"id"`
	SubjectAssignmentID string   `json:"subject_assignment_id"`
	Title               *string  `json:"title"`
	DurationMinutes     *int     `json:"duration_minutes"`
	TotalScore          *float64 `json:"total_score"`
	IsPublished         bool     `json:"is_published"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createAssessmentSetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	isPublished := false
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateAssessmentSetInput{SubjectAssignmentID: subjectAssignmentID, Title: req.Title, DurationMinutes: req.DurationMinutes, TotalScore: req.TotalScore, IsPublished: isPublished})
	if err != nil {
		log.Errf("assessment-sets.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAssessmentSetResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	subjectAssignmentID, err := utils.ParseQueryUUID(ctx.Query("subject_assignment_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	onlyPublished, err := strconv.ParseBool(ctx.DefaultQuery("only_published", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListAssessmentSetsInput{SubjectAssignmentID: subjectAssignmentID, OnlyPublished: onlyPublished})
	if err != nil {
		log.Errf("assessment-sets.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]assessmentSetResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toAssessmentSetResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseAssessmentSetID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AssessmentSetNotFound, nil)
			return
		}
		log.Errf("assessment-sets.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAssessmentSetResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseAssessmentSetID(ctx)
	if !ok {
		return
	}

	var req updateAssessmentSetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	isPublished := false
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateAssessmentSetInput{SubjectAssignmentID: subjectAssignmentID, Title: req.Title, DurationMinutes: req.DurationMinutes, TotalScore: req.TotalScore, IsPublished: isPublished})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AssessmentSetNotFound, nil)
			return
		}
		log.Errf("assessment-sets.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAssessmentSetResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseAssessmentSetID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("assessment-sets.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseAssessmentSetID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func toAssessmentSetResponse(item *ent.AssessmentSet) assessmentSetResponse {
	return assessmentSetResponse{ID: item.ID.String(), SubjectAssignmentID: item.SubjectAssignmentID.String(), Title: item.Title, DurationMinutes: item.DurationMinutes, TotalScore: item.TotalScore, IsPublished: item.IsPublished}
}
