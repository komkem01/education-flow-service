package staffs

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
	SchoolID   string  `json:"school_id" binding:"required,uuid"`
	Email      string  `json:"email" binding:"required,email,max=255"`
	Password   string  `json:"password" binding:"required,min=6,max=255"`
	GenderID   *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID   *string `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName  *string `json:"first_name" binding:"omitempty,max=255"`
	LastName   *string `json:"last_name" binding:"omitempty,max=255"`
	Phone      *string `json:"phone" binding:"omitempty,max=50"`
	Department *string `json:"department" binding:"omitempty,max=255"`
	IsActive   *bool   `json:"is_active"`
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

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	member, staff, err := c.svc.Register(ctx.Request.Context(), &RegisterStaffInput{
		SchoolID:   schoolID,
		Email:      req.Email,
		Password:   req.Password,
		GenderID:   genderID,
		PrefixID:   prefixID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Phone:      req.Phone,
		Department: req.Department,
		IsActive:   isActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
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

	staffs, err := c.svc.List(ctx.Request.Context(), &ListStaffsInput{MemberID: memberID, OnlyActive: onlyActive})
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
