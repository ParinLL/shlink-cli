package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the Shlink API client.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Debug      bool
}

// New creates a new Shlink API client.
func New(baseURL, apiKey string, debug bool) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Debug: debug,
	}
}

func (c *Client) doRequest(method, path string, query url.Values, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	u := c.BaseURL + path
	if query != nil && len(query) > 0 {
		u += "?" + query.Encode()
	}

	if c.Debug {
		fmt.Printf("[DEBUG] %s %s\n", method, u)
	}

	req, err := http.NewRequest(method, u, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.APIKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if c.Debug {
		fmt.Printf("[DEBUG] Status: %d\n", resp.StatusCode)
		fmt.Printf("[DEBUG] Response: %s\n", string(data))
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Detail != "" {
			return nil, resp.StatusCode, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Detail)
		}
		return nil, resp.StatusCode, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(data))
	}

	return data, resp.StatusCode, nil
}
