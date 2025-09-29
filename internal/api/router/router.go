package router

import (
	"github.com/aliskhannn/sales-tracker/internal/analytics"
	"github.com/aliskhannn/sales-tracker/internal/config"
	"github.com/aliskhannn/sales-tracker/internal/item"
	"github.com/aliskhannn/sales-tracker/internal/middleware"
	"github.com/gin-gonic/gin"
)

// New creates a new Gin engine and sets up routes for the SalesTracker API.
func New() *gin.Engine {
	r := gin.New()

	// Middlewares: Logger, Recovery, optional CORS
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(cfg)) // твой CORS middleware или можно использовать gin-contrib

	// --- Item routes ---
	itemGroup := r.Group("/api/items")
	{
		// Public routes (optional)
		itemGroup.GET("", func(c *gin.Context) {
			// Example: parse query params and call itemHandler.List(...)
		})
		itemGroup.GET("/:id", func(c *gin.Context) {
			// Example: call itemHandler.GetByID(...)
		})

		// Protected routes
		itemGroup.Use(middleware.Auth(cfg.JWT.Secret, cfg.JWT.TTL))
		{
			itemGroup.POST("", func(c *gin.Context) {
				// call itemHandler.Create(...)
			})
			itemGroup.PUT("/:id", func(c *gin.Context) {
				// call itemHandler.Update(...)
			})
			itemGroup.DELETE("/:id", func(c *gin.Context) {
				// call itemHandler.Delete(...)
			})
		}
	}

	// --- Analytics routes ---
	analyticsGroup := r.Group("/api/analytics")
	{
		// Protected route for aggregated data
		analyticsGroup.Use(middleware.Auth(cfg.JWT.Secret, cfg.JWT.TTL))
		{
			analyticsGroup.GET("/sum", func(c *gin.Context) {
				// parse from/to/category/kind query params
				// call analyticsHandler.Sum(...)
			})
			analyticsGroup.GET("/avg", func(c *gin.Context) {
				// call analyticsHandler.Avg(...)
			})
			analyticsGroup.GET("/count", func(c *gin.Context) {
				// call analyticsHandler.Count(...)
			})
			analyticsGroup.GET("/median", func(c *gin.Context) {
				// call analyticsHandler.Median(...)
			})
			analyticsGroup.GET("/percentile/:p", func(c *gin.Context) {
				// call analyticsHandler.Percentile(...), p from URL param (e.g., 0.9 for 90th percentile)
			})
		}
	}

	return r
}

// --- Example CORS middleware ---
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // настройка по необходимости
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
