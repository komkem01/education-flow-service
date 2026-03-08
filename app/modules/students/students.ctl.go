package students

import (
	"database/sql"
	"errors"
	"strconv"

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
	DefaultStudentNo   *int    `json:"default_student_no" binding:"omitempty,min=1"`
	FirstName          *string `json:"first_name" binding:"omitempty,max=255"`
	LastName           *string `json:"last_name" binding:"omitempty,max=255"`
	CitizenID          *string `json:"citizen_id" binding:"omitempty,max=13"`
	Phone              *string `json:"phone" binding:"omitempty,max=50"`
	IsActive           *bool   `json:"is_active"`
}

type updateStudentRequest = createStudentRequest

type registerStudentRequest struct {
	SchoolID           string                        `json:"school_id" binding:"required,uuid"`
	Email              string                        `json:"email" binding:"required,email,max=255"`
	Password           string                        `json:"password" binding:"required,min=6,max=255"`
	GenderID           *string                       `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID           *string                       `json:"prefix_id" binding:"omitempty,uuid"`
	AdvisorTeacherID   *string                       `json:"advisor_teacher_id" binding:"omitempty,uuid"`
	CurrentClassroomID *string                       `json:"current_classroom_id" binding:"omitempty,uuid"`
	StudentCode        *string                       `json:"student_code" binding:"omitempty,max=255"`
	DefaultStudentNo   *int                          `json:"default_student_no" binding:"omitempty,min=1"`
	FirstName          *string                       `json:"first_name" binding:"omitempty,max=255"`
	LastName           *string                       `json:"last_name" binding:"omitempty,max=255"`
	CitizenID          *string                       `json:"citizen_id" binding:"omitempty,max=13"`
	Phone              *string                       `json:"phone" binding:"omitempty,max=50"`
	IsActive           *bool                         `json:"is_active"`
	Parent             *registerStudentParentRequest `json:"parent" binding:"omitempty"`
}

type registerStudentParentRequest struct {
	Email          string  `json:"email" binding:"required,email,max=255"`
	Password       string  `json:"password" binding:"required,min=6,max=255"`
	GenderID       *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID       *string `json:"prefix_id" binding:"omitempty,uuid"`
	ParentCode     *string `json:"parent_code" binding:"omitempty,max=255"`
	FirstName      *string `json:"first_name" binding:"omitempty,max=255"`
	LastName       *string `json:"last_name" binding:"omitempty,max=255"`
	Phone          *string `json:"phone" binding:"omitempty,max=50"`
	Relationship   string  `json:"relationship" binding:"omitempty,oneof=father mother guardian"`
	IsMainGuardian *bool   `json:"is_main_guardian"`
	IsActive       *bool   `json:"is_active"`
}

type studentResponse struct {
	ID                 string  `json:"id"`
	MemberID           string  `json:"member_id"`
	GenderID           *string `json:"gender_id"`
	PrefixID           *string `json:"prefix_id"`
	AdvisorTeacherID   *string `json:"advisor_teacher_id"`
	CurrentClassroomID *string `json:"current_classroom_id"`
	StudentCode        *string `json:"student_code"`
	DefaultStudentNo   *int    `json:"default_student_no"`
	FirstName          *string `json:"first_name"`
	LastName           *string `json:"last_name"`
	CitizenID          *string `json:"citizen_id"`
	Phone              *string `json:"phone"`
	IsActive           bool    `json:"is_active"`
}

type registerStudentResponse struct {
	Member  memberRegisterResponse  `json:"member"`
	Student studentResponse         `json:"student"`
	Parent  *registerParentResponse `json:"parent,omitempty"`
}

type registerParentResponse struct {
	Member         memberRegisterResponse `json:"member"`
	Parent         parentResponseLite     `json:"parent"`
	Relationship   string                 `json:"relationship"`
	IsMainGuardian bool                   `json:"is_main_guardian"`
}

type parentResponseLite struct {
	ID         string  `json:"id"`
	MemberID   string  `json:"member_id"`
	GenderID   *string `json:"gender_id"`
	PrefixID   *string `json:"prefix_id"`
	ParentCode *string `json:"parent_code"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Phone      *string `json:"phone"`
	IsActive   bool    `json:"is_active"`
}

