package routes

import (
	"education-flow/app/modules"
	"education-flow/app/modules/entities/ent"

	"github.com/gin-gonic/gin"
)

func apiMasterData(r *gin.RouterGroup, mod *modules.Modules) {
	protected := r.Group("")
	protected.Use(requireAuth(mod), requireRoles(ent.MemberRoleAdmin, ent.MemberRoleStaff))

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
}
