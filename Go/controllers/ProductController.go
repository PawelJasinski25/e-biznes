package controllers

import (
	"net/http"
	"strconv"

	"Go/models"
	"github.com/labstack/echo/v4"
)

var products = []models.Product{
	{ID: 1, Name: "Laptop", CategoryID: 1, Price: 3500.00},
	{ID: 2, Name: "Telefon", CategoryID: 1, Price: 2500.00},
	{ID: 3, Name: "Tablet", CategoryID: 1, Price: 1500.00},
	{ID: 4, Name: "T-shirt", CategoryID: 2, Price: 15.99},
	{ID: 5, Name: "Jacket", CategoryID: 2, Price: 85.00},
	{ID: 6, Name: "Apple", CategoryID: 3, Price: 2.50},
}

var idCounter uint = 7

func GetProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, products)
}

func GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, p := range products {
		if int(p.ID) == id {
			return c.JSON(http.StatusOK, p)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
}

func CreateProduct(c echo.Context) error {
	var newProduct models.Product
	if err := c.Bind(&newProduct); err != nil {
		return err
	}
	newProduct.ID = idCounter
	idCounter++
	products = append(products, newProduct)
	return c.JSON(http.StatusCreated, newProduct)
}

func UpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var updatedProduct models.Product
	if err := c.Bind(&updatedProduct); err != nil {
		return err
	}

	for i, p := range products {
		if int(p.ID) == id {
			updatedProduct.ID = p.ID
			products[i] = updatedProduct
			return c.JSON(http.StatusOK, updatedProduct)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
}

func DeleteProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, p := range products {
		if int(p.ID) == id {
			products = append(products[:i], products[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"message": "Product not found"})
}
