package handlers

import (
	"errors"
	"log"
	"net/http"
	"person-project/db"
	"person-project/models"
	"person-project/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePerson(c *gin.Context) {
	var person models.Person
	err := c.ShouldBindJSON(&person)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Обогащение данных для имени %s", person.Name)
	person.Age = services.GetAge(person.Name)
	person.Gender = services.GetGender(person.Name)
	person.Nationality = services.GetNationality(person.Name)

	if err := db.DB.Create(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения данных" + err.Error()})
	}
	log.Printf("Человек %s успешно создан с обогащенными данными", person.Name)
	c.JSON(http.StatusCreated, person)

}

func DeletePerson(c *gin.Context) {
	var person models.Person
	strId := c.Param("id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный ID"})
	}
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if err := db.DB.Delete(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удалён"})

}

func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var updateDataPerson models.Person
	err := c.ShouldBindJSON(&updateDataPerson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateDataPerson.Name != "" {
		person.Name = updateDataPerson.Name

	}
	if updateDataPerson.Age != 0 {
		person.Age = updateDataPerson.Age

	}
	if updateDataPerson.Gender != "" {
		person.Gender = updateDataPerson.Gender
	}
	if updateDataPerson.Nationality != "" {
		person.Nationality = updateDataPerson.Nationality
	}
	if updateDataPerson.Surname != "" {
		person.Surname = updateDataPerson.Surname
	}
	if updateDataPerson.Patronymic != nil && *updateDataPerson.Patronymic != "" {
		person.Patronymic = updateDataPerson.Patronymic
	}
	if err := db.DB.Save(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}
	c.JSON(http.StatusOK, person)

}

func GetPersons(c *gin.Context) {

	var persons []models.Person
	db := db.DB.Model(&models.Person{})

	if name := c.Query("name"); name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname := c.Query("surname"); surname != "" {
		db = db.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if gender := c.Query("gender"); gender != "" {
		db = db.Where("gender ILIKE ?", gender)
	}
	if age := c.Query("age"); age != "" {
		db = db.Where("age = ?", age)
	}
	if minAge := c.Query("age_min"); minAge != "" {
		db = db.Where("age >= ?", minAge)
	}
	if maxAge := c.Query("age_max"); maxAge != "" {
		db = db.Where("age <= ?", maxAge)
	}

	strLimit := c.DefaultQuery("limit", "10")
	strOffset := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(strLimit)
	offset, _ := strconv.Atoi(strOffset)

	if err := db.Limit(limit).Offset(offset).Find(&persons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении людей", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, persons)
}
func GetPersonByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется id"})
	}
	var person models.Person
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Человек не найден"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, person)

}
