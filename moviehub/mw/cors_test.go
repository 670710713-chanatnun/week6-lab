package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflightAllowsConfiguredOrigin(t *testing.T) {
	handler := CORS([]string{"http://localhost:3000"}, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("preflight request should not reach the next handler")
	}))
	req := httptest.NewRequest(http.MethodOptions, "/api/movies", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:3000")
	}
}

func TestCORSDoesNotAllowUnconfiguredOrigin(t *testing.T) {
	handler := CORS([]string{"http://localhost:3000"}, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/movies", nil)
	req.Header.Set("Origin", "https://untrusted.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}
