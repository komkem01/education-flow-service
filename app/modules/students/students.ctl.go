package students

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
	return &Controller{tracer: trace, svc: svc}
}

type createStudentRequest struct {
	MemberID           string  `json:"member_id" binding:"required,uuid"`
	GenderID           *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID           *string `json:"prefix_id" binding:"omitempty,uuid"`
	AdvisorTeacherID   *string `json:"advisor_teacher_id" binding:"omitempty,uuid"`
	CurrentClassroomID *string `json:"current_classroom_id" binding:"omitempty,uuid"`
	StudentCode        *string `json:"student_code" binding:"omitempty,max=255"`
	FirstName          *string `json:"first_name" binding:"omitempty,max=255"`
	LastName           *string `json:"last_name" binding:"omitempty,max=255"`
	CitizenID          *string `json:"citizen_id" binding:"omitempty,max=13"`
	Phone              *string `json:"phone" binding:"omitempty,max=50"`
	IsActive           *bool   `json:"is_active"`
}

type updateStudentRequest = createStudentRequest

type studentResponse struct {
	ID                 string  `json:"id"`
	MemberID           string  `json:"member_id"`
	GenderID           *string `json:"gender_id"`
	PrefixID           *string `json:"prefix_id"`
	AdvisorTeacherID   *string `json:"advisor_teacher_id"`
	CurrentClassroomID *string `json:"current_classroom_id"`
	StudentCode        *string `json:"student_code"`
	FirstName          *string `json:"first_name"`
	LastName           *string `json:"last_name"`
	CitizenID          *string `json:"citizen_id"`
	Phone              *string `json:"phone"`
	IsActive           bool    `json:"is_active"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, advisorTeacherID, currentClassroomID, ok := parseStudentCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.AdvisorTeacherID, req.CurrentClassroomID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	student, err := c.svc.Create(ctx.Request.Context(), &CreateStudentInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, StudentCode: req.StudentCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.StudentCodeDuplicate, nil)
			return
		}
		log.Errf("students.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toStudentResponse(student))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	memberID, err := utils.ParseQueryUUID(ctx.Query("member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	advisorTeacherID, err := utils.ParseQueryUUID(ctx.Query("advisor_teacher_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	currentClassroomID, err := utils.ParseQueryUUID(ctx.Query("current_classroom_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	students, err := c.svc.List(ctx.Request.Context(), &ListStudentsInput{MemberID: memberID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("students.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	response := make([]studentResponse, 0, len(students))
	for _, student := range students {
		response = append(response, toStudentResponse(student))
	}
	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, ok := parseStudentID(ctx)
	if !ok {
		return
	}
	student, err := c.svc.GetByID(ctx.Request.Context(), studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		log.Errf("students.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toStudentResponse(student))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, ok := parseStudentID(ctx)
	if !ok {
		return
	}
	var req updateStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, advisorTeacherID, currentClassroomID, ok := parseStudentCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.AdvisorTeacherID, req.CurrentClassroomID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	student, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, &UpdateStudentInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, StudentCode: req.StudentCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.StudentCodeDuplicate, nil)
			return
		}
		log.Errf("students.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toStudentResponse(student))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, ok := parseStudentID(ctx)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), studentID); err != nil {
		log.Errf("students.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": studentID.String()})
}

func parseStudentID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func parseStudentCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string, advisorTeacherIDRaw *string, currentClassroomIDRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
	memberID, err := uuid.Parse(memberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, nil, false
	}
	genderID, err := utils.ParseUUIDPtr(genderIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, nil, false
	}
	prefixID, err := utils.ParseUUIDPtr(prefixIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, nil, false
	}
	advisorTeacherID, err := utils.ParseUUIDPtr(advisorTeacherIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, nil, false
	}
	currentClassroomID, err := utils.ParseUUIDPtr(currentClassroomIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, nil, false
	}
	return memberID, genderID, prefixID, advisorTeacherID, currentClassroomID, true
}

func toStudentResponse(student *ent.MemberStudent) studentResponse {
	return studentResponse{ID: student.ID.String(), MemberID: student.MemberID.String(), GenderID: utils.UUIDToStringPtr(student.GenderID), PrefixID: utils.UUIDToStringPtr(student.PrefixID), AdvisorTeacherID: utils.UUIDToStringPtr(student.AdvisorTeacherID), CurrentClassroomID: utils.UUIDToStringPtr(student.CurrentClassroomID), StudentCode: student.StudentCode, FirstName: student.FirstName, LastName: student.LastName, CitizenID: student.CitizenID, Phone: student.Phone, IsActive: student.IsActive}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
