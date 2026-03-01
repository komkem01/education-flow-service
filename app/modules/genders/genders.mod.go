package genders

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

func New(conf *config.Config[Config], db entitiesinf.GenderEntity) *Module {
	tracer := otel.Tracer("education-flow.modules.genders")
	svc := newService(&Options{
		Config: conf,
		tracer: tracer,
		db:     db,
	})

	return &Module{
		tracer: tracer,
		Svc:    svc,
		Ctl:    newController(tracer, svc),
	}
}
