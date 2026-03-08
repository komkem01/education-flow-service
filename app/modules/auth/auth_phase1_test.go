package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type mockAuthDB struct {
	member *ent.Member
	roles  []ent.MemberRole
}

func (m *mockAuthDB) CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAuthDB) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	if m.member == nil || m.member.ID != id {
		return nil, sql.ErrNoRows
	}
	return m.member, nil
}
func (m *mockAuthDB) GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error) {
	return nil, sql.ErrNoRows
}
func (m *mockAuthDB) UpdateMemberByID(ctx context.Context, id uuid.UUID, member *ent.Member) (*ent.Member, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAuthDB) UpdateMemberLastLoginByID(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockAuthDB) DeleteMemberByID(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}
func (m *mockAuthDB) ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAuthDB) AddMemberRole(ctx context.Context, memberID uuid.UUID, role ent.MemberRole) error {
	return errors.New("not implemented")
}
func (m *mockAuthDB) RemoveMemberRole(ctx context.Context, memberID uuid.UUID, role ent.MemberRole) error {
	return errors.New("not implemented")
}
func (m *mockAuthDB) ListMemberRolesByMemberID(ctx context.Context, memberID uuid.UUID) ([]ent.MemberRole, error) {
	if m.member == nil || m.member.ID != memberID {
		return nil, sql.ErrNoRows
	}
	return append([]ent.MemberRole{}, m.roles...), nil
}
func (m *mockAuthDB) MemberHasAnyRole(ctx context.Context, memberID uuid.UUID, roles []ent.MemberRole) (bool, error) {
	return false, nil
}
func (m *mockAuthDB) CreateSchool(ctx context.Context, school *ent.School) (*ent.School, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAuthDB) GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error) {
	return &ent.School{ID: id, Name: "Test School", Address: "Address"}, nil
}
func (m *mockAuthDB) UpdateSchoolByID(ctx context.Context, id uuid.UUID, school *ent.School) (*ent.School, error) {
	return nil, errors.New("not implemented")
}
func (m *mockAuthDB) DeleteSchoolByID(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}
func (m *mockAuthDB) ListSchools(ctx context.Context) ([]*ent.School, error) {
	return nil, errors.New("not implemented")
}

func TestPhase1SwitchRoleSuccess(t *testing.T) {
	memberID := uuid.New()
	schoolID := uuid.New()
	db := &mockAuthDB{
		member: &ent.Member{ID: memberID, SchoolID: schoolID, IsActive: true},
		roles:  []ent.MemberRole{ent.MemberRoleTeacher, ent.MemberRoleStaff},
	}

	svc := newService(&Options{db: db, appKey: "phase1-test-key"})
	result, err := svc.SwitchRole(context.Background(), &SwitchRoleInput{
		Claims: &TokenClaims{MemberID: memberID, SchoolID: schoolID, Roles: []ent.MemberRole{ent.MemberRoleTeacher, ent.MemberRoleStaff}},
		Role:   ent.MemberRoleStaff,
	})
	if err != nil {
		t.Fatalf("SwitchRole returned error: %v", err)
	}
	if result.Role != ent.MemberRoleStaff {
		t.Fatalf("expected switched role to be staff, got %s", result.Role)
	}
	if len(result.Roles) == 0 || result.Roles[0] != ent.MemberRoleStaff {
		t.Fatalf("expected primary role to be staff after switch, got %+v", result.Roles)
	}
	if result.AccessToken == "" {
		t.Fatalf("expected non-empty access token")
	}
}

func TestPhase1SwitchRoleForbiddenWhenRoleNotOwned(t *testing.T) {
	memberID := uuid.New()
	schoolID := uuid.New()
	db := &mockAuthDB{
		member: &ent.Member{ID: memberID, SchoolID: schoolID, IsActive: true},
		roles:  []ent.MemberRole{ent.MemberRoleTeacher},
	}

	svc := newService(&Options{db: db, appKey: "phase1-test-key"})
	_, err := svc.SwitchRole(context.Background(), &SwitchRoleInput{
		Claims: &TokenClaims{MemberID: memberID, SchoolID: schoolID, Roles: []ent.MemberRole{ent.MemberRoleTeacher}},
		Role:   ent.MemberRoleAdmin,
	})
	if !errors.Is(err, ErrRoleNotAllowed) {
		t.Fatalf("expected ErrRoleNotAllowed, got %v", err)
	}
}

func TestPhase1AccessTokenRoundTripWithRoles(t *testing.T) {
	memberID := uuid.New()
	schoolID := uuid.New()
	member := &ent.Member{ID: memberID, SchoolID: schoolID, IsActive: true}
	roles := []ent.MemberRole{ent.MemberRoleAdmin, ent.MemberRoleTeacher}

	svc := newService(&Options{db: &mockAuthDB{}, appKey: "phase1-test-key"})
	now := time.Now().UTC().Truncate(time.Second)
	expiresAt := now.Add(30 * time.Minute)
	token, err := svc.generateAccessToken(member, roles, now, expiresAt)
	if err != nil {
		t.Fatalf("generateAccessToken error: %v", err)
	}

	claims, err := svc.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken error: %v", err)
	}
	if claims.MemberID != memberID || claims.SchoolID != schoolID {
		t.Fatalf("claims mismatch: member=%s school=%s", claims.MemberID, claims.SchoolID)
	}
	if claims.Role != ent.MemberRoleAdmin {
		t.Fatalf("expected primary admin role, got %s", claims.Role)
	}
	if len(claims.Roles) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(claims.Roles))
	}
}

func TestPhase1SwitchSchoolSuccessForSuperAdmin(t *testing.T) {
	memberID := uuid.New()
	originalSchoolID := uuid.New()
	targetSchoolID := uuid.New()
	db := &mockAuthDB{
		member: &ent.Member{ID: memberID, SchoolID: originalSchoolID, IsActive: true},
		roles:  []ent.MemberRole{ent.MemberRoleSuperAdmin, ent.MemberRoleAdmin},
	}

	svc := newService(&Options{db: db, appKey: "phase1-test-key"})
	result, err := svc.SwitchSchool(context.Background(), &SwitchSchoolInput{
		Claims:   &TokenClaims{MemberID: memberID, SchoolID: originalSchoolID, Role: ent.MemberRoleSuperAdmin, Roles: []ent.MemberRole{ent.MemberRoleSuperAdmin, ent.MemberRoleAdmin}},
		SchoolID: targetSchoolID,
	})
	if err != nil {
		t.Fatalf("SwitchSchool returned error: %v", err)
	}
	if result.SchoolID != targetSchoolID {
		t.Fatalf("expected switched school %s, got %s", targetSchoolID, result.SchoolID)
	}
	if result.AccessToken == "" {
		t.Fatalf("expected non-empty access token")
	}
}
