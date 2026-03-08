package studentgraderecords

import (
	"context"
	"database/sql"
	"testing"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

type mockGradeRecordDB struct {
	enrollmentBelongs bool
	gradeItemBelongs  bool
	created           *ent.GradeRecord
}

func (m *mockGradeRecordDB) CreateGradeRecord(ctx context.Context, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error) {
	m.created = gradeRecord
	return gradeRecord, nil
}
func (m *mockGradeRecordDB) UpdateGradeRecordByID(ctx context.Context, id uuid.UUID, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error) {
	return gradeRecord, nil
}
func (m *mockGradeRecordDB) DeleteGradeRecordByID(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockGradeRecordDB) ListGradeRecordsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeRecord, error) {
	return nil, nil
}
func (m *mockGradeRecordDB) GradeRecordBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return true, nil
}
func (m *mockGradeRecordDB) EnrollmentBelongsToStudent(ctx context.Context, enrollmentID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return m.enrollmentBelongs, nil
}
func (m *mockGradeRecordDB) GradeItemBelongsToStudent(ctx context.Context, gradeItemID uuid.UUID, studentID uuid.UUID) (bool, error) {
	return m.gradeItemBelongs, nil
}
func (m *mockGradeRecordDB) GetStudentByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error) {
	return &ent.MemberStudent{MemberID: uuid.New()}, nil
}
func (m *mockGradeRecordDB) GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error) {
	return &ent.Member{SchoolID: uuid.New()}, nil
}

func TestPhase4CreateRejectsMismatchedEnrollment(t *testing.T) {
	svc := newService(&Options{db: &mockGradeRecordDB{enrollmentBelongs: false, gradeItemBelongs: true}})
	_, err := svc.Create(context.Background(), &CreateInput{StudentID: uuid.New(), EnrollmentID: uuid.New(), GradeItemID: uuid.New()})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows for mismatched enrollment, got %v", err)
	}
}

func TestPhase4CreateRejectsMismatchedGradeItem(t *testing.T) {
	svc := newService(&Options{db: &mockGradeRecordDB{enrollmentBelongs: true, gradeItemBelongs: false}})
	_, err := svc.Create(context.Background(), &CreateInput{StudentID: uuid.New(), EnrollmentID: uuid.New(), GradeItemID: uuid.New()})
	if err != sql.ErrNoRows {
		t.Fatalf("expected sql.ErrNoRows for mismatched grade item, got %v", err)
	}
}

func TestPhase4CreateSuccessWhenOwnershipValid(t *testing.T) {
	note := "  good progress  "
	db := &mockGradeRecordDB{enrollmentBelongs: true, gradeItemBelongs: true}
	svc := newService(&Options{db: db})
	item, err := svc.Create(context.Background(), &CreateInput{StudentID: uuid.New(), EnrollmentID: uuid.New(), GradeItemID: uuid.New(), TeacherNote: &note})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if item.TeacherNote == nil || *item.TeacherNote != "good progress" {
		t.Fatalf("expected trimmed teacher note, got %+v", item.TeacherNote)
	}
}
