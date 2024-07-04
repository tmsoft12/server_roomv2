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
	app.Get("/test", controller.Test)
	// app.Get("/private", controller.Private)

}
