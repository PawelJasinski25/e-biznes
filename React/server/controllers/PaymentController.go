package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"server/models"
)

var currentPayment models.Payment

func AddPayment(c echo.Context) error {
	var payment models.Payment
	if err := c.Bind(&payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	payment.Amount = calculateCartTotal(currentCart)
	currentPayment = payment

	return c.JSON(http.StatusCreated, currentPayment)
}

func GetPayment(c echo.Context) error {
	return c.JSON(http.StatusOK, currentPayment)
}
