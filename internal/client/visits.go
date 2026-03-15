package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// GetVisitsOverview returns general visit stats.
func (c *Client) GetVisitsOverview() (*VisitsOverview, error) {
	data, _, err := c.doRequest("GET", "/rest/v3/visits", nil, nil)
	if err != nil {
		return nil, err
	}

	var result VisitsOverview
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// GetShortURLVisits returns visits for a specific short URL.
func (c *Client) GetShortURLVisits(shortCode, domain string, page, itemsPerPage int) (*VisitsResponse, error) {
	q := url.Values{}
	if domain != "" {
		q.Set("domain", domain)
	}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if itemsPerPage > 0 {
		q.Set("itemsPerPage", strconv.Itoa(itemsPerPage))
	}

	data, _, err := c.doRequest("GET", "/rest/v3/short-urls/"+shortCode+"/visits", q, nil)
	if err != nil {
		return nil, err
	}

	var result VisitsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// GetTagVisits returns visits for a specific tag.
func (c *Client) GetTagVisits(tag string, page, itemsPerPage int) (*VisitsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if itemsPerPage > 0 {
		q.Set("itemsPerPage", strconv.Itoa(itemsPerPage))
	}

	data, _, err := c.doRequest("GET", "/rest/v3/tags/"+tag+"/visits", q, nil)
	if err != nil {
		return nil, err
	}

	var result VisitsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// GetOrphanVisits returns orphan visits.
func (c *Client) GetOrphanVisits(page, itemsPerPage int) (*VisitsResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if itemsPerPage > 0 {
		q.Set("itemsPerPage", strconv.Itoa(itemsPerPage))
	}

	data, _, err := c.doRequest("GET", "/rest/v3/visits/orphan", q, nil)
	if err != nil {
		return nil, err
	}

	var result VisitsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}
