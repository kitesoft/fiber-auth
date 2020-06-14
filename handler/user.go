package handler

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/model"

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
	id := c.Params("id")
	db := database.DB

	var user model.User
	db.First(&user, id)
	if user.Username == "" {
		c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
		return
	}
	db.Delete(&user)
	c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
