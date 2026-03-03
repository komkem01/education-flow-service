package auth

import (
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
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
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=255"`
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
