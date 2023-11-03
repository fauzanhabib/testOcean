package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb()(*gorm.DB, error){
	// Replace with your PostgreSQL database connection information
	dsn := "user=postgres password=postgres dbname=oceantest sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
