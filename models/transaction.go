package models

import "time"

type TransactionIn struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"` // Relasi ke Product
	Date      string    `json:"date"` // Format YYYY-MM-DD
	Quantity  float64   `json:"quantity"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type TransactionOut struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"` // Relasi ke Product
	Date      string    `json:"date"` // Format YYYY-MM-DD
	Quantity  float64   `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}