package storages

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
	SchoolID            string         `json:"school_id" binding:"required,uuid"`
	BucketName          string         `json:"bucket_name" binding:"required,max=255"`
	ObjectKey           string         `json:"object_key" binding:"required"`
	OriginalName        *string        `json:"original_name" binding:"omitempty,max=255"`
	Extension           *string        `json:"extension" binding:"omitempty,max=20"`
	MIMEType            *string        `json:"mime_type" binding:"omitempty,max=255"`
	SizeBytes           int64          `json:"size_bytes"`
	ChecksumSHA256      *string        `json:"checksum_sha256" binding:"omitempty,max=64"`
	ETag                *string        `json:"etag" binding:"omitempty,max=255"`
	Visibility          string         `json:"visibility" binding:"omitempty,oneof=private public signed"`
	Status              string         `json:"status" binding:"omitempty,oneof=pending active obsolete deleted"`
	VirusScanStatus     string         `json:"virus_scan_status" binding:"omitempty,oneof=pending clean infected failed"`
	UploadedByMemberID  *string        `json:"uploaded_by_member_id" binding:"omitempty,uuid"`
	VersionNo           int            `json:"version_no"`
	ReplacedByStorageID *string        `json:"replaced_by_storage_id" binding:"omitempty,uuid"`
	Metadata            map[string]any `json:"metadata"`
}

type updateRequest = createRequest

type response struct {
	ID                  string         `json:"id"`
	SchoolID            string         `json:"school_id"`
	BucketName          string         `json:"bucket_name"`
	ObjectKey           string         `json:"object_key"`
	OriginalName        *string        `json:"original_name"`
	Extension           *string        `json:"extension"`
	MIMEType            *string        `json:"mime_type"`
	SizeBytes           int64          `json:"size_bytes"`
	ChecksumSHA256      *string        `json:"checksum_sha256"`
	ETag                *string        `json:"etag"`
	Visibility          string         `json:"visibility"`
	Status              string         `json:"status"`
	VirusScanStatus     string         `json:"virus_scan_status"`
	UploadedByMemberID  *string        `json:"uploaded_by_member_id"`
	VersionNo           int            `json:"version_no"`
	ReplacedByStorageID *string        `json:"replaced_by_storage_id"`
	Metadata            map[string]any `json:"metadata"`
	CreatedAt           string         `json:"created_at"`
	UpdatedAt           string         `json:"updated_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	schoolID, uploadedByID, replacedByID, ok := parseCreateUpdateFields(ctx, req.SchoolID, req.UploadedByMemberID, req.ReplacedByStorageID)
	if !ok {
		return
	}
	visibility := ent.StorageVisibilityPrivate
	if req.Visibility != "" {
		visibility = ent.ToStorageVisibility(req.Visibility)
	}
	status := ent.StorageStatusPending
	if req.Status != "" {
		status = ent.ToStorageStatus(req.Status)
	}
	virusScanStatus := ent.StorageVirusScanStatusPending
	if req.VirusScanStatus != "" {
		virusScanStatus = ent.ToStorageVirusScanStatus(req.VirusScanStatus)
	}
	versionNo := req.VersionNo
	if versionNo <= 0 {
		versionNo = 1
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{SchoolID: schoolID, BucketName: req.BucketName, ObjectKey: req.ObjectKey, OriginalName: req.OriginalName, Extension: req.Extension, MIMEType: req.MIMEType, SizeBytes: req.SizeBytes, ChecksumSHA256: req.ChecksumSHA256, ETag: req.ETag, Visibility: visibility, Status: status, VirusScanStatus: virusScanStatus, UploadedByMemberID: uploadedByID, VersionNo: versionNo, ReplacedByStorageID: replacedByID, Metadata: req.Metadata})
	if err != nil {
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrUploadedByMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrReplacedByStorageNotFound) {
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		log.Errf("storages.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	schoolID, err := utils.ParseQueryUUID(ctx.Query("school_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	uploadedByID, err := utils.ParseQueryUUID(ctx.Query("uploaded_by_member_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	var status *ent.StorageStatus
	if raw := ctx.Query("status"); raw != "" {
		v := ent.ToStorageStatus(raw)
		status = &v
	}
	var visibility *ent.StorageVisibility
	if raw := ctx.Query("visibility"); raw != "" {
		v := ent.ToStorageVisibility(raw)
		visibility = &v
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListInput{SchoolID: schoolID, UploadedByMemberID: uploadedByID, Status: status, Visibility: visibility})
	if err != nil {
		log.Errf("storages.list.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		log.Errf("storages.get.error: %v", err)
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
	schoolID, uploadedByID, replacedByID, ok := parseCreateUpdateFields(ctx, req.SchoolID, req.UploadedByMemberID, req.ReplacedByStorageID)
	if !ok {
		return
	}
	versionNo := req.VersionNo
	if versionNo <= 0 {
		versionNo = 1
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{SchoolID: schoolID, BucketName: req.BucketName, ObjectKey: req.ObjectKey, OriginalName: req.OriginalName, Extension: req.Extension, MIMEType: req.MIMEType, SizeBytes: req.SizeBytes, ChecksumSHA256: req.ChecksumSHA256, ETag: req.ETag, Visibility: ent.ToStorageVisibility(req.Visibility), Status: ent.ToStorageStatus(req.Status), VirusScanStatus: ent.ToStorageVirusScanStatus(req.VirusScanStatus), UploadedByMemberID: uploadedByID, VersionNo: versionNo, ReplacedByStorageID: replacedByID, Metadata: req.Metadata})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrUploadedByMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrReplacedByStorageNotFound) {
			base.ValidateFailed(ctx, ci18n.StorageNotFound, nil)
			return
		}
		log.Errf("storages.update.error: %v", err)
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
		log.Errf("storages.delete.error: %v", err)
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

func parseCreateUpdateFields(ctx *gin.Context, schoolIDRaw string, uploadedByRaw *string, replacedByRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, bool) {
	schoolID, err := uuid.Parse(schoolIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	uploadedByID, err := utils.ParseUUIDPtr(uploadedByRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	replacedByID, err := utils.ParseUUIDPtr(replacedByRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, false
	}
	return schoolID, uploadedByID, replacedByID, true
}

func toResponse(item *ent.Storage) response {
	return response{ID: item.ID.String(), SchoolID: item.SchoolID.String(), BucketName: item.BucketName, ObjectKey: item.ObjectKey, OriginalName: item.OriginalName, Extension: item.Extension, MIMEType: item.MIMEType, SizeBytes: item.SizeBytes, ChecksumSHA256: item.ChecksumSHA256, ETag: item.ETag, Visibility: string(item.Visibility), Status: string(item.Status), VirusScanStatus: string(item.VirusScanStatus), UploadedByMemberID: utils.UUIDToStringPtr(item.UploadedByMemberID), VersionNo: item.VersionNo, ReplacedByStorageID: utils.UUIDToStringPtr(item.ReplacedByStorageID), Metadata: item.Metadata, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout), UpdatedAt: item.UpdatedAt.UTC().Format(dateTimeLayout)}
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
		if strings.Contains(strings.ToLower(item.BucketName), needle) ||
			strings.Contains(strings.ToLower(item.ObjectKey), needle) ||
			(item.OriginalName != nil && strings.Contains(strings.ToLower(*item.OriginalName), needle)) {
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
		case "updated_at":
			less = items[i].UpdatedAt < items[j].UpdatedAt
		case "size_bytes":
			less = items[i].SizeBytes < items[j].SizeBytes
		case "version_no":
			less = items[i].VersionNo < items[j].VersionNo
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
