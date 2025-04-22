package routes

import (
	"Go/controllers"
	"github.com/labstack/echo/v4"
)

func CartRoutes(e *echo.Echo) {
	e.GET("/carts", controllers.GetCarts)
	e.GET("/carts/:id", controllers.GetCart)
	e.POST("/carts", controllers.CreateCart)
	e.PUT("/carts/:id", controllers.UpdateCart)
	e.DELETE("/carts/:id", controllers.DeleteCart)
}
