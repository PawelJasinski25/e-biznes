package controllers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"server/database"
	"server/models"
)

func Register(c echo.Context) error {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Surname  string `json:"surname" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&input); err != nil {
		c.Logger().Errorf("Rejestracja: błąd bindowania danych: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("Próba rejestracji dla email: %s\n", input.Email)

	user := models.User{
		Email:   input.Email,
		Name:    input.Name,
		Surname: input.Surname,
	}

	if err := user.HashPassword(input.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się zahaszować hasła"})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć użytkownika"})
	}

	c.Logger().Infof("Rejestracja zakończona sukcesem dla %s", input.Email)
	return c.JSON(http.StatusOK, map[string]string{"message": "Rejestracja zakończona sukcesem"})
}

func Login(c echo.Context) error {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	log.Printf("Próba logowania dla email: %s", input.Email)

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		log.Printf("Logowanie: użytkownik %s nie został znaleziony", input.Email)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny login lub hasło"})
	}

	if !user.CheckPassword(input.Password) {
		log.Printf("Logowanie: błędne hasło dla email: %s", input.Email)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny login lub hasło"})
	}

	log.Printf("Logowanie zakończone sukcesem dla %s", input.Email)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Zalogowano pomyślnie",
		"user": map[string]interface{}{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"surname": user.Surname,
		},
	})
}
