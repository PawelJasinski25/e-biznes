package models

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"-"`
}
