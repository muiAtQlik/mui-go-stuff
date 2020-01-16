// https://www.thepolyglotdeveloper.com/2017/07/consume-restful-api-endpoints-golang-application/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetAllItems(t *testing.T) {

	type event struct {
		ID          string `json:"ID"`
		Title       string `json:"Title"`
		Description string `json:"Description"`
	}

	response, err := http.Get("http://localhost:8080/events")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %v but got %v", http.StatusOK, response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read the body %s\n", err)
		t.Fatal(err)
	}

	jsonStr := string(data)

	var eventsMaps []map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &eventsMaps)
	if err != nil {
		fmt.Println("Failed to unmarshal json string")
		t.Fatal(err)
	}

	// Verify the first default event
	id := eventsMaps[0]["ID"]
	title := eventsMaps[0]["Title"]
	description := eventsMaps[0]["Description"]

	if id != "1" {
		t.Errorf("Expected ID: %v but was: %v", "1", id)
	}

	if title != "Introduction to Golang" {
		t.Errorf("Expected Title: %v but was %v", "Introduction to Golang", title)
	}

	if description != "intro stuff" {
		t.Errorf("Expected Title: %v but was %v", "intro stuff", description)
	}

	// Verify the second default event
	id = eventsMaps[1]["ID"]
	title = eventsMaps[1]["Title"]
	description = eventsMaps[1]["Description"]

	if id != "2" {
		t.Errorf("Expected ID: %v but was: %v", "1", id)
	}

	if title != "A nice title" {
		t.Errorf("Expected Title: %v but was %v", "A nice title", title)
	}

	if description != "Some description goes here" {
		t.Errorf("Expected Title: %v but was %v", "Some description goes here", description)
	}
}

func TestAddAndDeleteEvent(t *testing.T) {
	// To be done
}
