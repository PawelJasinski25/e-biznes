package models

type Product struct {
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	CategoryID uint     `json:"-"`
	Price      float64  `json:"price"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`
}
