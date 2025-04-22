package routes

import (
	"Go/controllers"
	"github.com/labstack/echo/v4"
)

func CategoryRoutes(e *echo.Echo) {
	e.GET("/categories", controllers.GetCategories)
	e.GET("/categories/:id", controllers.GetCategory)
	e.POST("/categories", controllers.CreateCategory)
	e.PUT("/categories/:id", controllers.UpdateCategory)
	e.DELETE("/categories/:id", controllers.DeleteCategory)
}
