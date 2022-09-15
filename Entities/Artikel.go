package Entities

import "gorm.io/gorm"

type Artikel struct {
	gorm.Model
	Title   string `gorm:"size:255" json:"title"`
	Slug    string `gorm:"size:255" json:"slug"`
	Image   string `gorm:"size:255" json:"image"`
	Excerpt string `gorm:"size:255" json:"excerpt"`
	Body    string `gorm:"type:longtext" json:"body"`
}
