from fastapi.testclient import TestClient
from main import app

client = TestClient(app)


def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"


def test_features_evaluation():
    response = client.get("/api/features?user_id=usr-999&tier=pro")
    assert response.status_code == 200
    data = response.json()
    assert data["user_id"] == "usr-999"
    assert "flags" in data
    assert "ai_vector_search" in data["flags"]
    assert "theme" in data["flags"]


def test_recommendations():
    response = client.get("/api/recommendations?user_id=usr-999")
    assert response.status_code == 200
    data = response.json()
    assert "algorithm" in data
    assert len(data["results"]) == 2

