package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type MemberRole string

const (
	MemberRoleStudent MemberRole = "student"
	MemberRoleTeacher MemberRole = "teacher"
	MemberRoleAdmin   MemberRole = "admin"
	MemberRoleStaff   MemberRole = "staff"
	MemberRoleParent  MemberRole = "parent"
)

type Member struct {
	bun.BaseModel `bun:"table:members,alias:mem"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID  uuid.UUID  `bun:"school_id,type:uuid,notnull"`
	Email     string     `bun:"email,notnull"`
	Password  string     `bun:"password,notnull"`
	Role      MemberRole `bun:"role,notnull"`
	IsActive  bool       `bun:"is_active,notnull"`
	LastLogin *time.Time `bun:"last_login"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
}

func ToMemberRole(value string) MemberRole {
	switch value {
	case "student":
		return MemberRoleStudent
	case "teacher":
		return MemberRoleTeacher
	case "staff":
		return MemberRoleStaff
	case "parent":
		return MemberRoleParent
	default:
		return MemberRoleAdmin
	}
}
