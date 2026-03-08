package reports

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"education-flow/app/modules/entities/ent"
	"education-flow/internal/config"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	tracer trace.Tracer
	db     *bun.DB
}

type Options struct {
	*config.Config[Config]
	tracer trace.Tracer
	db     *bun.DB
}

type SummaryInput struct {
	SchoolID       uuid.UUID
	AcademicYearID *uuid.UUID
	SemesterNo     *int
}

type SummaryOutput struct {
	TeachersTotal     int64 `json:"teachers_total"`
	TeachersActive    int64 `json:"teachers_active"`
	StudentsTotal     int64 `json:"students_total"`
	StudentsActive    int64 `json:"students_active"`
	TransferInTotal   int64 `json:"transfer_in_total"`
	TransferOutTotal  int64 `json:"transfer_out_total"`
	PromotionPending  int64 `json:"promotion_pending_total"`
	SubjectsTotal     int64 `json:"subjects_total"`
	CoursesTotal      int64 `json:"courses_total"`
	GradeRecordsTotal int64 `json:"grade_records_total"`
	GradePassTotal    int64 `json:"grade_pass_total"`
	GradeFailTotal    int64 `json:"grade_fail_total"`
	AttendanceTotal   int64 `json:"attendance_total"`
	AttendancePresent int64 `json:"attendance_present"`
	AttendanceAbsent  int64 `json:"attendance_absent"`
	BehaviorTotal     int64 `json:"behavior_total"`
	BehaviorGood      int64 `json:"behavior_good"`
	BehaviorBad       int64 `json:"behavior_bad"`
	DocumentPending   int64 `json:"document_pending_total"`
	ParentNotifyToday int64 `json:"parent_notifications_today"`
}

type AcademicYearFilter struct {
	ID    uuid.UUID `json:"id"`
	Year  string    `json:"year"`
	Term  string    `json:"term"`
	Label string    `json:"label"`
}

type FilterOutput struct {
	AcademicYears []AcademicYearFilter `json:"academic_years"`
	Semesters     []int                `json:"semesters"`
}

type ApprovalType string

const (
	ApprovalTypeTeacherProfile ApprovalType = "teacher_profile_request"
	ApprovalTypeTeacherLeave   ApprovalType = "teacher_leave"
	ApprovalTypeInventory      ApprovalType = "inventory_request"
)

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

type ApprovalItem struct {
	Type          ApprovalType   `json:"type"`
	ID            uuid.UUID      `json:"id"`
	RequesterName string         `json:"requester_name"`
	RequesterRole string         `json:"requester_role"`
	Title         string         `json:"title"`
	Detail        string         `json:"detail"`
	Status        ApprovalStatus `json:"status"`
	Comment       *string        `json:"comment"`
	CreatedAt     time.Time      `json:"created_at"`
}

type ListApprovalsInput struct {
	SchoolID uuid.UUID
	Type     *ApprovalType
	Status   *ApprovalStatus
}

type UpdateApprovalInput struct {
	SchoolID uuid.UUID
	Type     ApprovalType
	ID       uuid.UUID
	Status   ApprovalStatus
	Comment  *string
}

type RoleMemberItem struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	ActiveRole string    `json:"active_role"`
	Name       string    `json:"name"`
	Roles      []string  `json:"roles"`
}

type roleMemberRow struct {
	ID         uuid.UUID `bun:"id"`
	Email      string    `bun:"email"`
	ActiveRole string    `bun:"active_role"`
	Name       *string   `bun:"name"`
	Role       *string   `bun:"role"`
}

type teacherProfileApprovalRow struct {
	ID        uuid.UUID `bun:"id"`
	Status    string    `bun:"status"`
	Reason    *string   `bun:"reason"`
	Comment   *string   `bun:"comment"`
	CreatedAt time.Time `bun:"created_at"`
	FirstName *string   `bun:"first_name"`
	LastName  *string   `bun:"last_name"`
}

type teacherLeaveApprovalRow struct {
	ID        uuid.UUID  `bun:"id"`
	Status    string     `bun:"status"`
	Reason    *string    `bun:"reason"`
	Type      string     `bun:"type"`
	StartDate *time.Time `bun:"start_date"`
	EndDate   *time.Time `bun:"end_date"`
	CreatedAt time.Time  `bun:"created_at"`
	FirstName *string    `bun:"first_name"`
	LastName  *string    `bun:"last_name"`
}

