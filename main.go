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
	r.GET("/ws", route.HandleWebSocket)
	r.GET("/users/", route.AuthMiddleware())
	r.GET("/users/profile", route.GetProfile)
	r.PUT("/users/profile", route.UpdateUserProfile)
	r.GET("/users/chats", route.AuthMiddleware(), route.GetChats)
	r.GET("/users/chats/messages/:userID", route.AuthMiddleware(), route.GetMessages)
	r.GET("/users/:userID", route.GetUserInfo)
	config.ConnectDatabase()

	err := r.Run(":8000")
	if err != nil {
		return
	}
}
