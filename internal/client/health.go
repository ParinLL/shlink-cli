package client

import (
	"encoding/json"
	"fmt"
)

// Health checks the health of the Shlink instance.
func (c *Client) Health() (*HealthResponse, error) {
	data, _, err := c.doRequest("GET", "/rest/health", nil, nil)
	if err != nil {
		return nil, err
	}

	var result HealthResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
