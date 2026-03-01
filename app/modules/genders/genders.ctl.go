package genders

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

type genderURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createGenderRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=50"`
	IsActive bool   `json:"is_active"`
}

type updateGenderRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=50"`
	IsActive bool   `json:"is_active"`
}

type genderResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createGenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	gender, err := c.svc.Create(ctx.Request.Context(), &CreateGenderInput{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.GenderNameDuplicate, nil)
			return
		}

		log.Errf("genders.create.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toGenderResponse(gender))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}
	genders, err := c.svc.List(ctx.Request.Context(), onlyActive)
	if err != nil {
		log.Errf("genders.list.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	response := make([]genderResponse, 0, len(genders))
	for _, gender := range genders {
		response = append(response, toGenderResponse(gender))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseGenderID(ctx)
	if !ok {
		return
	}

	gender, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "gender-not-found", nil)
			return
		}

		log.Errf("genders.get.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toGenderResponse(gender))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseGenderID(ctx)
	if !ok {
		return
	}

	var req updateGenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	gender, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateGenderInput{
		Name:     req.Name,
		IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "gender-not-found", nil)
			return
		}

		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.GenderNameDuplicate, nil)
			return
		}

		log.Errf("genders.update.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toGenderResponse(gender))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseGenderID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("genders.delete.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseGenderID(ctx *gin.Context) (uuid.UUID, bool) {
	var req genderURIRequest
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

func toGenderResponse(gender *ent.Gender) genderResponse {
	return genderResponse{
		ID:       gender.ID.String(),
		Name:     gender.Name,
		IsActive: gender.IsActive,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
