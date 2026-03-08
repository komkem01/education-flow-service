package studentgradeitems

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/auth"
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
	SubjectAssignmentID string   `json:"subject_assignment_id" binding:"required,uuid"`
	Name                *string  `json:"name"`
	MaxScore            *float64 `json:"max_score"`
	WeightPercentage    *float64 `json:"weight_percentage"`
}

type updateRequest = createRequest

type response struct {
	ID                  string   `json:"id"`
	SubjectAssignmentID string   `json:"subject_assignment_id"`
	Name                *string  `json:"name"`
	MaxScore            *float64 `json:"max_score"`
	WeightPercentage    *float64 `json:"weight_percentage"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{SchoolID: claims.SchoolID, StudentID: studentID, SubjectAssignmentID: subjectAssignmentID, Name: req.Name, MaxScore: req.MaxScore, WeightPercentage: req.WeightPercentage})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentNotFound, nil)
			return
		}
		log.Errf("student-grade-items.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	items, err := c.svc.ListByStudentID(ctx.Request.Context(), claims.SchoolID, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.Success(ctx, []response{})
			return
		}
		log.Errf("student-grade-items.list.error: %v", err)
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), claims.SchoolID, studentID, childID, &UpdateInput{SchoolID: claims.SchoolID, StudentID: studentID, SubjectAssignmentID: subjectAssignmentID, Name: req.Name, MaxScore: req.MaxScore, WeightPercentage: req.WeightPercentage})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.GradeItemNotFound, nil)
			return
		}
		log.Errf("student-grade-items.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), claims.SchoolID, studentID, childID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.GradeItemNotFound, nil)
			return
		}
		log.Errf("student-grade-items.delete.error: %v", err)
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

func toResponse(item *ent.GradeItem) response {
	return response{ID: item.ID.String(), SubjectAssignmentID: item.SubjectAssignmentID.String(), Name: item.Name, MaxScore: item.MaxScore, WeightPercentage: item.WeightPercentage}
}
