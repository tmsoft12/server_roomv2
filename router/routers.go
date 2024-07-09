package router

import (
	"tm/controller"
	"tm/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, db *controller.Database) {
	// CORS middleware
	app.Use(middleware.CORSMiddleware())

	// for login
	app.Get("/", controller.TestConnection)
	app.Post("/login", controller.Login)
	app.Post("/register", controller.Register)

	// JWT middleware
	app.Use(middleware.JWTMiddleware)

	// for telefon
	app.Get("/phone", controller.GetPhone)
	app.Put("/phone/:id", controller.UpdatePhone)

	// for device
	app.Put("/open_door/:id", controller.OpenDoor)
	app.Put("/movement_alert/:id", controller.MovementAlert)
	app.Put("/fire_alrt/:id", controller.FireAler)

	// WebSocket
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		controller.DeviceState(c, db)
	}))

}
