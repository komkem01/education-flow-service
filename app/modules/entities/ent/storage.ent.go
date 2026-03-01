package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StorageVisibility string

type StorageStatus string

type StorageVirusScanStatus string

const (
	StorageVisibilityPrivate StorageVisibility = "private"
	StorageVisibilityPublic  StorageVisibility = "public"
	StorageVisibilitySigned  StorageVisibility = "signed"
)

const (
	StorageStatusPending  StorageStatus = "pending"
	StorageStatusActive   StorageStatus = "active"
	StorageStatusObsolete StorageStatus = "obsolete"
	StorageStatusDeleted  StorageStatus = "deleted"
)

const (
	StorageVirusScanStatusPending  StorageVirusScanStatus = "pending"
	StorageVirusScanStatusClean    StorageVirusScanStatus = "clean"
	StorageVirusScanStatusInfected StorageVirusScanStatus = "infected"
	StorageVirusScanStatusFailed   StorageVirusScanStatus = "failed"
)

type Storage struct {
	bun.BaseModel `bun:"table:storages,alias:stg"`

	ID                  uuid.UUID              `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID            uuid.UUID              `bun:"school_id,type:uuid,notnull"`
	BucketName          string                 `bun:"bucket_name,notnull"`
	ObjectKey           string                 `bun:"object_key,notnull"`
	OriginalName        *string                `bun:"original_name"`
	Extension           *string                `bun:"extension"`
	MIMEType            *string                `bun:"mime_type"`
	SizeBytes           int64                  `bun:"size_bytes,notnull"`
	ChecksumSHA256      *string                `bun:"checksum_sha256"`
	ETag                *string                `bun:"etag"`
	Visibility          StorageVisibility      `bun:"visibility,notnull"`
	Status              StorageStatus          `bun:"status,notnull"`
	VirusScanStatus     StorageVirusScanStatus `bun:"virus_scan_status,notnull"`
	VirusScanAt         *time.Time             `bun:"virus_scan_at"`
	UploadedByMemberID  *uuid.UUID             `bun:"uploaded_by_member_id,type:uuid"`
	VersionNo           int                    `bun:"version_no,notnull"`
	ReplacedByStorageID *uuid.UUID             `bun:"replaced_by_storage_id,type:uuid"`
	Metadata            map[string]any         `bun:"metadata,type:jsonb"`
	CreatedAt           time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt           time.Time              `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt           *time.Time             `bun:"deleted_at,soft_delete"`
}

type StorageLink struct {
	bun.BaseModel `bun:"table:storage_links,alias:slk"`

	ID         uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StorageID  uuid.UUID `bun:"storage_id,type:uuid,notnull"`
	EntityType string    `bun:"entity_type,notnull"`
	EntityID   uuid.UUID `bun:"entity_id,type:uuid,notnull"`
	FieldName  *string   `bun:"field_name"`
	SortOrder  int       `bun:"sort_order,notnull"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToStorageVisibility(value string) StorageVisibility {
	switch value {
	case "public":
		return StorageVisibilityPublic
	case "signed":
		return StorageVisibilitySigned
	default:
		return StorageVisibilityPrivate
	}
}

func ToStorageStatus(value string) StorageStatus {
	switch value {
	case "active":
		return StorageStatusActive
	case "obsolete":
		return StorageStatusObsolete
	case "deleted":
		return StorageStatusDeleted
	default:
		return StorageStatusPending
	}
}

func ToStorageVirusScanStatus(value string) StorageVirusScanStatus {
	switch value {
	case "clean":
		return StorageVirusScanStatusClean
	case "infected":
		return StorageVirusScanStatusInfected
	case "failed":
		return StorageVirusScanStatusFailed
	default:
		return StorageVirusScanStatusPending
	}
}
