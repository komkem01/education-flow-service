package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiBackOffice(r *gin.RouterGroup, mod *modules.Modules) {
	protected := r.Group("/back-office")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	protected.GET("/auth/me", mod.Auth.Ctl.Me)
	protected.GET("/auth/permissions", mod.Auth.Ctl.Permissions)

	registerCRUD(protected, "/academic-years", mod.AcademicYear.Ctl.List, mod.AcademicYear.Ctl.Get, mod.AcademicYear.Ctl.Create, mod.AcademicYear.Ctl.Update, mod.AcademicYear.Ctl.Delete)
	registerCRUD(protected, "/schools", mod.School.Ctl.List, mod.School.Ctl.Get, mod.School.Ctl.Create, mod.School.Ctl.Update, mod.School.Ctl.Delete)
	registerCRUD(protected, "/genders", mod.Gender.Ctl.List, mod.Gender.Ctl.Get, mod.Gender.Ctl.Create, mod.Gender.Ctl.Update, mod.Gender.Ctl.Delete)
	registerCRUD(protected, "/prefixes", mod.Prefix.Ctl.List, mod.Prefix.Ctl.Get, mod.Prefix.Ctl.Create, mod.Prefix.Ctl.Update, mod.Prefix.Ctl.Delete)
	registerCRUD(protected, "/classrooms", mod.Classroom.Ctl.List, mod.Classroom.Ctl.Get, mod.Classroom.Ctl.Create, mod.Classroom.Ctl.Update, mod.Classroom.Ctl.Delete)
	registerCRUD(protected, "/subjects", mod.Subject.Ctl.List, mod.Subject.Ctl.Get, mod.Subject.Ctl.Create, mod.Subject.Ctl.Update, mod.Subject.Ctl.Delete)
	registerCRUD(protected, "/subject-assignments", mod.SubjectAssignment.Ctl.List, mod.SubjectAssignment.Ctl.Get, mod.SubjectAssignment.Ctl.Create, mod.SubjectAssignment.Ctl.Update, mod.SubjectAssignment.Ctl.Delete)
	registerCRUD(protected, "/schedules", mod.Schedule.Ctl.List, mod.Schedule.Ctl.Get, mod.Schedule.Ctl.Create, mod.Schedule.Ctl.Update, mod.Schedule.Ctl.Delete)
	registerCRUD(protected, "/question-banks", mod.QuestionBank.Ctl.List, mod.QuestionBank.Ctl.Get, mod.QuestionBank.Ctl.Create, mod.QuestionBank.Ctl.Update, mod.QuestionBank.Ctl.Delete)
	registerCRUD(protected, "/question-choices", mod.QuestionChoice.Ctl.List, mod.QuestionChoice.Ctl.Get, mod.QuestionChoice.Ctl.Create, mod.QuestionChoice.Ctl.Update, mod.QuestionChoice.Ctl.Delete)
	registerCRUD(protected, "/assessment-sets", mod.AssessmentSet.Ctl.List, mod.AssessmentSet.Ctl.Get, mod.AssessmentSet.Ctl.Create, mod.AssessmentSet.Ctl.Update, mod.AssessmentSet.Ctl.Delete)

	registerCRUD(protected, "/members", mod.Member.Ctl.List, mod.Member.Ctl.Get, mod.Member.Ctl.Create, mod.Member.Ctl.Update, mod.Member.Ctl.Delete)
	protected.GET("/members/:id/roles", mod.Member.Ctl.ListRoles)
	protected.POST("/members/:id/roles", mod.Member.Ctl.AddRole)
	protected.DELETE("/members/:id/roles/:role", mod.Member.Ctl.RemoveRole)

	registerCRUD(protected, "/admins", mod.Admin.Ctl.List, mod.Admin.Ctl.Get, mod.Admin.Ctl.Create, mod.Admin.Ctl.Update, mod.Admin.Ctl.Delete)
	registerCRUD(protected, "/staffs", mod.Staff.Ctl.List, mod.Staff.Ctl.Get, mod.Staff.Ctl.Create, mod.Staff.Ctl.Update, mod.Staff.Ctl.Delete)
	registerCRUD(protected, "/teachers", mod.Teacher.Ctl.List, mod.Teacher.Ctl.Get, mod.Teacher.Ctl.Create, mod.Teacher.Ctl.Update, mod.Teacher.Ctl.Delete)
	registerCRUD(protected, "/students", mod.Student.Ctl.List, mod.Student.Ctl.Get, mod.Student.Ctl.Create, mod.Student.Ctl.Update, mod.Student.Ctl.Delete)
	protected.GET("/parents", mod.Parent.Ctl.List)
	protected.POST("/parents", mod.Parent.Ctl.Create)
	protected.GET("/parents/:id", mod.Parent.Ctl.Get)
	protected.PATCH("/parents/:id", mod.Parent.Ctl.Update)
	protected.DELETE("/parents/:id", mod.Parent.Ctl.Delete)
	protected.GET("/parents/:id/students", mod.ParentStudents.Ctl.List)
	protected.POST("/parents/:id/students", mod.ParentStudents.Ctl.Create)
	protected.PATCH("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Update)
	protected.DELETE("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Delete)

	registerCRUD(protected, "/inventory-items", mod.InventoryItem.Ctl.List, mod.InventoryItem.Ctl.Get, mod.InventoryItem.Ctl.Create, mod.InventoryItem.Ctl.Update, mod.InventoryItem.Ctl.Delete)
	registerCRUD(protected, "/inventory-requests", mod.InventoryRequest.Ctl.List, mod.InventoryRequest.Ctl.Get, mod.InventoryRequest.Ctl.Create, mod.InventoryRequest.Ctl.Update, mod.InventoryRequest.Ctl.Delete)
	registerCRUD(protected, "/document-tracking", mod.DocumentTracking.Ctl.List, mod.DocumentTracking.Ctl.Get, mod.DocumentTracking.Ctl.Create, mod.DocumentTracking.Ctl.Update, mod.DocumentTracking.Ctl.Delete)
	registerCRUD(protected, "/school-announcements", mod.SchoolAnnouncement.Ctl.List, mod.SchoolAnnouncement.Ctl.Get, mod.SchoolAnnouncement.Ctl.Create, mod.SchoolAnnouncement.Ctl.Update, mod.SchoolAnnouncement.Ctl.Delete)
	registerCRUD(protected, "/storages", mod.Storage.Ctl.List, mod.Storage.Ctl.Get, mod.Storage.Ctl.Create, mod.Storage.Ctl.Update, mod.Storage.Ctl.Delete)
	registerCRUD(protected, "/storage-links", mod.StorageLink.Ctl.List, mod.StorageLink.Ctl.Get, mod.StorageLink.Ctl.Create, mod.StorageLink.Ctl.Update, mod.StorageLink.Ctl.Delete)

	protected.GET("/students/:id/enrollments", mod.StudentEnrollments.Ctl.List)
	protected.POST("/students/:id/enrollments", mod.StudentEnrollments.Ctl.Create)
	protected.PATCH("/students/:id/enrollments/:child_id", mod.StudentEnrollments.Ctl.Update)
	protected.DELETE("/students/:id/enrollments/:child_id", mod.StudentEnrollments.Ctl.Delete)

	protected.GET("/students/:id/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.List)
	protected.POST("/students/:id/assessment-submissions", mod.StudentAssessmentSubmissions.Ctl.Create)
	protected.PATCH("/students/:id/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Update)
	protected.DELETE("/students/:id/assessment-submissions/:child_id", mod.StudentAssessmentSubmissions.Ctl.Delete)

	protected.GET("/students/:id/invoices", mod.StudentInvoices.Ctl.List)
	protected.POST("/students/:id/invoices", mod.StudentInvoices.Ctl.Create)
	protected.PATCH("/students/:id/invoices/:child_id", mod.StudentInvoices.Ctl.Update)
	protected.DELETE("/students/:id/invoices/:child_id", mod.StudentInvoices.Ctl.Delete)

	protected.GET("/students/:id/attendance-logs", mod.StudentAttendanceLogs.Ctl.List)
	protected.POST("/students/:id/attendance-logs", mod.StudentAttendanceLogs.Ctl.Create)
	protected.PATCH("/students/:id/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Update)
	protected.DELETE("/students/:id/attendance-logs/:child_id", mod.StudentAttendanceLogs.Ctl.Delete)

	protected.GET("/students/:id/payment-transactions", mod.PaymentTransactions.Ctl.List)
	protected.POST("/students/:id/payment-transactions", mod.PaymentTransactions.Ctl.Create)
	protected.PATCH("/students/:id/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Update)
	protected.DELETE("/students/:id/payment-transactions/:child_id", mod.PaymentTransactions.Ctl.Delete)

	protected.GET("/students/:id/fee-categories", mod.StudentFeeCategories.Ctl.List)
	protected.POST("/students/:id/fee-categories", mod.StudentFeeCategories.Ctl.Create)
	protected.PATCH("/students/:id/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Update)
	protected.DELETE("/students/:id/fee-categories/:child_id", mod.StudentFeeCategories.Ctl.Delete)

	protected.GET("/students/:id/grade-items", mod.StudentGradeItems.Ctl.List)
	protected.POST("/students/:id/grade-items", mod.StudentGradeItems.Ctl.Create)
	protected.PATCH("/students/:id/grade-items/:child_id", mod.StudentGradeItems.Ctl.Update)
	protected.DELETE("/students/:id/grade-items/:child_id", mod.StudentGradeItems.Ctl.Delete)

	protected.GET("/students/:id/grade-records", mod.StudentGradeRecords.Ctl.List)
	protected.POST("/students/:id/grade-records", mod.StudentGradeRecords.Ctl.Create)
	protected.PATCH("/students/:id/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Update)
	protected.DELETE("/students/:id/grade-records/:child_id", mod.StudentGradeRecords.Ctl.Delete)

	protected.GET("/teachers/:id/subject-assignments", mod.SubjectAssignment.Ctl.ListByTeacher)
	protected.POST("/teachers/:id/subject-assignments", mod.SubjectAssignment.Ctl.CreateByTeacher)
	protected.PATCH("/teachers/:id/subject-assignments/:child_id", mod.SubjectAssignment.Ctl.UpdateByTeacher)
	protected.DELETE("/teachers/:id/subject-assignments/:child_id", mod.SubjectAssignment.Ctl.DeleteByTeacher)

	protected.GET("/teachers/:id/educations", mod.TeacherEducations.Ctl.List)
	protected.POST("/teachers/:id/educations", mod.TeacherEducations.Ctl.Create)
	protected.PATCH("/teachers/:id/educations/:child_id", mod.TeacherEducations.Ctl.Update)
	protected.DELETE("/teachers/:id/educations/:child_id", mod.TeacherEducations.Ctl.Delete)

	protected.GET("/teachers/:id/work-experiences", mod.TeacherWorkExperiences.Ctl.List)
	protected.POST("/teachers/:id/work-experiences", mod.TeacherWorkExperiences.Ctl.Create)
	protected.PATCH("/teachers/:id/work-experiences/:child_id", mod.TeacherWorkExperiences.Ctl.Update)
	protected.DELETE("/teachers/:id/work-experiences/:child_id", mod.TeacherWorkExperiences.Ctl.Delete)

	protected.GET("/teachers/:id/profile-requests", mod.TeacherProfileRequests.Ctl.List)
	protected.POST("/teachers/:id/profile-requests", mod.TeacherProfileRequests.Ctl.Create)
	protected.PATCH("/teachers/:id/profile-requests/:child_id", mod.TeacherProfileRequests.Ctl.Update)

	protected.GET("/teachers/:id/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.List)
	protected.POST("/teachers/:id/performance-agreements", mod.TeacherPerformanceAgreements.Ctl.Create)
	protected.PATCH("/teachers/:id/performance-agreements/:child_id", mod.TeacherPerformanceAgreements.Ctl.Update)

	protected.GET("/teachers/:id/pda-logs", mod.TeacherPDALogs.Ctl.List)
	protected.POST("/teachers/:id/pda-logs", mod.TeacherPDALogs.Ctl.Create)
	protected.DELETE("/teachers/:id/pda-logs/:child_id", mod.TeacherPDALogs.Ctl.Delete)

	protected.GET("/teachers/:id/leave-logs", mod.TeacherLeaveLogs.Ctl.List)
	protected.POST("/teachers/:id/leave-logs", mod.TeacherLeaveLogs.Ctl.Create)
	protected.PATCH("/teachers/:id/leave-logs/:child_id", mod.TeacherLeaveLogs.Ctl.Update)

	protected.GET("/system-audit-logs", mod.SystemAuditLog.Ctl.List)
	protected.GET("/system-audit-logs/:id", mod.SystemAuditLog.Ctl.Get)
	protected.POST("/system-audit-logs", mod.SystemAuditLog.Ctl.Create)

	protected.GET("/data-change-logs", mod.DataChangeLog.Ctl.List)
	protected.GET("/data-change-logs/:id", mod.DataChangeLog.Ctl.Get)
	protected.POST("/data-change-logs", mod.DataChangeLog.Ctl.Create)
}
