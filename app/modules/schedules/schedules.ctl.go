package schedules

import (
	"database/sql"
	"errors"
	"time"

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

type createScheduleRequest struct {
	SubjectAssignmentID string  `json:"subject_assignment_id" binding:"required,uuid"`
	DayOfWeek           string  `json:"day_of_week" binding:"omitempty,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime           *string `json:"start_time" binding:"omitempty,max=8"`
	EndTime             *string `json:"end_time" binding:"omitempty,max=8"`
	PeriodNo            *int    `json:"period_no"`
}

type updateScheduleRequest struct {
	SubjectAssignmentID string  `json:"subject_assignment_id" binding:"required,uuid"`
	DayOfWeek           string  `json:"day_of_week" binding:"required,oneof=monday tuesday wednesday thursday friday saturday sunday"`
	StartTime           *string `json:"start_time" binding:"omitempty,max=8"`
	EndTime             *string `json:"end_time" binding:"omitempty,max=8"`
	PeriodNo            *int    `json:"period_no"`
}

type scheduleResponse struct {
	ID                  string  `json:"id"`
	SubjectAssignmentID string  `json:"subject_assignment_id"`
	DayOfWeek           string  `json:"day_of_week"`
	StartTime           *string `json:"start_time"`
	EndTime             *string `json:"end_time"`
	PeriodNo            *int    `json:"period_no"`
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

	item, err := c.svc.Create(ctx.Request.Context(), &CreateScheduleInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: dayOfWeek, StartTime: startTime, EndTime: endTime, PeriodNo: req.PeriodNo})
	if err != nil {
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
		value := ent.ToScheduleDayOfWeek(raw)
		dayOfWeek = &value
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListSchedulesInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: dayOfWeek})
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

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateScheduleInput{SubjectAssignmentID: subjectAssignmentID, DayOfWeek: ent.ToScheduleDayOfWeek(req.DayOfWeek), StartTime: startTime, EndTime: endTime, PeriodNo: req.PeriodNo})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ScheduleNotFound, nil)
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
		return &value, nil
	}
	value, err := time.Parse("15:04", *raw)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func toClockStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format("15:04:05")
	return &formatted
}

func toScheduleResponse(item *ent.Schedule) scheduleResponse {
	return scheduleResponse{ID: item.ID.String(), SubjectAssignmentID: item.SubjectAssignmentID.String(), DayOfWeek: string(item.DayOfWeek), StartTime: toClockStringPtr(item.StartTime), EndTime: toClockStringPtr(item.EndTime), PeriodNo: item.PeriodNo}
}
