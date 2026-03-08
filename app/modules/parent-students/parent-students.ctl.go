package parentstudents

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
	StudentID      string `json:"student_id" binding:"required,uuid"`
	Relationship   string `json:"relationship" binding:"omitempty,oneof=father mother guardian"`
	IsMainGuardian *bool  `json:"is_main_guardian"`
}

type updateRequest struct {
	StudentID      string `json:"student_id" binding:"required,uuid"`
	Relationship   string `json:"relationship" binding:"required,oneof=father mother guardian"`
	IsMainGuardian *bool  `json:"is_main_guardian"`
}

type response struct {
	ID             string `json:"id"`
	ParentID       string `json:"parent_id"`
	StudentID      string `json:"student_id"`
	Relationship   string `json:"relationship"`
	IsMainGuardian bool   `json:"is_main_guardian"`
	CreatedAt      string `json:"created_at"`
}

type studentParentResponse struct {
	ID              string  `json:"id"`
	ParentID        string  `json:"parent_id"`
	StudentID       string  `json:"student_id"`
	Relationship    string  `json:"relationship"`
	IsMainGuardian  bool    `json:"is_main_guardian"`
	CreatedAt       string  `json:"created_at"`
	ParentMemberID  string  `json:"parent_member_id"`
	ParentGenderID  *string `json:"parent_gender_id"`
	ParentPrefixID  *string `json:"parent_prefix_id"`
	ParentCode      *string `json:"parent_code"`
	ParentFirstName *string `json:"parent_first_name"`
	ParentLastName  *string `json:"parent_last_name"`
	ParentPhone     *string `json:"parent_phone"`
	ParentIsActive  bool    `json:"parent_is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	relationship := ent.ParentRelationshipGuardian
	if req.Relationship != "" {
		relationship = ent.ToParentRelationship(req.Relationship)
	}
	isMainGuardian := false
	if req.IsMainGuardian != nil {
		isMainGuardian = *req.IsMainGuardian
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{ParentID: parentID, StudentID: studentID, Relationship: relationship, IsMainGuardian: isMainGuardian})
	if err != nil {
		if errors.Is(err, ErrParentNotFound) {
			base.ValidateFailed(ctx, ci18n.ParentNotFound, nil)
			return
		}
		if errors.Is(err, ErrStudentNotFound) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		log.Errf("parent-students.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	items, err := c.svc.ListByParentID(ctx.Request.Context(), parentID)
	if err != nil {
		log.Errf("parent-students.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}

	base.Success(ctx, responseList)
}

func (c *Controller) ListByStudent(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.ListByStudentID(ctx.Request.Context(), studentID)
	if err != nil {
		log.Errf("parent-students.list-by-student.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	responseList := make([]studentParentResponse, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toStudentParentResponse(item))
	}

	base.Success(ctx, responseList)
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	isMainGuardian := false
	if req.IsMainGuardian != nil {
		isMainGuardian = *req.IsMainGuardian
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), parentID, childID, &UpdateInput{StudentID: studentID, Relationship: ent.ToParentRelationship(req.Relationship), IsMainGuardian: isMainGuardian})
	if err != nil {
		if errors.Is(err, ErrParentNotFound) {
			base.ValidateFailed(ctx, ci18n.ParentNotFound, nil)
			return
		}
		if errors.Is(err, ErrStudentNotFound) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ParentStudentNotFound, nil)
			return
		}
		log.Errf("parent-students.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), parentID, childID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ParentStudentNotFound, nil)
			return
		}
		log.Errf("parent-students.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": childID.String()})
}

func parseIDs(ctx *gin.Context, childRequired bool) (uuid.UUID, uuid.UUID, bool) {
	parentID, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	if !childRequired {
		return parentID, uuid.Nil, true
	}
	childID, err := utils.ParsePathUUID(ctx, "child_id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}

	return parentID, childID, true
}

func toResponse(item *ent.MemberParentStudent) response {
	return response{ID: item.ID.String(), ParentID: item.ParentID.String(), StudentID: item.StudentID.String(), Relationship: string(item.Relationship), IsMainGuardian: item.IsMainGuardian, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func toStudentParentResponse(item *ent.StudentParent) studentParentResponse {
	return studentParentResponse{
		ID:              item.ID.String(),
		ParentID:        item.ParentID.String(),
		StudentID:       item.StudentID.String(),
		Relationship:    string(item.Relationship),
		IsMainGuardian:  item.IsMainGuardian,
		CreatedAt:       item.CreatedAt.UTC().Format(dateTimeLayout),
		ParentMemberID:  item.ParentMemberID.String(),
		ParentGenderID:  utils.UUIDToStringPtr(item.ParentGenderID),
		ParentPrefixID:  utils.UUIDToStringPtr(item.ParentPrefixID),
		ParentCode:      item.ParentCode,
		ParentFirstName: item.ParentFirstName,
		ParentLastName:  item.ParentLastName,
		ParentPhone:     item.ParentPhone,
		ParentIsActive:  item.ParentIsActive,
	}
}
