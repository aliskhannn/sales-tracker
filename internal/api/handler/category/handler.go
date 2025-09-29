package category

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/sales-tracker/internal/api/request"
	"github.com/aliskhannn/sales-tracker/internal/api/response"
	"github.com/aliskhannn/sales-tracker/internal/model"
	"github.com/aliskhannn/sales-tracker/internal/repository/category"
)

type service interface {
	// Create adds a new category.
	// parentID can be nil if the category has no parent.
	Create(ctx context.Context, name, description string, parentID *uuid.UUID) (uuid.UUID, error)

	// GetByID returns a category by its ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error)

	// List returns all categories.
	List(ctx context.Context) ([]model.Category, error)

	// Update modifies an existing category identified by id.
	// parentID can be nil if the category should not have a parent.
	Update(ctx context.Context, id uuid.UUID, name, description string, parentID *uuid.UUID) error

	// Delete removes a category by its ID.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Handler defines the HTTP layer for categories.
type Handler struct {
	service   service
	validator *validator.Validate
}

// NewHandler creates a new category handler.
func NewHandler(s service, v *validator.Validate) *Handler {
	return &Handler{
		service:   s,
		validator: v,
	}
}

// CreateRequest represents the JSON body for creating a category.
type CreateRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

// UpdateRequest represents the JSON body for updating a category.
type UpdateRequest struct {
	Name        string     `json:"name" validate:"required"`
	Description string     `json:"description" validate:"required"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}

// Create handles POST /categories.
func (h *Handler) Create(c *gin.Context) {
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

	id, err := h.service.Create(c.Request.Context(), req.Name, req.Description, req.ParentID)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create category")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.Created(c, map[string]string{"id": id.String()})
}

// GetByID handles GET /categories/:id.
func (h *Handler) GetByID(c *gin.Context) {
	id, err := request.ParseUUIDParam(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	cat, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		// If category not found, return 404 Not Found.
		if errors.Is(err, category.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		// Internal Server Error.
		zlog.Logger.Error().Err(err).Msg("failed to get category")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]*model.Category{"category": cat})
}

// List handles GET /categories.
func (h *Handler) List(c *gin.Context) {
	categories, err := h.service.List(c.Request.Context())
	if err != nil {
		// If category not found, return 404 Not Found.
		if errors.Is(err, category.ErrNoCategoriesFound) {
			zlog.Logger.Error().Err(err).Msg("categories not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		// Internal Server Error.
		zlog.Logger.Error().Err(err).Msg("failed to get categories")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string][]model.Category{"categories": categories})
}

// Update handles PUT /categories/:id.
func (h *Handler) Update(c *gin.Context) {
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

	if err := h.service.Update(c.Request.Context(), id, req.Name, req.Description, req.ParentID); err != nil {
		// If category not found, return 404 Not Found.
		if errors.Is(err, category.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		// Internal Server Error.
		zlog.Logger.Error().Err(err).Msg("failed to update category")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"message": "category updated"})
}

// Delete handles DELETE /categories/:id.
func (h *Handler) Delete(c *gin.Context) {
	id, err := request.ParseUUIDParam(c, "id")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		// If category not found, return 404 Not Found.
		if errors.Is(err, category.ErrCategoryNotFound) {
			zlog.Logger.Error().Err(err).Msg("category not found")
			response.Fail(c, http.StatusNotFound, err)
			return
		}

		// Internal Server Error.
		zlog.Logger.Error().Err(err).Msg("failed to delete category")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"message": "category deleted"})
}
