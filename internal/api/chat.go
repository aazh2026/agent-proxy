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

type ChatCompletionsHandler struct {
	forwardingStage *pipeline.ForwardingStage
	streamingProxy  *pipeline.StreamingProxy
	tokenResolver   *token.TokenResolver
	routingHandler  *routing.RequestHandler
}

func NewChatCompletionsHandler(forwardingStage *pipeline.ForwardingStage, tokenResolver *token.TokenResolver, routingHandler *routing.RequestHandler) *ChatCompletionsHandler {
	return &ChatCompletionsHandler{
		forwardingStage: forwardingStage,
		streamingProxy:  pipeline.NewStreamingProxy(forwardingStage),
		tokenResolver:   tokenResolver,
		routingHandler:  routingHandler,
	}
}

func (h *ChatCompletionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	var req ChatCompletionRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Failed to read request body", "invalid_request_error")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_request_error")
		return
	}

	if err := validateChatCompletionRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		userID = "default"
	}

	providerName := resolveProvider(req.Model)
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
		h.handleWithFailover(w, r, userID, providerName, req.Model, transformedBody, req.Stream)
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
		Stream:      req.Stream,
	}

	if req.Stream {
		h.handleStreaming(w, r, pipelineReq)
	} else {
		h.handleNonStreaming(w, r, pipelineReq, providerName)
	}
}

func (h *ChatCompletionsHandler) handleNonStreaming(w http.ResponseWriter, r *http.Request, req *pipeline.Request, providerName string) {
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

func (h *ChatCompletionsHandler) handleStreaming(w http.ResponseWriter, r *http.Request, req *pipeline.Request) {
	ctx := r.Context()

	config := &pipeline.ProviderConfig{
		BaseURL: h.getBaseURL(req.Provider),
		Timeout: 0,
	}
	h.forwardingStage.RegisterProvider(req.Provider, config)

	if err := h.streamingProxy.ProxyStreamWithContext(ctx, w, req); err != nil {
		if !strings.Contains(err.Error(), "context canceled") {
			WriteError(w, http.StatusBadGateway, fmt.Sprintf("Streaming failed: %v", err), "server_error")
		}
	}
}

func (h *ChatCompletionsHandler) transformRequest(providerName string, body []byte) ([]byte, error) {
	switch providerName {
	case "openai":
		return body, nil
	case "anthropic":
		var openaiReq provider.OpenAIChatRequest
		if err := json.Unmarshal(body, &openaiReq); err != nil {
			return nil, err
		}
		anthropicReq := provider.ConvertOpenAIToAnthropic(&openaiReq)
		return json.Marshal(anthropicReq)
	case "google":
		var openaiReq provider.OpenAIChatRequest
		if err := json.Unmarshal(body, &openaiReq); err != nil {
			return nil, err
		}
		geminiReq := provider.ConvertOpenAIToGemini(&openaiReq)
		return json.Marshal(geminiReq)
	default:
		return body, nil
	}
}

func (h *ChatCompletionsHandler) getBaseURL(providerName string) string {
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

func resolveProvider(model string) string {
	model = strings.ToLower(model)
	switch {
	case strings.HasPrefix(model, "gpt-"):
		return "openai"
	case strings.HasPrefix(model, "claude-"):
		return "anthropic"
	case strings.HasPrefix(model, "gemini-"):
		return "google"
	default:
		return ""
	}
}

func validateChatCompletionRequest(req *ChatCompletionRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}

	if len(req.Messages) == 0 {
		return fmt.Errorf("messages is required and must contain at least one message")
	}

	for i, msg := range req.Messages {
		if msg.Role == "" {
			return fmt.Errorf("message[%d].role is required", i)
		}
		validRoles := map[string]bool{
			"system":    true,
			"user":      true,
			"assistant": true,
		}
		if !validRoles[msg.Role] {
			return fmt.Errorf("message[%d].role must be one of: system, user, assistant", i)
		}
	}

	if req.Temperature != nil && (*req.Temperature < 0 || *req.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	if req.TopP != nil && (*req.TopP < 0 || *req.TopP > 1) {
		return fmt.Errorf("top_p must be between 0 and 1")
	}

	if req.MaxTokens != nil && *req.MaxTokens < 1 {
		return fmt.Errorf("max_tokens must be greater than 0")
	}

	return nil
}

func (h *ChatCompletionsHandler) handleWithFailover(w http.ResponseWriter, r *http.Request, userID, primaryProvider, model string, transformedBody []byte, stream bool) {
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
			Stream:      stream,
		}

		if stream {
			h.handleStreaming(w, r, pipelineReq)
		} else {
			h.handleNonStreaming(w, r, pipelineReq, provider)
		}
		return nil
	})

	if err != nil {
		WriteError(w, http.StatusBadGateway, fmt.Sprintf("All providers failed: %v", err), "server_error")
	}
}
