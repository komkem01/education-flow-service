package routes

import (
	"fmt"
	"net/http"

	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func WarpH(router *gin.RouterGroup, prefix string, handler http.Handler) {
	router.Any(fmt.Sprintf("%s/*w", prefix), gin.WrapH(http.StripPrefix(fmt.Sprintf("%s%s", router.BasePath(), prefix), handler)))
}

// api keeps all v1 route registrations centralized in this file.
func api(r *gin.RouterGroup, mod *modules.Modules) {
	// Example routes.
	r.GET("/example/:id", mod.Example.Ctl.Get)
	r.GET("/example-http", mod.Example.Ctl.GetHttpReq)
	r.POST("/example", mod.Example.Ctl.Create)

	// Domain route groups.
	apiAuth(r, mod)
	apiMasterData(r, mod)
	apiMemberTeachers(r, mod)
	apiMemberStudents(r, mod)
	apiMemberStaffs(r, mod)
	apiMemberAdmins(r, mod)
	apiBackOffice(r, mod)
}

// registerCRUD provides a standard REST mapping for modules using list/get/create/update/delete handlers.
func registerCRUD(r *gin.RouterGroup, path string, list gin.HandlerFunc, get gin.HandlerFunc, create gin.HandlerFunc, update gin.HandlerFunc, del gin.HandlerFunc) {
	r.GET(path, list)
	r.GET(path+"/:id", get)
	r.POST(path, create)

	if update != nil {
		r.PATCH(path+"/:id", update)
	}

	if del != nil {
		r.DELETE(path+"/:id", del)
	}
}

// apiAuth registers authentication endpoints.
func apiAuth(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/auth/login", mod.Auth.Ctl.Login)

	authProtected := r.Group("/auth")
	authProtected.Use(requireAuth(mod))
	authProtected.GET("/me", mod.Auth.Ctl.Me)
	authProtected.GET("/permissions", mod.Auth.Ctl.Permissions)
	authProtected.POST("/logout", mod.Auth.Ctl.Logout)
	authProtected.POST("/refresh", mod.Auth.Ctl.Refresh)
	authProtected.POST("/switch-role", mod.Auth.Ctl.SwitchRole)
	authProtected.POST("/switch-school", mod.Auth.Ctl.SwitchSchool)
}

// apiMasterData registers master-data endpoints for admin/staff roles.
func apiMasterData(r *gin.RouterGroup, mod *modules.Modules) {
	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	registerCRUD(protected, "/academic-years", mod.AcademicYear.Ctl.List, mod.AcademicYear.Ctl.Get, mod.AcademicYear.Ctl.Create, mod.AcademicYear.Ctl.Update, mod.AcademicYear.Ctl.Delete)
	registerCRUD(protected, "/schools", mod.School.Ctl.List, mod.School.Ctl.Get, mod.School.Ctl.Create, mod.School.Ctl.Update, mod.School.Ctl.Delete)
	registerCRUD(protected, "/genders", mod.Gender.Ctl.List, mod.Gender.Ctl.Get, mod.Gender.Ctl.Create, mod.Gender.Ctl.Update, mod.Gender.Ctl.Delete)
	registerCRUD(protected, "/prefixes", mod.Prefix.Ctl.List, mod.Prefix.Ctl.Get, mod.Prefix.Ctl.Create, mod.Prefix.Ctl.Update, mod.Prefix.Ctl.Delete)
	registerCRUD(protected, "/classrooms", mod.Classroom.Ctl.List, mod.Classroom.Ctl.Get, mod.Classroom.Ctl.Create, mod.Classroom.Ctl.Update, mod.Classroom.Ctl.Delete)
	registerCRUD(protected, "/subjects", mod.Subject.Ctl.List, mod.Subject.Ctl.Get, mod.Subject.Ctl.Create, mod.Subject.Ctl.Update, mod.Subject.Ctl.Delete)
	registerCRUD(protected, "/subject-groups", mod.SubjectGroup.Ctl.List, mod.SubjectGroup.Ctl.Get, mod.SubjectGroup.Ctl.Create, mod.SubjectGroup.Ctl.Update, mod.SubjectGroup.Ctl.Delete)
	registerCRUD(protected, "/subject-subgroups", mod.SubjectSubgroup.Ctl.List, mod.SubjectSubgroup.Ctl.Get, mod.SubjectSubgroup.Ctl.Create, mod.SubjectSubgroup.Ctl.Update, mod.SubjectSubgroup.Ctl.Delete)
	registerCRUD(protected, "/subject-assignments", mod.SubjectAssignment.Ctl.List, mod.SubjectAssignment.Ctl.Get, mod.SubjectAssignment.Ctl.Create, mod.SubjectAssignment.Ctl.Update, mod.SubjectAssignment.Ctl.Delete)
	registerCRUD(protected, "/schedules", mod.Schedule.Ctl.List, mod.Schedule.Ctl.Get, mod.Schedule.Ctl.Create, mod.Schedule.Ctl.Update, mod.Schedule.Ctl.Delete)
	registerCRUD(protected, "/question-banks", mod.QuestionBank.Ctl.List, mod.QuestionBank.Ctl.Get, mod.QuestionBank.Ctl.Create, mod.QuestionBank.Ctl.Update, mod.QuestionBank.Ctl.Delete)
	registerCRUD(protected, "/question-choices", mod.QuestionChoice.Ctl.List, mod.QuestionChoice.Ctl.Get, mod.QuestionChoice.Ctl.Create, mod.QuestionChoice.Ctl.Update, mod.QuestionChoice.Ctl.Delete)
	registerCRUD(protected, "/assessment-sets", mod.AssessmentSet.Ctl.List, mod.AssessmentSet.Ctl.Get, mod.AssessmentSet.Ctl.Create, mod.AssessmentSet.Ctl.Update, mod.AssessmentSet.Ctl.Delete)
}

// apiMemberAdmins registers admin/member/parent management endpoints.
func apiMemberAdmins(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/admins/register", mod.Admin.Ctl.Register)
	r.POST("/parents/register", mod.Parent.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	registerCRUD(protected, "/members", mod.Member.Ctl.List, mod.Member.Ctl.Get, mod.Member.Ctl.Create, mod.Member.Ctl.Update, mod.Member.Ctl.Delete)
	protected.GET("/members/:id/roles", mod.Member.Ctl.ListRoles)
	protected.POST("/members/:id/roles", mod.Member.Ctl.AddRole)
	protected.DELETE("/members/:id/roles/:role", mod.Member.Ctl.RemoveRole)
	registerCRUD(protected, "/admins", mod.Admin.Ctl.List, mod.Admin.Ctl.Get, mod.Admin.Ctl.Create, mod.Admin.Ctl.Update, mod.Admin.Ctl.Delete)

	protected.GET("/parents", mod.Parent.Ctl.List)
	protected.POST("/parents", mod.Parent.Ctl.Create)

	parentOwned := protected.Group("/parents/:id")
	parentOwned.Use(requireParentResourceOwnerOrRoles(mod, ent.MemberRoleAdmin, ent.MemberRoleStaff))
	parentOwned.GET("", mod.Parent.Ctl.Get)
	parentOwned.PATCH("", mod.Parent.Ctl.Update)
	parentOwned.DELETE("", mod.Parent.Ctl.Delete)

	parentOwned.GET("/students", mod.ParentStudents.Ctl.List)
	parentOwned.POST("/students", mod.ParentStudents.Ctl.Create)
	parentOwned.PATCH("/students/:child_id", mod.ParentStudents.Ctl.Update)
	parentOwned.DELETE("/students/:child_id", mod.ParentStudents.Ctl.Delete)
}

// apiMemberStaffs registers staff and operation inventory/document endpoints.
func apiMemberStaffs(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/staffs/register", mod.Staff.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

	registerCRUD(protected, "/staffs", mod.Staff.Ctl.List, mod.Staff.Ctl.Get, mod.Staff.Ctl.Create, mod.Staff.Ctl.Update, mod.Staff.Ctl.Delete)
	registerCRUD(protected, "/inventory-items", mod.InventoryItem.Ctl.List, mod.InventoryItem.Ctl.Get, mod.InventoryItem.Ctl.Create, mod.InventoryItem.Ctl.Update, mod.InventoryItem.Ctl.Delete)
	registerCRUD(protected, "/inventory-requests", mod.InventoryRequest.Ctl.List, mod.InventoryRequest.Ctl.Get, mod.InventoryRequest.Ctl.Create, mod.InventoryRequest.Ctl.Update, mod.InventoryRequest.Ctl.Delete)
	registerCRUD(protected, "/document-tracking", mod.DocumentTracking.Ctl.List, mod.DocumentTracking.Ctl.Get, mod.DocumentTracking.Ctl.Create, mod.DocumentTracking.Ctl.Update, mod.DocumentTracking.Ctl.Delete)
	registerCRUD(protected, "/school-announcements", mod.SchoolAnnouncement.Ctl.List, mod.SchoolAnnouncement.Ctl.Get, mod.SchoolAnnouncement.Ctl.Create, mod.SchoolAnnouncement.Ctl.Update, mod.SchoolAnnouncement.Ctl.Delete)
	registerCRUD(protected, "/storages", mod.Storage.Ctl.List, mod.Storage.Ctl.Get, mod.Storage.Ctl.Create, mod.Storage.Ctl.Update, mod.Storage.Ctl.Delete)
	registerCRUD(protected, "/storage-links", mod.StorageLink.Ctl.List, mod.StorageLink.Ctl.Get, mod.StorageLink.Ctl.Create, mod.StorageLink.Ctl.Update, mod.StorageLink.Ctl.Delete)
}

// apiMemberStudents registers student endpoints and student-owned resources.
func apiMemberStudents(r *gin.RouterGroup, mod *modules.Modules) {
	r.POST("/students/register", mod.Student.Ctl.Register)

	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher, ent.MemberRoleStudent, ent.MemberRoleParent))

	protected.GET("/students", mod.Student.Ctl.List)
	protected.POST("/students", mod.Student.Ctl.Create)

	// Read-only metadata endpoints used by student-facing UI.
	protected.GET("/students-meta/schools/:id", mod.School.Ctl.Get)
	protected.GET("/students-meta/academic-years/:id", mod.AcademicYear.Ctl.Get)
	protected.GET("/students-meta/classrooms/:id", mod.Classroom.Ctl.Get)
	protected.GET("/students-meta/genders/:id", mod.Gender.Ctl.Get)
	protected.GET("/students-meta/subjects/:id", mod.Subject.Ctl.Get)
	protected.GET("/students-meta/teachers/:id", mod.Teacher.Ctl.Get)
	protected.GET("/students-meta/prefixes/:id", mod.Prefix.Ctl.Get)
	protected.GET("/students-meta/subject-assignments", mod.SubjectAssignment.Ctl.List)
	protected.GET("/students-meta/subject-assignments/:id", mod.SubjectAssignment.Ctl.Get)
	protected.GET("/students-meta/schedules", mod.Schedule.Ctl.List)

	studentOwned := protected.Group("/students/:id")
	studentOwned.Use(requireStudentResourceAccess(mod, ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleTeacher))
	studentOwned.GET("", mod.Student.Ctl.Get)
	studentOwned.PATCH("", mod.Student.Ctl.Update)
	studentOwned.DELETE("", mod.Student.Ctl.Delete)

	studentOwned.GET("/enrollments", mod.StudentEnrollments.Ctl.List)
	studentOwned.GET("/parents", mod.ParentStudents.Ctl.ListByStudent)
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

// apiMemberTeachers registers teacher endpoints and teacher-owned resources.
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

// apiBackOffice registers back-office endpoints for admin/staff/super-admin roles.
func apiBackOffice(r *gin.RouterGroup, mod *modules.Modules) {
	protected := r.Group("/back-office")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff, ent.MemberRoleSuperAdmin))

	// Auth and report endpoints.
	registerBackOfficeAuthAndReports(protected, mod)

	// Onboarding endpoints.
	registerBackOfficeOnboarding(protected, mod)

	// Master data endpoints.
	registerBackOfficeMasterData(protected, mod)

	// Member and role endpoints.
	registerBackOfficeMembers(protected, mod)

	// Operations and storage endpoints.
	registerBackOfficeOperations(protected, mod)

	// Student nested resources.
	registerBackOfficeStudentNested(protected, mod)

	// Teacher nested resources.
	registerBackOfficeTeacherNested(protected, mod)

	// Audit and change-log endpoints.
	registerBackOfficeAudit(protected, mod)
}

