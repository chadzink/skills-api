package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"skills-api/handlers"
)

// Gets the response body as a string and removes the created_at & updated_at fields from the response
//
//	This is important becaus ethe created_at & updated_at fields cannot be compared to the golden files if they always change
func GetResponseBodyString(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObject handlers.ResponseResult[map[string]interface{}]
		if err := json.Unmarshal(resBodyBytes, &responseObject); err != nil {
			t.Error(err)
		}

		// remove the created_at & updated_at fields from the response
		responseObject.Data = RemoveDatesFromResponseData(responseObject.Data)

		jsonString, err2 := json.MarshalIndent(responseObject, "", "	")
		if err2 != nil {
			t.Error(err2)
		}
		return string(jsonString)
	}
	return ""
}

// Gets the response body for multipel records as a string and removes the created_at & updated_at fields from the response
//
//	This is important becaus ethe created_at & updated_at fields cannot be compared to the golden files if they always change
func GetResponsesBodyString(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()
	if resBodyBytes, err := io.ReadAll(resp.Body); err != nil {
		t.Error(err)
	} else {
		// parse the response body
		var responseObjects handlers.ResponseResults[map[string]interface{}]
		if err := json.Unmarshal(resBodyBytes, &responseObjects); err != nil {
			t.Error(err)
		}

		for d, _ := range responseObjects.Data {
			// remove the created_at & updated_at fields from the response
			responseObjects.Data[d] = RemoveDatesFromResponseData(responseObjects.Data[d])
		}

		jsonString, err2 := json.MarshalIndent(responseObjects, "", "	")
		if err2 != nil {
			t.Error(err2)
		}
		return string(jsonString)
	}
	return ""
}

// Removes the created_at & updated_at fields from the response map[string]interface{} object
func RemoveDatesFromResponseData(data map[string]interface{}) map[string]interface{} {
	// remove the created_at & updated_at fields from the response
	delete(data, "CreatedAt")
	delete(data, "UpdatedAt")
	delete(data, "DeletedAt")

	checkForSubKeys := []string{"skills", "categories", "person_skills"}
	for _, key := range checkForSubKeys {
		if _, ok := data[key]; ok {
			// Check for nil values
			if data[key] != nil {
				for i, _ := range data[key].([]interface{}) {
					delete(data[key].([]interface{})[i].(map[string]interface{}), "CreatedAt")
					delete(data[key].([]interface{})[i].(map[string]interface{}), "UpdatedAt")
					delete(data[key].([]interface{})[i].(map[string]interface{}), "DeletedAt")

					checkForSubSubKeys := []string{"skill", "expertise", "person", "category"}
					for _, subKey := range checkForSubSubKeys {
						if _, ok := data[key].([]interface{})[i].(map[string]interface{})[subKey]; ok {
							// Check for nil values
							if data[key].([]interface{})[i].(map[string]interface{})[subKey] != nil {
								delete(data[key].([]interface{})[i].(map[string]interface{})[subKey].(map[string]interface{}), "CreatedAt")
								delete(data[key].([]interface{})[i].(map[string]interface{})[subKey].(map[string]interface{}), "UpdatedAt")
								delete(data[key].([]interface{})[i].(map[string]interface{})[subKey].(map[string]interface{}), "DeletedAt")
							}
						}
					}
				}
			}
		}
	}

	return data
}
