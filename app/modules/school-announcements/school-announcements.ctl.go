package schoolannouncements

import (
	"database/sql"
	"errors"
	"sort"
	"strconv"
	"strings"

	"education-flow/app/modules/auth"
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
	SchoolID       string  `json:"school_id" binding:"required,uuid"`
	AuthorMemberID string  `json:"author_member_id" binding:"required,uuid"`
	Title          *string `json:"title" binding:"omitempty,max=255"`
	Content        *string `json:"content"`
	TargetRole     *string `json:"target_role" binding:"omitempty,oneof=student teacher admin staff parent"`
	IsPinned       bool    `json:"is_pinned"`
}

type updateRequest = createRequest

type response struct {
	ID             string  `json:"id"`
	SchoolID       string  `json:"school_id"`
	AuthorMemberID string  `json:"author_member_id"`
	Title          *string `json:"title"`
	Content        *string `json:"content"`
	TargetRole     *string `json:"target_role"`
	IsPinned       bool    `json:"is_pinned"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, authorID, role, ok := parseCreateUpdateFields(ctx, req.SchoolID, req.AuthorMemberID, req.TargetRole)
	if !ok {
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{SchoolID: schoolID, AuthorMemberID: authorID, Title: req.Title, Content: req.Content, TargetRole: role, IsPinned: req.IsPinned})
	if err != nil {
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrAuthorMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrInvalidAuthorRole) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		log.Errf("school-announcements.create.error: %v", err)
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
	if claims, ok := auth.GetClaimsFromGin(ctx); ok {
		if schoolID != nil && *schoolID != claims.SchoolID {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		schoolID = &claims.SchoolID
	}
	var role *ent.MemberRole
	if raw := ctx.Query("target_role"); raw != "" {
		v := ent.ToMemberRole(raw)
		role = &v
	}

	items, err := c.svc.List(ctx.Request.Context(), &ListInput{SchoolID: schoolID, TargetRole: role, OnlyPinned: ctx.Query("only_pinned") == "true"})
	if err != nil {
		log.Errf("school-announcements.list.error: %v", err)
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
			base.ValidateFailed(ctx, ci18n.SchoolAnnouncementNotFound, nil)
			return
		}
		log.Errf("school-announcements.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && item.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
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
	schoolID, authorID, role, ok := parseCreateUpdateFields(ctx, req.SchoolID, req.AuthorMemberID, req.TargetRole)
	if !ok {
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && claims.SchoolID != schoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	existing, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SchoolAnnouncementNotFound, nil)
			return
		}
		log.Errf("school-announcements.get-before-update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && existing.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), id, &UpdateInput{SchoolID: schoolID, AuthorMemberID: authorID, Title: req.Title, Content: req.Content, TargetRole: role, IsPinned: req.IsPinned})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SchoolAnnouncementNotFound, nil)
			return
		}
		if errors.Is(err, ErrSchoolNotFound) {
			base.ValidateFailed(ctx, ci18n.SchoolNotFound, nil)
			return
		}
		if errors.Is(err, ErrAuthorMemberNotFound) {
			base.ValidateFailed(ctx, ci18n.MemberNotFound, nil)
			return
		}
		if errors.Is(err, ErrInvalidAuthorRole) {
			base.Forbidden(ctx, ci18n.Forbidden, nil)
			return
		}
		log.Errf("school-announcements.update.error: %v", err)
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

	item, err := c.svc.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.SchoolAnnouncementNotFound, nil)
			return
		}
		log.Errf("school-announcements.get-before-delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	if claims, ok := auth.GetClaimsFromGin(ctx); ok && item.SchoolID != claims.SchoolID {
		base.Forbidden(ctx, ci18n.Forbidden, nil)
		return
	}

	if err := c.svc.DeleteByID(ctx.Request.Context(), id); err != nil {
		log.Errf("school-announcements.delete.error: %v", err)
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

func parseCreateUpdateFields(ctx *gin.Context, schoolIDRaw string, authorIDRaw string, roleRaw *string) (uuid.UUID, uuid.UUID, *ent.MemberRole, bool) {
	schoolID, err := uuid.Parse(schoolIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, nil, false
	}
	authorID, err := uuid.Parse(authorIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, nil, false
	}
	var role *ent.MemberRole
	if roleRaw != nil {
		v := ent.ToMemberRole(*roleRaw)
		role = &v
	}
	return schoolID, authorID, role, true
}

func toResponse(item *ent.SchoolAnnouncement) response {
	return response{ID: item.ID.String(), SchoolID: item.SchoolID.String(), AuthorMemberID: item.AuthorMemberID.String(), Title: item.Title, Content: item.Content, TargetRole: memberRoleToStringPtr(item.TargetRole), IsPinned: item.IsPinned, CreatedAt: item.CreatedAt.UTC().Format(dateTimeLayout), UpdatedAt: item.UpdatedAt.UTC().Format(dateTimeLayout)}
}

func memberRoleToStringPtr(role *ent.MemberRole) *string {
	if role == nil {
		return nil
	}
	v := string(*role)
	return &v
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
		title := ""
		if item.Title != nil {
			title = strings.ToLower(*item.Title)
		}
		content := ""
		if item.Content != nil {
			content = strings.ToLower(*item.Content)
		}
		if strings.Contains(title, needle) || strings.Contains(content, needle) {
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
		case "title":
			a := ""
			if items[i].Title != nil {
				a = strings.ToLower(*items[i].Title)
			}
			b := ""
			if items[j].Title != nil {
				b = strings.ToLower(*items[j].Title)
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
