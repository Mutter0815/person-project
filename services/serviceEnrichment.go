package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"person-project/config"
	"person-project/logger"
	"person-project/models"
)

func GetAge(name string) int {
	logger.Log.Debug("Запрос возраста в Agify API", "name", name)
	resp, err := http.Get(config.Cfg.API.Agify_API + fmt.Sprintf("/?name=%s", name))
	if err != nil {
		logger.Log.Error("Ошибка Agify API", "error", err)
		return 0
	}
	defer resp.Body.Close()
	var result struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log.Error("Ошибка парсинга ответа", "error", err)
		return 0
	}
	logger.Log.Info("Получен возраст", "name ", name, "age ", result.Age)
	return result.Age
}

func GetGender(name string) models.Gender {
	logger.Log.Debug("Запрос Пола в Genderize API", "name", name)
	resp, err := http.Get(config.Cfg.API.Genderize_API + fmt.Sprintf("/?name=%s", name))
	if err != nil {
		logger.Log.Error("Ошибка Genderize API", "error", err)
		return ""
	}
	defer resp.Body.Close()
	var result struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log.Error("Ошибка парсига ответа",
			"error", err)
		return ""
	}
	logger.Log.Info("Получен Пол пользователя", "name", name)
	return models.Gender(result.Gender)

}

func GetNationality(name string) string {
	logger.Log.Debug("Запрос Национальности в Nationality API",
		"name", name)
	resp, err := http.Get(config.Cfg.API.Nationalize_API + fmt.Sprintf("/?name=%s", name))
	if err != nil {
		return "null"
	}
	defer resp.Body.Close()
	var result struct {
		Count   int `json:"count"`
		Country []struct {
			CountryId   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Log.Error("Ошибка парсинга ответа",
			"error", err)
		return ""
	}

	if result.Count == 0 {
		logger.Log.Warn("Не удалось определить национальность", "name", name)
		return ""
	}
	max := result.Country[0]
	for _, country := range result.Country {
		if country.Probability > max.Probability {
			max = country
		}

	}
	logger.Log.Info("Получена максимально вероятная национальность пользователя",
		"name", name,
		"Country", max.CountryId)
	return max.CountryId

}
