package model

import "gorm.io/gorm"

// Device struct
type Device struct {
	gorm.Model
	Title        string `gorm:"not null" json:"title"`
	Description  string `gorm:"not null" json:"description"`
	Status       int    `gorm:"not null, default:0" json:"status"`
	OwnerID      uint   `gorm:"not null" json:"owner_id"`
	SecretID     string `gorm:"not null" json:"secret_id"`
	RoomID       uint   `gorm:"not null" json:"room_id"`
	DeviceTypeID int    `gorm:"not null" json:"device_type_id"`
	LastSeen     string `gorm:"default null" json:"last_seen"`
	IP           string `gorm:"default null" json:"ip"`

	// Relationship
	Owner      User       `gorm:"foreignKey:OwnerID" json:"owner"`
	Room       Room       `gorm:"foreignKey:RoomID" json:"room"`
	DeviceType DeviceType `gorm:"foreignKey:DeviceTypeID" json:"device_type"`
}
