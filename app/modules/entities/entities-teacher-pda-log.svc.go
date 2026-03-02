package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherPDALogEntity = (*Service)(nil)

func (s *Service) CreateTeacherPDALog(ctx context.Context, pdaLog *ent.TeacherPDALog) (*ent.TeacherPDALog, error) {
	if _, err := s.db.NewInsert().Model(pdaLog).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return pdaLog, nil
}

func (s *Service) DeleteTeacherPDALogByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.TeacherPDALog)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListTeacherPDALogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPDALog, error) {
	var logs []*ent.TeacherPDALog
	if err := s.db.NewSelect().Model(&logs).Where("teacher_id = ?", teacherID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *Service) TeacherPDALogBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.TeacherPDALog)(nil)).
		Where("id = ?", id).
		Where("teacher_id = ?", teacherID).
		Exists(ctx)
}
