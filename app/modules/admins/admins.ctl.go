package admins

import (
	"database/sql"
	"errors"
	"strconv"
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

const adminDateLayoutOnly = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createAdminRequest struct {
	MemberID  string  `json:"member_id" binding:"required,uuid"`
	GenderID  *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID  *string `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName *string `json:"first_name" binding:"omitempty,max=255"`
	LastName  *string `json:"last_name" binding:"omitempty,max=255"`
	Phone     *string `json:"phone" binding:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active"`
}

type updateAdminRequest = createAdminRequest

type registerAdminRequest struct {
	SchoolID        string                               `json:"school_id" binding:"required,uuid"`
	Email           string                               `json:"email" binding:"required,email,max=255"`
	Password        string                               `json:"password" binding:"required,min=6,max=255"`
	GenderID        *string                              `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID        *string                              `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName       *string                              `json:"first_name" binding:"omitempty,max=255"`
	LastName        *string                              `json:"last_name" binding:"omitempty,max=255"`
	Phone           *string                              `json:"phone" binding:"omitempty,max=50"`
	IsActive        *bool                                `json:"is_active"`
	Educations      []registerAdminEducationRequest      `json:"educations" binding:"required,min=1,dive"`
	WorkExperiences []registerAdminWorkExperienceRequest `json:"work_experiences" binding:"required,min=1,dive"`
}

type registerAdminEducationRequest struct {
	DegreeLevel    *string `json:"degree_level" binding:"omitempty,max=100"`
	DegreeName     *string `json:"degree_name" binding:"omitempty,max=255"`
	Major          *string `json:"major" binding:"omitempty,max=255"`
	University     *string `json:"university" binding:"omitempty,max=255"`
	GraduationYear *string `json:"graduation_year" binding:"omitempty,max=10"`
}

type registerAdminWorkExperienceRequest struct {
	Organization *string `json:"organization" binding:"omitempty,max=255"`
	Position     *string `json:"position" binding:"omitempty,max=255"`
	StartDate    *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	IsCurrent    *bool   `json:"is_current"`
	Description  *string `json:"description"`
}

type adminResponse struct {
	ID        string  `json:"id"`
	MemberID  string  `json:"member_id"`
	GenderID  *string `json:"gender_id"`
	PrefixID  *string `json:"prefix_id"`
	AdminCode *string `json:"admin_code"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	IsActive  bool    `json:"is_active"`
}

type registerAdminResponse struct {
	Member memberRegisterResponse `json:"member"`
	Admin  adminResponse          `json:"admin"`
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
	var req registerAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
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

	educations := make([]RegisterAdminEducationInput, 0, len(req.Educations))
	for _, item := range req.Educations {
		educations = append(educations, RegisterAdminEducationInput{
			DegreeLevel:    item.DegreeLevel,
			DegreeName:     item.DegreeName,
			Major:          item.Major,
			University:     item.University,
			GraduationYear: item.GraduationYear,
		})
	}

	workExperiences := make([]RegisterAdminWorkExperienceInput, 0, len(req.WorkExperiences))
	for _, item := range req.WorkExperiences {
		workStartDate, err := parseAdminDatePtr(item.StartDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		workEndDate, err := parseAdminDatePtr(item.EndDate)
		if err != nil {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		if workStartDate != nil && workEndDate != nil && workEndDate.Before(*workStartDate) {
			base.ValidateFailed(ctx, ci18n.AdminInvalidDateRange, nil)
			return
		}

		isCurrent := false
		if item.IsCurrent != nil {
			isCurrent = *item.IsCurrent
		}

		workExperiences = append(workExperiences, RegisterAdminWorkExperienceInput{
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

	member, admin, err := c.svc.Register(ctx.Request.Context(), &RegisterAdminInput{
		SchoolID:        schoolID,
		Email:           req.Email,
		Password:        req.Password,
		GenderID:        genderID,
		PrefixID:        prefixID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		IsActive:        isActive,
		Educations:      educations,
		WorkExperiences: workExperiences,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}
		log.Errf("admins.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, registerAdminResponse{
		Member: toMemberRegisterResponse(member),
		Admin:  toAdminResponse(admin),
	})
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseAdminCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	admin, err := c.svc.Create(ctx.Request.Context(), &CreateAdminInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidAdminMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("admins.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAdminResponse(admin))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
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

	admins, err := c.svc.List(ctx.Request.Context(), &ListAdminsInput{MemberID: memberID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("admins.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]adminResponse, 0, len(admins))
	for _, admin := range admins {
		response = append(response, toAdminResponse(admin))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	adminID, ok := parseAdminID(ctx)
	if !ok {
		return
	}

	admin, err := c.svc.GetByID(ctx.Request.Context(), adminID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AdminNotFound, nil)
			return
		}
		log.Errf("admins.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAdminResponse(admin))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	adminID, ok := parseAdminID(ctx)
	if !ok {
		return
	}

	var req updateAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseAdminCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	admin, err := c.svc.UpdateByID(ctx.Request.Context(), adminID, &UpdateAdminInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidAdminMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.AdminNotFound, nil)
			return
		}
		log.Errf("admins.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toAdminResponse(admin))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	adminID, ok := parseAdminID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), adminID); err != nil {
		log.Errf("admins.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": adminID.String()})
}

func parseAdminID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseAdminCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
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

func toAdminResponse(admin *ent.MemberAdmin) adminResponse {
	return adminResponse{ID: admin.ID.String(), MemberID: admin.MemberID.String(), GenderID: utils.UUIDToStringPtr(admin.GenderID), PrefixID: utils.UUIDToStringPtr(admin.PrefixID), AdminCode: admin.AdminCode, FirstName: admin.FirstName, LastName: admin.LastName, Phone: admin.Phone, IsActive: admin.IsActive}
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

func parseAdminDatePtr(raw *string) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}
	parsed, err := time.Parse(adminDateLayoutOnly, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
