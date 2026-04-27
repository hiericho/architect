package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"

	"github.com/hiericho/architect"
)

type ProxyRequest struct {
	SessionID string            `json:"session_id"` // Persistent session key
	ProfileID string            `json:"profile_id"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	Proxy     string            `json:"proxy"`
}

type ProxyResponse struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Error   string            `json:"error"`
}

// Session store to maintain cookies and clients across requests
var (
	sessions   = make(map[string]*http.Client)
	sessionsMu sync.RWMutex
)

func getSessionClient(req ProxyRequest) *http.Client {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	// If no session ID, create a one-time client
	if req.SessionID == "" {
		tr := &architect.Transport{
			DefaultProfile:     getProfile(req.ProfileID),
			ProxyURL:           req.Proxy,
			InsecureSkipVerify: true,
		}
		return &http.Client{Transport: tr}
	}

	if client, ok := sessions[req.SessionID]; ok {
		return client
	}

	// Create new persistent session
	jar, _ := cookiejar.New(nil)
	tr := &architect.Transport{
		DefaultProfile:     getProfile(req.ProfileID),
		ProxyURL:           req.Proxy,
		InsecureSkipVerify: true,
	}
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}
	sessions[req.SessionID] = client
	return client
}

func getProfile(id string) architect.Profile {
	switch id {
	case "safari_17_ios":
		return architect.Safari17iOS
	default:
		return architect.Chrome124
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := getSessionClient(req)
	hReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
	if err != nil {
		sendError(w, err)
		return
	}

	for k, v := range req.Headers {
		hReq.Header.Set(k, v)
	}

	resp, err := client.Do(hReq)
	if err != nil {
		sendError(w, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	resHeaders := make(map[string]string)
	for k, v := range resp.Header {
		resHeaders[k] = strings.Join(v, ", ")
	}

	json.NewEncoder(w).Encode(ProxyResponse{
		Status:  resp.StatusCode,
		Body:    string(body),
		Headers: resHeaders,
	})
}

func sendError(w http.ResponseWriter, err error) {
	json.NewEncoder(w).Encode(ProxyResponse{Error: err.Error()})
}

func main() {
	port := flag.Int("port", 8082, "Port to listen on")
	flag.Parse()

	http.HandleFunc("/", handler)
	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Printf("Architect Proxy starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
