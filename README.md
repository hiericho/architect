# 🏗️ Architect 

> **"The Stealthy Network Chameleon"** 🦎  
> *Surgical network fingerprinting for the modern web.*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://github.com/Hiericho/architect)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![Asyncio Support](https://img.shields.io/badge/Asyncio-Supported-663399?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![PyPI Version](https://img.shields.io/pypi/v/architect-net.svg?style=flat-square)](https://pypi.org/project/architect-net/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://github.com/Hiericho/architect/blob/main/LICENSE)

**Architect** is an elite low-level network engine designed to blend in perfectly. It surgically manipulates your network stack across every layer to produce TLS signatures, HTTP/2 frames, and TCP signatures that are indistinguishable from real browsers.

Bypass enterprise-grade WAFs like **Cloudflare, Akamai, and DataDome** with ease. 🛡️✨

---

## 🧐 Why Architect?

Standard networking libraries (like Python's `requests` or Go's `net/http`) are easily flagged because they leave "digital fingerprints" at every layer. Architect wipes those fingerprints clean.

### 🦎 The Stealth Stack:
*   **Layer 3 (Network):** Spoofs **TTL** and **TCP Window Size** to match specific Operating Systems.
*   **Layer 4 (Transport):** Uses **uTLS** for perfect **JA3/JA4** signatures and **Encrypted Client Hello (ECH)** to hide SNI.
*   **Layer 7 (Application):** Wire-level **HTTP/2 Header Ordering** and frame manipulation to match browser behavior.
*   **Behavioral:** Full **Cookie Session** persistence and TLS session resumption.

---

## 🚀 Key Features

- **⚡ High Performance:** Go-powered sidecar engine for surgical precision.
- **🔄 Identity Rotation:** Switch fingerprints (Chrome, Safari, etc.) without restarting.
- **🔌 Proxy Dominance:** Full SOCKS5 and HTTP proxy support with session isolation.
- **⚡ Asyncio Native:** Designed for high-scale concurrency in Python.
- **🔍 Deep Visibility:** Stream real-time engine logs to see exactly how handshakes are performing.

---

## 📦 Installation

### 1. Build the Engine 🏗️
Architect requires the Go sidecar to be compiled for your platform:

```bash
# Cross-compiles for Windows, Linux, and macOS automatically
python build.py
```

### 2. Install Python Package 🐍
```bash
pip install architect-net
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
    
    # Peek at the digital wire:
    for log in session.get_logs():
        print(f"[ENGINE]: {log}")

asyncio.run(main())
```

### 🧬 Dynamic Custom Profiles
No need to recompile the Go engine. Define your own identity in pure Python:

```python
MY_IDENTITY = {
    "ID": "custom_m1_mac",
    "TLSID": 1,          # Chrome-based uTLS
    "TTL": 64,           # MacOS TTL
    "TCPWindow": 64240,  # MacOS Window Size
    "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)..."
}

client = architect.Client(profile=MY_IDENTITY)
```

---

## 📂 Project Structure

```text
architect/
├── engine/         # ⚙️ Go Proxy Engine (Sidecar Source)
├── architect/      # 🐍 Python Package (Public API)
│   └── bin/        # 📦 Cross-compiled Engine Binaries
├── core/           # 🧩 Low-level protocol logic
├── build.py        # 🏗️ Build automation
└── README.md       # 📖 You are here!
```

---

## 📜 Disclaimer
Architect is intended for **educational use and authorized security testing only**. Don't be a mean chameleon! The authors assume no liability for misuse. 

---
*Made with ❤️ by Hiericho*
