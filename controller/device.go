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
	messageStr := door.Door // Assuming fire.Fire is the message you want to send
	// Execute shell command
	cmd := exec.Command("./cmd_commands/sens.sh", messageStr, phoneNumStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error executing command: %v, output: %s", err, string(output)))
	}

	return c.JSON(door)
}

func MovementAlert(c *fiber.Ctx) error {

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
	messageStr := pir.Pir // Assuming fire.Fire is the message you want to send
	// Execute shell command
	cmd := exec.Command("./cmd_commands/sens.sh", messageStr, phoneNumStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return c.Status(500).SendString(fmt.Sprintf("Error executing command: %v, output: %s", err, string(output)))
	}
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

	return c.JSON(fire)
}

func TempUpdate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}
	var temp models.StateDev
	if err := c.BodyParser(&temp); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	_, err = config.DB.Exec("UPDATE statesensor SET temp = ?, hum = ? WHERE id = ?", temp.Temp, temp.Hum, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

// for flutter
type Database struct {
	DB *sql.DB
}

func DeviceState(c *websocket.Conn, db *Database) {
	rows, err := db.DB.Query("SELECT id,door, fire ,pir, temp,hum FROM statesensor")
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
