# 🏗️ Architect

**Architect** is a high-performance network orchestration engine designed for protocol research and client emulation. It provides surgical control over the network stack, enabling developers and security researchers to move beyond the rigid defaults of standard libraries.

By separating the high-level API from a specialized Go sidecar, Architect allows for precise manipulation of TLS signatures, HTTP/2 frame structures, and TCP/IP parameters.

---

## 🔬 Core Capabilities

*   **Layer 7 (Application):** Full control over HTTP/2 header normalization and frame ordering to match specific client specifications.
*   **Layer 4 (Transport):** Native **uTLS** integration for accurate **JA3/JA4** fingerprinting and **Encrypted Client Hello (ECH)** testing.
*   **Layer 3 (Network):** Granular adjustment of **TTL** and **TCP Window Size** for operating system emulation.
*   **Persistence:** Built-in support for cookie-aware sessions and TLS session resumption.

---

## 🚀 Technical Highlights

- **Zero-Dependency Core:** Pre-compiled Go binaries are bundled directly with the Python package.
- **Asynchronous First:** Designed for high-scale concurrency using `Asyncio`.
- **Method Compliance:** Full support for `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `HEAD`, and `OPTIONS`.
- **Research Ready:** Real-time engine logs provide deep visibility into every protocol handshake.
- **Flexible Proxying:** Robust support for SOCKS5 and HTTP proxies with isolated connection pooling.

---

## 📦 Installation

```bash
pip install architect-net
```

---

## 📖 Implementation

### Asynchronous Client Emulation
```python
import asyncio
import architect

async def main():
    # Emulate a specific client environment
    session = architect.AsyncSession(architect.CHROME_124)
    
    # Execute a request with precise protocol signatures
    response = await session.get("https://tls.peet.ws/api/all")
    
    if response.status_code == 200:
        data = response.json()
        print(f"Verified JA3 Hash: {data['tls']['ja3_hash']}")

if __name__ == "__main__":
    asyncio.run(main())
```

### Custom Protocol Profiling
Architect allows you to define custom network parameters for granular testing:

```python
RESEARCH_PROFILE = {
    "ID": "custom_research_node",
    "TLSID": 1,          # Maps to internal TLS profiles
    "TTL": 128,          # Define custom hop limits
    "TCPWindow": 65535,  # Define initial window size
    "UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
}

async def run_test():
    # Use AsyncClient for non-session-based requests
    client = architect.AsyncClient(profile=RESEARCH_PROFILE)
    response = await client.get("https://tls.peet.ws/api/all")
    print(response.json())
```

---

## 📜 Research Integrity
This project is intended for **authorized security research, protocol analysis, and educational use**. Architect was built to help engineers understand network behavior and build more resilient systems. Users are responsible for ensuring their use of this tool complies with all relevant laws and ethical standards.

---
*Maintained by Hiericho*
