package ent

import (
	"time"

	"github.com/google/uuid"
)

// StudentParent is a flattened relation row used by student->parents endpoints.
type StudentParent struct {
	ID             uuid.UUID
	StudentID      uuid.UUID
	ParentID       uuid.UUID
	Relationship   ParentRelationship
	IsMainGuardian bool
	CreatedAt      time.Time

	ParentMemberID  uuid.UUID
	ParentGenderID  *uuid.UUID
	ParentPrefixID  *uuid.UUID
	ParentCode      *string
	ParentFirstName *string
	ParentLastName  *string
	ParentPhone     *string
	ParentIsActive  bool
}
