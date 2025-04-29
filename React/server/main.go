package main

import (
	"github.com/labstack/echo/v4"
	"server/routes"
)

func main() {
	e := echo.New()
	routes.RegisterCartRoutes(e)
	routes.RegisterProductRoutes(e)
	routes.RegisterPaymentRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
