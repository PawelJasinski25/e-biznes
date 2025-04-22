package controllers

import (
	"net/http"
	"strconv"

	"Go/database"
	"Go/models"
	"github.com/labstack/echo/v4"
)

func GetProducts(c echo.Context) error {
	var products []models.Product
	result := database.DB.Preload("Category").Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product models.Product
	result := database.DB.Preload("Category").First(&product, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}
	return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) error {
	var newProduct models.Product
	if err := c.Bind(&newProduct); err != nil {
		return err
	}
	result := database.DB.Create(&newProduct)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	return c.JSON(http.StatusCreated, newProduct)
}

func UpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}

	if err := c.Bind(&product); err != nil {
		return err
	}
	database.DB.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func DeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	result := database.DB.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
	}
	return c.NoContent(http.StatusNoContent)
}
