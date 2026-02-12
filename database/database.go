package database

import (
	"log"

	"tcc-tech/queue-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate
	if err := db.AutoMigrate(&models.QueueState{}, &models.QueueTicket{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed singleton QueueState row if not exists
	var count int64
	db.Model(&models.QueueState{}).Count(&count)
	if count == 0 {
		db.Create(&models.QueueState{
			CurrentLetter: "",
			CurrentNumber: -1,
		})
	}

	log.Println("Database connected and migrated successfully")
	return db
}