func registerBackOfficeAuthAndReports(protected *gin.RouterGroup, mod *modules.Modules) {
	protected.GET("/auth/me", mod.Auth.Ctl.Me)
	protected.GET("/auth/permissions", mod.Auth.Ctl.Permissions)
	protected.POST("/auth/switch-role", mod.Auth.Ctl.SwitchRole)
	protected.POST("/auth/switch-school", mod.Auth.Ctl.SwitchSchool)
	protected.GET("/reports/filters", mod.Report.Ctl.ListFilters)
	protected.GET("/reports/summary", mod.Report.Ctl.Summary)
	protected.GET("/reports/approvals", mod.Report.Ctl.ListApprovals)
	protected.GET("/reports/roles-members", mod.Report.Ctl.ListRoleMembers)
	protected.GET("/reports/approvals/:type/:id", mod.Report.Ctl.GetApproval)
	protected.PATCH("/reports/approvals/:type/:id", mod.Report.Ctl.UpdateApproval)
}

func registerBackOfficeOnboarding(protected *gin.RouterGroup, mod *modules.Modules) {
	protected.POST("/onboarding/staffs/register", mod.Staff.Ctl.Register)
	protected.POST("/onboarding/teachers/register", mod.Teacher.Ctl.Register)
	protected.POST("/onboarding/students/register", mod.Student.Ctl.Register)
	protected.POST("/onboarding/parents/register", mod.Parent.Ctl.Register)
}

