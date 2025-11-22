// backend/config/database.go

package config

import (
	"fmt"
	"log"
	"os"
	"task-manager/models"
	"task-manager/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// بررسی متغیرهای محیطی
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration is incomplete. Check .env file.")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		utils.LogError("Failed to connect to database", err)
		log.Fatal("Failed to connect to database:", err)
	}

	utils.LogInfo("Database connected successfully")

	// Migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.Streak{},
		&models.Group{},
		&models.GroupMember{},
		&models.Task{},
		&models.TaskProgress{},
		&models.GroupTaskProgress{},
		&models.TaskAssignment{},
		&models.File{},
		&models.Notification{},
	)

	if err != nil {
		utils.LogError("Failed to migrate database", err)
		log.Fatal("Failed to migrate database:", err)
	}

	utils.LogInfo("Database migration completed successfully")
}

// GetDB برای دسترسی به database instance
func GetDB() *gorm.DB {
	return DB
}