type inventoryApprovalRow struct {
	ID        uuid.UUID `bun:"id"`
	Status    string    `bun:"status"`
	Reason    *string   `bun:"reason"`
	Quantity  *int      `bun:"quantity"`
	CreatedAt time.Time `bun:"created_at"`
	Email     string    `bun:"email"`
	Role      string    `bun:"role"`
}

func newService(opt *Options) *Service {
	return &Service{tracer: opt.tracer, db: opt.db}
}

func (s *Service) ListFilters(ctx context.Context) (*FilterOutput, error) {
	rows := make([]AcademicYearFilter, 0)
	if err := s.db.NewSelect().
		TableExpr("academic_years AS acy").
		ColumnExpr("acy.id").
		ColumnExpr("acy.year").
		ColumnExpr("acy.term").
		ColumnExpr("(acy.year || ' / เทอม ' || acy.term) AS label").
		Where("acy.is_active = true").
		OrderExpr("acy.start_date DESC").
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	return &FilterOutput{AcademicYears: rows, Semesters: []int{1, 2}}, nil
}

func (s *Service) GetSummary(ctx context.Context, input *SummaryInput) (*SummaryOutput, error) {
	out := &SummaryOutput{}
	if err := s.db.NewSelect().
		TableExpr("member_teachers AS mtr").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("mem.school_id = ?", input.SchoolID).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.TeachersTotal); err != nil {
		return nil, err
	}
	if err := s.db.NewSelect().
		TableExpr("member_teachers AS mtr").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("mem.school_id = ?", input.SchoolID).
		Where("mtr.is_active = true").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.TeachersActive); err != nil {
		return nil, err
	}

	if err := s.db.NewSelect().
		TableExpr("member_students AS mst").
		Join("JOIN members AS mem ON mem.id = mst.member_id").
		Where("mem.school_id = ?", input.SchoolID).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.StudentsTotal); err != nil {
		return nil, err
	}
	if err := s.db.NewSelect().
		TableExpr("member_students AS mst").
		Join("JOIN members AS mem ON mem.id = mst.member_id").
		Where("mem.school_id = ?", input.SchoolID).
		Where("mst.is_active = true").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.StudentsActive); err != nil {
		return nil, err
	}

	transferInQuery := s.db.NewSelect().
		TableExpr("student_enrollments AS sen").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("sen.status = ?", ent.StudentEnrollmentStatusActive)
	applyAcademicFilter(transferInQuery, input.AcademicYearID, input.SemesterNo)
	if err := transferInQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.TransferInTotal); err != nil {
		return nil, err
	}

	transferOutQuery := s.db.NewSelect().
		TableExpr("student_enrollments AS sen").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("sen.status = ?", ent.StudentEnrollmentStatusDropped)
	applyAcademicFilter(transferOutQuery, input.AcademicYearID, input.SemesterNo)
	if err := transferOutQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.TransferOutTotal); err != nil {
		return nil, err
	}

	promotionPendingQuery := s.db.NewSelect().
		TableExpr("student_enrollments AS sen").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("sen.status = ?", ent.StudentEnrollmentStatusIncomplete)
	applyAcademicFilter(promotionPendingQuery, input.AcademicYearID, input.SemesterNo)
	if err := promotionPendingQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.PromotionPending); err != nil {
		return nil, err
	}

	if err := s.db.NewSelect().
		TableExpr("subjects AS sub").
		Where("sub.school_id = ?", input.SchoolID).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.SubjectsTotal); err != nil {
		return nil, err
	}

	courseQuery := s.db.NewSelect().
		TableExpr("subject_assignments AS sas").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID)
	applyAcademicFilter(courseQuery, input.AcademicYearID, input.SemesterNo)
	if err := courseQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.CoursesTotal); err != nil {
		return nil, err
	}

	gradeQuery := s.db.NewSelect().
		TableExpr("grade_records AS grr").
		Join("JOIN student_enrollments AS sen ON sen.id = grr.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID)
	applyAcademicFilter(gradeQuery, input.AcademicYearID, input.SemesterNo)
	if err := gradeQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.GradeRecordsTotal); err != nil {
		return nil, err
	}

	gradePassQuery := s.db.NewSelect().
		TableExpr("grade_records AS grr").
		Join("JOIN student_enrollments AS sen ON sen.id = grr.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("grr.score >= 50")
	applyAcademicFilter(gradePassQuery, input.AcademicYearID, input.SemesterNo)
	if err := gradePassQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.GradePassTotal); err != nil {
		return nil, err
	}

	gradeFailQuery := s.db.NewSelect().
		TableExpr("grade_records AS grr").
		Join("JOIN student_enrollments AS sen ON sen.id = grr.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("grr.score < 50")
	applyAcademicFilter(gradeFailQuery, input.AcademicYearID, input.SemesterNo)
	if err := gradeFailQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.GradeFailTotal); err != nil {
		return nil, err
	}

	attendanceQuery := s.db.NewSelect().
		TableExpr("student_attendance_logs AS sal").
		Join("JOIN student_enrollments AS sen ON sen.id = sal.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID)
	applyAcademicFilter(attendanceQuery, input.AcademicYearID, input.SemesterNo)
	if err := attendanceQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.AttendanceTotal); err != nil {
		return nil, err
	}

	attendancePresentQuery := s.db.NewSelect().
		TableExpr("student_attendance_logs AS sal").
		Join("JOIN student_enrollments AS sen ON sen.id = sal.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("sal.status = ?", ent.StudentAttendanceStatusPresent)
	applyAcademicFilter(attendancePresentQuery, input.AcademicYearID, input.SemesterNo)
	if err := attendancePresentQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.AttendancePresent); err != nil {
		return nil, err
	}

	attendanceAbsentQuery := s.db.NewSelect().
		TableExpr("student_attendance_logs AS sal").
		Join("JOIN student_enrollments AS sen ON sen.id = sal.enrollment_id").
		Join("JOIN subject_assignments AS sas ON sas.id = sen.subject_assignment_id").
		Join("JOIN classrooms AS cls ON cls.id = sas.classroom_id").
		Where("cls.school_id = ?", input.SchoolID).
		Where("sal.status = ?", ent.StudentAttendanceStatusAbsent)
	applyAcademicFilter(attendanceAbsentQuery, input.AcademicYearID, input.SemesterNo)
	if err := attendanceAbsentQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.AttendanceAbsent); err != nil {
		return nil, err
	}

	behaviorQuery := s.db.NewSelect().
		TableExpr("student_behaviors AS stb").
		Where("stb.school_id = ?", input.SchoolID).
		Where("stb.is_active = true")
	if err := behaviorQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.BehaviorTotal); err != nil {
		return nil, err
	}

	behaviorGoodQuery := s.db.NewSelect().
		TableExpr("student_behaviors AS stb").
		Where("stb.school_id = ?", input.SchoolID).
		Where("stb.is_active = true").
		Where("stb.behavior_type = ?", ent.StudentBehaviorTypeGood)
	if err := behaviorGoodQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.BehaviorGood); err != nil {
		return nil, err
	}

	behaviorBadQuery := s.db.NewSelect().
		TableExpr("student_behaviors AS stb").
		Where("stb.school_id = ?", input.SchoolID).
		Where("stb.is_active = true").
		Where("stb.behavior_type = ?", ent.StudentBehaviorTypeBad)
	if err := behaviorBadQuery.ColumnExpr("COUNT(*)").Scan(ctx, &out.BehaviorBad); err != nil {
		return nil, err
	}

	if err := s.db.NewSelect().
		TableExpr("document_tracking AS dtk").
		Where("dtk.school_id = ?", input.SchoolID).
		Where("dtk.status <> ?", ent.DocumentTrackingStatusProcessed).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.DocumentPending); err != nil {
		return nil, err
	}

	if err := s.db.NewSelect().
		TableExpr("school_announcements AS sca").
		Where("sca.school_id = ?", input.SchoolID).
		Where("sca.target_role = ?", ent.MemberRoleParent).
		Where("DATE(sca.created_at) = CURRENT_DATE").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &out.ParentNotifyToday); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *Service) ListApprovals(ctx context.Context, input *ListApprovalsInput) ([]*ApprovalItem, error) {
	items := make([]*ApprovalItem, 0)

	if input.Type == nil || *input.Type == ApprovalTypeTeacherProfile {
		rows, err := s.listTeacherProfileApprovals(ctx, input.SchoolID, input.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, rows...)
	}
	if input.Type == nil || *input.Type == ApprovalTypeTeacherLeave {
		rows, err := s.listTeacherLeaveApprovals(ctx, input.SchoolID, input.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, rows...)
	}
	if input.Type == nil || *input.Type == ApprovalTypeInventory {
		rows, err := s.listInventoryApprovals(ctx, input.SchoolID, input.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, rows...)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})

	return items, nil
}

