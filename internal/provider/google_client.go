package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GoogleClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewGoogleClient(baseURL, apiKey string, timeout time.Duration) *GoogleClient {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	}

	return &GoogleClient{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (c *GoogleClient) GenerateContent(ctx context.Context, model string, req *GeminiRequest) (*GeminiResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, model, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google API error (%d): %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &geminiResp, nil
}

func (c *GoogleClient) GenerateContentStream(ctx context.Context, model string, req *GeminiRequest) (io.ReadCloser, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?key=%s", c.baseURL, model, c.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Google API error (%d): %s", resp.StatusCode, string(body))
	}

	return resp.Body, nil
}

func ConvertOpenAIToGemini(openaiReq *OpenAIChatRequest) *GeminiRequest {
	geminiReq := &GeminiRequest{
		GenerationConfig: &GenerationConfig{
			Temperature:     openaiReq.Temperature,
			TopP:            openaiReq.TopP,
			MaxOutputTokens: openaiReq.MaxTokens,
		},
	}

	var contents []GeminiContent
	for _, msg := range openaiReq.Messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		if role == "system" {
			role = "user"
		}
		contents = append(contents, GeminiContent{
			Role: role,
			Parts: []GeminiPart{
				{Text: msg.Content},
			},
		})
	}
	geminiReq.Contents = contents

	return geminiReq
}

func ConvertGeminiToOpenAI(geminiResp *GeminiResponse) *OpenAIChatResponse {
	var content string
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	finishReason := "stop"
	if len(geminiResp.Candidates) > 0 && geminiResp.Candidates[0].FinishReason != "" {
		finishReason = geminiResp.Candidates[0].FinishReason
	}

	return &OpenAIChatResponse{
		ID:      "",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "",
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Message: OpenAIMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: finishReason,
			},
		},
		Usage: OpenAIUsage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
	}
}
