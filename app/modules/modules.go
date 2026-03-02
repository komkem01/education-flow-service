package modules

import (
	"log/slog"
	"sync"

	academicyears "education-flow/app/modules/academic-years"
	"education-flow/app/modules/entities"
	"education-flow/app/modules/example"
	"education-flow/app/modules/genders"
	"education-flow/app/modules/members"
	"education-flow/app/modules/prefixes"
	"education-flow/app/modules/schools"
	"education-flow/app/modules/sentry"
	"education-flow/app/modules/specs"
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
	Teacher                      *teachers.Module
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
	teacherMod := teachers.New(config.Conf[teachers.Config](confMod.Svc), entitiesMod.Svc)
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
		Teacher:                      teacherMod,
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
