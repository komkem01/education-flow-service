package studentenrollments

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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
	SubjectAssignmentID string `json:"subject_assignment_id" binding:"required,uuid"`
	StudentNo           *int   `json:"student_no" binding:"omitempty,min=1"`
	Status              string `json:"status" binding:"omitempty,oneof=active dropped incomplete"`
}

type updateRequest struct {
	SubjectAssignmentID string `json:"subject_assignment_id" binding:"required,uuid"`
	StudentNo           *int   `json:"student_no" binding:"omitempty,min=1"`
	Status              string `json:"status" binding:"required,oneof=active dropped incomplete"`
}

type response struct {
	ID                  string `json:"id"`
	StudentID           string `json:"student_id"`
	SubjectAssignmentID string `json:"subject_assignment_id"`
	StudentNo           *int   `json:"student_no"`
	Status              string `json:"status"`
	CreatedAt           string `json:"created_at"`
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
	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	status := ent.StudentEnrollmentStatusActive
	if req.Status != "" {
		status = ent.ToStudentEnrollmentStatus(req.Status)
	}
	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, SubjectAssignmentID: subjectAssignmentID, StudentNo: req.StudentNo, Status: status})
	if err != nil {
		if errors.Is(err, ErrSubjectAssignmentCapacityExceeded) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentCapacityExceeded, nil)
			return
		}
		if isEnrollmentStudentNoDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentStudentNoDuplicate, nil)
			return
		}
		log.Errf("student-enrollments.create.error: %v", err)
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
		log.Errf("student-enrollments.list.error: %v", err)
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
	subjectAssignmentID, err := uuid.Parse(req.SubjectAssignmentID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{SubjectAssignmentID: subjectAssignmentID, StudentNo: req.StudentNo, Status: ent.ToStudentEnrollmentStatus(req.Status)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentNotFound, nil)
			return
		}
		if errors.Is(err, ErrSubjectAssignmentCapacityExceeded) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentCapacityExceeded, nil)
			return
		}
		if isEnrollmentStudentNoDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentStudentNoDuplicate, nil)
			return
		}
		log.Errf("student-enrollments.update.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.StudentEnrollmentNotFound, nil)
			return
		}
		log.Errf("student-enrollments.delete.error: %v", err)
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

func toResponse(item *ent.StudentEnrollment) response {
	return response{ID: item.ID.String(), StudentID: item.StudentID.String(), SubjectAssignmentID: item.SubjectAssignmentID.String(), StudentNo: item.StudentNo, Status: string(item.Status), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func isEnrollmentStudentNoDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	return pgErr.ConstraintName == "uq_student_enrollments_subject_assignment_student_no"
}
