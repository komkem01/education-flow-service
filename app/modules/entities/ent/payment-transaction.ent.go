package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodTransfer PaymentMethod = "transfer"
	PaymentMethodQRCode   PaymentMethod = "qr_code"
)

type PaymentTransaction struct {
	bun.BaseModel `bun:"table:payment_transactions,alias:ptr"`

	ID                 uuid.UUID     `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	InvoiceID          uuid.UUID     `bun:"invoice_id,type:uuid,notnull"`
	AmountPaid         *float64      `bun:"amount_paid"`
	PaymentMethod      PaymentMethod `bun:"payment_method"`
	EvidenceURL        *string       `bun:"evidence_url"`
	TransactionDate    *time.Time    `bun:"transaction_date"`
	ProcessedByStaffID *uuid.UUID    `bun:"processed_by_staff_id,type:uuid"`
}

func ToPaymentMethod(value string) PaymentMethod {
	switch value {
	case "transfer":
		return PaymentMethodTransfer
	case "qr_code":
		return PaymentMethodQRCode
	default:
		return PaymentMethodCash
	}
}
