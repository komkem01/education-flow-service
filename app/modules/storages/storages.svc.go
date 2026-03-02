package storages

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
	ErrSchoolNotFound            = errors.New("school not found")
	ErrUploadedByMemberNotFound  = errors.New("uploaded by member not found")
	ErrReplacedByStorageNotFound = errors.New("replaced by storage not found")
)

type serviceDB interface {
	entitiesinf.StorageEntity
	entitiesinf.SchoolEntity
	entitiesinf.MemberEntity
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
	SchoolID            uuid.UUID
	BucketName          string
	ObjectKey           string
	OriginalName        *string
	Extension           *string
	MIMEType            *string
	SizeBytes           int64
	ChecksumSHA256      *string
	ETag                *string
	Visibility          ent.StorageVisibility
	Status              ent.StorageStatus
	VirusScanStatus     ent.StorageVirusScanStatus
	UploadedByMemberID  *uuid.UUID
	VersionNo           int
	ReplacedByStorageID *uuid.UUID
	Metadata            map[string]any
}

type UpdateInput = CreateInput

type ListInput struct {
	SchoolID           *uuid.UUID
	UploadedByMemberID *uuid.UUID
	Status             *ent.StorageStatus
	Visibility         *ent.StorageVisibility
}

func newService(opt *Options) *Service { return &Service{tracer: opt.tracer, db: opt.db} }

func (s *Service) Create(ctx context.Context, input *CreateInput) (*ent.Storage, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.UploadedByMemberID, input.ReplacedByStorageID); err != nil {
		return nil, err
	}

	item := &ent.Storage{SchoolID: input.SchoolID, BucketName: strings.TrimSpace(input.BucketName), ObjectKey: strings.TrimSpace(input.ObjectKey), OriginalName: trimStringPtr(input.OriginalName), Extension: trimStringPtr(input.Extension), MIMEType: trimStringPtr(input.MIMEType), SizeBytes: input.SizeBytes, ChecksumSHA256: trimStringPtr(input.ChecksumSHA256), ETag: trimStringPtr(input.ETag), Visibility: input.Visibility, Status: input.Status, VirusScanStatus: input.VirusScanStatus, UploadedByMemberID: input.UploadedByMemberID, VersionNo: input.VersionNo, ReplacedByStorageID: input.ReplacedByStorageID, Metadata: input.Metadata}
	return s.db.CreateStorage(ctx, item)
}

func (s *Service) List(ctx context.Context, input *ListInput) ([]*ent.Storage, error) {
	return s.db.ListStorages(ctx, input.SchoolID, input.UploadedByMemberID, input.Status, input.Visibility)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*ent.Storage, error) {
	return s.db.GetStorageByID(ctx, id)
}

func (s *Service) UpdateByID(ctx context.Context, id uuid.UUID, input *UpdateInput) (*ent.Storage, error) {
	if err := s.validateDependencies(ctx, input.SchoolID, input.UploadedByMemberID, input.ReplacedByStorageID); err != nil {
		return nil, err
	}

	item := &ent.Storage{SchoolID: input.SchoolID, BucketName: strings.TrimSpace(input.BucketName), ObjectKey: strings.TrimSpace(input.ObjectKey), OriginalName: trimStringPtr(input.OriginalName), Extension: trimStringPtr(input.Extension), MIMEType: trimStringPtr(input.MIMEType), SizeBytes: input.SizeBytes, ChecksumSHA256: trimStringPtr(input.ChecksumSHA256), ETag: trimStringPtr(input.ETag), Visibility: input.Visibility, Status: input.Status, VirusScanStatus: input.VirusScanStatus, UploadedByMemberID: input.UploadedByMemberID, VersionNo: input.VersionNo, ReplacedByStorageID: input.ReplacedByStorageID, Metadata: input.Metadata}
	return s.db.UpdateStorageByID(ctx, id, item)
}

func (s *Service) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.db.SoftDeleteStorageByID(ctx, id)
}

func (s *Service) validateDependencies(ctx context.Context, schoolID uuid.UUID, uploadedByMemberID *uuid.UUID, replacedByStorageID *uuid.UUID) error {
	if _, err := s.db.GetSchoolByID(ctx, schoolID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSchoolNotFound
		}
		return err
	}
	if uploadedByMemberID != nil {
		if _, err := s.db.GetMemberByID(ctx, *uploadedByMemberID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrUploadedByMemberNotFound
			}
			return err
		}
	}
	if replacedByStorageID != nil {
		if _, err := s.db.GetStorageByID(ctx, *replacedByStorageID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrReplacedByStorageNotFound
			}
			return err
		}
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
