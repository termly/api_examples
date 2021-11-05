package http

import (
	"crypto/hmac"
	"crypto/sha256"
)

func hmacSha256Hash(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
