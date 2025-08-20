package routes

import (
	"my-record-app/controllers"
	"my-record-app/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(router *gin.Engine, logger *zap.Logger) {
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	auth := router.Group("/")
	//Below is a protected route that requires a valid JWT token to access
	auth.Use(middleware.AuthRequired())
	{
		auth.POST("/records", controllers.CreateRecord(logger))
		auth.GET("/records", controllers.GetRecords(logger))
		//auth.PUT("/records/:id", controllers.UpdateRecord)
		auth.DELETE("/records/:id", controllers.DeleteRecord(logger))
		//auth.GET("/records/:id", controllers.GetRecord)
	}

	router.Run(":8080")
}
