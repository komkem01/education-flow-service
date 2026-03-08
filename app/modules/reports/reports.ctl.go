package reports

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"education-flow/app/modules/auth"
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

type summaryResponse struct {
	Summary *SummaryOutput `json:"summary"`
}

type approvalsResponse struct {
	Items []approvalItemResponse `json:"items"`
}

type roleMembersResponse struct {
	Items []roleMemberResponse `json:"items"`
}

type roleMemberResponse struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	ActiveRole string   `json:"active_role"`
	Name       string   `json:"name"`
	Roles      []string `json:"roles"`
}

type approvalItemResponse struct {
	Type          string  `json:"type"`
	ID            string  `json:"id"`
	RequesterName string  `json:"requester_name"`
	RequesterRole string  `json:"requester_role"`
	Title         string  `json:"title"`
	Detail        string  `json:"detail"`
	Status        string  `json:"status"`
	Comment       *string `json:"comment"`
	CreatedAt     string  `json:"created_at"`
}

type updateApprovalRequest struct {
	Status  string  `json:"status" binding:"required,oneof=pending approved rejected"`
	Comment *string `json:"comment"`
}

func (c *Controller) ListFilters(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	_, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	items, err := c.svc.ListFilters(ctx.Request.Context())
	if err != nil {
		log.Errf("reports.filters.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, items)
}

func (c *Controller) Summary(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	academicYearID, err := utils.ParseQueryUUID(ctx.Query("academic_year_id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var semesterNo *int
	if raw := strings.TrimSpace(ctx.Query("semester_no")); raw != "" {
		value, parseErr := strconv.Atoi(raw)
		if parseErr != nil || (value != 1 && value != 2) {
			base.BadRequest(ctx, ci18n.BadRequest, nil)
			return
		}
		semesterNo = &value
	}

	summary, err := c.svc.GetSummary(ctx.Request.Context(), &SummaryInput{
		SchoolID:       claims.SchoolID,
		AcademicYearID: academicYearID,
		SemesterNo:     semesterNo,
	})
	if err != nil {
		log.Errf("reports.summary.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, summaryResponse{Summary: summary})
}

func (c *Controller) ListApprovals(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	approvalType, err := parseApprovalTypePtr(ctx.Query("type"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	status, err := parseApprovalStatusPtr(ctx.Query("status"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	items, err := c.svc.ListApprovals(ctx.Request.Context(), &ListApprovalsInput{
		SchoolID: claims.SchoolID,
		Type:     approvalType,
		Status:   status,
	})
	if err != nil {
		log.Errf("reports.approvals.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]approvalItemResponse, 0, len(items))
	for _, item := range items {
		response = append(response, approvalItemResponse{
			Type:          string(item.Type),
			ID:            item.ID.String(),
			RequesterName: item.RequesterName,
			RequesterRole: item.RequesterRole,
			Title:         item.Title,
			Detail:        item.Detail,
			Status:        string(item.Status),
			Comment:       item.Comment,
			CreatedAt:     item.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	base.Success(ctx, approvalsResponse{Items: response})
}

func (c *Controller) ListRoleMembers(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	items, err := c.svc.ListRoleMembers(ctx.Request.Context(), claims.SchoolID)
	if err != nil {
		log.Errf("reports.roles-members.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	response := make([]roleMemberResponse, 0, len(items))
	for _, item := range items {
		response = append(response, roleMemberResponse{
			ID:         item.ID.String(),
			Email:      item.Email,
			ActiveRole: item.ActiveRole,
			Name:       item.Name,
			Roles:      item.Roles,
		})
	}

	base.Success(ctx, roleMembersResponse{Items: response})
}

func (c *Controller) GetApproval(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	approvalType, ok := parseApprovalTypeParam(ctx)
	if !ok {
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.GetApprovalByTypeAndID(ctx.Request.Context(), claims.SchoolID, approvalType, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.BadRequest, nil)
			return
		}
		log.Errf("reports.approvals.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, approvalItemResponse{
		Type:          string(item.Type),
		ID:            item.ID.String(),
		RequesterName: item.RequesterName,
		RequesterRole: item.RequesterRole,
		Title:         item.Title,
		Detail:        item.Detail,
		Status:        string(item.Status),
		Comment:       item.Comment,
		CreatedAt:     item.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
	})
}

func (c *Controller) UpdateApproval(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	claims, ok := auth.GetClaimsFromGin(ctx)
	if !ok {
		base.Unauthorized(ctx, ci18n.Unauthorized, nil)
		return
	}

	approvalType, ok := parseApprovalTypeParam(ctx)
	if !ok {
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	var req updateApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	err = c.svc.UpdateApproval(ctx.Request.Context(), &UpdateApprovalInput{
		SchoolID: claims.SchoolID,
		Type:     approvalType,
		ID:       id,
		Status:   ApprovalStatus(req.Status),
		Comment:  req.Comment,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.BadRequest, nil)
			return
		}
		log.Errf("reports.approvals.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, gin.H{"id": id.String(), "type": string(approvalType), "status": req.Status})
}

func parseApprovalTypePtr(raw string) (*ApprovalType, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}
	parsed := ApprovalType(value)
	if parsed != ApprovalTypeTeacherProfile && parsed != ApprovalTypeTeacherLeave && parsed != ApprovalTypeInventory {
		return nil, errors.New("invalid approval type")
	}
	return &parsed, nil
}

func parseApprovalStatusPtr(raw string) (*ApprovalStatus, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}
	parsed := ApprovalStatus(value)
	if parsed != ApprovalStatusPending && parsed != ApprovalStatusApproved && parsed != ApprovalStatusRejected {
		return nil, errors.New("invalid approval status")
	}
	return &parsed, nil
}

func parseApprovalTypeParam(ctx *gin.Context) (ApprovalType, bool) {
	parsed := ApprovalType(strings.TrimSpace(ctx.Param("type")))
	if parsed != ApprovalTypeTeacherProfile && parsed != ApprovalTypeTeacherLeave && parsed != ApprovalTypeInventory {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return "", false
	}
	return parsed, true
}
