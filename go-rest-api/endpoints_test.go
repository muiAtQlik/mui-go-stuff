// Unit test for events

package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"
)

func removeUnprintableChars(str string) string {
	return strings.TrimFunc(str, func(r rune) bool {
		return !unicode.IsGraphic(r)
	})
}

// TestGetAllEvents test to get all events
func TestGetAllEvents(t *testing.T) {
	req, err := http.NewRequest("GET", "/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllEvents)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"ID":"1","Title":"Introduction to Golang","Description":"intro stuff"},{"ID":"2","Title":"A nice title","Description":"Some description goes here"}]`
	//trim unprintable chars since the body of the response contanis a new line (just for testing stuff)
	body := removeUnprintableChars(rr.Body.String())

	if body != expected {
		t.Errorf("expected body:\n%v\nbut got:\n%v", expected, body)
	}
}

func TestGetOneEvent(t *testing.T) {
	req, err := http.NewRequest("GET", "/event/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getOneEvent)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"ID":"1","Title":"Introduction to Golang","Description":"intro stuff"}`
	body := removeUnprintableChars(rr.Body.String())

	if body != expected {
		t.Errorf("expected body:\n%v\nbut got:\n%v", expected, body)
	}
}

func TestCreateEvent(t *testing.T) {

	jsonStr := `{"ID":"3","Title":"666","Description":"3:rd event is evil"}`
	var jsonByteStr = []byte(jsonStr)
	req, err := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonByteStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(createEvent)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	actual := removeUnprintableChars(rr.Body.String())
	if actual != jsonStr {
		t.Errorf("expected body:\n%v\nbut got:\n%v", jsonStr, rr.Body.String())
	}
}
