package category

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
	ErrCategoryNotFound  = errors.New("category not found")
	ErrNoCategoriesFound = errors.New("no categories found")
)

// Repository provides methods to interact with categories.
type Repository struct {
	db *dbpg.DB
}

// NewRepository creates a new category repository.
func NewRepository(db *dbpg.DB) *Repository {
	return &Repository{db: db}
}

// Create adds a new category to the database.
func (r *Repository) Create(ctx context.Context, c *model.Category) (uuid.UUID, error) {
	query := `
		INSERT INTO categories (name, description, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := r.db.Master.QueryRowContext(ctx, query, c.Name, c.Description, c.ParentID).Scan(&c.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert category: %w", err)
	}

	return c.ID, nil
}

// GetByID retrieves a category by its ID.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM categories
		WHERE id = $1;
	`

	var c model.Category
	err := r.db.Master.QueryRowContext(ctx, query, id).Scan(
		c.ID, c.Name, c.Description, c.ParentID, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}

		return nil, fmt.Errorf("get category: %w", err)
	}

	return &c, nil
}

// List retrieves all categories from the database.
func (r *Repository) List(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, description, parent_id, created_at, updated_at
		FROM categories;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category

		if err = rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.ParentID, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("list categories: %w", err)
		}

		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	if len(categories) == 0 {
		return nil, ErrNoCategoriesFound
	}

	return categories, nil
}

// Update updates a category.
func (r *Repository) Update(ctx context.Context, c *model.Category) error {
	query := `
		UPDATE categories
		SET name = $1,
			description = $2,
			parent_id = $3,
			updated_at = NOW()
		WHERE id = $4;
	`

	res, err := r.db.ExecContext(ctx, query, c.Name, c.Description, c.ParentID, c.ID)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if n == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// Delete removes a category from the database.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM categories
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}

	if n == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
