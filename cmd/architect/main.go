package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hiericho/architect/core"
)

func main() {
	fmt.Println("=== Architect: Low-Level Network Engine ===")
	
	// 1. Target Definition
	targetDomain := "crypto.cloudflare.com"
	targetURL := fmt.Sprintf("https://%s/cdn-cgi/trace", targetDomain)

	// 2. Dynamic ECH Resolution via DoH (SNI Filtering Evasion)
	fmt.Printf("[+] Resolving ECH for %s...\n", targetDomain)
	echConfig, err := architect.FetchECHConfig(targetDomain)
	if err != nil {
		fmt.Printf("[!] ECH Warning: %v (Continuing with standard TLS 1.3)\n", err)
	} else {
		fmt.Printf("[+] ECH Config Injected (%d bytes)\n", len(echConfig))
	}

	// 3. Prepare Profiles with dynamic ECH
	chromeProfile := architect.Chrome124
	chromeProfile.ECHConfig = echConfig

	safariProfile := architect.Safari17iOS
	safariProfile.ECHConfig = echConfig

	// 4. Using the simplified Architect API
	// Phase 1: Chrome 124 (MacOS)
	fmt.Println("\n[Phase 1] Identity: Chrome 124 (MacOS)")
	executeRequest(chromeProfile, targetURL)

	// Phase 2: Safari 17 (iOS) - Rotation without restart
	fmt.Println("\n[Phase 2] Identity: Safari 17 (iOS)")
	executeRequest(safariProfile, targetURL)
}

func executeRequest(profile architect.Profile, url string) {
	// Simple API call
	client := architect.NewClient(profile)
	
	// The Transport now internally handles:
	// - TTL Normalization (Layer 3)
	// - uTLS JA3/JA4 (Layer 4)
	// - HTTP/2 Settings Interception (Layer 4/7)
	// - User-Agent Sync (Layer 7)

	req, _ := http.NewRequest("GET", url, nil)
	
	// Execution
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[!] Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Evasion Verification
	fmt.Println("--- Evasion Results ---")
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		// Look for TLS, SNI, and protocol indicators
		if containsAny(line, "tls=", "sni=", "http=", "uag=") {
			fmt.Printf("  %s\n", line)
		}
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
