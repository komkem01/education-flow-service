package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.TeacherProfileRequestEntity = (*Service)(nil)

func (s *Service) CreateTeacherProfileRequest(ctx context.Context, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error) {
	if _, err := s.db.NewInsert().Model(profileRequest).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return profileRequest, nil
}

func (s *Service) UpdateTeacherProfileRequestByID(ctx context.Context, id uuid.UUID, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error) {
	updated := new(ent.TeacherProfileRequest)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("requested_data = ?", profileRequest.RequestedData).
		Set("reason = ?", profileRequest.Reason).
		Set("status = ?", profileRequest.Status).
		Set("comment = ?", profileRequest.Comment).
		Set("processed_by_staff_id = ?", profileRequest.ProcessedByStaffID).
		Set("processed_at = ?", profileRequest.ProcessedAt).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) ListTeacherProfileRequestsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherProfileRequest, error) {
	var requests []*ent.TeacherProfileRequest
	if err := s.db.NewSelect().Model(&requests).Where("teacher_id = ?", teacherID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return requests, nil
}

func (s *Service) TeacherProfileRequestBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.TeacherProfileRequest)(nil)).
		Where("id = ?", id).
		Where("teacher_id = ?", teacherID).
		Exists(ctx)
}
