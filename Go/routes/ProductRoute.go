package routes

import (
	"Go/controllers"
	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo) {
	e.GET("/products", controllers.GetProducts)
	e.GET("/products/:id", controllers.GetProduct)
	e.POST("/products", controllers.CreateProduct)
	e.PUT("/products/:id", controllers.UpdateProduct)
	e.DELETE("/products/:id", controllers.DeleteProduct)
}