func registerBackOfficeMasterData(protected *gin.RouterGroup, mod *modules.Modules) {
	registerCRUD(protected, "/academic-years", mod.AcademicYear.Ctl.List, mod.AcademicYear.Ctl.Get, mod.AcademicYear.Ctl.Create, mod.AcademicYear.Ctl.Update, mod.AcademicYear.Ctl.Delete)
	registerCRUD(protected, "/schools", mod.School.Ctl.List, mod.School.Ctl.Get, mod.School.Ctl.Create, mod.School.Ctl.Update, mod.School.Ctl.Delete)
	registerCRUD(protected, "/genders", mod.Gender.Ctl.List, mod.Gender.Ctl.Get, mod.Gender.Ctl.Create, mod.Gender.Ctl.Update, mod.Gender.Ctl.Delete)
	registerCRUD(protected, "/prefixes", mod.Prefix.Ctl.List, mod.Prefix.Ctl.Get, mod.Prefix.Ctl.Create, mod.Prefix.Ctl.Update, mod.Prefix.Ctl.Delete)
	registerCRUD(protected, "/classrooms", mod.Classroom.Ctl.List, mod.Classroom.Ctl.Get, mod.Classroom.Ctl.Create, mod.Classroom.Ctl.Update, mod.Classroom.Ctl.Delete)
	registerCRUD(protected, "/departments", mod.Department.Ctl.List, mod.Department.Ctl.Get, mod.Department.Ctl.Create, mod.Department.Ctl.Update, mod.Department.Ctl.Delete)
	registerCRUD(protected, "/subjects", mod.Subject.Ctl.List, mod.Subject.Ctl.Get, mod.Subject.Ctl.Create, mod.Subject.Ctl.Update, mod.Subject.Ctl.Delete)
	registerCRUD(protected, "/courses", mod.SubjectGroup.Ctl.List, mod.SubjectGroup.Ctl.Get, mod.SubjectGroup.Ctl.Create, mod.SubjectGroup.Ctl.Update, mod.SubjectGroup.Ctl.Delete)
	registerCRUD(protected, "/subject-groups", mod.SubjectGroup.Ctl.List, mod.SubjectGroup.Ctl.Get, mod.SubjectGroup.Ctl.Create, mod.SubjectGroup.Ctl.Update, mod.SubjectGroup.Ctl.Delete)
	registerCRUD(protected, "/subject-subgroups", mod.SubjectSubgroup.Ctl.List, mod.SubjectSubgroup.Ctl.Get, mod.SubjectSubgroup.Ctl.Create, mod.SubjectSubgroup.Ctl.Update, mod.SubjectSubgroup.Ctl.Delete)
	registerCRUD(protected, "/subject-assignments", mod.SubjectAssignment.Ctl.List, mod.SubjectAssignment.Ctl.Get, mod.SubjectAssignment.Ctl.Create, mod.SubjectAssignment.Ctl.Update, mod.SubjectAssignment.Ctl.Delete)
	registerCRUD(protected, "/schedules", mod.Schedule.Ctl.List, mod.Schedule.Ctl.Get, mod.Schedule.Ctl.Create, mod.Schedule.Ctl.Update, mod.Schedule.Ctl.Delete)
	registerCRUD(protected, "/question-banks", mod.QuestionBank.Ctl.List, mod.QuestionBank.Ctl.Get, mod.QuestionBank.Ctl.Create, mod.QuestionBank.Ctl.Update, mod.QuestionBank.Ctl.Delete)
	registerCRUD(protected, "/question-choices", mod.QuestionChoice.Ctl.List, mod.QuestionChoice.Ctl.Get, mod.QuestionChoice.Ctl.Create, mod.QuestionChoice.Ctl.Update, mod.QuestionChoice.Ctl.Delete)
	registerCRUD(protected, "/assessment-sets", mod.AssessmentSet.Ctl.List, mod.AssessmentSet.Ctl.Get, mod.AssessmentSet.Ctl.Create, mod.AssessmentSet.Ctl.Update, mod.AssessmentSet.Ctl.Delete)
}

