package models

import "time"

type Product struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CategoryID uint      `json:"category_id"`
	Name       string    `json:"name"`
	Color      string    `json:"color"`
	Weight     float64   `json:"weight"`
	PricePerKg float64   `json:"price_per_kg"`
	Stock      int       `json:"stock"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
// Buat file terpisah untuk User dan Transaction dengan cara yang sama...