import threading
import json
import logging
import time
import urllib.request

logger = logging.getLogger("flagmanagment-sdk")


class Client:
    """SSE streaming client that connects to the FlagManagment backend,
    maintains a thread-safe in-memory cache of flag definitions,
    and supports exponential backoff reconnection."""

    def __init__(self, api_key: str, stream_url: str):
        self.api_key = api_key
        self.stream_url = stream_url
        self.flags: dict = {}
        self._lock = threading.RLock()
        self._running = False
        self._thread: threading.Thread | None = None

    def connect(self):
        """Start background SSE streaming."""
        if self._running:
            return
        self._running = True
        self._thread = threading.Thread(target=self._stream, daemon=True, name="flagmanagment-sse")
        self._thread.start()

    def shutdown(self):
        """Gracefully stop the SSE stream."""
        self._running = False

    def _stream(self):
        attempt = 0
        while self._running:
            try:
                req = urllib.request.Request(
                    self.stream_url,
                    headers={
                        "Authorization": f"Bearer {self.api_key}",
                        "Accept": "text/event-stream",
                    },
                )
                with urllib.request.urlopen(req, timeout=30) as resp:
                    if resp.status != 200:
                        attempt += 1
                        self._backoff(attempt)
                        continue

                    # Connection established — reset backoff
                    attempt = 0
                    logger.info("connected to %s", self.stream_url)

                    current_event = ""
                    for raw_line in resp:
                        if not self._running:
                            return
                        line = raw_line.decode("utf-8").strip()

                        if line.startswith("event:"):
                            current_event = line[6:].strip()
                        elif line.startswith("data:"):
                            data_str = line[5:].strip()
                            self._handle_event(current_event, data_str)
                        elif line == "":
                            current_event = ""

            except Exception as e:
                logger.warning("connection error: %s", e)

            if self._running:
                attempt += 1
                logger.info("reconnecting (attempt %d)...", attempt)
                self._backoff(attempt)

    def _handle_event(self, event_type: str, data_str: str):
        try:
            payload = json.loads(data_str)
        except json.JSONDecodeError as e:
            logger.warning("failed to parse event data: %s", e)
            return

        if event_type == "bootstrap":
            flags = payload.get("flags")
            if flags:
                with self._lock:
                    self.flags = flags
                logger.info("bootstrapped %d flags", len(flags))

        elif event_type == "flag_updated":
            flag_key = payload.get("flagKey")
            flag = payload.get("flag")
            if flag_key and flag:
                with self._lock:
                    self.flags[flag_key] = flag
                logger.info("updated flag: %s", flag_key)

        elif event_type == "ping":
            pass  # heartbeat — no action

    def _backoff(self, attempt: int):
        delay = min(attempt * attempt, 60)
        if delay < 1:
            delay = 1
        time.sleep(delay)

    def get_flag(self, key: str):
        """Return the flag definition for the given key, or None."""
        with self._lock:
            return self.flags.get(key)
