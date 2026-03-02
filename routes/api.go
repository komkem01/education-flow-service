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

	r.GET("/classrooms", mod.Classroom.Ctl.List)
	r.GET("/classrooms/:id", mod.Classroom.Ctl.Get)
	r.POST("/classrooms", mod.Classroom.Ctl.Create)
	r.PATCH("/classrooms/:id", mod.Classroom.Ctl.Update)
	r.DELETE("/classrooms/:id", mod.Classroom.Ctl.Delete)

	r.GET("/subjects", mod.Subject.Ctl.List)
	r.GET("/subjects/:id", mod.Subject.Ctl.Get)
	r.POST("/subjects", mod.Subject.Ctl.Create)
	r.PATCH("/subjects/:id", mod.Subject.Ctl.Update)
	r.DELETE("/subjects/:id", mod.Subject.Ctl.Delete)

	r.GET("/subject-assignments", mod.SubjectAssignment.Ctl.List)
	r.GET("/subject-assignments/:id", mod.SubjectAssignment.Ctl.Get)
	r.POST("/subject-assignments", mod.SubjectAssignment.Ctl.Create)
	r.PATCH("/subject-assignments/:id", mod.SubjectAssignment.Ctl.Update)
	r.DELETE("/subject-assignments/:id", mod.SubjectAssignment.Ctl.Delete)

	r.GET("/schedules", mod.Schedule.Ctl.List)
	r.GET("/schedules/:id", mod.Schedule.Ctl.Get)
	r.POST("/schedules", mod.Schedule.Ctl.Create)
	r.PATCH("/schedules/:id", mod.Schedule.Ctl.Update)
	r.DELETE("/schedules/:id", mod.Schedule.Ctl.Delete)

	r.GET("/question-banks", mod.QuestionBank.Ctl.List)
	r.GET("/question-banks/:id", mod.QuestionBank.Ctl.Get)
	r.POST("/question-banks", mod.QuestionBank.Ctl.Create)
	r.PATCH("/question-banks/:id", mod.QuestionBank.Ctl.Update)
	r.DELETE("/question-banks/:id", mod.QuestionBank.Ctl.Delete)

	r.GET("/question-choices", mod.QuestionChoice.Ctl.List)
	r.GET("/question-choices/:id", mod.QuestionChoice.Ctl.Get)
	r.POST("/question-choices", mod.QuestionChoice.Ctl.Create)
	r.PATCH("/question-choices/:id", mod.QuestionChoice.Ctl.Update)
	r.DELETE("/question-choices/:id", mod.QuestionChoice.Ctl.Delete)

	r.GET("/assessment-sets", mod.AssessmentSet.Ctl.List)
	r.GET("/assessment-sets/:id", mod.AssessmentSet.Ctl.Get)
	r.POST("/assessment-sets", mod.AssessmentSet.Ctl.Create)
	r.PATCH("/assessment-sets/:id", mod.AssessmentSet.Ctl.Update)
	r.DELETE("/assessment-sets/:id", mod.AssessmentSet.Ctl.Delete)

	r.GET("/members", mod.Member.Ctl.List)
	r.GET("/members/:id", mod.Member.Ctl.Get)
	r.POST("/members", mod.Member.Ctl.Create)
	r.PATCH("/members/:id", mod.Member.Ctl.Update)
	r.DELETE("/members/:id", mod.Member.Ctl.Delete)

	r.GET("/staffs", mod.Staff.Ctl.List)
	r.GET("/staffs/:id", mod.Staff.Ctl.Get)
	r.POST("/staffs", mod.Staff.Ctl.Create)
	r.PATCH("/staffs/:id", mod.Staff.Ctl.Update)
	r.DELETE("/staffs/:id", mod.Staff.Ctl.Delete)

	r.GET("/admins", mod.Admin.Ctl.List)
	r.GET("/admins/:id", mod.Admin.Ctl.Get)
	r.POST("/admins", mod.Admin.Ctl.Create)
	r.PATCH("/admins/:id", mod.Admin.Ctl.Update)
	r.DELETE("/admins/:id", mod.Admin.Ctl.Delete)

	r.GET("/parents", mod.Parent.Ctl.List)
	r.GET("/parents/:id", mod.Parent.Ctl.Get)
	r.POST("/parents", mod.Parent.Ctl.Create)
	r.PATCH("/parents/:id", mod.Parent.Ctl.Update)
	r.DELETE("/parents/:id", mod.Parent.Ctl.Delete)

	r.GET("/parents/:id/students", mod.ParentStudents.Ctl.List)
	r.POST("/parents/:id/students", mod.ParentStudents.Ctl.Create)
	r.PATCH("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Update)
	r.DELETE("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Delete)

	r.GET("/teachers", mod.Teacher.Ctl.List)
	r.GET("/teachers/:id", mod.Teacher.Ctl.Get)
	r.POST("/teachers", mod.Teacher.Ctl.Create)
	r.PATCH("/teachers/:id", mod.Teacher.Ctl.Update)
	r.DELETE("/teachers/:id", mod.Teacher.Ctl.Delete)

	r.GET("/students", mod.Student.Ctl.List)
	r.GET("/students/:id", mod.Student.Ctl.Get)
	r.POST("/students", mod.Student.Ctl.Create)
	r.PATCH("/students/:id", mod.Student.Ctl.Update)
	r.DELETE("/students/:id", mod.Student.Ctl.Delete)

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
