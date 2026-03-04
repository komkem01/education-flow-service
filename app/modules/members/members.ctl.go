package members

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
	return &Controller{
		tracer: trace,
		svc:    svc,
	}
}

type memberURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type memberRoleURIRequest struct {
	ID   string `uri:"id" binding:"required"`
	Role string `uri:"role" binding:"required"`
}

type addMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=student teacher admin staff parent"`
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

type memberRolesResponse struct {
	MemberID string   `json:"member_id"`
	Roles    []string `json:"roles"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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
	if schoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var schoolID *uuid.UUID
	schoolIDQuery := ctx.Query("school_id")
	if schoolIDQuery != "" {
		parsedSchoolID, err := uuid.Parse(schoolIDQuery)
		if err != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return
		}
		if parsedSchoolID != claims.SchoolID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
	}
	schoolID = &claims.SchoolID

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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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
	if member.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	base.Success(ctx, toMemberResponse(member))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	existing, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}

		log.Errf("members.update.precheck.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if existing.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
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
	if schoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
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
		if errors.Is(err, ErrStudentRoleExclusive) {
			base.ValidateFailed(ctx, ci18n.BadRequest, nil)
			return
		}

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
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

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

		log.Errf("members.delete.precheck.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if member.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("members.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func (c *Controller) ListRoles(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	roles, err := c.svc.ListRoles(ctx.Request.Context(), claims.SchoolID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}

		log.Errf("members.list-roles.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberRolesResponse(id, roles))
}

func (c *Controller) AddRole(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	id, ok := parseMemberID(ctx)
	if !ok {
		return
	}

	var req addMemberRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	roles, err := c.svc.AddRole(ctx.Request.Context(), claims.SchoolID, id, ent.ToMemberRole(req.Role))
	if err != nil {
		if errors.Is(err, ErrStudentRoleExclusive) {
			base.ValidateFailed(ctx, ci18n.BadRequest, nil)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}

		log.Errf("members.add-role.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberRolesResponse(id, roles))
}

func (c *Controller) RemoveRole(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	memberID, role, ok := parseMemberRoleIDAndRole(ctx)
	if !ok {
		return
	}

	roles, err := c.svc.RemoveRole(ctx.Request.Context(), claims.SchoolID, memberID, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberSchoolMismatch) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		if errors.Is(err, ErrMemberRoleRequired) {
			base.ValidateFailed(ctx, ci18n.BadRequest, nil)
			return
		}

		log.Errf("members.remove-role.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toMemberRolesResponse(memberID, roles))
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

func parseMemberRoleIDAndRole(ctx *gin.Context) (uuid.UUID, ent.MemberRole, bool) {
	var req memberRoleURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, "", false
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, "", false
	}

	role := ent.ToMemberRole(req.Role)
	if string(role) != req.Role {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, "", false
	}

	return id, role, true
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

func toMemberRolesResponse(memberID uuid.UUID, roles []ent.MemberRole) memberRolesResponse {
	responseRoles := make([]string, 0, len(roles))
	for _, role := range roles {
		responseRoles = append(responseRoles, string(role))
	}

	return memberRolesResponse{
		MemberID: memberID.String(),
		Roles:    responseRoles,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
