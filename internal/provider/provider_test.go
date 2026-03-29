package provider

import (
	"context"
	"io"
	"testing"
)

func TestRegistry_ResolveProvider_GPT(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "openai"}
	registry.Register(mockProvider)

	provider, name := registry.ResolveProvider("gpt-3.5-turbo")
	if provider == nil {
		t.Error("Expected provider for gpt-3.5-turbo, got nil")
	}
	if name != "openai" {
		t.Errorf("Expected provider name 'openai', got '%s'", name)
	}
}

func TestRegistry_ResolveProvider_Claude(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "anthropic"}
	registry.Register(mockProvider)

	provider, name := registry.ResolveProvider("claude-3-opus")
	if provider == nil {
		t.Error("Expected provider for claude-3-opus, got nil")
	}
	if name != "anthropic" {
		t.Errorf("Expected provider name 'anthropic', got '%s'", name)
	}
}

func TestRegistry_ResolveProvider_Gemini(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "google"}
	registry.Register(mockProvider)

	provider, name := registry.ResolveProvider("gemini-pro")
	if provider == nil {
		t.Error("Expected provider for gemini-pro, got nil")
	}
	if name != "google" {
		t.Errorf("Expected provider name 'google', got '%s'", name)
	}
}

func TestRegistry_ResolveProvider_Unknown(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "openai"}
	registry.Register(mockProvider)

	provider, name := registry.ResolveProvider("unknown-model")
	if provider != nil {
		t.Error("Expected nil provider for unknown model")
	}
	if name != "" {
		t.Errorf("Expected empty provider name for unknown model, got '%s'", name)
	}
}

func TestRegistry_ResolveProvider_CustomAlias(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "openai"}
	registry.Register(mockProvider)

	provider, name := registry.ResolveProvider("gpt-4")
	if provider == nil {
		t.Error("Expected provider for gpt-4, got nil")
	}
	if name != "openai" {
		t.Errorf("Expected provider name 'openai', got '%s'", name)
	}
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	mockProvider := &MockProvider{name: "test-provider"}
	registry.Register(mockProvider)

	provider, ok := registry.Get("test-provider")
	if !ok {
		t.Error("Expected to find registered provider")
	}
	if provider == nil {
		t.Error("Expected non-nil provider")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&MockProvider{name: "openai"})
	registry.Register(&MockProvider{name: "anthropic"})

	names := registry.List()
	if len(names) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(names))
	}
}

type MockProvider struct {
	name string
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) TransformRequest(req *ChatRequest) (interface{}, error) {
	return req, nil
}

func (m *MockProvider) TransformResponse(resp interface{}) (*ChatResponse, error) {
	return &ChatResponse{}, nil
}

func (m *MockProvider) TransformStreamChunk(chunk interface{}) (*StreamChunk, error) {
	return &StreamChunk{}, nil
}

func (m *MockProvider) InjectAuth(req interface{}, token string) (interface{}, error) {
	return req, nil
}

func (m *MockProvider) HandleStreaming(ctx context.Context, body io.ReadCloser) (io.ReadCloser, error) {
	return body, nil
}