func (s *Service) GetApprovalByTypeAndID(ctx context.Context, schoolID uuid.UUID, approvalType ApprovalType, id uuid.UUID) (*ApprovalItem, error) {
	list, err := s.ListApprovals(ctx, &ListApprovalsInput{SchoolID: schoolID, Type: &approvalType})
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (s *Service) UpdateApproval(ctx context.Context, input *UpdateApprovalInput) error {
	switch input.Type {
	case ApprovalTypeTeacherProfile:
		ok, err := s.teacherProfileBelongsToSchool(ctx, input.ID, input.SchoolID)
		if err != nil {
			return err
		}
		if !ok {
			return sql.ErrNoRows
		}
		_, err = s.db.NewUpdate().
			TableExpr("teacher_profile_requests").
			Set("status = ?", input.Status).
			Set("comment = ?", trimStringPtr(input.Comment)).
			Set("processed_at = current_timestamp").
			Where("id = ?", input.ID).
			Exec(ctx)
		return err
	case ApprovalTypeTeacherLeave:
		ok, err := s.teacherLeaveBelongsToSchool(ctx, input.ID, input.SchoolID)
		if err != nil {
			return err
		}
		if !ok {
			return sql.ErrNoRows
		}
		_, err = s.db.NewUpdate().
			TableExpr("teacher_leave_logs").
			Set("status = ?", input.Status).
			Where("id = ?", input.ID).
			Exec(ctx)
		return err
	case ApprovalTypeInventory:
		ok, err := s.inventoryRequestBelongsToSchool(ctx, input.ID, input.SchoolID)
		if err != nil {
			return err
		}
		if !ok {
			return sql.ErrNoRows
		}
		_, err = s.db.NewUpdate().
			TableExpr("inventory_requests").
			Set("status = ?", input.Status).
			Where("id = ?", input.ID).
			Exec(ctx)
		return err
	default:
		return fmt.Errorf("unsupported approval type: %s", input.Type)
	}
}

func (s *Service) ListRoleMembers(ctx context.Context, schoolID uuid.UUID) ([]*RoleMemberItem, error) {
	rows := make([]roleMemberRow, 0)
	if err := s.db.NewSelect().
		TableExpr("members AS mem").
		ColumnExpr("mem.id").
		ColumnExpr("mem.email").
		ColumnExpr(`COALESCE((
			SELECT mr.role
			FROM member_roles AS mr
			WHERE mr.member_id = mem.id
			ORDER BY CASE mr.role
				WHEN 'admin' THEN 1
				WHEN 'staff' THEN 2
				WHEN 'teacher' THEN 3
				WHEN 'parent' THEN 4
				WHEN 'student' THEN 5
				ELSE 99
			END, mr.created_at ASC
			LIMIT 1
		), 'admin') AS active_role`).
		ColumnExpr("COALESCE(mad.first_name || ' ' || mad.last_name, mst.first_name || ' ' || mst.last_name, mtr.first_name || ' ' || mtr.last_name, mss.first_name || ' ' || mss.last_name, mpa.first_name || ' ' || mpa.last_name) AS name").
		ColumnExpr("mrl.role").
		Join("LEFT JOIN member_admins AS mad ON mad.member_id = mem.id").
		Join("LEFT JOIN member_staffs AS mst ON mst.member_id = mem.id").
		Join("LEFT JOIN member_teachers AS mtr ON mtr.member_id = mem.id").
		Join("LEFT JOIN member_students AS mss ON mss.member_id = mem.id").
		Join("LEFT JOIN member_parents AS mpa ON mpa.member_id = mem.id").
		Join("LEFT JOIN member_roles AS mrl ON mrl.member_id = mem.id").
		Where("mem.school_id = ?", schoolID).
		OrderExpr("mem.created_at DESC").
		Scan(ctx, &rows); err != nil {
		return nil, err
	}

	itemsByID := make(map[uuid.UUID]*RoleMemberItem, len(rows))
	for _, row := range rows {
		item, ok := itemsByID[row.ID]
		if !ok {
			name := strings.TrimSpace(valueOrEmpty(row.Name))
			if name == "" {
				name = row.Email
			}
			item = &RoleMemberItem{
				ID:         row.ID,
				Email:      row.Email,
				ActiveRole: row.ActiveRole,
				Name:       name,
				Roles:      make([]string, 0, 2),
			}
			itemsByID[row.ID] = item
		}

		if row.Role == nil {
			continue
		}

		role := strings.TrimSpace(*row.Role)
		if role == "" {
			continue
		}

		exists := false
		for _, current := range item.Roles {
			if current == role {
				exists = true
				break
			}
		}
		if !exists {
			item.Roles = append(item.Roles, role)
		}
	}

	result := make([]*RoleMemberItem, 0, len(itemsByID))
	for _, row := range rows {
		item := itemsByID[row.ID]
		if item == nil {
			continue
		}
		result = append(result, item)
		delete(itemsByID, row.ID)
	}

	for _, item := range result {
		sort.Strings(item.Roles)
	}

	return result, nil
}

func (s *Service) listTeacherProfileApprovals(ctx context.Context, schoolID uuid.UUID, status *ApprovalStatus) ([]*ApprovalItem, error) {
	rows := make([]teacherProfileApprovalRow, 0)
	query := s.db.NewSelect().
		TableExpr("teacher_profile_requests AS tpr").
		ColumnExpr("tpr.id").
		ColumnExpr("tpr.status").
		ColumnExpr("tpr.reason").
		ColumnExpr("tpr.comment").
		ColumnExpr("tpr.created_at").
		ColumnExpr("mtr.first_name").
		ColumnExpr("mtr.last_name").
		Join("JOIN member_teachers AS mtr ON mtr.id = tpr.teacher_id").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("mem.school_id = ?", schoolID)
	if status != nil {
		query = query.Where("tpr.status = ?", *status)
	}
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, err
	}

	items := make([]*ApprovalItem, 0, len(rows))
	for _, row := range rows {
		detail := strings.TrimSpace(valueOrEmpty(row.Reason))
		if detail == "" {
			detail = "คำขอแก้ไขข้อมูลโปรไฟล์ครู"
		}
		items = append(items, &ApprovalItem{
			Type:          ApprovalTypeTeacherProfile,
			ID:            row.ID,
			RequesterName: fullName(row.FirstName, row.LastName),
			RequesterRole: "teacher",
			Title:         "คำขอแก้ไขโปรไฟล์ครู",
			Detail:        detail,
			Status:        ApprovalStatus(row.Status),
			Comment:       row.Comment,
			CreatedAt:     row.CreatedAt,
		})
	}

	return items, nil
}

