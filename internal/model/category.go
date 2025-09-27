package model

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a hierarchical category for items.
//
// Fields:
//   - ID: UUID primary key, generated in DB via gen_random_uuid()
//   - Name: human-readable category name
//   - Description: optional description text
//   - ParentID: optional parent category UUID for hierarchy
//   - CreatedAt, UpdatedAt: timestamps managed by DB
type Category struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description,omitempty" json:"description,omitempty"`
	ParentID    *uuid.UUID `db:"parent_id,omitempty" json:"parent_id,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
