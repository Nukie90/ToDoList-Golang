package handler

import (
	"fmt"
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
	deadLine := task.Deadline
	deadlineTime, err := time.Parse("2006-01-02", deadLine)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	now := time.Now()
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if deadlineTime.Before(todayMidnight) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Deadline is before today",
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
	deadLine := task.Deadline
	deadlineTime, _ := time.Parse("2006-01-02", deadLine)
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

func GetTasksOnSpecificPeriodOfTime(c *fiber.Ctx, db *gorm.DB) error {
	timeString := c.Params("time")
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
	switch timeString {
	case "today":
		username := claims["username"].(string)
		var tasks []data.Task
		now := time.Now()
		todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		fmt.Println(todayMidnight)
		// query for view
		query := db.Model(&data.Task{}).Where("deadline = ?", todayMidnight)
		//create view
		err := db.Migrator().CreateView("today_table", gorm.ViewOption{
			Query: query,
			Replace: true,
		})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		//select from view
		db.Raw(`
			select * from today_table
			where owner = ? or privacy = 'Public';
		`, username).Scan(&tasks)

		db.Migrator().DropView("today_table")

		return c.JSON(tasks)
	case "week":
		username := claims["username"].(string)
		var tasks []data.Task
		now := time.Now()
		weekEnd := now.AddDate(0, 0, 6)
		// query for view
		query := db.Model(&data.Task{}).Where("deadline between ? and ?", now, weekEnd)
		//create view
		err := db.Migrator().CreateView("week_table", gorm.ViewOption{
			Query: query,
			Replace: true,
		})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		//select from view
		db.Raw(`
			select * from week_table
			where owner = ? or privacy = 'Public';
		`, username).Scan(&tasks)

		db.Migrator().DropView("week_table")

		return c.JSON(tasks)

	case "month":
		username := claims["username"].(string)
		var tasks []data.Task
		now := time.Now()
		monthEnd := now.AddDate(0, 1, -1)
		// query for view
		query := db.Model(&data.Task{}).Where("deadline between ? and ?", now, monthEnd)
		//create view
		err := db.Migrator().CreateView("month_table", gorm.ViewOption{
			Query: query,
			Replace: true,
		})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		//select from view
		db.Raw(`
			select * from month_table
			where owner = ? or privacy = 'Public';
		`, username).Scan(&tasks)

		db.Migrator().DropView("month_table")

		return c.JSON(tasks)
	default:
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid time",
		})
	}
}
