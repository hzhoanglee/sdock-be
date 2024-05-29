package model

import "gorm.io/gorm"

// HomeSetting struct
type HomeSetting struct {
	gorm.Model
	HomeID   string `gorm:"not null" json:"home_id"`
	HomeInfo Home   `gorm:"foreignKey:HomeID" json:"home_info"`
	Key      string `gorm:"not null" json:"key"`
	Value    string `gorm:"not null" json:"value"`
}
