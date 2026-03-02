package schoolannouncements

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
			base.ValidateFailed(ctx, ci18n.SchoolAnnouncementNotFound, nil)
			return
		}
		log.Errf("school-announcements.get.error: %v", err)
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
	schoolID, authorID, role, ok := parseCreateUpdateFields(ctx, req.SchoolID, req.AuthorMemberID, req.TargetRole)
	if !ok {
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
