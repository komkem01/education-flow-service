package systemauditlogs

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
	MemberID        *string        `json:"member_id" binding:"omitempty,uuid"`
	Action          *string        `json:"action" binding:"omitempty,max=100"`
	Module          *string        `json:"module" binding:"omitempty,max=100"`
	Description     *string        `json:"description"`
	IPAddress       *string        `json:"ip_address" binding:"omitempty,max=100"`
	UserAgent       *string        `json:"user_agent"`
	ActorType       *string        `json:"actor_type" binding:"omitempty,max=100"`
	ActorIdentifier *string        `json:"actor_identifier" binding:"omitempty,max=255"`
	TraceID         *string        `json:"trace_id" binding:"omitempty,max=64"`
	SpanID          *string        `json:"span_id" binding:"omitempty,max=32"`
	RequestID       *string        `json:"request_id" binding:"omitempty,max=128"`
	HTTPMethod      *string        `json:"http_method" binding:"omitempty,max=10"`
	HTTPPath        *string        `json:"http_path"`
	RoutePath       *string        `json:"route_path"`
	QueryParams     map[string]any `json:"query_params"`
	RequestBody     map[string]any `json:"request_body"`
	ResponseStatus  *int           `json:"response_status"`
	ResponseBody    map[string]any `json:"response_body"`
	ErrorMessage    *string        `json:"error_message"`
	Outcome         *string        `json:"outcome" binding:"omitempty,max=50"`
	ResourceType    *string        `json:"resource_type" binding:"omitempty,max=100"`
	ResourceID      *string        `json:"resource_id" binding:"omitempty,uuid"`
	DurationMS      *int64         `json:"duration_ms"`
}

type updateRequest = createRequest

type response struct {
	ID              string         `json:"id"`
	MemberID        *string        `json:"member_id"`
	Action          *string        `json:"action"`
	Module          *string        `json:"module"`
	Description     *string        `json:"description"`
	IPAddress       *string        `json:"ip_address"`
	UserAgent       *string        `json:"user_agent"`
	ActorType       *string        `json:"actor_type"`
	ActorIdentifier *string        `json:"actor_identifier"`
	TraceID         *string        `json:"trace_id"`
	SpanID          *string        `json:"span_id"`
	RequestID       *string        `json:"request_id"`
	HTTPMethod      *string        `json:"http_method"`
	HTTPPath        *string        `json:"http_path"`
	RoutePath       *string        `json:"route_path"`
	QueryParams     map[string]any `json:"query_params"`
	RequestBody     map[string]any `json:"request_body"`
	ResponseStatus  *int           `json:"response_status"`
	ResponseBody    map[string]any `json:"response_body"`
	ErrorMessage    *string        `json:"error_message"`
	Outcome         *string        `json:"outcome"`
	ResourceType    *string        `json:"resource_type"`
	ResourceID      *string        `json:"resource_id"`
	DurationMS      *int64         `json:"duration_ms"`
	CreatedAt       string         `json:"created_at"`
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
	resourceID, ok := parseResourceID(ctx, req.ResourceID)
	if !ok {
		return
	}

	enrichAuditContext(ctx, &req)

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{MemberID: memberID, Action: req.Action, Module: req.Module, Description: req.Description, IPAddress: req.IPAddress, UserAgent: req.UserAgent, ActorType: req.ActorType, ActorIdentifier: req.ActorIdentifier, TraceID: req.TraceID, SpanID: req.SpanID, RequestID: req.RequestID, HTTPMethod: req.HTTPMethod, HTTPPath: req.HTTPPath, RoutePath: req.RoutePath, QueryParams: req.QueryParams, RequestBody: req.RequestBody, ResponseStatus: req.ResponseStatus, ResponseBody: req.ResponseBody, ErrorMessage: req.ErrorMessage, Outcome: req.Outcome, ResourceType: req.ResourceType, ResourceID: resourceID, DurationMS: req.DurationMS})
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

func parseMemberID(ctx *gin.Context, raw *string) (*uuid.UUID, bool) {
	memberID, err := utils.ParseUUIDPtr(raw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return nil, false
	}
	return memberID, true
}

func toResponse(item *ent.SystemAuditLog) response {
	return response{ID: item.ID.String(), MemberID: utils.UUIDToStringPtr(item.MemberID), Action: item.Action, Module: item.Module, Description: item.Description, IPAddress: item.IPAddress, UserAgent: item.UserAgent, ActorType: item.ActorType, ActorIdentifier: item.ActorIdentifier, TraceID: item.TraceID, SpanID: item.SpanID, RequestID: item.RequestID, HTTPMethod: item.HTTPMethod, HTTPPath: item.HTTPPath, RoutePath: item.RoutePath, QueryParams: item.QueryParams, RequestBody: item.RequestBody, ResponseStatus: item.ResponseStatus, ResponseBody: item.ResponseBody, ErrorMessage: item.ErrorMessage, Outcome: item.Outcome, ResourceType: item.ResourceType, ResourceID: utils.UUIDToStringPtr(item.ResourceID), DurationMS: item.DurationMS, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout)}
}

func parseResourceID(ctx *gin.Context, raw *string) (*uuid.UUID, bool) {
	id, err := utils.ParseUUIDPtr(raw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return nil, false
	}
	return id, true
}

func enrichAuditContext(ctx *gin.Context, req *createRequest) {
	if req.IPAddress == nil {
		ip := ctx.ClientIP()
		req.IPAddress = &ip
	}
	if req.UserAgent == nil {
		ua := ctx.Request.UserAgent()
		req.UserAgent = &ua
	}
	if req.HTTPMethod == nil {
		method := ctx.Request.Method
		req.HTTPMethod = &method
	}
	if req.HTTPPath == nil {
		path := ctx.Request.URL.Path
		req.HTTPPath = &path
	}
	if req.RoutePath == nil {
		routePath := ctx.FullPath()
		req.RoutePath = &routePath
	}
	if req.QueryParams == nil {
		params := map[string]any{}
		for key, values := range ctx.Request.URL.Query() {
			if len(values) == 1 {
				params[key] = values[0]
			} else {
				copied := make([]string, len(values))
				copy(copied, values)
				params[key] = copied
			}
		}
		req.QueryParams = params
	}
	spanCtx := trace.SpanContextFromContext(ctx.Request.Context())
	if spanCtx.IsValid() {
		if req.TraceID == nil {
			traceID := spanCtx.TraceID().String()
			req.TraceID = &traceID
		}
		if req.SpanID == nil {
			spanID := spanCtx.SpanID().String()
			req.SpanID = &spanID
		}
	}
	if req.RequestID == nil {
		if value := strings.TrimSpace(ctx.GetHeader("X-Request-ID")); value != "" {
			req.RequestID = &value
		}
	}
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
		action := ""
		if item.Action != nil {
			action = strings.ToLower(*item.Action)
		}
		module := ""
		if item.Module != nil {
			module = strings.ToLower(*item.Module)
		}
		description := ""
		if item.Description != nil {
			description = strings.ToLower(*item.Description)
		}
		if strings.Contains(action, needle) || strings.Contains(module, needle) || strings.Contains(description, needle) {
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
		case "module":
			a := ""
			if items[i].Module != nil {
				a = strings.ToLower(*items[i].Module)
			}
			b := ""
			if items[j].Module != nil {
				b = strings.ToLower(*items[j].Module)
			}
			less = a < b
		case "action":
			a := ""
			if items[i].Action != nil {
				a = strings.ToLower(*items[i].Action)
			}
			b := ""
			if items[j].Action != nil {
				b = strings.ToLower(*items[j].Action)
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
