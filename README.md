<div align="center">

# 🏗️ Architect

**High-performance network orchestration engine for protocol research & client emulation.**

<br />

[![PyPI version](https://img.shields.io/pypi/v/architect-net?color=0ea5e9&style=flat-square&label=pypi)](https://pypi.org/project/architect-net/)
[![Python](https://img.shields.io/badge/python-3.10%2B-0ea5e9?style=flat-square&logo=python&logoColor=white)](https://python.org)
[![License](https://img.shields.io/badge/license-MIT-22c55e?style=flat-square)](LICENSE)
[![Stars](https://img.shields.io/github/stars/Hiericho/architect?color=f59e0b&style=flat-square)](https://github.com/Hiericho/architect/stargazers)
[![Research](https://img.shields.io/badge/use-security%20research-ef4444?style=flat-square)](#%EF%B8%8F-responsible-use)

<br />

> Move beyond the rigid defaults of standard libraries.  
> Architect gives you surgical control over TLS, HTTP/2, and TCP/IP — down to the byte.

<br />

[**Install**](#-installation) · [**Quick Start**](#-quick-start) · [**Capabilities**](#-capabilities) · [**Custom Profiles**](#-custom-protocol-profiles)

</div>

---

## ✨ Why Architect?

Most HTTP clients are black boxes. You send a request; something happens on the wire; you get a response. **Architect exposes everything in between.**

By separating a clean Python API from a specialized Go sidecar, Architect lets you precisely control:

- 🔐 **TLS fingerprints** — spoof or study JA3/JA4 signatures with native uTLS
- 📡 **HTTP/2 internals** — frame ordering, header normalization, SETTINGS parameters
- 🌐 **OS-level TCP stack** — TTL, window size, and ECH for deep emulation fidelity
- 🍪 **Session state** — cookie-aware persistence + TLS session resumption

---

## 🔬 Capabilities

<table>
<tr>
<td width="50%">

### 🟦 Layer 7 — Application
Full control over HTTP/2 header normalization and frame ordering to match specific client specifications (Chrome, Firefox, Safari, curl…)

</td>
<td width="50%">

### 🟨 Layer 4 — Transport
Native **uTLS** integration for accurate **JA3 / JA4** fingerprinting and **Encrypted Client Hello (ECH)** testing.

</td>
</tr>
<tr>
<td width="50%">

### 🟩 Layer 3 — Network
Granular **TTL** and **TCP Window Size** adjustment for operating system-level emulation.

</td>
<td width="50%">

### 🟪 Persistence
Built-in cookie-aware sessions and **TLS session resumption** across requests.

</td>
</tr>
</table>

---

## 🚀 Technical Highlights

| Feature | Detail |
|---|---|
| ⚡ **Zero-dependency core** | Pre-compiled Go binaries ship with the package — no Go toolchain needed |
| 🔄 **Async-first** | Designed from the ground up for `asyncio` and high-scale concurrency |
| 🛠️ **Full HTTP methods** | `GET` `POST` `PUT` `DELETE` `PATCH` `HEAD` `OPTIONS` |
| 🔍 **Deep visibility** | Real-time engine logs expose every handshake, frame, and negotiation |
| 🔀 **Flexible proxying** | SOCKS5 & HTTP proxies with per-session isolated connection pools |

---

## 📦 Installation

```bash
pip install architect-net
```

> **Requirements:** Python 3.10+

---

## ⚡ Quick Start

### Session-based client emulation

```python
import asyncio
import architect

async def main():
    # Emulate a Chrome 124 network environment
    session = architect.AsyncSession(architect.CHROME_124)

    response = await session.get("https://tls.peet.ws/api/all")

    if response.status_code == 200:
        data = response.json()
        print(f"✅ Verified JA3 Hash: {data['tls']['ja3_hash']}")

asyncio.run(main())
```

---

## 🧪 Custom Protocol Profiles

Define your own network fingerprint for granular testing and research:

```python
import architect

RESEARCH_PROFILE = {
    "ID":         "custom_research_node",
    "TLSID":      1,        # Maps to an internal TLS profile
    "TTL":        128,      # Custom hop limit (128 = Windows)
    "TCPWindow":  65535,    # Initial TCP window size
    "UserAgent":  "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
}

async def run_test():
    # AsyncClient is stateless — no session/cookie persistence
    client = architect.AsyncClient(profile=RESEARCH_PROFILE)
    response = await client.get("https://tls.peet.ws/api/all")
    print(response.json())
```

<details>
<summary><b>📋 Profile field reference</b></summary>

<br />

| Field | Type | Description |
|---|---|---|
| `ID` | `str` | Human-readable identifier for this profile |
| `TLSID` | `int` | Index into internal TLS client hello profiles |
| `TTL` | `int` | IP time-to-live value (`64` = Linux, `128` = Windows) |
| `TCPWindow` | `int` | Initial TCP receive window size |
| `UserAgent` | `str` | HTTP `User-Agent` header value |

</details>

---

## 🔄 AsyncSession vs AsyncClient

| | `AsyncSession` | `AsyncClient` |
|---|---|---|
| **Cookie persistence** | ✅ Yes | ❌ No |
| **TLS session resumption** | ✅ Yes | ❌ No |
| **Stateless requests** | ❌ | ✅ Yes |
| **Best for** | Multi-request flows, login flows | One-off probes, benchmarking |

---

## ⚖️ Responsible Use

Architect is built for **authorized security research, protocol analysis, and educational use**. It exists to help engineers understand network behavior and build more resilient systems.

> ⚠️ Users are solely responsible for ensuring their use of this tool complies with all applicable laws and ethical standards. **Do not use Architect against systems you do not own or have explicit permission to test.**

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Check the [issues page](https://github.com/Hiericho/architect/issues) to get started.

---

<div align="center">

Made with 🖤 by [Hiericho](https://github.com/Hiericho)

<br />

⭐ **Star this repo if Architect saved you time** ⭐

</div>