package studentattendancelogs

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

const dateLayoutOnly = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	EnrollmentID string  `json:"enrollment_id" binding:"required,uuid"`
	ScheduleID   string  `json:"schedule_id" binding:"required,uuid"`
	CheckDate    *string `json:"check_date" binding:"omitempty,datetime=2006-01-02"`
	Status       string  `json:"status" binding:"omitempty,oneof=present absent late sick business"`
	Note         *string `json:"note"`
}

type updateRequest struct {
	EnrollmentID string  `json:"enrollment_id" binding:"required,uuid"`
	ScheduleID   string  `json:"schedule_id" binding:"required,uuid"`
	CheckDate    *string `json:"check_date" binding:"omitempty,datetime=2006-01-02"`
	Status       string  `json:"status" binding:"required,oneof=present absent late sick business"`
	Note         *string `json:"note"`
}

type response struct {
	ID           string  `json:"id"`
	EnrollmentID string  `json:"enrollment_id"`
	ScheduleID   string  `json:"schedule_id"`
	CheckDate    *string `json:"check_date"`
	Status       string  `json:"status"`
	Note         *string `json:"note"`
	CreatedAt    string  `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	scheduleID, err := uuid.Parse(req.ScheduleID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	checkDate, err := utils.ParseTimePtrWithLayout(req.CheckDate, dateLayoutOnly)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	status := ent.StudentAttendanceStatusPresent
	if req.Status != "" {
		status = ent.ToStudentAttendanceStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, EnrollmentID: enrollmentID, ScheduleID: scheduleID, CheckDate: checkDate, Status: status, Note: req.Note})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentNotFound, nil)
			return
		}
		log.Errf("student-attendance-logs.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	items, err := c.svc.ListByStudentID(ctx.Request.Context(), studentID)
	if err != nil {
		log.Errf("student-attendance-logs.list.error: %v", err)
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
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	enrollmentID, err := uuid.Parse(req.EnrollmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	scheduleID, err := uuid.Parse(req.ScheduleID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	checkDate, err := utils.ParseTimePtrWithLayout(req.CheckDate, dateLayoutOnly)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{EnrollmentID: enrollmentID, ScheduleID: scheduleID, CheckDate: checkDate, Status: ent.ToStudentAttendanceStatus(req.Status), Note: req.Note})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentAttendanceLogNotFound, nil)
			return
		}
		log.Errf("student-attendance-logs.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), studentID, childID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentAttendanceLogNotFound, nil)
			return
		}
		log.Errf("student-attendance-logs.delete.error: %v", err)
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

func toResponse(item *ent.StudentAttendanceLog) response {
	return response{ID: item.ID.String(), EnrollmentID: item.EnrollmentID.String(), ScheduleID: item.ScheduleID.String(), CheckDate: utils.TimeToStringPtrWithLayout(item.CheckDate, dateLayoutOnly), Status: string(item.Status), Note: item.Note, CreatedAt: item.CreatedAt.UTC().Format(time.RFC3339)}
}
