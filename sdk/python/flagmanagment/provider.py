from openfeature.provider import AbstractProvider
from openfeature.provider import Metadata
from openfeature.flag_evaluation import FlagResolutionDetails, Reason, ErrorCode
from openfeature.evaluation_context import EvaluationContext
from typing import Optional, Any

from .client import Client
from .evaluator import Evaluator


class FlagManagmentProvider(AbstractProvider):
    """OpenFeature provider for FlagManagment. Evaluates flags locally
    using the in-memory cache maintained by Client and MurmurHash3
    bucketing via Evaluator."""

    def __init__(self, client: Client):
        self._client = client

    def get_metadata(self) -> Metadata:
        return Metadata(name="FlagManagment-Python-Provider")

    def _evaluate(
        self,
        flag_key: str,
        default_value: Any,
        expected_type: Optional[type] = None,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails:
        flag = self._client.get_flag(flag_key)
        if flag is None:
            return FlagResolutionDetails(
                value=default_value,
                reason=Reason.ERROR,
                error_code=ErrorCode.FLAG_NOT_FOUND,
            )

        # Extract targeting key from context
        targeting_key = None
        if evaluation_context is not None:
            targeting_key = evaluation_context.targeting_key

        # Delegate to the local evaluator
        value, variant, reason_str = Evaluator.evaluate_flag(flag, targeting_key)

        if value is None:
            value = default_value

        if expected_type is not None:
            if expected_type is bool and not isinstance(value, bool):
                return FlagResolutionDetails(value=default_value, reason=Reason.ERROR, error_code=ErrorCode.TYPE_MISMATCH)
            elif expected_type is str and not isinstance(value, str):
                return FlagResolutionDetails(value=default_value, reason=Reason.ERROR, error_code=ErrorCode.TYPE_MISMATCH)
            elif expected_type is int and (not isinstance(value, int) or isinstance(value, bool)):
                return FlagResolutionDetails(value=default_value, reason=Reason.ERROR, error_code=ErrorCode.TYPE_MISMATCH)
            elif expected_type is float and (not isinstance(value, (float, int)) or isinstance(value, bool)):
                return FlagResolutionDetails(value=default_value, reason=Reason.ERROR, error_code=ErrorCode.TYPE_MISMATCH)

        # Map reason string to OpenFeature Reason enum
        reason_map = {
            "DISABLED": Reason.DISABLED,
            "DEFAULT": Reason.DEFAULT,
            "TARGETING_MATCH": Reason.TARGETING_MATCH,
        }
        reason = reason_map.get(reason_str, Reason.STATIC)

        return FlagResolutionDetails(value=value, variant=variant, reason=reason)

    def resolve_boolean_details(
        self,
        flag_key: str,
        default_value: bool,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails[bool]:
        return self._evaluate(flag_key, default_value, bool, evaluation_context)

    def resolve_string_details(
        self,
        flag_key: str,
        default_value: str,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails[str]:
        return self._evaluate(flag_key, default_value, str, evaluation_context)

    def resolve_integer_details(
        self,
        flag_key: str,
        default_value: int,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails[int]:
        return self._evaluate(flag_key, default_value, int, evaluation_context)

    def resolve_float_details(
        self,
        flag_key: str,
        default_value: float,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails[float]:
        return self._evaluate(flag_key, default_value, float, evaluation_context)

    def resolve_object_details(
        self,
        flag_key: str,
        default_value: dict,
        evaluation_context: Optional[EvaluationContext] = None,
    ) -> FlagResolutionDetails[dict]:
        return self._evaluate(flag_key, default_value, dict, evaluation_context)
