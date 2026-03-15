package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ListTags lists all tags.
func (c *Client) ListTags() (*TagsList, error) {
	data, _, err := c.doRequest("GET", "/rest/v3/tags", nil, nil)
	if err != nil {
		return nil, err
	}

	var result TagsList
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// TagsWithStats lists tags with visit stats.
func (c *Client) TagsWithStats() (*TagsWithStats, error) {
	data, _, err := c.doRequest("GET", "/rest/v3/tags/stats", nil, nil)
	if err != nil {
		return nil, err
	}

	var result TagsWithStats
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// RenameTag renames a tag.
func (c *Client) RenameTag(oldName, newName string) error {
	_, _, err := c.doRequest("PUT", "/rest/v3/tags", nil, &RenameTagRequest{
		OldName: oldName,
		NewName: newName,
	})
	return err
}

// DeleteTags deletes one or more tags.
func (c *Client) DeleteTags(tags []string) error {
	q := url.Values{}
	q.Set("tags[]", strings.Join(tags, ","))
	_, _, err := c.doRequest("DELETE", "/rest/v3/tags", q, nil)
	return err
}
