package database

import (
	"context"

	dto "education-flow/internal/database/dto"
	"education-flow/internal/provider"
)

type DatabaseModule struct {
	Svc *DatabaseService
}

var _ provider.Close = (*DatabaseModule)(nil)

func New(opts map[string]*dto.Option) *DatabaseModule {
	service := newService(opts)
	return &DatabaseModule{
		Svc: service,
	}
}

func (db *DatabaseModule) Close(ctx context.Context) error {
	return db.Svc.close(ctx)
}
