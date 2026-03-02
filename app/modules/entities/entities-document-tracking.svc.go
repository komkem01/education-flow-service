package entities

import (
	"context"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.DocumentTrackingEntity = (*Service)(nil)

func (s *Service) CreateDocumentTracking(ctx context.Context, document *ent.DocumentTracking) (*ent.DocumentTracking, error) {
	if _, err := s.db.NewInsert().Model(document).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *Service) GetDocumentTrackingByID(ctx context.Context, id uuid.UUID) (*ent.DocumentTracking, error) {
	document := new(ent.DocumentTracking)
	if err := s.db.NewSelect().Model(document).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *Service) UpdateDocumentTrackingByID(ctx context.Context, id uuid.UUID, document *ent.DocumentTracking) (*ent.DocumentTracking, error) {
	updated := new(ent.DocumentTracking)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("school_id = ?", document.SchoolID).
		Set("doc_number = ?", document.DocNumber).
		Set("title = ?", document.Title).
		Set("content_summary = ?", document.ContentSummary).
		Set("priority = ?", document.Priority).
		Set("sender_member_id = ?", document.SenderMemberID).
		Set("receiver_member_id = ?", document.ReceiverMemberID).
		Set("file_url = ?", document.FileURL).
		Set("status = ?", document.Status).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteDocumentTrackingByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.DocumentTracking)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListDocumentTrackings(ctx context.Context, schoolID *uuid.UUID, senderMemberID *uuid.UUID, receiverMemberID *uuid.UUID, status *ent.DocumentTrackingStatus) ([]*ent.DocumentTracking, error) {
	var documents []*ent.DocumentTracking
	query := s.db.NewSelect().Model(&documents).Order("created_at DESC")

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	}
	if senderMemberID != nil {
		query = query.Where("sender_member_id = ?", *senderMemberID)
	}
	if receiverMemberID != nil {
		query = query.Where("receiver_member_id = ?", *receiverMemberID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return documents, nil
}