func registerBackOfficeMembers(protected *gin.RouterGroup, mod *modules.Modules) {
	registerCRUD(protected, "/members", mod.Member.Ctl.List, mod.Member.Ctl.Get, mod.Member.Ctl.Create, mod.Member.Ctl.Update, mod.Member.Ctl.Delete)
	protected.GET("/members/:id/roles", mod.Member.Ctl.ListRoles)
	protected.POST("/members/:id/roles", mod.Member.Ctl.AddRole)
	protected.DELETE("/members/:id/roles/:role", mod.Member.Ctl.RemoveRole)

	registerCRUD(protected, "/admins", mod.Admin.Ctl.List, mod.Admin.Ctl.Get, mod.Admin.Ctl.Create, mod.Admin.Ctl.Update, mod.Admin.Ctl.Delete)
	registerCRUD(protected, "/staffs", mod.Staff.Ctl.List, mod.Staff.Ctl.Get, mod.Staff.Ctl.Create, mod.Staff.Ctl.Update, mod.Staff.Ctl.Delete)
	registerCRUD(protected, "/teachers", mod.Teacher.Ctl.List, mod.Teacher.Ctl.Get, mod.Teacher.Ctl.Create, mod.Teacher.Ctl.Update, mod.Teacher.Ctl.Delete)
	registerCRUD(protected, "/students", mod.Student.Ctl.List, mod.Student.Ctl.Get, mod.Student.Ctl.Create, mod.Student.Ctl.Update, mod.Student.Ctl.Delete)
	protected.GET("/students/:id/parents", mod.ParentStudents.Ctl.ListByStudent)
	protected.GET("/parents", mod.Parent.Ctl.List)
	protected.POST("/parents", mod.Parent.Ctl.Create)
	protected.GET("/parents/:id", mod.Parent.Ctl.Get)
	protected.PATCH("/parents/:id", mod.Parent.Ctl.Update)
	protected.DELETE("/parents/:id", mod.Parent.Ctl.Delete)
	protected.GET("/parents/:id/students", mod.ParentStudents.Ctl.List)
	protected.POST("/parents/:id/students", mod.ParentStudents.Ctl.Create)
	protected.PATCH("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Update)
	protected.DELETE("/parents/:id/students/:child_id", mod.ParentStudents.Ctl.Delete)
}

func registerBackOfficeOperations(protected *gin.RouterGroup, mod *modules.Modules) {
	registerCRUD(protected, "/inventory-items", mod.InventoryItem.Ctl.List, mod.InventoryItem.Ctl.Get, mod.InventoryItem.Ctl.Create, mod.InventoryItem.Ctl.Update, mod.InventoryItem.Ctl.Delete)
	registerCRUD(protected, "/inventory-requests", mod.InventoryRequest.Ctl.List, mod.InventoryRequest.Ctl.Get, mod.InventoryRequest.Ctl.Create, mod.InventoryRequest.Ctl.Update, mod.InventoryRequest.Ctl.Delete)
	registerCRUD(protected, "/document-tracking", mod.DocumentTracking.Ctl.List, mod.DocumentTracking.Ctl.Get, mod.DocumentTracking.Ctl.Create, mod.DocumentTracking.Ctl.Update, mod.DocumentTracking.Ctl.Delete)
	registerCRUD(protected, "/school-announcements", mod.SchoolAnnouncement.Ctl.List, mod.SchoolAnnouncement.Ctl.Get, mod.SchoolAnnouncement.Ctl.Create, mod.SchoolAnnouncement.Ctl.Update, mod.SchoolAnnouncement.Ctl.Delete)
	registerCRUD(protected, "/school-calendar-events", mod.SchoolCalendarEvent.Ctl.List, mod.SchoolCalendarEvent.Ctl.Get, mod.SchoolCalendarEvent.Ctl.Create, mod.SchoolCalendarEvent.Ctl.Update, mod.SchoolCalendarEvent.Ctl.Delete)
	registerCRUD(protected, "/student-behaviors", mod.StudentBehavior.Ctl.List, mod.StudentBehavior.Ctl.Get, mod.StudentBehavior.Ctl.Create, mod.StudentBehavior.Ctl.Update, mod.StudentBehavior.Ctl.Delete)
	registerCRUD(protected, "/storages", mod.Storage.Ctl.List, mod.Storage.Ctl.Get, mod.Storage.Ctl.Create, mod.Storage.Ctl.Update, mod.Storage.Ctl.Delete)
	registerCRUD(protected, "/storage-links", mod.StorageLink.Ctl.List, mod.StorageLink.Ctl.Get, mod.StorageLink.Ctl.Create, mod.StorageLink.Ctl.Update, mod.StorageLink.Ctl.Delete)
}

func registerBackOfficeStudentNested(protected *gin.RouterGroup, mod *modules.Modules) {
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
}

func registerBackOfficeTeacherNested(protected *gin.RouterGroup, mod *modules.Modules) {
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
}

func registerBackOfficeAudit(protected *gin.RouterGroup, mod *modules.Modules) {
	protected.GET("/system-audit-logs", mod.SystemAuditLog.Ctl.List)
	protected.GET("/system-audit-logs/:id", mod.SystemAuditLog.Ctl.Get)
	protected.POST("/system-audit-logs", mod.SystemAuditLog.Ctl.Create)

	protected.GET("/data-change-logs", mod.DataChangeLog.Ctl.List)
	protected.GET("/data-change-logs/:id", mod.DataChangeLog.Ctl.Get)
	protected.POST("/data-change-logs", mod.DataChangeLog.Ctl.Create)
}
