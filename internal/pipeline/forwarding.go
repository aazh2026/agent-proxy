package pipeline

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ForwardingStage struct {
	httpClient *http.Client
	providers  map[string]*ProviderConfig
}

type ProviderConfig struct {
	BaseURL string
	Timeout time.Duration
}

func NewForwardingStage() *ForwardingStage {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	return &ForwardingStage{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   60 * time.Second,
		},
		providers: make(map[string]*ProviderConfig),
	}
}

func (s *ForwardingStage) RegisterProvider(name string, config *ProviderConfig) {
	s.providers[name] = config
}

func (s *ForwardingStage) Name() string {
	return "forwarding"
}

func (s *ForwardingStage) Process(ctx context.Context, req *Request) (*Request, error) {
	config, ok := s.providers[req.Provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", req.Provider)
	}

	url := s.buildURL(config.BaseURL, req.Provider, req.Model, req.Stream)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(httpReq, req.Provider, req.Token, req.Stream)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	req.Body = body
	return req, nil
}

func (s *ForwardingStage) ForwardStream(ctx context.Context, req *Request) (io.ReadCloser, http.Header, int, error) {
	config, ok := s.providers[req.Provider]
	if !ok {
		return nil, nil, 0, fmt.Errorf("unknown provider: %s", req.Provider)
	}

	url := s.buildURL(config.BaseURL, req.Provider, req.Model, true)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(req.Body))
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	s.setHeaders(httpReq, req.Provider, req.Token, true)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to send request: %w", err)
	}

	return resp.Body, resp.Header, resp.StatusCode, nil
}

func (s *ForwardingStage) buildURL(baseURL, provider, model string, stream bool) string {
	// Detect embeddings requests based on model name patterns
	isEmbedding := isEmbeddingModel(model)

	switch provider {
	case "openai":
		if isEmbedding {
			return baseURL + "/embeddings"
		}
		return baseURL + "/chat/completions"
	case "anthropic":
		// Anthropic doesn't have a separate embeddings endpoint in the same format
		// They use the messages API with specific models
		return baseURL + "/messages"
	case "google":
		if isEmbedding {
			return fmt.Sprintf("%s/models/%s:embedContent", baseURL, model)
		}
		return fmt.Sprintf("%s/models/%s:generateContent", baseURL, model)
	default:
		return baseURL + "/chat/completions"
	}
}

// isEmbeddingModel determines if a model is used for embeddings based on naming conventions
func isEmbeddingModel(model string) bool {
	modelLower := strings.ToLower(model)
	return strings.Contains(modelLower, "embedding") ||
		strings.Contains(modelLower, "embed") ||
		strings.HasPrefix(modelLower, "text-embedding") ||
		strings.HasPrefix(modelLower, "gemini-embedding")
}

func (s *ForwardingStage) setHeaders(req *http.Request, provider, token string, stream bool) {
	req.Header.Set("Content-Type", "application/json")

	switch provider {
	case "openai":
		req.Header.Set("Authorization", "Bearer "+token)
	case "anthropic":
		req.Header.Set("x-api-key", token)
		req.Header.Set("anthropic-version", "2023-06-01")
	case "google":
		q := req.URL.Query()
		q.Set("key", token)
		req.URL.RawQuery = q.Encode()
	}

	if stream {
		req.Header.Set("Accept", "text/event-stream")
	}
}
