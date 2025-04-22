package main

import (
	"Go/database"
	"Go/models"
	"Go/routes"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()

	if err := database.InitDatabase(); err != nil {
		fmt.Println("Error connecting with database:", err)
		os.Exit(1)
	}

	err := database.DB.AutoMigrate(
		&models.Product{},
		&models.Category{},
		&models.Cart{},
		&models.CartItem{},
	)
	if err != nil {
		return
	}

	routes.ProductRoutes(e)
	routes.CartRoutes(e)
	routes.CategoryRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
