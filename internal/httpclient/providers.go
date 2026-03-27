package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderGoogle    Provider = "google"
)

type ProviderClient struct {
	*Client
	provider Provider
}

func NewOpenAIClient(apiKey string, config *Config) *ProviderClient {
	if config == nil {
		config = DefaultConfig()
	}
	return &ProviderClient{
		Client:   NewClient("https://api.openai.com/v1", apiKey, "openai", config),
		provider: ProviderOpenAI,
	}
}

func NewAnthropicClient(apiKey string, config *Config) *ProviderClient {
	if config == nil {
		config = DefaultConfig()
	}
	return &ProviderClient{
		Client:   NewClient("https://api.anthropic.com/v1", apiKey, "anthropic", config),
		provider: ProviderAnthropic,
	}
}

func NewGoogleClient(apiKey string, config *Config) *ProviderClient {
	if config == nil {
		config = DefaultConfig()
	}
	return &ProviderClient{
		Client:   NewClient("https://generativelanguage.googleapis.com/v1beta", apiKey, "google", config),
		provider: ProviderGoogle,
	}
}

func (pc *ProviderClient) Provider() Provider {
	return pc.provider
}

func (pc *ProviderClient) ChatCompletion(ctx context.Context, req interface{}) (*http.Response, error) {
	switch pc.provider {
	case ProviderOpenAI:
		return pc.openAIChatCompletion(ctx, req)
	case ProviderAnthropic:
		return pc.anthropicChatCompletion(ctx, req)
	case ProviderGoogle:
		return pc.googleChatCompletion(ctx, req)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", pc.provider)
	}
}

func (pc *ProviderClient) Embedding(ctx context.Context, req interface{}) (*http.Response, error) {
	switch pc.provider {
	case ProviderOpenAI:
		return pc.openAIEmbedding(ctx, req)
	case ProviderGoogle:
		return pc.googleEmbedding(ctx, req)
	default:
		return nil, fmt.Errorf("provider %s does not support embeddings", pc.provider)
	}
}

func (pc *ProviderClient) openAIChatCompletion(ctx context.Context, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+pc.apiKey)

	return pc.Do(httpReq)
}

func (pc *ProviderClient) openAIEmbedding(ctx context.Context, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+pc.apiKey)

	return pc.Do(httpReq)
}

func (pc *ProviderClient) anthropicChatCompletion(ctx context.Context, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", pc.baseURL+"/messages", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", pc.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	return pc.Do(httpReq)
}

func (pc *ProviderClient) googleChatCompletion(ctx context.Context, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/models/gemini-pro:generateContent?key=%s", pc.baseURL, pc.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	return pc.Do(httpReq)
}

func (pc *ProviderClient) googleEmbedding(ctx context.Context, req interface{}) (*http.Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/models/embedding-001:embedContent?key=%s", pc.baseURL, pc.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	return pc.Do(httpReq)
}

func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ParseJSON(resp *http.Response, v interface{}) error {
	body, err := ReadBody(resp)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

type ProviderError struct {
	StatusCode int    `json:"-"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider error (%d): %s", e.StatusCode, e.Message)
}

func HandleErrorResponse(resp *http.Response) error {
	var providerErr ProviderError
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ProviderError{
			StatusCode: resp.StatusCode,
			Type:       "api_error",
			Message:    fmt.Sprintf("HTTP %d", resp.StatusCode),
		}
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &providerErr); err != nil {
		return &ProviderError{
			StatusCode: resp.StatusCode,
			Type:       "api_error",
			Message:    string(body),
		}
	}

	providerErr.StatusCode = resp.StatusCode
	return &providerErr
}

type RetryConfig struct {
	MaxRetries      int
	InitialInterval time.Duration
	MaxInterval     time.Duration
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
	}
}
