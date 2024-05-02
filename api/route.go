package api

import (
	"main/api/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    }) 
    //user
    UserRoutes(app, db)
    //task
    // TaskRoutes(app, db)

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

// func TaskRoutes(app *fiber.App, db *gorm.DB) {
//     panic("not implemented")
// }