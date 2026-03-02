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

	r.GET("/academic-years", mod.AcademicYear.Ctl.List)
	r.GET("/academic-years/:id", mod.AcademicYear.Ctl.Get)
	r.POST("/academic-years", mod.AcademicYear.Ctl.Create)
	r.PATCH("/academic-years/:id", mod.AcademicYear.Ctl.Update)
	r.DELETE("/academic-years/:id", mod.AcademicYear.Ctl.Delete)

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

	r.GET("/teachers", mod.Teacher.Ctl.List)
	r.GET("/teachers/:id", mod.Teacher.Ctl.Get)
	r.POST("/teachers", mod.Teacher.Ctl.Create)
	r.PATCH("/teachers/:id", mod.Teacher.Ctl.Update)
	r.DELETE("/teachers/:id", mod.Teacher.Ctl.Delete)

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
