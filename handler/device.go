package handler

import (
	"app/cmd"
	"app/database"
	"app/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// GetAllRoomDevice get all devices
// User Action
func GetAllRoomDevice(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	roomID := c.Params("room_id")
	deviceType := c.Query("kind")
	deviceType = strings.ToUpper(deviceType)
	if deviceType != "SENSOR" && deviceType != "SWITCH" {
		deviceType = ""
	}

	var room = model.Room{}
	db := database.DB
	db.Find(&room, roomID)
	if room.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Room found with ID", "data": nil})
	}
	if !cmd.CheckRoomPermission(uid, room) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this room", "data": nil})
	}
	var devices []model.Device
	db.Preload("Owner").Preload("Room").Preload("DeviceType").Find(&devices).Where("room_id = ?", roomID)
	if deviceType != "" {
		for i, device := range devices {
			if device.DeviceType.Kind != deviceType {
				devices = append(devices[:i], devices[i+1:]...)
			}
		}
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Devices", "data": devices})
}

func RegisterDevice(c *fiber.Ctx) error {
	var count int64
	db := database.DB
	db.Model(&model.DeviceType{}).Count(&count)
	if count == 0 {
		cmd.SeedDeviceType()
	}
	uid := cmd.GetUserIDFromToken(c)
	device := model.Device{}

	if err := c.BodyParser(&device); err != nil {
		return err
	}
	device.OwnerID = uid
	device.SecretID = utils.UUIDv4()

	var room = model.Room{}
	db.Find(&room, device.RoomID)
	if room.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Room found with ID", "data": nil})
	}
	if !cmd.CheckRoomPermission(uid, room) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this room", "data": nil})
	}

	var deviceType = model.DeviceType{}
	db.Find(&deviceType, device.DeviceTypeID)
	if deviceType.Name == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device Type found with ID", "data": nil})
	}
	device.DeviceTypeID = int(deviceType.ID)
	fmt.Println(device.IP)
	sendServerConfig(device.IP, device.SecretID)
	db.Create(&device)
	// Create device status
	deviceStatus := model.DeviceStatus{
		DeviceID: int(device.ID),
		Status:   0,
	}

	db.Create(&deviceStatus)
	if device.ID == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to create device", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Device created", "data": device})

}

func DeleteDevice(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	deviceID := c.Params("device_id")
	var device = model.Device{}
	db := database.DB
	db.Find(&device, deviceID)
	if device.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device found with ID", "data": nil})
	}
	if device.OwnerID != uid {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this device", "data": nil})
	}
	db.Delete(&device)
	return c.JSON(fiber.Map{"status": "success", "message": "Device deleted", "data": nil})
}

func GetStatusDevice(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	deviceID := c.Params("device_id")
	var device = model.Device{}
	db := database.DB
	db.Preload("Room").Find(&device, deviceID)
	if device.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device found with ID", "data": nil})
	}
	if !cmd.CheckRoomPermission(uid, device.Room) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this device", "data": nil})
	}
	deviceIDInt, _ := strconv.Atoi(deviceID)
	deviceStatus, _ := getDeviceStatus(deviceIDInt)
	return c.JSON(fiber.Map{"status": "success", "message": "Device status", "data": deviceStatus})
}

func GetLogDevice(c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	deviceID := c.Params("device_id")
	var device = model.Device{}
	db := database.DB
	db.Preload("Room").Find(&device, deviceID)
	if device.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device found with ID", "data": nil})
	}
	if !cmd.CheckRoomPermission(uid, device.Room) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this device", "data": nil})
	}
	var deviceLogs []model.DeviceLog
	db.Preload("Owner").Preload("Device").Preload("Device.Owner").Preload("Device.Room").Preload("Device.DeviceType").Preload("Device.Room.Owner").Where("device_id = ?", deviceID).Find(&deviceLogs)
	return c.JSON(fiber.Map{"status": "success", "message": "Device logs", "data": deviceLogs})
}

func SetStatusDevice(c *fiber.Ctx) error {
	//uid := cmd.GetUserIDFromToken(c)
	//payload := struct {
	//	Status   int `json:"status"`
	//	DeviceID int `json:"device_id"`
	//}{}
	//if err := c.BodyParser(&payload); err != nil {
	//	return err
	//}
	// Rethink later, don't know why but really slow
	//rdb := database.Rethink()
	//err := rdb.Insert(payload)
	//if err != nil {
	//	return err
	//}
	uid := cmd.GetUserIDFromToken(c)
	var DeviceStatus = model.DeviceStatus{}
	if err := c.BodyParser(&DeviceStatus); err != nil {
		return err
	}
	db := database.DB
	var Device = model.Device{}
	db.Preload("DeviceType").Find(&Device, DeviceStatus.DeviceID)
	if Device.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device found with ID", "data": nil})
	}
	if Device.OwnerID != uid && !cmd.CheckRoomPermission(uid, Device.Room) {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this device", "data": nil})
	}
	//return c.JSON(fiber.Map{"status": "success", "message": "Device status updated", "data": nil})
	if Device.DeviceType.Kind == "SWITCH" {
		if DeviceStatus.Status == 1 {
			DeviceStatus.Value = "on"
		} else {
			DeviceStatus.Value = "off"
		}
		err := setIOTStatusToDevice(DeviceStatus.DeviceID, DeviceStatus.Status, DeviceStatus.Value)
		if err != nil {
			return err
		}
	}
	d := setDeviceStatus(DeviceStatus.DeviceID, DeviceStatus.Status, DeviceStatus.Value, c)
	if d != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Failed to update device status", "data": d.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Device status updated", "data": nil})
}

