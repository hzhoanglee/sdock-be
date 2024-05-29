package model

import (
	"gorm.io/gorm"
)

// DeviceType struct
type DeviceType struct {
	gorm.Model
	Name         string   `gorm:"not null" json:"name"`
	Code         string   `gorm:"not null, unique" json:"code"`
	Kind         string   `gorm:"not null" json:"kind"`
	InitialValue string   `gorm:"null" json:"initial_value"`
	Icon         string   `gorm:"null" json:"icon"`
	Device       []Device `gorm:"foreignKey:DeviceTypeID" json:"device"`
}
