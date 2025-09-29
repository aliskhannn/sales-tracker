package request

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// ParseUUIDParam parses a UUID from the URL parameters and logs errors if invalid.
// Returns the UUID and an error if parsing fails.
func ParseUUIDParam(c *ginext.Context, key string) (uuid.UUID, error) {
	value := c.Param(key)
	id, err := uuid.Parse(value)
	if err != nil {
		zlog.Logger.Error().Err(err).Interface(key, value).Msg("failed to parse UUID param")
		return uuid.Nil, fmt.Errorf("invalid %s", key)
	}

	return id, nil
}

// ParseUUIDQuery parses a query parameter as *uuid.UUID.
// Returns nil if parameter is empty.
func ParseUUIDQuery(c *ginext.Context, key string) (*uuid.UUID, error) {
	value := c.Query(key)
	if value == "" {
		return nil, nil
	}

	id, err := uuid.Parse(value)
	if err != nil {
		zlog.Logger.Error().Err(err).Str(key, value).Msg("failed to parse UUID query")
		return nil, fmt.Errorf("invalid %s", key)
	}

	return &id, nil
}

// ParseTimeQuery parses a query parameter as *time.Time.
// Returns nil if parameter is empty.
func ParseTimeQuery(c *ginext.Context, key string, layout string) (*time.Time, error) {
	value := c.Query(key)
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(layout, value)
	if err != nil {
		zlog.Logger.Error().Err(err).Interface(key, value).Msg("failed to parse time query")
		return nil, fmt.Errorf("invalid time format for %s", key)
	}

	return &t, nil
}

// ParseStringQuery returns a string query parameter or defaultValue if empty.
func ParseStringQuery(c *ginext.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// ParseStringQueryPtr parses a query parameter as *string.
// Returns nil if parameter is empty.
func ParseStringQueryPtr(c *ginext.Context, key string) *string {
	value := c.Query(key)
	if value == "" {
		return nil
	}

	return &value
}

// ParseIntQuery parses a query parameter as int, returns defaultValue if empty.
// Returns error if value is present but not a valid integer.
func ParseIntQuery(c *ginext.Context, key string, defaultValue int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return defaultValue, nil
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		zlog.Logger.Error().Err(err).Str(key, value).Msg("failed to parse int query")
		return 0, fmt.Errorf("invalid int format for %s", key)
	}

	return n, nil
}
