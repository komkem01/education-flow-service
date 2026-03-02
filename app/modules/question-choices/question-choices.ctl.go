package questionchoices

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

type createQuestionChoiceRequest struct {
	QuestionID string  `json:"question_id" binding:"required,uuid"`
	Content    *string `json:"content"`
	IsCorrect  *bool   `json:"is_correct"`
	OrderNo    *int    `json:"order_no"`
}

type updateQuestionChoiceRequest = createQuestionChoiceRequest

type questionChoiceResponse struct {
	ID         string  `json:"id"`
	QuestionID string  `json:"question_id"`
	Content    *string `json:"content"`
	IsCorrect  *bool   `json:"is_correct"`
	OrderNo    *int    `json:"order_no"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createQuestionChoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateQuestionChoiceInput{QuestionID: questionID, Content: req.Content, IsCorrect: req.IsCorrect, OrderNo: req.OrderNo})
	if err != nil {
		log.Errf("question-choices.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionChoiceResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	questionID, err := utils.ParseQueryUUID(ctx.Query("question_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListQuestionChoicesInput{QuestionID: questionID})
	if err != nil {
		log.Errf("question-choices.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]questionChoiceResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toQuestionChoiceResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionChoiceID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.QuestionChoiceNotFound, nil)
			return
		}
		log.Errf("question-choices.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionChoiceResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionChoiceID(ctx)
	if !ok {
		return
	}

	var req updateQuestionChoiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	questionID, err := uuid.Parse(req.QuestionID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateQuestionChoiceInput{QuestionID: questionID, Content: req.Content, IsCorrect: req.IsCorrect, OrderNo: req.OrderNo})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.QuestionChoiceNotFound, nil)
			return
		}
		log.Errf("question-choices.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionChoiceResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionChoiceID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("question-choices.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseQuestionChoiceID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func toQuestionChoiceResponse(item *ent.QuestionChoice) questionChoiceResponse {
	return questionChoiceResponse{ID: item.ID.String(), QuestionID: item.QuestionID.String(), Content: item.Content, IsCorrect: item.IsCorrect, OrderNo: item.OrderNo}
}
