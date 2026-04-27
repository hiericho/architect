package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hiericho/architect/core"
)

func main() {
	// Initialize Architect with a Chrome 124 profile
	tr := &architect.Transport{
		DefaultProfile:     architect.Chrome124,
		InsecureSkipVerify: true,
	}
	client := &http.Client{Transport: tr}

	// Diagnostic endpoint for TLS/HTTP2 fingerprints
	url := "https://tls.peet.ws/api/all"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Modern browsers use specific header ordering and casing (H2 is lowercase, but pseudo-headers come first)
	// Go's http.Client handles pseudo-headers correctly for H2.
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Pretty print the result to see the JA3/JA4 and HTTP/2 Fingerprints
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// If not JSON, print raw
		fmt.Printf("Response: %s\n", body)
	} else {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(prettyJSON))
	}
}
