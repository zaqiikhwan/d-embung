package Entities

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	LinkFoto string `gorm:"size:255" json:"linkFoto"`
	Deskripsi string `gorm:"text" json:"deskripsi"`
}

type AdditionalInfo struct {
	gorm.Model
	Description string `gorm:"text" json:"desc"`
	LinkImages string `gorm:"text" json:"link"`
}
