package main

import (
	termlyhttp "api.termly.io/client/http"
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	envPublicKey  = "TERMLY_API_PUBLIC_KEY"
	envPrivateKey = "TERMLY_API_PRIVATE_KEY"
)

func main() {
	publicKey := getRequiredEnvVariable(envPublicKey)
	privateKey := getRequiredEnvVariable(envPrivateKey)

	req, err := termlyhttp.NewHttpGet(publicKey, privateKey, "authn", nil)
	if err != nil {
		log.Fatalf("unable to create GET request: %v", err)
	}

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Unable to hit authn endpoint - %v", err)
	}

	log.Println(resp.StatusCode)
}

func getRequiredEnvVariable(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("No environment variable with name %s found\n", key)
	}

	value = strings.TrimSpace(value)
	if value == "" {
		log.Fatalf("Value for %s key is empty\n", key)
	}

	return value
}
