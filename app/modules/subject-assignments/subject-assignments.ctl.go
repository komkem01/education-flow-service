package subjectassignments

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

type createSubjectAssignmentRequest struct {
	SubjectID      string  `json:"subject_id" binding:"required,uuid"`
	TeacherID      string  `json:"teacher_id" binding:"required,uuid"`
	ClassroomID    string  `json:"classroom_id" binding:"required,uuid"`
	AcademicYearID string  `json:"academic_year_id" binding:"required,uuid"`
	Section        *string `json:"section" binding:"omitempty,max=50"`
	SemesterNo     *int    `json:"semester_no" binding:"omitempty,gte=1,lte=3"`
	MaxStudents    *int    `json:"max_students" binding:"omitempty,gte=0"`
	StartDate      *string `json:"start_date" binding:"omitempty,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	Note           *string `json:"note" binding:"omitempty,max=4000"`
	IsActive       *bool   `json:"is_active"`
}

type updateSubjectAssignmentRequest = createSubjectAssignmentRequest

type createTeacherSubjectAssignmentRequest struct {
	SubjectID      string  `json:"subject_id" binding:"required,uuid"`
	ClassroomID    string  `json:"classroom_id" binding:"required,uuid"`
	AcademicYearID string  `json:"academic_year_id" binding:"required,uuid"`
	Section        *string `json:"section" binding:"omitempty,max=50"`
	SemesterNo     *int    `json:"semester_no" binding:"omitempty,gte=1,lte=3"`
	MaxStudents    *int    `json:"max_students" binding:"omitempty,gte=0"`
	StartDate      *string `json:"start_date" binding:"omitempty,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	Note           *string `json:"note" binding:"omitempty,max=4000"`
	IsActive       *bool   `json:"is_active"`
}

type updateTeacherSubjectAssignmentRequest = createTeacherSubjectAssignmentRequest

