import json
import os
import subprocess
import time
import requests
import httpx
import asyncio
import atexit
import platform
import socket
import uuid

# Pre-defined Profile Objects (Matches Go structs)
CHROME_124 = {
    "ID": "chrome_124_macos",
    "Name": "Chrome 124 MacOS",
    "TLSID": {"Client": "Chrome", "Version": "120"},
    "InitialWinSize": 6291456,
    "TTL": 64,
    "TCPWindow": 64240,
    "UserAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
    "PseudoHeaderOrder": [":method", ":authority", ":scheme", ":path"],
}

SAFARI_17 = {
    "ID": "safari_17_ios",
    "Name": "Safari 17 on iOS",
    "TLSID": {"Client": "Safari", "Version": "16.0"},
    "InitialWinSize": 2097152,
    "TTL": 64,
    "TCPWindow": 65535,
    "UserAgent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4 Mobile/15E148 Safari/604.1",
    "PseudoHeaderOrder": [":method", ":scheme", ":path", ":authority"],
}

class Response:
    def __init__(self, data):
        self.status_code = data.get("status")
        self.text = data.get("body")
        self.headers = data.get("headers")
        self.error = data.get("error")
        self._json = None
        if self.text:
            try:
                self._json = json.loads(self.text)
            except:
                pass

    def json(self):
        """Returns the parsed JSON body of the response."""
        return self._json

class BaseClient:
    _proxy_process = None
    _port = None

    def __init__(self, profile=CHROME_124, proxy=None):
        self.profile = profile
        self.proxy = proxy
        self.headers = {}
        self.session_id = None
        self._ensure_engine_running()

    def _get_free_port(self):
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.bind(('', 0))
        port = s.getsockname()[1]
        s.close()
        return port

    def _ensure_engine_running(self):
        if BaseClient._proxy_process is None:
            BaseClient._port = self._get_free_port()
            
            ext = ".exe" if os.name == "nt" else ""
            system = platform.system().lower()
            
            # Professional Architecture Detection
            machine = platform.machine().lower()
            if machine in ["amd64", "x86_64"]:
                arch = "amd64"
            elif machine in ["arm64", "aarch64"]:
                arch = "arm64"
            else:
                arch = machine # Fallback
            
            bin_name = f"architect_{system}_{arch}{ext}"
            bin_path = os.path.join(os.path.dirname(__file__), "bin", bin_name)
            
            if not os.path.exists(bin_path):
                bin_path = os.path.join(os.getcwd(), "architect", "architect_proxy.exe")

            try:
                BaseClient._proxy_process = subprocess.Popen(
                    [bin_path, "-port", str(BaseClient._port)],
                    stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL,
                    creationflags=subprocess.CREATE_NO_WINDOW if os.name == "nt" else 0
                )
                time.sleep(1.5)
                atexit.register(self._cleanup)
            except Exception as e:
                raise RuntimeError(f"Architect Engine failed to start: {e}")

    def _cleanup(self):
        if BaseClient._proxy_process:
            BaseClient._proxy_process.terminate()
            BaseClient._proxy_process = None

    def get_logs(self):
        """Fetches latest logs from the engine."""
        try:
            r = requests.get(f"http://127.0.0.1:{BaseClient._port}/logs")
            return r.json()
        except:
            return []

    def _prepare_payload(self, method, url, body, headers):
        profile = self.profile.copy()
        
        # Compatibility layer for integer TLSIDs
        if isinstance(profile.get("TLSID"), int):
            mapping = {
                1: {"Client": "Chrome", "Version": "120"},
                2: {"Client": "Firefox", "Version": "120"},
                3: {"Client": "Safari", "Version": "16.0"}
            }
            profile["TLSID"] = mapping.get(profile["TLSID"], profile["TLSID"])

        return {
            "session_id": self.session_id,
            "profile": profile,
            "url": url,
            "method": method,
            "headers": {**self.headers, **(headers or {})},
            "body": body if isinstance(body, str) else json.dumps(body),
            "proxy": self.proxy
        }

class Client(BaseClient):
    def get(self, url, headers=None):
        return self._request("GET", url, None, headers)

    def post(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return self._request("POST", url, body, headers)

    def put(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return self._request("PUT", url, body, headers)

    def delete(self, url, headers=None):
        return self._request("DELETE", url, None, headers)

    def patch(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return self._request("PATCH", url, body, headers)

    def head(self, url, headers=None):
        return self._request("HEAD", url, None, headers)

    def options(self, url, headers=None):
        return self._request("OPTIONS", url, None, headers)

    def _request(self, method, url, body, headers):
        payload = self._prepare_payload(method, url, body, headers)
        r = requests.post(f"http://127.0.0.1:{BaseClient._port}/", json=payload, timeout=60)
        return Response(r.json())

class AsyncClient(BaseClient):
    async def get(self, url, headers=None):
        return await self._request("GET", url, None, headers)

    async def post(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return await self._request("POST", url, body, headers)

    async def put(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return await self._request("PUT", url, body, headers)

    async def delete(self, url, headers=None):
        return await self._request("DELETE", url, None, headers)

    async def patch(self, url, data=None, json_data=None, headers=None):
        body = json_data if json_data else data
        return await self._request("PATCH", url, body, headers)

    async def head(self, url, headers=None):
        return await self._request("HEAD", url, None, headers)

    async def options(self, url, headers=None):
        return await self._request("OPTIONS", url, None, headers)

    async def _request(self, method, url, body, headers):
        payload = self._prepare_payload(method, url, body, headers)
        async with httpx.AsyncClient() as client:
            r = await client.post(f"http://127.0.0.1:{BaseClient._port}/", json=payload, timeout=60)
            return Response(r.json())

class Session(Client):
    def __init__(self, profile=CHROME_124, proxy=None):
        super().__init__(profile, proxy)
        self.session_id = str(uuid.uuid4())

class AsyncSession(AsyncClient):
    def __init__(self, profile=CHROME_124, proxy=None):
        super().__init__(profile, proxy)
        self.session_id = str(uuid.uuid4())
