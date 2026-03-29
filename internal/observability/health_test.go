package observability

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthChecker_HandleHealth(t *testing.T) {
	checker := NewHealthChecker(nil)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	checker.HandleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}
}

func TestHealthChecker_HandleReady(t *testing.T) {
	checker := NewHealthChecker(nil)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()

	checker.HandleReady(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ready" {
		t.Errorf("Expected status 'ready', got '%v'", response["status"])
	}
}

func TestHealthChecker_HandleLive(t *testing.T) {
	checker := NewHealthChecker(nil)

	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()

	checker.HandleLive(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "alive" {
		t.Errorf("Expected status 'alive', got '%s'", response["status"])
	}
}

func TestHealthChecker_MethodNotAllowed(t *testing.T) {
	checker := NewHealthChecker(nil)

	methods := []string{"POST", "PUT", "DELETE", "PATCH"}
	endpoints := []string{"/health", "/health/ready", "/health/live"}

	for _, method := range methods {
		for _, endpoint := range endpoints {
			req := httptest.NewRequest(method, endpoint, nil)
			w := httptest.NewRecorder()

			switch endpoint {
			case "/health":
				checker.HandleHealth(w, req)
			case "/health/ready":
				checker.HandleReady(w, req)
			case "/health/live":
				checker.HandleLive(w, req)
			}

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 405 for %s %s, got %d", method, endpoint, w.Code)
			}
		}
	}
}
