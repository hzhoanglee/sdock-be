package model

import "gorm.io/gorm"

// Home struct
type Home struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"null" json:"description"`
	Status      int    `gorm:"not null, default:0" json:"status"`
	OwnerID     uint   `gorm:"not null" json:"owner_id"`
	Long        string `gorm:"default:0" json:"long"`
	Lat         string `gorm:"default:0" json:"lat"`

	//	Relationship
	Owner       User          `gorm:"foreignKey:OwnerID" json:"owner"`
	Rooms       []Room        `gorm:"foreignKey:HomeID" json:"rooms"`
	HomeSetting []HomeSetting `gorm:"foreignKey:HomeID" json:"home_setting"`
}
