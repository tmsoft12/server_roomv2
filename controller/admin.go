package controller

import (
	"database/sql"
	"fmt"
	"server_room/models"

	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

func UpdatePhone(c *fiber.Ctx) error {

	// Request body'den JSON verisini oku
	var req models.PhoneModel
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	if req.ID == 0 || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).SendString("ID and phone are required")
	}

	// SQL UPDATE sorgusunu hazırla
	stmt, err := db.Prepare("UPDATE phone SET phone = ? WHERE id = ?")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	defer stmt.Close()

	// Sorguyu çalıştır
	res, err := stmt.Exec(req.Phone, req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Etkilenen satır sayısını al
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).SendString("No record found with the given ID")
	}

	return c.SendString(fmt.Sprintf("Record with ID %d updated successfully", req.ID))
}
