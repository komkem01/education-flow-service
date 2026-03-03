package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiMemberAdmins(r *gin.RouterGroup, mod *modules.Modules) {
	registerCRUD(r, "/members", mod.Member.Ctl.List, mod.Member.Ctl.Get, mod.Member.Ctl.Create, mod.Member.Ctl.Update, mod.Member.Ctl.Delete)
	r.POST("/admins/register", mod.Admin.Ctl.Register)
	registerCRUD(r, "/admins", mod.Admin.Ctl.List, mod.Admin.Ctl.Get, mod.Admin.Ctl.Create, mod.Admin.Ctl.Update, mod.Admin.Ctl.Delete)
	r.POST("/parents/register", mod.Parent.Ctl.Register)
	registerCRUD(r, "/parents", mod.Parent.Ctl.List, mod.Parent.Ctl.Get, mod.Parent.Ctl.Create, mod.Parent.Ctl.Update, mod.Parent.Ctl.Delete)

	r.GET("/parents/:id/students", mod.ParentStudents.Ctl.List)
	r.POST("/parents/:id/students", mod.ParentStudents.Ctl.Create)
	r.PATCH("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Update)
	r.DELETE("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Delete)
}
