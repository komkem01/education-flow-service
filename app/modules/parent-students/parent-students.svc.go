package parentstudents

import (
	"context"
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.MemberParentStudentEntity
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.MemberParentStudentEntity
}

type CreateInput struct {
	ParentID       uuid.UUID
	StudentID      uuid.UUID
	Relationship   ent.ParentRelationship
	IsMainGuardian bool
}

type UpdateInput struct {
	StudentID      uuid.UUID
	Relationship   ent.ParentRelationship
	IsMainGuardian bool
}

var (
	ErrParentNotFound  = errors.New("parent-not-found")
	ErrStudentNotFound = errors.New("student-not-found")
)

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.MemberParentStudent, error) {
	parentExists, err := s.db.ParentExistsByID(ctx, input.ParentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, ErrParentNotFound
	}

	studentExists, err := s.db.StudentExistsByID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !studentExists {
		return nil, ErrStudentNotFound
	}

	item := &ent.MemberParentStudent{
		ParentID:       input.ParentID,
		StudentID:      input.StudentID,
		Relationship:   input.Relationship,
		IsMainGuardian: input.IsMainGuardian,
	}
	return s.db.CreateParentStudent(ctx, item)
}

func (s *Service) ListByParentID(ctx context.Context, parentID uuid.UUID) ([]*ent.MemberParentStudent, error) {
	return s.db.ListParentStudentsByParentID(ctx, parentID)
}

func (s *Service) UpdateByID(ctx context.Context, parentID uuid.UUID, id uuid.UUID, input *UpdateInput) (*ent.MemberParentStudent, error) {
	parentExists, err := s.db.ParentExistsByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if !parentExists {
		return nil, ErrParentNotFound
	}

	studentExists, err := s.db.StudentExistsByID(ctx, input.StudentID)
	if err != nil {
		return nil, err
	}
	if !studentExists {
		return nil, ErrStudentNotFound
	}

	belongs, err := s.db.ParentStudentBelongsToParent(ctx, id, parentID)
	if err != nil {
		return nil, err
	}
	if !belongs {
		return nil, sql.ErrNoRows
	}

	item := &ent.MemberParentStudent{ParentID: parentID, StudentID: input.StudentID, Relationship: input.Relationship, IsMainGuardian: input.IsMainGuardian}
	return s.db.UpdateParentStudentByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, parentID uuid.UUID, id uuid.UUID) error {
	belongs, err := s.db.ParentStudentBelongsToParent(ctx, id, parentID)
	if err != nil {
		return err
	}
	if !belongs {
		return sql.ErrNoRows
	}

	return s.db.DeleteParentStudentByID(ctx, id)
}
