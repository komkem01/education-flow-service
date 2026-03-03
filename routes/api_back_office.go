package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiBackOffice(r *gin.RouterGroup, mod *modules.Modules) {
	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin))

	protected.GET("/system-audit-logs", mod.SystemAuditLog.Ctl.List)
	protected.GET("/system-audit-logs/:id", mod.SystemAuditLog.Ctl.Get)
	protected.POST("/system-audit-logs", mod.SystemAuditLog.Ctl.Create)

	protected.GET("/data-change-logs", mod.DataChangeLog.Ctl.List)
	protected.GET("/data-change-logs/:id", mod.DataChangeLog.Ctl.Get)
	protected.POST("/data-change-logs", mod.DataChangeLog.Ctl.Create)
}
