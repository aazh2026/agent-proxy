package provider

import (
	"context"
	"io"
)

type AnthropicProvider struct {
	baseURL    string
	apiVersion string
}

func NewAnthropicProvider(baseURL string) *AnthropicProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com"
	}
	return &AnthropicProvider{
		baseURL:    baseURL,
		apiVersion: "2023-06-01",
	}
}

func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

func (p *AnthropicProvider) TransformRequest(req *ChatRequest) (interface{}, error) {
	anthropicReq := &AnthropicRequest{
		Model:       req.Model,
		MaxTokens:   4096,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
	}

	if req.MaxTokens != nil {
		anthropicReq.MaxTokens = *req.MaxTokens
	}

	var system string
	var messages []AnthropicMessage
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			system = msg.Content
		} else {
			messages = append(messages, AnthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	anthropicReq.System = system
	anthropicReq.Messages = messages

	return anthropicReq, nil
}

func (p *AnthropicProvider) TransformResponse(resp interface{}) (*ChatResponse, error) {
	anthropicResp, ok := resp.(*AnthropicResponse)
	if !ok {
		return nil, nil
	}

	var content string
	for _, block := range anthropicResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return &ChatResponse{
		ID:      anthropicResp.ID,
		Object:  "chat.completion",
		Created: 0,
		Model:   anthropicResp.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: anthropicResp.StopReason,
			},
		},
		Usage: Usage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
	}, nil
}

func (p *AnthropicProvider) TransformStreamChunk(chunk interface{}) (*StreamChunk, error) {
	anthropicChunk, ok := chunk.(*AnthropicStreamChunk)
	if !ok {
		return nil, nil
	}

	var content string
	var finishReason *string
	if anthropicChunk.Type == "content_block_delta" && anthropicChunk.Delta != nil {
		content = anthropicChunk.Delta.Text
	}
	if anthropicChunk.Type == "message_delta" && anthropicChunk.Delta != nil {
		finishReason = &anthropicChunk.Delta.StopReason
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

func (p *AnthropicProvider) InjectAuth(req interface{}, token string) (interface{}, error) {
	return req, nil
}

func (p *AnthropicProvider) HandleStreaming(ctx context.Context, body io.ReadCloser) (io.ReadCloser, error) {
	return body, nil
}

type AnthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	System      string             `json:"system,omitempty"`
	Messages    []AnthropicMessage `json:"messages"`
	Temperature *float64           `json:"temperature,omitempty"`
	TopP        *float64           `json:"top_p,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Role       string             `json:"role"`
	Content    []AnthropicContent `json:"content"`
	Model      string             `json:"model"`
	StopReason string             `json:"stop_reason"`
	Usage      AnthropicUsage     `json:"usage"`
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type AnthropicStreamChunk struct {
	Type  string                `json:"type"`
	Index int                   `json:"index"`
	Delta *AnthropicStreamDelta `json:"delta,omitempty"`
}

type AnthropicStreamDelta struct {
	Type       string `json:"type,omitempty"`
	Text       string `json:"text,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
}
