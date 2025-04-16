package models

type Product struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price"`
}