func (s *Service) listTeacherLeaveApprovals(ctx context.Context, schoolID uuid.UUID, status *ApprovalStatus) ([]*ApprovalItem, error) {
	rows := make([]teacherLeaveApprovalRow, 0)
	query := s.db.NewSelect().
		TableExpr("teacher_leave_logs AS tll").
		ColumnExpr("tll.id").
		ColumnExpr("tll.status").
		ColumnExpr("tll.reason").
		ColumnExpr("tll.type").
		ColumnExpr("tll.start_date").
		ColumnExpr("tll.end_date").
		ColumnExpr("tll.created_at").
		ColumnExpr("mtr.first_name").
		ColumnExpr("mtr.last_name").
		Join("JOIN member_teachers AS mtr ON mtr.id = tll.teacher_id").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("mem.school_id = ?", schoolID)
	if status != nil {
		query = query.Where("tll.status = ?", *status)
	}
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, err
	}

	items := make([]*ApprovalItem, 0, len(rows))
	for _, row := range rows {
		parts := make([]string, 0, 3)
		parts = append(parts, fmt.Sprintf("ประเภทการลา: %s", row.Type))
		if row.StartDate != nil {
			parts = append(parts, fmt.Sprintf("เริ่ม: %s", row.StartDate.Format("2006-01-02")))
		}
		if row.EndDate != nil {
			parts = append(parts, fmt.Sprintf("สิ้นสุด: %s", row.EndDate.Format("2006-01-02")))
		}
		if reason := strings.TrimSpace(valueOrEmpty(row.Reason)); reason != "" {
			parts = append(parts, reason)
		}

		items = append(items, &ApprovalItem{
			Type:          ApprovalTypeTeacherLeave,
			ID:            row.ID,
			RequesterName: fullName(row.FirstName, row.LastName),
			RequesterRole: "teacher",
			Title:         "คำขอลางานครู",
			Detail:        strings.Join(parts, " | "),
			Status:        ApprovalStatus(row.Status),
			CreatedAt:     row.CreatedAt,
		})
	}

	return items, nil
}

