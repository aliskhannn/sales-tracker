package analytics

import (
	"context"
	"fmt"

	"github.com/wb-go/wbf/dbpg"

	"github.com/aliskhannn/sales-tracker/internal/model"
)

// Repository provides methods to interact with analytics.
type Repository struct {
	db *dbpg.DB
}

// NewRepository creates a new analytics repository.
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// Sum calculates the total amount of items matching the filter.
func (r *Repository) Sum(ctx context.Context, filter *model.ItemFilter) (string, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM items
		WHERE ($1::timestamptz IS NULL OR occurred_at >= $1)
		  AND ($2::timestamptz IS NULL OR occurred_at <= $2)
		  AND ($3::uuid IS NULL OR category_id = $3)
		  AND ($4::text IS NULL OR kind = $4);
	`

	var total string
	err := r.db.Master.QueryRowContext(ctx, query,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
	).Scan(&total)
	if err != nil {
		return "", fmt.Errorf("sum items: %w", err)
	}

	return total, nil
}

// Avg calculates the average amount of items matching the filter.
func (r *Repository) Avg(ctx context.Context, filter *model.ItemFilter) (string, error) {
	query := `
		SELECT COALESCE(AVG(amount), 0)
		FROM items
		WHERE ($1::timestamptz IS NULL OR occurred_at >= $1)
		  AND ($2::timestamptz IS NULL OR occurred_at <= $2)
		  AND ($3::uuid IS NULL OR category_id = $3)
		  AND ($4::text IS NULL OR kind = $4);
	`

	var avg string
	err := r.db.Master.QueryRowContext(ctx, query,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
	).Scan(&avg)
	if err != nil {
		return "", fmt.Errorf("avg items: %w", err)
	}

	return avg, nil
}

// Count returns the number of items matching the filter.
func (r *Repository) Count(ctx context.Context, filter *model.ItemFilter) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM items
		WHERE ($1::timestamptz IS NULL OR occurred_at >= $1)
		  AND ($2::timestamptz IS NULL OR occurred_at <= $2)
		  AND ($3::uuid IS NULL OR category_id = $3)
		  AND ($4::text IS NULL OR kind = $4);
	`

	var cnt int64
	err := r.db.Master.QueryRowContext(ctx, query,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
	).Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("count items: %w", err)
	}

	return cnt, nil
}

// Median calculates the median amount of items matching the filter.
func (r *Repository) Median(ctx context.Context, filter *model.ItemFilter) (string, error) {
	query := `
		SELECT COALESCE(
			percentile_cont(0.5) WITHIN GROUP (ORDER BY amount),
			0
		)
		FROM items
		WHERE ($1::timestamptz IS NULL OR occurred_at >= $1)
		  AND ($2::timestamptz IS NULL OR occurred_at <= $2)
		  AND ($3::uuid IS NULL OR category_id = $3)
		  AND ($4::text IS NULL OR kind = $4);
	`

	var median string
	err := r.db.Master.QueryRowContext(ctx, query,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
	).Scan(&median)
	if err != nil {
		return "", fmt.Errorf("median items: %w", err)
	}

	return median, nil
}

// Percentile calculates the N-th percentile (0.0–1.0) of items matching the filter.
func (r *Repository) Percentile(ctx context.Context, filter *model.ItemFilter, percentile float64) (string, error) {
	query := `
		SELECT COALESCE(
			percentile_cont($1) WITHIN GROUP (ORDER BY amount),
			0
		)
		FROM items
		WHERE ($2::timestamptz IS NULL OR occurred_at >= $2)
		  AND ($3::timestamptz IS NULL OR occurred_at <= $3)
		  AND ($4::uuid IS NULL OR category_id = $4)
		  AND ($5::text IS NULL OR kind = $5);
	`

	var value string
	err := r.db.Master.QueryRowContext(ctx, query,
		percentile,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
	).Scan(&value)
	if err != nil {
		return "", fmt.Errorf("percentile items: %w", err)
	}

	return value, nil
}
