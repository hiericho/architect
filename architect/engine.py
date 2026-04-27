import json
import os
import subprocess
import time
import requests
import atexit
import platform
import socket

# Identity Constants
CHROME_124 = "chrome_124_macos"
SAFARI_17 = "safari_17_ios"

class Response:
    def __init__(self, data):
        self.status_code = data.get("status")
        self.text = data.get("body")
        self.headers = data.get("headers")
        self.error = data.get("error")

class Client:
    _proxy_process = None
    _port = None

    def __init__(self, profile=CHROME_124):
        self.profile = profile
        self.headers = {}
        self._ensure_engine_running()

    def _get_free_port(self):
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.bind(('', 0))
        port = s.getsockname()[1]
        s.close()
        return port

    def _ensure_engine_running(self):
        if Client._proxy_process is None:
            Client._port = self._get_free_port()
            
            # Locate binary based on OS
            ext = ".exe" if os.name == "nt" else ""
            system = platform.system().lower()
            arch = "amd64" # Simplified for this example
            
            bin_name = f"architect_{system}_{arch}{ext}"
            bin_path = os.path.join(os.path.dirname(__file__), "bin", bin_name)
            
            if not os.path.exists(bin_path):
                # Fallback for local development
                bin_path = os.path.join(os.path.dirname(__file__), "architect_proxy.exe")

            try:
                Client._proxy_process = subprocess.Popen(
                    [bin_path, "-port", str(Client._port)],
                    stdout=subprocess.DEVNULL,
                    stderr=subprocess.DEVNULL,
                    creationflags=subprocess.CREATE_NO_WINDOW if os.name == "nt" else 0
                )
                time.sleep(1.5) # Wait for engine warm-up
                atexit.register(self._cleanup)
            except Exception as e:
                raise RuntimeError(f"Architect Engine failed to start: {e}")

    def _cleanup(self):
        if Client._proxy_process:
            Client._proxy_process.terminate()
            Client._proxy_process = None

    def get(self, url, headers=None):
        return self._request("GET", url, headers=headers)

    def post(self, url, data=None, headers=None):
        return self._request("POST", url, body=data, headers=headers)

    def _request(self, method, url, body="", headers=None):
        payload = {
            "profile_id": self.profile,
            "url": url,
            "method": method,
            "headers": {**self.headers, **(headers or {})},
            "body": body if isinstance(body, str) else json.dumps(body)
        }
        
        try:
            r = requests.post(f"http://127.0.0.1:{Client._port}/", json=payload, timeout=60)
            res_data = r.json()
        except Exception as e:
            raise ConnectionError(f"Architect Engine Communication Error: {e}")

        response = Response(res_data)
        if response.error:
            raise ValueError(f"Architect Engine Error: {response.error}")
            
        return response
