package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hiericho/architect"
)

func main() {
	manager := architect.NewProfileManager()
	
	// Create a single Transport that uses the manager
	tr := &architect.Transport{
		Manager:            manager,
		DefaultProfile:     architect.Chrome124,
		InsecureSkipVerify: true,
	}
	client := &http.Client{Transport: tr}

	url := "https://tls.peet.ws/api/all"

	// 1. Request with Chrome (Default)
	fmt.Println("--- Requesting with Chrome (Default) ---")
	performRequest(client, context.Background(), architect.Chrome124, url)

	// 2. Request with Safari iOS (Dynamic Rotation)
	fmt.Println("\n--- Requesting with Safari iOS (Rotation) ---")
	ctx := architect.WithProfile(context.Background(), architect.Safari17iOS)
	performRequest(client, ctx, architect.Safari17iOS, url)
}

func performRequest(client *http.Client, ctx context.Context, p architect.Profile, url string) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	
	// Ensure User-Agent matches the profile
	req.Header.Set("user-agent", p.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	// Extract fingerprints to confirm rotation
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err == nil {
		fmt.Printf("Profile: %s\n", p.Name)
		if tls, ok := data["tls"].(map[string]interface{}); ok {
			fmt.Printf("  JA3 Hash: %v\n", tls["ja3_hash"])
			fmt.Printf("  JA4:      %v\n", tls["ja4"])
		}
	} else {
		fmt.Println("Failed to parse JSON response")
	}
}
