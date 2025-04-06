package main

import (
	"log"
	"person-project/db"
	"person-project/handlers"

	"github.com/gin-gonic/gin"
)

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
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
