package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	datachangelogs "education-flow/app/modules/data-change-logs"
	"education-flow/app/modules"
	systemauditlogs "education-flow/app/modules/system-audit-logs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const maxAuditBodyBytes = 8192
const auditQueueSize = 4096
const auditWriteTimeout = 2 * time.Second

var (
	auditWorkerOnce sync.Once
	auditQueue      chan *auditQueueItem
	auditMetrics    = &auditQueueMetrics{}
)

type auditQueueItem struct {
	auditInput *systemauditlogs.CreateInput
	oldValues  map[string]any
}

type auditQueueMetrics struct {
	enqueued       atomic.Uint64
	written        atomic.Uint64
	dropped        atomic.Uint64
	writeErrors    atomic.Uint64
	workerStarted  atomic.Bool
	lastWriteError atomic.Value
}

type auditQueueHealth struct {
	WorkerStarted bool   `json:"worker_started"`
	QueueLength   int    `json:"queue_length"`
	QueueCapacity int    `json:"queue_capacity"`
	EnqueuedTotal uint64 `json:"enqueued_total"`
	WrittenTotal  uint64 `json:"written_total"`
	DroppedTotal  uint64 `json:"dropped_total"`
	WriteErrors   uint64 `json:"write_errors"`
	LastWriteErr  string `json:"last_write_error,omitempty"`
}

type responseCaptureWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCaptureWriter) Write(data []byte) (int, error) {
	if w.body.Len() < maxAuditBodyBytes {
		remaining := maxAuditBodyBytes - w.body.Len()
		if len(data) > remaining {
			w.body.Write(data[:remaining])
		} else {
			w.body.Write(data)
		}
	}
	return w.ResponseWriter.Write(data)
}

func auditLogMiddleware(mod *modules.Modules) gin.HandlerFunc {
	initAuditWorker(mod)

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		resourceType, resourceID := extractResource(path)
		oldValues := captureOldValues(ctx.Request.Context(), mod, resourceType, resourceID, ctx.Request.Method)

		requestBody := readRequestBody(ctx)
		queryParams := buildQueryParams(ctx)

		writer := &responseCaptureWriter{ResponseWriter: ctx.Writer, body: bytes.NewBuffer(nil)}
		ctx.Writer = writer

		startedAt := time.Now()
		ctx.Next()
		duration := time.Since(startedAt).Milliseconds()

		responseStatus := ctx.Writer.Status()
		responseBody := parseJSONBody(writer.body.Bytes())

		module := resourceType
		action := strings.ToUpper(ctx.Request.Method)
		httpMethod := ctx.Request.Method
		httpPath := path
		routePath := ctx.FullPath()
		ipAddress := ctx.ClientIP()
		userAgent := ctx.Request.UserAgent()
		outcome := "success"
		if responseStatus >= 400 {
			outcome = "error"
		}

		var errorMessage *string
		if len(ctx.Errors) > 0 {
			value := strings.TrimSpace(ctx.Errors.String())
			if value != "" {
				errorMessage = &value
			}
		}

		var memberID *uuid.UUID
		if raw := strings.TrimSpace(ctx.GetHeader("X-Member-ID")); raw != "" {
			if parsed, err := uuid.Parse(raw); err == nil {
				memberID = &parsed
			}
		}

		var requestID *string
		if raw := strings.TrimSpace(ctx.GetHeader("X-Request-ID")); raw != "" {
			requestID = &raw
		}

		spanCtx := trace.SpanContextFromContext(ctx.Request.Context())
		var traceID *string
		var spanID *string
		if spanCtx.IsValid() {
			t := spanCtx.TraceID().String()
			s := spanCtx.SpanID().String()
			traceID = &t
			spanID = &s
		}

		actorType := "anonymous"
		if memberID != nil {
			actorType = "member"
		}

		enqueueAuditLog(&auditQueueItem{auditInput: &systemauditlogs.CreateInput{
			MemberID:        memberID,
			Action:          &action,
			Module:          &module,
			Description:     nil,
			IPAddress:       &ipAddress,
			UserAgent:       &userAgent,
			ActorType:       &actorType,
			ActorIdentifier: nil,
			TraceID:         traceID,
			SpanID:          spanID,
			RequestID:       requestID,
			HTTPMethod:      &httpMethod,
			HTTPPath:        &httpPath,
			RoutePath:       &routePath,
			QueryParams:     queryParams,
			RequestBody:     requestBody,
			ResponseStatus:  &responseStatus,
			ResponseBody:    responseBody,
			ErrorMessage:    errorMessage,
			Outcome:         &outcome,
			ResourceType:    &resourceType,
			ResourceID:      resourceID,
			DurationMS:      &duration,
		}, oldValues: oldValues})
	}
}

