package members

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
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}

type memberURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createMemberRequest struct {
	SchoolID string `json:"school_id" binding:"required,uuid"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=255"`
	Role     string `json:"role" binding:"required,oneof=student teacher admin staff parent"`
	IsActive bool   `json:"is_active"`
}

type updateMemberRequest struct {
	SchoolID string `json:"school_id" binding:"required,uuid"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=255"`
	Role     string `json:"role" binding:"required,oneof=student teacher admin staff parent"`
	IsActive bool   `json:"is_active"`
}

type memberResponse struct {
	ID        string  `json:"id"`
	SchoolID  string  `json:"school_id"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	IsActive  bool    `json:"is_active"`
	LastLogin *string `json:"last_login"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	member, err := c.svc.Create(ctx.Request.Context(), &CreateMemberInput{
		SchoolID: schoolID,
		Email:    req.Email,
		Password: req.Password,
		Role:     ent.ToMemberRole(req.Role),
		IsActive: req.IsActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}

		log.Errf("members.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberResponse(member))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)

	var schoolID *uuid.UUID
	schoolIDQuery := ctx.Query("school_id")
	if schoolIDQuery != "" {
		parsedSchoolID, err := uuid.Parse(schoolIDQuery)
		if err != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return
		}
		schoolID = &parsedSchoolID
	}

	var role *ent.MemberRole
	roleQuery := ctx.Query("role")
	if roleQuery != "" {
		parsedRole := ent.ToMemberRole(roleQuery)
		if string(parsedRole) != roleQuery {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		role = &parsedRole
	}

	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	members, err := c.svc.List(ctx.Request.Context(), &ListMembersInput{
		SchoolID:   schoolID,
		Role:       role,
		OnlyActive: onlyActive,
	})
	if err != nil {
		log.Errf("members.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]memberResponse, 0, len(members))
	for _, member := range members {
		response = append(response, toMemberResponse(member))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	member, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}

		log.Errf("members.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberResponse(member))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	var req updateMemberRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	member, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateMemberInput{
		SchoolID: schoolID,
		Email:    req.Email,
		Password: req.Password,
		Role:     ent.ToMemberRole(req.Role),
		IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}

		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}

		log.Errf("members.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberResponse(member))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("members.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseMemberID(ctx *gin.Context) (uuid.UUID, bool) {
	var req memberURIRequest
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

func toMemberResponse(member *ent.Member) memberResponse {
	var lastLogin *string
	if member.LastLogin != nil {
		formatted := member.LastLogin.UTC().Format("2006-01-02T15:04:05Z")
		lastLogin = &formatted
	}

	return memberResponse{
		ID:        member.ID.String(),
		SchoolID:  member.SchoolID.String(),
		Email:     member.Email,
		Role:      string(member.Role),
		IsActive:  member.IsActive,
		LastLogin: lastLogin,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
