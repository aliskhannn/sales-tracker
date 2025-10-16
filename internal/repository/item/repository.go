package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/dbpg"

	"github.com/aliskhannn/sales-tracker/internal/model"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrNoItemsFound = errors.New("no items found")
)

// Repository provides methods to interact with items.
type Repository struct {
	db *dbpg.DB
}

// NewRepository creates a new item repository.
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// Create adds a new item to the database.
func (r *Repository) Create(ctx context.Context, i *model.Item) (uuid.UUID, error) {
	query := `
		INSERT INTO items (
		    kind, title, amount, currency, occurred_at, category_id, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
	`

	err := r.db.QueryRowContext(ctx, query,
		i.Kind, i.Title, i.Amount, i.Currency, i.OccurredAt, i.CategoryID, i.Metadata,
	).Scan(&i.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert item: %w", err)
	}

	return i.ID, nil
}

// GetByID retrieves an item by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	query := `
		SELECT id, kind, title, amount, currency, occurred_at, category_id, metadata, created_at, updated_at
		FROM items
		WHERE id = $1;
	`

	var i model.Item
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID, &i.Kind, &i.Title, &i.Amount, &i.Currency, &i.OccurredAt,
		&i.CategoryID, &i.Metadata, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}

		return nil, fmt.Errorf("get item: %w", err)
	}

	return &i, nil
}

// List retrieves items from the database applying optional filters.
// Filters can include date range (From, To), category, kind, pagination (Limit, Offset),
// and sort order (SortBy).
func (r *Repository) List(ctx context.Context, filter *model.ItemFilter) ([]model.Item, error) {
	query := `
		SELECT id, kind, title, amount, currency, occurred_at, category_id, metadata, created_at, updated_at
		FROM items
		WHERE ($1::timestamptz IS NULL OR occurred_at >= $1)
		  AND ($2::timestamptz IS NULL OR occurred_at <= $2)
		  AND ($3::uuid IS NULL OR category_id = $3)
		  AND ($4::item_kind IS NULL OR kind = $4)
		ORDER BY occurred_at DESC
		LIMIT $5 OFFSET $6;
	`

	rows, err := r.db.QueryContext(ctx, query,
		filter.From,
		filter.To,
		filter.CategoryID,
		filter.Kind,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var i model.Item
		if err = rows.Scan(
			&i.ID, &i.Kind, &i.Title, &i.Amount, &i.Currency, &i.OccurredAt,
			&i.CategoryID, &i.Metadata, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list items: %w", err)
		}

		items = append(items, i)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	return items, nil
}

// Update updates an item.
func (r *Repository) Update(ctx context.Context, i *model.Item) error {
	query := `
		UPDATE items
		SET
    		kind = $1,
    		title = $2,
    		amount = $3,
    		currency = $4,
    		occurred_at = $5,
    		category_id = $6,
    		metadata = $7,
    		updated_at = NOW()
		WHERE id = $8;
	`

	res, err := r.db.ExecContext(ctx, query,
		i.Kind,
		i.Title,
		i.Amount,
		i.Currency,
		i.OccurredAt,
		i.CategoryID,
		i.Metadata,
		i.ID,
	)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if n == 0 {
		return ErrItemNotFound
	}

	return nil
}

// Delete removes an item from the database.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM items
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if n == 0 {
		return ErrItemNotFound
	}

	return nil
}
