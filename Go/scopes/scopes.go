package scopes

import (
	"gorm.io/gorm"
)

func PreloadCategory(db *gorm.DB) *gorm.DB {
	return db.Preload("Category")
}

func PreloadItemsProduct(db *gorm.DB) *gorm.DB {
	return db.Preload("Items.Product")
}

func PreloadProducts(db *gorm.DB) *gorm.DB {
	return db.Preload("Products")
}
