// handler/automation.go
package handler

import (
	"app/cmd"
	"app/database"
	"app/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

// GetAllAutomations retrieves all automations for a user
func GetAllAutomations(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)

	var automations []model.Automation
	db := database.DB

	db.Preload("Owner").Preload("Conditions").Preload("Actions").
		Where("owner_id = ?", uid).Find(&automations)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "All automations",
		"data":    automations,
	})
}

// CreateAutomation creates a new automation rule
func CreateAutomation(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)

	var automation model.Automation
	if err := c.BodyParser(&automation); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	automation.OwnerID = uid

	// Validate conditions and actions
	if err := validateAutomation(&automation, uid); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	db := database.DB
	if err := db.Create(&automation).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create automation",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Automation created",
		"data":    automation,
	})
}

// UpdateAutomation updates an existing automation
func UpdateAutomation(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	automationID := c.Params("automation_id")

	var automation model.Automation
	db := database.DB

	if err := db.First(&automation, automationID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Automation not found",
			"data":    nil,
		})
	}

	if automation.OwnerID != uid {
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not the owner of this automation",
			"data":    nil,
		})
	}

	var updates model.Automation
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	// Update fields
	automation.Title = updates.Title
	automation.Description = updates.Description
	automation.IsActive = updates.IsActive

	if err := db.Save(&automation).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update automation",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Automation updated",
		"data":    automation,
	})
}

// DeleteAutomation deletes an automation
func DeleteAutomation(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	automationID := c.Params("automation_id")

	var automation model.Automation
	db := database.DB

	if err := db.First(&automation, automationID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Automation not found",
			"data":    nil,
		})
	}

	if automation.OwnerID != uid {
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not the owner of this automation",
			"data":    nil,
		})
	}

	// Delete related conditions and actions
	db.Where("automation_id = ?", automationID).Delete(&model.AutomationCondition{})
	db.Where("automation_id = ?", automationID).Delete(&model.AutomationAction{})
	db.Delete(&automation)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Automation deleted",
		"data":    nil,
	})
}

// GetAutomationLogs retrieves logs for an automation
func GetAutomationLogs(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	automationID := c.Params("automation_id")

	var automation model.Automation
	db := database.DB

	if err := db.First(&automation, automationID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Automation not found",
			"data":    nil,
		})
	}

	if automation.OwnerID != uid {
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "You are not the owner of this automation",
			"data":    nil,
		})
	}

	var logs []model.AutomationLog
	db.Where("automation_id = ?", automationID).
		Order("created_at DESC").
		Limit(100).
		Find(&logs)

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Automation logs",
		"data":    logs,
	})
}

// validateAutomation validates automation conditions and actions
func validateAutomation(automation *model.Automation, uid uint) error {
	db := database.DB

	// Validate conditions
	for _, condition := range automation.Conditions {
		// Validate condition type
		validTypes := []string{"TIME", "TEMPERATURE", "HUMIDITY", "DEVICE_STATUS"}
		if !contains(validTypes, condition.Type) {
			return fmt.Errorf("invalid condition type: %s", condition.Type)
		}

		// Validate operator
		validOperators := []string{"EQUALS", "GREATER_THAN", "LESS_THAN", "BETWEEN"}
		if !contains(validOperators, condition.Operator) {
			return fmt.Errorf("invalid operator: %s", condition.Operator)
		}

		// Validate device access if device-based condition
		if condition.DeviceID != nil {
			var device model.Device
			db.Preload("Room").First(&device, *condition.DeviceID)
			if device.ID == 0 {
				return fmt.Errorf("device not found")
			}
			if !cmd.CheckRoomPermission(uid, device.Room) {
				return fmt.Errorf("no permission for device")
			}
		}
	}

	// Validate actions
	for _, action := range automation.Actions {
		var device model.Device
		db.Preload("Room").First(&device, action.DeviceID)
		if device.ID == 0 {
			return fmt.Errorf("device not found")
		}
		if !cmd.CheckRoomPermission(uid, device.Room) {
			return fmt.Errorf("no permission for device")
		}

		// Validate action type
		validActions := []string{"ON", "OFF", "SET_VALUE"}
		if !contains(validActions, action.Action) {
			return fmt.Errorf("invalid action: %s", action.Action)
		}
	}

	return nil
}

