package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiMemberTeachers(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/teachers/register", mod.Teacher.Ctl.Register)
	registerCRUD(r, "/teachers", mod.Teacher.Ctl.List, mod.Teacher.Ctl.Get, mod.Teacher.Ctl.Create, mod.Teacher.Ctl.Update, mod.Teacher.Ctl.Delete)

	r.GET("/teachers/:id/educations", mod.TeacherEducations.Ctl.List)
	r.POST("/teachers/:id/educations", mod.TeacherEducations.Ctl.Create)
	r.PATCH("/teachers/:id/educations/:child_id", mod.TeacherEducations.Ctl.Update)
	r.DELETE("/teachers/:id/educations/:child_id", mod.TeacherEducations.Ctl.Delete)

	r.GET("/teachers/:id/profile-requests", mod.TeacherProfileRequests.Ctl.List)
	r.POST("/teachers/:id/profile-requests", mod.TeacherProfileRequests.Ctl.Create)
	r.PATCH("/teachers/:id/profile-requests/:child_id", mod.TeacherProfileRequests.Ctl.Update)

	r.GET("/teachers/:id/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.List)
	r.POST("/teachers/:id/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.Create)
	r.PATCH("/teachers/:id/performance-agreements/:child_id", mod.TeacherPerformanceAgreements.Ctl.Update)

	r.GET("/teachers/:id/pda-logs", mod.TeacherPDALogs.Ctl.List)
	r.POST("/teachers/:id/pda-logs", mod.TeacherPDALogs.Ctl.Create)
	r.DELETE("/teachers/:id/pda-logs/:child_id", mod.TeacherPDALogs.Ctl.Delete)

	r.GET("/teachers/:id/leave-logs", mod.TeacherLeaveLogs.Ctl.List)
	r.POST("/teachers/:id/leave-logs", mod.TeacherLeaveLogs.Ctl.Create)
	r.PATCH("/teachers/:id/leave-logs/:child_id", mod.TeacherLeaveLogs.Ctl.Update)
}
