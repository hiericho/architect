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
	SessionID string            `json:"session_id"`
	Profile   architect.Profile `json:"profile"`
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

var (
	sessions   = make(map[string]*http.Client)
	sessionsMu sync.RWMutex
	logBuffer  []string
	logMu      sync.Mutex
)

// Custom logger to capture engine events
func engineLog(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logMu.Lock()
	logBuffer = append(logBuffer, msg)
	if len(logBuffer) > 100 { // Keep last 100 logs
		logBuffer = logBuffer[1:]
	}
	logMu.Unlock()
	log.Println(msg)
}

func getSessionClient(req ProxyRequest) *http.Client {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	profile := req.Profile
	if profile.ID == "" {
		profile = architect.Chrome124
	}

	if req.SessionID == "" {
		tr := &architect.Transport{
			DefaultProfile:     profile,
			ProxyURL:           req.Proxy,
			InsecureSkipVerify: true,
		}
		return &http.Client{Transport: tr}
	}

	if client, ok := sessions[req.SessionID]; ok {
		return client
	}

	jar, _ := cookiejar.New(nil)
	tr := &architect.Transport{
		DefaultProfile:     profile,
		ProxyURL:           req.Proxy,
		InsecureSkipVerify: true,
	}
	client := &http.Client{
		Transport: tr,
		Jar:       jar,
	}
	sessions[req.SessionID] = client
	engineLog("Created persistent session: %s", req.SessionID)
	return client
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/logs" {
		logMu.Lock()
		json.NewEncoder(w).Encode(logBuffer)
		logBuffer = []string{} // Clear after read
		logMu.Unlock()
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	engineLog("Request: %s %s [Profile: %s]", req.Method, req.URL, req.Profile.Name)

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
		engineLog("Error: %v", err)
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
