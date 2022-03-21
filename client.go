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
	url := fmt.Sprintf("%v/jwtmock/generate-jwt", c.URL)

	var jwtResp jwtResponse
	err := c.jsonRequest(ctx, url, claims, http.StatusOK, &jwtResp)
	if err != nil {
		return "", err
	}

	return jwtResp.Token, nil
}

// RegisterClient register a new client
func (c *Client) RegisterClient(ctx context.Context, registration ClientRegistration) error {
	url := fmt.Sprintf("%v/jwtmock/clients", c.URL)

	return c.jsonRequest(ctx, url, registration, http.StatusAccepted, nil)
}

// WithHTTPClient option is used to set the http client
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.Client = hc
	}
}

// jsonRequest sends a JSON request with expected status and an instance for populating response
func (c *Client) jsonRequest(ctx context.Context, url string, reqBody interface{},
	expectedStatus int, response interface{}) error {
	claimsJSON, err := json.Marshal(reqBody)
	if err != nil {
		if err != nil {
			return fmt.Errorf("marshal JSON: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(claimsJSON))
	if err != nil {
		return fmt.Errorf("create reqBody: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("got HTTP status: %v", resp.StatusCode)
	}

	if response == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("unmarshal JSON: %w", err)
	}

	return nil
}
