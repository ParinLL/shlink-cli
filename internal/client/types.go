package client

import "time"

// APIError represents a Shlink API error response.
type APIError struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

// ShortURL represents a short URL entity.
type ShortURL struct {
	ShortCode    string     `json:"shortCode"`
	ShortURL     string     `json:"shortUrl"`
	LongURL      string     `json:"longUrl"`
	DateCreated  time.Time  `json:"dateCreated"`
	Tags         []string   `json:"tags"`
	Domain       *string    `json:"domain"`
	Title        *string    `json:"title"`
	Crawlable    bool       `json:"crawlable"`
	ForwardQuery bool       `json:"forwardQuery"`
	VisitsSummary *VisitsSummary `json:"visitsSummary"`
}

// ShortURLList represents a paginated list of short URLs.
type ShortURLList struct {
	ShortUrls struct {
		Data       []ShortURL `json:"data"`
		Pagination Pagination `json:"pagination"`
	} `json:"shortUrls"`
}

// Pagination holds pagination info.
type Pagination struct {
	CurrentPage   int `json:"currentPage"`
	PagesCount    int `json:"pagesCount"`
	ItemsPerPage  int `json:"itemsPerPage"`
	ItemsInCurrentPage int `json:"itemsInCurrentPage"`
	TotalItems    int `json:"totalItems"`
}

// CreateShortURLRequest is the request body for creating a short URL.
type CreateShortURLRequest struct {
	LongURL      string   `json:"longUrl"`
	CustomSlug   string   `json:"customSlug,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Title        string   `json:"title,omitempty"`
	Domain       string   `json:"domain,omitempty"`
	Crawlable    bool     `json:"crawlable,omitempty"`
	ForwardQuery bool     `json:"forwardQuery,omitempty"`
	FindIfExists bool     `json:"findIfExists,omitempty"`
	ValidSince   string   `json:"validSince,omitempty"`
	ValidUntil   string   `json:"validUntil,omitempty"`
	MaxVisits    int      `json:"maxVisits,omitempty"`
}

// EditShortURLRequest is the request body for editing a short URL.
type EditShortURLRequest struct {
	LongURL      *string  `json:"longUrl,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Title        *string  `json:"title,omitempty"`
	Crawlable    *bool    `json:"crawlable,omitempty"`
	ForwardQuery *bool    `json:"forwardQuery,omitempty"`
	ValidSince   *string  `json:"validSince,omitempty"`
	ValidUntil   *string  `json:"validUntil,omitempty"`
	MaxVisits    *int     `json:"maxVisits,omitempty"`
}

// VisitsSummary holds visit count info.
type VisitsSummary struct {
	Total    int `json:"total"`
	NonBots  int `json:"nonBots"`
	Bots     int `json:"bots"`
}

// VisitsResponse represents a paginated list of visits.
type VisitsResponse struct {
	Visits struct {
		Data       []Visit    `json:"data"`
		Pagination Pagination `json:"pagination"`
	} `json:"visits"`
}

// Visit represents a single visit.
type Visit struct {
	Referer   string    `json:"referer"`
	Date      time.Time `json:"date"`
	UserAgent string    `json:"userAgent"`
	Potentialbot bool   `json:"potentialBot"`
}

// VisitsOverview represents general visit stats.
type VisitsOverview struct {
	Visits struct {
		NonOrphanVisits *VisitsSummary `json:"nonOrphanVisits"`
		OrphanVisits    *VisitsSummary `json:"orphanVisits"`
	} `json:"visits"`
}

// TagsList represents a list of tags.
type TagsList struct {
	Tags struct {
		Data []string `json:"data"`
	} `json:"tags"`
}

// TagsWithStats represents tags with visit stats.
type TagsWithStats struct {
	Tags struct {
		Data []TagWithStats `json:"data"`
	} `json:"tags"`
}

// TagWithStats represents a single tag with stats.
type TagWithStats struct {
	Tag            string         `json:"tag"`
	ShortUrlsCount int            `json:"shortUrlsCount"`
	VisitsSummary  *VisitsSummary `json:"visitsSummary"`
}

// RenameTagRequest is the request body for renaming a tag.
type RenameTagRequest struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

// DomainsList represents a list of domains.
type DomainsList struct {
	Domains struct {
		Data          []DomainItem `json:"data"`
		DefaultRedirects *DomainRedirects `json:"defaultRedirects"`
	} `json:"domains"`
}

// DomainItem represents a domain entry.
type DomainItem struct {
	Domain       string           `json:"domain"`
	IsDefault    bool             `json:"isDefault"`
	Redirects    *DomainRedirects `json:"redirects"`
}

// DomainRedirects holds redirect URLs for a domain.
type DomainRedirects struct {
	BaseUrlRedirect          *string `json:"baseUrlRedirect"`
	Regular404Redirect       *string `json:"regular404Redirect"`
	InvalidShortUrlRedirect  *string `json:"invalidShortUrlRedirect"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
