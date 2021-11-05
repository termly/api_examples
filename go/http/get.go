package http

import (
	"fmt"
	syshttp "net/http"
	"net/url"
	"strings"
)

func NewHttpGet(publicKey, privateKey, endpoint string, query interface{}) (*syshttp.Request, error) {
	return NewHttpGetWithBaseURL(getBaseUrl(), publicKey, privateKey, endpoint, query)
}

func NewHttpGetWithBaseURL(
	baseURL *url.URL, publicKey, privateKey, endpoint string, query interface{},
) (*syshttp.Request, error) {
	endpointURL, err := resolveEndpoint(baseURL, endpoint, query)
	if err != nil {
		return nil, fmt.Errorf("newHttpGet @resolveEndpoint %w", err)
	}

	timestamp := newTimestamp()

	signature, err := CreateSignature(privateKey, syshttp.MethodGet, timestamp, endpointURL, strings.NewReader(""))
	if err != nil {
		return nil, fmt.Errorf("newHttpGet %w", err)
	}

	httpReq, err := syshttp.NewRequest(syshttp.MethodGet, endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("newHttpGet @http.NewRequest %w", err)
	}

	httpReq.Header.Add(httpTimestampHeader, timestamp)
	httpReq.Header.Add(httpAuthHeader, createAuthValue(publicKey, signature))

	return httpReq, nil
}
