package category

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/aliskhannn/sales-tracker/internal/model"
)

// repository provides methods to interact with categories.
type repository interface {
	// Create adds a new category to the database.
	Create(ctx context.Context, c *model.Category) (uuid.UUID, error)

	// GetByID retrieves a category by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error)

	// List retrieves all categories from the database.
	List(ctx context.Context) ([]model.Category, error)

	// Update updates a category.
	Update(ctx context.Context, c *model.Category) error

	// Delete removes a category from the database.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Service provides category-related business logic.
type Service struct {
	repository repository
}

// NewService creates a new category service.
func NewService(r repository) *Service {
	return &Service{repository: r}
}

// Create adds a new category.
// parentID can be nil if the category has no parent.
func (s *Service) Create(ctx context.Context, name, description string, parentID *uuid.UUID) (uuid.UUID, error) {
	c := &model.Category{
		Name:        name,
		Description: &description,
		ParentID:    parentID,
	}

	id, err := s.repository.Create(ctx, c)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create category: %w", err)
	}

	return id, nil
}

// GetByID returns a category by its ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	c, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get category: %w", err)
	}

	return c, nil
}

// List returns all categories.
func (s *Service) List(ctx context.Context) ([]model.Category, error) {
	categories, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	return categories, nil
}

// Update modifies an existing category identified by id.
// parentID can be nil if the category should not have a parent.
func (s *Service) Update(ctx context.Context, id uuid.UUID, name, description string, parentID uuid.UUID) error {
	c := &model.Category{
		ID:          id,
		Name:        name,
		Description: &description,
		ParentID:    &parentID,
	}

	err := s.repository.Update(ctx, c)
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	return nil
}

// Delete removes a category by its ID.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}

	return nil
}
