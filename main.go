
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	// Limit request body size.
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
		return
	}

	// Reject duplicate JSON object keys before decoding
	// the request into the application model.
	hasDuplicates, err := hasDuplicateJSONKeys(body)
	if err != nil || hasDuplicates {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
		return
	}

	var request chatbot.ChatRequest

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&request)
	if err != nil {
		writeJSONError(
			w,
			http.StatusBadRequest,
			"invalid_request",
		)
		return
	}

	// Require exactly one JSON value in the request body.
	var extra interface{}

	err = decoder.Decode(&extra)
	if err != io.EOF {
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

func hasDuplicateJSONKeys(data []byte) (bool, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))

	hasDuplicates, err := scanJSONValue(decoder)
	if err != nil {
		return false, err
	}

	return hasDuplicates, nil
}

func scanJSONValue(decoder *json.Decoder) (bool, error) {
	token, err := decoder.Token()
	if err != nil {
		return false, err
	}

	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			seenKeys := make(map[string]struct{})

			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return false, err
				}

				key, ok := keyToken.(string)
				if !ok {
					return false, fmt.Errorf("invalid JSON object key")
				}

				if _, exists := seenKeys[key]; exists {
					return true, nil
				}

				seenKeys[key] = struct{}{}

				hasDuplicates, err := scanJSONValue(decoder)
				if err != nil {
					return false, err
				}

				if hasDuplicates {
					return true, nil
				}
			}

			_, err := decoder.Token()
			if err != nil {
				return false, err
			}

		case '[':
			for decoder.More() {
				hasDuplicates, err := scanJSONValue(decoder)
				if err != nil {
					return false, err
				}

				if hasDuplicates {
					return true, nil
				}
			}

			_, err := decoder.Token()
			if err != nil {
				return false, err
			}
		}
	}

	return false, nil
}