func (s *Service) listInventoryApprovals(ctx context.Context, schoolID uuid.UUID, status *ApprovalStatus) ([]*ApprovalItem, error) {
	rows := make([]inventoryApprovalRow, 0)
	query := s.db.NewSelect().
		TableExpr("inventory_requests AS ivr").
		ColumnExpr("ivr.id").
		ColumnExpr("ivr.status").
		ColumnExpr("ivr.reason").
		ColumnExpr("ivr.quantity").
		ColumnExpr("ivr.created_at").
		ColumnExpr("mem.email").
		ColumnExpr(`COALESCE((
			SELECT mr.role
			FROM member_roles AS mr
			WHERE mr.member_id = mem.id
			ORDER BY CASE mr.role
				WHEN 'admin' THEN 1
				WHEN 'staff' THEN 2
				WHEN 'teacher' THEN 3
				WHEN 'parent' THEN 4
				WHEN 'student' THEN 5
				ELSE 99
			END, mr.created_at ASC
			LIMIT 1
		), 'admin') AS role`).
		Join("JOIN members AS mem ON mem.id = ivr.requester_member_id").
		Where("mem.school_id = ?", schoolID)
	if status != nil {
		query = query.Where("ivr.status = ?", *status)
	}
	if err := query.Scan(ctx, &rows); err != nil {
		return nil, err
	}

	items := make([]*ApprovalItem, 0, len(rows))
	for _, row := range rows {
		detail := strings.TrimSpace(valueOrEmpty(row.Reason))
		if row.Quantity != nil {
			qty := fmt.Sprintf("จำนวน: %d", *row.Quantity)
			if detail == "" {
				detail = qty
			} else {
				detail = qty + " | " + detail
			}
		}
		if detail == "" {
			detail = "คำขอเบิกอุปกรณ์"
		}

		items = append(items, &ApprovalItem{
			Type:          ApprovalTypeInventory,
			ID:            row.ID,
			RequesterName: row.Email,
			RequesterRole: row.Role,
			Title:         "คำขอเบิกพัสดุ",
			Detail:        detail,
			Status:        ApprovalStatus(row.Status),
			CreatedAt:     row.CreatedAt,
		})
	}

	return items, nil
}

