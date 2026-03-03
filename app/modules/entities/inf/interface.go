package entitiesinf

import (
	"context"

	"education-flow/app/modules/entities/ent"

	"github.com/google/uuid"
)

// ObjectEntity defines the interface for object entity operations such as create, retrieve, update, and soft delete.
type ExampleEntity interface {
	CreateExample(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
	GetExampleByID(ctx context.Context, id uuid.UUID) (*ent.Example, error)
	UpdateExampleByID(ctx context.Context, id uuid.UUID, status ent.ExampleStatus) (*ent.Example, error)
	SoftDeleteExampleByID(ctx context.Context, id uuid.UUID) error
	ListExamplesByStatus(ctx context.Context, status ent.ExampleStatus) ([]*ent.Example, error)
}
type ExampleTwoEntity interface {
	CreateExampleTwo(ctx context.Context, userID uuid.UUID) (*ent.Example, error)
}

type StorageEntity interface {
	CreateStorage(ctx context.Context, storage *ent.Storage) (*ent.Storage, error)
	GetStorageByID(ctx context.Context, id uuid.UUID) (*ent.Storage, error)
	GetStorageByObjectKey(ctx context.Context, bucketName, objectKey string) (*ent.Storage, error)
	UpdateStorageByID(ctx context.Context, id uuid.UUID, storage *ent.Storage) (*ent.Storage, error)
	UpdateStorageStatusByID(ctx context.Context, id uuid.UUID, status ent.StorageStatus) (*ent.Storage, error)
	ListStorages(ctx context.Context, schoolID *uuid.UUID, uploadedByMemberID *uuid.UUID, status *ent.StorageStatus, visibility *ent.StorageVisibility) ([]*ent.Storage, error)
	SoftDeleteStorageByID(ctx context.Context, id uuid.UUID) error
	DeleteStorageByID(ctx context.Context, id uuid.UUID) error
	CreateStorageLink(ctx context.Context, link *ent.StorageLink) (*ent.StorageLink, error)
	GetStorageLinkByID(ctx context.Context, id uuid.UUID) (*ent.StorageLink, error)
	UpdateStorageLinkByID(ctx context.Context, id uuid.UUID, link *ent.StorageLink) (*ent.StorageLink, error)
	DeleteStorageLinkByID(ctx context.Context, id uuid.UUID) error
	ListStorageLinks(ctx context.Context, storageID *uuid.UUID, entityType *string, entityID *uuid.UUID) ([]*ent.StorageLink, error)
	ListStorageLinksByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*ent.StorageLink, error)
	DeleteStorageLinksByStorageID(ctx context.Context, storageID uuid.UUID) error
}

type SchoolAnnouncementEntity interface {
	CreateSchoolAnnouncement(ctx context.Context, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error)
	GetSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) (*ent.SchoolAnnouncement, error)
	UpdateSchoolAnnouncementByID(ctx context.Context, id uuid.UUID, announcement *ent.SchoolAnnouncement) (*ent.SchoolAnnouncement, error)
	DeleteSchoolAnnouncementByID(ctx context.Context, id uuid.UUID) error
	ListSchoolAnnouncements(ctx context.Context, schoolID *uuid.UUID, targetRole *ent.MemberRole, onlyPinned bool) ([]*ent.SchoolAnnouncement, error)
}

type SystemAuditLogEntity interface {
	CreateSystemAuditLog(ctx context.Context, auditLog *ent.SystemAuditLog) (*ent.SystemAuditLog, error)
	GetSystemAuditLogByID(ctx context.Context, id uuid.UUID) (*ent.SystemAuditLog, error)
	UpdateSystemAuditLogByID(ctx context.Context, id uuid.UUID, auditLog *ent.SystemAuditLog) (*ent.SystemAuditLog, error)
	DeleteSystemAuditLogByID(ctx context.Context, id uuid.UUID) error
	ListSystemAuditLogs(ctx context.Context, memberID *uuid.UUID, module *string) ([]*ent.SystemAuditLog, error)
}

