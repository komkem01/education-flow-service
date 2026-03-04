package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiAuth(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/auth/login", mod.Auth.Ctl.Login)

	authProtected := r.Group("/auth")
	authProtected.Use(requireAuth(mod))
	authProtected.GET("/me", mod.Auth.Ctl.Me)
	authProtected.GET("/permissions", mod.Auth.Ctl.Permissions)
	authProtected.POST("/switch-role", mod.Auth.Ctl.SwitchRole)
}
