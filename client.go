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

	url string
}

// ClientOption allows setting options on the client
type ClientOption func(*Client)

// NewClient creates a new client with the base URL and given options
func NewClient(url string, options ...ClientOption) *Client {
	c := &Client{url: url}
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
	url := fmt.Sprintf("%v/generate-jwt", c.url)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		if err != nil {
			return "", fmt.Errorf("claims JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(claimsJSON))
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

// URL is the base URL
func (c *Client) URL() string {
	return c.url
}

// WithHTTPClient option is used to set the http client
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.Client = hc
	}
}
