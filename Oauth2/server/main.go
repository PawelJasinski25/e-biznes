package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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
	
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))

	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)

	e.GET("/auth/google/login", controllers.GoogleLogin)
	e.GET("/auth/google/callback", controllers.GoogleCallback)

	e.GET("/auth/github/login", controllers.GithubLogin)
	e.GET("/auth/github/callback", controllers.GithubCallback)

	e.Logger.Info("Serwer uruchomiony na porcie 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
