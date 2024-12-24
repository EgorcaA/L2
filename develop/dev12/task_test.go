package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func TestCreateEvent(t *testing.T) {
	server := httptest.NewServer(loggingMiddleware(http.HandlerFunc(createEvent)))
	defer server.Close()

	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Test Event")
	data.Set("description", "Description")
	data.Set("date", "2024-12-25")

	resp, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Result != "event created" {
		t.Errorf("Unexpected result: %v", res.Result)
	}
}

func TestUpdateEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(createEvent))
	defer server.Close()

	// First, create an event
	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Original Event")
	data.Set("description", "Original Description")
	data.Set("date", "2024-12-25")
	_, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Now, update it
	updateServer := httptest.NewServer(http.HandlerFunc(updateEvent))
	defer updateServer.Close()

	updateData := url.Values{}
	updateData.Set("id", "1")
	updateData.Set("title", "Updated Event")

	req, err := http.NewRequest(http.MethodPost, updateServer.URL, strings.NewReader(updateData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Result != "event updated" {
		t.Errorf("Unexpected result: %v", res.Result)
	}
}

func TestDeleteEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(createEvent))
	defer server.Close()

	// Create an event
	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Event to delete")
	data.Set("description", "Description")
	data.Set("date", "2024-12-25")
	_, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Delete the event
	deleteServer := httptest.NewServer(http.HandlerFunc(deleteEvent))
	defer deleteServer.Close()

	deleteData := url.Values{}
	deleteData.Set("id", "1")

	req, err := http.NewRequest(http.MethodPost, deleteServer.URL, strings.NewReader(deleteData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if res.Result != "event deleted" {
		t.Errorf("Unexpected result: %v", res.Result)
	}
}

func TestEventsForDay(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(createEvent))
	defer server.Close()

	// Create an event
	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Day Event")
	data.Set("description", "Day Description")
	data.Set("date", "2024-12-25")
	_, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Fetch events for the day
	dayServer := httptest.NewServer(http.HandlerFunc(eventsForDay))
	defer dayServer.Close()

	resp, err := http.Get(dayServer.URL + "?date=2024-12-25")
	if err != nil {
		t.Fatalf("Failed to fetch events for the day: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(res["result"].([]interface{})) != 1 {
		t.Errorf("Expected 1 event, got %v", len(res["result"].([]interface{})))
	}
}

func TestEventsForWeek(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(createEvent))
	defer server.Close()

	// Create an event
	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Week Event")
	data.Set("description", "Week Description")
	data.Set("date", "2024-12-26")
	_, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Fetch events for the week
	weekServer := httptest.NewServer(http.HandlerFunc(eventsForWeek))
	defer weekServer.Close()

	resp, err := http.Get(weekServer.URL + "?date=2024-12-25")
	if err != nil {
		t.Fatalf("Failed to fetch events for the week: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(res["result"].([]interface{})) != 1 {
		t.Errorf("Expected 1 event, got %v", len(res["result"].([]interface{})))
	}
}

func TestEventsForMonth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(createEvent))
	defer server.Close()

	// Create an event
	data := url.Values{}
	data.Set("user_id", "1")
	data.Set("title", "Month Event")
	data.Set("description", "Month Description")
	data.Set("date", "2024-12-31")
	_, err := http.PostForm(server.URL, data)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	// Fetch events for the month
	monthServer := httptest.NewServer(http.HandlerFunc(eventsForMonth))
	defer monthServer.Close()

	resp, err := http.Get(monthServer.URL + "?date=2024-12-01")
	if err != nil {
		t.Fatalf("Failed to fetch events for the month: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(res["result"].([]interface{})) != 1 {
		t.Errorf("Expected 1 event, got %v", len(res["result"].([]interface{})))
	}
}
