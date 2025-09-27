-- +goose Up
-- +goose StatementBegin
-- Materialized view that aggregates daily sums/counts per category/kind.
-- Refresh this periodically (e.g., nightly) or on demand.

CREATE MATERIALIZED VIEW IF NOT EXISTS mv_daily_aggregates AS
SELECT date_trunc('day', occurred_at) AS day,
       kind,
       category_id,
       count(*)                       AS cnt,
       sum(amount)                    AS total_amount,
       avg(amount)                    AS avg_amount
FROM items
GROUP BY date_trunc('day', occurred_at), kind, category_id;

-- Index on the materialized view for faster lookups:
CREATE INDEX IF NOT EXISTS idx_mv_daily_aggregates_day ON mv_daily_aggregates (day);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_mv_daily_aggregates_day;
DROP MATERIALIZED VIEW IF EXISTS mv_daily_aggregates;
-- +goose StatementEnd
