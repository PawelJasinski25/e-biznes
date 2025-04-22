package models

type Cart struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Items []CartItem `json:"items"`
}
