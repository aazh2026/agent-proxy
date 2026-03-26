package provider

import (
	"context"
	"io"
)

type Provider interface {
	Name() string
	TransformRequest(req *ChatRequest) (interface{}, error)
	TransformResponse(resp interface{}) (*ChatResponse, error)
	TransformStreamChunk(chunk interface{}) (*StreamChunk, error)
	InjectAuth(req interface{}, token string) (interface{}, error)
	HandleStreaming(ctx context.Context, body io.ReadCloser) (io.ReadCloser, error)
}

type ChatRequest struct {
	Model            string
	Messages         []Message
	Temperature      *float64
	TopP             *float64
	MaxTokens        *int
	Stream           bool
	PresencePenalty  *float64
	FrequencyPenalty *float64
}

type Message struct {
	Role    string
	Content string
}

type ChatResponse struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []Choice
	Usage   Usage
}

type Choice struct {
	Index        int
	Message      Message
	FinishReason string
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type StreamChunk struct {
	ID      string
	Object  string
	Created int64
	Model   string
	Choices []StreamChoice
}

type StreamChoice struct {
	Index        int
	Delta        Message
	FinishReason *string
}

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(provider Provider) {
	r.providers[provider.Name()] = provider
}

func (r *Registry) Get(name string) (Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) List() []string {
	var names []string
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

func (r *Registry) ResolveProvider(model string) (Provider, string) {
	switch {
	case startsWith(model, "gpt-"):
		if p, ok := r.providers["openai"]; ok {
			return p, "openai"
		}
	case startsWith(model, "claude-"):
		if p, ok := r.providers["anthropic"]; ok {
			return p, "anthropic"
		}
	case startsWith(model, "gemini-"):
		if p, ok := r.providers["google"]; ok {
			return p, "google"
		}
	}
	return nil, ""
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
