package http

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func sha256HashReader(reader io.Reader) ([]byte, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("sha256HashReader@read: %w", err)
	}

	return sha256Hash(buf), nil
}

func sha256Hash(data []byte) []byte {
	// Sum256 returns an array not a slice
	sum := sha256.Sum256(data)
	return sum[:]
}
