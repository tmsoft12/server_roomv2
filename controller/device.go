package controller

import (
	"database/sql"
	"log"
	"os/exec"
	"tm/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func TestConnection(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"messgae": "oki",
	})
}
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

type Database struct {
	DB *sql.DB
}

func GetUsers(c *websocket.Conn, db *Database) {
	rows, err := db.DB.Query("SELECT id,door, fire,pir FROM statesensor")
	if err != nil {
		log.Println("failed to query users:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dev models.StateDev
		err := rows.Scan(&dev.ID, &dev.Door, &dev.Fire, &dev.Pir)
		if err != nil {
			log.Println("failed to scan user:", err)
			return
		}
		if err := c.WriteJSON(dev); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
