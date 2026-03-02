package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiBackOffice(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/system-audit-logs", mod.SystemAuditLog.Ctl.List)
	r.GET("/system-audit-logs/:id", mod.SystemAuditLog.Ctl.Get)
	r.POST("/system-audit-logs", mod.SystemAuditLog.Ctl.Create)

	r.GET("/data-change-logs", mod.DataChangeLog.Ctl.List)
	r.GET("/data-change-logs/:id", mod.DataChangeLog.Ctl.Get)
	r.POST("/data-change-logs", mod.DataChangeLog.Ctl.Create)
}
