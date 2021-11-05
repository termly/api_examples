package http

import (
	"fmt"
	"io"
	syshttp "net/http"
	"net/url"
)

func NewHttpPut(publicKey, privateKey, endpoint string, body io.ReadSeeker) (*syshttp.Request, error) {
	return NewHttpPutWithBaseURL(getBaseUrl(), publicKey, privateKey, endpoint, body)
}

func NewHttpPutWithBaseURL(
	baseURL *url.URL, publicKey, privateKey, endpoint string, body io.ReadSeeker,
) (*syshttp.Request, error) {
	endpointURL, err := resolveEndpoint(baseURL, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("newHttpPut @resolveEndpoint %w", err)
	}

	timestamp := newTimestamp()

	signature, err := CreateSignature(privateKey, syshttp.MethodPut, timestamp, endpointURL, body)
	if err != nil {
		return nil, fmt.Errorf("newHttpPut %w", err)
	}

	_, err = body.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("newHttpPut @seek %w", err)
	}

	httpReq, err := syshttp.NewRequest(syshttp.MethodPut, endpointURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("newHttpPut @http.NewRequest %w", err)
	}

	httpReq.Header.Add(httpTimestampHeader, timestamp)
	httpReq.Header.Add(httpAuthHeader, createAuthValue(publicKey, signature))

	return httpReq, nil
}
