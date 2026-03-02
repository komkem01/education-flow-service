package ent

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type DocumentPriority string
type DocumentTrackingStatus string

const (
	DocumentPriorityNormal    DocumentPriority = "normal"
	DocumentPriorityUrgent    DocumentPriority = "urgent"
	DocumentPriorityTopUrgent DocumentPriority = "top_urgent"
)

const (
	DocumentTrackingStatusSent      DocumentTrackingStatus = "sent"
	DocumentTrackingStatusRead      DocumentTrackingStatus = "read"
	DocumentTrackingStatusProcessed DocumentTrackingStatus = "processed"
)

type DocumentTracking struct {
	bun.BaseModel `bun:"table:document_tracking,alias:dtk"`

	ID               uuid.UUID              `bun:"id,pk,type:uuid,default:gen_random_uuid()"`
	SchoolID         uuid.UUID              `bun:"school_id,type:uuid,notnull"`
	DocNumber        *string                `bun:"doc_number"`
	Title            *string                `bun:"title"`
	ContentSummary   *string                `bun:"content_summary"`
	Priority         DocumentPriority       `bun:"priority"`
	SenderMemberID   *uuid.UUID             `bun:"sender_member_id,type:uuid"`
	ReceiverMemberID *uuid.UUID             `bun:"receiver_member_id,type:uuid"`
	FileURL          *string                `bun:"file_url"`
	Status           DocumentTrackingStatus `bun:"status"`
	CreatedAt        time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func ToDocumentPriority(value string) DocumentPriority {
	switch value {
	case "urgent":
		return DocumentPriorityUrgent
	case "top_urgent":
		return DocumentPriorityTopUrgent
	default:
		return DocumentPriorityNormal
	}
}

func ToDocumentTrackingStatus(value string) DocumentTrackingStatus {
	switch value {
	case "read":
		return DocumentTrackingStatusRead
	case "processed":
		return DocumentTrackingStatusProcessed
	default:
		return DocumentTrackingStatusSent
	}
}
