package client

import (
	"encoding/json"
	"fmt"
)

// ListDomains lists all configured domains.
func (c *Client) ListDomains() (*DomainsList, error) {
	data, _, err := c.doRequest("GET", "/rest/v3/domains", nil, nil)
	if err != nil {
		return nil, err
	}

	var result DomainsList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// SetDomainRedirects sets the "not found" redirects for a domain.
func (c *Client) SetDomainRedirects(domain string, redirects *DomainRedirects) error {
	body := map[string]interface{}{
		"domain": domain,
	}
	if redirects.BaseUrlRedirect != nil {
		body["baseUrlRedirect"] = *redirects.BaseUrlRedirect
	}
	if redirects.Regular404Redirect != nil {
		body["regular404Redirect"] = *redirects.Regular404Redirect
	}
	if redirects.InvalidShortUrlRedirect != nil {
		body["invalidShortUrlRedirect"] = *redirects.InvalidShortUrlRedirect
	}

	_, _, err := c.doRequest("PATCH", "/rest/v3/domains/redirects", nil, body)
	return err
}
