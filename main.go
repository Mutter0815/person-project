package main

import (
	"log"
	"person-project/db"
	"person-project/handlers"

	_ "person-project/docs"

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

	db.Connect()

	db.Migrate()
	r := gin.Default()

	person := r.Group("/people")
	{
		person.POST("", handlers.CreatePerson)
		person.GET("", handlers.GetPersons)
		person.GET("/:id", handlers.GetPersonByID)
		person.PUT("/:id", handlers.UpdatePerson)
		person.DELETE("/:id", handlers.DeletePerson)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
