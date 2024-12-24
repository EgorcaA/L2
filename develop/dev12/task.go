package main

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Event represents a calendar event.
type Event struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

var (
	events  = make(map[int]Event)
	mu      sync.Mutex
	eventID = 0
)

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Helper to parse form values
func parseFormValue(r *http.Request, key string) (string, error) {
	value := r.FormValue(key)
	if value == "" {
		return "", fmt.Errorf("missing parameter: %s", key)
	}
	return value, nil
}

func parseInt(value string, key string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer for %s", key)
	}
	return result, nil
}

func parseDate(value string, key string) (time.Time, error) {
	result, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date for %s", key)
	}
	return result, nil
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	userIDStr, err := parseFormValue(r, "user_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	userID, err := parseInt(userIDStr, "user_id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	title, err := parseFormValue(r, "title")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	description, err := parseFormValue(r, "description")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	dateStr, err := parseFormValue(r, "date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	date, err := parseDate(dateStr, "date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	eventID++
	events[eventID] = Event{
		ID:          eventID,
		UserID:      userID,
		Title:       title,
		Description: description,
		Date:        date,
	}

	writeJSON(w, http.StatusOK, map[string]string{"result": "event created"})
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	idStr, err := parseFormValue(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseInt(idStr, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	event, exists := events[id]
	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	if userIDStr := r.FormValue("user_id"); userIDStr != "" {
		userID, err := parseInt(userIDStr, "user_id")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		event.UserID = userID
	}

	if title := r.FormValue("title"); title != "" {
		event.Title = title
	}

	if description := r.FormValue("description"); description != "" {
		event.Description = description
	}

	if dateStr := r.FormValue("date"); dateStr != "" {
		date, err := parseDate(dateStr, "date")
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		event.Date = date
	}

	events[id] = event
	writeJSON(w, http.StatusOK, map[string]string{"result": "event updated"})
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid input"})
		return
	}

	idStr, err := parseFormValue(r, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := parseInt(idStr, "id")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := events[id]; !exists {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "event not found"})
		return
	}

	delete(events, id)
	writeJSON(w, http.StatusOK, map[string]string{"result": "event deleted"})
}

func eventsForDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := parseDate(dateStr, "date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()
	result := []Event{}
	for _, event := range events {
		if event.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": result})
}

func eventsForWeek(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	startDate, err := parseDate(dateStr, "date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	endDate := startDate.AddDate(0, 0, 7)

	mu.Lock()
	defer mu.Unlock()
	result := []Event{}
	for _, event := range events {
		if event.Date.After(startDate.Add(-time.Nanosecond)) && event.Date.Before(endDate) {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": result})
}

func eventsForMonth(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	startDate, err := parseDate(dateStr, "date")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	endDate := startDate.AddDate(0, 1, 0)

	mu.Lock()
	defer mu.Unlock()
	result := []Event{}
	for _, event := range events {
		if event.Date.After(startDate.Add(-time.Nanosecond)) && event.Date.Before(endDate) {
			result = append(result, event)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{"result": result})
}

func main() {
	http.HandleFunc("/create_event", createEvent)
	http.HandleFunc("/update_event", updateEvent)
	http.HandleFunc("/delete_event", deleteEvent)
	http.HandleFunc("/events_for_day", eventsForDay)
	http.HandleFunc("/events_for_week", eventsForWeek)
	http.HandleFunc("/events_for_month", eventsForMonth)

	handler := loggingMiddleware(http.DefaultServeMux)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
