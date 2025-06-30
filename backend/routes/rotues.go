package routes

import (
	"github.com/gin-gonic/gin"
	"my-record-app/controllers"
	"my-record-app/middleware"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/login", controllers.Login)
	router.POST("/register", controllers.Register)

	auth := router.Group("/")
	//Below is a protected route that requires a valid JWT token to access
	auth.Use(middleware.AuthRequired())
	{
		//auth.POST("/records", controllers.CreateRecord)
		auth.GET("/records", controllers.GetRecords)
		//auth.PUT("/records/:id", controllers.UpdateRecord)
		//auth.DELETE("/records/:id", controllers.DeleteRecord)
		//auth.GET("/records/:id", controllers.GetRecord)
	}

	router.Run(":8080")
}