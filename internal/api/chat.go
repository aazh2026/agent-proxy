package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChatCompletionsHandler struct {
}

func NewChatCompletionsHandler() *ChatCompletionsHandler {
	return &ChatCompletionsHandler{}
}

func (h *ChatCompletionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	var req ChatCompletionRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_request_error")
		return
	}

	if err := validateChatCompletionRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	if req.Stream {
		h.handleStreaming(w, r, &req)
	} else {
		h.handleNonStreaming(w, r, &req)
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

func (h *ChatCompletionsHandler) handleNonStreaming(w http.ResponseWriter, r *http.Request, req *ChatCompletionRequest) {
	response := ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   req.Model,
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: "This is a placeholder response. Provider integration not yet implemented.",
				},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     10,
			CompletionTokens: 15,
			TotalTokens:      25,
		},
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *ChatCompletionsHandler) handleStreaming(w http.ResponseWriter, r *http.Request, req *ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "Streaming not supported", "server_error")
		return
	}

	id := fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano())
	created := time.Now().Unix()

	chunk := ChatCompletionChunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   req.Model,
		Choices: []ChatCompletionDelta{
			{
				Index: 0,
				Delta: ChatMessage{
					Role:    "assistant",
					Content: "This is a placeholder streaming response.",
				},
			},
		},
	}

	data, _ := json.Marshal(chunk)
	fmt.Fprintf(w, "data: %s\n\n", data)
	flusher.Flush()

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
