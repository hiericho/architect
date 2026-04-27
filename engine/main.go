package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hiericho/architect"
)

type ProxyRequest struct {
	ProfileID string            `json:"profile_id"`
	URL       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
}

type ProxyResponse struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
	Error   string            `json:"error"`
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

	// Choose profile
	var profile architect.Profile
	switch req.ProfileID {
	case "safari_17_ios":
		profile = architect.Safari17iOS
	default:
		profile = architect.Chrome124
	}

	client := architect.NewClient(profile)
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
