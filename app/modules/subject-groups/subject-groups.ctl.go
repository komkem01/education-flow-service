package subjectgroups

import (
	"database/sql"
	"errors"
	"strconv"
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

type createRequest struct {
	Code        string  `json:"code" binding:"required,min=1,max=50"`
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Head        *string `json:"head" binding:"omitempty,max=255"`
	Description *string `json:"description" binding:"omitempty,max=4000"`
	IsActive    *bool   `json:"is_active"`
}

type updateRequest = createRequest

type response struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Head        *string `json:"head"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectGroupInput{Code: req.Code, Name: req.Name, Head: req.Head, Description: req.Description, IsActive: isActive})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, "subject-group-duplicate", nil)
			return
		}
		log.Errf("subject-groups.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), onlyActive)
	if err != nil {
		log.Errf("subject-groups.list.error: %v", err)
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
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-group-not-found", nil)
			return
		}
		log.Errf("subject-groups.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateSubjectGroupInput{Code: req.Code, Name: req.Name, Head: req.Head, Description: req.Description, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "subject-group-not-found", nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, "subject-group-duplicate", nil)
			return
		}
		log.Errf("subject-groups.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("subject-groups.delete.error: %v", err)
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

func toResponse(item *ent.SubjectGroup) response {
	return response{
		ID:          item.ID.String(),
		Code:        item.Code,
		Name:        item.Name,
		Head:        item.Head,
		Description: item.Description,
		IsActive:    item.IsActive,
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
	return strings.Contains(constraint, "subject_groups")
}
