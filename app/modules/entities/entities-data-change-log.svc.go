package entities

import (
	"context"
	"strings"

	"education-flow/app/modules/entities/ent"
	entitiesinf "education-flow/app/modules/entities/inf"

	"github.com/google/uuid"
)

var _ entitiesinf.DataChangeLogEntity = (*Service)(nil)

func (s *Service) CreateDataChangeLog(ctx context.Context, changeLog *ent.DataChangeLog) (*ent.DataChangeLog, error) {
	if _, err := s.db.NewInsert().Model(changeLog).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return changeLog, nil
}

func (s *Service) GetDataChangeLogByID(ctx context.Context, id uuid.UUID) (*ent.DataChangeLog, error) {
	item := new(ent.DataChangeLog)
	if err := s.db.NewSelect().Model(item).Where("id = ?", id).Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) UpdateDataChangeLogByID(ctx context.Context, id uuid.UUID, changeLog *ent.DataChangeLog) (*ent.DataChangeLog, error) {
	updated := new(ent.DataChangeLog)
	if err := s.db.NewUpdate().
		Model(updated).
		Set("audit_log_id = ?", changeLog.AuditLogID).
		Set("table_name = ?", changeLog.TableName).
		Set("record_id = ?", changeLog.RecordID).
		Set("old_values = ?", changeLog.OldValues).
		Set("new_values = ?", changeLog.NewValues).
		Where("id = ?", id).
		Returning("*").
		Scan(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteDataChangeLogByID(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.NewDelete().Model((*ent.DataChangeLog)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (s *Service) ListDataChangeLogs(ctx context.Context, auditLogID *uuid.UUID, tableName *string, recordID *uuid.UUID) ([]*ent.DataChangeLog, error) {
	var items []*ent.DataChangeLog
	query := s.db.NewSelect().Model(&items).Order("created_at DESC")

	if auditLogID != nil {
		query = query.Where("audit_log_id = ?", *auditLogID)
	}
	if tableName != nil {
		query = query.Where("table_name = ?", strings.TrimSpace(*tableName))
	}
	if recordID != nil {
		query = query.Where("record_id = ?", *recordID)
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
