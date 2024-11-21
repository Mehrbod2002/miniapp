package main

import (
	"mini-telegram/config"
	"mini-telegram/models"
	"mini-telegram/routes"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})
	router := routes.SetupRoutes()
	router.Run(":8080")
}
