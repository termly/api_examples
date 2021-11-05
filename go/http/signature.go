package http

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
)

const (
	defaultService      = "default"
	defaultServiceOwner = "termly"
	parameterQuery      = "query"
	parameterPaging     = "paging"
)

func CreateSignature(privateKey, method, timestamp string, endpoint *url.URL, body io.Reader) (string, error) {
	canonicalRequest, err := createCanonicalRequest(method, timestamp, endpoint, body)
	if err != nil {
		return "", fmt.Errorf("createSignature @canonicalRequest %w", err)
	}

	derivedKey := createDerivedKey(privateKey, timestamp)

	signature := hmacSha256Hash(derivedKey, []byte(canonicalRequest))
	return hex.EncodeToString(signature), nil
}

func createDerivedKey(privateKey, timestamp string) []byte {
	key := []byte(privateKey)
	key = hmacSha256Hash(key, []byte(timestamp))
	key = hmacSha256Hash(key, []byte(defaultService))
	key = hmacSha256Hash(key, []byte(defaultServiceOwner))

	return key
}

func createCanonicalRequest(method, timestamp string, endpoint *url.URL, body io.Reader) (string, error) {
	bodyHash, err := sha256HashReader(body)
	if err != nil {
		return "", fmt.Errorf("@sha256HashReader %w", err)
	}

	hexBody := hex.EncodeToString(bodyHash)

	queryValue := getQueryOrPaging(endpoint)

	pieces := []string{
		method,
		endpoint.Hostname(),
		endpoint.EscapedPath(),
		queryValue,
		timestamp,
		hexBody,
	}

	return strings.Join(pieces, "\n"), nil
}

func getQueryOrPaging(endpoint *url.URL) string {
	values := endpoint.Query()
	if values.Has(parameterQuery) {
		return url.QueryEscape(values.Get(parameterQuery))
	}

	if values.Has(parameterPaging) {
		return url.QueryEscape(values.Get(parameterPaging))
	}

	return ""
}
