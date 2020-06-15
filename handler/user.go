package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateUser new user
func CreateUser(c *fiber.Ctx) {
	type NewUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	db := database.DB
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
		return
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
		return
	}

	user.Password = hash
	if err := db.Create(&user).Error; err != nil {
		c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
		return
	}

	newUser := NewUser{
		Email:    user.Email,
		Username: user.Username,
	}

	c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": newUser})
}

// DeleteUser delete user
func DeleteUser(c *fiber.Ctx) {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
		return
	}
	id := c.Params("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn'n read id from params", "data": err})
		return
	}

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))
	db := database.DB
	var user model.User

	if uid != n {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Couldn't delete user", "data": nil})
		return
	}

	db.First(&user, id)
	if user.Username == "" {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
		return
	}

	if !CheckPasswordHash(pi.Password, user.Password) {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
		return
	}

	db.Delete(&user)
	c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
