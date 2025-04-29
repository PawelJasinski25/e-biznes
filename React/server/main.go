package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"server/routes"
)

func main() {
	e := echo.New()

	// 🔥 Middleware CORS 🔥
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Rejestracja tras
	routes.RegisterProductRoutes(e)
	routes.RegisterCartRoutes(e)
	routes.RegisterPaymentRoutes(e)

	e.Start(":8080")
}
