package models

type CartItem struct {
	ID        uint    `json:"id"`
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Amount    uint    `json:"amount"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}
