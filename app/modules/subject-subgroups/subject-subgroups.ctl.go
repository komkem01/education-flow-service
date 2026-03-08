package subjectsubgroups

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"education-flow/app/modules/auth"
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

type createRequest struct {
	SubjectGroupID string  `json:"subject_group_id" binding:"required,uuid"`
	Code           string  `json:"code" binding:"required,min=1,max=50"`
	Name           string  `json:"name" binding:"required,min=1,max=255"`
	Description    *string `json:"description" binding:"omitempty,max=4000"`
	IsActive       *bool   `json:"is_active"`
}

type updateRequest = createRequest

type response struct {
	ID             string  `json:"id"`
	SubjectGroupID string  `json:"subject_group_id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	IsActive       bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectGroupID, err := uuid.Parse(req.SubjectGroupID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectSubgroupInput{SchoolID: claims.SchoolID, SubjectGroupID: subjectGroupID, Code: req.Code, Name: req.Name, Description: req.Description, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-group-not-found", nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, "subject-subgroup-duplicate", nil)
			return
		}
		log.Errf("subject-subgroups.create.error: %v", err)
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

	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectGroupID, err := utils.ParseQueryUUID(ctx.Query("subject_group_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListSubjectSubgroupsInput{SchoolID: claims.SchoolID, SubjectGroupID: subjectGroupID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("subject-subgroups.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}
	base.Success(ctx, responseList)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByIDInSchool(ctx.Request.Context(), claims.SchoolID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-subgroup-not-found", nil)
			return
		}
		log.Errf("subject-subgroups.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectGroupID, err := uuid.Parse(req.SubjectGroupID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByIDInSchool(ctx.Request.Context(), claims.SchoolID, id, &UpdateSubjectSubgroupInput{SubjectGroupID: subjectGroupID, Code: req.Code, Name: req.Name, Description: req.Description, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-subgroup-not-found", nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, "subject-subgroup-duplicate", nil)
			return
		}
		log.Errf("subject-subgroups.update.error: %v", err)
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

	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByIDInSchool(ctx.Request.Context(), claims.SchoolID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-subgroup-not-found", nil)
			return
		}
		log.Errf("subject-subgroups.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func toResponse(item *ent.SubjectSubgroup) response {
	return response{
		ID:             item.ID.String(),
		SubjectGroupID: item.SubjectGroupID.String(),
		Code:           item.Code,
		Name:           item.Name,
		Description:    item.Description,
		IsActive:       item.IsActive,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}
	if pgErr.Code != "23505" {
		return false
	}
	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "subject_subgroups")
}
