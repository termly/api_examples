package http

import (
	"fmt"
	syshttp "net/http"
	"net/url"
	"strings"
)

func NewHttpDelete(publicKey, privateKey, endpoint string, query interface{}) (*syshttp.Request, error) {
	return NewHttpDeleteWithBaseURL(getBaseUrl(), publicKey, privateKey, endpoint, query)
}

func NewHttpDeleteWithBaseURL(
	baseURL *url.URL, publicKey, privateKey, endpoint string, query interface{},
) (*syshttp.Request, error) {
	endpointURL, err := resolveEndpoint(baseURL, endpoint, query)
	if err != nil {
		return nil, fmt.Errorf("newHttpDelete @resolveEndpoint %w", err)
	}

	timestamp := newTimestamp()

	signature, err := CreateSignature(privateKey, syshttp.MethodDelete, timestamp, endpointURL, strings.NewReader(""))
	if err != nil {
		return nil, fmt.Errorf("newHttpDelete %w", err)
	}

	httpReq, err := syshttp.NewRequest(syshttp.MethodDelete, endpointURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("newHttpDelete @http.NewRequest %w", err)
	}

	httpReq.Header.Add(httpTimestampHeader, timestamp)
	httpReq.Header.Add(httpAuthHeader, createAuthValue(publicKey, signature))

	return httpReq, nil
}
