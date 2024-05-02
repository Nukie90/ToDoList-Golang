package handler

import (
	"main/data"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func CreateTask(c *fiber.Ctx, db *gorm.DB) error {
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
	userID := claims["username"].(string)
	task := new(data.Task)

	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	task.Owner = userID
	deadLine := task.Deadline.String() // Convert deadLine to string
	deadlineTime, _ := time.Parse("2006/01/02", deadLine)
	if deadlineTime.Before(time.Now()) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Have passed deadline",
		})
	}

	task.Status = "Incompleted"
	switch task.Privacy {
		case strconv.Itoa(1):
			task.Privacy = "Public"
		case strconv.Itoa(2):
			task.Privacy = "Private"
		default:
			task.Privacy = "Public"
	}
	db.Create(&task)
	return c.JSON(task)
}

func GetTasks(c *fiber.Ctx, db *gorm.DB) error {
	var tasks []data.Task
	db.Find(&tasks)
	return c.JSON(tasks)
}


func UpdateTask(c *fiber.Ctx, db *gorm.DB) error {
		id := c.Params("id")
		task := new(data.Task)
		if err := db.First(&task, id).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		deadLine := task.Deadline.String()
		deadlineTime, _ := time.Parse("2006/01/02", deadLine)
		if deadlineTime.Before(time.Now()) {
			return c.Status(400).JSON(fiber.Map{
				"message": "Have passed deadline , cannot update task",
			})
		}

		switch task.Status {
			case strconv.Itoa(1):
				task.Status = "Completed"
			case strconv.Itoa(2):
				task.Status = "Cancelled"
			default:
				task.Status = "Incompleted"
		}
		db.Save(&task)
		return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("id")
	task := new(data.Task)
	if err := db.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Task not found",
		})
	}
	db.Delete(&task)
	return c.SendString("Task deleted")
}

func GetTasksOnSpecificUser(c *fiber.Ctx, db *gorm.DB) error {
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
	var tasks []data.Task
	db.Where("owner = ? or privacy = 'Public'", username).Find(&tasks)
	return c.JSON(tasks)
}