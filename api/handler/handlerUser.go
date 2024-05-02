package handler

import (
	"main/data"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx, db *gorm.DB) error {
	user := new(data.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid input",
		})
	}
	db.Create(&user)
	return c.JSON(user)
}

func GetUsers(c *fiber.Ctx, db *gorm.DB) error {
	var users []data.User
	db.Find(&users)
	return c.JSON(users)
}

func Login(c *fiber.Ctx, db *gorm.DB) error {
	user := new(data.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid input",
		})
	}
	var userDB data.User
	if err := db.Where("username = ?", user.Username).First(&userDB).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "User not found",
		})
	}
	if user.Password != userDB.Password {
		return c.Status(400).JSON(fiber.Map{
			"password": user.Password,
			"db": userDB.Password,
			"message": "Invalid password",
		})
	}

	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	//set cookie
	c.Cookie(&fiber.Cookie{
		Name: "jwt",
		Value: tokenString,
		Expires: time.Now().Add(time.Hour * 72),
	})

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}

func GetUser(c *fiber.Ctx, db *gorm.DB) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	var user data.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("jwt")

	return c.JSON(fiber.Map{
		"message": "Logout successfully",
	})
}