package http

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	defaultApiUrl = "https://api.termly.io"
	rootPath      = "/v1/"
)

var (
	lookupEnv = os.LookupEnv
)

func getBaseUrl() *url.URL {
	value, ok := lookupEnv("TERMLY_API_URL")
	if !ok {
		value = defaultApiUrl
	}

	value = strings.TrimSpace(value)
	if value == "" {
		value = defaultApiUrl
	}

	baseUrl, err := url.Parse(value)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse url %s - %v\n", value, err))
	}

	baseUrl.Path = rootPath
	baseUrl.RawQuery = ""
	baseUrl.Opaque = ""

	return baseUrl
}

func resolveEndpoint(baseURL *url.URL, endpoint string, query interface{}) (*url.URL, error) {
	// they should just pass in things like "websites", "collaborators", etc. But
	// if they do pass in more "/collaborators", "v1/collaborators", or "/v1/collaborators"
	// we trim it off.
	resource := strings.TrimPrefix(endpoint, "/")
	resource = strings.TrimPrefix(resource, "v1/")
	resource = strings.TrimSuffix(resource, "/")

	newURL, err := baseURL.Parse(resource)
	if err != nil {
		return nil, fmt.Errorf("@parse: %w", err)
	}

	if query != nil {
		queryBuf, err := json.Marshal(query)
		if err != nil {
			return nil, fmt.Errorf("@marshal: %w", err)
		}

		queryValues := url.Values{}
		queryValues.Set("query", string(queryBuf))
		newURL.RawQuery = queryValues.Encode()
	}

	return newURL, nil
}
