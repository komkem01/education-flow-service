package teachers

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
)

const dateLayoutOnly = "2006-01-02"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type teacherURIRequest struct {
	ID string `uri:"id" binding:"required"`
}

type createTeacherRequest struct {
	MemberID                string  `json:"member_id" binding:"required,uuid"`
	GenderID                *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID                *string `json:"prefix_id" binding:"omitempty,uuid"`
	TeacherCode             *string `json:"teacher_code" binding:"omitempty,max=255"`
	FirstName               *string `json:"first_name" binding:"omitempty,max=255"`
	LastName                *string `json:"last_name" binding:"omitempty,max=255"`
	CitizenID               *string `json:"citizen_id" binding:"omitempty,max=13"`
	Phone                   *string `json:"phone" binding:"omitempty,max=50"`
	CurrentPosition         *string `json:"current_position" binding:"omitempty,max=255"`
	CurrentAcademicStanding *string `json:"current_academic_standing" binding:"omitempty,max=255"`
	Department              *string `json:"department" binding:"omitempty,max=255"`
	StartDate               *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	IsActive                *bool   `json:"is_active"`
}

type updateTeacherRequest = createTeacherRequest

type registerTeacherRequest struct {
	SchoolID                string  `json:"school_id" binding:"required,uuid"`
	Email                   string  `json:"email" binding:"required,email,max=255"`
	Password                string  `json:"password" binding:"required,min=6,max=255"`
	GenderID                *string `json:"gender_id" binding:"omitempty,uuid"`
	PrefixID                *string `json:"prefix_id" binding:"omitempty,uuid"`
	TeacherCode             *string `json:"teacher_code" binding:"omitempty,max=255"`
	FirstName               *string `json:"first_name" binding:"omitempty,max=255"`
	LastName                *string `json:"last_name" binding:"omitempty,max=255"`
	CitizenID               *string `json:"citizen_id" binding:"omitempty,max=13"`
	Phone                   *string `json:"phone" binding:"omitempty,max=50"`
	CurrentPosition         *string `json:"current_position" binding:"omitempty,max=255"`
	CurrentAcademicStanding *string `json:"current_academic_standing" binding:"omitempty,max=255"`
	Department              *string `json:"department" binding:"omitempty,max=255"`
	StartDate               *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	IsActive                *bool   `json:"is_active"`
}

type teacherResponse struct {
	ID                      string  `json:"id"`
	MemberID                string  `json:"member_id"`
	GenderID                *string `json:"gender_id"`
	PrefixID                *string `json:"prefix_id"`
	TeacherCode             *string `json:"teacher_code"`
	FirstName               *string `json:"first_name"`
	LastName                *string `json:"last_name"`
	CitizenID               *string `json:"citizen_id"`
	Phone                   *string `json:"phone"`
	CurrentPosition         *string `json:"current_position"`
	CurrentAcademicStanding *string `json:"current_academic_standing"`
	Department              *string `json:"department"`
	StartDate               *string `json:"start_date"`
	IsActive                bool    `json:"is_active"`
}

type registerTeacherResponse struct {
	Member memberRegisterResponse `json:"member"`
	Teacher teacherResponse       `json:"teacher"`
}

