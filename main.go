package main

import (
	"log"
	"tm/config"
	"tm/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.InitDatabase()
	defer config.DB.Close()

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":3002"))
}
