package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"

	"github.com/aliskhannn/sales-tracker/internal/api/handler/analytics"
	"github.com/aliskhannn/sales-tracker/internal/api/handler/category"
	"github.com/aliskhannn/sales-tracker/internal/api/handler/item"
)

// New creates a new Gin engine and sets up routes for the SalesTracker API.
func New(
	categoryHandler *category.Handler,
	itemHandler item.Handler,
	analyticsHandler analytics.Handler,
) *ginext.Engine {
	r := ginext.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.POST("", categoryHandler.Create)
			categories.GET("", categoryHandler.List)
			categories.GET("/:id", categoryHandler.GetByID)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		items := api.Group("/items")
		{
			items.POST("", itemHandler.Create)
			items.GET("", itemHandler.List)
			items.GET("/:id", itemHandler.GetByID)
			items.PUT("/:id", itemHandler.Update)
			items.DELETE("/:id", itemHandler.Delete)
		}

		analyticsGroup := api.Group("/analytics")
		{
			analyticsGroup.GET("/sum", analyticsHandler.Sum)
			analyticsGroup.GET("/avg", analyticsHandler.Avg)
			analyticsGroup.GET("/count", analyticsHandler.Count)
			analyticsGroup.GET("/median", analyticsHandler.Median)
			analyticsGroup.GET("/percentile", analyticsHandler.Percentile)
		}
	}

	return r
}
