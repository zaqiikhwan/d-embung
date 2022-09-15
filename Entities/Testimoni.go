package Entities

import "gorm.io/gorm"

type Testimoni struct {
	gorm.Model
	Identitas string `gorm:"size:255" json:"identitas"`
	Testimoni string `json:"testimoni"`
}