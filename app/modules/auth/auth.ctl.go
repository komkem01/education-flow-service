package auth

import (
	"errors"
	"sort"

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

type loginRequest struct {
	Email    string `json:"email" binding:"required,max=255"`
	Password string `json:"password" binding:"required,max=255"`
}

type loginResponse struct {
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	ExpiresAt   string         `json:"expires_at"`
	Member      memberResponse `json:"member"`
}

type memberResponse struct {
	ID       string   `json:"id"`
	SchoolID string   `json:"school_id"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"is_active"`
}

type meResponse struct {
	MemberID  string   `json:"member_id"`
	SchoolID  string   `json:"school_id"`
	Role      string   `json:"role"`
	Roles     []string `json:"roles"`
	IssuedAt  string   `json:"issued_at"`
	ExpiresAt string   `json:"expires_at"`
}

type permissionsResponse struct {
	MemberID     string   `json:"member_id"`
	SchoolID     string   `json:"school_id"`
	Roles        []string `json:"roles"`
	Permissions  []string `json:"permissions"`
	BackOffice   bool     `json:"back_office"`
	PrimaryRole  string   `json:"primary_role"`
	TokenExpires string   `json:"token_expires"`
}

type switchRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type switchRoleResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresAt   string   `json:"expires_at"`
	Role        string   `json:"role"`
	Roles       []string `json:"roles"`
}

type switchSchoolRequest struct {
	SchoolID string `json:"school_id" binding:"required,uuid"`
}

type switchSchoolResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresAt   string   `json:"expires_at"`
	Role        string   `json:"role"`
	Roles       []string `json:"roles"`
	SchoolID    string   `json:"school_id"`
}

type refreshResponse = switchRoleResponse

func (c *Controller) Login(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)

	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	result, err := c.svc.Login(ctx.Request.Context(), &LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrInactiveMember) {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			return
		}
		log.Errf("auth.login.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, loginResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Member: memberResponse{
			ID:       result.Member.ID.String(),
			SchoolID: result.Member.SchoolID.String(),
			Email:    result.Member.Email,
			Role:     string(result.Member.Role),
			Roles:    toRoleStrings(result.Roles),
			IsActive: result.Member.IsActive,
		},
	})
}

func (c *Controller) Me(ctx *gin.Context) {
	claims, ok := GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	base.Success(ctx, meResponse{
		MemberID:  claims.MemberID.String(),
		SchoolID:  claims.SchoolID.String(),
		Role:      string(primaryRoleFromClaims(claims)),
		Roles:     toRoleStrings(claims.Roles),
		IssuedAt:  claims.IssuedAt.UTC().Format("2006-01-02T15:04:05Z"),
		ExpiresAt: claims.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

func (c *Controller) Permissions(ctx *gin.Context) {
	claims, ok := GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	permissions := permissionsFromRoles(claims.Roles)
	base.Success(ctx, permissionsResponse{
		MemberID:     claims.MemberID.String(),
		SchoolID:     claims.SchoolID.String(),
		Roles:        toRoleStrings(claims.Roles),
		Permissions:  permissions,
		BackOffice:   hasRole(claims.Roles, ent.MemberRoleAdmin) || hasRole(claims.Roles, ent.MemberRoleStaff) || hasRole(claims.Roles, ent.MemberRoleSuperAdmin),
		PrimaryRole:  string(primaryRoleFromClaims(claims)),
		TokenExpires: claims.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
	})
}

func (c *Controller) Logout(ctx *gin.Context) {
	// JWT auth is stateless in this service. Client-side token disposal completes logout.
	base.Success(ctx, gin.H{"success": true})
}

func (c *Controller) Refresh(ctx *gin.Context) {
	claims, ok := GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	result, err := c.svc.RefreshAccessToken(ctx.Request.Context(), &RefreshTokenInput{Claims: claims})
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken):
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		default:
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		}
		return
	}

	base.Success(ctx, refreshResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Role:        string(result.Role),
		Roles:       toRoleStrings(result.Roles),
	})
}

func (c *Controller) SwitchRole(ctx *gin.Context) {
	claims, ok := GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var req switchRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	role, parsed := parseKnownMemberRole(req.Role)
	if !parsed {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	result, err := c.svc.SwitchRole(ctx.Request.Context(), &SwitchRoleInput{
		Claims: claims,
		Role:   role,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrRoleNotAllowed):
			base.Forbidden(ctx, ci18n.Forbidden, nil)
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken):
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		default:
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		}
		return
	}

	base.Success(ctx, switchRoleResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Role:        string(result.Role),
		Roles:       toRoleStrings(result.Roles),
	})
}

func (c *Controller) SwitchSchool(ctx *gin.Context) {
	claims, ok := GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	var req switchSchoolRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	result, err := c.svc.SwitchSchool(ctx.Request.Context(), &SwitchSchoolInput{
		Claims:   claims,
		SchoolID: schoolID,
	})
	if err != nil {
		switch {
		case errors.Is(err, ErrRoleNotAllowed):
			base.Forbidden(ctx, ci18n.Forbidden, nil)
		case errors.Is(err, ErrSchoolNotFound):
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
		case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpiredToken):
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		default:
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		}
		return
	}

	base.Success(ctx, switchSchoolResponse{
		AccessToken: result.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   result.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Role:        string(result.Role),
		Roles:       toRoleStrings(result.Roles),
		SchoolID:    result.SchoolID.String(),
	})
}

func toRoleStrings(roles []ent.MemberRole) []string {
	out := make([]string, 0, len(roles))
	for _, role := range roles {
		out = append(out, string(role))
	}
	return out
}

func primaryRoleFromClaims(claims *TokenClaims) ent.MemberRole {
	if len(claims.Roles) > 0 {
		return claims.Roles[0]
	}

	return claims.Role
}

func permissionsFromRoles(roles []ent.MemberRole) []string {
	set := map[string]struct{}{
		"auth:me": {},
	}

	for _, role := range roles {
		switch role {
		case ent.MemberRoleAdmin:
			set["members:read"] = struct{}{}
			set["members:write"] = struct{}{}
			set["members:roles:write"] = struct{}{}
			set["admins:read"] = struct{}{}
			set["admins:write"] = struct{}{}
			set["backoffice:read"] = struct{}{}
		case ent.MemberRoleStaff:
			set["members:read"] = struct{}{}
			set["members:write"] = struct{}{}
			set["members:roles:write"] = struct{}{}
			set["admins:read"] = struct{}{}
			set["backoffice:read"] = struct{}{}
		case ent.MemberRoleTeacher:
			set["teacher:self:read"] = struct{}{}
			set["teacher:self:write"] = struct{}{}
		case ent.MemberRoleStudent:
			set["student:self:read"] = struct{}{}
		case ent.MemberRoleParent:
			set["parent:self:read"] = struct{}{}
		case ent.MemberRoleSuperAdmin:
			set["backoffice:read"] = struct{}{}
			set["backoffice:write"] = struct{}{}
			set["platform:schools:write"] = struct{}{}
			set["platform:admins:write"] = struct{}{}
			set["platform:scope:switch"] = struct{}{}
		}
	}

	result := make([]string, 0, len(set))
	for permission := range set {
		result = append(result, permission)
	}
	sort.Strings(result)

	return result
}

func hasRole(roles []ent.MemberRole, role ent.MemberRole) bool {
	for _, item := range roles {
		if item == role {
			return true
		}
	}

	return false
}
