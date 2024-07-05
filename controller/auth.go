package controller

import (
	"time"
	"tm/config"
	"tm/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Register(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password required"})
	}

	_, err := config.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not register user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user registered successfully"})
}

func Login(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username and password required"})
	}

	// Kullanıcıyı veritabanında ara
	row := config.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", user.Username)

	storedUser := new(models.User)
	err := row.Scan(&storedUser.ID, &storedUser.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}

	// Şifre kontrolü
	if user.Password != storedUser.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid username or password"})
	}

	// JWT token oluşturma
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = storedUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Token'ı imzala ve string'e dönüştür
	t, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    t,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
		SameSite: "Strict", // Opsiyonel: Cross-site request koruması için
	}

	c.Cookie(&cookie)

	// Yanıtı JSON formatında döndürme
	return c.JSON(fiber.Map{
		"message": "login successful",
		"token":   t,
	})
}
