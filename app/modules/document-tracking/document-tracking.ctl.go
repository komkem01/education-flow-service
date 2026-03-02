package documenttracking

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

type createDocumentTrackingRequest struct {
	SchoolID         string  `json:"school_id" binding:"required,uuid"`
	DocNumber        *string `json:"doc_number" binding:"omitempty,max=100"`
	Title            *string `json:"title" binding:"omitempty,max=255"`
	ContentSummary   *string `json:"content_summary" binding:"omitempty,max=2000"`
	Priority         string  `json:"priority" binding:"omitempty,oneof=normal urgent top_urgent"`
	SenderMemberID   *string `json:"sender_member_id" binding:"omitempty,uuid"`
	ReceiverMemberID *string `json:"receiver_member_id" binding:"omitempty,uuid"`
	FileURL          *string `json:"file_url" binding:"omitempty,max=1024"`
	Status           string  `json:"status" binding:"omitempty,oneof=sent read processed"`
}

type updateDocumentTrackingRequest struct {
	SchoolID         string  `json:"school_id" binding:"required,uuid"`
	DocNumber        *string `json:"doc_number" binding:"omitempty,max=100"`
	Title            *string `json:"title" binding:"omitempty,max=255"`
	ContentSummary   *string `json:"content_summary" binding:"omitempty,max=2000"`
	Priority         string  `json:"priority" binding:"required,oneof=normal urgent top_urgent"`
	SenderMemberID   *string `json:"sender_member_id" binding:"omitempty,uuid"`
	ReceiverMemberID *string `json:"receiver_member_id" binding:"omitempty,uuid"`
	FileURL          *string `json:"file_url" binding:"omitempty,max=1024"`
	Status           string  `json:"status" binding:"required,oneof=sent read processed"`
}

type documentTrackingResponse struct {
	ID               string  `json:"id"`
	SchoolID         string  `json:"school_id"`
	DocNumber        *string `json:"doc_number"`
	Title            *string `json:"title"`
	ContentSummary   *string `json:"content_summary"`
	Priority         string  `json:"priority"`
	SenderMemberID   *string `json:"sender_member_id"`
	ReceiverMemberID *string `json:"receiver_member_id"`
	FileURL          *string `json:"file_url"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createDocumentTrackingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, senderMemberID, receiverMemberID, ok := parseDocumentTrackingCreateUpdateFields(ctx, req.SchoolID, req.SenderMemberID, req.ReceiverMemberID)
	if !ok {
		return
	}

	priority := ent.DocumentPriorityNormal
	if req.Priority != "" {
		priority = ent.ToDocumentPriority(req.Priority)
	}
	status := ent.DocumentTrackingStatusSent
	if req.Status != "" {
		status = ent.ToDocumentTrackingStatus(req.Status)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateDocumentTrackingInput{SchoolID: schoolID, DocNumber: req.DocNumber, Title: req.Title, ContentSummary: req.ContentSummary, Priority: priority, SenderMemberID: senderMemberID, ReceiverMemberID: receiverMemberID, FileURL: req.FileURL, Status: status})
	if err != nil {
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}

		log.Errf("document-tracking.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toDocumentTrackingResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	senderMemberID, err := utils.ParseQueryUUID(ctx.Query("sender_member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	receiverMemberID, err := utils.ParseQueryUUID(ctx.Query("receiver_member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var status *ent.DocumentTrackingStatus
	if raw := ctx.Query("status"); raw != "" {
		value := ent.ToDocumentTrackingStatus(raw)
		status = &value
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListDocumentTrackingsInput{SchoolID: schoolID, SenderMemberID: senderMemberID, ReceiverMemberID: receiverMemberID, Status: status})
	if err != nil {
		log.Errf("document-tracking.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]documentTrackingResponse, 0, len(items))
	for _, item := range items {
		response = append(response, toDocumentTrackingResponse(item))
	}

	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseDocumentTrackingID(ctx)
	if !ok {
		return
	}

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.DocumentTrackingNotFound, nil)
			return
		}
		log.Errf("document-tracking.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toDocumentTrackingResponse(item))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseDocumentTrackingID(ctx)
	if !ok {
		return
	}

	var req updateDocumentTrackingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, senderMemberID, receiverMemberID, ok := parseDocumentTrackingCreateUpdateFields(ctx, req.SchoolID, req.SenderMemberID, req.ReceiverMemberID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateDocumentTrackingInput{SchoolID: schoolID, DocNumber: req.DocNumber, Title: req.Title, ContentSummary: req.ContentSummary, Priority: ent.ToDocumentPriority(req.Priority), SenderMemberID: senderMemberID, ReceiverMemberID: receiverMemberID, FileURL: req.FileURL, Status: ent.ToDocumentTrackingStatus(req.Status)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.DocumentTrackingNotFound, nil)
			return
		}
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("document-tracking.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, toDocumentTrackingResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	id, ok := parseDocumentTrackingID(ctx)
	if !ok {
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("document-tracking.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String()})
}

func parseDocumentTrackingID(ctx *gin.Context) (uuid.UUID, bool) {
	id, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}

	return id, true
}

func parseDocumentTrackingCreateUpdateFields(ctx *gin.Context, schoolIDRaw string, senderMemberIDRaw *string, receiverMemberIDRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
	schoolID, err := uuid.Parse(schoolIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	senderMemberID, err := utils.ParseUUIDPtr(senderMemberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	receiverMemberID, err := utils.ParseUUIDPtr(receiverMemberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}

	return schoolID, senderMemberID, receiverMemberID, true
}

func toDocumentTrackingResponse(item *ent.DocumentTracking) documentTrackingResponse {
	return documentTrackingResponse{ID: item.ID.String(), SchoolID: item.SchoolID.String(), DocNumber: item.DocNumber, Title: item.Title, ContentSummary: item.ContentSummary, Priority: string(item.Priority), SenderMemberID: utils.UUIDToStringPtr(item.SenderMemberID), ReceiverMemberID: utils.UUIDToStringPtr(item.ReceiverMemberID), FileURL: item.FileURL, Status: string(item.Status), CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}
