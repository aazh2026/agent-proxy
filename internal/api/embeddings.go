package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/pipeline"
	"github.com/openclaw/agent-proxy/internal/provider"
	"github.com/openclaw/agent-proxy/internal/routing"
	"github.com/openclaw/agent-proxy/internal/token"
)

type EmbeddingRequest struct {
	Model          string      `json:"model"`
	Input          interface{} `json:"input"`
	EncodingFormat string      `json:"encoding_format,omitempty"`
	User           string      `json:"user,omitempty"`
}

type EmbeddingResponse struct {
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
	Model  string          `json:"model"`
	Usage  Usage           `json:"usage"`
}

type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

type EmbeddingsHandler struct {
	forwardingStage *pipeline.ForwardingStage
	tokenResolver   *token.TokenResolver
	routingHandler  *routing.RequestHandler
}

func NewEmbeddingsHandler(forwardingStage *pipeline.ForwardingStage, tokenResolver *token.TokenResolver, routingHandler *routing.RequestHandler) *EmbeddingsHandler {
	return &EmbeddingsHandler{
		forwardingStage: forwardingStage,
		tokenResolver:   tokenResolver,
		routingHandler:  routingHandler,
	}
}

func (h *EmbeddingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	var req EmbeddingRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Failed to read request body", "invalid_request_error")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_request_error")
		return
	}

	if err := validateEmbeddingRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		userID = "default"
	}

	providerName := resolveEmbeddingProvider(req.Model)
	if providerName == "" {
		WriteError(w, http.StatusNotFound, fmt.Sprintf("No provider found for model: %s", req.Model), "invalid_request_error")
		return
	}

	transformedBody, err := h.transformRequest(providerName, body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to transform request", "server_error")
		return
	}

	if h.routingHandler != nil {
		h.handleEmbeddingWithFailover(w, r, userID, providerName, req.Model, transformedBody)
		return
	}

	resolvedToken, err := h.tokenResolver.ResolveToken(userID, providerName, req.Model)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, fmt.Sprintf("No available token for provider %s: %v", providerName, err), "authentication_error")
		return
	}
	defer resolvedToken.Clear()

	pipelineReq := &pipeline.Request{
		HTTPRequest: r,
		UserID:      userID,
		Model:       req.Model,
		Provider:    providerName,
		Token:       resolvedToken.AccessToken,
		Body:        transformedBody,
		Stream:      false,
	}

	h.handleRequest(w, r, pipelineReq, providerName)
}

func (h *EmbeddingsHandler) handleRequest(w http.ResponseWriter, r *http.Request, req *pipeline.Request, providerName string) {
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	config := &pipeline.ProviderConfig{
		BaseURL: h.getBaseURL(providerName),
		Timeout: 60 * time.Second,
	}
	h.forwardingStage.RegisterProvider(providerName, config)

	forwardedReq, err := h.forwardingStage.Process(ctx, req)
	if err != nil {
		WriteError(w, http.StatusBadGateway, fmt.Sprintf("Failed to forward request: %v", err), "server_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(forwardedReq.Body)
}

func (h *EmbeddingsHandler) transformRequest(providerName string, body []byte) ([]byte, error) {
	switch providerName {
	case "openai":
		return body, nil
	case "anthropic":
		var openaiReq provider.OpenAIEmbeddingRequest
		if err := json.Unmarshal(body, &openaiReq); err != nil {
			return nil, err
		}
		anthropicReq := convertOpenAIToAnthropicEmbedding(&openaiReq)
		return json.Marshal(anthropicReq)
	case "google":
		var openaiReq provider.OpenAIEmbeddingRequest
		if err := json.Unmarshal(body, &openaiReq); err != nil {
			return nil, err
		}
		geminiReq := convertOpenAIToGeminiEmbedding(&openaiReq)
		return json.Marshal(geminiReq)
	default:
		return body, nil
	}
}

func (h *EmbeddingsHandler) getBaseURL(providerName string) string {
	switch providerName {
	case "openai":
		return "https://api.openai.com/v1"
	case "anthropic":
		return "https://api.anthropic.com/v1"
	case "google":
		return "https://generativelanguage.googleapis.com/v1beta"
	default:
		return ""
	}
}

func resolveEmbeddingProvider(model string) string {
	model = strings.ToLower(model)
	switch {
	case strings.Contains(model, "claude-embedding"):
		return "anthropic"
	case strings.Contains(model, "embedding") || strings.Contains(model, "embed"):
		return "openai"
	case strings.Contains(model, "gemini"):
		return "google"
	default:
		return "openai"
	}
}

func validateEmbeddingRequest(req *EmbeddingRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}

	if req.Input == nil {
		return fmt.Errorf("input is required")
	}

	validFormats := map[string]bool{
		"float":  true,
		"base64": true,
		"":       true,
	}
	if !validFormats[req.EncodingFormat] {
		return fmt.Errorf("encoding_format must be 'float' or 'base64'")
	}

	return nil
}

type GeminiEmbeddingRequest struct {
	Model   string `json:"model"`
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

func convertOpenAIToGeminiEmbedding(openaiReq *provider.OpenAIEmbeddingRequest) *GeminiEmbeddingRequest {
	geminiReq := &GeminiEmbeddingRequest{
		Model: "models/embedding-001",
	}

	var text string
	if len(openaiReq.Input) > 0 {
		text = openaiReq.Input[0]
	}

	geminiReq.Content.Parts = []struct {
		Text string `json:"text"`
	}{
		{Text: text},
	}

	return geminiReq
}

type AnthropicEmbeddingRequest struct {
	model          string `json:"model"`
	input          string `json:"input"`
	encodingFormat string `json:"encoding_format"`
}

func convertOpenAIToAnthropicEmbedding(openaiReq *provider.OpenAIEmbeddingRequest) *AnthropicEmbeddingRequest {
	var text string
	if len(openaiReq.Input) > 0 {
		text = openaiReq.Input[0]
	}

	return &AnthropicEmbeddingRequest{
		model:          openaiReq.Model,
		input:          text,
		encodingFormat: openaiReq.EncodingFormat,
	}
}

func (h *EmbeddingsHandler) handleEmbeddingWithFailover(w http.ResponseWriter, r *http.Request, userID, primaryProvider, model string, transformedBody []byte) {
	ctx := r.Context()

	err := h.routingHandler.ExecuteWithFailover(ctx, userID, model, primaryProvider, func(provider string, tok *token.Token) error {
		accessToken, err := h.tokenResolver.DecryptToken(tok.AccessTokenEncrypted)
		if err != nil {
			return err
		}

		pipelineReq := &pipeline.Request{
			HTTPRequest: r,
			UserID:      userID,
			Model:       model,
			Provider:    provider,
			Token:       accessToken,
			Body:        transformedBody,
			Stream:      false,
		}

		h.handleRequest(w, r, pipelineReq, provider)
		return nil
	})

	if err != nil {
		WriteError(w, http.StatusBadGateway, fmt.Sprintf("All providers failed: %v", err), "server_error")
	}
}
