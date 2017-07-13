package npms

import (
	"fmt"
	"net/http"
	"strings"
)

// QueryFloat helpfully produces a SearchQualifiers-compatible float pointer
func QueryFloat(f float32) *float32 {
	g := f
	return (*float32)(&g)
}

// QualifierFilter defines filter flags for a query
type QualifierFilter uint32

// MaxSuggestions are the v2/search/suggestions API's default limit
const MaxSuggestions = 100

const (
	NotDeprecated QualifierFilter = 1 << (32 - 1 - iota)
	NotUnstable
	NotInsecure
	IsDeprecated
	IsUnstable
	IsInsecure
)

// SearchQualifiers represents available flags for the v2/search API
type SearchQualifiers struct {
	Author     string
	BoostExact bool
	Filters    QualifierFilter

	// Keywords may be excluded py prefixing with a dash: '-test'
	Keywords          []string
	Maintainer        string
	MaintenanceWeight *float32
	PopularityWeight  *float32
	QualityWeight     *float32
	ScoreEffect       *float32
}

func floatParam(key string, val *float32) string {
	return fmt.Sprintf("%s:%.2f", key, *val)
}

// String converts the SearchQualifiers to a query string
func (q SearchQualifiers) String() string {
	var qualifiers []string

	if q.Author != "" {
		qualifiers = append(qualifiers, "author:"+q.Author)
	}

	if q.BoostExact {
		qualifiers = append(qualifiers, "boost-exact:true")
	}

	if len(q.Keywords) > 0 {
		qualifiers = append(qualifiers, "keywords:"+strings.Join(q.Keywords, ","))
	}

	if q.Filters&IsDeprecated != 0 {
		qualifiers = append(qualifiers, "is:deprecated")
	}

	if q.Filters&IsInsecure != 0 {
		qualifiers = append(qualifiers, "is:insecure")
	}

	if q.Filters&IsUnstable != 0 {
		qualifiers = append(qualifiers, "is:unstable")
	}

	if q.Maintainer != "" {
		qualifiers = append(qualifiers, "maintainer:"+q.Maintainer)
	}

	if q.MaintenanceWeight != nil {
		qualifiers = append(qualifiers, floatParam("maintenance-weight", q.MaintenanceWeight))
	}

	if q.Filters&NotDeprecated != 0 {
		qualifiers = append(qualifiers, "not:deprecated")
	}

	if q.Filters&NotInsecure != 0 {
		qualifiers = append(qualifiers, "not:insecure")
	}

	if q.Filters&NotUnstable != 0 {
		qualifiers = append(qualifiers, "not:unstable")
	}

	if q.PopularityWeight != nil {
		qualifiers = append(qualifiers, floatParam("popularity-weight", q.PopularityWeight))
	}

	if q.QualityWeight != nil {
		qualifiers = append(qualifiers, floatParam("quality-weight", q.QualityWeight))
	}

	if q.ScoreEffect != nil {
		qualifiers = append(qualifiers, floatParam("score-effect", q.ScoreEffect))
	}

	return strings.Join(qualifiers, ",")
}

// SearchQuery builds a `?q=...` param for the v2/search API
func SearchQuery(q string, qualifiers *SearchQualifiers) string {
	var parts []string
	if q != "" {
		parts = append(parts, q)
	}
	if qualifiers != nil {
		qualStr := (*qualifiers).String()
		if len(qualStr) > 0 {
			parts = append(parts, qualStr)
		}
	}

	return strings.Join(parts, " ")
}

// SearchParams bound a query to the v2/search API
type SearchParams struct {
	Q    string `url:"q"`
	Size string `url:"size,omitempty"`
	From string `url:"from,omitempty"`
}

// SearchResult is a single search result from the v2/search API
type SearchResult struct {
	Flags       map[string]interface{} `json:"flags"`
	SearchScore float32                `json:"searchScore"`
	Package     `json:"package"`
	Score       score `json:"score"`
}

// SearchResults wraps the collection of results from the v2/search API
type SearchResults struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// SearchService implements npms.io's v2/search API
type SearchService struct {
	client *Client
}

// Query invokes the v2/search API
func (s SearchService) Query(params *SearchParams) (results *SearchResults, resp *http.Response, err error) {
	results = new(SearchResults)
	resp, err = s.client.base().Get("search").QueryStruct(params).ReceiveSuccess(results)
	return
}

// SuggestionsResult wraps a single suggestion from the v2/search/suggestions API
type SuggestionsResult struct {
	SearchResult
	Highlight *string `json:"highlight"`
}

// Suggestions invokes the v2/search/suggestions API
func (s SearchService) Suggestions(params *SearchParams) (results []SuggestionsResult, resp *http.Response, err error) {
	results = make([]SuggestionsResult, MaxSuggestions)
	resp, err = s.client.base().Get("search/suggestions").QueryStruct(params).ReceiveSuccess(&results)
	return
}
