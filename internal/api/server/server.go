package server

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/aliskhannn/sales-tracker/internal/config"
)

// New creates a new HTTP server with the specified address and router.
//
// Parameters:
//   - addr: the address (host:port) where the server will listen.
//   - router: the Gin engine (or ginext.Engine) that will handle incoming requests.
//
// Returns:
//   - an *http.Server configured with the given address and router.
func New(cfg *config.Config, router *ginext.Engine) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.HTTPPort,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      router,
	}
}
