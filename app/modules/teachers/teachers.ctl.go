package teachers

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

const dateLayoutOnly = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type teacherURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createTeacherRequest struct {
	MemberID                string                  `json:"member_id" binding:"required,uuid"`
	GenderID                *string                 `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID                *string                 `json:"prefix_id" binding:"omitempty,uuid"`
	TeacherCode             *string                 `json:"teacher_code" binding:"omitempty,max=255"`
	FirstName               *string                 `json:"first_name" binding:"omitempty,max=255"`
	LastName                *string                 `json:"last_name" binding:"omitempty,max=255"`
	CitizenID               *string                 `json:"citizen_id" binding:"omitempty,max=13"`
	Phone                   *string                 `json:"phone" binding:"omitempty,max=50"`
	CurrentPosition         *string                 `json:"current_position" binding:"omitempty,max=255"`
	CurrentAcademicStanding *string                 `json:"current_academic_standing" binding:"omitempty,max=255"`
	Department              *string                 `json:"department" binding:"omitempty,max=255"`
	StartDate               *string                 `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	Addresses               *[]memberAddressRequest `json:"addresses" binding:"omitempty,dive"`
	IsActive                *bool                   `json:"is_active"`
}

type updateTeacherRequest = createTeacherRequest

type registerTeacherRequest struct {
	SchoolID                string                                 `json:"school_id" binding:"required,uuid"`
	Email                   string                                 `json:"email" binding:"required,email,max=255"`
	Password                string                                 `json:"password" binding:"required,min=6,max=255"`
	GenderID                *string                                `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID                *string                                `json:"prefix_id" binding:"omitempty,uuid"`
	TeacherCode             *string                                `json:"teacher_code" binding:"omitempty,max=255"`
	FirstName               *string                                `json:"first_name" binding:"omitempty,max=255"`
	LastName                *string                                `json:"last_name" binding:"omitempty,max=255"`
	CitizenID               *string                                `json:"citizen_id" binding:"omitempty,max=13"`
	Phone                   *string                                `json:"phone" binding:"omitempty,max=50"`
	CurrentPosition         *string                                `json:"current_position" binding:"omitempty,max=255"`
	CurrentAcademicStanding *string                                `json:"current_academic_standing" binding:"omitempty,max=255"`
	Department              *string                                `json:"department" binding:"omitempty,max=255"`
	StartDate               *string                                `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	Addresses               *[]memberAddressRequest                `json:"addresses" binding:"omitempty,dive"`
	IsActive                *bool                                  `json:"is_active"`
	Educations              []registerTeacherEducationRequest      `json:"educations" binding:"required,min=1,dive"`
	WorkExperiences         []registerTeacherWorkExperienceRequest `json:"work_experiences" binding:"required,min=1,dive"`
}

type memberAddressRequest struct {
	Label       *string `json:"label" binding:"omitempty,max=100"`
	AddressLine string  `json:"address_line" binding:"required,max=1000"`
	IsPrimary   *bool   `json:"is_primary"`
	SortOrder   *int    `json:"sort_order" binding:"omitempty,min=0,max=1000"`
}

type memberAddressResponse struct {
	ID          string  `json:"id"`
	MemberID    string  `json:"member_id"`
	Label       *string `json:"label"`
	AddressLine string  `json:"address_line"`
	IsPrimary   bool    `json:"is_primary"`
	SortOrder   int     `json:"sort_order"`
}

type registerTeacherEducationRequest struct {
	DegreeLevel    *string `json:"degree_level" binding:"omitempty,max=100"`
	DegreeName     *string `json:"degree_name" binding:"omitempty,max=255"`
	Major          *string `json:"major" binding:"omitempty,max=255"`
	University     *string `json:"university" binding:"omitempty,max=255"`
	GraduationYear *string `json:"graduation_year" binding:"omitempty,max=10"`
}

type registerTeacherWorkExperienceRequest struct {
	Organization *string `json:"organization" binding:"omitempty,max=255"`
	Position     *string `json:"position" binding:"omitempty,max=255"`
	StartDate    *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	IsCurrent    *bool   `json:"is_current"`
	Description  *string `json:"description"`
}

type teacherResponse struct {
	ID                      string                  `json:"id"`
	MemberID                string                  `json:"member_id"`
	GenderID                *string                 `json:"gender_id"`
	PrefixID                *string                 `json:"prefix_id"`
	TeacherCode             *string                 `json:"teacher_code"`
	FirstName               *string                 `json:"first_name"`
	LastName                *string                 `json:"last_name"`
	CitizenID               *string                 `json:"citizen_id"`
	Phone                   *string                 `json:"phone"`
	CurrentPosition         *string                 `json:"current_position"`
	CurrentAcademicStanding *string                 `json:"current_academic_standing"`
	Department              *string                 `json:"department"`
	StartDate               *string                 `json:"start_date"`
	Addresses               []memberAddressResponse `json:"addresses"`
	IsActive                bool                    `json:"is_active"`
}

