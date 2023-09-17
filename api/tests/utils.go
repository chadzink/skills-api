package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/chadzink/skills-api/handlers"
)

func GetResponseBodyString(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult
		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// remove the created_at & updated_at fields from the response
		delete(responseObject.Data, "CreatedAt")
		delete(responseObject.Data, "UpdatedAt")
		delete(responseObject.Data, "DeletedAt")

		checkForSubKeys := []string{"skills", "categories", "person_skills"}
		for _, key := range checkForSubKeys {
			if _, ok := responseObject.Data[key]; ok {
				// Check fo rnil values
				if responseObject.Data[key] != nil {
					for i, _ := range responseObject.Data[key].([]interface{}) {
						delete(responseObject.Data[key].([]interface{})[i].(map[string]interface{}), "CreatedAt")
						delete(responseObject.Data[key].([]interface{})[i].(map[string]interface{}), "UpdatedAt")
						delete(responseObject.Data[key].([]interface{})[i].(map[string]interface{}), "DeletedAt")
					}
				}
			}
		}

		jsonString, err2 := json.MarshalIndent(responseObject, "", "	")
		if err2 != nil {
			t.Error(err2)
		}
		return string(jsonString)
	}
	return ""
}

func GetResponsesBodyString(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObjects handlers.ResponseResults
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for d, _ := range responseObjects.Data {
			// remove the created_at & updated_at fields from the response
			delete(responseObjects.Data[d], "CreatedAt")
			delete(responseObjects.Data[d], "UpdatedAt")
			delete(responseObjects.Data[d], "DeletedAt")

			checkForSubKeys := []string{"skills", "categories", "person_skills"}
			for _, key := range checkForSubKeys {
				if _, ok := responseObjects.Data[d][key]; ok {
					// Check fo rnil values
					if responseObjects.Data[d][key] != nil {
						for i, _ := range responseObjects.Data[d][key].([]interface{}) {
							delete(responseObjects.Data[d][key].([]interface{})[i].(map[string]interface{}), "CreatedAt")
							delete(responseObjects.Data[d][key].([]interface{})[i].(map[string]interface{}), "UpdatedAt")
							delete(responseObjects.Data[d][key].([]interface{})[i].(map[string]interface{}), "DeletedAt")
						}
					}
				}
			}
		}

		jsonString, err2 := json.MarshalIndent(responseObjects, "", "	")
		if err2 != nil {
			t.Error(err2)
		}
		return string(jsonString)
	}
	return ""
}
