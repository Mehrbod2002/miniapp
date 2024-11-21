package routes

import (
	"mini-telegram/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/register", controllers.Register)
	router.GET("/points/:username", controllers.GetPoints)
	router.POST("/points/:username", controllers.AddPoints)

	return router
}
