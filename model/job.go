package model

import "gorm.io/gorm"

// Job struct
type Job struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Task     string `gorm:"not null" json:"task"`
	Retry    int    `gorm:"not null, default: 0" json:"retry"`
	MaxRetry int    `gorm:"not null, default: 3" json:"max_retry"`
	Status   int    `gorm:"not null, default: 0" json:"status"`
	JobUUID  string `gorm:"not null" json:"job_uuid"`
}
