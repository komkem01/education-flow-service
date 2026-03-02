package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SchoolAnnouncementEntity = (*Service)(nil)

func (s *Service) CreateSchoolAnnouncement(ctx context.Context, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	if _, err := s.db.NewInsert().Model(announcement).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return announcement, nil
}

func (s *Service) GetSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error) {
	announcement := new(ent.SchoolAnnouncement)
	if err := s.db.NewSelect().Model(announcement).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return announcement, nil
}

func (s *Service) UpdateSchoolAnnouncementByID(ctx context.Context, id uuid.UUID, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	updated := new(ent.SchoolAnnouncement)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", announcement.SchoolID).
		Set("author_member_id = ?", announcement.AuthorMemberID).
		Set("title = ?", announcement.Title).
		Set("content = ?", announcement.Content).
		Set("target_role = ?", announcement.TargetRole).
		Set("is_pinned = ?", announcement.IsPinned).
		Set("updated_at = current_timestamp").
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SchoolAnnouncement)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSchoolAnnouncements(ctx context.Context, schoolID *uuid.UUID, targetRole *ent.MemberRole, onlyPinned bool) ([]*ent.SchoolAnnouncement, error) {
	var items []*ent.SchoolAnnouncement
	query := s.db.NewSelect().Model(&items).Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if targetRole != nil {
		query = query.Where("target_role = ?", *targetRole)
	}
	if onlyPinned {
		query = query.Where("is_pinned = true")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
