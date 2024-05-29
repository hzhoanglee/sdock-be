package model

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

// Room struct
type Room struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Status      int    `gorm:"not null, default:0" json:"status"`
	HomeID      uint   `gorm:"not null" json:"home_id"`
	OwnerID     uint   `gorm:"not null" json:"owner_id"`
	Image       string `gorm:"default: /images/home.jpg" json:"image"`
	Owner       User   `gorm:"foreignKey:OwnerID" json:"owner"`
}

func (r *Room) FullImageURL(domain string) string {
	// Check if Image already has a full URL
	if r.Image != "" && (strings.HasPrefix(r.Image, "http://") || strings.HasPrefix(r.Image, "https://")) {
		return r.Image
	}
	// Otherwise, concatenate domain and image path
	return fmt.Sprintf("%s%s", domain, r.Image)
}
