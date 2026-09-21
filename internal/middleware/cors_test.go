package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSAllowsApprovedOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(next)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	request.Header.Set("Origin", "http://localhost:5173")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Origin"); actual != "http://localhost:5173" {
		t.Fatalf(
			"expected Access-Control-Allow-Origin to be %q, got %q",
			"http://localhost:5173",
			actual,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Methods"); actual != "POST, OPTIONS" {
		t.Fatalf(
			"expected Access-Control-Allow-Methods to be %q, got %q",
			"POST, OPTIONS",
			actual,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Headers"); actual != "Content-Type" {
		t.Fatalf(
			"expected Access-Control-Allow-Headers to be %q, got %q",
			"Content-Type",
			actual,
		)
	}

	if actual := recorder.Header().Get("Vary"); actual != "Origin" {
		t.Fatalf(
			"expected Vary to be %q, got %q",
			"Origin",
			actual,
		)
	}
}

func TestCORSRejectsUnapprovedPreflightOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for rejected preflight")
	})

	handler := CORS(next)

	request := httptest.NewRequest(
		http.MethodOptions,
		"/api/chat",
		nil,
	)

	request.Header.Set("Origin", "http://malicious.example")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusForbidden,
			recorder.Code,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Origin"); actual != "" {
		t.Fatalf(
			"expected no Access-Control-Allow-Origin header, got %q",
			actual,
		)
	}
}

func TestCORSAllowsApprovedPreflight(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called for preflight")
	})

	handler := CORS(next)

	request := httptest.NewRequest(
		http.MethodOptions,
		"/api/chat",
		nil,
	)

	request.Header.Set("Origin", "http://localhost:5173")
	request.Header.Set("Access-Control-Request-Method", "POST")
	request.Header.Set("Access-Control-Request-Headers", "Content-Type")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNoContent,
			recorder.Code,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Origin"); actual != "http://localhost:5173" {
		t.Fatalf(
			"expected Access-Control-Allow-Origin to be %q, got %q",
			"http://localhost:5173",
			actual,
		)
	}
}

func TestCORSDoesNotGrantHeadersToUnapprovedOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORS(next)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	request.Header.Set("Origin", "http://malicious.example")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}

	if actual := recorder.Header().Get("Access-Control-Allow-Origin"); actual != "" {
		t.Fatalf(
			"expected no Access-Control-Allow-Origin header, got %q",
			actual,
		)
	}
}
