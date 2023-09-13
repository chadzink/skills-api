package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/chadzink/skills-api/handlers"
	"github.com/chadzink/skills-api/models"
)

func MapToSkill(m map[string]interface{}) models.Skill {
	return models.Skill{
		Name:        m["name"].(string),
		Description: m["description"].(string),
		ShortKey:    m["short_key"].(string),
		Active:      m["active"].(bool),
	}
}

func MapToCategory(m map[string]interface{}) models.Category {
	return models.Category{
		Name:        m["name"].(string),
		Description: m["description"].(string),
		ShortKey:    m["short_key"].(string),
		Active:      m["active"].(bool),
	}
}

func MapToPerson(m map[string]interface{}) models.Person {
	return models.Person{
		Name:    m["name"].(string),
		Email:   m["email"].(string),
		Phone:   m["phone"].(string),
		Profile: m["profile"].(string),
	}
}

func parseResponseData(t *testing.T, resp *http.Response, i interface{}) interface{} {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult

		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// Use a case statement to determine the type of object to return
		switch i.(type) {
		case models.Skill:
			responseDataAsSkill := MapToSkill(responseObject.Data)
			responseDataAsSkill.ID = uint(responseObject.Data["ID"].(float64))
			return responseDataAsSkill
		case models.Category:
			responseDataAsCategory := MapToCategory(responseObject.Data)
			responseDataAsCategory.ID = uint(responseObject.Data["ID"].(float64))
			return responseDataAsCategory
		case models.Person:
			responseDataAsPerson := MapToPerson(responseObject.Data)
			responseDataAsPerson.ID = uint(responseObject.Data["ID"].(float64))
			return responseDataAsPerson
		default:
			return nil
		}
	}

	return models.Skill{}
}
