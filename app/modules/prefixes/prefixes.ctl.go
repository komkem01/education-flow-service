package prefixes

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
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
)

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}

type prefixURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createPrefixRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	IsActive bool   `json:"is_active"`
}

type updatePrefixRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=20"`
	IsActive bool   `json:"is_active"`
}

type prefixResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createPrefixRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	prefix, err := c.svc.Create(ctx.Request.Context(), &CreatePrefixInput{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.PrefixNameDuplicate, nil)
			return
		}

		log.Errf("prefixes.create.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toPrefixResponse(prefix))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}
	prefixes, err := c.svc.List(ctx.Request.Context(), onlyActive)
	if err != nil {
		log.Errf("prefixes.list.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	response := make([]prefixResponse, 0, len(prefixes))
	for _, prefix := range prefixes {
		response = append(response, toPrefixResponse(prefix))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parsePrefixID(ctx)
	if !ok {
		return
	}

	prefix, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "prefix-not-found", nil)
			return
		}

		log.Errf("prefixes.get.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toPrefixResponse(prefix))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parsePrefixID(ctx)
	if !ok {
		return
	}

	var req updatePrefixRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	prefix, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdatePrefixInput{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "prefix-not-found", nil)
			return
		}

		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.PrefixNameDuplicate, nil)
			return
		}

		log.Errf("prefixes.update.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toPrefixResponse(prefix))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parsePrefixID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("prefixes.delete.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parsePrefixID(ctx *gin.Context) (uuid.UUID, bool) {
	var req prefixURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return uuid.Nil, false
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, "invalid-id", nil)
		return uuid.Nil, false
	}

	return id, true
}

func toPrefixResponse(prefix *ent.Prefix) prefixResponse {
	return prefixResponse{
		ID:       prefix.ID.String(),
		Name:     prefix.Name,
		IsActive: prefix.IsActive,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
