package provider

import (
	"context"
	"io"
)

type GoogleProvider struct {
	baseURL string
}

func NewGoogleProvider(baseURL string) *GoogleProvider {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com"
	}
	return &GoogleProvider{baseURL: baseURL}
}

func (p *GoogleProvider) Name() string {
	return "google"
}

func (p *GoogleProvider) TransformRequest(req *ChatRequest) (interface{}, error) {
	geminiReq := &GeminiRequest{
		GenerationConfig: &GenerationConfig{
			Temperature:     req.Temperature,
			TopP:            req.TopP,
			MaxOutputTokens: req.MaxTokens,
		},
	}

	var contents []GeminiContent
	for _, msg := range req.Messages {
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, GeminiContent{
			Role: role,
			Parts: []GeminiPart{
				{Text: msg.Content},
			},
		})
	}
	geminiReq.Contents = contents

	return geminiReq, nil
}

func (p *GoogleProvider) TransformResponse(resp interface{}) (*ChatResponse, error) {
	geminiResp, ok := resp.(*GeminiResponse)
	if !ok {
		return nil, nil
	}

	var content string
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		content = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	finishReason := "stop"
	if len(geminiResp.Candidates) > 0 && geminiResp.Candidates[0].FinishReason != "" {
		finishReason = geminiResp.Candidates[0].FinishReason
	}

	return &ChatResponse{
		ID:      "",
		Object:  "chat.completion",
		Created: 0,
		Model:   "",
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: finishReason,
			},
		},
		Usage: Usage{
			PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
			CompletionTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
		},
	}, nil
}

func (p *GoogleProvider) TransformStreamChunk(chunk interface{}) (*StreamChunk, error) {
	geminiChunk, ok := chunk.(*GeminiResponse)
	if !ok {
		return nil, nil
	}

	var content string
	var finishReason *string
	if len(geminiChunk.Candidates) > 0 && len(geminiChunk.Candidates[0].Content.Parts) > 0 {
		content = geminiChunk.Candidates[0].Content.Parts[0].Text
	}
	if len(geminiChunk.Candidates) > 0 && geminiChunk.Candidates[0].FinishReason != "" {
		reason := geminiChunk.Candidates[0].FinishReason
		finishReason = &reason
	}

	return &StreamChunk{
		ID:      "",
		Object:  "chat.completion.chunk",
		Created: 0,
		Model:   "",
		Choices: []StreamChoice{
			{
				Index: 0,
				Delta: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: finishReason,
			},
		},
	}, nil
}

func (p *GoogleProvider) InjectAuth(req interface{}, token string) (interface{}, error) {
	return req, nil
}

func (p *GoogleProvider) HandleStreaming(ctx context.Context, body io.ReadCloser) (io.ReadCloser, error) {
	return body, nil
}

type GeminiRequest struct {
	Contents         []GeminiContent   `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature     *float64 `json:"temperature,omitempty"`
	TopP            *float64 `json:"topP,omitempty"`
	MaxOutputTokens *int     `json:"maxOutputTokens,omitempty"`
}

type GeminiResponse struct {
	Candidates    []GeminiCandidate `json:"candidates"`
	UsageMetadata GeminiUsage       `json:"usageMetadata"`
}

type GeminiCandidate struct {
	Content      GeminiContent `json:"content"`
	FinishReason string        `json:"finishReason"`
}

type GeminiUsage struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}
