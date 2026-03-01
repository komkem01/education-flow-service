package schools

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type schoolURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createSchoolRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,url,max=2048"`
	ThemeColor  *string `json:"theme_color" binding:"omitempty,startswith=#,len=7"`
	Address     string  `json:"address" binding:"required,min=1,max=500"`
	Description *string `json:"description"`
}

type updateSchoolRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,url,max=2048"`
	ThemeColor  *string `json:"theme_color" binding:"omitempty,startswith=#,len=7"`
	Address     string  `json:"address" binding:"required,min=1,max=500"`
	Description *string `json:"description"`
}

type schoolResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	LogoURL     *string `json:"logo_url"`
	ThemeColor  *string `json:"theme_color"`
	Address     string  `json:"address"`
	Description *string `json:"description"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	school, err := c.svc.Create(ctx.Request.Context(), &CreateSchoolInput{
		Name:        req.Name,
		LogoURL:     req.LogoURL,
		ThemeColor:  req.ThemeColor,
		Address:     req.Address,
		Description: req.Description,
	})
	if err != nil {
		log.Errf("schools.create.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toSchoolResponse(school))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schools, err := c.svc.List(ctx.Request.Context())
	if err != nil {
		log.Errf("schools.list.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	response := make([]schoolResponse, 0, len(schools))
	for _, school := range schools {
		response = append(response, toSchoolResponse(school))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSchoolID(ctx)
	if !ok {
		return
	}

	school, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "school-not-found", nil)
			return
		}

		log.Errf("schools.get.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toSchoolResponse(school))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSchoolID(ctx)
	if !ok {
		return
	}

	var req updateSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	school, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateSchoolInput{
		Name:        req.Name,
		LogoURL:     req.LogoURL,
		ThemeColor:  req.ThemeColor,
		Address:     req.Address,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, "school-not-found", nil)
			return
		}

		log.Errf("schools.update.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toSchoolResponse(school))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSchoolID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("schools.delete.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseSchoolID(ctx *gin.Context) (uuid.UUID, bool) {
	var req schoolURIRequest
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

func toSchoolResponse(school *ent.School) schoolResponse {
	return schoolResponse{
		ID:          school.ID.String(),
		Name:        school.Name,
		LogoURL:     school.LogoURL,
		ThemeColor:  school.ThemeColor,
		Address:     school.Address,
		Description: school.Description,
	}
}
