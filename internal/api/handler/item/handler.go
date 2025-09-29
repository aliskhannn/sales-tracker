package item

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/sales-tracker/internal/api/request"
	"github.com/aliskhannn/sales-tracker/internal/api/response"
	"github.com/aliskhannn/sales-tracker/internal/model"
	"github.com/aliskhannn/sales-tracker/internal/repository/item"
)

// service defines business logic for items.
type service interface {
	// Create adds a new item with the given fields.
	Create(ctx context.Context, kind, title string, amount decimal.Decimal, currency string, occurredAt time.Time, categoryID *uuid.UUID, metadata json.RawMessage) (uuid.UUID, error)

	// GetByID returns an item by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Item, error)

	// List returns items applying the given filters such as date range,
	// category, kind, pagination, and sort order.
	List(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string, limit, offset int, sortBy string) ([]model.Item, error)

	// Update modifies an existing item by its ID.
	Update(ctx context.Context, id uuid.UUID, kind, title string, amount decimal.Decimal, currency string, occurredAt time.Time, categoryID *uuid.UUID, metadata json.RawMessage) error

	// Delete removes an item by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Handler defines HTTP layer for items.
type Handler struct {
	service   service
	validator *validator.Validate
}

// NewHandler creates a new item handler.
func NewHandler(s service, v *validator.Validate) *Handler {
	return &Handler{service: s, validator: v}
}

// CreateRequest JSON body for creating an item.
type CreateRequest struct {
	Kind       string          `json:"kind" validate:"required"`
	Title      string          `json:"title" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required"`
	Currency   string          `json:"currency" validate:"required"`
	OccurredAt time.Time       `json:"occurred_at" validate:"required"`
	CategoryID *uuid.UUID      `json:"category_id,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

// UpdateRequest JSON body for updating an item.
type UpdateRequest struct {
	Kind       string          `json:"kind" validate:"required"`
	Title      string          `json:"title" validate:"required"`
	Amount     decimal.Decimal `json:"amount" validate:"required"`
	Currency   string          `json:"currency" validate:"required"`
	OccurredAt time.Time       `json:"occurred_at" validate:"required"`
	CategoryID *uuid.UUID      `json:"category_id,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
}

// Create handles POST /items.
func (h *Handler) Create(c *ginext.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind create request")
		response.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to validate request")
		response.Fail(c, http.StatusBadRequest, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	if len(req.Metadata) == 0 {
		req.Metadata = json.RawMessage(`{}`)
	}

	id, err := h.service.Create(c.Request.Context(), req.Kind, req.Title, req.Amount, req.Currency, req.OccurredAt, req.CategoryID, req.Metadata)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create item")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.Created(c, map[string]string{"id": id.String()})
}

// GetByID handles GET /items/:id.
func (h *Handler) GetByID(c *ginext.Context) {
	id, err := request.ParseUUIDParam(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	i, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, item.ErrItemNotFound) {
			zlog.Logger.Error().Err(err).Msg("item not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to get item")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]*model.Item{"item": i})
}

// List handles GET /items.
func (h *Handler) List(c *ginext.Context) {
	from, err := request.ParseTimeQuery(c, "from", time.RFC3339)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	to, err := request.ParseTimeQuery(c, "to", time.RFC3339)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	categoryID, err := request.ParseUUIDQuery(c, "category_id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	kind := request.ParseStringQueryPtr(c, "kind")

	limit, err := request.ParseIntQuery(c, "limit", 20) // default = 20
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	offset, err := request.ParseIntQuery(c, "offset", 0)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	sortBy := request.ParseStringQuery(c, "sort_by", "occurred_at")

	items, err := h.service.List(c.Request.Context(), from, to, categoryID, kind, limit, offset, sortBy)
	if err != nil {
		if errors.Is(err, item.ErrNoItemsFound) {
			zlog.Logger.Error().Err(err).Msg("items not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to list items")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string][]model.Item{"items": items})
}

// Update handles PUT /items/:id.
func (h *Handler) Update(c *ginext.Context) {
	id, err := request.ParseUUIDParam(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to bind update request")
		response.Fail(c, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if err := h.validator.Struct(req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to validate request")
		response.Fail(c, http.StatusBadRequest, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	if len(req.Metadata) == 0 {
		req.Metadata = json.RawMessage(`{}`)
	}

	if err := h.service.Update(c.Request.Context(), id, req.Kind, req.Title, req.Amount, req.Currency, req.OccurredAt, req.CategoryID, req.Metadata); err != nil {
		if errors.Is(err, item.ErrItemNotFound) {
			zlog.Logger.Error().Err(err).Msg("item not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to update item")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"message": "item updated"})
}

// Delete handles DELETE /items/:id.
func (h *Handler) Delete(c *ginext.Context) {
	id, err := request.ParseUUIDParam(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, item.ErrItemNotFound) {
			zlog.Logger.Error().Err(err).Msg("item not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		zlog.Logger.Error().Err(err).Msg("failed to delete item")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"message": "item deleted"})
}
