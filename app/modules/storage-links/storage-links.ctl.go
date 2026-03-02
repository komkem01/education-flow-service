package storagelinks

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
	StorageID  string  `json:"storage_id" binding:"required,uuid"`
	EntityType string  `json:"entity_type" binding:"required,max=100"`
	EntityID   string  `json:"entity_id" binding:"required,uuid"`
	FieldName  *string `json:"field_name" binding:"omitempty,max=100"`
	SortOrder  int     `json:"sort_order"`
}

type updateRequest = createRequest

type response struct {
	ID         string  `json:"id"`
	StorageID  string  `json:"storage_id"`
	EntityType string  `json:"entity_type"`
	EntityID   string  `json:"entity_id"`
	FieldName  *string `json:"field_name"`
	SortOrder  int     `json:"sort_order"`
	CreatedAt  string  `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	storageID, entityID, ok := parseCreateUpdateFields(ctx, req.StorageID, req.EntityID)
	if !ok {
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StorageID: storageID, EntityType: req.EntityType, EntityID: entityID, FieldName: req.FieldName, SortOrder: req.SortOrder})
	if err != nil {
		if errors.Is(err, ErrStorageNotFound) {
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		log.Errf("storage-links.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	storageID, err := utils.ParseQueryUUID(ctx.Query("storage_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	entityID, err := utils.ParseQueryUUID(ctx.Query("entity_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	entityType := trimmedQueryPtr(ctx.Query("entity_type"))

	items, err := c.svc.List(ctx.Request.Context(), &ListInput{StorageID: storageID, EntityType: entityType, EntityID: entityID})
	if err != nil {
		log.Errf("storage-links.list.error: %v", err)
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

	sortBy := strings.TrimSpace(ctx.DefaultQuery("sort_by", "sort_order"))
	orderBy := strings.ToLower(strings.TrimSpace(ctx.DefaultQuery("order_by", "asc")))
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
			base.ValidateFailed(ctx, ci18n.StorageLinkNotFound, nil)
			return
		}
		log.Errf("storage-links.get.error: %v", err)
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
	storageID, entityID, ok := parseCreateUpdateFields(ctx, req.StorageID, req.EntityID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{StorageID: storageID, EntityType: req.EntityType, EntityID: entityID, FieldName: req.FieldName, SortOrder: req.SortOrder})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StorageLinkNotFound, nil)
			return
		}
		if errors.Is(err, ErrStorageNotFound) {
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		log.Errf("storage-links.update.error: %v", err)
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
		log.Errf("storage-links.delete.error: %v", err)
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

func parseCreateUpdateFields(ctx *gin.Context, storageIDRaw string, entityIDRaw string) (uuid.UUID, uuid.UUID, bool) {
	storageID, err := uuid.Parse(storageIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	entityID, err := uuid.Parse(entityIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	return storageID, entityID, true
}

func toResponse(item *ent.StorageLink) response {
	return response{ID: item.ID.String(), StorageID: item.StorageID.String(), EntityType: item.EntityType, EntityID: item.EntityID.String(), FieldName: item.FieldName, SortOrder: item.SortOrder, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
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
		if strings.Contains(strings.ToLower(item.EntityType), needle) ||
			(item.FieldName != nil && strings.Contains(strings.ToLower(*item.FieldName), needle)) {
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
		case "created_at":
			less = items[i].CreatedAt < items[j].CreatedAt
		case "entity_type":
			less = strings.ToLower(items[i].EntityType) < strings.ToLower(items[j].EntityType)
		default:
			less = items[i].SortOrder < items[j].SortOrder
		}

		if asc {
			return less
		}
		return !less
	})

	return items
}
