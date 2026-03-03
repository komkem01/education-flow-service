package routes

import (
	"database/sql"

	"education-flow/app/modules"
	"education-flow/app/modules/auth"
	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func requireTeacherResourceOwnerOrRoles(mod *modules.Modules, roles ...ent.MemberRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := auth.GetClaimsFromGin(ctx)
		if !ok {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			ctx.Abort()
			return
		}
		if claimsHasAnyRole(claims, roles...) {
			ctx.Next()
			return
		}
		if !claimsHasAnyRole(claims, ent.MemberRoleTeacher) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		teacherID, ok := parseIDParam(ctx)
		if !ok {
			return
		}

		teacher, err := mod.ENT.Svc.GetTeacherByID(ctx.Request.Context(), teacherID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.Next()
				return
			}
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
			ctx.Abort()
			return
		}

		if teacher.MemberID != claims.MemberID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func requireStudentResourceAccess(mod *modules.Modules, roles ...ent.MemberRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := auth.GetClaimsFromGin(ctx)
		if !ok {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			ctx.Abort()
			return
		}
		if claimsHasAnyRole(claims, roles...) {
			ctx.Next()
			return
		}

		studentID, ok := parseIDParam(ctx)
		if !ok {
			return
		}

		if claimsHasAnyRole(claims, ent.MemberRoleStudent) {
			student, err := mod.ENT.Svc.GetStudentByID(ctx.Request.Context(), studentID)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.Next()
					return
				}
				base.InternalServerError(ctx, ci18n.InternalServerError, nil)
				ctx.Abort()
				return
			}
			if student.MemberID != claims.MemberID {
				base.Forbidden(ctx, ci18n.Forbidden, nil)
				ctx.Abort()
				return
			}
			ctx.Next()
			return
		}

		if claimsHasAnyRole(claims, ent.MemberRoleParent) {
			parents, err := mod.ENT.Svc.ListParents(ctx.Request.Context(), &claims.MemberID, false)
			if err != nil {
				base.InternalServerError(ctx, ci18n.InternalServerError, nil)
				ctx.Abort()
				return
			}
			if len(parents) == 0 {
				base.Forbidden(ctx, ci18n.Forbidden, nil)
				ctx.Abort()
				return
			}

			links, err := mod.ENT.Svc.ListParentStudentsByParentID(ctx.Request.Context(), parents[0].ID)
			if err != nil {
				base.InternalServerError(ctx, ci18n.InternalServerError, nil)
				ctx.Abort()
				return
			}
			for _, item := range links {
				if item.StudentID == studentID {
					ctx.Next()
					return
				}
			}

			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		base.Forbidden(ctx, ci18n.Forbidden, nil)
		ctx.Abort()
		return
	}
}

func requireParentResourceOwnerOrRoles(mod *modules.Modules, roles ...ent.MemberRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, ok := auth.GetClaimsFromGin(ctx)
		if !ok {
			base.Unauthorized(ctx, ci18n.Unauthorized, nil)
			ctx.Abort()
			return
		}
		if claimsHasAnyRole(claims, roles...) {
			ctx.Next()
			return
		}
		if !claimsHasAnyRole(claims, ent.MemberRoleParent) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		parentID, ok := parseIDParam(ctx)
		if !ok {
			return
		}

		parent, err := mod.ENT.Svc.GetParentByID(ctx.Request.Context(), parentID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.Next()
				return
			}
			base.InternalServerError(ctx, ci18n.InternalServerError, nil)
			ctx.Abort()
			return
		}

		if parent.MemberID != claims.MemberID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func parseIDParam(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		ctx.Abort()
		return uuid.Nil, false
	}
	return id, true
}

func toRoleSet(roles ...ent.MemberRole) map[ent.MemberRole]struct{} {
	m := make(map[ent.MemberRole]struct{}, len(roles))
	for _, role := range roles {
		m[role] = struct{}{}
	}
	return m
}

func claimsHasAnyRole(claims *auth.TokenClaims, roles ...ent.MemberRole) bool {
	allow := toRoleSet(roles...)

	for _, role := range claims.Roles {
		if _, ok := allow[role]; ok {
			return true
		}
	}

	return false
}
