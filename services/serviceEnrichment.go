package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetAge(name string) int {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	var result struct {
		Age int `json:"age"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Age
}

func GetGender(name string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var result struct {
		Gender string `json:"gender"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Gender
}
func GetNationality(name string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "null"
	}
	defer resp.Body.Close()
	var result struct {
		Country []struct {
			CountryId string `json:"country_id"`
		} `json:"country"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result.Country) > 0 {
		return result.Country[0].CountryId
	}
	return "null"

}
