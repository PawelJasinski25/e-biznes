package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"server/controllers"
	"server/database"
)

func main() {
	database.ConnectDatabase()

	e := echo.New()

	// Middleware CORS – upewnij się, że frontend (np. localhost:5173) jest dozwolony.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Endpointy klasyczne (rejestracja i logowanie) – 3.0, 3.5:
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	// Endpointy OAuth2 – wymaganie 4.0 (logowanie przez Google)
	e.GET("/auth/google/login", controllers.GoogleLogin)
	e.GET("/auth/google/callback", controllers.GoogleCallback)

	e.Logger.Info("Serwer uruchomiony na porcie 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
