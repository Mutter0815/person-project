package handlers

import (
	"errors"
	"net/http"
	"person-project/db"
	"person-project/dto"
	"person-project/logger"
	"person-project/models"
	"person-project/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// CreatePerson godoc
// @Summary Создание нового человека
// @Description Создаёт нового человека с обогащением данных
// @Tags person
// @Accept json
// @Produce json
// @Param person body dto.CreatePersonRequest true "Данные нового человека (имя и фамилия обязательны)"
// @Success 201 {object} models.Person "Пользователь успешно создан"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people [post]
func CreatePerson(c *gin.Context) {

	var personRequest dto.CreatePersonRequest
	err := c.ShouldBindJSON(&personRequest)
	if err != nil {
		logger.Log.Error("Ошибка парсинга JSON",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Debug("Создание пользователя",
		"name", personRequest.Name,
		"surname", personRequest.Surname)
	var person models.Person
	person.Name = personRequest.Name
	person.Surname = personRequest.Surname
	person.Patronymic = personRequest.Patronymic
	person.Age = services.GetAge(person.Name)
	person.Gender = services.GetGender(person.Name)
	person.Nationality = services.GetNationality(person.Name)
	logger.Log.Debug("Сохранение пользователя в бд",
		"name", person.Name,
		"id", person.ID)
	if err := db.DB.Create(&person).Error; err != nil {
		logger.Log.Error("Ошибка сохранения в БД", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения данных" + err.Error()})
	}
	logger.Log.Info("Пользователь успешно создан", "id", person.ID)
	c.JSON(http.StatusCreated, person)

}

// DeletePerson godoc
// @Summary Удаление пользователя
// @Description Удаляет пользователя по id
// @Tags person
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.Person "Успешное удаление пользователя"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /people/{id} [delete]
func DeletePerson(c *gin.Context) {
	var person models.Person
	strId := c.Param("id")
	logger.Log.Debug("Попытка удаления пользователя",
		"id", strId)
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Log.Error("Ошибка,недействительный ID",
			"id", id,
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный ID"})
		return
	}

	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Warn("Пользователь не найден",
				"id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Debug("Удаления пользователя",
		"id", person.ID,
		"name", person.Name)

	if err := db.DB.Delete(&person).Error; err != nil {
		logger.Log.Error("Ошибка при удалении пользователя",
			"error", err,
			"id", person.ID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Пользователь успешно удалён")
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удалён"})

}

// UpdatePerson godoc
// @Summary Обновление данных пользователя
// @Description Обновляет данные пользователя по id
// @Tags person
// @Accept json
// @Param id path string true "ID пользователя"
// @Param input body dto.UpdatePersonRequest false "Данные для обновления"
// @Produce json
// @Success 200 {object} models.Person "Обновленные данные пользователя"
// @Failrue 404 {object} dto.ErrorResponse
// @Failrue 500 {object} dto.ErrorResponse
// @Router /people/{id} [patch]
func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	logger.Log.Debug("Попытка обновления данных пользователя", "id", id)
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("Пользователь не найден", "id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var updateDataPerson dto.UpdatePersonRequest

	err := c.ShouldBindJSON(&updateDataPerson)
	if err != nil {
		logger.Log.Error("Ошибка парсинга обновленных данных JSON",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Debug("Валидация данных",
		"person", person)
	validate := validator.New()
	if err := validate.Struct(updateDataPerson); err != nil {
		logger.Log.Error("Ошибка валидации",
			"error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if updateDataPerson.Name != nil {
		person.Name = *updateDataPerson.Name

	}
	if updateDataPerson.Age != nil {
		person.Age = *updateDataPerson.Age

	}
	if updateDataPerson.Gender != nil {
		person.Gender = models.Gender(*updateDataPerson.Gender)
	}
	if updateDataPerson.Nationality != nil {
		person.Nationality = *updateDataPerson.Nationality
	}
	if updateDataPerson.Surname != nil {
		person.Surname = *updateDataPerson.Surname
	}
	if updateDataPerson.Patronymic != nil && *updateDataPerson.Patronymic != "" {
		person.Patronymic = updateDataPerson.Patronymic
	}
	if err := db.DB.Save(&person).Error; err != nil {
		logger.Log.Error("Ошибка при обновлении данных пользователя",
			"name", person.Name,
			"id", person.ID,
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}
	logger.Log.Info("Данные пользователя успешно обновлены",
		"id", person.ID,
		"name", person.Name)
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
// @Success 200  {array} models.Person "Cписок пользователей"
// @Failrue 500 {object} dto.ErrorResponse "Ошибка сервера"
// @Router /people [get]
func GetPersons(c *gin.Context) {

	var persons []models.Person
	db := db.DB.Model(&models.Person{})
	logger.Log.Debug("Попытка вывода пользователей по фильтрам")
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
	logger.Log.Debug("Попытка выполнения запроса в БД")
	if err := db.Limit(limit).Offset(offset).Find(&persons).Error; err != nil {
		logger.Log.Error("Ошибка при попытке выполнения запроса",
			"error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей", "details": err.Error()})
		return
	}
	logger.Log.Info("Выполнение запроса прошло успешно",
		"persons", persons)
	c.JSON(http.StatusOK, persons)
}

// GetPersonByID godoc
// @Summary Возращает данные о пользователе
// @Description Возращает данные пользователе по id
// @Tags person
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} models.Person "Данные о пользователе"
// @Failrue 404 dto.ErrorResponse "Пользователь не найден"
// @Router /people/{id} [get]
func GetPersonByID(c *gin.Context) {
	logger.Log.Debug("Попытка вывода Пользователя по id")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.Error("Неверный формат id",
			"id", id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Требуется id"})
	}
	var person models.Person
	logger.Log.Debug("Поиск пользователя в бд",
		"id", id)
	if err := db.DB.First(&person, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Log.Error("Пользователь не найден",
				"id", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Пользователь успешно найден",
		"id", id,
		"name", person.Name)
	c.JSON(http.StatusOK, person)

}
