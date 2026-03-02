package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherEducationEntity = (*Service)(nil)

func (s *Service) CreateTeacherEducation(ctx context.Context, education *ent.TeacherEducation) (*ent.TeacherEducation, error) {
	if _, err := s.db.NewInsert().Model(education).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return education, nil
}

func (s *Service) UpdateTeacherEducationByID(ctx context.Context, id uuid.UUID, education *ent.TeacherEducation) (*ent.TeacherEducation, error) {
	updated := new(ent.TeacherEducation)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("degree_level = ?", education.DegreeLevel).
		Set("degree_name = ?", education.DegreeName).
		Set("major = ?", education.Major).
		Set("university = ?", education.University).
		Set("graduation_year = ?", education.GraduationYear).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteTeacherEducationByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.TeacherEducation)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListTeacherEducationsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherEducation, error) {
	var educations []*ent.TeacherEducation
	if err := s.db.NewSelect().Model(&educations).Where("teacher_id = ?", teacherID).Scan(ctx); err != nil {
		return nil, err
	}

	return educations, nil
}
