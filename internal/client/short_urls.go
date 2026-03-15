package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// ListShortURLs lists short URLs with optional pagination.
func (c *Client) ListShortURLs(page, itemsPerPage int, searchTerm, orderBy string, tags []string) (*ShortURLList, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if itemsPerPage > 0 {
		q.Set("itemsPerPage", strconv.Itoa(itemsPerPage))
	}
	if searchTerm != "" {
		q.Set("searchTerm", searchTerm)
	}
	if orderBy != "" {
		q.Set("orderBy", orderBy)
	}
	for _, tag := range tags {
		q.Add("tags[]", tag)
	}

	data, _, err := c.doRequest("GET", "/rest/v3/short-urls", q, nil)
	if err != nil {
		return nil, err
	}

	var result ShortURLList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// CreateShortURL creates a new short URL.
func (c *Client) CreateShortURL(req *CreateShortURLRequest) (*ShortURL, error) {
	data, _, err := c.doRequest("POST", "/rest/v3/short-urls", nil, req)
	if err != nil {
		return nil, err
	}

	var result ShortURL
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// GetShortURL retrieves a short URL by its short code.
func (c *Client) GetShortURL(shortCode, domain string) (*ShortURL, error) {
	q := url.Values{}
	if domain != "" {
		q.Set("domain", domain)
	}

	data, _, err := c.doRequest("GET", "/rest/v3/short-urls/"+shortCode, q, nil)
	if err != nil {
		return nil, err
	}

	var result ShortURL
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// EditShortURL edits an existing short URL.
func (c *Client) EditShortURL(shortCode, domain string, req *EditShortURLRequest) (*ShortURL, error) {
	q := url.Values{}
	if domain != "" {
		q.Set("domain", domain)
	}

	data, _, err := c.doRequest("PATCH", "/rest/v3/short-urls/"+shortCode, q, req)
	if err != nil {
		return nil, err
	}

	var result ShortURL
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// DeleteShortURL deletes a short URL by its short code.
func (c *Client) DeleteShortURL(shortCode, domain string) error {
	q := url.Values{}
	if domain != "" {
		q.Set("domain", domain)
	}

	_, _, err := c.doRequest("DELETE", "/rest/v3/short-urls/"+shortCode, q, nil)
	return err
}
