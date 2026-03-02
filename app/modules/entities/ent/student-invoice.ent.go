package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type StudentInvoiceStatus string

const (
	StudentInvoiceStatusUnpaid    StudentInvoiceStatus = "unpaid"
	StudentInvoiceStatusPaid      StudentInvoiceStatus = "paid"
	StudentInvoiceStatusPartial   StudentInvoiceStatus = "partial"
	StudentInvoiceStatusCancelled StudentInvoiceStatus = "cancelled"
)

type StudentInvoice struct {
	bun.BaseModel `bun:"table:student_invoices,alias:siv"`

	ID             uuid.UUID            `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	StudentID      uuid.UUID            `bun:"student_id,type:uuid,notnull"`
	FeeCategoryID  uuid.UUID            `bun:"fee_category_id,type:uuid,notnull"`
	AcademicYearID uuid.UUID            `bun:"academic_year_id,type:uuid,notnull"`
	Amount         *float64             `bun:"amount"`
	DueDate        *time.Time           `bun:"due_date,type:date"`
	Status         StudentInvoiceStatus `bun:"status"`
	CreatedAt      time.Time            `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToStudentInvoiceStatus(value string) StudentInvoiceStatus {
	switch value {
	case "paid":
		return StudentInvoiceStatusPaid
	case "partial":
		return StudentInvoiceStatusPartial
	case "cancelled":
		return StudentInvoiceStatusCancelled
	default:
		return StudentInvoiceStatusUnpaid
	}
}
