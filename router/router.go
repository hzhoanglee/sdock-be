package router

import (
	"app/handler"
	"app/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
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
	deviceParing := api.Group("/device_paring", middleware.Protected())
	deviceParing.Get("/scan", handler.DoScanDevice)
	deviceParing.Get("/type", handler.GetTypeDevice)

	// Static
	app.Static("/", "./public")

	// Product
	//product := api.Group("/product")
	//product.Get("/", handler.GetAllProducts)
	//product.Get("/:id", handler.GetProduct)
	//product.Post("/", middleware.Protected(), handler.CreateProduct)
	//product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)
}
