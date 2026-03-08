package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiMemberStudents(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/students/register", mod.Student.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher, ent.MemberRoleStudent, ent.MemberRoleParent))

	protected.GET("/students", mod.Student.Ctl.List)
	protected.POST("/students", mod.Student.Ctl.Create)

	studentOwned := protected.Group("/students/:id")
	studentOwned.Use(requireStudentResourceAccess(mod, ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher))
	studentOwned.GET("", mod.Student.Ctl.Get)
	studentOwned.PATCH("", mod.Student.Ctl.Update)
	studentOwned.DELETE("", mod.Student.Ctl.Delete)

	studentOwned.GET("/enrollments", mod.StudentEnrollments.Ctl.List)
	studentOwned.POST("/enrollments", mod.StudentEnrollments.Ctl.Create)
	studentOwned.PATCH("/enrollments/:child_id", mod.StudentEnrollments.Ctl.Update)
	studentOwned.DELETE("/enrollments/:child_id", mod.StudentEnrollments.Ctl.Delete)

	studentOwned.GET("/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.List)
	studentOwned.POST("/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.Create)
	studentOwned.PATCH("/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Update)
	studentOwned.DELETE("/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Delete)

	studentOwned.GET("/invoices", mod.StudentInvoices.Ctl.List)
	studentOwned.POST("/invoices", mod.StudentInvoices.Ctl.Create)
	studentOwned.PATCH("/invoices/:child_id", mod.StudentInvoices.Ctl.Update)
	studentOwned.DELETE("/invoices/:child_id", mod.StudentInvoices.Ctl.Delete)

	studentOwned.GET("/attendance-logs", mod.StudentAttendanceLogs.Ctl.List)
	studentOwned.POST("/attendance-logs", mod.StudentAttendanceLogs.Ctl.Create)
	studentOwned.PATCH("/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Update)
	studentOwned.DELETE("/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Delete)

	studentOwned.GET("/payment-transactions", mod.PaymentTransactions.Ctl.List)
	studentOwned.POST("/payment-transactions", mod.PaymentTransactions.Ctl.Create)
	studentOwned.PATCH("/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Update)
	studentOwned.DELETE("/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Delete)

	studentOwned.GET("/fee-categories", mod.StudentFeeCategories.Ctl.List)
	studentOwned.POST("/fee-categories", mod.StudentFeeCategories.Ctl.Create)
	studentOwned.PATCH("/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Update)
	studentOwned.DELETE("/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Delete)

	studentOwned.GET("/grade-items", mod.StudentGradeItems.Ctl.List)
	studentOwned.POST("/grade-items", mod.StudentGradeItems.Ctl.Create)
	studentOwned.PATCH("/grade-items/:child_id", mod.StudentGradeItems.Ctl.Update)
	studentOwned.DELETE("/grade-items/:child_id", mod.StudentGradeItems.Ctl.Delete)

	studentOwned.GET("/grade-records", mod.StudentGradeRecords.Ctl.List)
	studentOwned.POST("/grade-records", mod.StudentGradeRecords.Ctl.Create)
	studentOwned.PATCH("/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Update)
	studentOwned.DELETE("/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Delete)

	studentOwned.GET("/announcements", mod.SchoolAnnouncement.Ctl.List)
}
