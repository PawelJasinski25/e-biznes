package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"server/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie udało się połączyć z bazą danych:", err)
	}
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Błąd migracji:", err)
	}
	log.Println("Połączono z bazą.")
}
