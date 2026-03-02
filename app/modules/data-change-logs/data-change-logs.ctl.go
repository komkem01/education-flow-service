package datachangelogs

import (
	"database/sql"
	"errors"
	"sort"
	"strconv"
	"strings"

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

	search := strings.TrimSpace(ctx.Query("search"))
	if search != "" {
		responseList = filterResponses(responseList, search)
	}

	sortBy := strings.TrimSpace(ctx.DefaultQuery("sort_by", "created_at"))
	orderBy := strings.ToLower(strings.TrimSpace(ctx.DefaultQuery("order_by", "desc")))
	responseList = sortResponses(responseList, sortBy, orderBy)

	page, size := parsePageSize(ctx)
	total := int64(len(responseList))
	responseList = paginateResponses(responseList, page, size)

	base.Paginate(ctx, responseList, &base.ResponsePaginate{Page: int64(page), Size: int64(size), Total: total})
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
	base.Forbidden(ctx, ci18n.Forbidden, nil)
}

func (c *Controller) Delete(ctx *gin.Context) {
	base.Forbidden(ctx, ci18n.Forbidden, nil)
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

func parsePageSize(ctx *gin.Context) (int, int) {
	page := 1
	size := 20

	if raw := strings.TrimSpace(ctx.Query("page")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if raw := strings.TrimSpace(ctx.Query("size")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			if parsed > 100 {
				parsed = 100
			}
			size = parsed
		}
	}

	return page, size
}

func paginateResponses(items []response, page int, size int) []response {
	if len(items) == 0 {
		return items
	}

	start := (page - 1) * size
	if start >= len(items) {
		return []response{}
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

func filterResponses(items []response, search string) []response {
	needle := strings.ToLower(search)
	filtered := make([]response, 0, len(items))
	for _, item := range items {
		table := ""
		if item.TableName != nil {
			table = strings.ToLower(*item.TableName)
		}
		if strings.Contains(table, needle) || strings.Contains(item.AuditLogID, needle) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func sortResponses(items []response, sortBy string, orderBy string) []response {
	if len(items) < 2 {
		return items
	}

	asc := orderBy == "asc"
	sort.Slice(items, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "table_name":
			a := ""
			if items[i].TableName != nil {
				a = strings.ToLower(*items[i].TableName)
			}
			b := ""
			if items[j].TableName != nil {
				b = strings.ToLower(*items[j].TableName)
			}
			less = a < b
		default:
			less = items[i].CreatedAt < items[j].CreatedAt
		}

		if asc {
			return less
		}
		return !less
	})

	return items
}
