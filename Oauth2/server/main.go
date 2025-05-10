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

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	e.Logger.Info("Serwer uruchomiony na porcie 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