type memberRegisterResponse struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) Register(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req registerStudentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	if claims, ok := auth.GetClaimsFromGin(ctx); ok && claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	genderID, err := utils.ParseUUIDPtr(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	prefixID, err := utils.ParseUUIDPtr(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	advisorTeacherID, err := utils.ParseUUIDPtr(req.AdvisorTeacherID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	currentClassroomID, err := utils.ParseUUIDPtr(req.CurrentClassroomID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	parentInput, ok := parseRegisterParentInput(ctx, req.Parent)
	if !ok {
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	result, err := c.svc.Register(ctx.Request.Context(), &RegisterStudentInput{
		SchoolID:           schoolID,
		Email:              req.Email,
		Password:           req.Password,
		GenderID:           genderID,
		PrefixID:           prefixID,
		AdvisorTeacherID:   advisorTeacherID,
		CurrentClassroomID: currentClassroomID,
		StudentCode:        req.StudentCode,
		DefaultStudentNo:   req.DefaultStudentNo,
		FirstName:          req.FirstName,
		LastName:           req.LastName,
		CitizenID:          req.CitizenID,
		Phone:              req.Phone,
		IsActive:           isActive,
		Parent:             parentInput,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_members_school_id") {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_students_gender_id") || isForeignKeyConstraintError(err, "fk_member_parents_gender_id") {
			base.ValidateFailed(ctx, ci18n.GenderNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_students_prefix_id") || isForeignKeyConstraintError(err, "fk_member_parents_prefix_id") {
			base.ValidateFailed(ctx, ci18n.PrefixNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_students_advisor_teacher_id") {
			base.ValidateFailed(ctx, ci18n.TeacherNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_students_current_classroom_id") {
			base.ValidateFailed(ctx, ci18n.ClassroomNotFound, nil)
			return
		}
		log.Errf("students.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := registerStudentResponse{
		Member:  toMemberRegisterResponse(result.StudentMember),
		Student: toStudentResponse(result.Student),
	}
	if result.ParentMember != nil && result.Parent != nil && result.ParentStudent != nil {
		response.Parent = &registerParentResponse{
			Member:         toMemberRegisterResponse(result.ParentMember),
			Parent:         toParentResponseLite(result.Parent),
			Relationship:   string(result.ParentStudent.Relationship),
			IsMainGuardian: result.ParentStudent.IsMainGuardian,
		}
	}

	base.Success(ctx, response)
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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
	student, err := c.svc.Create(ctx.Request.Context(), &CreateStudentInput{SchoolID: claims.SchoolID, MemberID: memberID, GenderID: genderID, PrefixID: prefixID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, StudentCode: req.StudentCode, DefaultStudentNo: req.DefaultStudentNo, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrStudentSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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
	students, err := c.svc.List(ctx.Request.Context(), &ListStudentsInput{SchoolID: claims.SchoolID, MemberID: memberID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, OnlyActive: onlyActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.Success(ctx, []studentResponse{})
			return
		}
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, ok := parseStudentID(ctx)
	if !ok {
		return
	}
	student, err := c.svc.GetByIDInSchool(ctx.Request.Context(), claims.SchoolID, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		if errors.Is(err, ErrStudentSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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
	student, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, &UpdateStudentInput{SchoolID: claims.SchoolID, MemberID: memberID, GenderID: genderID, PrefixID: prefixID, AdvisorTeacherID: advisorTeacherID, CurrentClassroomID: currentClassroomID, StudentCode: req.StudentCode, DefaultStudentNo: req.DefaultStudentNo, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentNotFound, nil)
			return
		}
		if errors.Is(err, ErrStudentSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	studentID, ok := parseStudentID(ctx)
	if !ok {
		return
	}
	if err := c.svc.DeleteByIDInSchool(ctx.Request.Context(), claims.SchoolID, studentID); err != nil {
		if errors.Is(err, ErrStudentSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
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
	return studentResponse{ID: student.ID.String(), MemberID: student.MemberID.String(), GenderID: utils.UUIDToStringPtr(student.GenderID), PrefixID: utils.UUIDToStringPtr(student.PrefixID), AdvisorTeacherID: utils.UUIDToStringPtr(student.AdvisorTeacherID), CurrentClassroomID: utils.UUIDToStringPtr(student.CurrentClassroomID), StudentCode: student.StudentCode, DefaultStudentNo: student.DefaultStudentNo, FirstName: student.FirstName, LastName: student.LastName, CitizenID: student.CitizenID, Phone: student.Phone, IsActive: student.IsActive}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func isForeignKeyConstraintError(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23503" && pgErr.ConstraintName == constraint
}

func toMemberRegisterResponse(member *ent.Member) memberRegisterResponse {
	return memberRegisterResponse{
		ID:       member.ID.String(),
		SchoolID: member.SchoolID.String(),
		Email:    member.Email,
		Role:     string(member.Role),
		IsActive: member.IsActive,
	}
}

func toParentResponseLite(parent *ent.MemberParent) parentResponseLite {
	return parentResponseLite{
		ID:         parent.ID.String(),
		MemberID:   parent.MemberID.String(),
		GenderID:   utils.UUIDToStringPtr(parent.GenderID),
		PrefixID:   utils.UUIDToStringPtr(parent.PrefixID),
		ParentCode: parent.ParentCode,
		FirstName:  parent.FirstName,
		LastName:   parent.LastName,
		Phone:      parent.Phone,
		IsActive:   parent.IsActive,
	}
}

func parseRegisterParentInput(ctx *gin.Context, req *registerStudentParentRequest) (*RegisterParentInput, bool) {
	if req == nil {
		return nil, true
	}

	genderID, err := utils.ParseUUIDPtr(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return nil, false
	}
	prefixID, err := utils.ParseUUIDPtr(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return nil, false
	}

	relationship := ent.ParentRelationshipGuardian
	if req.Relationship != "" {
		relationship = ent.ToParentRelationship(req.Relationship)
	}

	isMainGuardian := true
	if req.IsMainGuardian != nil {
		isMainGuardian = *req.IsMainGuardian
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &RegisterParentInput{
		Email:          req.Email,
		Password:       req.Password,
		GenderID:       genderID,
		PrefixID:       prefixID,
		ParentCode:     req.ParentCode,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		Relationship:   relationship,
		IsMainGuardian: isMainGuardian,
		IsActive:       isActive,
	}, true
}
