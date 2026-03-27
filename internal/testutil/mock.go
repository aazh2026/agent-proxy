package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type MockOpenAIServer struct {
	Server *httptest.Server
}

func NewMockOpenAIServer() *MockOpenAIServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      "chatcmpl-test",
			"object":  "chat.completion",
			"created": 1234567890,
			"model":   "gpt-4",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hello from mock OpenAI!",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]interface{}{
				"prompt_tokens":     10,
				"completion_tokens": 20,
				"total_tokens":      30,
			},
		})
	})

	mux.HandleFunc("/v1/embeddings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"object": "list",
			"data": []map[string]interface{}{
				{
					"object":    "embedding",
					"embedding": []float64{0.1, 0.2, 0.3, 0.4, 0.5},
					"index":     0,
				},
			},
			"model": "text-embedding-ada-002",
			"usage": map[string]interface{}{
				"prompt_tokens": 5,
				"total_tokens":  5,
			},
		})
	})

	return &MockOpenAIServer{
		Server: httptest.NewServer(mux),
	}
}

func (m *MockOpenAIServer) Close() {
	m.Server.Close()
}

func (m *MockOpenAIServer) URL() string {
	return m.Server.URL
}

type MockAnthropicServer struct {
	Server *httptest.Server
}

func NewMockAnthropicServer() *MockAnthropicServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   "msg_test",
			"type": "message",
			"role": "assistant",
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "Hello from mock Anthropic!",
				},
			},
			"model":       "claude-3-opus-20240229",
			"stop_reason": "end_turn",
			"usage": map[string]interface{}{
				"input_tokens":  10,
				"output_tokens": 20,
			},
		})
	})

	return &MockAnthropicServer{
		Server: httptest.NewServer(mux),
	}
}

func (m *MockAnthropicServer) Close() {
	m.Server.Close()
}

func (m *MockAnthropicServer) URL() string {
	return m.Server.URL
}

type MockGoogleServer struct {
	Server *httptest.Server
}

func NewMockGoogleServer() *MockGoogleServer {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"candidates": []map[string]interface{}{
				{
					"content": map[string]interface{}{
						"parts": []map[string]interface{}{
							{"text": "Hello from mock Google!"},
						},
						"role": "model",
					},
					"finishReason": "STOP",
				},
			},
			"usageMetadata": map[string]interface{}{
				"promptTokenCount":     10,
				"candidatesTokenCount": 20,
				"totalTokenCount":      30,
			},
		})
	})

	return &MockGoogleServer{
		Server: httptest.NewServer(mux),
	}
}

func (m *MockGoogleServer) Close() {
	m.Server.Close()
}

func (m *MockGoogleServer) URL() string {
	return m.Server.URL
}

func AssertNoError(t interface{ Fatalf(string, ...interface{}) }, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func AssertEqual(t interface{ Errorf(string, ...interface{}) }, expected, actual interface{}, msg string) {
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

func FormatJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func CreateTestRequest(method, path string, body interface{}) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, path, nil)
	}

	_, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	return http.NewRequest(method, path, nil)
}
