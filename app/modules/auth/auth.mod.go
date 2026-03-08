package auth

import (
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Module struct {
	tracer trace.Tracer
	Svc    *Service
	Ctl    *Controller
}

type Config struct {
	AccessTokenTTLMinutes int `conf:"optional"`
}

func New(conf *config.Config[Config], db interface {
	entitiesinf.MemberEntity
	entitiesinf.MemberRoleEntity
	entitiesinf.SchoolEntity
}, appKey string) *Module {
	tracer := otel.Tracer("education-flow.modules.auth")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		db:     db,
		appKey: appKey,
	})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
