package controller

import (
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func SendCommand(ctx *fiber.Ctx) error {
	cmd := exec.Command("python3", "python_scripts/t.py")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	// fmt.Println(string(out))
	return ctx.JSON(fiber.Map{

		"message": string(out),
	})

}
