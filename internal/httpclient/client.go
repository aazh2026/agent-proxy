package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	TLSHandshakeTimeout time.Duration
	DisableKeepAlives   bool
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:             60 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
	}
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	provider   string
}

func NewClient(baseURL, apiKey, provider string, config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
		TLSHandshakeTimeout: config.TLSHandshakeTimeout,
		DisableKeepAlives:   config.DisableKeepAlives,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
		provider:   provider,
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

func (c *Client) DoWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return c.httpClient.Do(req)
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (c *Client) Provider() string {
	return c.provider
}

func (c *Client) CloseIdleConnections() {
	c.httpClient.CloseIdleConnections()
}
