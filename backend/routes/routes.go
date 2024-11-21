package routes

import (
	"mini-telegram/controllers"
	"mini-telegram/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.New()

	// Use logging and recovery middleware
	router.Use(middleware.LoggerMiddleware(), gin.Recovery())

	api := router.Group("/backend")
	{
		api.POST("/register", controllers.Register)
		api.POST("/log", controllers.LogHandler)

		// Protected routes
		api.GET("/points/:telegram_id", controllers.GetPoints)
		api.POST("/points/:telegram_id", controllers.AddPoints)
	}

	return router
}
