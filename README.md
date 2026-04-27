# 🦎 Architect

> **"The Friendly Neighborhood Network Engine"** ✨  
> *Making advanced network fingerprinting as easy as pie!*

[![PyPI Version](https://img.shields.io/pypi/v/architect-net.svg?style=flat-square&color=blue)](https://pypi.org/project/architect-net/)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://github.com/Hiericho/architect)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB?style=flat-square&logo=python)](https://github.com/Hiericho/architect)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://github.com/Hiericho/architect/blob/main/LICENSE)

Hello! **Architect** is a super-powered but friendly network engine. It helps you talk to the internet just like a real browser would. Under the hood, it does some very clever math to make your network signatures (like TLS and HTTP/2) look perfect, but it keeps things simple for you on the outside. 🍬

Whether you're a security researcher, a developer, or just curious about how the web works, Architect is here to help you blend in! 🌈

---

## 🧐 How does it work?

Most libraries (like `requests`) are a bit "loud" and tell websites exactly who they are. Architect is like a talented actor—it can play the role of any browser (Chrome, Safari, and more!) by changing its "network costume" at every layer.

### 🧩 The Magic Layers:
*   **Layer 3 (Network):** Adjusts **TTL** and **TCP Window Size** to match your favorite OS.
*   **Layer 4 (Transport):** Uses **uTLS** to create perfect **JA3/JA4** signatures. It even hides its destination using **Encrypted Client Hello (ECH)**! 🤫
*   **Layer 7 (Application):** Perfectly orders **HTTP/2 Headers** so they look just like the real deal.
*   **Memory:** Remembers **Cookies** and TLS sessions so you stay logged in and fast! 🍪

---

## 🚀 Why you'll love it:

- **⚡ Easy-Peasy:** Just `pip install`. We've already compiled the heavy Go engine for you!
- **⚡ Super Fast:** Built with `Asyncio` to handle hundreds of tasks at once without breaking a sweat.
- **🌐 All the Methods:** `GET`, `POST`, `PUT`, `DELETE`... you name it, we support it!
- **📦 Smart JSON:** Responses come with a built-in `.json()` method for easy snack-sized data.
- **🔄 Quick Change:** Swap browser identities in a heartbeat or make your own!
- **🔌 Proxy Friendly:** Works beautifully with SOCKS5 and HTTP proxies to keep you moving.

---

## 📦 Getting Started

### Just the basics 💅
```bash
pip install architect-net
```

### For the builders 🏗️
```bash
# 1. Clone and build the sidecar
python scripts/build_engine.py

# 2. Install and start creating!
pip install -e .
```

---

## 📖 Let's see it in action!

### 🐍 Simple Async Example
```python
import asyncio
import architect

async def main():
    # Let's pretend we're on Chrome today! 🎩
    session = architect.AsyncSession(architect.CHROME_124)
    
    # Send a friendly request
    response = await session.get("https://tls.peet.ws/api/all")
    print(f"Success! Status: {response.status_code} 🎉")
    
    # Grab our JA3 signature from the JSON
    data = response.json()
    print(f"Our secret JA3 code: {data['tls']['ja3_hash']}")

if __name__ == "__main__":
    asyncio.run(main())
```

### 🧬 Making your own "Costume" (Custom Profile)
```python
# Create a custom Windows identity! 🪟
MY_PROFILE = {
    "ID": "custom_win_chrome",
    "TLSID": 1,          # Matches Chrome's TLS fingerprint
    "TTL": 128,          # Windows standard TTL
    "TCPWindow": 65535,
    "UserAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
}

async def run_test():
    client = architect.AsyncClient(profile=MY_PROFILE)
    response = await client.get("https://tls.peet.ws/api/all")
    print("Everything looks great! ✅")
```

---

## 📜 Being a Good Neighbor (Disclaimer)
Architect is for **authorized security testing and learning only**. We believe in keeping the internet safe and fun for everyone. Please use your new powers for good! 💖

---
*Made with ❤️ and integrity by Hiericho*
