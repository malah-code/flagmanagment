import threading
import requests
import json
import time
from datetime import datetime, timezone
import logging

logger = logging.getLogger("flagmanagment-sdk")

class MetricsSync:
    def __init__(self, api_key: str, endpoint_url: str = None):
        self.api_key = api_key
        self.endpoint_url = endpoint_url or "http://localhost:8080"
        self.buffer = []
        self._lock = threading.Lock()
        self._stop_event = threading.Event()
        self._thread = None

    def start(self, interval_sec: int = 15):
        if self._thread and self._thread.is_alive():
            return
        
        self._stop_event.clear()
        self._thread = threading.Thread(target=self._run, args=(interval_sec,), daemon=True)
        self._thread.start()

    def stop(self):
        self._stop_event.set()
        if self._thread:
            self._thread.join(timeout=2.0)
        self.flush()

    def record_evaluation(self, flag_key: str):
        with self._lock:
            self.buffer.append({
                "flagKey": flag_key,
                "timestamp": datetime.now(timezone.utc).isoformat()
            })

    def _run(self, interval_sec: int):
        while not self._stop_event.wait(interval_sec):
            self.flush()

    def flush(self):
        with self._lock:
            if not self.buffer:
                return
            to_send = list(self.buffer)
            self.buffer.clear()

        payload = {
            "evaluations": to_send
        }

        url = f"{self.endpoint_url}/api/v1/sdk/metrics"
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.api_key}"
        }

        try:
            resp = requests.post(url, json=payload, headers=headers, timeout=5)
            if resp.status_code not in (200, 202):
                logger.warning(f"Unexpected status syncing metrics: {resp.status_code}")
        except Exception as e:
            logger.warning(f"Failed to sync metrics: {e}")
