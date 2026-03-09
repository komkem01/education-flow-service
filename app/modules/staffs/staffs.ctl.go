package staffs

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

const staffDateLayoutOnly = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createStaffRequest struct {
	MemberID   string  `json:"member_id" binding:"required,uuid"`
	GenderID   *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID   *string `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName  *string `json:"first_name" binding:"omitempty,max=255"`
	LastName   *string `json:"last_name" binding:"omitempty,max=255"`
	Phone      *string `json:"phone" binding:"omitempty,max=50"`
	Department *string `json:"department" binding:"omitempty,max=255"`
	IsActive   *bool   `json:"is_active"`
}

type updateStaffRequest = createStaffRequest

type registerStaffRequest struct {
	SchoolID        string                               `json:"school_id" binding:"required,uuid"`
	Email           string                               `json:"email" binding:"required,email,max=255"`
	Password        string                               `json:"password" binding:"required,min=6,max=255"`
	GenderID        *string                              `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID        *string                              `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName       *string                              `json:"first_name" binding:"omitempty,max=255"`
	LastName        *string                              `json:"last_name" binding:"omitempty,max=255"`
	Phone           *string                              `json:"phone" binding:"omitempty,max=50"`
	Department      *string                              `json:"department" binding:"omitempty,max=255"`
	IsActive        *bool                                `json:"is_active"`
	Educations      []registerStaffEducationRequest      `json:"educations" binding:"required,min=1,dive"`
	WorkExperiences []registerStaffWorkExperienceRequest `json:"work_experiences" binding:"required,min=1,dive"`
}

type registerStaffEducationRequest struct {
	DegreeLevel    *string `json:"degree_level" binding:"omitempty,max=100"`
	DegreeName     *string `json:"degree_name" binding:"omitempty,max=255"`
	Major          *string `json:"major" binding:"omitempty,max=255"`
	University     *string `json:"university" binding:"omitempty,max=255"`
	GraduationYear *string `json:"graduation_year" binding:"omitempty,max=10"`
}

type registerStaffWorkExperienceRequest struct {
	Organization *string `json:"organization" binding:"omitempty,max=255"`
	Position     *string `json:"position" binding:"omitempty,max=255"`
	StartDate    *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	IsCurrent    *bool   `json:"is_current"`
	Description  *string `json:"description"`
}

type staffResponse struct {
	ID         string  `json:"id"`
	MemberID   string  `json:"member_id"`
	GenderID   *string `json:"gender_id"`
	PrefixID   *string `json:"prefix_id"`
	StaffCode  *string `json:"staff_code"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Phone      *string `json:"phone"`
	Department *string `json:"department"`
	IsActive   bool    `json:"is_active"`
}

type registerStaffResponse struct {
	Member memberRegisterResponse `json:"member"`
	Staff  staffResponse          `json:"staff"`
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
	var req registerStaffRequest
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

	educations := make([]RegisterStaffEducationInput, 0, len(req.Educations))
	for _, item := range req.Educations {
		educations = append(educations, RegisterStaffEducationInput{
			DegreeLevel:    item.DegreeLevel,
			DegreeName:     item.DegreeName,
			Major:          item.Major,
			University:     item.University,
			GraduationYear: item.GraduationYear,
		})
	}

	workExperiences := make([]RegisterStaffWorkExperienceInput, 0, len(req.WorkExperiences))
	for _, item := range req.WorkExperiences {
		workStartDate, err := parseStaffDatePtr(item.StartDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		workEndDate, err := parseStaffDatePtr(item.EndDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		if workStartDate != nil && workEndDate != nil && workEndDate.Before(*workStartDate) {
			base.ValidateFailed(ctx, ci18n.StaffInvalidDateRange, nil)
			return
		}

		isCurrent := false
		if item.IsCurrent != nil {
			isCurrent = *item.IsCurrent
		}

		workExperiences = append(workExperiences, RegisterStaffWorkExperienceInput{
			Organization: item.Organization,
			Position:     item.Position,
			StartDate:    workStartDate,
			EndDate:      workEndDate,
			IsCurrent:    isCurrent,
			Description:  item.Description,
		})
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	member, staff, err := c.svc.Register(ctx.Request.Context(), &RegisterStaffInput{
		SchoolID:        schoolID,
		Email:           req.Email,
		Password:        req.Password,
		GenderID:        genderID,
		PrefixID:        prefixID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Department:      req.Department,
		IsActive:        isActive,
		Educations:      educations,
		WorkExperiences: workExperiences,
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
		if isForeignKeyConstraintError(err, "fk_member_staffs_gender_id") {
			base.ValidateFailed(ctx, ci18n.GenderNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_staffs_prefix_id") {
			base.ValidateFailed(ctx, ci18n.PrefixNotFound, nil)
			return
		}
		log.Errf("staffs.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, registerStaffResponse{
		Member: toMemberRegisterResponse(member),
		Staff:  toStaffResponse(staff),
	})
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createStaffRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseStaffCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	staff, err := c.svc.Create(ctx.Request.Context(), &CreateStaffInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, Department: req.Department, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidStaffMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("staffs.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toStaffResponse(staff))
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
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	staffs, err := c.svc.List(ctx.Request.Context(), &ListStaffsInput{SchoolID: &claims.SchoolID, MemberID: memberID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("staffs.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]staffResponse, 0, len(staffs))
	for _, staff := range staffs {
		response = append(response, toStaffResponse(staff))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	staffID, ok := parseStaffID(ctx)
	if !ok {
		return
	}

	staff, err := c.svc.GetByID(ctx.Request.Context(), staffID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StaffNotFound, nil)
			return
		}
		log.Errf("staffs.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toStaffResponse(staff))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	staffID, ok := parseStaffID(ctx)
	if !ok {
		return
	}

	var req updateStaffRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseStaffCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	staff, err := c.svc.UpdateByID(ctx.Request.Context(), staffID, &UpdateStaffInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, Department: req.Department, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidStaffMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StaffNotFound, nil)
			return
		}
		log.Errf("staffs.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toStaffResponse(staff))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	staffID, ok := parseStaffID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), staffID); err != nil {
		log.Errf("staffs.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": staffID.String()})
}

func parseStaffID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseStaffCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
	memberID, err := uuid.Parse(memberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	genderID, err := utils.ParseUUIDPtr(genderIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	prefixID, err := utils.ParseUUIDPtr(prefixIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}

	return memberID, genderID, prefixID, true
}

func toStaffResponse(staff *ent.MemberStaff) staffResponse {
	return staffResponse{ID: staff.ID.String(), MemberID: staff.MemberID.String(), GenderID: utils.UUIDToStringPtr(staff.GenderID), PrefixID: utils.UUIDToStringPtr(staff.PrefixID), StaffCode: staff.StaffCode, FirstName: staff.FirstName, LastName: staff.LastName, Phone: staff.Phone, Department: staff.Department, IsActive: staff.IsActive}
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

func parseStaffDatePtr(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	parsed, err := time.Parse(staffDateLayoutOnly, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
