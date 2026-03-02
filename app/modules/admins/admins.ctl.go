package admins

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
	"go.opentelemetry.io/otel/trace"
)

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

type adminResponse struct {
	ID        string  `json:"id"`
	MemberID  string  `json:"member_id"`
	GenderID  *string `json:"gender_id"`
	PrefixID  *string `json:"prefix_id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	IsActive  bool    `json:"is_active"`
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
	return adminResponse{ID: admin.ID.String(), MemberID: admin.MemberID.String(), GenderID: utils.UUIDToStringPtr(admin.GenderID), PrefixID: utils.UUIDToStringPtr(admin.PrefixID), FirstName: admin.FirstName, LastName: admin.LastName, Phone: admin.Phone, IsActive: admin.IsActive}
}
