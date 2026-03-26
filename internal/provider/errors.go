package provider

import (
	"fmt"
	"net/http"
)

type ErrorTranslator interface {
	TranslateError(statusCode int, body []byte, provider string) *ProviderError
}

type ProviderError struct {
	StatusCode int
	Message    string
	Type       string
	Code       string
	Provider   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("provider error (%s): %s", e.Provider, e.Message)
}

type DefaultErrorTranslator struct{}

func NewDefaultErrorTranslator() *DefaultErrorTranslator {
	return &DefaultErrorTranslator{}
}

func (t *DefaultErrorTranslator) TranslateError(statusCode int, body []byte, provider string) *ProviderError {
	err := &ProviderError{
		StatusCode: statusCode,
		Provider:   provider,
	}

	switch statusCode {
	case http.StatusUnauthorized:
		err.Message = "Invalid API key"
		err.Type = "authentication_error"
		err.Code = "invalid_api_key"
	case http.StatusForbidden:
		err.Message = "Access denied"
		err.Type = "permission_error"
		err.Code = "access_denied"
	case http.StatusNotFound:
		err.Message = "Resource not found"
		err.Type = "invalid_request_error"
		err.Code = "not_found"
	case http.StatusTooManyRequests:
		err.Message = "Rate limit exceeded"
		err.Type = "rate_limit_error"
		err.Code = "rate_limit_exceeded"
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		err.Message = "Provider service error"
		err.Type = "server_error"
		err.Code = "provider_error"
	default:
		err.Message = "Unknown error"
		err.Type = "api_error"
		err.Code = "unknown"
	}

	return err
}