type DataChangeLogEntity interface {
	CreateDataChangeLog(ctx context.Context, changeLog *ent.DataChangeLog) (*ent.DataChangeLog, error)
	GetDataChangeLogByID(ctx context.Context, id uuid.UUID) (*ent.DataChangeLog, error)
	UpdateDataChangeLogByID(ctx context.Context, id uuid.UUID, changeLog *ent.DataChangeLog) (*ent.DataChangeLog, error)
	DeleteDataChangeLogByID(ctx context.Context, id uuid.UUID) error
	ListDataChangeLogs(ctx context.Context, auditLogID *uuid.UUID, tableName *string, recordID *uuid.UUID) ([]*ent.DataChangeLog, error)
}

type SchoolEntity interface {
	CreateSchool(ctx context.Context, school *ent.School) (*ent.School, error)
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*ent.School, error)
	UpdateSchoolByID(ctx context.Context, id uuid.UUID, school *ent.School) (*ent.School, error)
	DeleteSchoolByID(ctx context.Context, id uuid.UUID) error
	ListSchools(ctx context.Context) ([]*ent.School, error)
}

type GenderEntity interface {
	CreateGender(ctx context.Context, gender *ent.Gender) (*ent.Gender, error)
	GetGenderByID(ctx context.Context, id uuid.UUID) (*ent.Gender, error)
	UpdateGenderByID(ctx context.Context, id uuid.UUID, gender *ent.Gender) (*ent.Gender, error)
	DeleteGenderByID(ctx context.Context, id uuid.UUID) error
	ListGenders(ctx context.Context, onlyActive bool) ([]*ent.Gender, error)
}

type PrefixEntity interface {
	CreatePrefix(ctx context.Context, prefix *ent.Prefix) (*ent.Prefix, error)
	GetPrefixByID(ctx context.Context, id uuid.UUID) (*ent.Prefix, error)
	UpdatePrefixByID(ctx context.Context, id uuid.UUID, prefix *ent.Prefix) (*ent.Prefix, error)
	DeletePrefixByID(ctx context.Context, id uuid.UUID) error
	ListPrefixes(ctx context.Context, onlyActive bool) ([]*ent.Prefix, error)
}

type AcademicYearEntity interface {
	CreateAcademicYear(ctx context.Context, academicYear *ent.AcademicYear) (*ent.AcademicYear, error)
	GetAcademicYearByID(ctx context.Context, id uuid.UUID) (*ent.AcademicYear, error)
	UpdateAcademicYearByID(ctx context.Context, id uuid.UUID, academicYear *ent.AcademicYear) (*ent.AcademicYear, error)
	DeleteAcademicYearByID(ctx context.Context, id uuid.UUID) error
	ListAcademicYears(ctx context.Context, onlyActive bool, onlyCurrent bool) ([]*ent.AcademicYear, error)
}

type ClassroomEntity interface {
	CreateClassroom(ctx context.Context, classroom *ent.Classroom) (*ent.Classroom, error)
	GetClassroomByID(ctx context.Context, id uuid.UUID) (*ent.Classroom, error)
	UpdateClassroomByID(ctx context.Context, id uuid.UUID, classroom *ent.Classroom) (*ent.Classroom, error)
	DeleteClassroomByID(ctx context.Context, id uuid.UUID) error
	ListClassrooms(ctx context.Context, schoolID *uuid.UUID) ([]*ent.Classroom, error)
}

type SubjectEntity interface {
	CreateSubject(ctx context.Context, subject *ent.Subject) (*ent.Subject, error)
	GetSubjectByID(ctx context.Context, id uuid.UUID) (*ent.Subject, error)
	UpdateSubjectByID(ctx context.Context, id uuid.UUID, subject *ent.Subject) (*ent.Subject, error)
	DeleteSubjectByID(ctx context.Context, id uuid.UUID) error
	ListSubjects(ctx context.Context, schoolID *uuid.UUID) ([]*ent.Subject, error)
}

