package controller

import (
	"strconv"
	"tm/config"
	"tm/models"

	"github.com/gofiber/fiber/v2"
)

func GetPhone(c *fiber.Ctx) error {
	rows, err := config.DB.Query("SELECT id ,phone FROM phone")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	var phones []models.Phone
	for rows.Next() {
		var phone models.Phone
		if err := rows.Scan(&phone.ID, &phone.Phone); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		phones = append(phones, phone)
	}
	return c.JSON(phones)
}
func UpdatePhone(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).SendString("Invalid ID")
	}

	var phone models.Phone
	if err := c.BodyParser(&phone); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	db := config.DB
	_, err = db.Exec("UPDATE phone SET phone = ? WHERE id = ?", phone.Phone, id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	phone.ID = id
	return c.JSON(phone)
}