func initAuditWorker(mod *modules.Modules) {
	auditWorkerOnce.Do(func() {
		auditQueue = make(chan *auditQueueItem, auditQueueSize)
		auditMetrics.workerStarted.Store(true)
		go func() {
			for item := range auditQueue {
				if item == nil || item.auditInput == nil {
					continue
				}

				input := item.auditInput
				writeCtx, cancel := context.WithTimeout(context.Background(), auditWriteTimeout)
				auditLog, err := mod.SystemAuditLog.Svc.Create(writeCtx, input)
				cancel()

				if err != nil {
					auditMetrics.writeErrors.Add(1)
					auditMetrics.lastWriteError.Store(err.Error())
					continue
				}

				if shouldCreateDataChangeLog(input) && mod.DataChangeLog != nil {
					dataChangeInput := toDataChangeLogInput(auditLog.ID, input, item.oldValues)
					if dataChangeInput != nil {
						changeCtx, changeCancel := context.WithTimeout(context.Background(), auditWriteTimeout)
						_, dataErr := mod.DataChangeLog.Svc.Create(changeCtx, dataChangeInput)
						changeCancel()
						if dataErr != nil {
							auditMetrics.writeErrors.Add(1)
							auditMetrics.lastWriteError.Store(dataErr.Error())
						}
					}
				}

				auditMetrics.written.Add(1)
				auditMetrics.lastWriteError.Store("")
			}
		}()
	})
}

func enqueueAuditLog(input *auditQueueItem) {
	if auditQueue == nil {
		auditMetrics.dropped.Add(1)
		return
	}

	select {
	case auditQueue <- input:
		auditMetrics.enqueued.Add(1)
	default:
		auditMetrics.dropped.Add(1)
	}
}

func auditQueueHealthSnapshot() auditQueueHealth {
	queueLength := 0
	queueCapacity := 0
	if auditQueue != nil {
		queueLength = len(auditQueue)
		queueCapacity = cap(auditQueue)
	}

	health := auditQueueHealth{
		WorkerStarted: auditMetrics.workerStarted.Load(),
		QueueLength:   queueLength,
		QueueCapacity: queueCapacity,
		EnqueuedTotal: auditMetrics.enqueued.Load(),
		WrittenTotal:  auditMetrics.written.Load(),
		DroppedTotal:  auditMetrics.dropped.Load(),
		WriteErrors:   auditMetrics.writeErrors.Load(),
	}

	if value := auditMetrics.lastWriteError.Load(); value != nil {
		if text, ok := value.(string); ok {
			health.LastWriteErr = strings.TrimSpace(text)
		}
	}

	return health
}

func auditHealthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, auditQueueHealthSnapshot())
	}
}

func readRequestBody(ctx *gin.Context) map[string]any {
	if ctx.Request.Body == nil {
		return nil
	}

	raw, err := io.ReadAll(io.LimitReader(ctx.Request.Body, maxAuditBodyBytes))
	if err != nil {
		return nil
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
	return parseJSONBody(raw)
}

func parseJSONBody(raw []byte) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil
	}

	var mapped map[string]any
	if err := json.Unmarshal(trimmed, &mapped); err == nil {
		return mapped
	}

	return nil
}

