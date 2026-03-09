package schoolannouncements

import (
	"context"
	"database/sql"
	"testing"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type mockSchoolAnnouncementDB struct {
	schoolExists bool
	memberExists bool
	isAdmin      bool
	isStaff      bool
}

func (m *mockSchoolAnnouncementDB) CreateSchoolAnnouncement(ctx context.Context, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	return announcement, nil
}
func (m *mockSchoolAnnouncementDB) GetSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error) {
	return nil, sql.ErrNoRows
}
func (m *mockSchoolAnnouncementDB) UpdateSchoolAnnouncementByID(ctx context.Context, id uuid.UUID, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error) {
	return announcement, nil
}
func (m *mockSchoolAnnouncementDB) DeleteSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) ListSchoolAnnouncements(ctx context.Context, schoolID *uuid.UUID, targetRole *ent.MemberRole, onlyPinned bool) ([]*ent.SchoolAnnouncement, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) CreateSchool(ctx context.Context, school *ent.School) (*ent.School, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	if !m.schoolExists {
		return nil, sql.ErrNoRows
	}
	return &ent.School{ID: id}, nil
}
func (m *mockSchoolAnnouncementDB) UpdateSchoolByID(ctx context.Context, id uuid.UUID, school *ent.School) (*ent.School, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) DeleteSchoolByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) ListSchools(ctx context.Context) ([]*ent.School, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	if !m.memberExists {
		return nil, sql.ErrNoRows
	}
	return &ent.Member{ID: id}, nil
}
func (m *mockSchoolAnnouncementDB) GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error) {
	return nil, sql.ErrNoRows
}
func (m *mockSchoolAnnouncementDB) UpdateMemberByID(ctx context.Context, id uuid.UUID, member *ent.Member) (*ent.Member, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) UpdateMemberLastLoginByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) CreateAdmin(ctx context.Context, admin *ent.MemberAdmin) (*ent.MemberAdmin, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) GetAdminByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID) (*ent.MemberAdmin, error) {
	return nil, sql.ErrNoRows
}
func (m *mockSchoolAnnouncementDB) UpdateAdminByID(ctx context.Context, id uuid.UUID, admin *ent.MemberAdmin) (*ent.MemberAdmin, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) DeleteAdminByID(ctx context.Context, id uuid.UUID, schoolID *uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) ListAdmins(ctx context.Context, schoolID *uuid.UUID, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberAdmin, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) MemberHasAdminRole(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return m.isAdmin, nil
}
func (m *mockSchoolAnnouncementDB) CreateStaff(ctx context.Context, staff *ent.MemberStaff) (*ent.MemberStaff, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) GetStaffByID(ctx context.Context, id uuid.UUID) (*ent.MemberStaff, error) {
	return nil, sql.ErrNoRows
}
func (m *mockSchoolAnnouncementDB) UpdateStaffByID(ctx context.Context, id uuid.UUID, staff *ent.MemberStaff) (*ent.MemberStaff, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) DeleteStaffByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockSchoolAnnouncementDB) ListStaffs(ctx context.Context, schoolID *uuid.UUID, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberStaff, error) {
	return nil, nil
}
func (m *mockSchoolAnnouncementDB) MemberHasStaffRole(ctx context.Context, memberID uuid.UUID) (bool, error) {
	return m.isStaff, nil
}

func TestPhase5NotificationRejectsUnknownSchool(t *testing.T) {
	svc := newService(&Options{db: &mockSchoolAnnouncementDB{schoolExists: false, memberExists: true, isAdmin: true}})
	_, err := svc.Create(context.Background(), &CreateInput{SchoolID: uuid.New(), AuthorMemberID: uuid.New()})
	if err != ErrSchoolNotFound {
		t.Fatalf("expected ErrSchoolNotFound, got %v", err)
	}
}

func TestPhase5NotificationRejectsAuthorWithoutAdminOrStaffRole(t *testing.T) {
	svc := newService(&Options{db: &mockSchoolAnnouncementDB{schoolExists: true, memberExists: true, isAdmin: false, isStaff: false}})
	_, err := svc.Create(context.Background(), &CreateInput{SchoolID: uuid.New(), AuthorMemberID: uuid.New()})
	if err != ErrInvalidAuthorRole {
		t.Fatalf("expected ErrInvalidAuthorRole, got %v", err)
	}
}
