package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SubjectEntity = (*Service)(nil)

func (s *Service) CreateSubject(ctx context.Context, subject *ent.Subject) (*ent.Subject, error) {
	if _, err := s.db.NewInsert().Model(subject).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *Service) GetSubjectByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error) {
	subject := new(ent.Subject)
	if err := s.db.NewSelect().Model(subject).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *Service) UpdateSubjectByID(ctx context.Context, id uuid.UUID, subject *ent.Subject) (*ent.Subject, error) {
	updated := new(ent.Subject)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", subject.SchoolID).
		Set("subject_code = ?", subject.SubjectCode).
		Set("name = ?", subject.Name).
		Set("credits = ?", subject.Credits).
		Set("type = ?", subject.Type).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteSubjectByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Subject)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSubjects(ctx context.Context, schoolID *uuid.UUID) ([]*ent.Subject, error) {
	var subjects []*ent.Subject
	query := s.db.NewSelect().Model(&subjects).Order("name ASC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return subjects, nil
}
