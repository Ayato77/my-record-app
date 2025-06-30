package models

import "time"

type Record struct {
	ID uint `gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    ImageURL  string    `json:"image_url"`
    Rating    int       `json:"rating"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}