package controllers

import (
	"net/http"
	"strconv"

	"Go/database"
	"Go/models"
	"Go/scopes"
	"github.com/labstack/echo/v4"
)

func GetCarts(c echo.Context) error {
	var carts []models.Cart
	err := database.DB.Scopes(scopes.PreloadItemsProduct).Find(&carts).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	var result []map[string]interface{}

	for _, cart := range carts {
		cartData := map[string]interface{}{
			"id":    cart.ID,
			"items": []map[string]interface{}{},
		}

		for _, item := range cart.Items {
			itemData := map[string]interface{}{
				"product_id":   item.ProductID,
				"product_name": item.Product.Name,
				"amount":       item.Amount,
			}
			cartData["items"] = append(cartData["items"].([]map[string]interface{}), itemData)
		}

		result = append(result, cartData)
	}

	return c.JSON(http.StatusOK, result)
}

func GetCart(c echo.Context) error {
	id := c.Param("id")

	var cart models.Cart
	err := database.DB.Scopes(scopes.PreloadItemsProduct).First(&cart, id).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}

	cartData := map[string]interface{}{
		"id":    cart.ID,
		"items": []map[string]interface{}{},
	}

	for _, item := range cart.Items {
		itemData := map[string]interface{}{
			"product_id":   item.ProductID,
			"product_name": item.Product.Name,
			"amount":       item.Amount,
		}
		cartData["items"] = append(cartData["items"].([]map[string]interface{}), itemData)
	}

	return c.JSON(http.StatusOK, cartData)
}

type CreateCartInput struct {
	Items []struct {
		ProductID uint `json:"product_id"`
		Amount    int  `json:"amount"`
	} `json:"items"`
}

func CreateCart(c echo.Context) error {
	var input CreateCartInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	cart := models.Cart{}
	if err := database.DB.Create(&cart).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	var items []map[string]interface{}
	for _, item := range input.Items {
		cartItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: item.ProductID,
			Amount:    uint(item.Amount),
		}
		database.DB.Create(&cartItem)

		var product models.Product
		database.DB.Scopes(scopes.PreloadCategory).First(&product, item.ProductID)

		items = append(items, map[string]interface{}{
			"product_id":   item.ProductID,
			"product_name": product.Name,
			"amount":       item.Amount,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    cart.ID,
		"items": items,
	})
}

func UpdateCart(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var cart models.Cart
	if err := database.DB.First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Cart not found"})
	}

	var input CreateCartInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	var items []map[string]interface{}
	for _, item := range input.Items {
		newItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: item.ProductID,
			Amount:    uint(item.Amount),
		}
		database.DB.Create(&newItem)

		var product models.Product
		database.DB.First(&product, item.ProductID)

		items = append(items, map[string]interface{}{
			"product_id":   item.ProductID,
			"product_name": product.Name,
			"amount":       item.Amount,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":    cart.ID,
		"items": items,
	})
}

func DeleteCart(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	database.DB.Where("cart_id = ?", id).Delete(&models.CartItem{})
	result := database.DB.Delete(&models.Cart{}, id)
	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Cart not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
