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

// CreatePerson godoc
// @Summary Создание нового человека
// @Description Создаёт нового человека с обогащением данных
// @Tags person
// @Accept json
// @Produce json
// @Param person body models.CreatePersonRequest true "Данные нового человека (имя и фамилия обязательны)"
// @Success 201 {object} models.Person "Пользователь успешно создан"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /people [post]
func CreatePerson(c *gin.Context) {
	var personRequest models.Person
	err := c.ShouldBindJSON(&personRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Обогащение данных для имени %s", personRequest.Name)
	var person models.Person
	person.Name = personRequest.Name
	person.Surname = personRequest.Surname
	person.Patronymic = personRequest.Patronymic
	person.Age = services.GetAge(person.Name)
	person.Gender = services.GetGender(person.Name)
	person.Nationality = services.GetNationality(person.Name)

	if err := db.DB.Create(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения данных" + err.Error()})
	}
	log.Printf("Человек %s успешно создан с обогащенными данными", person.Name)
	c.JSON(http.StatusCreated, person)

}

// DeletePerson godoc
// @Summary Удаление пользователя
// @Description Удаляет пользователя по id
// @Tags person
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /people/{id} [delete]
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

// UpdatePerson godoc
// @Summary Обновление данных пользователя
// @Description Обновляет данные пользователя по id
// @Tags person
// @Accept json
// @Param id path string true "ID пользователя"
// @Param input body models.Person false "Данные для обновления"
// @Produce json
// @Success 200 {object} models.Person "Обновленные данные пользователя"
// @Failrue 404 {object} models.ErrorResponse
// @Failrue 500 {object} models.ErrorResponse
// @Router /people/{id} [put]
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

// GetPersons godoc
// @Summary Возращает список пользователей
// @Description Возращает список пользователей, c возможностью фильтрации по имени,фамилии возрасту,и полу,также с пагинацией
// @Tags person
// @Produce json
// @Param name query string false "Фильтр по имени"
// @Param surname query string false "Фильтр по фамилии"
// @Param gender query string false "Фильтр по полу"
// @Param age query integer false "Фильтр по точному возрасту"
// @Param age_min query integer false "Минимальный возраст (включительно)"
// @Param age_max query integer false "Максимальный возраст (включительно)"
// @Param limit query integer false "Лимит записей (по умолчанию 10)"
// @Param offset query integer false "Смещение (по умолчанию 0)"
// @Success 200  {array} models.Person "Cписок людей"
// @Failrue 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /people [get]
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

// GetPersonByID godoc
// @Summary Возращает данные о пользователе
// @Description Возращает данные пользователе по id
// @Tags person
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} models.Person"Данные о пользователе"
// @Failrue 404 models.ErrorResponse "Пользователь не найден"
// @Router /people/{id} [get]
func GetPersonByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется id"})
	}
	var person models.Person
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, person)

}
