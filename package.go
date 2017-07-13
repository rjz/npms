package npms

import (
	"fmt"
	"net/http"
)

type collected struct {
	// npm, metadata, github, source
	NPM      *map[string]interface{} `json:"npm"`
	Metadata *map[string]interface{} `json:"metadata"`
	Github   *map[string]interface{} `json:"github"`
	Source   *map[string]interface{} `json:"source"`
}

type evaluation struct {
	Quality     map[string]interface{} `json:"quality"`
	Popularity  map[string]interface{} `json:"popularity"`
	Maintenance map[string]interface{} `json:"maintenance"`
}

type scoreDetail struct {
	Quality     float32 `json:"quality"`
	Popularity  float32 `json:"popularity"`
	Maintenance float32 `json:"maintenance"`
}

type score struct {
	Final  float64     `json:"final"`
	Detail scoreDetail `json:"detail"`
}

// PackageResult holds a single result from the v2/package API
type PackageResult struct {
	AnalyzedAt string      `json:"analyzedAt"`
	Collected  collected   `json:"collected"`
	Evaluation evaluation  `json:"evaluation"`
	Score      score       `json:"score"`
	Error      interface{} `json:"error"`
}

// PackageService implements npms.io's v2/package API
type PackageService struct {
	client *Client
}

// Get returns a single package from the v2/package API
func (s PackageService) Get(name string) (pkg *PackageResult, resp *http.Response, err error) {
	path := fmt.Sprintf("package/%s", name)
	pkg = new(PackageResult)
	resp, err = s.client.base().Get(path).ReceiveSuccess(pkg)
	return
}

// PackageMap indexes results (e.g. by package name)
type PackageMap map[string]PackageResult

// MGet returns multiple Package definitions using the v2/package/mget API
func (s PackageService) MGet(names ...string) (results *PackageMap, resp *http.Response, err error) {
	results = new(PackageMap)
	resp, err = s.client.base().Post("package/mget").BodyJSON(names).ReceiveSuccess(results)
	return
}
