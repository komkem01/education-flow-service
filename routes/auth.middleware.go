package routes

import (
	"errors"
	"strings"

	"education-flow/app/modules"
	"education-flow/app/modules/auth"
	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
)

func requireAuth(mod *modules.Modules) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, ok := extractBearerToken(ctx.GetHeader("Authorization"))
		if !ok {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			ctx.Abort()
			return
		}

		claims, err := mod.Auth.Svc.ParseAccessToken(token)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
				base.Unauthorized(ctx, ci18n.Unauthorized, nil)
				ctx.Abort()
				return
			}
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
			ctx.Abort()
			return
		}

		auth.SetClaimsToGin(ctx, claims)
		ctx.Next()
	}
}

func requireRoles(roles ...ent.MemberRole) gin.HandlerFunc {
	allow := make(map[ent.MemberRole]struct{}, len(roles))
	for _, role := range roles {
		allow[role] = struct{}{}
	}

	return func(ctx *gin.Context) {
		claims, ok := auth.GetClaimsFromGin(ctx)
		if !ok {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			ctx.Abort()
			return
		}

		if !claimsHasAnyRoleSet(claims, allow) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func claimsHasAnyRoleSet(claims *auth.TokenClaims, allow map[ent.MemberRole]struct{}) bool {
	for _, role := range claims.Roles {
		if _, ok := allow[role]; ok {
			return true
		}
	}

	return false
}

func extractBearerToken(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return "", false
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", false
	}

	return token, true
}
