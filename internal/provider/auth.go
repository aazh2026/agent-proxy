package provider

import (
	"fmt"
	"net/http"
)

type AuthInjector interface {
	InjectAuth(req *http.Request, token string, provider string) error
}

type DefaultAuthInjector struct{}

func NewDefaultAuthInjector() *DefaultAuthInjector {
	return &DefaultAuthInjector{}
}

func (i *DefaultAuthInjector) InjectAuth(req *http.Request, token string, provider string) error {
	switch provider {
	case "openai":
		req.Header.Set("Authorization", "Bearer "+token)
	case "anthropic":
		req.Header.Set("x-api-key", token)
		req.Header.Set("anthropic-version", "2023-06-01")
	case "google":
		req.Header.Set("Authorization", "Bearer "+token)
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
	return nil
}
