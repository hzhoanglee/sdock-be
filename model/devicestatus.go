package model

import "gorm.io/gorm"

// DeviceStatus struct
type DeviceStatus struct {
	gorm.Model
	DeviceID int    `gorm:"not null" json:"device_id"`
	Status   int    `gorm:"not null, default:0" json:"status"`
	Value    string `gorm:"default null" json:"value"`

	// Relationship
	Device Device `gorm:"foreignKey:DeviceID" json:"device"`
}
