package architect

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// DoHResponse represents the JSON response from a DoH provider.
type DoHResponse struct {
	Status int         `json:"Status"`
	Answer []DoHAnswer `json:"Answer"`
}

type DoHAnswer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

// FetchECHConfig queries a DoH provider for the HTTPS (type 65) record of a domain
// and extracts the Encrypted Client Hello config list.
func FetchECHConfig(domain string) ([]byte, error) {
	// Using Cloudflare DoH JSON API for simplicity and reliability
	url := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=HTTPS", domain)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH query failed with status: %d", resp.StatusCode)
	}

	var dohResp DoHResponse
	if err := json.NewDecoder(resp.Body).Decode(&dohResp); err != nil {
		return nil, err
	}

	if dohResp.Status != 0 {
		return nil, fmt.Errorf("DNS query returned non-zero status: %d", dohResp.Status)
	}

	// Search for the HTTPS record (Type 65)
	for _, answer := range dohResp.Answer {
		if answer.Type == 65 {
			// Cloudflare's Data field for HTTPS usually looks like:
			// 1 . alpn="h2,http/1.1" ipv4hint="1.1.1.1" ech="AEX+DQBB..."
			parts := strings.Fields(answer.Data)
			for _, part := range parts {
				if strings.HasPrefix(part, "ech=") {
					b64 := strings.TrimPrefix(part, "ech=")
					b64 = strings.Trim(b64, `"`)
					
					decoded, err := base64.StdEncoding.DecodeString(b64)
					if err != nil {
						return nil, fmt.Errorf("failed to decode ECH config: %v", err)
					}
					return decoded, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no ECH config found in HTTPS record for %s", domain)
}
