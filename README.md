# 🏗️ Architect 

> **"The Stealthy Network Chameleon"** 🦎  
> *Surgical network fingerprinting for the modern web.*

[![PyPI Version](https://img.shields.io/pypi/v/architect-net.svg?style=flat-square&color=blue)](https://pypi.org/project/architect-net/)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://github.com/Hiericho/architect)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![Asyncio Support](https://img.shields.io/badge/Asyncio-Supported-663399?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://github.com/Hiericho/architect/blob/main/LICENSE)

**Architect** is an elite low-level network engine designed to blend in perfectly. It surgically manipulates your network stack across every layer to produce TLS signatures, HTTP/2 frames, and TCP signatures that are indistinguishable from real browsers.

Bypass enterprise-grade WAFs like **Cloudflare, Akamai, and DataDome** with ease. 🛡️✨

---

## 🧐 Why Architect?

Standard networking libraries (like Python's `requests` or Go's `net/http`) are easily flagged because they leave "digital fingerprints" at every layer. Architect wipes those fingerprints clean by running a high-performance Go sidecar that handles the raw protocol work.

### 🦎 The Stealth Stack:
*   **Layer 3 (Network):** Spoofs **TTL** and **TCP Window Size** to match specific Operating Systems.
*   **Layer 4 (Transport):** Uses **uTLS** for perfect **JA3/JA4** signatures and **Encrypted Client Hello (ECH)** to hide SNI.
*   **Layer 7 (Application):** Wire-level **HTTP/2 Header Ordering** and frame manipulation to match browser behavior.
*   **Behavioral:** Full **Cookie Session** persistence and TLS session resumption.

---

## 🚀 Key Features

- **⚡ Zero-Friction User Experience:** Just `pip install`. Pre-compiled binaries for Windows, Linux, and macOS are bundled in the wheel.
- **⚡ Asyncio Native:** Designed for high-scale concurrency with `AsyncClient` and `AsyncSession`.
- **🌐 Full HTTP Method Support:** Effortlessly execute `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS` requests.
- **📦 Native JSON Support:** All responses include a `.json()` method for seamless data parsing.
- **🔄 Identity Rotation:** Switch fingerprints (Chrome, Safari, etc.) or inject custom profiles on the fly.
- **🔌 Proxy Dominance:** Full SOCKS5 and HTTP proxy support with isolated connection pools.
- **🔍 Deep Visibility:** Stream real-time engine logs to see exactly how handshakes are performing.

---

## 📦 Installation

### The Easy Way (Users) 💅
```bash
pip install architect-net
```
*Architect bundles pre-compiled Go binaries—no Go installation required!*

### Building from Source 🏗️
```bash
# 1. Clone and compile the engine sidecars
python scripts/build_engine.py

# 2. Install in editable mode
pip install -e .
```

---

## 📖 Usage Examples

### 🐍 Basic Async Session
```python
import asyncio
import architect

async def main():
    # Maintains cookies & TLS state automatically!
    session = architect.AsyncSession(architect.CHROME_124)
    
    # 🦎 Perfect emulation through a residential proxy
    session.proxy = "socks5://user:pass@p.proxy.net:8000"
    
    response = await session.get("https://tls.peet.ws/api/all")
    print(f"Bypassed! Status: {response.status_code} 🎉")
    
    # Direct JSON access
    data = response.json()
    print(f"JA3 Hash: {data['tls']['ja3_hash']}")
    
    # Peek at the digital wire:
    for log in session.get_logs():
        print(f"[ENGINE]: {log}")

if __name__ == "__main__":
    asyncio.run(main())
```

### 🧬 Dynamic Custom Profiles
Define your own identity in pure Python. Architect supports both dictionary-based `TLSID` and legacy integer IDs for backward compatibility:

```python
MY_IDENTITY = {
    "ID": "custom_win_chrome",
    "TLSID": 1,          # Maps to Chrome 120 automatically
    "TTL": 128,          # Windows TTL
    "TCPWindow": 65535,  # Windows Window Size
    "UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
}

# Note: Use AsyncClient within async functions!
async def run():
    client = architect.AsyncClient(profile=MY_IDENTITY)
    response = await client.get("https://tls.peet.ws/api/all")
    print(response.json())
```

### 🐹 Go (Pure Performance)
```go
import "github.com/hiericho/architect/core"

func main() {
    client := architect.NewClient(architect.Chrome124)
    resp, _ := client.Get("https://tls.peet.ws/api/all")
}
```

---

## 📂 Project Structure

```text
architect/
├── engine/         # ⚙️ Go Proxy Engine (Sidecar Source)
├── architect/      # 🐍 Python Package (Public API)
│   └── bin/        # 📦 Bundled Cross-compiled Engine Binaries
├── scripts/        # 🏗️ Build and automation scripts
├── core/           # 🧩 Core Go Logic
└── README.md       # 📖 You are here!
```

---

## 📜 Disclaimer
Architect is intended for **educational use and authorized security testing only**. Don't be a mean chameleon! The authors assume no liability for misuse. 

---
*Made with ❤️ by Hiericho*
