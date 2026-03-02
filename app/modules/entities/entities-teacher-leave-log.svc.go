package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherLeaveLogEntity = (*Service)(nil)

func (s *Service) CreateTeacherLeaveLog(ctx context.Context, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error) {
	if _, err := s.db.NewInsert().Model(leaveLog).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return leaveLog, nil
}

func (s *Service) UpdateTeacherLeaveLogByID(ctx context.Context, id uuid.UUID, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error) {
	updated := new(ent.TeacherLeaveLog)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("type = ?", leaveLog.Type).
		Set("start_date = ?", leaveLog.StartDate).
		Set("end_date = ?", leaveLog.EndDate).
		Set("reason = ?", leaveLog.Reason).
		Set("status = ?", leaveLog.Status).
		Set("approved_by_staff_id = ?", leaveLog.ApprovedByStaffID).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) ListTeacherLeaveLogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherLeaveLog, error) {
	var logs []*ent.TeacherLeaveLog
	if err := s.db.NewSelect().Model(&logs).Where("teacher_id = ?", teacherID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return logs, nil
}
