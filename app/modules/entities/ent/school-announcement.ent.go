package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SchoolAnnouncement struct {
	bun.BaseModel `bun:"table:school_announcements,alias:sca"`

	ID             uuid.UUID   `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID       uuid.UUID   `bun:"school_id,type:uuid,notnull"`
	AuthorMemberID uuid.UUID   `bun:"author_member_id,type:uuid,notnull"`
	Title          *string     `bun:"title"`
	Content        *string     `bun:"content"`
	Category       *string     `bun:"category"`
	Status         string      `bun:"status,notnull"`
	AnnouncedAt    *time.Time  `bun:"announced_at,type:date"`
	PublishedAt    *time.Time  `bun:"published_at,type:date"`
	ExpiresAt      *time.Time  `bun:"expires_at,type:date"`
	CreatedByName  *string     `bun:"created_by_name"`
	TargetRole     *MemberRole `bun:"target_role"`
	IsPinned       bool        `bun:"is_pinned,notnull"`
	CreatedAt      time.Time   `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt      *time.Time  `bun:"deleted_at,soft_delete"`
}