func (s *Service) teacherProfileBelongsToSchool(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		TableExpr("teacher_profile_requests AS tpr").
		Join("JOIN member_teachers AS mtr ON mtr.id = tpr.teacher_id").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("tpr.id = ?", id).
		Where("mem.school_id = ?", schoolID).
		Exists(ctx)
}

func (s *Service) teacherLeaveBelongsToSchool(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		TableExpr("teacher_leave_logs AS tll").
		Join("JOIN member_teachers AS mtr ON mtr.id = tll.teacher_id").
		Join("JOIN members AS mem ON mem.id = mtr.member_id").
		Where("tll.id = ?", id).
		Where("mem.school_id = ?", schoolID).
		Exists(ctx)
}

func (s *Service) inventoryRequestBelongsToSchool(ctx context.Context, id uuid.UUID, schoolID uuid.UUID) (bool, error) {
	return s.db.NewSelect().
		TableExpr("inventory_requests AS ivr").
		Join("JOIN members AS mem ON mem.id = ivr.requester_member_id").
		Where("ivr.id = ?", id).
		Where("mem.school_id = ?", schoolID).
		Exists(ctx)
}

func applyAcademicFilter(query *bun.SelectQuery, academicYearID *uuid.UUID, semesterNo *int) {
	if academicYearID != nil {
		query.Where("sas.academic_year_id = ?", *academicYearID)
	}
	if semesterNo != nil {
		query.Where("sas.semester_no = ?", *semesterNo)
	}
}

func fullName(firstName *string, lastName *string) string {
	first := strings.TrimSpace(valueOrEmpty(firstName))
	last := strings.TrimSpace(valueOrEmpty(lastName))
	name := strings.TrimSpace(first + " " + last)
	if name == "" {
		return "ไม่ระบุชื่อ"
	}
	return name
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func trimStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
