import os
from typing import Optional
from fastapi import FastAPI, Query
from openfeature import api
from openfeature.evaluation_context import EvaluationContext
from flagmanagment import Client, Evaluator
from flagmanagment.provider import FlagManagmentProvider

app = FastAPI(
    title="FlagManagment Python FastAPI Sample",
    description="Demonstrates feature flag evaluation with OpenFeature in Python",
    version="1.0.0",
)

SERVER_KEY = os.getenv("FLAGMANAGMENT_SERVER_KEY", "sample-python-server-key")
BASE_URL = os.getenv("FLAGMANAGMENT_API_URL", "http://localhost:8080")

# Initialize OpenFeature provider
provider = FlagManagmentProvider(server_key=SERVER_KEY, base_url=BASE_URL)
api.set_provider(provider)
client = api.get_client()


@app.get("/health")
def health():
    return {"status": "ok", "service": "python-fastapi-sample"}


@app.get("/api/features")
def get_features(
    user_id: str = Query("user-101", description="Targeting Key for User"),
    tier: str = Query("enterprise", description="Customer Tier"),
    country: str = Query("US", description="User Country"),
):
    """
    Evaluates flags using OpenFeature evaluation context.
    """
    context = EvaluationContext(
        targeting_key=user_id,
        attributes={"tier": tier, "country": country},
    )

    # Boolean Flag with Context
    enable_ai_search = client.get_boolean_value("ai-vector-search", False, context)

    # String / Multivariate Flag
    theme = client.get_string_value("app-theme", "dark", context)

    # Number / Limit Flag
    rate_limit = client.get_integer_value("api-rate-limit", 100, context)

    return {
        "user_id": user_id,
        "context": {"tier": tier, "country": country},
        "flags": {
            "ai_vector_search": enable_ai_search,
            "theme": theme,
            "rate_limit": rate_limit,
        },
    }


@app.get("/api/recommendations")
def get_recommendations(user_id: str = "anon-user"):
    """
    Simulates dynamic algorithm selection via feature flag.
    """
    context = EvaluationContext(targeting_key=user_id)
    algo_variant = client.get_string_value("recommendation-model", "collaborative-filtering", context)

    return {
        "algorithm": algo_variant,
        "results": [
            {"id": "item-1", "title": "Advanced Microservices Architecture"},
            {"id": "item-2", "title": "Production Kubernetes Handbook"},
        ],
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

