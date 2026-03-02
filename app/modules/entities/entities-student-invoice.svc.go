package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.StudentInvoiceEntity = (*Service)(nil)

func (s *Service) CreateStudentInvoice(ctx context.Context, invoice *ent.StudentInvoice) (*ent.StudentInvoice, error) {
	if _, err := s.db.NewInsert().Model(invoice).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (s *Service) UpdateStudentInvoiceByID(ctx context.Context, id uuid.UUID, invoice *ent.StudentInvoice) (*ent.StudentInvoice, error) {
	updated := new(ent.StudentInvoice)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("fee_category_id = ?", invoice.FeeCategoryID).
		Set("academic_year_id = ?", invoice.AcademicYearID).
		Set("amount = ?", invoice.Amount).
		Set("due_date = ?", invoice.DueDate).
		Set("status = ?", invoice.Status).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteStudentInvoiceByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.StudentInvoice)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListStudentInvoicesByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentInvoice, error) {
	var invoices []*ent.StudentInvoice
	if err := s.db.NewSelect().Model(&invoices).Where("student_id = ?", studentID).Order("created_at DESC").Scan(ctx); err != nil {
		return nil, err
	}

	return invoices, nil
}

func (s *Service) StudentInvoiceBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		Model((*ent.StudentInvoice)(nil)).
		Where("id = ?", id).
		Where("student_id = ?", studentID).
		Exists(ctx)
}
