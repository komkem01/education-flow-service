package modules

import (
	"log/slog"
	"sync"

	academicyears "education-flow/app/modules/academic-years"
	"education-flow/app/modules/admins"
	"education-flow/app/modules/entities"
	"education-flow/app/modules/example"
	"education-flow/app/modules/genders"
	"education-flow/app/modules/members"
	parentstudents "education-flow/app/modules/parent-students"
	"education-flow/app/modules/parents"
	paymenttransactions "education-flow/app/modules/payment-transactions"
	"education-flow/app/modules/prefixes"
	"education-flow/app/modules/schools"
	"education-flow/app/modules/sentry"
	"education-flow/app/modules/specs"
	"education-flow/app/modules/staffs"
	studentassessmentsubmissions "education-flow/app/modules/student-assessment-submissions"
	studentattendancelogs "education-flow/app/modules/student-attendance-logs"
	studentenrollments "education-flow/app/modules/student-enrollments"
	studentfeecategories "education-flow/app/modules/student-fee-categories"
	studentgradeitems "education-flow/app/modules/student-grade-items"
	studentgraderecords "education-flow/app/modules/student-grade-records"
	studentinvoices "education-flow/app/modules/student-invoices"
	"education-flow/app/modules/students"
	teachereducations "education-flow/app/modules/teacher-educations"
	teacherleavelogs "education-flow/app/modules/teacher-leave-logs"
	teacherpdalogs "education-flow/app/modules/teacher-pda-logs"
	teacherperformanceagreements "education-flow/app/modules/teacher-performance-agreements"
	teacherprofilerequests "education-flow/app/modules/teacher-profile-requests"
	"education-flow/app/modules/teachers"
	"education-flow/internal/config"
	"education-flow/internal/database"
	"education-flow/internal/log"
	"education-flow/internal/otel/collector"

	exampletwo "education-flow/app/modules/example-two"

	appConf "education-flow/config"
	// "education-flow/app/modules/kafka"
)

type Modules struct {
	Conf                         *config.Module[appConf.Config]
	Specs                        *specs.Module
	Log                          *log.Module
	OTEL                         *collector.Module
	Sentry                       *sentry.Module
	DB                           *database.DatabaseModule
	ENT                          *entities.Module
	AcademicYear                 *academicyears.Module
	School                       *schools.Module
	Gender                       *genders.Module
	Prefix                       *prefixes.Module
	Member                       *members.Module
	Admin                        *admins.Module
	Staff                        *staffs.Module
	Parent                       *parents.Module
	ParentStudents               *parentstudents.Module
	Teacher                      *teachers.Module
	Student                      *students.Module
	StudentEnrollments           *studentenrollments.Module
	StudentAssessmentSubmissions *studentassessmentsubmissions.Module
	StudentInvoices              *studentinvoices.Module
	StudentAttendanceLogs        *studentattendancelogs.Module
	PaymentTransactions          *paymenttransactions.Module
	StudentFeeCategories         *studentfeecategories.Module
	StudentGradeItems            *studentgradeitems.Module
	StudentGradeRecords          *studentgraderecords.Module
	TeacherEducations            *teachereducations.Module
	TeacherProfileRequests       *teacherprofilerequests.Module
	TeacherPerformanceAgreements *teacherperformanceagreements.Module
	TeacherPDALogs               *teacherpdalogs.Module
	TeacherLeaveLogs             *teacherleavelogs.Module
	// Kafka *kafka.Module
	Example  *example.Module
	Example2 *exampletwo.Module
}

