package db

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SoftDeletedAt is a nullable timestamp for soft-delete logic.
// Replaces gorm.DeletedAt without pulling in the GORM dependency.
type SoftDeletedAt struct {
	time.Time
	Valid bool
}

func (s SoftDeletedAt) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.Time, nil
}

func (s *SoftDeletedAt) Scan(src any) error {
	if src == nil {
		s.Time, s.Valid = time.Time{}, false
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		s.Time, s.Valid = v, true
		return nil
	default:
		return fmt.Errorf("unsupported type for SoftDeletedAt: %T", src)
	}
}

// Model provides auto-managed timestamps.
type Model struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SoftDeleteModel provides timestamps and a soft-delete flag.
type SoftDeleteModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt SoftDeletedAt  `json:"deletedAt"`
}

// AuditableModel extends SoftDeleteModel with audit fields.
type AuditableModel struct {
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt SoftDeletedAt  `json:"deletedAt"`
	CreatedBy uuid.UUID      `json:"createdBy"`
	UpdatedBy uuid.UUID      `json:"updatedBy"`
	DeletedBy *uuid.UUID     `json:"deletedBy"`
}
