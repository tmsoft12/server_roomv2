package router

import (
	"tm/controller"
	"tm/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Use(middleware.CORSMiddleware())

	app.Post("/login", controller.Login)
	app.Post("/register", controller.Register)

	app.Use(middleware.JWTMiddleware)

	//phone
	app.Get("/phone", controller.GetPhone)
	app.Post("phone/:id", controller.UpdatePhone)

	// Device
	app.Get("/open_door", controller.OpenDoor)

}
