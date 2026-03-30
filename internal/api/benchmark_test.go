package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/openclaw/agent-proxy/internal/pipeline"
	"github.com/openclaw/agent-proxy/internal/routing"
	"github.com/openclaw/agent-proxy/internal/token"
)

func BenchmarkChatCompletions(b *testing.B) {
	forwardingStage := pipeline.NewForwardingStage()
	tokenStore := &token.TokenStore{}
	tokenResolver := token.NewTokenResolver(tokenStore, nil)
	routingHandler := routing.NewRequestHandler(tokenResolver, 3, 100, 5000, routing.StrategyRoundRobin)
	handler := NewChatCompletionsHandler(forwardingStage, tokenResolver, routingHandler)

	body := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}
	bodyBytes, _ := json.Marshal(body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/v1/chat/completions", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkEmbeddings(b *testing.B) {
	forwardingStage := pipeline.NewForwardingStage()
	tokenStore := &token.TokenStore{}
	tokenResolver := token.NewTokenResolver(tokenStore, nil)
	routingHandler := routing.NewRequestHandler(tokenResolver, 3, 100, 5000, routing.StrategyRoundRobin)
	handler := NewEmbeddingsHandler(forwardingStage, tokenResolver, routingHandler)

	body := map[string]interface{}{
		"model": "text-embedding-ada-002",
		"input": "Hello world",
	}
	bodyBytes, _ := json.Marshal(body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/v1/embeddings", bytes.NewReader(bodyBytes))
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkRequestParsing(b *testing.B) {
	body := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}
	bodyBytes, _ := json.Marshal(body)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var req ChatCompletionRequest
		json.Unmarshal(bodyBytes, &req)
	}
}

func BenchmarkResponseWriting(b *testing.B) {
	response := ChatCompletionResponse{
		ID:      "test",
		Object:  "chat.completion",
		Created: 1234567890,
		Model:   "gpt-4",
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: "Hello!",
				},
				FinishReason: "stop",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(response)
	}
}
