# Architect

> A low-level network engine for precise browser fingerprint emulation.

Architect surgically manipulates the network stack across multiple layers to produce TLS signatures, HTTP/2 frames, and header ordering that are indistinguishable from real browsers — enabling it to pass inspection by enterprise-grade WAFs such as Cloudflare, Akamai, and DataDome.

---

## How It Works

Architect operates at three distinct layers of the network stack:

| Layer | Scope | What Architect Does |
|---|---|---|
| **L3/L4 — Network/Transport** | TCP + TLS | TTL normalization and uTLS integration for accurate JA3/JA4 signatures |
| **L4 — Security** | TLS Handshake | Automated Encrypted Client Hello (ECH) resolution via DNS-over-HTTPS (DoH) |
| **L7 — Application** | HTTP/2 | Full control over SETTINGS frames and header ordering |

A high-performance Go engine handles all low-level work. A clean Python API sits on top for easy integration.

---

## Installation

### 1. Build the Go Engine

```bash
go build -o architect/bin/architect_win_amd64.exe ./engine/main.go
```

### 2. Install the Python Package

```bash
pip install -e .
```

> **Requirements:** Go 1.21+ and Python 3.9+

---

## Usage

### Python

```python
import architect

client = architect.Client(architect.CHROME_124)
response = client.get("https://tls.peet.ws/api/all")
print(response.status_code)  # 200
```

Browser identities can be rotated at runtime without restarting the engine.

### Go

```go
import "github.com/hiericho/architect"

client := architect.NewClient(architect.Chrome124)
resp, err := client.Get("https://tls.peet.ws/api/all")
```

---

## Supported Profiles

- `CHROME_124` / `Chrome124`
- *(additional profiles go here)*

---

## Project Structure

```
architect/
├── engine/         # Go network engine
│   └── main.go
├── architect/      # Python bindings
│   └── bin/        # Compiled engine binaries
└── README.md
```

---

## Disclaimer

Architect is intended for **educational use and authorized security testing only**. Do not use it against systems you do not own or have explicit permission to test. The authors assume no liability for misuse.