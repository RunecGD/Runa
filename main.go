package main

import (
	"Runa/config" // Убедитесь, что путь правильный
	"Runa/route"  // Убедитесь, что путь правильный
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	router := gin.Default()
	router.POST("/register", route.Register)
	router.POST("/login", route.Login)
	router.GET("/users", route.AuthMiddleware(), route.GetUsers)
	router.Run(":8000")
}
