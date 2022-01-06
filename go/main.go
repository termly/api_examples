package main

import (
	termlyhttp "api.termly.io/client/http"
	"context"
	"io/ioutil"
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

	query := []struct {
		AccountID string `json:"account_id"`
		WebsiteID string `json:"website_id"`
		ID        string `json:"id,omitempty"`
	}{
		{
			AccountID: "acct_1234546",
			WebsiteID: "web_12345678-dead-beef-cafe-123456789012",
			ID:        "rpt_123456789012",
		},
	}

	req, err := termlyhttp.NewHttpGet(publicKey, privateKey, "websites/scan_report", &query)
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
	defer func() {
		if resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	log.Println(resp.StatusCode)

	if resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Unable to read body - %s", err)
		}

		log.Println(string(body))
	}
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
