package main

import (
	"log"
	"tm/config"
	"tm/middleware"
	"tm/router"

	"github.com/gofiber/fiber/v2"
)

func main() {

	config.InitDatabase()
	defer config.DB.Close()

	app := fiber.New()
	app.Use(middleware.CookieMiddleware)
	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
