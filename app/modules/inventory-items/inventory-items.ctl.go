package inventoryitems

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

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createInventoryItemRequest struct {
	SchoolID          string  `json:"school_id" binding:"required,uuid"`
	Name              *string `json:"name" binding:"omitempty,max=255"`
	Category          *string `json:"category" binding:"omitempty,max=255"`
	QuantityAvailable *int    `json:"quantity_available"`
	Unit              *string `json:"unit" binding:"omitempty,max=100"`
}

type updateInventoryItemRequest = createInventoryItemRequest

type inventoryItemResponse struct {
	ID                string  `json:"id"`
	SchoolID          string  `json:"school_id"`
	Name              *string `json:"name"`
	Category          *string `json:"category"`
	QuantityAvailable *int    `json:"quantity_available"`
	Unit              *string `json:"unit"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createInventoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInventoryItemInput{SchoolID: schoolID, Name: req.Name, Category: req.Category, QuantityAvailable: req.QuantityAvailable, Unit: req.Unit})
	if err != nil {
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}

		log.Errf("inventory-items.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryItemResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListInventoryItemsInput{SchoolID: schoolID})
	if err != nil {
		log.Errf("inventory-items.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]inventoryItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toInventoryItemResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryItemID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.InventoryItemNotFound, nil)
			return
		}
		log.Errf("inventory-items.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryItemResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryItemID(ctx)
	if !ok {
		return
	}

	var req updateInventoryItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInventoryItemInput{SchoolID: schoolID, Name: req.Name, Category: req.Category, QuantityAvailable: req.QuantityAvailable, Unit: req.Unit})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.InventoryItemNotFound, nil)
			return
		}
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		log.Errf("inventory-items.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryItemResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryItemID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("inventory-items.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseInventoryItemID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func toInventoryItemResponse(item *ent.InventoryItem) inventoryItemResponse {
	return inventoryItemResponse{ID: item.ID.String(), SchoolID: item.SchoolID.String(), Name: item.Name, Category: item.Category, QuantityAvailable: item.QuantityAvailable, Unit: item.Unit}
}
