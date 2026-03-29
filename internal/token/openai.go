package token

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/openclaw/agent-proxy/internal/auth"
	"github.com/openclaw/agent-proxy/internal/crypto"
)

type OpenAIOAuthHandler struct {
	clientID      string
	clientSecret  string
	redirectURI   string
	authURL       string
	tokenEndpoint string
	tokenStore    *TokenStore
	encryptor     *crypto.Encryptor
	stateMu       sync.Mutex
	stateToUser   map[string]string
}

func NewOpenAIOAuthHandler(clientID, clientSecret, redirectURI string, tokenStore *TokenStore, encryptor *crypto.Encryptor) *OpenAIOAuthHandler {
	return &OpenAIOAuthHandler{
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURI:   redirectURI,
		authURL:       "https://auth.openai.com/oauth/authorize",
		tokenEndpoint: "https://auth.openai.com/oauth/token",
		tokenStore:    tokenStore,
		encryptor:     encryptor,
		stateToUser:   make(map[string]string),
	}
}

func (h *OpenAIOAuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.GetUserID(r.Context())
	if userID == "" {
		userID = r.Header.Get("X-User-ID")
	}
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if h.clientID == "" || h.clientSecret == "" || h.redirectURI == "" {
		http.Error(w, "OpenAI OAuth client is not configured", http.StatusInternalServerError)
		return
	}

	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	h.stateMu.Lock()
	h.stateToUser[state] = userID
	h.stateMu.Unlock()

	u, err := url.Parse(h.authURL)
	if err != nil {
		http.Error(w, "Invalid OpenAI auth URL", http.StatusInternalServerError)
		return
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", h.clientID)
	q.Set("redirect_uri", h.redirectURI)
	q.Set("scope", "openid offline_access")
	q.Set("state", state)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (h *OpenAIOAuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Error(w, "Missing state or code", http.StatusBadRequest)
		return
	}

	h.stateMu.Lock()
	userID, ok := h.stateToUser[state]
	if ok {
		delete(h.stateToUser, state)
	}
	h.stateMu.Unlock()

	if !ok || userID == "" {
		http.Error(w, "Invalid or expired state", http.StatusBadRequest)
		return
	}

	if h.clientID == "" || h.clientSecret == "" || h.redirectURI == "" {
		http.Error(w, "OpenAI OAuth client is not configured", http.StatusInternalServerError)
		return
	}

	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("client_id", h.clientID)
	values.Set("client_secret", h.clientSecret)
	values.Set("redirect_uri", h.redirectURI)

	resp, err := http.PostForm(h.tokenEndpoint, values)
	if err != nil {
		http.Error(w, "Failed to exchange code", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("OpenAI token exchange failed: %s", string(body)), http.StatusBadGateway)
		return
	}

	var exchangeResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		http.Error(w, "Failed to decode OpenAI token response", http.StatusInternalServerError)
		return
	}

	if exchangeResp.AccessToken == "" {
		http.Error(w, "OpenAI token response missing access_token", http.StatusInternalServerError)
		return
	}

	accessEncrypted, err := h.encryptor.Encrypt([]byte(exchangeResp.AccessToken))
	if err != nil {
		http.Error(w, "Failed to encrypt access token", http.StatusInternalServerError)
		return
	}

	refreshEncrypted := []byte{}
	if exchangeResp.RefreshToken != "" {
		refreshEncrypted, err = h.encryptor.Encrypt([]byte(exchangeResp.RefreshToken))
		if err != nil {
			http.Error(w, "Failed to encrypt refresh token", http.StatusInternalServerError)
			return
		}
	}

	tokenID, err := generateTokenID()
	if err != nil {
		http.Error(w, "Failed to generate token ID", http.StatusInternalServerError)
		return
	}

	tok := &Token{
		TokenID:               tokenID,
		UserID:                userID,
		Provider:              "openai",
		Type:                  "oauth2",
		AccessTokenEncrypted:  accessEncrypted,
		RefreshTokenEncrypted: refreshEncrypted,
		ExpiresAt:             time.Now().Unix() + exchangeResp.ExpiresIn,
		Status:                "enabled",
	}

	if err := h.tokenStore.CreateToken(tok); err != nil {
		http.Error(w, "Failed to persist OpenAI token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "ok",
		"token_id":    tokenID,
		"user_id":     userID,
		"provider":    "openai",
		"expires_at":  tok.ExpiresAt,
		"refreshable": exchangeResp.RefreshToken != "",
	})
}

func generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
