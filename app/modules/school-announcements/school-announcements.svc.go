package schoolannouncements

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrSchoolNotFound       = errors.New("school not found")
	ErrAuthorMemberNotFound = errors.New("author member not found")
	ErrInvalidAuthorRole    = errors.New("invalid author role")
)

type serviceDB interface {
	entitiesinf.SchoolAnnouncementEntity
	entitiesinf.SchoolEntity
	entitiesinf.MemberEntity
	entitiesinf.MemberAdminEntity
	entitiesinf.MemberStaffEntity
}

type Service struct {
	tracer trace.Tracer
	db     serviceDB
}

type Config struct{}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     serviceDB
}

type CreateInput struct {
	SchoolID       uuid.UUID
	AuthorMemberID uuid.UUID
	Title          *string
	Content        *string
	TargetRole     *ent.MemberRole
	IsPinned       bool
}

type UpdateInput = CreateInput

type ListInput struct {
	SchoolID   *uuid.UUID
	TargetRole *ent.MemberRole
	OnlyPinned bool
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.SchoolAnnouncement, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.AuthorMemberID); err != nil {
		return nil, err
	}

	item := &ent.SchoolAnnouncement{SchoolID: input.SchoolID, AuthorMemberID: input.AuthorMemberID, Title: trimStringPtr(input.Title), Content: trimStringPtr(input.Content), TargetRole: input.TargetRole, IsPinned: input.IsPinned}
	return s.db.CreateSchoolAnnouncement(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInput) ([]*ent.SchoolAnnouncement, error) {
	return s.db.ListSchoolAnnouncements(ctx, input.SchoolID, input.TargetRole, input.OnlyPinned)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error) {
	return s.db.GetSchoolAnnouncementByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.SchoolAnnouncement, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.AuthorMemberID); err != nil {
		return nil, err
	}

	item := &ent.SchoolAnnouncement{SchoolID: input.SchoolID, AuthorMemberID: input.AuthorMemberID, Title: trimStringPtr(input.Title), Content: trimStringPtr(input.Content), TargetRole: input.TargetRole, IsPinned: input.IsPinned}
	return s.db.UpdateSchoolAnnouncementByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.DeleteSchoolAnnouncementByID(ctx, id)
}

func (s *Service) validateDependencies(ctx context.Context, schoolID uuid.UUID, authorMemberID uuid.UUID) error {
	if _, err := s.db.GetSchoolByID(ctx, schoolID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSchoolNotFound
		}
		return err
	}
	if _, err := s.db.GetMemberByID(ctx, authorMemberID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAuthorMemberNotFound
		}
		return err
	}

	isAdmin, err := s.db.MemberHasAdminRole(ctx, authorMemberID)
	if err != nil {
		return err
	}
	isStaff, err := s.db.MemberHasStaffRole(ctx, authorMemberID)
	if err != nil {
		return err
	}
	if !isAdmin && !isStaff {
		return ErrInvalidAuthorRole
	}
	return nil
}

func trimStringPtr(input *string) *string {
	if input == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*input)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
