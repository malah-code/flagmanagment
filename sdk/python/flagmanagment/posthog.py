import logging
import requests
import threading
from datetime import datetime, timezone
from openfeature.hook import Hook
from openfeature.evaluation_context import EvaluationContext
from openfeature.hook_context import HookContext
from openfeature.hook_hints import HookHints
from openfeature.flag_evaluation import FlagEvaluationDetails

logger = logging.getLogger("flagmanagment-posthog")

class PostHogHook(Hook):
    def __init__(self, posthog_api_key: str, endpoint: str = "https://us.i.posthog.com/capture/"):
        self.posthog_api_key = posthog_api_key
        self.endpoint = endpoint

    def after(self, hook_context: HookContext, details: FlagEvaluationDetails, hints: HookHints):
        context: EvaluationContext = hook_context.evaluation_context
        distinct_id = context.targeting_key if context and context.targeting_key else "anonymous-server-user"

        payload = {
            "api_key": self.posthog_api_key,
            "event": "$feature_flag_called",
            "properties": {
                "distinct_id": distinct_id,
                "$feature_flag": hook_context.flag_key,
                "$feature_flag_response": details.value,
                "timestamp": datetime.now(timezone.utc).isoformat()
            }
        }
        self._send_event_async(payload)

    def error(self, hook_context: HookContext, exception: Exception, hints: HookHints):
        context: EvaluationContext = hook_context.evaluation_context
        distinct_id = context.targeting_key if context and context.targeting_key else "anonymous-server-user"

        payload = {
            "api_key": self.posthog_api_key,
            "event": "Feature Flag Error",
            "properties": {
                "distinct_id": distinct_id,
                "$feature_flag": hook_context.flag_key,
                "error": str(exception),
                "timestamp": datetime.now(timezone.utc).isoformat()
            }
        }
        self._send_event_async(payload)

    def _send_event_async(self, payload: dict):
        def _send():
            try:
                requests.post(self.endpoint, json=payload, timeout=5)
            except Exception as e:
                logger.warning(f"Failed to send PostHog event: {e}")
        
        threading.Thread(target=_send, daemon=True).start()
