package questionbanks

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

type createQuestionBankRequest struct {
	SubjectID       string  `json:"subject_id" binding:"required,uuid"`
	TeacherID       string  `json:"teacher_id" binding:"required,uuid"`
	Content         *string `json:"content"`
	Type            string  `json:"type" binding:"omitempty,oneof=multiple_choice true_false short_answer essay"`
	DifficultyLevel *int    `json:"difficulty_level" binding:"omitempty,min=1,max=5"`
	IndicatorCode   *string `json:"indicator_code" binding:"omitempty,max=100"`
	Tags            *string `json:"tags" binding:"omitempty,max=255"`
}

type updateQuestionBankRequest struct {
	SubjectID       string  `json:"subject_id" binding:"required,uuid"`
	TeacherID       string  `json:"teacher_id" binding:"required,uuid"`
	Content         *string `json:"content"`
	Type            string  `json:"type" binding:"required,oneof=multiple_choice true_false short_answer essay"`
	DifficultyLevel *int    `json:"difficulty_level" binding:"omitempty,min=1,max=5"`
	IndicatorCode   *string `json:"indicator_code" binding:"omitempty,max=100"`
	Tags            *string `json:"tags" binding:"omitempty,max=255"`
}

type questionBankResponse struct {
	ID              string  `json:"id"`
	SubjectID       string  `json:"subject_id"`
	TeacherID       string  `json:"teacher_id"`
	Content         *string `json:"content"`
	Type            string  `json:"type"`
	DifficultyLevel *int    `json:"difficulty_level"`
	IndicatorCode   *string `json:"indicator_code"`
	Tags            *string `json:"tags"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createQuestionBankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, teacherID, ok := parseQuestionBankCreateUpdateFields(ctx, req.SubjectID, req.TeacherID)
	if !ok {
		return
	}

	questionType := ent.QuestionBankTypeMultipleChoice
	if req.Type != "" {
		questionType = ent.ToQuestionBankType(req.Type)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateQuestionBankInput{SubjectID: subjectID, TeacherID: teacherID, Content: req.Content, Type: questionType, DifficultyLevel: req.DifficultyLevel, IndicatorCode: req.IndicatorCode, Tags: req.Tags})
	if err != nil {
		log.Errf("question-banks.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionBankResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	subjectID, err := utils.ParseQueryUUID(ctx.Query("subject_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	teacherID, err := utils.ParseQueryUUID(ctx.Query("teacher_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var questionType *ent.QuestionBankType
	if raw := ctx.Query("type"); raw != "" {
		value := ent.ToQuestionBankType(raw)
		questionType = &value
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListQuestionBanksInput{SubjectID: subjectID, TeacherID: teacherID, Type: questionType})
	if err != nil {
		log.Errf("question-banks.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]questionBankResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toQuestionBankResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionBankID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.QuestionBankNotFound, nil)
			return
		}
		log.Errf("question-banks.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionBankResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionBankID(ctx)
	if !ok {
		return
	}

	var req updateQuestionBankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, teacherID, ok := parseQuestionBankCreateUpdateFields(ctx, req.SubjectID, req.TeacherID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateQuestionBankInput{SubjectID: subjectID, TeacherID: teacherID, Content: req.Content, Type: ent.ToQuestionBankType(req.Type), DifficultyLevel: req.DifficultyLevel, IndicatorCode: req.IndicatorCode, Tags: req.Tags})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.QuestionBankNotFound, nil)
			return
		}
		log.Errf("question-banks.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toQuestionBankResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseQuestionBankID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("question-banks.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseQuestionBankID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseQuestionBankCreateUpdateFields(ctx *gin.Context, subjectIDRaw string, teacherIDRaw string) (uuid.UUID, uuid.UUID, bool) {
	subjectID, err := uuid.Parse(subjectIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	teacherID, err := uuid.Parse(teacherIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return subjectID, teacherID, true
}

func toQuestionBankResponse(item *ent.QuestionBank) questionBankResponse {
	return questionBankResponse{ID: item.ID.String(), SubjectID: item.SubjectID.String(), TeacherID: item.TeacherID.String(), Content: item.Content, Type: string(item.Type), DifficultyLevel: item.DifficultyLevel, IndicatorCode: item.IndicatorCode, Tags: item.Tags}
}
