package academicyears

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

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

const dateLayout = "2006-01-02"

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

type academicYearURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createAcademicYearRequest struct {
	Year      string `json:"year" binding:"required,min=1,max=9"`
	Term      string `json:"term" binding:"required,min=1,max=20"`
	IsCurrent bool   `json:"is_current"`
	IsActive  bool   `json:"is_active"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02"`
}

type updateAcademicYearRequest struct {
	Year      string `json:"year" binding:"required,min=1,max=9"`
	Term      string `json:"term" binding:"required,min=1,max=20"`
	IsCurrent bool   `json:"is_current"`
	IsActive  bool   `json:"is_active"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" binding:"required,datetime=2006-01-02"`
}

type academicYearResponse struct {
	ID        string `json:"id"`
	Year      string `json:"year"`
	Term      string `json:"term"`
	IsCurrent bool   `json:"is_current"`
	IsActive  bool   `json:"is_active"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var req createAcademicYearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate)
	if !ok {
		return
	}

	academicYear, err := c.svc.Create(ctx.Request.Context(), &CreateAcademicYearInput{
		SchoolID:  claims.SchoolID,
		Year:      req.Year,
		Term:      req.Term,
		IsCurrent: req.IsCurrent,
		IsActive:  req.IsActive,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		if errors.Is(err, errAcademicYearYearOutOfRange) {
			base.ValidateFailed(ctx, ci18n.AcademicYearYearOutOfRange, nil)
			return
		}

		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.AcademicYearDuplicate, nil)
			return
		}

		log.Errf("academic-years.create.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toAcademicYearResponse(academicYear))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	if schoolIDQuery := ctx.Query("school_id"); schoolIDQuery != "" {
		schoolID, err := uuid.Parse(schoolIDQuery)
		if err != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return
		}
		if schoolID != claims.SchoolID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
	}

	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}
	onlyCurrent, err := strconv.ParseBool(ctx.DefaultQuery("only_current", "false"))
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	academicYears, err := c.svc.List(ctx.Request.Context(), claims.SchoolID, onlyActive, onlyCurrent)
	if err != nil {
		log.Errf("academic-years.list.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	response := make([]academicYearResponse, 0, len(academicYears))
	for _, academicYear := range academicYears {
		response = append(response, toAcademicYearResponse(academicYear))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseAcademicYearID(ctx)
	if !ok {
		return
	}

	academicYear, err := c.svc.GetByIDInSchool(ctx.Request.Context(), claims.SchoolID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AcademicYearNotFound, nil)
			return
		}

		log.Errf("academic-years.get.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toAcademicYearResponse(academicYear))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseAcademicYearID(ctx)
	if !ok {
		return
	}

	var req updateAcademicYearRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate)
	if !ok {
		return
	}

	academicYear, err := c.svc.UpdateByIDInSchool(ctx.Request.Context(), id, &UpdateAcademicYearInput{
		SchoolID:  claims.SchoolID,
		Year:      req.Year,
		Term:      req.Term,
		IsCurrent: req.IsCurrent,
		IsActive:  req.IsActive,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		if errors.Is(err, errAcademicYearYearOutOfRange) {
			base.ValidateFailed(ctx, ci18n.AcademicYearYearOutOfRange, nil)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AcademicYearNotFound, nil)
			return
		}

		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.AcademicYearDuplicate, nil)
			return
		}

		log.Errf("academic-years.update.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, toAcademicYearResponse(academicYear))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseAcademicYearID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByIDInSchool(ctx.Request.Context(), claims.SchoolID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AcademicYearNotFound, nil)
			return
		}
		log.Errf("academic-years.delete.error: %v", err)
		base.InternalServerError(ctx, "internal-server-error", nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseAcademicYearID(ctx *gin.Context) (uuid.UUID, bool) {
	var req academicYearURIRequest
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

func parseDateRange(ctx *gin.Context, startDateRaw string, endDateRaw string) (time.Time, time.Time, bool) {
	startDate, err := time.Parse(dateLayout, startDateRaw)
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return time.Time{}, time.Time{}, false
	}

	endDate, err := time.Parse(dateLayout, endDateRaw)
	if err != nil {
		base.BadRequest(ctx, "bad-request", nil)
		return time.Time{}, time.Time{}, false
	}

	if endDate.Before(startDate) {
		base.ValidateFailed(ctx, "academic-year-invalid-date-range", nil)
		return time.Time{}, time.Time{}, false
	}

	return startDate, endDate, true
}

func toAcademicYearResponse(academicYear *ent.AcademicYear) academicYearResponse {
	return academicYearResponse{
		ID:        academicYear.ID.String(),
		Year:      academicYear.Year,
		Term:      academicYear.Term,
		IsCurrent: academicYear.IsCurrent,
		IsActive:  academicYear.IsActive,
		StartDate: academicYear.StartDate.Format(dateLayout),
		EndDate:   academicYear.EndDate.Format(dateLayout),
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
