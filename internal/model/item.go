package model

import (
	"time"

	"github.com/google/uuid"
)

// ItemKind enumerates the supported kinds of an item/transaction.
//
// Common values:
//
//	"income"  - money coming in
//	"expense" - money going out
//	"refund"  - refunded amount
//	"transfer" - internal transfer (not counted as revenue)
//
// Note: The DB uses a Postgres ENUM "item_kind" to enforce allowed values.
type ItemKind string

// Item represents a financial record / sale / transaction.
//
// Fields:
//   - ID: UUID primary key (DB default gen_random_uuid())
//   - Kind: type of transaction (income/expense/refund/transfer)
//   - Title: short human-friendly description
//   - Amount: non-negative monetary amount (NUMERIC in DB)
//   - Currency: 3-letter ISO currency code, e.g. "USD"
//   - OccurredAt: the timestamp when the transaction actually happened
//   - CategoryID: optional FK to categories table
//   - Metadata: JSONB for extensible attributes
//   - CreatedAt, UpdatedAt: DB-managed timestamps
type Item struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	Kind       ItemKind   `db:"kind" json:"kind"`
	Title      string     `db:"title" json:"title"`
	Amount     string     `db:"amount" json:"amount"` // as string to preserve precision; parse with decimal libs if needed
	Currency   string     `db:"currency" json:"currency"`
	OccurredAt time.Time  `db:"occurred_at" json:"occurred_at"`
	CategoryID *uuid.UUID `db:"category_id,omitempty" json:"category_id,omitempty"`
	Metadata   []byte     `db:"metadata" json:"metadata"` // store raw JSONB bytes; unmarshal when needed
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}
