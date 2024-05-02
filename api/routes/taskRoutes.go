package routes

import (
	"main/api/handler"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TaskRoutes(app *fiber.App, db *gorm.DB) {
    task := fiber.New()

    app.Mount("/task", task)

    //create task
    task.Post("/create", func(c *fiber.Ctx) error {
        return handler.CreateTask(c, db)
    })

    //get all tasks
    task.Get("/all", func(c *fiber.Ctx) error {
        return handler.GetTasks(c, db)
    })

	//get current user tasks
	task.Get("/seetasks", func(c *fiber.Ctx) error {
		return handler.GetTasksOnSpecificUser(c, db)
	})

	//update task
	task.Put("/update/:id", func(c *fiber.Ctx) error {
		return handler.UpdateTask(c, db)
	})

	//delete task
	task.Delete("/delete/:id", func(c *fiber.Ctx) error {
		return handler.DeleteTask(c, db)
	})

	//get current user tasks on specific period of time
	task.Get("/seetasks/:time", func(c *fiber.Ctx) error {
		return handler.GetTasksOnSpecificPeriodOfTime(c, db)
	})
}