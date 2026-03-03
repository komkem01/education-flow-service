package subjects

import (
	"database/sql"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
	SchoolID           string   `json:"school_id" binding:"required,uuid"`
	SubjectCode        *string  `json:"subject_code" binding:"omitempty,max=50"`
	Name               string   `json:"name" binding:"required,min=1,max=255"`
	NameEN             *string  `json:"name_en" binding:"omitempty,max=255"`
	Description        *string  `json:"description" binding:"omitempty,max=4000"`
	LearningObjectives *string  `json:"learning_objectives" binding:"omitempty,max=4000"`
	LearningOutcomes   *string  `json:"learning_outcomes" binding:"omitempty,max=4000"`
	AssessmentCriteria *string  `json:"assessment_criteria" binding:"omitempty,max=4000"`
	GradeLevel         *string  `json:"grade_level" binding:"omitempty,max=50"`
	Category           *string  `json:"category" binding:"omitempty,max=100"`
	Credits            *float64 `json:"credits" binding:"omitempty,gte=0"`
	Type               string   `json:"type" binding:"omitempty,oneof=core elective activity"`
	IsActive           *bool    `json:"is_active"`
}

type updateSubjectRequest struct {
	SchoolID           string   `json:"school_id" binding:"required,uuid"`
	SubjectCode        *string  `json:"subject_code" binding:"omitempty,max=50"`
	Name               string   `json:"name" binding:"required,min=1,max=255"`
	NameEN             *string  `json:"name_en" binding:"omitempty,max=255"`
	Description        *string  `json:"description" binding:"omitempty,max=4000"`
	LearningObjectives *string  `json:"learning_objectives" binding:"omitempty,max=4000"`
	LearningOutcomes   *string  `json:"learning_outcomes" binding:"omitempty,max=4000"`
	AssessmentCriteria *string  `json:"assessment_criteria" binding:"omitempty,max=4000"`
	GradeLevel         *string  `json:"grade_level" binding:"omitempty,max=50"`
	Category           *string  `json:"category" binding:"omitempty,max=100"`
	Credits            *float64 `json:"credits" binding:"omitempty,gte=0"`
	Type               string   `json:"type" binding:"required,oneof=core elective activity"`
	IsActive           *bool    `json:"is_active"`
}

type subjectResponse struct {
	ID                 string   `json:"id"`
	SchoolID           string   `json:"school_id"`
	SubjectCode        *string  `json:"subject_code"`
	Name               string   `json:"name"`
	NameEN             *string  `json:"name_en"`
	Description        *string  `json:"description"`
	LearningObjectives *string  `json:"learning_objectives"`
	LearningOutcomes   *string  `json:"learning_outcomes"`
	AssessmentCriteria *string  `json:"assessment_criteria"`
	GradeLevel         *string  `json:"grade_level"`
	Category           *string  `json:"category"`
	Credits            *float64 `json:"credits"`
	Type               string   `json:"type"`
	IsActive           bool     `json:"is_active"`
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectInput{SchoolID: schoolID, SubjectCode: req.SubjectCode, Name: req.Name, NameEN: req.NameEN, Description: req.Description, LearningObjectives: req.LearningObjectives, LearningOutcomes: req.LearningOutcomes, AssessmentCriteria: req.AssessmentCriteria, GradeLevel: req.GradeLevel, Category: req.Category, Credits: req.Credits, Type: subjectType, IsActive: isActive})
	if err != nil {
		if isSubjectCodeDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectCodeDuplicate, nil)
			return
		}
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateSubjectInput{SchoolID: schoolID, SubjectCode: req.SubjectCode, Name: req.Name, NameEN: req.NameEN, Description: req.Description, LearningObjectives: req.LearningObjectives, LearningOutcomes: req.LearningOutcomes, AssessmentCriteria: req.AssessmentCriteria, GradeLevel: req.GradeLevel, Category: req.Category, Credits: req.Credits, Type: ent.ToSubjectType(req.Type), IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectNotFound, nil)
			return
		}
		if isSubjectCodeDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectCodeDuplicate, nil)
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
	return subjectResponse{
		ID:                 item.ID.String(),
		SchoolID:           item.SchoolID.String(),
		SubjectCode:        item.SubjectCode,
		Name:               item.Name,
		NameEN:             item.NameEN,
		Description:        item.Description,
		LearningObjectives: item.LearningObjectives,
		LearningOutcomes:   item.LearningOutcomes,
		AssessmentCriteria: item.AssessmentCriteria,
		GradeLevel:         item.GradeLevel,
		Category:           item.Category,
		Credits:            item.Credits,
		Type:               string(item.Type),
		IsActive:           item.IsActive,
	}
}

func isSubjectCodeDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "uq_subjects_school_subject_code")
}