type SubjectAssignmentEntity interface {
	CreateSubjectAssignment(ctx context.Context, subjectAssignment *ent.SubjectAssignment) (*ent.SubjectAssignment, error)
	GetSubjectAssignmentByID(ctx context.Context, id uuid.UUID) (*ent.SubjectAssignment, error)
	UpdateSubjectAssignmentByID(ctx context.Context, id uuid.UUID, subjectAssignment *ent.SubjectAssignment) (*ent.SubjectAssignment, error)
	DeleteSubjectAssignmentByID(ctx context.Context, id uuid.UUID) error
	ListSubjectAssignments(ctx context.Context, subjectID *uuid.UUID, teacherID *uuid.UUID, classroomID *uuid.UUID, academicYearID *uuid.UUID) ([]*ent.SubjectAssignment, error)
}

type ScheduleEntity interface {
	CreateSchedule(ctx context.Context, schedule *ent.Schedule) (*ent.Schedule, error)
	GetScheduleByID(ctx context.Context, id uuid.UUID) (*ent.Schedule, error)
	UpdateScheduleByID(ctx context.Context, id uuid.UUID, schedule *ent.Schedule) (*ent.Schedule, error)
	DeleteScheduleByID(ctx context.Context, id uuid.UUID) error
	ListSchedules(ctx context.Context, subjectAssignmentID *uuid.UUID, dayOfWeek *ent.ScheduleDayOfWeek) ([]*ent.Schedule, error)
}

type QuestionBankEntity interface {
	CreateQuestionBank(ctx context.Context, question *ent.QuestionBank) (*ent.QuestionBank, error)
	GetQuestionBankByID(ctx context.Context, id uuid.UUID) (*ent.QuestionBank, error)
	UpdateQuestionBankByID(ctx context.Context, id uuid.UUID, question *ent.QuestionBank) (*ent.QuestionBank, error)
	DeleteQuestionBankByID(ctx context.Context, id uuid.UUID) error
	ListQuestionBanks(ctx context.Context, subjectID *uuid.UUID, teacherID *uuid.UUID, questionType *ent.QuestionBankType) ([]*ent.QuestionBank, error)
}

type QuestionChoiceEntity interface {
	CreateQuestionChoice(ctx context.Context, choice *ent.QuestionChoice) (*ent.QuestionChoice, error)
	GetQuestionChoiceByID(ctx context.Context, id uuid.UUID) (*ent.QuestionChoice, error)
	UpdateQuestionChoiceByID(ctx context.Context, id uuid.UUID, choice *ent.QuestionChoice) (*ent.QuestionChoice, error)
	DeleteQuestionChoiceByID(ctx context.Context, id uuid.UUID) error
	ListQuestionChoices(ctx context.Context, questionID *uuid.UUID) ([]*ent.QuestionChoice, error)
}

type AssessmentSetEntity interface {
	CreateAssessmentSet(ctx context.Context, assessmentSet *ent.AssessmentSet) (*ent.AssessmentSet, error)
	GetAssessmentSetByID(ctx context.Context, id uuid.UUID) (*ent.AssessmentSet, error)
	UpdateAssessmentSetByID(ctx context.Context, id uuid.UUID, assessmentSet *ent.AssessmentSet) (*ent.AssessmentSet, error)
	DeleteAssessmentSetByID(ctx context.Context, id uuid.UUID) error
	ListAssessmentSets(ctx context.Context, subjectAssignmentID *uuid.UUID, onlyPublished bool) ([]*ent.AssessmentSet, error)
}

type InventoryItemEntity interface {
	CreateInventoryItem(ctx context.Context, item *ent.InventoryItem) (*ent.InventoryItem, error)
	GetInventoryItemByID(ctx context.Context, id uuid.UUID) (*ent.InventoryItem, error)
	UpdateInventoryItemByID(ctx context.Context, id uuid.UUID, item *ent.InventoryItem) (*ent.InventoryItem, error)
	DeleteInventoryItemByID(ctx context.Context, id uuid.UUID) error
	ListInventoryItems(ctx context.Context, schoolID *uuid.UUID) ([]*ent.InventoryItem, error)
}

