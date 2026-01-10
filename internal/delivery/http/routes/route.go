package routes

import (
	"github.com/RehanAthallahAzhar/tokohobby-orders/internal/handlers"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, orderHandler *handlers.OrderHandler, authMiddleware echo.MiddlewareFunc) {
	apiV1 := e.Group("/api/v1")

	apiV1.Use(authMiddleware)

	orderGroup := apiV1.Group("/orders")
	{
		orderGroup.POST("/", orderHandler.CreateOrder())
		orderGroup.GET("/", orderHandler.GetUserOrders())
		orderGroup.GET("/:id", orderHandler.GetOrderDetails())
		orderGroup.POST("/:id/cancel", orderHandler.CancelOrder())
		orderGroup.POST("/reset-caches", orderHandler.ResetAllOrderCaches())
	}
}
