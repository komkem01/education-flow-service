package parents

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

type createParentRequest struct {
	MemberID  string  `json:"member_id" binding:"required,uuid"`
	GenderID  *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID  *string `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName *string `json:"first_name" binding:"omitempty,max=255"`
	LastName  *string `json:"last_name" binding:"omitempty,max=255"`
	Phone     *string `json:"phone" binding:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active"`
}

type updateParentRequest = createParentRequest

type registerParentRequest struct {
	SchoolID  string  `json:"school_id" binding:"required,uuid"`
	Email     string  `json:"email" binding:"required,email,max=255"`
	Password  string  `json:"password" binding:"required,min=6,max=255"`
	GenderID  *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID  *string `json:"prefix_id" binding:"omitempty,uuid"`
	FirstName *string `json:"first_name" binding:"omitempty,max=255"`
	LastName  *string `json:"last_name" binding:"omitempty,max=255"`
	Phone     *string `json:"phone" binding:"omitempty,max=50"`
	IsActive  *bool   `json:"is_active"`
}

type parentResponse struct {
	ID        string  `json:"id"`
	MemberID  string  `json:"member_id"`
	GenderID  *string `json:"gender_id"`
	PrefixID  *string `json:"prefix_id"`
	ParentCode *string `json:"parent_code"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	IsActive  bool    `json:"is_active"`
}

type registerParentResponse struct {
	Member memberRegisterResponse `json:"member"`
	Parent parentResponse         `json:"parent"`
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
	var req registerParentRequest
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

	member, parent, err := c.svc.Register(ctx.Request.Context(), &RegisterParentInput{
		SchoolID:  schoolID,
		Email:     req.Email,
		Password:  req.Password,
		GenderID:  genderID,
		PrefixID:  prefixID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  isActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}
		log.Errf("parents.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, registerParentResponse{
		Member: toMemberRegisterResponse(member),
		Parent: toParentResponse(parent),
	})
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createParentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseParentCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	parent, err := c.svc.Create(ctx.Request.Context(), &CreateParentInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidParentMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("parents.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toParentResponse(parent))
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

	parents, err := c.svc.List(ctx.Request.Context(), &ListParentsInput{MemberID: memberID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("parents.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]parentResponse, 0, len(parents))
	for _, parent := range parents {
		response = append(response, toParentResponse(parent))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, ok := parseParentID(ctx)
	if !ok {
		return
	}

	parent, err := c.svc.GetByID(ctx.Request.Context(), parentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ParentNotFound, nil)
			return
		}
		log.Errf("parents.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toParentResponse(parent))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, ok := parseParentID(ctx)
	if !ok {
		return
	}

	var req updateParentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	memberID, genderID, prefixID, ok := parseParentCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	parent, err := c.svc.UpdateByID(ctx.Request.Context(), parentID, &UpdateParentInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, FirstName: req.FirstName, LastName: req.LastName, Phone: req.Phone, IsActive: isActive})
	if err != nil {
		if errors.Is(err, ErrInvalidParentMemberRole) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.ParentNotFound, nil)
			return
		}
		log.Errf("parents.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toParentResponse(parent))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	parentID, ok := parseParentID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), parentID); err != nil {
		log.Errf("parents.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": parentID.String()})
}

func parseParentID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseParentCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
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

func toParentResponse(parent *ent.MemberParent) parentResponse {
	return parentResponse{ID: parent.ID.String(), MemberID: parent.MemberID.String(), GenderID: utils.UUIDToStringPtr(parent.GenderID), PrefixID: utils.UUIDToStringPtr(parent.PrefixID), ParentCode: parent.ParentCode, FirstName: parent.FirstName, LastName: parent.LastName, Phone: parent.Phone, IsActive: parent.IsActive}
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
