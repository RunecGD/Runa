package main

import (
	"Runa/api/config"
	"Runa/api/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Authorization", "Content-Type"},
	}))

	// Ваши маршруты
	r.POST("/register", route.Register)
	r.POST("/login", route.Login)
	r.GET("/users", route.AuthMiddleware(), route.GetUsers)
	r.GET("/ws", route.AuthMiddleware(), route.HandleWebSocket)
	config.ConnectDatabase()

	err := r.Run(":8000")
	if err != nil {
		return
	}
}
