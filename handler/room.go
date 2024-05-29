package handler

import (
	"app/cmd"
	"app/database"
	"app/model"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAllRooms query all Rooms
func GetAllRooms(c *fiber.Ctx) error {
	uid := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	homeID := c.Params("home_id")
	db := database.DB
	var Home model.Home
	db.Find(&Home, homeID)
	if Home.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})
	}
	if uid != Home.OwnerID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this home", "data": nil})
	}

	var Rooms []model.Room
	db.Find(&Rooms).Where("home_id = ?", homeID)
	return c.JSON(fiber.Map{"status": "success", "message": "All Rooms", "data": Rooms})
}

// GetRoom query Room
func GetRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var Room model.Room
	db.Find(&Room, id)
	if Room.OwnerID != uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64)) {
		var RoomShare model.RoomShare
		db.Find(&RoomShare, "room_id = ? AND user_id = ?", id, c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
		if RoomShare.ID == 0 {
			return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this room", "data": nil})
		}
	}
	if Room.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Room found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Room found", "data": Room})
}

func GetAllDeviceRoom(c *fiber.Ctx) error {
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
	var outDevices []model.Device
	db.Preload("Owner").Preload("Room").Preload("DeviceType").Where("room_id = ?", roomID).Find(&devices)
	if deviceType != "" {
		for i, device := range devices {
			if device.DeviceType.Kind == deviceType {
				outDevices = append(outDevices, devices[i])
			}
		}
	}

	return c.JSON(fiber.Map{"status": "success", "message": "All Devices", "data": outDevices})
}

// CreateRoom new Room
func CreateRoom(c *fiber.Ctx) error {
	db := database.DB
	uid := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	homeID := c.Params("home_id")
	var Home model.Home
	db.Find(&Home, homeID)
	if Home.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})
	}
	if uid != Home.OwnerID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this home", "data": nil})
	}
	fmt.Println(homeID)
	Room := new(model.Room)
	Room.OwnerID = uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	Room.Status = 1
	Room.HomeID = Home.ID
	if err := c.BodyParser(Room); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create Room", "data": err})
	}
	db.Create(&Room)
	if Room.ID == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create Room", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Created Room", "data": Room})
}

func ShareRoom(c *fiber.Ctx) error {
	payload := struct {
		RoomID    int    `json:"room_id"`
		UserEmail string `json:"email"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	db := database.DB

	var Room model.Room
	db.Find(&Room, payload.RoomID)
	if Room.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Room found with ID: " + payload.UserEmail, "data": nil})
	}
	uid := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	if uid != Room.OwnerID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this room", "data": nil})
	}
	var User model.User
	db.Find(&User, "email = ?", payload.UserEmail)
	if User.ID == 0 {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No User found with Email: " + payload.UserEmail, "data": nil})
	}
	var Share model.RoomShare
	Share.RoomID = Room.ID
	Share.UserID = User.ID
	db.Find(&Share, "room_id = ? AND user_id = ?", Room.ID, User.ID)
	if Share.ID != 0 {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Room already shared with this user", "data": nil})
	}
	db.Create(&Share)
	if Share.ID == 0 {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't share Room", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Shared Room", "data": Share})

}

// DeleteRoom delete Room
func DeleteRoom(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var Room model.Room
	db.First(&Room, id)
	if Room.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Room found with ID", "data": nil})

	}
	db.Delete(&Room)
	return c.JSON(fiber.Map{"status": "success", "message": "Room successfully deleted", "data": nil})
}
