package npms_test

import (
	"encoding/json"
	"fmt"
	"github.com/rjz/npms"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	tc     *npms.Client
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	tc = npms.NewClient()
	tc.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

func expectStringEqual(t *testing.T, desc, expected, actual string) {
	if expected != actual {
		t.Error(fmt.Sprintf("Expected %s to be '%s'; got '%s'", desc, expected, actual))
	}
}

func expectQueryParam(t *testing.T, r *http.Request, key, val string) {
	expectStringEqual(t, fmt.Sprintf("query param '%s'", key), val, r.URL.Query().Get(key))
}

func expectMethod(t *testing.T, r *http.Request, method string) {
	expectStringEqual(t, "method", method, r.Method)
}

func expectHeader(t *testing.T, r *http.Request, key, expected string) {
	expectStringEqual(t, "header "+key, expected, r.Header.Get(key))
}

func ExamplePackageService_get() {
	client := npms.NewClient()
	pkg, _, err := client.Package.Get("express")
	if err != nil {
		fmt.Println("Request failed", err)
	}

	fmt.Println("score is", pkg.Score.Final)
}

func TestPackageGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/package/fzbz", func(w http.ResponseWriter, r *http.Request) {
		expectMethod(t, r, "GET")
		fmt.Fprintf(w, `{}`)
	})

	tc.Package.Get("fzbz")
}

func TestPackageMGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/package/mget", func(w http.ResponseWriter, r *http.Request) {
		expectMethod(t, r, "POST")
		expectHeader(t, r, "Content-type", "application/json")
		expectHeader(t, r, "Accept", "application/json")
		body, _ := ioutil.ReadAll(r.Body)
		names := make([]string, 1)
		json.Unmarshal(body, &names)
		if len(names) != 1 || names[0] != "fzbz" {
			t.Errorf("expected json, didn't get it")
		}

		fmt.Fprintf(w, `{}`)
	})

	tc.Package.MGet("fzbz")
}

func TestSearchQuery(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		expectMethod(t, r, "GET")
		expectQueryParam(t, r, "q", "fzbz author:rjz")
		fmt.Fprintf(w, `{}`)
	})

	tc.Search.Query(&npms.SearchParams{Q: "fzbz author:rjz"})
}

func TestSearchSuggestions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/suggestions", func(w http.ResponseWriter, r *http.Request) {
		expectMethod(t, r, "GET")
		expectQueryParam(t, r, "q", "fzbz")
		fmt.Fprintf(w, `{}`)
	})

	tc.Search.Suggestions(&npms.SearchParams{Q: "fzbz"})
}
