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

type AcademicYearEntity interface {
	CreateAcademicYear(ctx context.Context, academicYear *ent.AcademicYear) (*ent.AcademicYear, error)
	GetAcademicYearByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error)
	UpdateAcademicYearByID(ctx context.Context, id uuid.UUID, academicYear *ent.AcademicYear) (*ent.AcademicYear, error)
	DeleteAcademicYearByID(ctx context.Context, id uuid.UUID) error
	ListAcademicYears(ctx context.Context, onlyActive bool, onlyCurrent bool) ([]*ent.AcademicYear, error)
}

type MemberEntity interface {
	CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error)
	UpdateMemberByID(ctx context.Context, id uuid.UUID, member *ent.Member) (*ent.Member, error)
	DeleteMemberByID(ctx context.Context, id uuid.UUID) error
	ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error)
}

type MemberTeacherEntity interface {
	CreateTeacher(ctx context.Context, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error)
	GetTeacherByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error)
	UpdateTeacherByID(ctx context.Context, id uuid.UUID, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error)
	DeleteTeacherByID(ctx context.Context, id uuid.UUID) error
	ListTeachers(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberTeacher, error)
}

type TeacherEducationEntity interface {
	CreateTeacherEducation(ctx context.Context, education *ent.TeacherEducation) (*ent.TeacherEducation, error)
	UpdateTeacherEducationByID(ctx context.Context, id uuid.UUID, education *ent.TeacherEducation) (*ent.TeacherEducation, error)
	DeleteTeacherEducationByID(ctx context.Context, id uuid.UUID) error
	ListTeacherEducationsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherEducation, error)
}

type TeacherProfileRequestEntity interface {
	CreateTeacherProfileRequest(ctx context.Context, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error)
	UpdateTeacherProfileRequestByID(ctx context.Context, id uuid.UUID, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error)
	ListTeacherProfileRequestsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherProfileRequest, error)
}

type TeacherPerformanceAgreementEntity interface {
	CreateTeacherPerformanceAgreement(ctx context.Context, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error)
	UpdateTeacherPerformanceAgreementByID(ctx context.Context, id uuid.UUID, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error)
	ListTeacherPerformanceAgreementsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPerformanceAgreement, error)
}

type TeacherPDALogEntity interface {
	CreateTeacherPDALog(ctx context.Context, pdaLog *ent.TeacherPDALog) (*ent.TeacherPDALog, error)
	DeleteTeacherPDALogByID(ctx context.Context, id uuid.UUID) error
	ListTeacherPDALogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPDALog, error)
}

type TeacherLeaveLogEntity interface {
	CreateTeacherLeaveLog(ctx context.Context, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error)
	UpdateTeacherLeaveLogByID(ctx context.Context, id uuid.UUID, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error)
	ListTeacherLeaveLogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherLeaveLog, error)
}
