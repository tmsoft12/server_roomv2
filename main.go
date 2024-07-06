package main

import (
	"log"
	"tm/config"
	"tm/controller"
	"tm/middleware"
	"tm/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Veritabanı bağlantısını başlat
	config.InitDatabase()
	defer config.DB.Close()

	// Fiber uygulamasını oluştur
	app := fiber.New()

	// Middleware'leri uygula
	app.Use(logger.New())
	app.Use(middleware.CookieMiddleware)

	// Kontrolcü veritabanı bağlantısını ayarla
	dbController := &controller.Database{DB: config.DB}

	// Rotaları yapılandır
	router.SetupRoutes(app, dbController)

	// Fiber uygulamasını dinle
	log.Fatal(app.Listen(":3000"))
}
