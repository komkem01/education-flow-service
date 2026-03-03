package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherWorkExperienceEntity = (*Service)(nil)

func (s *Service) CreateTeacherWorkExperience(ctx context.Context, work *ent.TeacherWorkExperience) (*ent.TeacherWorkExperience, error) {
	if _, err := s.db.NewInsert().Model(work).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return work, nil
}

func (s *Service) UpdateTeacherWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.TeacherWorkExperience) (*ent.TeacherWorkExperience, error) {
	updated := new(ent.TeacherWorkExperience)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("organization = ?", work.Organization).
		Set("position = ?", work.Position).
		Set("start_date = ?", work.StartDate).
		Set("end_date = ?", work.EndDate).
		Set("is_current = ?", work.IsCurrent).
		Set("description = ?", work.Description).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteTeacherWorkExperienceByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.TeacherWorkExperience)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListTeacherWorkExperiencesByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherWorkExperience, error) {
	var works []*ent.TeacherWorkExperience
	if err := s.db.NewSelect().Model(&works).Where("teacher_id = ?", teacherID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return works, nil
}

func (s *Service) TeacherWorkExperienceBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.TeacherWorkExperience)(nil)).
		Where("id = ?", id).
		Where("teacher_id = ?", teacherID).
		Exists(ctx)
}
