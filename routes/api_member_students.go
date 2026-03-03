package routes

import (
	"education-flow/app/modules"

	"github.com/gin-gonic/gin"
)

func apiMemberStudents(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/students/register", mod.Student.Ctl.Register)
	registerCRUD(r, "/students", mod.Student.Ctl.List, mod.Student.Ctl.Get, mod.Student.Ctl.Create, mod.Student.Ctl.Update, mod.Student.Ctl.Delete)

	r.GET("/students/:id/enrollments", mod.StudentEnrollments.Ctl.List)
	r.POST("/students/:id/enrollments", mod.StudentEnrollments.Ctl.Create)
	r.PATCH("/students/:id/enrollments/:child_id", mod.StudentEnrollments.Ctl.Update)
	r.DELETE("/students/:id/enrollments/:child_id", mod.StudentEnrollments.Ctl.Delete)

	r.GET("/students/:id/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.List)
	r.POST("/students/:id/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.Create)
	r.PATCH("/students/:id/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Update)
	r.DELETE("/students/:id/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Delete)

	r.GET("/students/:id/invoices", mod.StudentInvoices.Ctl.List)
	r.POST("/students/:id/invoices", mod.StudentInvoices.Ctl.Create)
	r.PATCH("/students/:id/invoices/:child_id", mod.StudentInvoices.Ctl.Update)
	r.DELETE("/students/:id/invoices/:child_id", mod.StudentInvoices.Ctl.Delete)

	r.GET("/students/:id/attendance-logs", mod.StudentAttendanceLogs.Ctl.List)
	r.POST("/students/:id/attendance-logs", mod.StudentAttendanceLogs.Ctl.Create)
	r.PATCH("/students/:id/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Update)
	r.DELETE("/students/:id/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Delete)

	r.GET("/students/:id/payment-transactions", mod.PaymentTransactions.Ctl.List)
	r.POST("/students/:id/payment-transactions", mod.PaymentTransactions.Ctl.Create)
	r.PATCH("/students/:id/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Update)
	r.DELETE("/students/:id/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Delete)

	r.GET("/students/:id/fee-categories", mod.StudentFeeCategories.Ctl.List)
	r.POST("/students/:id/fee-categories", mod.StudentFeeCategories.Ctl.Create)
	r.PATCH("/students/:id/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Update)
	r.DELETE("/students/:id/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Delete)

	r.GET("/students/:id/grade-items", mod.StudentGradeItems.Ctl.List)
	r.POST("/students/:id/grade-items", mod.StudentGradeItems.Ctl.Create)
	r.PATCH("/students/:id/grade-items/:child_id", mod.StudentGradeItems.Ctl.Update)
	r.DELETE("/students/:id/grade-items/:child_id", mod.StudentGradeItems.Ctl.Delete)

	r.GET("/students/:id/grade-records", mod.StudentGradeRecords.Ctl.List)
	r.POST("/students/:id/grade-records", mod.StudentGradeRecords.Ctl.Create)
	r.PATCH("/students/:id/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Update)
	r.DELETE("/students/:id/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Delete)
}
