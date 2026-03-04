package teacherprofilerequests

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type mockTeacherProfileRequestDB struct {
	belongsResult bool
	updated       *ent.TeacherProfileRequest
}

func (m *mockTeacherProfileRequestDB) CreateTeacherProfileRequest(ctx context.Context, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error) {
	return profileRequest, nil
}
func (m *mockTeacherProfileRequestDB) UpdateTeacherProfileRequestByID(ctx context.Context, id uuid.UUID, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error) {
	m.updated = profileRequest
	profileRequest.ID = id
	return profileRequest, nil
}
func (m *mockTeacherProfileRequestDB) ListTeacherProfileRequestsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherProfileRequest, error) {
	return nil, nil
}
func (m *mockTeacherProfileRequestDB) TeacherProfileRequestBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error) {
	return m.belongsResult, nil
}

func TestPhase3UpdateByIDReturnsNotFoundWhenNotBelongingTeacher(t *testing.T) {
	svc := newService(&Options{db: &mockTeacherProfileRequestDB{belongsResult: false}})
	_, err := svc.UpdateByID(context.Background(), uuid.New(), uuid.New(), &UpdateInput{Status: ent.TeacherProfileRequestStatusPending})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

func TestPhase3UpdateByIDAllowsApproveFlowForOwner(t *testing.T) {
	db := &mockTeacherProfileRequestDB{belongsResult: true}
	svc := newService(&Options{db: db})
	now := time.Now().UTC()
	staffID := uuid.New()
	comment := " approved "
	updated, err := svc.UpdateByID(context.Background(), uuid.New(), uuid.New(), &UpdateInput{
		RequestedData:      map[string]any{"phone": "0900000000"},
		Reason:             nil,
		Status:             ent.TeacherProfileRequestStatusApproved,
		Comment:            &comment,
		ProcessedByStaffID: &staffID,
		ProcessedAt:        &now,
	})
	if err != nil {
		t.Fatalf("UpdateByID returned error: %v", err)
	}
	if updated.Status != ent.TeacherProfileRequestStatusApproved {
		t.Fatalf("expected approved status, got %s", updated.Status)
	}
	if db.updated == nil || db.updated.Comment == nil || *db.updated.Comment != "approved" {
		t.Fatalf("expected trimmed approval comment, got %+v", db.updated)
	}
}
