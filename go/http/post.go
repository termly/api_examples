package http

import (
	"fmt"
	"io"
	syshttp "net/http"
	"net/url"
)

func NewHttpPost(publicKey, privateKey, endpoint string, body io.ReadSeeker) (*syshttp.Request, error) {
	return NewHttpPostWithBaseURL(getBaseUrl(), publicKey, privateKey, endpoint, body)
}

func NewHttpPostWithBaseURL(
	baseURL *url.URL, publicKey, privateKey, endpoint string, body io.ReadSeeker,
) (*syshttp.Request, error) {
	endpointURL, err := resolveEndpoint(baseURL, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("newHttpPost @resolveEndpoint %w", err)
	}

	timestamp := newTimestamp()

	signature, err := CreateSignature(privateKey, syshttp.MethodPost, timestamp, endpointURL, body)
	if err != nil {
		return nil, fmt.Errorf("newHttpPost %w", err)
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("newHttpPost @seek %w", err)
	}

	httpReq, err := syshttp.NewRequest(syshttp.MethodPost, endpointURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("newHttpPost @http.NewRequest %w", err)
	}

	httpReq.Header.Add(httpTimestampHeader, timestamp)
	httpReq.Header.Add(httpAuthHeader, createAuthValue(publicKey, signature))

	return httpReq, nil
}
