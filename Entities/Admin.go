package Entities

type Admin struct {
	ID uint `gorm:"primaryKey" json:"id"`
	NamaLengkap string `gorm:"size:255" json:"nama_lengkap"`
	Password string `gorm:"size:255" json:"password"`
}