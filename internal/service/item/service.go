package item

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/aliskhannn/sales-tracker/internal/model"
)

type repository interface {
	// Create adds a new item to the database.
	Create(ctx context.Context, i *model.Item) (uuid.UUID, error)

	// GetByID retrieves an item by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error)

	// List retrieves items from the database applying optional filters.
	// Filters can include date range (From, To), category, kind, pagination (Limit, Offset),
	// and sort order (SortBy).
	List(ctx context.Context, filter *model.ItemFilter) ([]model.Item, error)

	// Update updates an item.
	Update(ctx context.Context, i *model.Item) error

	// Delete removes an item from the database.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Service provides item-related business logic.
type Service struct {
	repository repository
}

// NewService creates a new item service.
func NewService(r repository) *Service {
	return &Service{repository: r}
}

// Create adds a new item with the given fields.
func (s *Service) Create(
	ctx context.Context,
	kind string,
	title string,
	amount decimal.Decimal,
	currency string,
	occurredAt time.Time,
	metadata []byte,
) (uuid.UUID, error) {
	i := &model.Item{
		Kind:       kind,
		Title:      title,
		Amount:     amount,
		Currency:   currency,
		OccurredAt: occurredAt,
		Metadata:   metadata,
	}

	id, err := s.repository.Create(ctx, i)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create item: %w", err)
	}

	return id, nil
}

// GetByID returns an item by its ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error) {
	i, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	return i, nil
}

// List returns items applying the given filters such as date range,
// category, kind, pagination, and sort order.
func (s *Service) List(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
	limit, offset int,
	sortBy string,
) ([]model.Item, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
		Limit:      limit,
		Offset:     offset,
		SortBy:     sortBy,
	}

	items, err := s.repository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	return items, nil
}

// Update modifies an existing item by its ID.
func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	kind string,
	title string,
	amount decimal.Decimal,
	currency string,
	occurredAt time.Time,
	metadata []byte,
) error {
	i := &model.Item{
		ID:         id,
		Kind:       kind,
		Title:      title,
		Amount:     amount,
		Currency:   currency,
		OccurredAt: occurredAt,
		Metadata:   metadata,
	}

	err := s.repository.Update(ctx, i)
	if err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

// Delete removes an item by its ID.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	return nil
}
