package routes

import (
	"fmt"
	"net/http"

	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

func api(r *gin.RouterGroup, mod *modules.Modules) {
	r.GET("/example/:id", mod.Example.Ctl.Get)
	r.GET("/example-http", mod.Example.Ctl.GetHttpReq)
	r.POST("/example", mod.Example.Ctl.Create)

	r.GET("/schools", mod.School.Ctl.List)
	r.GET("/schools/:id", mod.School.Ctl.Get)
	r.POST("/schools", mod.School.Ctl.Create)
	r.PATCH("/schools/:id", mod.School.Ctl.Update)
	r.DELETE("/schools/:id", mod.School.Ctl.Delete)

	r.GET("/genders", mod.Gender.Ctl.List)
	r.GET("/genders/:id", mod.Gender.Ctl.Get)
	r.POST("/genders", mod.Gender.Ctl.Create)
	r.PATCH("/genders/:id", mod.Gender.Ctl.Update)
	r.DELETE("/genders/:id", mod.Gender.Ctl.Delete)

	r.GET("/prefixes", mod.Prefix.Ctl.List)
	r.GET("/prefixes/:id", mod.Prefix.Ctl.Get)
	r.POST("/prefixes", mod.Prefix.Ctl.Create)
	r.PATCH("/prefixes/:id", mod.Prefix.Ctl.Update)
	r.DELETE("/prefixes/:id", mod.Prefix.Ctl.Delete)

	r.GET("/members", mod.Member.Ctl.List)
	r.GET("/members/:id", mod.Member.Ctl.Get)
	r.POST("/members", mod.Member.Ctl.Create)
	r.PATCH("/members/:id", mod.Member.Ctl.Update)
	r.DELETE("/members/:id", mod.Member.Ctl.Delete)
}
