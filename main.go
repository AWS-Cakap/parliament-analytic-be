package main

import (
	"fmt"
	"parliament-analytic-be/config"
	"parliament-analytic-be/models"
	"parliament-analytic-be/routes"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Initialize DB
	config.InitDB()

	// Auto migrate tables
	config.DB.AutoMigrate(
		&models.Admin{},
		&models.Party{},
		&models.Tweet{},
	)

	seedAdmin() // call admin function

	println("Database OK!")

	// Run routing
	r := gin.Default()
	r.Static("/media", "./media")
	routes.SetupRoutes(r)
	r.Run(":8080") // run in localhost
}

// Insert admin
func seedAdmin() {
	var count int64
	config.DB.Model(&models.Admin{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.Admin{
			Username: "admin",
			Password: string(hashedPassword),
		}
		config.DB.Create(&admin)
		fmt.Println("Admin has been created: admin/admin123")
	} else {
		fmt.Println("Default Admin already exists")
	}
}
