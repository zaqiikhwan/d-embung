package Entities

import (
	"gorm.io/gorm"
)

type Operasional struct {
	gorm.Model 
	HariOperasional string `gorm:"size:255" json:"day"`
	JamOperasional 	string `gorm:"size:255" json:"hour"`
	Harga           string `gorm:"size:255" json:"price"`
}
