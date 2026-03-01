package entitiesinf

import (
	"context"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

// ObjectEntity defines the interface for object entity operations such as create, retrieve, update, and soft delete.
type ExampleEntity interface {
	CreateExample(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
	GetExampleByID(ctx context.Context, id uuid.UUID) (*ent.Example, error)
	UpdateExampleByID(ctx context.Context, id uuid.UUID, status ent.ExampleStatus) (*ent.Example, error)
	SoftDeleteExampleByID(ctx context.Context, id uuid.UUID) error
	ListExamplesByStatus(ctx context.Context, status ent.ExampleStatus) ([]*ent.Example, error)
}
type ExampleTwoEntity interface {
	CreateExampleTwo(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
}

type StorageEntity interface {
	CreateStorage(ctx context.Context, storage *ent.Storage) (*ent.Storage, error)
	GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.Storage, error)
	GetStorageByObjectKey(ctx context.Context, bucketName, objectKey string) (*ent.Storage, error)
	UpdateStorageStatusByID(ctx context.Context, id uuid.UUID, status ent.StorageStatus) (*ent.Storage, error)
	SoftDeleteStorageByID(ctx context.Context, id uuid.UUID) error
	CreateStorageLink(ctx context.Context, link *ent.StorageLink) (*ent.StorageLink, error)
	ListStorageLinksByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*ent.StorageLink, error)
	DeleteStorageLinksByStorageID(ctx context.Context, storageID uuid.UUID) error
}

type SchoolEntity interface {
	CreateSchool(ctx context.Context, school *ent.School) (*ent.School, error)
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error)
	UpdateSchoolByID(ctx context.Context, id uuid.UUID, school *ent.School) (*ent.School, error)
	DeleteSchoolByID(ctx context.Context, id uuid.UUID) error
	ListSchools(ctx context.Context) ([]*ent.School, error)
}

type GenderEntity interface {
	CreateGender(ctx context.Context, gender *ent.Gender) (*ent.Gender, error)
	GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error)
	UpdateGenderByID(ctx context.Context, id uuid.UUID, gender *ent.Gender) (*ent.Gender, error)
	DeleteGenderByID(ctx context.Context, id uuid.UUID) error
	ListGenders(ctx context.Context, onlyActive bool) ([]*ent.Gender, error)
}

type PrefixEntity interface {
	CreatePrefix(ctx context.Context, prefix *ent.Prefix) (*ent.Prefix, error)
	GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error)
	UpdatePrefixByID(ctx context.Context, id uuid.UUID, prefix *ent.Prefix) (*ent.Prefix, error)
	DeletePrefixByID(ctx context.Context, id uuid.UUID) error
	ListPrefixes(ctx context.Context, onlyActive bool) ([]*ent.Prefix, error)
}
