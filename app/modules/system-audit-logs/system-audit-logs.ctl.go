package systemauditlogs

import (
	"database/sql"
	"errors"
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
	MemberID    *string `json:"member_id" binding:"omitempty,uuid"`
	Action      *string `json:"action" binding:"omitempty,max=100"`
	Module      *string `json:"module" binding:"omitempty,max=100"`
	Description *string `json:"description"`
	IPAddress   *string `json:"ip_address" binding:"omitempty,max=100"`
	UserAgent   *string `json:"user_agent"`
}

type updateRequest = createRequest

type response struct {
	ID          string  `json:"id"`
	MemberID    *string `json:"member_id"`
	Action      *string `json:"action"`
	Module      *string `json:"module"`
	Description *string `json:"description"`
	IPAddress   *string `json:"ip_address"`
	UserAgent   *string `json:"user_agent"`
	CreatedAt   string  `json:"created_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, ok := parseMemberID(ctx, req.MemberID)
	if !ok {
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{MemberID: memberID, Action: req.Action, Module: req.Module, Description: req.Description, IPAddress: req.IPAddress, UserAgent: req.UserAgent})
	if err != nil {
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("system-audit-logs.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	memberID, err := utils.ParseQueryUUID(ctx.Query("member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	var module *string
	if raw := strings.TrimSpace(ctx.Query("module")); raw != "" {
		module = &raw
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListInput{MemberID: memberID, Module: module})
	if err != nil {
		log.Errf("system-audit-logs.list.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.SystemAuditLogNotFound, nil)
			return
		}
		log.Errf("system-audit-logs.get.error: %v", err)
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
	memberID, ok := parseMemberID(ctx, req.MemberID)
	if !ok {
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{MemberID: memberID, Action: req.Action, Module: req.Module, Description: req.Description, IPAddress: req.IPAddress, UserAgent: req.UserAgent})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SystemAuditLogNotFound, nil)
			return
		}
		if errors.Is(err, ErrMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		log.Errf("system-audit-logs.update.error: %v", err)
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
		log.Errf("system-audit-logs.delete.error: %v", err)
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

func parseMemberID(ctx *gin.Context, raw *string) (*uuid.UUID, bool) {
	memberID, err := utils.ParseUUIDPtr(raw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return nil, false
	}
	return memberID, true
}

func toResponse(item *ent.SystemAuditLog) response {
	return response{ID: item.ID.String(), MemberID: utils.UUIDToStringPtr(item.MemberID), Action: item.Action, Module: item.Module, Description: item.Description, IPAddress: item.IPAddress, UserAgent: item.UserAgent, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}
