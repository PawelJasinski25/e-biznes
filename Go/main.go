package main

import (
	"Go/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	routes.ProductRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
