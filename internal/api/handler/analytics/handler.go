package analytics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	"github.com/aliskhannn/sales-tracker/internal/api/request"
	"github.com/aliskhannn/sales-tracker/internal/api/response"
)

type service interface {
	// Sum returns the total amount of items matching the filter.
	Sum(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string) (string, error)

	// Avg returns the average amount of items matching the filter.
	Avg(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string) (string, error)

	// Count returns the number of items matching the filter.
	Count(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string) (int64, error)

	// Median returns the median amount of items matching the filter.
	Median(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string) (string, error)

	// Percentile returns the N-th percentile amount of items matching the filter.
	Percentile(ctx context.Context, from, to *time.Time, categoryID *uuid.UUID, kind *string, percentile float64) (string, error)
}

// Handler provides HTTP handlers for analytics.
type Handler struct {
	service service
}

// NewHandler creates a new analytics handler.
func NewHandler(s service) *Handler {
	return &Handler{service: s}
}

// Query represents query parameters for analytics endpoints.
type Query struct {
	From       *time.Time
	To         *time.Time
	CategoryID *uuid.UUID
	Kind       *string
	Percentile float64
}

// Sum handles GET /analytics/sum.
func (h *Handler) Sum(c *ginext.Context) {
	q, err := parseQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	total, err := h.service.Sum(c, q.From, q.To, q.CategoryID, q.Kind)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate sum")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"sum": total})
}

// Avg handles GET /analytics/avg.
func (h *Handler) Avg(c *ginext.Context) {
	q, err := parseQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	avg, err := h.service.Avg(c, q.From, q.To, q.CategoryID, q.Kind)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate average")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"avg": avg})
}

// Count handles GET /analytics/count.
func (h *Handler) Count(c *ginext.Context) {
	q, err := parseQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	cnt, err := h.service.Count(c, q.From, q.To, q.CategoryID, q.Kind)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate count")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]int64{"count": cnt})
}

// Median handles GET /analytics/median.
func (h *Handler) Median(c *ginext.Context) {
	q, err := parseQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	median, err := h.service.Median(c, q.From, q.To, q.CategoryID, q.Kind)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate median")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"median": median})
}

// Percentile handles GET /analytics/percentile.
func (h *Handler) Percentile(c *gin.Context) {
	q, err := parseQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err)
		return
	}

	value, err := h.service.Percentile(c, q.From, q.To, q.CategoryID, q.Kind, q.Percentile)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to calculate percentile")
		response.Fail(c, http.StatusInternalServerError, fmt.Errorf("internal server error"))
		return
	}

	response.OK(c, map[string]string{"percentile": value})
}

// parseQuery parses common analytics query parameters.
func parseQuery(c *ginext.Context) (*Query, error) {
	from, err := request.ParseTimeQuery(c, "from", time.RFC3339)
	if err != nil {
		return nil, err
	}

	to, err := request.ParseTimeQuery(c, "to", time.RFC3339)
	if err != nil {
		return nil, err
	}

	categoryID, err := request.ParseUUIDQuery(c, "category_id")
	if err != nil {
		return nil, err
	}

	kind := request.ParseStringQueryPtr(c, "kind")

	percentile, err := request.ParseFloatQuery(c, "percentile", 0.9)
	if err != nil {
		return nil, err
	}

	return &Query{
		From:       from,
		To:         to,
		CategoryID: categoryID,
		Kind:       kind,
		Percentile: percentile,
	}, nil
}
