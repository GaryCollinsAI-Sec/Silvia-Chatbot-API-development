package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllowsRequestsWithinLimit(t *testing.T) {
	rateLimiter := NewRateLimiter(3, time.Minute)

	requestCount := 0

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiter.Middleware(next)

	for i := 0; i < 3; i++ {
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/chat",
			nil,
		)

		request.RemoteAddr = "192.0.2.1:12345"

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"request %d: expected status %d, got %d",
				i+1,
				http.StatusOK,
				recorder.Code,
			)
		}
	}

	if requestCount != 3 {
		t.Fatalf(
			"expected next handler to be called 3 times, got %d",
			requestCount,
		)
	}
}

func TestRateLimiterBlocksRequestsOverLimit(t *testing.T) {
	rateLimiter := NewRateLimiter(3, time.Minute)

	requestCount := 0

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiter.Middleware(next)

	for i := 0; i < 3; i++ {
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/chat",
			nil,
		)

		request.RemoteAddr = "192.0.2.2:12345"

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	request.RemoteAddr = "192.0.2.2:12345"

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusTooManyRequests,
			recorder.Code,
		)
	}

	if requestCount != 3 {
		t.Fatalf(
			"expected next handler to be called 3 times, got %d",
			requestCount,
		)
	}
}

func TestRateLimiterTracksClientsSeparately(t *testing.T) {
	rateLimiter := NewRateLimiter(2, time.Minute)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiter.Middleware(next)

	firstClient := "192.0.2.10:12345"
	secondClient := "192.0.2.11:12345"

	for i := 0; i < 2; i++ {
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/chat",
			nil,
		)

		request.RemoteAddr = firstClient

		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusOK {
			t.Fatalf(
				"first client request %d: expected status %d, got %d",
				i+1,
				http.StatusOK,
				recorder.Code,
			)
		}
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	request.RemoteAddr = secondClient

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"second client: expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	rateLimiter := NewRateLimiter(1, 50*time.Millisecond)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimiter.Middleware(next)

	client := "192.0.2.20:12345"

	firstRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	firstRequest.RemoteAddr = client

	firstRecorder := httptest.NewRecorder()

	handler.ServeHTTP(firstRecorder, firstRequest)

	if firstRecorder.Code != http.StatusOK {
		t.Fatalf(
			"first request: expected status %d, got %d",
			http.StatusOK,
			firstRecorder.Code,
		)
	}

	secondRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	secondRequest.RemoteAddr = client

	secondRecorder := httptest.NewRecorder()

	handler.ServeHTTP(secondRecorder, secondRequest)

	if secondRecorder.Code != http.StatusTooManyRequests {
		t.Fatalf(
			"second request: expected status %d, got %d",
			http.StatusTooManyRequests,
			secondRecorder.Code,
		)
	}

	time.Sleep(60 * time.Millisecond)

	thirdRequest := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	thirdRequest.RemoteAddr = client

	thirdRecorder := httptest.NewRecorder()

	handler.ServeHTTP(thirdRecorder, thirdRequest)

	if thirdRecorder.Code != http.StatusOK {
		t.Fatalf(
			"third request after window: expected status %d, got %d",
			http.StatusOK,
			thirdRecorder.Code,
		)
	}
}

func TestRateLimiterRejectsInvalidRemoteAddress(t *testing.T) {
	rateLimiter := NewRateLimiter(3, time.Minute)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := rateLimiter.Middleware(next)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/chat",
		nil,
	)

	request.RemoteAddr = "invalid-address"

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}
}
