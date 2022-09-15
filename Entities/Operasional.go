package Entities

import (
	"gorm.io/gorm"
)

type Operasional struct {
	gorm.Model 
	HariOperasional string `gorm:"size:255" json:"hariOperasional"`
	JamOperasional string `gorm:"size:255" json:"jamOperasional"`
	Harga           string `gorm:"size:255" json:"harga"`
}