type InventoryRequestEntity interface {
	CreateInventoryRequest(ctx context.Context, request *ent.InventoryRequest) (*ent.InventoryRequest, error)
	GetInventoryRequestByID(ctx context.Context, id uuid.UUID) (*ent.InventoryRequest, error)
	UpdateInventoryRequestByID(ctx context.Context, id uuid.UUID, request *ent.InventoryRequest) (*ent.InventoryRequest, error)
	DeleteInventoryRequestByID(ctx context.Context, id uuid.UUID) error
	ListInventoryRequests(ctx context.Context, itemID *uuid.UUID, requesterMemberID *uuid.UUID, status *ent.InventoryRequestStatus) ([]*ent.InventoryRequest, error)
}

type DocumentTrackingEntity interface {
	CreateDocumentTracking(ctx context.Context, document *ent.DocumentTracking) (*ent.DocumentTracking, error)
	GetDocumentTrackingByID(ctx context.Context, id uuid.UUID) (*ent.DocumentTracking, error)
	UpdateDocumentTrackingByID(ctx context.Context, id uuid.UUID, document *ent.DocumentTracking) (*ent.DocumentTracking, error)
	DeleteDocumentTrackingByID(ctx context.Context, id uuid.UUID) error
	ListDocumentTrackings(ctx context.Context, schoolID *uuid.UUID, senderMemberID *uuid.UUID, receiverMemberID *uuid.UUID, status *ent.DocumentTrackingStatus) ([]*ent.DocumentTracking, error)
}

type MemberEntity interface {
	CreateMember(ctx context.Context, member *ent.Member) (*ent.Member, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*ent.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*ent.Member, error)
	UpdateMemberByID(ctx context.Context, id uuid.UUID, member *ent.Member) (*ent.Member, error)
	UpdateMemberLastLoginByID(ctx context.Context, id uuid.UUID) error
	DeleteMemberByID(ctx context.Context, id uuid.UUID) error
	ListMembers(ctx context.Context, schoolID *uuid.UUID, role *ent.MemberRole, onlyActive bool) ([]*ent.Member, error)
}

type MemberRoleEntity interface {
	AddMemberRole(ctx context.Context, memberID uuid.UUID, role ent.MemberRole) error
	ListMemberRolesByMemberID(ctx context.Context, memberID uuid.UUID) ([]ent.MemberRole, error)
	MemberHasAnyRole(ctx context.Context, memberID uuid.UUID, roles []ent.MemberRole) (bool, error)
}