// CheckAutomations checks all active automations and triggers actions if conditions are met
func CheckAutomations() {
	db := database.DB

	var automations []model.Automation
	db.Preload("Conditions.Device").Preload("Conditions.Room").
		Preload("Actions.Device").
		Where("is_active = ?", true).
		Find(&automations)

	for _, automation := range automations {
		if checkAutomationConditions(automation) {
			executeAutomationActions(automation)
		}
	}
}

// checkAutomationConditions checks if all conditions for an automation are met
func checkAutomationConditions(automation model.Automation) bool {
	for _, condition := range automation.Conditions {
		if !checkCondition(condition) {
			return false
		}
	}
	return true
}

// checkCondition checks if a single condition is met
func checkCondition(condition model.AutomationCondition) bool {
	switch condition.Type {
	case "TIME":
		return checkTimeCondition(condition)
	case "TEMPERATURE":
		return checkSensorCondition(condition, "temperature")
	case "HUMIDITY":
		return checkSensorCondition(condition, "humidity")
	case "DEVICE_STATUS":
		return checkDeviceStatusCondition(condition)
	default:
		return false
	}
}

// checkTimeCondition checks time-based conditions
func checkTimeCondition(condition model.AutomationCondition) bool {
	now := time.Now()

	switch condition.Operator {
	case "EQUALS":
		// Format: "15:30" for 3:30 PM
		conditionTime, err := time.Parse("15:04", condition.Value)
		if err != nil {
			return false
		}
		return now.Format("15:04") == conditionTime.Format("15:04")

	case "BETWEEN":
		// Format: "08:00" and "17:00"
		startTime, err1 := time.Parse("15:04", condition.Value)
		endTime, err2 := time.Parse("15:04", condition.Value2)
		if err1 != nil || err2 != nil {
			return false
		}

		currentTime := now.Format("15:04")
		return currentTime >= startTime.Format("15:04") && currentTime <= endTime.Format("15:04")

	default:
		return false
	}
}

// checkSensorCondition checks sensor-based conditions
func checkSensorCondition(condition model.AutomationCondition, sensorType string) bool {
	if condition.DeviceID == nil {
		return false
	}

	deviceStatus, err := getDeviceStatus(int(*condition.DeviceID))
	if err != nil {
		return false
	}

	// Parse sensor value
	currentValue, err := strconv.ParseFloat(deviceStatus.Value, 64)
	if err != nil {
		return false
	}

	conditionValue, err := strconv.ParseFloat(condition.Value, 64)
	if err != nil {
		return false
	}

	switch condition.Operator {
	case "EQUALS":
		return currentValue == conditionValue
	case "GREATER_THAN":
		return currentValue > conditionValue
	case "LESS_THAN":
		return currentValue < conditionValue
	case "BETWEEN":
		conditionValue2, err := strconv.ParseFloat(condition.Value2, 64)
		if err != nil {
			return false
		}
		return currentValue >= conditionValue && currentValue <= conditionValue2
	default:
		return false
	}
}

// checkDeviceStatusCondition checks device status conditions
func checkDeviceStatusCondition(condition model.AutomationCondition) bool {
	if condition.DeviceID == nil {
		return false
	}

	deviceStatus, err := getDeviceStatus(int(*condition.DeviceID))
	if err != nil {
		return false
	}

	switch condition.Operator {
	case "EQUALS":
		return strings.ToLower(deviceStatus.Value) == strings.ToLower(condition.Value)
	default:
		return false
	}
}

// executeAutomationActions executes all actions for an automation
func executeAutomationActions(automation model.Automation) {
	for _, action := range automation.Actions {
		var status int
		var value string

		switch action.Action {
		case "ON":
			status = 1
			value = "on"
		case "OFF":
			status = 0
			value = "off"
		case "SET_VALUE":
			status = 1
			value = action.Value
		}

		// Execute the action
		err := setIOTStatusToDevice(int(action.DeviceID), status, value)
		if err != nil {
			// Log error
			logAutomation(automation.ID, "FAILED", fmt.Sprintf("Failed to execute action: %v", err))
			continue
		}

		// Update device status
		_ = setDeviceStatus(int(action.DeviceID), status, value, nil)
	}

	// Log successful execution
	logAutomation(automation.ID, "TRIGGERED", "Automation executed successfully")
}

// logAutomation creates a log entry for automation execution
func logAutomation(automationID uint, status, message string) {
	db := database.DB

	log := model.AutomationLog{
		AutomationID: automationID,
		Status:       status,
		Message:      message,
	}

	db.Create(&log)
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
