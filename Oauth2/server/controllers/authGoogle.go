package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server/config"
	jwtutil "server/utils"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"server/database"
	"server/models"
)

var googleOauthConfig *oauth2.Config
var oauthStateString = "pseudo-random"

func init() {
	config.LoadEnv()
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type GoogleUserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GoogleLogin(c echo.Context) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateString {
		log.Printf("Niepoprawny stan autoryzacji, oczekiwano '%s', otrzymano '%s'\n", oauthStateString, state)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawny stan autoryzacji"})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Brak kodu autoryzacyjnego"})
	}

	oauthToken, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Wymiana kodu nie powiodła się: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd wymiany kodu autoryzacyjnego"})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + oauthToken.AccessToken)
	if err != nil {
		log.Printf("Błąd pobierania danych użytkownika: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd pobierania danych"})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Błąd odczytu odpowiedzi: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd odczytu odpowiedzi"})
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		log.Printf("Błąd parsowania danych użytkownika: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd przetwarzania danych"})
	}

	log.Printf("Otrzymane dane z Google: %+v\n", userInfo)

	var user models.User
	err = database.DB.Where("email = ?", userInfo.Email).First(&user).Error
	if err != nil {
		user = models.User{
			Email:   userInfo.Email,
			Name:    userInfo.GivenName,
			Surname: userInfo.FamilyName,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("Błąd zapisu użytkownika do bazy: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd zapisu do bazy"})
		}
	}
	user.GoogleAccessToken = oauthToken.AccessToken
	user.GoogleTokenExpires = oauthToken.Expiry

	jwtToken, err := jwtutil.GenerateJWT(user)
	if err != nil {
		log.Printf("Błąd generowania JWT: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd generowania tokena"})
	}
	user.JWT = jwtToken
	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("Błąd zapisania danych w bazie: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd zapisu danych w bazie"})
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Printf("Błąd pobierania sesji: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd sesji"})
	}
	sess.Values["jwt"] = user.JWT
	sess.Values["userID"] = user.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Printf("Błąd zapisywania sesji: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd zapisu sesji"})
	}

	return c.Redirect(http.StatusTemporaryRedirect, "http://localhost:5173/home")
}
