package api

import (
	"main/api/handler"
    "main/api/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    }) 
    //user
    routes.UserRoutes(app, db)
    //task
    routes.TaskRoutes(app, db)

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Post("/login", func(c *fiber.Ctx) error {
        return handler.Login(c, db)
    })

    app.Get("/logout", func(c *fiber.Ctx) error {
        return handler.Logout(c)
    })
}