package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server/config"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"server/database"
	"server/models"
)

var githubOauthConfig *oauth2.Config
var oauthStateStringGithub = "pseudo-random-github"

func init() {
	config.LoadEnv()
	githubOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/github/callback",
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GithubLogin(c echo.Context) error {
	url := githubOauthConfig.AuthCodeURL(oauthStateStringGithub)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func GithubCallback(c echo.Context) error {
	state := c.QueryParam("state")
	if state != oauthStateStringGithub {
		log.Printf("Niepoprawny stan autoryzacji GitHub, oczekiwano '%s', otrzymano '%s'\n", oauthStateStringGithub, state)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawny stan autoryzacji"})
	}

	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Brak kodu autoryzacyjnego"})
	}

	token, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Wymiana kodu nie powiodła się (GitHub): %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd wymiany kodu autoryzacyjnego"})
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("Błąd tworzenia zapytania do GitHub: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd tworzenia zapytania"})
	}
	req.Header.Set("Authorization", "token "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Błąd pobierania danych użytkownika z GitHub: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd pobierania danych"})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Błąd odczytu odpowiedzi (GitHub): %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd odczytu odpowiedzi"})
	}

	var githubUser GithubUser
	if err := json.Unmarshal(body, &githubUser); err != nil {
		log.Printf("Błąd parsowania danych użytkownika z GitHub: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd przetwarzania danych"})
	}

	log.Printf("Otrzymane dane z GitHub: %+v\n", githubUser)

	if githubUser.Email == "" {
		emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emailReq.Header.Set("Authorization", "token "+token.AccessToken)
		emailResp, err := http.DefaultClient.Do(emailReq)
		if err != nil {
			log.Printf("Błąd pobierania emaili z GitHub: %s\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd pobierania emaili"})
		}
		defer emailResp.Body.Close()
		emailBody, _ := ioutil.ReadAll(emailResp.Body)
		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err := json.Unmarshal(emailBody, &emails); err == nil && len(emails) > 0 {
			githubUser.Email = emails[0].Email
		}
	}

	var user models.User
	err = database.DB.Where("email = ?", githubUser.Email).First(&user).Error
	if err != nil {
		user = models.User{
			Email:   githubUser.Email,
			Name:    githubUser.Name,
			Surname: "",
		}
		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("Błąd zapisu użytkownika do bazy (GitHub): %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd zapisu do bazy"})
		}
	}

	ownToken := fmt.Sprintf("own-token-for-user-%d", user.ID)
	redirectURL := fmt.Sprintf("http://localhost:5173/home?token=%s", ownToken)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
