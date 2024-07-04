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

	var phones []models.Phone
	if err := c.BodyParser(&phones); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	db := config.DB

	tx, err := db.Begin()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer tx.Rollback()

	for _, phone := range phones {
		stmt, err := tx.Prepare("UPDATE phone SET phone = ? WHERE id = ?")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer stmt.Close()

		_, err = stmt.Exec(phone.Phone, id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	phones[0].ID = id
	return c.JSON(phones[0])
}
