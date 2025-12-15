package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email" gorm:"unique"`
	FullName  string    `json:"full_name"`
	Avatar    string    `json:"avatar"`
	Role      string    `json:"role"`   // 'admin' or 'staff'
	Status    string    `json:"status"` // 'active' or 'pending'
	CreatedAt time.Time `json:"created_at"`
}