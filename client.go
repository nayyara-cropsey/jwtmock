package jwtmock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type jwtResponse struct {
	Token string `json:"token"`
}

// Client is a wrapper for an existing JWT mock server
type Client struct {
	*http.Client

	URL string
}

// ClientOption allows setting options on the client
type ClientOption func(*Client)

// NewClient creates a new client with the base URL and given options
func NewClient(url string, options ...ClientOption) *Client {
	c := &Client{URL: url}
	for _, option := range options {
		option(c)
	}

	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	return c
}

// GenerateJWT generates a JWT token for use in authorization header.
func (c *Client) GenerateJWT(ctx context.Context, claims Claims) (string, error) {
	url := fmt.Sprintf("%v/generate-jwt", c.URL)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		if err != nil {
			return "", fmt.Errorf("claims JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(claimsJSON))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got HTTP status: %v", resp.StatusCode)
	}

	var jwtResp jwtResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwtResp); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	return jwtResp.Token, nil
}

// RegisterClient register a new client
func (c *Client) RegisterClient(ctx context.Context, registration ClientRegistration) (string, error) {
	url := fmt.Sprintf("%v/clients", c.URL)

	claimsJSON, err := json.Marshal(registration)
	if err != nil {
		if err != nil {
			return "", fmt.Errorf("client registration JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(claimsJSON))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("got HTTP status: %v", resp.StatusCode)
	}

	var jwtResp jwtResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwtResp); err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}

	return jwtResp.Token, nil
}

// WithHTTPClient option is used to set the http client
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.Client = hc
	}
}
