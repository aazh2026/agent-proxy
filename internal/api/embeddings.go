package api

import (
	"fmt"
	"net/http"
	"time"
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
}

func NewEmbeddingsHandler() *EmbeddingsHandler {
	return &EmbeddingsHandler{}
}

func (h *EmbeddingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "invalid_request_error")
		return
	}

	var req EmbeddingRequest
	if err := DecodeJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid JSON body", "invalid_request_error")
		return
	}

	if err := validateEmbeddingRequest(&req); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error(), "invalid_request_error")
		return
	}

	h.handleRequest(w, r, &req)
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

func (h *EmbeddingsHandler) handleRequest(w http.ResponseWriter, r *http.Request, req *EmbeddingRequest) {
	var inputs []string
	switch v := req.Input.(type) {
	case string:
		inputs = []string{v}
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				inputs = append(inputs, s)
			}
		}
	}

	data := make([]EmbeddingData, len(inputs))
	for i := range inputs {
		data[i] = EmbeddingData{
			Object:    "embedding",
			Embedding: []float64{0.1, 0.2, 0.3, 0.4, 0.5},
			Index:     i,
		}
	}

	response := EmbeddingResponse{
		Object: "list",
		Data:   data,
		Model:  req.Model,
		Usage: Usage{
			PromptTokens: len(inputs) * 10,
			TotalTokens:  len(inputs) * 10,
		},
	}

	_ = time.Now()
	WriteJSON(w, http.StatusOK, response)
}