func buildQueryParams(ctx *gin.Context) map[string]any {
	params := map[string]any{}
	for key, values := range ctx.Request.URL.Query() {
		if len(values) == 1 {
			params[key] = values[0]
			continue
		}
		copied := make([]string, len(values))
		copy(copied, values)
		params[key] = copied
	}
	if len(params) == 0 {
		return nil
	}
	return params
}

func extractResource(path string) (string, *uuid.UUID) {
	resourceType := strings.Trim(path, "/")
	if resourceType == "" {
		resourceType = "root"
	}
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 3 {
		resourceType = parts[2]
	} else if len(parts) >= 1 && parts[0] != "" {
		resourceType = parts[len(parts)-1]
	}

	for _, part := range parts {
		id, err := uuid.Parse(part)
		if err == nil {
			return resourceType, &id
		}
	}

	return resourceType, nil
}

func shouldCreateDataChangeLog(input *systemauditlogs.CreateInput) bool {
	if input == nil || input.HTTPMethod == nil {
		return false
	}

	if input.ResponseStatus == nil || *input.ResponseStatus >= 400 {
		return false
	}

	method := strings.ToUpper(strings.TrimSpace(*input.HTTPMethod))
	return method == "POST" || method == "PUT" || method == "PATCH" || method == "DELETE"
}

func toDataChangeLogInput(auditLogID uuid.UUID, input *systemauditlogs.CreateInput, oldValues map[string]any) *datachangelogs.CreateInput {
	if input == nil || input.HTTPMethod == nil {
		return nil
	}

	operation := strings.ToUpper(strings.TrimSpace(*input.HTTPMethod))
	if operation == "" {
		return nil
	}

	tableName := input.ResourceType
	recordID := deriveRecordID(input)
	changedFields := collectChangedFields(input.RequestBody)
	source := "audit-middleware"

	return &datachangelogs.CreateInput{
		AuditLogID:        auditLogID,
		TableName:         tableName,
		RecordID:          recordID,
		Operation:         &operation,
		ChangedFields:     changedFields,
		ChangedByMemberID: input.MemberID,
		Source:            &source,
		Reason:            nil,
		OldValues:         oldValues,
		NewValues:         input.RequestBody,
	}
}

func collectChangedFields(requestBody map[string]any) []string {
	if len(requestBody) == 0 {
		return nil
	}

	fields := make([]string, 0, len(requestBody))
	for key := range requestBody {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}
		fields = append(fields, trimmed)
	}

	if len(fields) == 0 {
		return nil
	}

	sort.Strings(fields)
	return fields
}

func deriveRecordID(input *systemauditlogs.CreateInput) *uuid.UUID {
	if input == nil {
		return nil
	}

	if input.ResourceID != nil {
		return input.ResourceID
	}

	if input.ResponseBody == nil {
		return nil
	}

	raw, ok := input.ResponseBody["id"]
	if !ok {
		return nil
	}

	value, ok := raw.(string)
	if !ok {
		return nil
	}

	parsed, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return nil
	}

	return &parsed
}

func captureOldValues(ctx context.Context, mod *modules.Modules, resourceType string, resourceID *uuid.UUID, method string) map[string]any {
	if mod == nil || resourceID == nil {
		return nil
	}

	normalizedMethod := strings.ToUpper(strings.TrimSpace(method))
	if normalizedMethod != "PATCH" && normalizedMethod != "PUT" && normalizedMethod != "DELETE" {
		return nil
	}

	resource := strings.ToLower(strings.TrimSpace(resourceType))
	var value any
	var err error

	switch resource {
	case "prefixes":
		value, err = mod.Prefix.Svc.GetByID(ctx, *resourceID)
	case "genders":
		value, err = mod.Gender.Svc.GetByID(ctx, *resourceID)
	case "schools":
		value, err = mod.School.Svc.GetByID(ctx, *resourceID)
	default:
		return nil
	}

	if err != nil || value == nil {
		return nil
	}

	encoded, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		return nil
	}

	var mapped map[string]any
	if unmarshalErr := json.Unmarshal(encoded, &mapped); unmarshalErr != nil {
		return nil
	}

	if len(mapped) == 0 {
		return nil
	}

	return mapped
}