type memberRegisterResponse struct {
	ID       string `json:"id"`
	SchoolID string `json:"school_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

func (c *Controller) Register(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req registerTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	schoolID, err := uuid.Parse(req.SchoolID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	genderID, err := parseUUIDPtr(req.GenderID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	prefixID, err := parseUUIDPtr(req.PrefixID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	startDate, err := parseDatePtr(req.StartDate)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	member, teacher, err := c.svc.Register(ctx.Request.Context(), &RegisterTeacherInput{
		SchoolID:                schoolID,
		Email:                   req.Email,
		Password:                req.Password,
		GenderID:                genderID,
		PrefixID:                prefixID,
		TeacherCode:             req.TeacherCode,
		FirstName:               req.FirstName,
		LastName:                req.LastName,
		CitizenID:               req.CitizenID,
		Phone:                   req.Phone,
		CurrentPosition:         req.CurrentPosition,
		CurrentAcademicStanding: req.CurrentAcademicStanding,
		Department:              req.Department,
		StartDate:               startDate,
		IsActive:                isActive,
	})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.MemberEmailDuplicate, nil)
			return
		}
		log.Errf("teachers.register.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}

	base.Success(ctx, registerTeacherResponse{
		Member: toMemberRegisterResponse(member),
		Teacher: toTeacherResponse(teacher),
	})
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var req createTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, startDate, ok := parseTeacherCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.StartDate)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	teacher, err := c.svc.Create(ctx.Request.Context(), &CreateTeacherInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, TeacherCode: req.TeacherCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, CurrentPosition: req.CurrentPosition, CurrentAcademicStanding: req.CurrentAcademicStanding, Department: req.Department, StartDate: startDate, IsActive: isActive})
	if err != nil {
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.TeacherCodeDuplicate, nil)
			return
		}
		log.Errf("teachers.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	var memberID *uuid.UUID
	if memberIDRaw := ctx.Query("member_id"); memberIDRaw != "" {
		parsedMemberID, err := uuid.Parse(memberIDRaw)
		if err != nil {
			base.BadRequest(ctx, ci18n.InvalidID, nil)
			return
		}
		memberID = &parsedMemberID
	}
	onlyActive, err := strconv.ParseBool(ctx.DefaultQuery("only_active", "false"))
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	teachers, err := c.svc.List(ctx.Request.Context(), &ListTeachersInput{MemberID: memberID, OnlyActive: onlyActive})
	if err != nil {
		log.Errf("teachers.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	response := make([]teacherResponse, 0, len(teachers))
	for _, teacher := range teachers {
		response = append(response, toTeacherResponse(teacher))
	}
	base.Success(ctx, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	teacher, err := c.svc.GetByID(ctx.Request.Context(), teacherID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherNotFound, nil)
			return
		}
		log.Errf("teachers.get.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher))
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	var req updateTeacherRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	memberID, genderID, prefixID, startDate, ok := parseTeacherCreateUpdateFields(ctx, req.MemberID, req.GenderID, req.PrefixID, req.StartDate)
	if !ok {
		return
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	teacher, err := c.svc.UpdateByID(ctx.Request.Context(), teacherID, &UpdateTeacherInput{MemberID: memberID, GenderID: genderID, PrefixID: prefixID, TeacherCode: req.TeacherCode, FirstName: req.FirstName, LastName: req.LastName, CitizenID: req.CitizenID, Phone: req.Phone, CurrentPosition: req.CurrentPosition, CurrentAcademicStanding: req.CurrentAcademicStanding, Department: req.Department, StartDate: startDate, IsActive: isActive})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.TeacherNotFound, nil)
			return
		}
		if isDuplicateKeyError(err) {
			base.ValidateFailed(ctx, ci18n.TeacherCodeDuplicate, nil)
			return
		}
		log.Errf("teachers.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toTeacherResponse(teacher))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	teacherID, ok := parseTeacherID(ctx)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), teacherID); err != nil {
		log.Errf("teachers.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": teacherID.String()})
}

func parseTeacherID(ctx *gin.Context) (uuid.UUID, bool) {
	var req teacherURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, false
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, false
	}
	return id, true
}

func parseTeacherCreateUpdateFields(ctx *gin.Context, memberIDRaw string, genderIDRaw *string, prefixIDRaw *string, startDateRaw *string) (uuid.UUID, *uuid.UUID, *uuid.UUID, *time.Time, bool) {
	memberID, err := uuid.Parse(memberIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	genderID, err := parseUUIDPtr(genderIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	prefixID, err := parseUUIDPtr(prefixIDRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	startDate, err := parseDatePtr(startDateRaw)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return uuid.Nil, nil, nil, nil, false
	}
	return memberID, genderID, prefixID, startDate, true
}

func parseUUIDPtr(raw *string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func parseDatePtr(raw *string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(dateLayoutOnly, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func toTeacherResponse(teacher *ent.MemberTeacher) teacherResponse {
	return teacherResponse{ID: teacher.ID.String(), MemberID: teacher.MemberID.String(), GenderID: uuidToStringPtr(teacher.GenderID), PrefixID: uuidToStringPtr(teacher.PrefixID), TeacherCode: teacher.TeacherCode, FirstName: teacher.FirstName, LastName: teacher.LastName, CitizenID: teacher.CitizenID, Phone: teacher.Phone, CurrentPosition: teacher.CurrentPosition, CurrentAcademicStanding: teacher.CurrentAcademicStanding, Department: teacher.Department, StartDate: dateToStringPtr(teacher.StartDate), IsActive: teacher.IsActive}
}

func uuidToStringPtr(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	parsed := value.String()
	return &parsed
}

func dateToStringPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	parsed := value.Format(dateLayoutOnly)
	return &parsed
}

func toMemberRegisterResponse(member *ent.Member) memberRegisterResponse {
	return memberRegisterResponse{
		ID:       member.ID.String(),
		SchoolID: member.SchoolID.String(),
		Email:    member.Email,
		Role:     string(member.Role),
		IsActive: member.IsActive,
	}
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