func modulesInit() {
	confMod := config.New(&appConf.App)
	specsMod := specs.New(config.Conf[specs.Config](confMod.Svc))
	conf := confMod.Svc.Config()

	logMod := log.New(config.Conf[log.Option](confMod.Svc))
	otel := collector.New(config.Conf[collector.Config](confMod.Svc))
	log := log.With(slog.String("module", "modules"))

	sentryMod := sentry.New(config.Conf[sentry.Config](confMod.Svc))

	db := database.New(conf.Database.Sql)
	entitiesMod := entities.New(db.Svc.DB())
	academicYearMod := academicyears.New(config.Conf[academicyears.Config](confMod.Svc), entitiesMod.Svc)
	schoolMod := schools.New(config.Conf[schools.Config](confMod.Svc), entitiesMod.Svc)
	genderMod := genders.New(config.Conf[genders.Config](confMod.Svc), entitiesMod.Svc)
	prefixMod := prefixes.New(config.Conf[prefixes.Config](confMod.Svc), entitiesMod.Svc)
	memberMod := members.New(config.Conf[members.Config](confMod.Svc), entitiesMod.Svc)
	adminMod := admins.New(config.Conf[admins.Config](confMod.Svc), entitiesMod.Svc)
	staffMod := staffs.New(config.Conf[staffs.Config](confMod.Svc), entitiesMod.Svc)
	parentMod := parents.New(config.Conf[parents.Config](confMod.Svc), entitiesMod.Svc)
	parentStudentsMod := parentstudents.New(config.Conf[parentstudents.Config](confMod.Svc), entitiesMod.Svc)
	teacherMod := teachers.New(config.Conf[teachers.Config](confMod.Svc), entitiesMod.Svc)
	studentMod := students.New(config.Conf[students.Config](confMod.Svc), entitiesMod.Svc)
	studentEnrollmentsMod := studentenrollments.New(config.Conf[studentenrollments.Config](confMod.Svc), entitiesMod.Svc)
	studentAssessmentSubmissionsMod := studentassessmentsubmissions.New(config.Conf[studentassessmentsubmissions.Config](confMod.Svc), entitiesMod.Svc)
	studentInvoicesMod := studentinvoices.New(config.Conf[studentinvoices.Config](confMod.Svc), entitiesMod.Svc)
	studentAttendanceLogsMod := studentattendancelogs.New(config.Conf[studentattendancelogs.Config](confMod.Svc), entitiesMod.Svc)
	paymentTransactionsMod := paymenttransactions.New(config.Conf[paymenttransactions.Config](confMod.Svc), entitiesMod.Svc)
	studentFeeCategoriesMod := studentfeecategories.New(config.Conf[studentfeecategories.Config](confMod.Svc), entitiesMod.Svc)
	studentGradeItemsMod := studentgradeitems.New(config.Conf[studentgradeitems.Config](confMod.Svc), entitiesMod.Svc)
	studentGradeRecordsMod := studentgraderecords.New(config.Conf[studentgraderecords.Config](confMod.Svc), entitiesMod.Svc)
	teacherEducationsMod := teachereducations.New(config.Conf[teachereducations.Config](confMod.Svc), entitiesMod.Svc)
	teacherProfileRequestsMod := teacherprofilerequests.New(config.Conf[teacherprofilerequests.Config](confMod.Svc), entitiesMod.Svc)
	teacherPerformanceAgreementsMod := teacherperformanceagreements.New(config.Conf[teacherperformanceagreements.Config](confMod.Svc), entitiesMod.Svc)
	teacherPDALogsMod := teacherpdalogs.New(config.Conf[teacherpdalogs.Config](confMod.Svc), entitiesMod.Svc)
	teacherLeaveLogsMod := teacherleavelogs.New(config.Conf[teacherleavelogs.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod := example.New(config.Conf[example.Config](confMod.Svc), entitiesMod.Svc)
	exampleMod2 := exampletwo.New(config.Conf[exampletwo.Config](confMod.Svc), entitiesMod.Svc)
	// kafka := kafka.New(&conf.Kafka)
	mod = &Modules{
		Conf:                         confMod,
		Specs:                        specsMod,
		Log:                          logMod,
		OTEL:                         otel,
		Sentry:                       sentryMod,
		DB:                           db,
		ENT:                          entitiesMod,
		AcademicYear:                 academicYearMod,
		School:                       schoolMod,
		Gender:                       genderMod,
		Prefix:                       prefixMod,
		Member:                       memberMod,
		Admin:                        adminMod,
		Staff:                        staffMod,
		Parent:                       parentMod,
		ParentStudents:               parentStudentsMod,
		Teacher:                      teacherMod,
		Student:                      studentMod,
		StudentEnrollments:           studentEnrollmentsMod,
		StudentAssessmentSubmissions: studentAssessmentSubmissionsMod,
		StudentInvoices:              studentInvoicesMod,
		StudentAttendanceLogs:        studentAttendanceLogsMod,
		PaymentTransactions:          paymentTransactionsMod,
		StudentFeeCategories:         studentFeeCategoriesMod,
		StudentGradeItems:            studentGradeItemsMod,
		StudentGradeRecords:          studentGradeRecordsMod,
		TeacherEducations:            teacherEducationsMod,
		TeacherProfileRequests:       teacherProfileRequestsMod,
		TeacherPerformanceAgreements: teacherPerformanceAgreementsMod,
		TeacherPDALogs:               teacherPDALogsMod,
		TeacherLeaveLogs:             teacherLeaveLogsMod,
		Example:                      exampleMod,
		Example2:                     exampleMod2,
	}

	log.Infof("all modules initialized")
}

var (
	once sync.Once
	mod  *Modules
)

func Get() *Modules {
	once.Do(modulesInit)

	return mod
}
