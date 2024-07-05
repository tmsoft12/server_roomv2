package controller

import (
	"os/exec"

	"github.com/gofiber/fiber/v2"
)

func OpenDoor(c *fiber.Ctx) error {
	cmd := exec.Command("python3", "python_scripts/t.py")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	// fmt.Println(string(out))
	return c.JSON(fiber.Map{
		"message": string(out),
	})
}
func MovementAlert(c *fiber.Ctx) error {
	cmd := exec.Command("python3", "python_scripts/t.py")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	// fmt.Println(string(out))
	return c.JSON(fiber.Map{
		"message": string(out),
	})
}