type registerTeacherResponse struct {
	Member  memberRegisterResponse `json:"member"`
	Teacher teacherResponse        `json:"teacher"`
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
	var req registerTeacherRequest
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

	genderID, err := parseUUIDPtr(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	prefixID, err := parseUUIDPtr(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	startDate, err := parseDatePtr(req.StartDate)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	educations := make([]RegisterTeacherEducationInput, 0, len(req.Educations))
	for _, item := range req.Educations {
		educations = append(educations, RegisterTeacherEducationInput{
			DegreeLevel:    item.DegreeLevel,
			DegreeName:     item.DegreeName,
			Major:          item.Major,
			University:     item.University,
			GraduationYear: item.GraduationYear,
		})
	}

	workExperiences := make([]RegisterTeacherWorkExperienceInput, 0, len(req.WorkExperiences))
	for _, item := range req.WorkExperiences {
		workStartDate, err := parseDatePtr(item.StartDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		workEndDate, err := parseDatePtr(item.EndDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		if workStartDate != nil && workEndDate != nil && workEndDate.Before(*workStartDate) {
			base.ValidateFailed(ctx, ci18n.TeacherInvalidDateRange, nil)
			return
		}

		isCurrent := false
		if item.IsCurrent != nil {
			isCurrent = *item.IsCurrent
		}

		workExperiences = append(workExperiences, RegisterTeacherWorkExperienceInput{
			Organization: item.Organization,
			Position:     item.Position,
			StartDate:    workStartDate,
			EndDate:      workEndDate,
			IsCurrent:    isCurrent,
			Description:  item.Description,
		})
	}
	addresses := parseMemberAddressInputs(req.Addresses)

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	member, teacher, err := c.svc.Register(ctx.Request.Context(), &RegisterTeacherInput{
		SchoolID:                schoolID,
		Email:                   req.Email,
		Password:                req.Password,
		GenderID:                genderID,
		PrefixID:                prefixID,
		TeacherCode:             req.TeacherCode,
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		CitizenID:               req.CitizenID,
		Phone:                   req.Phone,
		CurrentPosition:         req.CurrentPosition,
		CurrentAcademicStanding: req.CurrentAcademicStanding,
		Department:              req.Department,
		StartDate:               startDate,
		Addresses:               addresses,
		IsActive:                isActive,
		Educations:              educations,
		WorkExperiences:         workExperiences,
	})
	if err != nil {
		if isTeacherCodeDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.TeacherCodeDuplicate, nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_members_school_id") {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_teachers_gender_id") {
			base.ValidateFailed(ctx, ci18n.GenderNotFound, nil)
			return
		}
		if isForeignKeyConstraintError(err, "fk_member_teachers_prefix_id") {
			base.ValidateFailed(ctx, ci18n.PrefixNotFound, nil)
			return
		}
		log.Errf("teachers.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	teacherAddresses, err := c.svc.ListAddressesByMemberID(ctx.Request.Context(), teacher.MemberID)
	if err != nil {
		log.Errf("teachers.register.addresses.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, registerTeacherResponse{
		Member:  toMemberRegisterResponse(member),
		Teacher: toTeacherResponse(teacher, teacherAddresses),
	})
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var req createTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, startDate, ok := parseTeacherCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.StartDate)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	addresses := parseMemberAddressInputs(req.Addresses)
	teacher, err := c.svc.Create(ctx.Request.Context(), &CreateTeacherInput{SchoolID: claims.SchoolID, MemberID: memberID, GenderID: genderID, PrefixID: prefixID, TeacherCode: req.TeacherCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, CurrentPosition: req.CurrentPosition, CurrentAcademicStanding: req.CurrentAcademicStanding, Department: req.Department, StartDate: startDate, Addresses: addresses, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidTeacherMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrTeacherSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.TeacherCodeDuplicate, nil)
			return
		}
		log.Errf("teachers.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	teacherAddresses, err := c.svc.ListAddressesByMemberID(ctx.Request.Context(), teacher.MemberID)
	if err != nil {
		log.Errf("teachers.create.addresses.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher, teacherAddresses))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var memberID *uuid.UUID
	if memberIDRaw := ctx.Query("member_id"); memberIDRaw != "" {
		parsedMemberID, err := uuid.Parse(memberIDRaw)
		if err != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return
		}
		memberID = &parsedMemberID
	}
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	teachers, err := c.svc.List(ctx.Request.Context(), &ListTeachersInput{SchoolID: claims.SchoolID, MemberID: memberID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("teachers.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	response := make([]teacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		addresses, addrErr := c.svc.ListAddressesByMemberID(ctx.Request.Context(), teacher.MemberID)
		if addrErr != nil {
			log.Errf("teachers.list.addresses.error: %v", addrErr)
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
			return
		}
		response = append(response, toTeacherResponse(teacher, addresses))
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

	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	teacher, err := c.svc.GetByIDInSchool(ctx.Request.Context(), claims.SchoolID, teacherID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherNotFound, nil)
			return
		}
		if errors.Is(err, ErrTeacherSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		log.Errf("teachers.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	addresses, err := c.svc.ListAddressesByMemberID(ctx.Request.Context(), teacher.MemberID)
	if err != nil {
		log.Errf("teachers.get.addresses.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher, addresses))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	var req updateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, startDate, ok := parseTeacherCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.StartDate)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	addresses := parseMemberAddressInputs(req.Addresses)
	teacher, err := c.svc.UpdateByID(ctx.Request.Context(), teacherID, &UpdateTeacherInput{SchoolID: claims.SchoolID, MemberID: memberID, GenderID: genderID, PrefixID: prefixID, TeacherCode: req.TeacherCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, CurrentPosition: req.CurrentPosition, CurrentAcademicStanding: req.CurrentAcademicStanding, Department: req.Department, StartDate: startDate, Addresses: addresses, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidTeacherMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherNotFound, nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.TeacherCodeDuplicate, nil)
			return
		}
		if errors.Is(err, ErrTeacherSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		log.Errf("teachers.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	teacherAddresses, err := c.svc.ListAddressesByMemberID(ctx.Request.Context(), teacher.MemberID)
	if err != nil {
		log.Errf("teachers.update.addresses.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher, teacherAddresses))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	if err := c.svc.DeleteByIDInSchool(ctx.Request.Context(), claims.SchoolID, teacherID); err != nil {
		if errors.Is(err, ErrTeacherSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		log.Errf("teachers.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": teacherID.String()})
}

func parseTeacherID(ctx *gin.Context) (uuid.UUID, bool) {
	var req teacherURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, false
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func parseTeacherCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string, startDateRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, *time.Time, bool) {
	memberID, err := uuid.Parse(memberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	genderID, err := parseUUIDPtr(genderIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	prefixID, err := parseUUIDPtr(prefixIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	startDate, err := parseDatePtr(startDateRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	return memberID, genderID, prefixID, startDate, true
}

func parseUUIDPtr(raw *string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func parseDatePtr(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateLayoutOnly, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func toTeacherResponse(teacher *ent.MemberTeacher, addresses []*ent.MemberAddress) teacherResponse {
	addressResponse := make([]memberAddressResponse, 0, len(addresses))
	for _, item := range addresses {
		addressResponse = append(addressResponse, memberAddressResponse{
			ID:          item.ID.String(),
			MemberID:    item.MemberID.String(),
			Label:       item.Label,
			AddressLine: item.AddressLine,
			IsPrimary:   item.IsPrimary,
			SortOrder:   item.SortOrder,
		})
	}

	return teacherResponse{ID: teacher.ID.String(), MemberID: teacher.MemberID.String(), GenderID: uuidToStringPtr(teacher.GenderID), PrefixID: uuidToStringPtr(teacher.PrefixID), TeacherCode: teacher.TeacherCode, FirstName: teacher.FirstName, LastName: teacher.LastName, CitizenID: teacher.CitizenID, Phone: teacher.Phone, CurrentPosition: teacher.CurrentPosition, CurrentAcademicStanding: teacher.CurrentAcademicStanding, Department: teacher.Department, StartDate: dateToStringPtr(teacher.StartDate), Addresses: addressResponse, IsActive: teacher.IsActive}
}

func parseMemberAddressInputs(raw *[]memberAddressRequest) []MemberAddressInput {
	if raw == nil {
		return []MemberAddressInput{}
	}

	items := make([]MemberAddressInput, 0, len(*raw))
	for i, item := range *raw {
		isPrimary := false
		if item.IsPrimary != nil {
			isPrimary = *item.IsPrimary
		}
		sortOrder := i
		if item.SortOrder != nil {
			sortOrder = *item.SortOrder
		}

		items = append(items, MemberAddressInput{
			Label:       item.Label,
			AddressLine: item.AddressLine,
			IsPrimary:   isPrimary,
			SortOrder:   sortOrder,
		})
	}

	return items
}

func uuidToStringPtr(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	parsed := value.String()
	return &parsed
}

func dateToStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	parsed := value.Format(dateLayoutOnly)
	return &parsed
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

func isTeacherCodeDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Code != "23505" {
		return false
	}

	constraint := pgErr.ConstraintName
	return constraint == "uq_member_teachers_teacher_code" || constraint == "member_teachers_teacher_code_key"
}
