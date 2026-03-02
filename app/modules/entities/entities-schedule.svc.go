package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.ScheduleEntity = (*Service)(nil)

func (s *Service) CreateSchedule(ctx context.Context, schedule *ent.Schedule) (*ent.Schedule, error) {
	if _, err := s.db.NewInsert().Model(schedule).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *Service) GetScheduleByID(ctx context.Context, id uuid.UUID) (*ent.Schedule, error) {
	schedule := new(ent.Schedule)
	if err := s.db.NewSelect().Model(schedule).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *Service) UpdateScheduleByID(ctx context.Context, id uuid.UUID, schedule *ent.Schedule) (*ent.Schedule, error) {
	updated := new(ent.Schedule)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("subject_assignment_id = ?", schedule.SubjectAssignmentID).
		Set("day_of_week = ?", schedule.DayOfWeek).
		Set("start_time = ?", schedule.StartTime).
		Set("end_time = ?", schedule.EndTime).
		Set("period_no = ?", schedule.PeriodNo).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteScheduleByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.Schedule)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSchedules(ctx context.Context, subjectAssignmentID *uuid.UUID, dayOfWeek *ent.ScheduleDayOfWeek) ([]*ent.Schedule, error) {
	var schedules []*ent.Schedule
	query := s.db.NewSelect().Model(&schedules)

	if subjectAssignmentID != nil {
		query = query.Where("subject_assignment_id = ?", *subjectAssignmentID)
	}
	if dayOfWeek != nil {
		query = query.Where("day_of_week = ?", *dayOfWeek)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return schedules, nil
}