func GetTypeDevice(c *fiber.Ctx) error {
	var kind = c.Query("kind")
	if kind != "SENSOR" && kind != "SWITCH" {
		kind = ""
	}
	kind = strings.ToUpper(kind)
	var deviceTypes []model.DeviceType
	db := database.DB
	if kind != "" {
		db.Where("kind = ?", kind).Find(&deviceTypes)
	} else {
		db.Find(&deviceTypes)
	}
	return c.JSON(fiber.Map{"status": "success", "message": "All Device Types", "data": deviceTypes})
}

// Updating Functions
func getDeviceStatus(deviceId int) (model.DeviceStatus, error) {
	db := database.DB
	var deviceStatus = model.DeviceStatus{}
	db.Preload("Device").Preload("Device.Room").Preload("Device.DeviceType").Where("device_id = ?", deviceId).First(&deviceStatus)
	return deviceStatus, nil
}

func storeDeviceLog(deviceId, uid int, value string) error {
	db := database.DB
	deviceLog := model.DeviceLog{
		DeviceID: uint(deviceId),
		OwnerID:  uint(uid),
		Value:    value,
	}
	err := db.Create(&deviceLog).Error
	if err != nil {
		return err
	}
	return nil
}

// User and Device actions
func setDeviceStatus(deviceId, status int, value string, c *fiber.Ctx) error {
	uid := cmd.GetUserIDFromToken(c)
	db := database.DB
	var deviceStatus model.DeviceStatus
	err := db.Where("device_id = ?", deviceId).First(&deviceStatus, model.DeviceStatus{DeviceID: deviceId}).Error
	if err != nil {
		return err
	}
	// if old status is the same as new status, don't update
	if deviceStatus.Status == status && deviceStatus.Value == value {
		return nil
	}
	deviceStatus.Status = status
	deviceStatus.Value = value
	err = db.Save(&deviceStatus).Error
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	_ = storeDeviceLog(deviceId, int(uid), value)

	return nil
}

// PingFromDevice Handle IOT
func PingFromDevice(c *fiber.Ctx) error {
	secret := c.FormValue("secret")
	ip := c.FormValue("ip")
	if secret == "" || ip == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Missing required fields"})
	}
	var device = model.Device{}
	db := database.DB
	db.Find(&device, model.Device{SecretID: secret})
	if device.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Device found with Secret", "data": nil})
	}
	device.Status = 1
	device.LastSeen = cmd.GetCurrentTime()
	device.IP = ip
	db.Save(&device)
	return c.JSON(fiber.Map{"status": "success", "message": "Pong"})
}

func setIOTStatusToDevice(deviceId, status int, value string) error {
	device := model.Device{}
	db := database.DB
	db.Find(&device, deviceId)
	if device.Title == "" {
		return errors.New("no Device found with ID")
	}
	deviceSecret := device.SecretID
	deviceIP := device.IP
	requestURL := "http://" + deviceIP + "/lamp?state=" + value
	println(requestURL)
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth("sdock", deviceSecret)

	resp, err := http.DefaultClient.Do(req)
	println(resp)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	return nil
}

func DoScanDevice(c *fiber.Ctx) error {
	devices := ScanForDevices()
	return c.JSON(fiber.Map{"status": "success", "message": "Devices found", "data": devices})
}

func sendServerConfig(remoteDeviceIP, secretID string) string {
	// CCU Server = http://<localIP>:3000/api/v1/local_ip
	// perform run to get local IP
	thisDeviceIP := cmd.GetLocalIP()
	fmt.Println("Local IP: ", thisDeviceIP)
	ccuEndpoint := "http://" + thisDeviceIP + ":3000/api/v1/local_ip"
	fmt.Println(secretID)
	// Create data for the POST request
	data := map[string]string{"server": ccuEndpoint, "secret": secretID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "Error in JSON marshalling"
	}
	// Create a new request
	req, err := http.NewRequest("POST", "http://"+remoteDeviceIP+":80/server-config", bytes.NewBuffer(jsonData))
	if err != nil {
		return "Error in creating request"
	}
	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Error in sending request"
	}
	defer resp.Body.Close()
	return "Request sent successfully"
}
