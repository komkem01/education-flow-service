package entities

import (
	"context"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.SystemAuditLogEntity = (*Service)(nil)

func (s *Service) CreateSystemAuditLog(ctx context.Context, auditLog *ent.SystemAuditLog) (*ent.SystemAuditLog, error) {
	if _, err := s.db.NewInsert().Model(auditLog).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return auditLog, nil
}

func (s *Service) GetSystemAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.SystemAuditLog, error) {
	item := new(ent.SystemAuditLog)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) UpdateSystemAuditLogByID(ctx context.Context, id uuid.UUID, auditLog *ent.SystemAuditLog) (*ent.SystemAuditLog, error) {
	updated := new(ent.SystemAuditLog)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("member_id = ?", auditLog.MemberID).
		Set("action = ?", auditLog.Action).
		Set("module = ?", auditLog.Module).
		Set("description = ?", auditLog.Description).
		Set("ip_address = ?", auditLog.IPAddress).
		Set("user_agent = ?", auditLog.UserAgent).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteSystemAuditLogByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.SystemAuditLog)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListSystemAuditLogs(ctx context.Context, memberID *uuid.UUID, module *string) ([]*ent.SystemAuditLog, error) {
	var items []*ent.SystemAuditLog
	query := s.db.NewSelect().Model(&items).Order("created_at DESC")

	if memberID != nil {
		query = query.Where("member_id = ?", *memberID)
	}
	if module != nil {
		query = query.Where("module = ?", strings.TrimSpace(*module))
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
