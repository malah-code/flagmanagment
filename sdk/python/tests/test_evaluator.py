import pytest
from flagmanagment.evaluator import Evaluator

def test_evaluator_flag():
    flag = {
        "enabled": True,
        "defaultVariant": "off",
        "variants": {
            "on": True,
            "off": False
        }
    }
    val, variant, reason = Evaluator.evaluate_flag(flag, "user_123")
    assert val is False
    assert variant == "off"
    assert reason == "DEFAULT"
