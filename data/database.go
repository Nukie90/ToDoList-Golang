package data

import (
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() (*gorm.DB, error) {
	des := "postgres://nukie:Gamersking0@localhost:5432/todolist?sslmode=disable"
	sqlDB, err := sql.Open("pgx", des)
	if err != nil {
		return nil, errors.New("failed to connect database #1")
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database #2")
	}
	
	fmt.Println("Database connected")
	return gormDB, nil
}
