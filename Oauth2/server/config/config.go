package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nie udało się załadować pliku .env:", err)
	} else {
		log.Println("Plik .env został załadowany")
	}
}
