package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email              string    `json:"email" gorm:"unique"`
	Name               string    `json:"name"`
	Surname            string    `json:"surname"`
	Password           string    `json:"password"`
	GoogleAccessToken  string    `json:"googleAccessToken"`
	GoogleTokenExpires time.Time `json:"googleTokenExpires"`
	JWT                string    `json:"jwt"`
	GithubAccessToken  string    `json:"githubAccessToken"`
	GithubTokenExpires time.Time `json:"githubTokenExpires"`
}

func (u *User) HashPassword(password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
