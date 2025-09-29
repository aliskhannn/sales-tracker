package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
	ID         uuid.UUID       `db:"id" json:"id"`
	Kind       string          `db:"kind" json:"kind"`
	Title      string          `db:"title" json:"title"`
	Amount     decimal.Decimal `db:"amount" json:"amount"` // as string to preserve precision; parse with decimal libs if needed
	Currency   string          `db:"currency" json:"currency"`
	OccurredAt time.Time       `db:"occurred_at" json:"occurred_at"`
	CategoryID *uuid.UUID      `db:"category_id,omitempty" json:"category_id,omitempty"`
	Metadata   json.RawMessage `db:"metadata" json:"metadata"` // store raw JSONB bytes
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
}
