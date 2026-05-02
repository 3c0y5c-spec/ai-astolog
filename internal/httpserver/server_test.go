package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealthz(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handleHealthz(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "{\"status\":\"ok\"}\n" {
		t.Fatalf("body = %q, want status ok JSON", response.Body.String())
	}
}
