package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/cors"
)

func TestCORS(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	cfg := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	CORS(cfg).Then(mux).ServeHTTP(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	vary := rec.Header().Get("Vary")
	if origin != "*" {
		t.Errorf("want %s, got %s", "*", origin)
	}

	if vary != "Origin" {
		t.Errorf("want %s, got %s", "Origin", vary)
	}
}
