import json
import os
import subprocess
import time
import requests
import atexit
import platform
import socket
import uuid

# Identity Constants
CHROME_124 = "chrome_124_macos"
SAFARI_17 = "safari_17_ios"

class Response:
    def __init__(self, data):
        self.status_code = data.get("status")
        self.text = data.get("body")
        self.headers = data.get("headers")
        self.error = data.get("error")

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
            arch = "amd64" 
            
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

    def get(self, url, headers=None):
        return self._request("GET", url, headers=headers)

    def post(self, url, data=None, headers=None):
        return self._request("POST", url, body=data, headers=headers)

    def _request(self, method, url, body="", headers=None):
        payload = {
            "session_id": self.session_id,
            "profile_id": self.profile,
            "url": url,
            "method": method,
            "headers": {**self.headers, **(headers or {})},
            "body": body if isinstance(body, str) else json.dumps(body),
            "proxy": self.proxy
        }
        
        try:
            r = requests.post(f"http://127.0.0.1:{BaseClient._port}/", json=payload, timeout=60)
            res_data = r.json()
        except Exception as e:
            raise ConnectionError(f"Architect Engine Communication Error: {e}")

        response = Response(res_data)
        if response.error:
            raise ValueError(f"Architect Engine Error: {response.error}")
            
        return response

class Client(BaseClient):
    """A single-request client (no persistent cookies)."""
    pass

class Session(BaseClient):
    """A persistent session (maintains cookies and TLS session)."""
    def __init__(self, profile=CHROME_124, proxy=None):
        super().__init__(profile, proxy)
        self.session_id = str(uuid.uuid4())
