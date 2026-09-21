package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/GaryCollinsAI-Sec/silver-dragons-chatbot-api/internal/chatbot"
	"github.com/GaryCollinsAI-Sec/silver-dragons-chatbot-api/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type APIError struct {
	Error string `json:"error"`
}

func main() {
	router := chi.NewRouter()

	rateLimiter := middleware.NewRateLimiter(20, time.Minute)

	// Global security middleware
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.CORS)

	// API routes
	router.Get("/api/health", healthHandler)
	router.With(rateLimiter.Middleware).Post("/api/chat", chatHandler)

	fmt.Println("Silvia API running on http://localhost:8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}

func writeJSONError(w http.ResponseWriter, status int, errorCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := APIError{
		Error: errorCode,
	}

	json.NewEncoder(w).Encode(response)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		writeJSONError(
			w,
			http.StatusMethodNotAllowed,
			"method_not_allowed",
		)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(
			w,
			http.StatusUnsupportedMediaType,
			"unsupported_media_type",
		)
		return
	}

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	var request chatbot.ChatRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)
	if err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
		return
	}

	message := strings.TrimSpace(request.Message)

	if message == "" {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"empty_message",
		)
		return
	}

	if len(message) > 1000 {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"message_too_long",
		)
		return
	}

	reply := chatbot.Answer(message)

	response := chatbot.ChatResponse{
		Assistant: "Silvia",
		Reply:     reply,
	}

	json.NewEncoder(w).Encode(response)
}
