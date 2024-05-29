package model

import "gorm.io/gorm"

// DeviceLog struct
type DeviceLog struct {
	gorm.Model
	DeviceID uint   `gorm:"not null" json:"device_id"`
	OwnerID  uint   `gorm:"default:0" json:"owner_id"`
	Value    string `gorm:"not null" json:"value"`

	Owner  User   `gorm:"foreignKey:OwnerID" json:"owner"`
	Device Device `gorm:"foreignKey:DeviceID" json:"device"`
}
