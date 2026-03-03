package schedules

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

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

type createScheduleRequest struct {
	SubjectAssignmentID string  `json:"subject_assignment_id" binding:"required,uuid"`
	DayOfWeek           string  `json:"day_of_week" binding:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime           *string `json:"start_time" binding:"omitempty,max=8"`
	EndTime             *string `json:"end_time" binding:"omitempty,max=8"`
	PeriodNo            *int    `json:"period_no" binding:"omitempty,gte=1"`
	Note                *string `json:"note" binding:"omitempty,max=4000"`
	IsActive            *bool   `json:"is_active"`
}

type updateScheduleRequest struct {
	SubjectAssignmentID string  `json:"subject_assignment_id" binding:"required,uuid"`
	DayOfWeek           string  `json:"day_of_week" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime           *string `json:"start_time" binding:"omitempty,max=8"`
	EndTime             *string `json:"end_time" binding:"omitempty,max=8"`
	PeriodNo            *int    `json:"period_no" binding:"omitempty,gte=1"`
	Note                *string `json:"note" binding:"omitempty,max=4000"`
	IsActive            *bool   `json:"is_active"`
}

type scheduleResponse struct {
	ID                  string  `json:"id"`
	SubjectAssignmentID string  `json:"subject_assignment_id"`
	DayOfWeek           string  `json:"day_of_week"`
	StartTime           *string `json:"start_time"`
	EndTime             *string `json:"end_time"`
	PeriodNo            *int    `json:"period_no"`
	Note                *string `json:"note"`
	IsActive            bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	startTime, err := parseClockPtr(req.StartTime)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	endTime, err := parseClockPtr(req.EndTime)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	dayOfWeek := ent.ScheduleDayMonday
	if req.DayOfWeek != "" {
		dayOfWeek = ent.ToScheduleDayOfWeek(req.DayOfWeek)
	}

	if startTime != nil && endTime != nil && !endTime.After(*startTime) {
		base.ValidateFailed(ctx, ci18n.ScheduleInvalidTimeRange, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateScheduleInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: dayOfWeek, StartTime: startTime, EndTime: endTime, PeriodNo: req.PeriodNo, Note: req.Note, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrScheduleTeacherConflict) {
			base.ValidateFailed(ctx, ci18n.ScheduleTeacherConflict, nil)
			return
		}
		if errors.Is(err, ErrScheduleClassroomConflict) {
			base.ValidateFailed(ctx, ci18n.ScheduleClassroomConflict, nil)
			return
		}
		if isScheduleDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.ScheduleDuplicate, nil)
			return
		}
		if isScheduleTimeRangeError(err) {
			base.ValidateFailed(ctx, ci18n.ScheduleInvalidTimeRange, nil)
			return
		}
		log.Errf("schedules.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toScheduleResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	subjectAssignmentID, err := utils.ParseQueryUUID(ctx.Query("subject_assignment_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var dayOfWeek *ent.ScheduleDayOfWeek
	if raw := ctx.Query("day_of_week"); raw != "" {
		value, ok := parseScheduleDayOfWeek(raw)
		if !ok {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		dayOfWeek = &value
	}

	var onlyActive *bool
	if raw := strings.TrimSpace(ctx.Query("only_active")); raw != "" {
		value, convErr := strconv.ParseBool(raw)
		if convErr != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		onlyActive = &value
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListSchedulesInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: dayOfWeek, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("schedules.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]scheduleResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toScheduleResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseScheduleID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ScheduleNotFound, nil)
			return
		}
		log.Errf("schedules.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toScheduleResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseScheduleID(ctx)
	if !ok {
		return
	}

	var req updateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	startTime, err := parseClockPtr(req.StartTime)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	endTime, err := parseClockPtr(req.EndTime)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	if startTime != nil && endTime != nil && !endTime.After(*startTime) {
		base.ValidateFailed(ctx, ci18n.ScheduleInvalidTimeRange, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateScheduleInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: ent.ToScheduleDayOfWeek(req.DayOfWeek), StartTime: startTime, EndTime: endTime, PeriodNo: req.PeriodNo, Note: req.Note, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ScheduleNotFound, nil)
			return
		}
		if errors.Is(err, ErrScheduleTeacherConflict) {
			base.ValidateFailed(ctx, ci18n.ScheduleTeacherConflict, nil)
			return
		}
		if errors.Is(err, ErrScheduleClassroomConflict) {
			base.ValidateFailed(ctx, ci18n.ScheduleClassroomConflict, nil)
			return
		}
		if isScheduleDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.ScheduleDuplicate, nil)
			return
		}
		if isScheduleTimeRangeError(err) {
			base.ValidateFailed(ctx, ci18n.ScheduleInvalidTimeRange, nil)
			return
		}
		log.Errf("schedules.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toScheduleResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseScheduleID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("schedules.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseScheduleID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseClockPtr(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	if value, err := time.Parse("15:04:05", *raw); err == nil {
		normalized := normalizeClock(value)
		return &normalized, nil
	}
	value, err := time.Parse("15:04", *raw)
	if err != nil {
		return nil, err
	}
	normalized := normalizeClock(value)
	return &normalized, nil
}

func normalizeClock(value time.Time) time.Time {
	return time.Date(1970, 1, 1, value.Hour(), value.Minute(), value.Second(), 0, time.UTC)
}

func toClockStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format("15:04:05")
	return &formatted
}

func toScheduleResponse(item *ent.Schedule) scheduleResponse {
	return scheduleResponse{ID: item.ID.String(), SubjectAssignmentID: item.SubjectAssignmentID.String(), DayOfWeek: string(item.DayOfWeek), StartTime: toClockStringPtr(item.StartTime), EndTime: toClockStringPtr(item.EndTime), PeriodNo: item.PeriodNo, Note: item.Note, IsActive: item.IsActive}
}

func parseScheduleDayOfWeek(raw string) (ent.ScheduleDayOfWeek, bool) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday":
		return ent.ToScheduleDayOfWeek(value), true
	default:
		return "", false
	}
}

func isScheduleDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "uq_schedules_assignment_day_period") || strings.Contains(constraint, "uq_schedules_assignment_day_timerange")
}

func isScheduleTimeRangeError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23514" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "chk_schedules_time_range")
}
