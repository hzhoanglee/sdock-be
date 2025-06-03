package cmd

import (
	"app/database"
	"app/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net"
	"time"
)

func GetUserIDFromToken(c *fiber.Ctx) uint {
	return uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
}

func CheckRoomPermission(uid uint, r model.Room) bool {
	if r.OwnerID == uid {
		return true
	}
	// query in roomshare table
	roomShare := database.DB.Where("room_id = ? AND user_id = ?", r.ID, uid).First(&model.RoomShare{})
	if roomShare.RowsAffected > 0 {
		return true
	}
	return false
}

// SeedDeviceType func
func SeedDeviceType() {
	db := database.DB
	deviceTypes := []model.DeviceType{
		{Name: "Temperature Sensor", Code: "SENSOR_TEMPERATURE", Kind: "SENSOR", InitialValue: "0", Icon: "temperature"},
		{Name: "Humidity Sensor", Code: "SENSOR_HUMIDITY", Kind: "SENSOR", InitialValue: "0", Icon: "humidity"},
		{Name: "Light Sensor", Code: "SENSOR_LIGHT", Kind: "SENSOR", InitialValue: "0", Icon: "light"},
		{Name: "Door Sensor", Code: "SENSOR_DOOR", Kind: "SENSOR", InitialValue: "0", Icon: "light"},

		{Name: "Light Switch", Code: "SWITCH_LIGHT", Kind: "SWITCH", InitialValue: "0", Icon: "light"},
		{Name: "Fan Switch", Code: "SWITCH_FAN", Kind: "SWITCH", InitialValue: "0", Icon: "fan"},
		{Name: "AC Switch", Code: "SWITCH_AC", Kind: "SWITCH", InitialValue: "0", Icon: "ac"},
		{Name: "TV Switch", Code: "SWITCH_TV", Kind: "SWITCH", InitialValue: "0", Icon: "tv"},
		{Name: "Door Switch", Code: "SWITCH_DOOR", Kind: "SWITCH", InitialValue: "0", Icon: "door"},
	}
	for _, deviceType := range deviceTypes {
		db.Create(&deviceType)
	}
}

func GetCurrentTime() string {
	return time.Now().Format(time.DateTime)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
