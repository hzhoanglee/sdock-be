package model

import "gorm.io/gorm"

// RoomShare struct
type RoomShare struct {
	gorm.Model
	RoomID uint `gorm:"not null" json:"room_id"`
	UserID uint `gorm:"not null" json:"user_id"`
	Role   int  `gorm:"not null, default:2" json:"role"`

	//	Relationship
	Room Room `gorm:"foreignKey:RoomID" json:"room"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}
