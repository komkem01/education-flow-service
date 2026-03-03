package schedules

import (
	"context"
	"errors"
	"time"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     entitiesinf.ScheduleEntity
	saDB   entitiesinf.SubjectAssignmentEntity
}

var (
	ErrScheduleTeacherConflict   = errors.New("schedule-teacher-conflict")
	ErrScheduleClassroomConflict = errors.New("schedule-classroom-conflict")
)

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     entitiesinf.ScheduleEntity
	saDB   entitiesinf.SubjectAssignmentEntity
}

type CreateScheduleInput struct {
	SubjectAssignmentID uuid.UUID
	DayOfWeek           ent.ScheduleDayOfWeek
	StartTime           *time.Time
	EndTime             *time.Time
	PeriodNo            *int
	Note                *string
	IsActive            bool
}

type UpdateScheduleInput = CreateScheduleInput

type ListSchedulesInput struct {
	SubjectAssignmentID *uuid.UUID
	DayOfWeek           *ent.ScheduleDayOfWeek
	OnlyActive          *bool
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db, saDB: opt.saDB}
}

func (s *Service) Create(ctx context.Context, input *CreateScheduleInput) (*ent.Schedule, error) {
	if err := s.validateConflicts(ctx, nil, input); err != nil {
		return nil, err
	}

	item := &ent.Schedule{SubjectAssignmentID: input.SubjectAssignmentID, DayOfWeek: input.DayOfWeek, StartTime: input.StartTime, EndTime: input.EndTime, PeriodNo: input.PeriodNo, Note: input.Note, IsActive: input.IsActive}
	return s.db.CreateSchedule(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListSchedulesInput) ([]*ent.Schedule, error) {
	items, err := s.db.ListSchedules(ctx, input.SubjectAssignmentID, input.DayOfWeek)
	if err != nil {
		return nil, err
	}

	if input.OnlyActive == nil {
		return items, nil
	}

	filtered := make([]*ent.Schedule, 0, len(items))
	for _, item := range items {
		if item.IsActive != *input.OnlyActive {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Schedule, error) {
	return s.db.GetScheduleByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateScheduleInput) (*ent.Schedule, error) {
	if err := s.validateConflicts(ctx, &id, input); err != nil {
		return nil, err
	}

	item := &ent.Schedule{SubjectAssignmentID: input.SubjectAssignmentID, DayOfWeek: input.DayOfWeek, StartTime: input.StartTime, EndTime: input.EndTime, PeriodNo: input.PeriodNo, Note: input.Note, IsActive: input.IsActive}
	return s.db.UpdateScheduleByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteScheduleByID(ctx, id)
}

func (s *Service) validateConflicts(ctx context.Context, ignoreScheduleID *uuid.UUID, input *CreateScheduleInput) error {
	if !input.IsActive {
		return nil
	}

	assignment, err := s.saDB.GetSubjectAssignmentByID(ctx, input.SubjectAssignmentID)
	if err != nil {
		return err
	}

	if !assignment.IsActive {
		return nil
	}

	teacherAssignments, err := s.saDB.ListSubjectAssignments(ctx, nil, &assignment.TeacherID, nil, &assignment.AcademicYearID)
	if err != nil {
		return err
	}

	if hasConflict, err := s.hasScheduleConflict(ctx, input, assignment, teacherAssignments, ignoreScheduleID); err != nil {
		return err
	} else if hasConflict {
		return ErrScheduleTeacherConflict
	}

	classroomAssignments, err := s.saDB.ListSubjectAssignments(ctx, nil, nil, &assignment.ClassroomID, &assignment.AcademicYearID)
	if err != nil {
		return err
	}

	if hasConflict, err := s.hasScheduleConflict(ctx, input, assignment, classroomAssignments, ignoreScheduleID); err != nil {
		return err
	} else if hasConflict {
		return ErrScheduleClassroomConflict
	}

	return nil
}

func (s *Service) hasScheduleConflict(ctx context.Context, input *CreateScheduleInput, baseAssignment *ent.SubjectAssignment, assignments []*ent.SubjectAssignment, ignoreScheduleID *uuid.UUID) (bool, error) {
	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}
		if assignment.ID == baseAssignment.ID {
			continue
		}
		if !assignment.IsActive {
			continue
		}
		if assignment.SemesterNo != baseAssignment.SemesterNo {
			continue
		}

		targetID := assignment.ID
		schedules, err := s.db.ListSchedules(ctx, &targetID, &input.DayOfWeek)
		if err != nil {
			return false, err
		}

		for _, schedule := range schedules {
			if schedule == nil || !schedule.IsActive {
				continue
			}
			if ignoreScheduleID != nil && schedule.ID == *ignoreScheduleID {
				continue
			}
			if !isScheduleSlotOverlapped(input.StartTime, input.EndTime, input.PeriodNo, schedule.StartTime, schedule.EndTime, schedule.PeriodNo) {
				continue
			}
			return true, nil
		}
	}

	return false, nil
}

func isScheduleSlotOverlapped(startA *time.Time, endA *time.Time, periodA *int, startB *time.Time, endB *time.Time, periodB *int) bool {
	if periodA != nil && periodB != nil {
		return *periodA == *periodB
	}

	if startA == nil || endA == nil || startB == nil || endB == nil {
		return false
	}

	return startA.Before(*endB) && startB.Before(*endA)
}
