// model/automation.go
package model

import (
	"gorm.io/gorm"
)

type Automation struct {
	gorm.Model
	Title       string                `json:"title"`
	Description string                `json:"description"`
	IsActive    bool                  `json:"is_active" gorm:"default:true"`
	OwnerID     uint                  `json:"owner_id"`
	Owner       User                  `json:"owner" gorm:"foreignKey:OwnerID"`
	Conditions  []AutomationCondition `json:"conditions" gorm:"foreignKey:AutomationID"`
	Actions     []AutomationAction    `json:"actions" gorm:"foreignKey:AutomationID"`
}

type AutomationCondition struct {
	gorm.Model
	AutomationID uint   `json:"automation_id"`
	Type         string `json:"type"`     // TIME, TEMPERATURE, HUMIDITY, DEVICE_STATUS
	Operator     string `json:"operator"` // EQUALS, GREATER_THAN, LESS_THAN, BETWEEN
	Value        string `json:"value"`
	Value2       string `json:"value2,omitempty"` // For BETWEEN operator
	DeviceID     *uint  `json:"device_id,omitempty"`
	Device       Device `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	RoomID       *uint  `json:"room_id,omitempty"`
	Room         Room   `json:"room,omitempty" gorm:"foreignKey:RoomID"`
}

type AutomationAction struct {
	gorm.Model
	AutomationID uint   `json:"automation_id"`
	DeviceID     uint   `json:"device_id"`
	Device       Device `json:"device" gorm:"foreignKey:DeviceID"`
	Action       string `json:"action"` // ON, OFF, SET_VALUE
	Value        string `json:"value,omitempty"`
}

type AutomationLog struct {
	gorm.Model
	AutomationID uint       `json:"automation_id"`
	Automation   Automation `json:"automation" gorm:"foreignKey:AutomationID"`
	Status       string     `json:"status"` // TRIGGERED, FAILED, SKIPPED
	Message      string     `json:"message"`
}
