package models

type CartItem struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Amount    uint `json:"amount"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
