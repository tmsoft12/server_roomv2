package controller

import (
	"database/sql"
	"log"
	"os/exec"
	"strconv"
	"tm/config"
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
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invaild ID")
	}
	var door models.StateDev
	if err := c.BodyParser(&door); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db := config.DB
	_, err = db.Exec("UPDATE statesensor SET door = ? WHERE id = ?", door.Door, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	door.ID = id

	return c.JSON(door)
}

func MovementAlert(c *fiber.Ctx) error {
	cmd := exec.Command("python3", "python_scripts/t.py")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invaild ID")
	}
	var pir models.StateDev
	if err := c.BodyParser(&pir); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db := config.DB
	_, err = db.Exec("UPDATE statesensor SET pir = ? WHERE id = ?", pir.Pir, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	pir.ID = id

	return c.JSON(pir)
}

func FireAler(c *fiber.Ctx) error {
	cmd := exec.Command("python3", "python_scripts/t.py")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invaild ID")
	}
	var fire models.StateDev
	if err := c.BodyParser(&fire); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db := config.DB
	_, err = db.Exec("UPDATE statesensor SET fire = ? WHERE id = ?", fire.Fire, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	fire.ID = id

	return c.JSON(fire)
}

// for flutter
type Database struct {
	DB *sql.DB
}

func DeviceState(c *websocket.Conn, db *Database) {
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
