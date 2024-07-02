package main

import (
	"server_room/database"
	"server_room/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.InitDatabase()
	app := fiber.New()
	routers.InitRouter(app)

	app.Listen(":3000")
}
