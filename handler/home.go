package handler

import (
	"app/database"
	"app/model"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

// GetAllHomes query all Homes
func GetAllHomes(c *fiber.Ctx) error {
	db := database.DB
	var Homes []model.Home
	db.Where("owner_id = ?", c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64)).Preload("HomeSetting").Preload("Rooms").Find(&Homes)
	return c.JSON(fiber.Map{"status": "success", "message": "All Homes", "data": Homes})
}

// GetHome query Home
func GetHome(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var Home model.Home
	db.Find(&Home, id)
	db.Preload("Owner").Preload("HomeSetting").Preload("Rooms").Find(&Home, id)
	if Home.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})
	}
	uid := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	if uid != Home.OwnerID {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "You are not the owner of this home", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Home found", "data": Home})
}

// CreateHome new Home
func CreateHome(c *fiber.Ctx) error {
	db := database.DB
	Home := new(model.Home)
	Home.OwnerID = uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	if err := c.BodyParser(Home); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create Home", "data": err})
	}
	db.Create(&Home)
	return c.JSON(fiber.Map{"status": "success", "message": "Created Home", "data": Home})
}

// GetHomeSettings query Home
func GetHomeSettings(c *fiber.Ctx) error {
	id := c.Params("id")
	uid := uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	db := database.DB
	var Home model.Home
	db.Find(&Home, id).Where("owner_id = ?", uid)
	if Home.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Home found", "data": Home})
}

// DeleteHome delete Home
func DeleteHome(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	var Home model.Home
	db.First(&Home, id)
	if Home.Title == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})

	}
	db.Delete(&Home)
	return c.JSON(fiber.Map{"status": "success", "message": "Home successfully deleted", "data": nil})
}

//func HomeTransfer(c *fiber.Ctx) error {
//	id := c.Params("id")
//	db := database.DB
//
//	var Home model.Home
//	db.First(&Home, id)
//	if Home.Title == "" {
//		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No Home found with ID", "data": nil})
//
//	}
//	chec
//	Home.OwnerID = uint(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
//	db.Save(&Home)
//	return c.JSON(fiber.Map{"status": "success", "message": "Home successfully transfered", "data": Home})
//}
