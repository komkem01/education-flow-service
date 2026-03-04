package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiMemberAdmins(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/admins/register", mod.Admin.Ctl.Register)
	r.POST("/parents/register", mod.Parent.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	registerCRUD(protected, "/members", mod.Member.Ctl.List, mod.Member.Ctl.Get, mod.Member.Ctl.Create, mod.Member.Ctl.Update, mod.Member.Ctl.Delete)
	protected.GET("/members/:id/roles", mod.Member.Ctl.ListRoles)
	protected.POST("/members/:id/roles", mod.Member.Ctl.AddRole)
	protected.DELETE("/members/:id/roles/:role", mod.Member.Ctl.RemoveRole)
	registerCRUD(protected, "/admins", mod.Admin.Ctl.List, mod.Admin.Ctl.Get, mod.Admin.Ctl.Create, mod.Admin.Ctl.Update, mod.Admin.Ctl.Delete)

	protected.GET("/parents", mod.Parent.Ctl.List)
	protected.POST("/parents", mod.Parent.Ctl.Create)

	parentOwned := protected.Group("/parents/:id")
	parentOwned.Use(requireParentResourceOwnerOrRoles(mod, ent.MemberRoleAdmin, ent.MemberRoleStaff))
	parentOwned.GET("", mod.Parent.Ctl.Get)
	parentOwned.PATCH("", mod.Parent.Ctl.Update)
	parentOwned.DELETE("", mod.Parent.Ctl.Delete)

	parentOwned.GET("/students", mod.ParentStudents.Ctl.List)
	parentOwned.POST("/students", mod.ParentStudents.Ctl.Create)
	parentOwned.PATCH("/students/:child_id", mod.ParentStudents.Ctl.Update)
	parentOwned.DELETE("/students/:child_id", mod.ParentStudents.Ctl.Delete)
}
