// Package npms implements a client for the npms.io API
//
// See: https://api-docs.npms.io
package npms

import (
	"github.com/dghubble/sling"
)

const (
	defaultBaseURL = "https://api.npms.io/v2/"
)

type links struct {
	NPM        *string `json:"npm"`
	Homepage   *string `json:"homepage"`
	Repository *string `json:"repository"`
	Bugs       *string `json:"bugs"`
}

// Package contains a package description returned by the v2/search API
type Package struct {
	Name        string       `json:"name"`
	Score       *string      `json:"scope"`
	Version     string       `json:"version"`
	Description *string      `json:"description"`
	Keywords    []string     `json:"keywords"`
	Links       *links       `json:"links"`
	Author      *interface{} `json:"author"`
	License     *string      `json:"license"`
	Date        *string      `json:"date"`
}

// Client interacts with the npms.io v2 API
type Client struct {
	BaseURL string
	Search  SearchService
	Package PackageService
}

func (c *Client) base() *sling.Sling {
	return sling.New().Base(c.BaseURL).Set("Accept", "application/json")
}

// NewClient produces a new Client using sensible defaults
func NewClient() *Client {
	c := Client{BaseURL: defaultBaseURL}
	c.Search = SearchService{&c}
	c.Package = PackageService{&c}
	return &c
}
