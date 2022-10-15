package Entities

type Admin struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Nickname string `gorm:"size:255" json:"nickname" binding:"required"`
	Password string `gorm:"size:255" json:"password" binding:"required"`
}