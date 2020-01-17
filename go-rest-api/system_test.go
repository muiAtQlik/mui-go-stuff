// https://www.thepolyglotdeveloper.com/2017/07/consume-restful-api-endpoints-golang-application/

package main

import (
	"bytes"
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

	// Create an event
	jsonData := map[string]string{"ID": "666", "Title": "Evil Number", "Description": "NUmber of the beast"}
	jsonDataBytes, _ := json.Marshal(jsonData)
	response, err := http.Post("http://localhost:8080/event", "application/json", bytes.NewBuffer(jsonDataBytes))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code: %v but got %v", http.StatusCreated, response.Status)
	}

	data, _ := ioutil.ReadAll(response.Body)
	expected := `{"ID":"666","Title":"Evil Number","Description":"NUmber of the beast"}`

	body := removeUnprintableChars(string(data))

	if body != expected {
		t.Errorf("Expected body: %v but got %v", expected, body)
	}

	// fetch all event and verify
	response, err = http.Get("http://localhost:8080/events")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %v but got %v", http.StatusOK, response.Status)
	}

	data, _ = ioutil.ReadAll(response.Body)
	body = removeUnprintableChars(string(data))

	expected = `[{"ID":"1","Title":"Introduction to Golang","Description":"intro stuff"},{"ID":"2","Title":"A nice title","Description":"Some description goes here"},{"ID":"666","Title":"Evil Number","Description":"NUmber of the beast"}]`

	if body != expected {
		t.Errorf("Expected\n%v\nbut got\n%v", expected, body)
	}

	// delete an event
	client := &http.Client{}

	request, error := http.NewRequest("DELETE", "http://localhost:8080/event/666", nil)

	if error != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		t.Fatal(err)
	}

	response, err = client.Do(request)

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %v but got %v", http.StatusOK, response.Status)
	}

	// fetch all event and verify that the new event is deleted

	response, err = http.Get("http://localhost:8080/events")

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code: %v but got %v", http.StatusOK, response.Status)
	}

	data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read the body %s\n", err)
		t.Fatal(err)
	}

	actual := removeUnprintableChars(string(data))
	expected = `[{"ID":"1","Title":"Introduction to Golang","Description":"intro stuff"},{"ID":"2","Title":"A nice title","Description":"Some description goes here"}]`

	if expected != actual {
		t.Errorf("Expected:\n%v\nbut got:\n%v", expected, actual)
	}
}
