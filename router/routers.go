package router

import (
	"tm/controller"
	"tm/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App, db *controller.Database) {
	// CORS middleware'ini uygula
	app.Use(middleware.CORSMiddleware())

	// Temel rotaları ayarla
	app.Get("/", controller.TestConnection)
	app.Post("/login", controller.Login)
	app.Post("/register", controller.Register)

	// JWT middleware'ini uygula
	app.Use(middleware.JWTMiddleware)

	// Telefon rotalarını ayarla
	app.Get("/phone", controller.GetPhone)
	app.Put("/phone/:id", controller.UpdatePhone)

	// Cihaz rotalarını ayarla
	app.Get("/open_door", controller.OpenDoor)
	app.Get("/movement_alert", controller.MovementAlert)

	// WebSocket rotasını ayarla
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		controller.DeviceState(c, db)
	}))
}
