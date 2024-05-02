package api

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// SetupMiddleware function is used to set up the middleware for the application
func SetupMiddleware(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey :jwtware.SigningKey{Key: []byte("secret")},
	}))
}

