package main

import (
	"log"
	"person-project/config"
	"person-project/db"
	_ "person-project/docs"
	"person-project/handlers"
	"person-project/logger"
	"person-project/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title People API
// @version 1.0
// @description API для создания и получения информации о людях
// @host localhost:8080
// @BasePath /
func main() {

	config.LoadConfig()
	logger.Init()
	db.Connect()

	db.Migrate()
	r := gin.Default()
	r.Use(middleware.LoggingMiddleware())

	person := r.Group("/people")
	{
		person.POST("", handlers.CreatePerson)
		person.GET("", handlers.GetPersons)
		person.GET("/:id", handlers.GetPersonByID)
		person.PATCH("/:id", handlers.UpdatePerson)
		person.DELETE("/:id", handlers.DeletePerson)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
