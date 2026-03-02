package datachangelogs

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const dateTimeLayout = "2006-01-02T15:04:05Z07:00"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	AuditLogID string         `json:"audit_log_id" binding:"required,uuid"`
	TableName  *string        `json:"table_name" binding:"omitempty,max=255"`
	RecordID   *string        `json:"record_id" binding:"omitempty,uuid"`
	OldValues  map[string]any `json:"old_values"`
	NewValues  map[string]any `json:"new_values"`
}

type updateRequest = createRequest

type response struct {
	ID         string         `json:"id"`
	AuditLogID string         `json:"audit_log_id"`
	TableName  *string        `json:"table_name"`
	RecordID   *string        `json:"record_id"`
	OldValues  map[string]any `json:"old_values"`
	NewValues  map[string]any `json:"new_values"`
	CreatedAt  string         `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	auditLogID, recordID, ok := parseCreateUpdateFields(ctx, req.AuditLogID, req.RecordID)
	if !ok {
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{AuditLogID: auditLogID, TableName: req.TableName, RecordID: recordID, OldValues: req.OldValues, NewValues: req.NewValues})
	if err != nil {
		if errors.Is(err, ErrAuditLogNotFound) {
			base.ValidateFailed(ctx, ci18n.SystemAuditLogNotFound, nil)
			return
		}
		log.Errf("data-change-logs.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	auditLogID, err := utils.ParseQueryUUID(ctx.Query("audit_log_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	recordID, err := utils.ParseQueryUUID(ctx.Query("record_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	tableName := trimmedQueryPtr(ctx.Query("table_name"))

	items, err := c.svc.List(ctx.Request.Context(), &ListInput{AuditLogID: auditLogID, TableName: tableName, RecordID: recordID})
	if err != nil {
		log.Errf("data-change-logs.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}
	base.Success(ctx, responseList)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.DataChangeLogNotFound, nil)
			return
		}
		log.Errf("data-change-logs.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	auditLogID, recordID, ok := parseCreateUpdateFields(ctx, req.AuditLogID, req.RecordID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{AuditLogID: auditLogID, TableName: req.TableName, RecordID: recordID, OldValues: req.OldValues, NewValues: req.NewValues})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.DataChangeLogNotFound, nil)
			return
		}
		if errors.Is(err, ErrAuditLogNotFound) {
			base.ValidateFailed(ctx, ci18n.SystemAuditLogNotFound, nil)
			return
		}
		log.Errf("data-change-logs.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("data-change-logs.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": id.String()})
}

func parseID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func parseCreateUpdateFields(ctx *gin.Context, auditLogIDRaw string, recordIDRaw *string) (uuid.UUID, *uuid.UUID, bool) {
	auditLogID, err := uuid.Parse(auditLogIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, false
	}
	recordID, err := utils.ParseUUIDPtr(recordIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, false
	}
	return auditLogID, recordID, true
}

func toResponse(item *ent.DataChangeLog) response {
	return response{ID: item.ID.String(), AuditLogID: item.AuditLogID.String(), TableName: item.TableName, RecordID: utils.UUIDToStringPtr(item.RecordID), OldValues: item.OldValues, NewValues: item.NewValues, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func trimmedQueryPtr(raw string) *string {
	if raw == "" {
		return nil
	}
	value := raw
	return &value
}
