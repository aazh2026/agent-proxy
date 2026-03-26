package token

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/crypto"
)

type TokenHandler struct {
	tokenStore *TokenStore
	encryptor  *crypto.Encryptor
}

func NewTokenHandler(db *sql.DB, encryptor *crypto.Encryptor) *TokenHandler {
	return &TokenHandler{
		tokenStore: NewTokenStore(db),
		encryptor:  encryptor,
	}
}

type CreateTokenRequest struct {
	Provider      string   `json:"provider"`
	Type          string   `json:"type"`
	AccessToken   string   `json:"access_token"`
	RefreshToken  string   `json:"refresh_token,omitempty"`
	ExpiresAt     int64    `json:"expires_at,omitempty"`
	Status        string   `json:"status,omitempty"`
	Priority      int      `json:"priority,omitempty"`
	AllowedModels []string `json:"allowed_models,omitempty"`
}

type TokenResponse struct {
	TokenID       string   `json:"token_id"`
	UserID        string   `json:"user_id"`
	Provider      string   `json:"provider"`
	Type          string   `json:"type"`
	ExpiresAt     int64    `json:"expires_at"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	Status        string   `json:"status"`
	Priority      int      `json:"priority"`
	AllowedModels []string `json:"allowed_models"`
}

func (h *TokenHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Provider == "" || req.AccessToken == "" {
		writeError(w, http.StatusBadRequest, "Provider and access_token are required")
		return
	}

	accessTokenEncrypted, err := h.encryptor.Encrypt([]byte(req.AccessToken))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to encrypt token")
		return
	}

	var refreshTokenEncrypted []byte
	if req.RefreshToken != "" {
		refreshTokenEncrypted, err = h.encryptor.Encrypt([]byte(req.RefreshToken))
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to encrypt refresh token")
			return
		}
	}

	tokenID, err := generateTokenID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to generate token ID")
		return
	}

	token := &Token{
		TokenID:               tokenID,
		UserID:                userID,
		Provider:              req.Provider,
		Type:                  req.Type,
		AccessTokenEncrypted:  accessTokenEncrypted,
		RefreshTokenEncrypted: refreshTokenEncrypted,
		ExpiresAt:             req.ExpiresAt,
		Status:                "enabled",
		Priority:              req.Priority,
		AllowedModels:         req.AllowedModels,
	}

	if err := h.tokenStore.CreateToken(token); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TokenResponse{
		TokenID:       token.TokenID,
		UserID:        token.UserID,
		Provider:      token.Provider,
		Type:          token.Type,
		ExpiresAt:     token.ExpiresAt,
		CreatedAt:     token.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     token.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Status:        token.Status,
		Priority:      token.Priority,
		AllowedModels: token.AllowedModels,
	})
}

func (h *TokenHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokens, err := h.tokenStore.GetTokensByUser(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get tokens")
		return
	}

	var response []TokenResponse
	for _, token := range tokens {
		response = append(response, TokenResponse{
			TokenID:       token.TokenID,
			UserID:        token.UserID,
			Provider:      token.Provider,
			Type:          token.Type,
			ExpiresAt:     token.ExpiresAt,
			CreatedAt:     token.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:     token.UpdatedAt.Format("2006-01-02T15:04:05Z"),
			Status:        token.Status,
			Priority:      token.Priority,
			AllowedModels: token.AllowedModels,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TokenHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenID := r.URL.Query().Get("token_id")
	if tokenID == "" {
		writeError(w, http.StatusBadRequest, "token_id is required")
		return
	}

	token, err := h.tokenStore.GetToken(tokenID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get token")
		return
	}
	if token == nil || token.UserID != userID {
		writeError(w, http.StatusNotFound, "Token not found")
		return
	}

	var req CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Status != "" {
		token.Status = req.Status
	}
	if req.Priority != 0 {
		token.Priority = req.Priority
	}
	if req.AllowedModels != nil {
		token.AllowedModels = req.AllowedModels
	}

	if err := h.tokenStore.UpdateToken(token); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TokenHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenID := r.URL.Query().Get("token_id")
	if tokenID == "" {
		writeError(w, http.StatusBadRequest, "token_id is required")
		return
	}

	token, err := h.tokenStore.GetToken(tokenID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get token")
		return
	}
	if token == nil || token.UserID != userID {
		writeError(w, http.StatusNotFound, "Token not found")
		return
	}

	if err := h.tokenStore.DeleteToken(tokenID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "tk_" + hex.EncodeToString(bytes), nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
			"type":    "api_error",
		},
	})
}
