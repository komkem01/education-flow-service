package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentAttendanceLogEntity = (*Service)(nil)

func (s *Service) CreateStudentAttendanceLog(ctx context.Context, attendanceLog *ent.StudentAttendanceLog) (*ent.StudentAttendanceLog, error) {
	if _, err := s.db.NewInsert().Model(attendanceLog).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return attendanceLog, nil
}

func (s *Service) UpdateStudentAttendanceLogByID(ctx context.Context, id uuid.UUID, attendanceLog *ent.StudentAttendanceLog) (*ent.StudentAttendanceLog, error) {
	updated := new(ent.StudentAttendanceLog)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("enrollment_id = ?", attendanceLog.EnrollmentID).
		Set("schedule_id = ?", attendanceLog.ScheduleID).
		Set("check_date = ?", attendanceLog.CheckDate).
		Set("status = ?", attendanceLog.Status).
		Set("note = ?", attendanceLog.Note).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStudentAttendanceLogByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StudentAttendanceLog)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudentAttendanceLogsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentAttendanceLog, error) {
	enrollmentIDs := s.db.NewSelect().Model((*ent.StudentEnrollment)(nil)).Column("id").Where("student_id = ?", studentID)

	var logs []*ent.StudentAttendanceLog
	if err := s.db.NewSelect().
		Model(&logs).
		Where("enrollment_id IN (?)", enrollmentIDs).
		Order("created_at DESC").
		Scan(ctx); err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *Service) StudentAttendanceLogBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentAttendanceLog)(nil)).
		Join("JOIN student_enrollments AS sen ON sen.id = sal.enrollment_id").
		Where("sal.id = ?", id).
		Where("sen.student_id = ?", studentID).
		Exists(ctx)
}

func (s *Service) EnrollmentBelongsToStudent(ctx context.Context, enrollmentID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentEnrollment)(nil)).
		Where("sen.id = ?", enrollmentID).
		Where("sen.student_id = ?", studentID).
		Exists(ctx)
}
