package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SubjectAssignment struct {
	bun.BaseModel `bun:"table:subject_assignments,alias:sas"`

	ID             uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SubjectID      uuid.UUID  `bun:"subject_id,type:uuid,notnull"`
	TeacherID      uuid.UUID  `bun:"teacher_id,type:uuid,notnull"`
	ClassroomID    uuid.UUID  `bun:"classroom_id,type:uuid,notnull"`
	AcademicYearID uuid.UUID  `bun:"academic_year_id,type:uuid,notnull"`
	Section        *string    `bun:"section"`
	SemesterNo     int        `bun:"semester_no,notnull,default:1"`
	MaxStudents    *int       `bun:"max_students"`
	StartDate      *time.Time `bun:"start_date,type:date"`
	EndDate        *time.Time `bun:"end_date,type:date"`
	Note           *string    `bun:"note"`
	IsActive       bool       `bun:"is_active,notnull,default:true"`
}
