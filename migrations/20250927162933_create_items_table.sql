-- +goose Up
-- +goose StatementBegin
CREATE TYPE item_kind AS ENUM ('income', 'expense', 'refund', 'transfer');

CREATE TABLE IF NOT EXISTS items
(
    id          UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    kind        item_kind      NOT NULL,                     -- income / expense / ...
    title       TEXT           NOT NULL,                     -- e.g. "Order #1234", "Coffee"
    amount      NUMERIC(18, 2) NOT NULL CHECK (amount >= 0), -- stored as non-negative decimal
    currency    VARCHAR(3)     NOT NULL DEFAULT 'USD',       -- ISO currency code, adjust as needed
    occurred_at TIMESTAMPTZ    NOT NULL,                     -- date/time of the transaction
    category_id UUID           REFERENCES categories (id) ON DELETE SET NULL,
    metadata    JSONB                   DEFAULT '{}'::jsonb, -- free-form metadata (eg. {"source":"stripe"})
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

-- Indexes to support analytics and filtering
CREATE INDEX IF NOT EXISTS idx_items_occurred_at ON items (occurred_at);
CREATE INDEX IF NOT EXISTS idx_items_amount ON items (amount);
CREATE INDEX IF NOT EXISTS idx_items_category ON items (category_id);
CREATE INDEX IF NOT EXISTS idx_items_kind_occurred_at ON items (kind, occurred_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS item_kind;
DROP TABLE IF EXISTS items;
DROP INDEX IF EXISTS idx_items_occurred_at;
DROP INDEX IF EXISTS idx_items_amount;
DROP INDEX IF EXISTS idx_items_category;
DROP INDEX IF EXISTS idx_items_kind_occurred_at;
-- +goose StatementEnd
