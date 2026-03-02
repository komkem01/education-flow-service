package ent

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectAssignment struct {
	bun.BaseModel `bun:"table:subject_assignments,alias:sas"`

	ID             uuid.UUID `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectID      uuid.UUID `bun:"subject_id,type:uuid,notnull"`
	TeacherID      uuid.UUID `bun:"teacher_id,type:uuid,notnull"`
	ClassroomID    uuid.UUID `bun:"classroom_id,type:uuid,notnull"`
	AcademicYearID uuid.UUID `bun:"academic_year_id,type:uuid,notnull"`
}
