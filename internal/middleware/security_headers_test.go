package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{
			name:     "X-Content-Type-Options",
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
		{
			name:     "X-Frame-Options",
			header:   "X-Frame-Options",
			expected: "DENY",
		},
		{
			name:     "Referrer-Policy",
			header:   "Referrer-Policy",
			expected: "no-referrer",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := recorder.Header().Get(test.header)

			if actual != test.expected {
				t.Fatalf(
					"expected %s to be %q, got %q",
					test.header,
					test.expected,
					actual,
				)
			}
		})
	}
}
