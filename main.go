package main

import (
	"fmt"
	"log"
	"main/api"
	"main/data"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func main() {
	db, err := data.DbConnection()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&data.Task{})
	db.AutoMigrate(&data.User{})

	app := fiber.New()

	api.SetupRoutes(app, db)
	api.SetupMiddleware(app)

	fmt.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}