type MemberTeacherEntity interface {
	CreateTeacher(ctx context.Context, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error)
	GetTeacherByID(ctx context.Context, id uuid.UUID) (*ent.MemberTeacher, error)
	UpdateTeacherByID(ctx context.Context, id uuid.UUID, teacher *ent.MemberTeacher) (*ent.MemberTeacher, error)
	DeleteTeacherByID(ctx context.Context, id uuid.UUID) error
	ListTeachers(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberTeacher, error)
	MemberHasTeacherRole(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type MemberStudentEntity interface {
	CreateStudent(ctx context.Context, student *ent.MemberStudent) (*ent.MemberStudent, error)
	GetStudentByID(ctx context.Context, id uuid.UUID) (*ent.MemberStudent, error)
	UpdateStudentByID(ctx context.Context, id uuid.UUID, student *ent.MemberStudent) (*ent.MemberStudent, error)
	DeleteStudentByID(ctx context.Context, id uuid.UUID) error
	ListStudents(ctx context.Context, memberID *uuid.UUID, advisorTeacherID *uuid.UUID, currentClassroomID *uuid.UUID, onlyActive bool) ([]*ent.MemberStudent, error)
}

type MemberStaffEntity interface {
	CreateStaff(ctx context.Context, staff *ent.MemberStaff) (*ent.MemberStaff, error)
	GetStaffByID(ctx context.Context, id uuid.UUID) (*ent.MemberStaff, error)
	UpdateStaffByID(ctx context.Context, id uuid.UUID, staff *ent.MemberStaff) (*ent.MemberStaff, error)
	DeleteStaffByID(ctx context.Context, id uuid.UUID) error
	ListStaffs(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberStaff, error)
	MemberHasStaffRole(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type MemberAdminEntity interface {
	CreateAdmin(ctx context.Context, admin *ent.MemberAdmin) (*ent.MemberAdmin, error)
	GetAdminByID(ctx context.Context, id uuid.UUID) (*ent.MemberAdmin, error)
	UpdateAdminByID(ctx context.Context, id uuid.UUID, admin *ent.MemberAdmin) (*ent.MemberAdmin, error)
	DeleteAdminByID(ctx context.Context, id uuid.UUID) error
	ListAdmins(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberAdmin, error)
	MemberHasAdminRole(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type StaffEducationEntity interface {
	CreateStaffEducation(ctx context.Context, education *ent.StaffEducation) (*ent.StaffEducation, error)
	UpdateStaffEducationByID(ctx context.Context, id uuid.UUID, education *ent.StaffEducation) (*ent.StaffEducation, error)
	DeleteStaffEducationByID(ctx context.Context, id uuid.UUID) error
	ListStaffEducationsByStaffID(ctx context.Context, staffID uuid.UUID) ([]*ent.StaffEducation, error)
	StaffEducationBelongsToStaff(ctx context.Context, id uuid.UUID, staffID uuid.UUID) (bool, error)
}

type StaffWorkExperienceEntity interface {
	CreateStaffWorkExperience(ctx context.Context, work *ent.StaffWorkExperience) (*ent.StaffWorkExperience, error)
	UpdateStaffWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.StaffWorkExperience) (*ent.StaffWorkExperience, error)
	DeleteStaffWorkExperienceByID(ctx context.Context, id uuid.UUID) error
	ListStaffWorkExperiencesByStaffID(ctx context.Context, staffID uuid.UUID) ([]*ent.StaffWorkExperience, error)
	StaffWorkExperienceBelongsToStaff(ctx context.Context, id uuid.UUID, staffID uuid.UUID) (bool, error)
}

type AdminEducationEntity interface {
	CreateAdminEducation(ctx context.Context, education *ent.AdminEducation) (*ent.AdminEducation, error)
	UpdateAdminEducationByID(ctx context.Context, id uuid.UUID, education *ent.AdminEducation) (*ent.AdminEducation, error)
	DeleteAdminEducationByID(ctx context.Context, id uuid.UUID) error
	ListAdminEducationsByAdminID(ctx context.Context, adminID uuid.UUID) ([]*ent.AdminEducation, error)
	AdminEducationBelongsToAdmin(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (bool, error)
}

type AdminWorkExperienceEntity interface {
	CreateAdminWorkExperience(ctx context.Context, work *ent.AdminWorkExperience) (*ent.AdminWorkExperience, error)
	UpdateAdminWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.AdminWorkExperience) (*ent.AdminWorkExperience, error)
	DeleteAdminWorkExperienceByID(ctx context.Context, id uuid.UUID) error
	ListAdminWorkExperiencesByAdminID(ctx context.Context, adminID uuid.UUID) ([]*ent.AdminWorkExperience, error)
	AdminWorkExperienceBelongsToAdmin(ctx context.Context, id uuid.UUID, adminID uuid.UUID) (bool, error)
}

type MemberParentEntity interface {
	CreateParent(ctx context.Context, parent *ent.MemberParent) (*ent.MemberParent, error)
	GetParentByID(ctx context.Context, id uuid.UUID) (*ent.MemberParent, error)
	UpdateParentByID(ctx context.Context, id uuid.UUID, parent *ent.MemberParent) (*ent.MemberParent, error)
	DeleteParentByID(ctx context.Context, id uuid.UUID) error
	ListParents(ctx context.Context, memberID *uuid.UUID, onlyActive bool) ([]*ent.MemberParent, error)
	MemberHasParentRole(ctx context.Context, memberID uuid.UUID) (bool, error)
}

type MemberParentStudentEntity interface {
	CreateParentStudent(ctx context.Context, parentStudent *ent.MemberParentStudent) (*ent.MemberParentStudent, error)
	UpdateParentStudentByID(ctx context.Context, id uuid.UUID, parentStudent *ent.MemberParentStudent) (*ent.MemberParentStudent, error)
	DeleteParentStudentByID(ctx context.Context, id uuid.UUID) error
	ListParentStudentsByParentID(ctx context.Context, parentID uuid.UUID) ([]*ent.MemberParentStudent, error)
	ParentStudentBelongsToParent(ctx context.Context, id uuid.UUID, parentID uuid.UUID) (bool, error)
	ParentExistsByID(ctx context.Context, parentID uuid.UUID) (bool, error)
	StudentExistsByID(ctx context.Context, studentID uuid.UUID) (bool, error)
}

type StudentEnrollmentEntity interface {
	CreateStudentEnrollment(ctx context.Context, enrollment *ent.StudentEnrollment) (*ent.StudentEnrollment, error)
	UpdateStudentEnrollmentByID(ctx context.Context, id uuid.UUID, enrollment *ent.StudentEnrollment) (*ent.StudentEnrollment, error)
	DeleteStudentEnrollmentByID(ctx context.Context, id uuid.UUID) error
	ListStudentEnrollmentsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentEnrollment, error)
	StudentEnrollmentBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
}

type StudentAssessmentSubmissionEntity interface {
	CreateStudentAssessmentSubmission(ctx context.Context, submission *ent.StudentAssessmentSubmission) (*ent.StudentAssessmentSubmission, error)
	UpdateStudentAssessmentSubmissionByID(ctx context.Context, id uuid.UUID, submission *ent.StudentAssessmentSubmission) (*ent.StudentAssessmentSubmission, error)
	DeleteStudentAssessmentSubmissionByID(ctx context.Context, id uuid.UUID) error
	ListStudentAssessmentSubmissionsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentAssessmentSubmission, error)
	StudentAssessmentSubmissionBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
}

type StudentInvoiceEntity interface {
	CreateStudentInvoice(ctx context.Context, invoice *ent.StudentInvoice) (*ent.StudentInvoice, error)
	UpdateStudentInvoiceByID(ctx context.Context, id uuid.UUID, invoice *ent.StudentInvoice) (*ent.StudentInvoice, error)
	DeleteStudentInvoiceByID(ctx context.Context, id uuid.UUID) error
	ListStudentInvoicesByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentInvoice, error)
	StudentInvoiceBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
	FeeCategoryBelongsToStudent(ctx context.Context, feeCategoryID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type StudentAttendanceLogEntity interface {
	CreateStudentAttendanceLog(ctx context.Context, attendanceLog *ent.StudentAttendanceLog) (*ent.StudentAttendanceLog, error)
	UpdateStudentAttendanceLogByID(ctx context.Context, id uuid.UUID, attendanceLog *ent.StudentAttendanceLog) (*ent.StudentAttendanceLog, error)
	DeleteStudentAttendanceLogByID(ctx context.Context, id uuid.UUID) error
	ListStudentAttendanceLogsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.StudentAttendanceLog, error)
	StudentAttendanceLogBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
	EnrollmentBelongsToStudent(ctx context.Context, enrollmentID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type PaymentTransactionEntity interface {
	CreatePaymentTransaction(ctx context.Context, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error)
	UpdatePaymentTransactionByID(ctx context.Context, id uuid.UUID, transaction *ent.PaymentTransaction) (*ent.PaymentTransaction, error)
	DeletePaymentTransactionByID(ctx context.Context, id uuid.UUID) error
	ListPaymentTransactionsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.PaymentTransaction, error)
	PaymentTransactionBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
	InvoiceBelongsToStudent(ctx context.Context, invoiceID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type FeeCategoryEntity interface {
	CreateFeeCategory(ctx context.Context, feeCategory *ent.FeeCategory) (*ent.FeeCategory, error)
	UpdateFeeCategoryByID(ctx context.Context, id uuid.UUID, feeCategory *ent.FeeCategory) (*ent.FeeCategory, error)
	DeleteFeeCategoryByID(ctx context.Context, id uuid.UUID) error
	ListFeeCategoriesByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.FeeCategory, error)
	ResolveSchoolIDByStudentID(ctx context.Context, studentID uuid.UUID) (uuid.UUID, error)
	FeeCategoryBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
}

type GradeItemEntity interface {
	CreateGradeItem(ctx context.Context, gradeItem *ent.GradeItem) (*ent.GradeItem, error)
	UpdateGradeItemByID(ctx context.Context, id uuid.UUID, gradeItem *ent.GradeItem) (*ent.GradeItem, error)
	DeleteGradeItemByID(ctx context.Context, id uuid.UUID) error
	ListGradeItemsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeItem, error)
	GradeItemBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
	SubjectAssignmentBelongsToStudent(ctx context.Context, subjectAssignmentID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type GradeRecordEntity interface {
	CreateGradeRecord(ctx context.Context, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error)
	UpdateGradeRecordByID(ctx context.Context, id uuid.UUID, gradeRecord *ent.GradeRecord) (*ent.GradeRecord, error)
	DeleteGradeRecordByID(ctx context.Context, id uuid.UUID) error
	ListGradeRecordsByStudentID(ctx context.Context, studentID uuid.UUID) ([]*ent.GradeRecord, error)
	GradeRecordBelongsToStudent(ctx context.Context, id uuid.UUID, studentID uuid.UUID) (bool, error)
	EnrollmentBelongsToStudent(ctx context.Context, enrollmentID uuid.UUID, studentID uuid.UUID) (bool, error)
	GradeItemBelongsToStudent(ctx context.Context, gradeItemID uuid.UUID, studentID uuid.UUID) (bool, error)
}

type TeacherEducationEntity interface {
	CreateTeacherEducation(ctx context.Context, education *ent.TeacherEducation) (*ent.TeacherEducation, error)
	UpdateTeacherEducationByID(ctx context.Context, id uuid.UUID, education *ent.TeacherEducation) (*ent.TeacherEducation, error)
	DeleteTeacherEducationByID(ctx context.Context, id uuid.UUID) error
	ListTeacherEducationsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherEducation, error)
	TeacherEducationBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}

type TeacherWorkExperienceEntity interface {
	CreateTeacherWorkExperience(ctx context.Context, work *ent.TeacherWorkExperience) (*ent.TeacherWorkExperience, error)
	UpdateTeacherWorkExperienceByID(ctx context.Context, id uuid.UUID, work *ent.TeacherWorkExperience) (*ent.TeacherWorkExperience, error)
	DeleteTeacherWorkExperienceByID(ctx context.Context, id uuid.UUID) error
	ListTeacherWorkExperiencesByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherWorkExperience, error)
	TeacherWorkExperienceBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}

type TeacherProfileRequestEntity interface {
	CreateTeacherProfileRequest(ctx context.Context, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error)
	UpdateTeacherProfileRequestByID(ctx context.Context, id uuid.UUID, profileRequest *ent.TeacherProfileRequest) (*ent.TeacherProfileRequest, error)
	ListTeacherProfileRequestsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherProfileRequest, error)
	TeacherProfileRequestBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}

type TeacherPerformanceAgreementEntity interface {
	CreateTeacherPerformanceAgreement(ctx context.Context, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error)
	UpdateTeacherPerformanceAgreementByID(ctx context.Context, id uuid.UUID, agreement *ent.TeacherPerformanceAgreement) (*ent.TeacherPerformanceAgreement, error)
	ListTeacherPerformanceAgreementsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPerformanceAgreement, error)
	TeacherPerformanceAgreementBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}

type TeacherPDALogEntity interface {
	CreateTeacherPDALog(ctx context.Context, pdaLog *ent.TeacherPDALog) (*ent.TeacherPDALog, error)
	DeleteTeacherPDALogByID(ctx context.Context, id uuid.UUID) error
	ListTeacherPDALogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherPDALog, error)
	TeacherPDALogBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}

type TeacherLeaveLogEntity interface {
	CreateTeacherLeaveLog(ctx context.Context, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error)
	UpdateTeacherLeaveLogByID(ctx context.Context, id uuid.UUID, leaveLog *ent.TeacherLeaveLog) (*ent.TeacherLeaveLog, error)
	ListTeacherLeaveLogsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]*ent.TeacherLeaveLog, error)
	TeacherLeaveLogBelongsToTeacher(ctx context.Context, id uuid.UUID, teacherID uuid.UUID) (bool, error)
}
