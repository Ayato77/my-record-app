package main

import (
	"my-record-app/routes"
	"my-record-app/database"
	"my-record-app/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

	router := gin.Default()
	database.ConnectDB()
	database.DB.AutoMigrate(&models.Record{})// auto migration for record
	database.DB.AutoMigrate(&models.User{})//auto migration for user
	routes.SetupRoutes(router)
	router.Run(":8080")
}