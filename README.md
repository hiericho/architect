# Architect

Architect is a low-level network engine designed for perfect browser fingerprint emulation. It bypasses enterprise-grade WAFs (Cloudflare, Akamai, Datadome) by surgically manipulating the network stack across multiple layers.

## Features

- **Layer 3/4 (Network/Transport)**: TCP TTL normalization and uTLS integration for perfect JA3/JA4 signatures.
- **Layer 4 (Security)**: Automated **Encrypted Client Hello (ECH)** resolution via DNS-over-HTTPS (DoH).
- **Layer 7 (Application)**: Full control over HTTP/2 SETTINGS frames and header ordering.
- **Python Integration**: Clean Python API that manages a high-performance Go engine in the background.

## Quick Start (Python)

```python
import architect

# Identity rotation without restarting
client = architect.Client(architect.CHROME_124)
response = client.get("https://tls.peet.ws/api/all")

print(f"Bypassed: {response.status_code}")
```

## Low-Level Access (Go)

```go
import "github.com/hiericho/architect"

client := architect.NewClient(architect.Chrome124)
resp, _ := client.Get("https://target.com")
```

## Installation

1. **Build the Engine**:
   `go build -o architect/bin/architect_win_amd64.exe ./engine/main.go`
   
2. **Install Python Package**:
   `pip install -e .`

## Disclaimer
Architect is for educational and authorized security testing purposes only.
