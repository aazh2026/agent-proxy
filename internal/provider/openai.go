package provider

import (
	"context"
	"io"
)

type OpenAIProvider struct {
	baseURL string
}

func NewOpenAIProvider(baseURL string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	return &OpenAIProvider{baseURL: baseURL}
}

func (p *OpenAIProvider) Name() string {
	return "openai"
}

func (p *OpenAIProvider) TransformRequest(req *ChatRequest) (interface{}, error) {
	return req, nil
}

func (p *OpenAIProvider) TransformResponse(resp interface{}) (*ChatResponse, error) {
	if chatResp, ok := resp.(*ChatResponse); ok {
		return chatResp, nil
	}
	return nil, nil
}

func (p *OpenAIProvider) TransformStreamChunk(chunk interface{}) (*StreamChunk, error) {
	if streamChunk, ok := chunk.(*StreamChunk); ok {
		return streamChunk, nil
	}
	return nil, nil
}

func (p *OpenAIProvider) InjectAuth(req interface{}, token string) (interface{}, error) {
	return req, nil
}

func (p *OpenAIProvider) HandleStreaming(ctx context.Context, body io.ReadCloser) (io.ReadCloser, error) {
	return body, nil
}
