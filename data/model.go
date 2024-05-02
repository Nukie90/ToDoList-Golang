package data

import (
	"time"

	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Task struct {
	Model
	Title    string `gorm:"type:varchar(50);not null" json:"title"`
	Desc     string `gorm:"type:varchar(100);not null" json:"desc"`
	Status   string `gorm:"type:varchar(32);not null" json:"status"`
	Deadline string `gorm:"type:date;not null" json:"deadline"`
	Privacy  string `gorm:"type:varchar(32);not null" json:"privacy"`
	Owner    string `gorm:"type:varchar(100);not null" json:"owner"`
}

type User struct {
	Model
	Username string `gorm:"type:varchar(100);not null" json:"username"`
	Password string `gorm:"type:varchar(100);not null" json:"password"`
}
