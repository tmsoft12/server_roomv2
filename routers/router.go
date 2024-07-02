package routers

import (
	"server_room/controller"

	"github.com/gofiber/fiber/v2"
)

func InitRouter(app *fiber.App) {
	app.Get("/", controller.SendCommand)
	app.Put("phone/:id", controller.UpdatePhone)
}
