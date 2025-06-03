package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	//allow cors
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.Next()
	})

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Post("/api/v1/local_ip", handler.PingFromDevice)

	// Middleware
	api := app.Group("/api/v1", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.CreateUser)

	// User
	user := api.Group("/user")
	user.Get("/info/:id", handler.GetUser)
	user.Get("/me", middleware.Protected(), handler.GetUserMe)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Home
	home := api.Group("/home", middleware.Protected())
	home.Get("/info/:id", handler.GetHome)
	home.Get("/all", handler.GetAllHomes)
	home.Post("/create", handler.CreateHome)

	// Room
	room := api.Group("/room", middleware.Protected())
	room.Get("/all/:home_id", handler.GetAllRooms)
	room.Get("/info/:id", handler.GetRoom)
	room.Post("/create", handler.CreateRoom)
	room.Post("/share", handler.ShareRoom)
	room.Get("/devices/all/:room_id", handler.GetAllDeviceRoom)

	// Device
	device := api.Group("/device", middleware.Protected())
	device.Get("/all/:room_id", handler.GetAllRoomDevice)
	device.Post("/register", handler.RegisterDevice)
	device.Delete("/delete/:device_id", handler.DeleteDevice)
	device.Post("/status/set", handler.SetStatusDevice)
	device.Get("/status/get/:device_id", handler.GetStatusDevice)
	device.Get("/log/get/:device_id", handler.GetLogDevice)

	// DeviceParing
	devicePairing := api.Group("/device_pairing", middleware.Protected())
	devicePairing.Get("/scan", handler.DoScanDevice)
	devicePairing.Get("/type", handler.GetTypeDevice)

	// routes/automation.go
	app.Get("/api/v1/automations", handler.GetAllAutomations)
	app.Post("/api/v1/automations", handler.CreateAutomation)
	app.Put("/api/v1/automations/:automation_id", handler.UpdateAutomation)
	app.Delete("/api/v1/automations/:automation_id", handler.DeleteAutomation)
	app.Get("/api/v1/automations/:automation_id/logs", handler.GetAutomationLogs)

	// Static
	app.Static("/", "./public")

	// Product
	//product := api.Group("/product")
	//product.Get("/", handler.GetAllProducts)
	//product.Get("/:id", handler.GetProduct)
	//product.Post("/", middleware.Protected(), handler.CreateProduct)
	//product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)
}
