package inventoryrequests

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

type createInventoryRequestRequest struct {
	ItemID             string  `json:"item_id" binding:"required,uuid"`
	RequesterMemberID  string  `json:"requester_member_id" binding:"required,uuid"`
	Quantity           *int    `json:"quantity"`
	Reason             *string `json:"reason"`
	Status             string  `json:"status" binding:"omitempty,oneof=pending approved rejected completed"`
	ProcessedByStaffID *string `json:"processed_by_staff_id" binding:"omitempty,uuid"`
}

type updateInventoryRequestRequest struct {
	ItemID             string  `json:"item_id" binding:"required,uuid"`
	RequesterMemberID  string  `json:"requester_member_id" binding:"required,uuid"`
	Quantity           *int    `json:"quantity"`
	Reason             *string `json:"reason"`
	Status             string  `json:"status" binding:"required,oneof=pending approved rejected completed"`
	ProcessedByStaffID *string `json:"processed_by_staff_id" binding:"omitempty,uuid"`
}

type inventoryRequestResponse struct {
	ID                 string  `json:"id"`
	ItemID             string  `json:"item_id"`
	RequesterMemberID  string  `json:"requester_member_id"`
	Quantity           *int    `json:"quantity"`
	Reason             *string `json:"reason"`
	Status             string  `json:"status"`
	ProcessedByStaffID *string `json:"processed_by_staff_id"`
	CreatedAt          string  `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	itemID, requesterMemberID, processedByStaffID, ok := parseInventoryRequestCreateUpdateFields(ctx, req.ItemID, req.RequesterMemberID, req.ProcessedByStaffID)
	if !ok {
		return
	}

	status := ent.InventoryRequestStatusPending
	if req.Status != "" {
		status = ent.ToInventoryRequestStatus(req.Status)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInventoryRequestInput{ItemID: itemID, RequesterMemberID: requesterMemberID, Quantity: req.Quantity, Reason: req.Reason, Status: status, ProcessedByStaffID: processedByStaffID})
	if err != nil {
		if errors.Is(err, ErrInventoryItemNotFound) {
			base.ValidateFailed(ctx, ci18n.InventoryItemNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrStaffNotFound) {
			base.ValidateFailed(ctx, ci18n.StaffNotFound, nil)
			return
		}

		log.Errf("inventory-requests.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryRequestResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	itemID, err := utils.ParseQueryUUID(ctx.Query("item_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	requesterMemberID, err := utils.ParseQueryUUID(ctx.Query("requester_member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var status *ent.InventoryRequestStatus
	if raw := ctx.Query("status"); raw != "" {
		value := ent.ToInventoryRequestStatus(raw)
		status = &value
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListInventoryRequestsInput{ItemID: itemID, RequesterMemberID: requesterMemberID, Status: status})
	if err != nil {
		log.Errf("inventory-requests.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]inventoryRequestResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toInventoryRequestResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryRequestID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.InventoryRequestNotFound, nil)
			return
		}
		log.Errf("inventory-requests.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryRequestResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryRequestID(ctx)
	if !ok {
		return
	}

	var req updateInventoryRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	itemID, requesterMemberID, processedByStaffID, ok := parseInventoryRequestCreateUpdateFields(ctx, req.ItemID, req.RequesterMemberID, req.ProcessedByStaffID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInventoryRequestInput{ItemID: itemID, RequesterMemberID: requesterMemberID, Quantity: req.Quantity, Reason: req.Reason, Status: ent.ToInventoryRequestStatus(req.Status), ProcessedByStaffID: processedByStaffID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.InventoryRequestNotFound, nil)
			return
		}
		if errors.Is(err, ErrInventoryItemNotFound) {
			base.ValidateFailed(ctx, ci18n.InventoryItemNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrStaffNotFound) {
			base.ValidateFailed(ctx, ci18n.StaffNotFound, nil)
			return
		}
		log.Errf("inventory-requests.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toInventoryRequestResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseInventoryRequestID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("inventory-requests.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseInventoryRequestID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseInventoryRequestCreateUpdateFields(ctx *gin.Context, itemIDRaw string, requesterMemberIDRaw string, processedByStaffIDRaw *string) (uuid.UUID, uuid.UUID, *uuid.UUID, bool) {
	itemID, err := uuid.Parse(itemIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, nil, false
	}
	requesterMemberID, err := uuid.Parse(requesterMemberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, nil, false
	}
	processedByStaffID, err := utils.ParseUUIDPtr(processedByStaffIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, nil, false
	}

	return itemID, requesterMemberID, processedByStaffID, true
}

func toInventoryRequestResponse(item *ent.InventoryRequest) inventoryRequestResponse {
	return inventoryRequestResponse{ID: item.ID.String(), ItemID: item.ItemID.String(), RequesterMemberID: item.RequesterMemberID.String(), Quantity: item.Quantity, Reason: item.Reason, Status: string(item.Status), ProcessedByStaffID: utils.UUIDToStringPtr(item.ProcessedByStaffID), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}
