package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aliskhannn/sales-tracker/internal/model"
)

// repository defines the required behavior for analytics persistence.
type repository interface {
	// Sum calculates the total amount of items matching the filter.
	Sum(ctx context.Context, filter *model.ItemFilter) (string, error)

	// Avg calculates the average amount of items matching the filter.
	Avg(ctx context.Context, filter *model.ItemFilter) (string, error)

	// Count returns the number of items matching the filter.
	Count(ctx context.Context, filter *model.ItemFilter) (int64, error)

	// Median calculates the median amount of items matching the filter.
	Median(ctx context.Context, filter *model.ItemFilter) (string, error)

	// Percentile calculates the N-th percentile of items matching the filter.
	Percentile(ctx context.Context, filter *model.ItemFilter, percentile float64) (string, error)
}

// Service provides analytics-related business logic.
type Service struct {
	repository repository
}

// NewService creates a new analytics service.
func NewService(r repository) *Service {
	return &Service{repository: r}
}

// Sum returns the total amount of items matching the filter.
func (s *Service) Sum(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
) (string, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
	}

	total, err := s.repository.Sum(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("analytics sum: %w", err)
	}
	return total, nil
}

// Avg returns the average amount of items matching the filter.
func (s *Service) Avg(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
) (string, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
	}

	avg, err := s.repository.Avg(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("analytics avg: %w", err)
	}
	return avg, nil
}

// Count returns the number of items matching the filter.
func (s *Service) Count(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
) (int64, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
	}

	cnt, err := s.repository.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("analytics count: %w", err)
	}
	return cnt, nil
}

// Median returns the median amount of items matching the filter.
func (s *Service) Median(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
) (string, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
	}

	median, err := s.repository.Median(ctx, filter)
	if err != nil {
		return "", fmt.Errorf("analytics median: %w", err)
	}
	return median, nil
}

// Percentile returns the N-th percentile amount of items matching the filter.
func (s *Service) Percentile(
	ctx context.Context,
	from, to *time.Time,
	categoryID *uuid.UUID,
	kind *string,
	percentile float64,
) (string, error) {
	filter := &model.ItemFilter{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
	}

	value, err := s.repository.Percentile(ctx, filter, percentile)
	if err != nil {
		return "", fmt.Errorf("analytics percentile: %w", err)
	}
	return value, nil
}
