package routes

import (
	"main/api/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoutes(app *fiber.App, db *gorm.DB) {
    user := fiber.New()

    app.Mount("/user", user)

    //create user
    user.Post("/create", func(c *fiber.Ctx) error {
        return handler.CreateUser(c, db)
    })

    //get all users
    user.Get("/all", func(c *fiber.Ctx) error {
        return handler.GetUsers(c, db)
    })

    //get current user profile
    user.Get("/profile", func(c *fiber.Ctx) error {
        return handler.GetUser(c, db)
    })
}