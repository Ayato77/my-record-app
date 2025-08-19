package main

import (
	"my-record-app/database"
	"my-record-app/models"
	"my-record-app/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewDevelopment()

	if err := godotenv.Load(); err != nil {
		logger.Sugar().Fatalf("Error loading .env file:", err)
	}

	router := gin.Default()
	if err := database.ConnectDB(); err != nil {
		logger.Sugar().Fatalf("Error connecting the database failed ", err)
	}
	database.DB.AutoMigrate(&models.Record{}) // auto migration for record
	database.DB.AutoMigrate(&models.User{})   //auto migration for user
	routes.SetupRoutes(router, logger)
	router.Run(":8080")
}
