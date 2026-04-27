# 🏗️ Architect

> **"Advanced Network Fingerprint Research & Evasion Engine"** 🛡️  
> *Surgical network stack analysis for security researchers and developers.*

[![PyPI Version](https://img.shields.io/pypi/v/architect-net.svg?style=flat-square&color=blue)](https://pypi.org/project/architect-net/)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://github.com/Hiericho/architect)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://github.com/Hiericho/architect/blob/main/LICENSE)

**Architect** is a high-performance network engine designed for security research and protocol analysis. It provides granular control over the entire network stack—from Layer 3 TCP/IP parameters to Layer 7 HTTP/2 frame ordering—enabling researchers to emulate various client environments and study network fingerprinting behaviors.

Ideal for **Red Team operations, penetration testing, and verifying the robustness of WAF/IDS configurations**. 🔍✨

---

## 🔬 Why Architect?

Most networking libraries (like Python's `requests` or Go's `net/http`) produce rigid, easily identifiable digital signatures. Architect allows for "protocol-level transparency," moving beyond default behaviors to provide a customizable sidecar engine that handles the low-level protocol handshake.

### 🧩 The Research Stack:
*   **Layer 3 (Network):** Analyze the impact of **TTL** and **TCP Window Size** on network identification.
*   **Layer 4 (Transport):** Full **uTLS** implementation for studying **JA3/JA4** signatures and **Encrypted Client Hello (ECH)** adoption.
*   **Layer 7 (Application):** Wire-level **HTTP/2 Header Normalization** and frame analysis to match specific browser implementations.
*   **Stateful Analysis:** Built-in **Cookie Session** management and TLS session resumption testing.

---

## 🚀 Key Features

- **⚡ Professional Tooling:** Bundled cross-compiled Go sidecars. No Go installation required for Python users.
- **⚡ High-Scale Concurrency:** Native `AsyncClient` and `AsyncSession` support for high-throughput testing.
- **🌐 Full HTTP Method Support:** Native support for `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS`.
- **🔄 Profile Management:** Rapidly switch between pre-defined client profiles (Chrome, Safari, etc.) or inject custom research parameters.
- **🔌 Advanced Proxy Integration:** Robust SOCKS5 and HTTP proxy support with isolated connection pools.
- **🔍 Diagnostic Logs:** Real-time logging to inspect handshake details and protocol transitions.

---

## 📦 Installation

### Standard Installation 💅
```bash
pip install architect-net
```

### Development Setup 🏗️
```bash
# 1. Clone and compile the sidecar engine
python scripts/build_engine.py

# 2. Install in editable mode
pip install -e .
```

---

## 📖 Usage Examples

### 🧪 Client Emulation (Async)
```python
import asyncio
import architect

async def analyze_fingerprint():
    # Emulate a modern browser environment
    session = architect.AsyncSession(architect.CHROME_124)
    
    # Route through a diagnostic proxy if needed
    session.proxy = "socks5://user:pass@p.proxy.net:8000"
    
    # Analyze how the target sees your fingerprint
    response = await session.get("https://tls.peet.ws/api/all")
    print(f"Status: {response.status_code}")
    
    # Inspect the protocol details
    data = response.json()
    print(f"JA3 Signature: {data['tls']['ja3_hash']}")

if __name__ == "__main__":
    asyncio.run(analyze_fingerprint())
```

### 🧬 Custom Protocol Profiling
Define specific network parameters for testing:

```python
RESEARCH_PROFILE = {
    "ID": "win_10_research",
    "TLSID": 1,          # Maps to Chrome TLS fingerprint
    "TTL": 128,          # Windows-standard TTL
    "TCPWindow": 65535,
    "UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
}

async def run_test():
    client = architect.AsyncClient(profile=RESEARCH_PROFILE)
    response = await client.get("https://tls.peet.ws/api/all")
    print(response.json())
```

---

## 📜 Ethical Use & Disclaimer
Architect is a tool for **authorized security testing and educational research only**. It is designed to help professionals understand and secure network protocols. The authors do not condone or support unauthorized access or malicious use of this software.

---
*Maintained with integrity by Hiericho*
