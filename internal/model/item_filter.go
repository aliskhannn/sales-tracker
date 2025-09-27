package model

import (
	"time"

	"github.com/google/uuid"
)

// ItemFilter represents a query filter for retrieving items.
//
// Fields can be nil if not used.
// Limit/Offset provide pagination, SortBy defines order clause.
type ItemFilter struct {
	From       *time.Time `json:"from,omitempty"`
	To         *time.Time `json:"to,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Kind       *ItemKind  `json:"kind,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
	SortBy     string     `json:"sort_by,omitempty"` // e.g. "occurred_at desc"
}
