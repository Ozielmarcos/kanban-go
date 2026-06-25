package main

import (
	"time"

	"github.com/Ozielmarcos/mytodolist/app/internal/handler"
	"github.com/Ozielmarcos/mytodolist/app/internal/utils/middleware"
	"github.com/Ozielmarcos/mytodolist/app/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectWithRetry()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // Adicione as URLs do seu front aqui
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cacheia a resposta do OPTIONS por 12 horas
	}))

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	// Users
	r.POST("/register", handler.RegisterHandler)
	r.POST("/login", handler.Login)
	r.POST("/refresh", handler.RefreshTokenHandler)

	//Stories
	auth.POST("/stories", handler.CreateStory)
	auth.GET("/stories", handler.GetStoriesByUser)
	auth.GET("/stories/:story_id/tasks", handler.GetTasksByStory)
	auth.PUT("/story/:id", handler.UpdateStory)
	auth.DELETE("/story/:id", handler.RemoveStory)

	//Tasks
	auth.GET("/tasks", handler.GetTasksByUser)
	auth.GET("/task/:id", handler.GetTaskById)
	auth.POST("/tasks", handler.CreateTask)
	auth.PUT("/task/:id", handler.UpdateTask)
	auth.PATCH("/task/:id/status", handler.UpdateTaskStatus)

	//Timers
	auth.PUT("/task/:id/start", handler.StarTimerHandler)
	auth.PUT("/task/:id/pause", handler.PauseTimerHandler)
	auth.PUT("/task/:id/resume", handler.ResumeTimerHandler)

	r.Run(":3005")
}
