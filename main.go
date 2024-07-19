package main

import (
	"database/sql"
	"log"
	"tm/config"
	"tm/controller"
	"tm/middleware"
	"tm/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	config.InitDatabase()
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(config.DB)

	app := fiber.New()

	app.Use(logger.New())
	app.Use(middleware.CookieMiddleware)

	dbController := &controller.Database{DB: config.DB}

	router.SetupRoutes(app, dbController)

	log.Fatal(app.Listen(":3000"))
}
