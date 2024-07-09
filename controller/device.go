package controller

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
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
	cmd := exec.Command("python_scripts/sens.sh 'mesaj' '99361570538'")
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
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
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

	rows, err := config.DB.Query("SELECT id, phone FROM phone")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	var phones []models.Phone
	for rows.Next() {
		var phone models.Phone
		if err := rows.Scan(&phone.ID, &phone.Phone); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		phones = append(phones, phone)
	}

	// Prepare phone numbers for command execution
	var phoneNums []string
	for _, p := range phones {
		phoneNums = append(phoneNums, p.Phone)
	}
	phoneNumStr := strings.Join(phoneNums, ",")

	// Prepare message for shell command execution
	messageStr := fire.Fire // Assuming fire.Fire is the message you want to send
	// Execute shell command
	cmd := exec.Command("./cmd_commands/sens.sh", messageStr, phoneNumStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error executing command: %v, output: %s", err, string(output)))
	}

	// Return JSON response with fire object
	return c.JSON(fire)
}

// for flutter
type Database struct {
	DB *sql.DB
}

func DeviceState(c *websocket.Conn, db *Database) {
	rows, err := db.DB.Query("SELECT id,door, fire,pir, temp,hum FROM statesensor")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var dev models.StateDev
		err := rows.Scan(&dev.ID, &dev.Door, &dev.Fire, &dev.Pir, &dev.Temp, &dev.Hum)
		if err != nil {

			return
		}
		if err := c.WriteJSON(dev); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
