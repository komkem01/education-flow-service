package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParsePathUUID(ctx *gin.Context, param string) (uuid.UUID, error) {
	return uuid.Parse(ctx.Param(param))
}

func ParseUUIDPtr(raw *string) (*uuid.UUID, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(*raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func ParseQueryUUID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func UUIDToStringPtr(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}
	parsed := value.String()
	return &parsed
}

func ParseTimePtrWithLayout(raw *string, layout string) (*time.Time, error) {
	if raw == nil || *raw == "" {
		return nil, nil
	}
	parsed, err := time.Parse(layout, *raw)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func TimeToStringPtrWithLayout(value *time.Time, layout string) *string {
	if value == nil {
		return nil
	}
	parsed := value.UTC().Format(layout)
	return &parsed
}
