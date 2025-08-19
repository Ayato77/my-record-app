package models

import (
	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
	Rating   int    `json:"rating"`
	Tags     []Tag  `gorm:"many2many:record_tag;"`
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
}
