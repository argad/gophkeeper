package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gophkeeper/client/internal/config"
	"net/http"
	"os"
	"time"
)

var (
	// ServerURL can be configured via environment variable
	ServerURL = getServerURL()
)

func getServerURL() string {
	if url := os.Getenv("GOPHKEEPER_SERVER_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// Client is a GophKeeper API client.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new GophKeeper API client.
func NewClient() *Client {
	// Configure TLS settings
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Allow insecure connections for development with self-signed certificates
	// In production, this should be removed or controlled via environment variable
	if os.Getenv("GOPHKEEPER_INSECURE_TLS") == "true" {
		tlsConfig.InsecureSkipVerify = true
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}
}

// Request makes an HTTP request to the GophKeeper server without authentication.
func (c *Client) Request(method, path string, body interface{}) (*http.Response, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, ServerURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// AuthenticatedRequest makes an HTTP request to the GophKeeper server with the JWT token.
func (c *Client) AuthenticatedRequest(method, path string, body interface{}) (*http.Response, error) {
	token, err := config.LoadToken()
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, ServerURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}
