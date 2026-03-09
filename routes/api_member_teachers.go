package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiMemberTeachers(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/teachers/register", mod.Teacher.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher))

	protected.GET("/teachers", mod.Teacher.Ctl.List)
	protected.POST("/teachers", mod.Teacher.Ctl.Create)

	// Read-only metadata endpoints used by teacher-facing UI.
	protected.GET("/teachers-meta/schools/:id", mod.School.Ctl.Get)
	protected.GET("/teachers-meta/academic-years/:id", mod.AcademicYear.Ctl.Get)
	protected.GET("/teachers-meta/classrooms/:id", mod.Classroom.Ctl.Get)
	protected.GET("/teachers-meta/subjects/:id", mod.Subject.Ctl.Get)
	protected.GET("/teachers-meta/schedules", mod.Schedule.Ctl.List)
	protected.GET("/teachers-meta/school-announcements", mod.SchoolAnnouncement.Ctl.List)

	teacherOwned := protected.Group("/teachers/:id")
	teacherOwned.Use(requireTeacherResourceOwnerOrRoles(mod, ent.MemberRoleAdmin, ent.MemberRoleStaff))
	teacherOwned.GET("", mod.Teacher.Ctl.Get)
	teacherOwned.PATCH("", mod.Teacher.Ctl.Update)
	teacherOwned.DELETE("", mod.Teacher.Ctl.Delete)

	teacherOwned.GET("/educations", mod.TeacherEducations.Ctl.List)
	teacherOwned.POST("/educations", mod.TeacherEducations.Ctl.Create)
	teacherOwned.PATCH("/educations/:child_id", mod.TeacherEducations.Ctl.Update)
	teacherOwned.DELETE("/educations/:child_id", mod.TeacherEducations.Ctl.Delete)

	teacherOwned.GET("/work-experiences", mod.TeacherWorkExperiences.Ctl.List)
	teacherOwned.POST("/work-experiences", mod.TeacherWorkExperiences.Ctl.Create)
	teacherOwned.PATCH("/work-experiences/:child_id", mod.TeacherWorkExperiences.Ctl.Update)
	teacherOwned.DELETE("/work-experiences/:child_id", mod.TeacherWorkExperiences.Ctl.Delete)

	teacherOwned.GET("/profile-requests", mod.TeacherProfileRequests.Ctl.List)
	teacherOwned.POST("/profile-requests", mod.TeacherProfileRequests.Ctl.Create)
	teacherOwned.PATCH("/profile-requests/:child_id", mod.TeacherProfileRequests.Ctl.Update)

	teacherOwned.GET("/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.List)
	teacherOwned.POST("/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.Create)
	teacherOwned.PATCH("/performance-agreements/:child_id", mod.TeacherPerformanceAgreements.Ctl.Update)

	teacherOwned.GET("/pda-logs", mod.TeacherPDALogs.Ctl.List)
	teacherOwned.POST("/pda-logs", mod.TeacherPDALogs.Ctl.Create)
	teacherOwned.DELETE("/pda-logs/:child_id", mod.TeacherPDALogs.Ctl.Delete)

	teacherOwned.GET("/leave-logs", mod.TeacherLeaveLogs.Ctl.List)
	teacherOwned.POST("/leave-logs", mod.TeacherLeaveLogs.Ctl.Create)
	teacherOwned.PATCH("/leave-logs/:child_id", mod.TeacherLeaveLogs.Ctl.Update)

	teacherOwned.GET("/subject-assignments", mod.SubjectAssignment.Ctl.ListByTeacher)
	teacherOwned.POST("/subject-assignments", mod.SubjectAssignment.Ctl.CreateByTeacher)
	teacherOwned.PATCH("/subject-assignments/:child_id", mod.SubjectAssignment.Ctl.UpdateByTeacher)
	teacherOwned.DELETE("/subject-assignments/:child_id", mod.SubjectAssignment.Ctl.DeleteByTeacher)
}
