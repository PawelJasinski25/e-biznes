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

	oauthToken, err := githubOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Wymiana kodu nie powiodła się (GitHub): %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd wymiany kodu autoryzacyjnego"})
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("Błąd tworzenia zapytania do GitHub: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd tworzenia zapytania"})
	}
	req.Header.Set("Authorization", "token "+oauthToken.AccessToken)

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

	var ghUser GithubUser
	if err := json.Unmarshal(body, &ghUser); err != nil {
		log.Printf("Błąd parsowania danych użytkownika z GitHub: %s\n", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd przetwarzania danych"})
	}

	log.Printf("Otrzymane dane z GitHub: %+v\n", ghUser)

	if ghUser.Email == "" {
		emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emailReq.Header.Set("Authorization", "token "+oauthToken.AccessToken)
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
			ghUser.Email = emails[0].Email
		}
	}

	var user models.User
	err = database.DB.Where("email = ?", ghUser.Email).First(&user).Error
	if err != nil {
		user = models.User{
			Email:   ghUser.Email,
			Name:    ghUser.Name,
			Surname: "",
		}
		if err := database.DB.Create(&user).Error; err != nil {
			log.Printf("Błąd zapisu użytkownika do bazy (GitHub): %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd zapisu do bazy"})
		}
	}

	user.GithubAccessToken = oauthToken.AccessToken
	user.GithubTokenExpires = oauthToken.Expiry

	jwtToken, err := jwtutil.GenerateJWT(user)
	if err != nil {
		log.Printf("Błąd generowania JWT: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd generowania tokena"})
	}
	user.JWT = jwtToken

	if err := database.DB.Save(&user).Error; err != nil {
		log.Printf("Błąd zapisania danych w bazie (GitHub): %v\n", err)
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