type subjectAssignmentResponse struct {
	ID             string  `json:"id"`
	SubjectID      string  `json:"subject_id"`
	TeacherID      string  `json:"teacher_id"`
	ClassroomID    string  `json:"classroom_id"`
	AcademicYearID string  `json:"academic_year_id"`
	Section        *string `json:"section"`
	SemesterNo     int     `json:"semester_no"`
	MaxStudents    *int    `json:"max_students"`
	StartDate      *string `json:"start_date"`
	EndDate        *string `json:"end_date"`
	Note           *string `json:"note"`
	IsActive       bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createSubjectAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, teacherID, classroomID, academicYearID, ok := parseSubjectAssignmentCreateUpdateFields(ctx, req.SubjectID, req.TeacherID, req.ClassroomID, req.AcademicYearID)
	if !ok {
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate, ci18n.SubjectAssignmentInvalidDateRange)
	if !ok {
		return
	}

	semesterNo := 1
	if req.SemesterNo != nil {
		semesterNo = *req.SemesterNo
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectAssignmentInput{SubjectID: subjectID, TeacherID: teacherID, ClassroomID: classroomID, AcademicYearID: academicYearID, Section: req.Section, SemesterNo: semesterNo, MaxStudents: req.MaxStudents, StartDate: startDate, EndDate: endDate, Note: req.Note, IsActive: isActive})
	if err != nil {
		if isSubjectAssignmentDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentDuplicate, nil)
			return
		}
		if isSubjectAssignmentDateRangeError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentInvalidDateRange, nil)
			return
		}
		log.Errf("subject-assignments.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectAssignmentResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	subjectID, err := utils.ParseQueryUUID(ctx.Query("subject_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	teacherID, err := utils.ParseQueryUUID(ctx.Query("teacher_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	classroomID, err := utils.ParseQueryUUID(ctx.Query("classroom_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	academicYearID, err := utils.ParseQueryUUID(ctx.Query("academic_year_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var semesterNo *int
	if raw := strings.TrimSpace(ctx.Query("semester_no")); raw != "" {
		value, convErr := strconv.Atoi(raw)
		if convErr != nil || value < 1 || value > 3 {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		semesterNo = &value
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

	items, err := c.svc.List(ctx.Request.Context(), &ListSubjectAssignmentsInput{SubjectID: subjectID, TeacherID: teacherID, ClassroomID: classroomID, AcademicYearID: academicYearID, SemesterNo: semesterNo, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("subject-assignments.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]subjectAssignmentResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toSubjectAssignmentResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectAssignmentID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
			return
		}
		log.Errf("subject-assignments.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectAssignmentResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectAssignmentID(ctx)
	if !ok {
		return
	}

	var req updateSubjectAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, teacherID, classroomID, academicYearID, ok := parseSubjectAssignmentCreateUpdateFields(ctx, req.SubjectID, req.TeacherID, req.ClassroomID, req.AcademicYearID)
	if !ok {
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate, ci18n.SubjectAssignmentInvalidDateRange)
	if !ok {
		return
	}

	semesterNo := 1
	if req.SemesterNo != nil {
		semesterNo = *req.SemesterNo
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateSubjectAssignmentInput{SubjectID: subjectID, TeacherID: teacherID, ClassroomID: classroomID, AcademicYearID: academicYearID, Section: req.Section, SemesterNo: semesterNo, MaxStudents: req.MaxStudents, StartDate: startDate, EndDate: endDate, Note: req.Note, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
			return
		}
		if isSubjectAssignmentDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentDuplicate, nil)
			return
		}
		if isSubjectAssignmentDateRangeError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentInvalidDateRange, nil)
			return
		}
		log.Errf("subject-assignments.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectAssignmentResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseSubjectAssignmentID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("subject-assignments.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func (c *Controller) ListByTeacher(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherIDFromPath(ctx)
	if !ok {
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListSubjectAssignmentsInput{TeacherID: &teacherID})
	if err != nil {
		log.Errf("subject-assignments.list-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]subjectAssignmentResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toSubjectAssignmentResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) CreateByTeacher(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherIDFromPath(ctx)
	if !ok {
		return
	}

	var req createTeacherSubjectAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, classroomID, academicYearID, ok := parseTeacherSubjectAssignmentCreateUpdateFields(ctx, req.SubjectID, req.ClassroomID, req.AcademicYearID)
	if !ok {
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate, ci18n.SubjectAssignmentInvalidDateRange)
	if !ok {
		return
	}

	semesterNo := 1
	if req.SemesterNo != nil {
		semesterNo = *req.SemesterNo
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateSubjectAssignmentInput{SubjectID: subjectID, TeacherID: teacherID, ClassroomID: classroomID, AcademicYearID: academicYearID, Section: req.Section, SemesterNo: semesterNo, MaxStudents: req.MaxStudents, StartDate: startDate, EndDate: endDate, Note: req.Note, IsActive: isActive})
	if err != nil {
		if isSubjectAssignmentDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentDuplicate, nil)
			return
		}
		if isSubjectAssignmentDateRangeError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentInvalidDateRange, nil)
			return
		}
		log.Errf("subject-assignments.create-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectAssignmentResponse(item))
}

func (c *Controller) UpdateByTeacher(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherIDFromPath(ctx)
	if !ok {
		return
	}

	assignmentID, ok := parseSubjectAssignmentChildID(ctx)
	if !ok {
		return
	}

	existing, err := c.svc.GetByID(ctx.Request.Context(), assignmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
			return
		}
		log.Errf("subject-assignments.get-for-update-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if existing.TeacherID != teacherID {
		base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
		return
	}

	var req updateTeacherSubjectAssignmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	subjectID, classroomID, academicYearID, ok := parseTeacherSubjectAssignmentCreateUpdateFields(ctx, req.SubjectID, req.ClassroomID, req.AcademicYearID)
	if !ok {
		return
	}

	startDate, endDate, ok := parseDateRange(ctx, req.StartDate, req.EndDate, ci18n.SubjectAssignmentInvalidDateRange)
	if !ok {
		return
	}

	semesterNo := 1
	if req.SemesterNo != nil {
		semesterNo = *req.SemesterNo
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), assignmentID, &UpdateSubjectAssignmentInput{SubjectID: subjectID, TeacherID: teacherID, ClassroomID: classroomID, AcademicYearID: academicYearID, Section: req.Section, SemesterNo: semesterNo, MaxStudents: req.MaxStudents, StartDate: startDate, EndDate: endDate, Note: req.Note, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
			return
		}
		if isSubjectAssignmentDuplicateError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentDuplicate, nil)
			return
		}
		if isSubjectAssignmentDateRangeError(err) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentInvalidDateRange, nil)
			return
		}
		log.Errf("subject-assignments.update-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toSubjectAssignmentResponse(item))
}

func (c *Controller) DeleteByTeacher(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherIDFromPath(ctx)
	if !ok {
		return
	}

	assignmentID, ok := parseSubjectAssignmentChildID(ctx)
	if !ok {
		return
	}

	existing, err := c.svc.GetByID(ctx.Request.Context(), assignmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
			return
		}
		log.Errf("subject-assignments.get-for-delete-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if existing.TeacherID != teacherID {
		base.ValidateFailed(ctx, ci18n.SubjectAssignmentNotFound, nil)
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), assignmentID); err != nil {
		log.Errf("subject-assignments.delete-by-teacher.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": assignmentID.String()})
}

func parseSubjectAssignmentID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseTeacherIDFromPath(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseSubjectAssignmentChildID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "child_id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseSubjectAssignmentCreateUpdateFields(ctx *gin.Context, subjectIDRaw string, teacherIDRaw string, classroomIDRaw string, academicYearIDRaw string) (uuid.UUID, uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	subjectID, err := uuid.Parse(subjectIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	teacherID, err := uuid.Parse(teacherIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	classroomID, err := uuid.Parse(classroomIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	academicYearID, err := uuid.Parse(academicYearIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, uuid.Nil, false
	}

	return subjectID, teacherID, classroomID, academicYearID, true
}

func parseTeacherSubjectAssignmentCreateUpdateFields(ctx *gin.Context, subjectIDRaw string, classroomIDRaw string, academicYearIDRaw string) (uuid.UUID, uuid.UUID, uuid.UUID, bool) {
	subjectID, err := uuid.Parse(subjectIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	classroomID, err := uuid.Parse(classroomIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}
	academicYearID, err := uuid.Parse(academicYearIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, uuid.Nil, false
	}

	return subjectID, classroomID, academicYearID, true
}

func toSubjectAssignmentResponse(item *ent.SubjectAssignment) subjectAssignmentResponse {
	return subjectAssignmentResponse{
		ID:             item.ID.String(),
		SubjectID:      item.SubjectID.String(),
		TeacherID:      item.TeacherID.String(),
		ClassroomID:    item.ClassroomID.String(),
		AcademicYearID: item.AcademicYearID.String(),
		Section:        item.Section,
		SemesterNo:     item.SemesterNo,
		MaxStudents:    item.MaxStudents,
		StartDate:      toDateStringPtr(item.StartDate),
		EndDate:        toDateStringPtr(item.EndDate),
		Note:           item.Note,
		IsActive:       item.IsActive,
	}
}

func parseDatePtr(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	value, err := time.Parse("2006-01-02", *raw)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func parseDateRange(ctx *gin.Context, startDateRaw *string, endDateRaw *string, invalidMessage string) (*time.Time, *time.Time, bool) {
	startDate, err := parseDatePtr(startDateRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return nil, nil, false
	}

	endDate, err := parseDatePtr(endDateRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return nil, nil, false
	}

	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		base.ValidateFailed(ctx, invalidMessage, nil)
		return nil, nil, false
	}

	return startDate, endDate, true
}

func toDateStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format("2006-01-02")
	return &formatted
}

func isSubjectAssignmentDuplicateError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "uq_subject_assignments_unique_slot")
}

func isSubjectAssignmentDateRangeError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23514" {
		return false
	}

	constraint := strings.ToLower(pgErr.ConstraintName)
	return strings.Contains(constraint, "chk_subject_assignments_date_range")
}
