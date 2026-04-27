# 🏗️ Architect 

> **"The Stealthy Network Chameleon"** 🦎  
> *Surgical network fingerprinting for the modern web.*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat-square&logo=python)](https://www.python.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

**Architect** is a low-level network engine that helps you blend in. It surgically manipulates your network stack across multiple layers to produce TLS signatures, HTTP/2 frames, and header ordering that are indistinguishable from real browsers.

Say goodbye to "Access Denied" from Cloudflare, Akamai, and DataDome. 🛡️✨

---

## 🛠️ How the Magic Happens

Architect handles the heavy lifting at three critical layers of the stack:

| Layer | Scope | The Architect Touch |
| :--- | :--- | :--- |
| **L3/L4** | **TCP + TLS** | 🦎 **Cloaking:** TTL normalization & uTLS for perfect JA3/JA4 signatures. |
| **L4** | **Security** | 🔐 **Secret Handshake:** Automated Encrypted Client Hello (ECH) via DoH. |
| **L7** | **Application** | ✉️ **Perfect Delivery:** Full control over H2 SETTINGS & Header ordering. |

A high-performance **Go engine** handles the surgical precision, while a **clean Python API** keeps your code looking beautiful. 💅

---

## 🚀 Getting Started

### 1. Build the Engine 🏗️
Architect uses a background proxy to do the heavy lifting. Compile it for your platform:

```bash
# Windows
go build -o architect/bin/architect_win_amd64.exe ./engine/main.go

# Linux
go build -o architect/bin/architect_linux_amd64 ./engine/main.go
```

### 2. Install Python Package 🐍
```bash
pip install -e .
```

---

## 📖 Usage Examples

### 🐍 Python (Simple & Sweet)
```python
import architect

# Just pick a profile and go! 
client = architect.Client(architect.CHROME_124)

# Architect handles the background engine lifecycle automatically ✨
response = client.get("https://tls.peet.ws/api/all")

print(f"Status: {response.status_code} - We are in! 🎉")
```

### 🐹 Go (Pure Performance)
```go
import "github.com/hiericho/architect"

func main() {
    client := architect.NewClient(architect.Chrome124)
    resp, _ := client.Get("https://tls.peet.ws/api/all")
}
```

---

## 📂 Project Structure

```text
architect/
├── engine/         # ⚙️ Go Proxy Engine (The Brains)
├── architect/      # 🐍 Python Package (The Beauty)
│   └── bin/        # 📦 Compiled Binaries
├── core/           # 🧩 Go Library Logic
└── README.md       # 📖 You are here!
```

---

## 📜 Disclaimer
Architect is intended for **educational use and authorized security testing only**. Don't be a mean chameleon! The authors assume no liability for misuse. 

---
*Made with ❤️ by Hiericho*
