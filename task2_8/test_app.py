import pytest
from fastapi.testclient import TestClient
from app import app

client = TestClient(app)


def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


def test_root():
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"message": "Hello, World!"}


def test_data_valid():
    response = client.post("/data", json={"key": "value"})
    assert response.status_code == 200
    assert response.json() == {"key": "value"}


def test_data_empty():
    response = client.post("/data", json={})
    assert response.status_code == 400


def test_data_nested():
    response = client.post("/data", json={"user": {"name": "Alice", "age": 25}})
    assert response.status_code == 200
    assert response.json() == {"user": {"name": "Alice", "age": 25}}


def test_data_array():
    response = client.post("/data", json=[1, 2, 3])
    assert response.status_code == 422


def test_data_string_body():
    response = client.post("/data", content=b"hello", headers={"content-type": "application/json"})
    assert response.status_code == 422


def test_health_post():
    response = client.post("/health")
    assert response.status_code == 405


def test_root_post():
    response = client.post("/")
    assert response.status_code == 405


def test_data_get():
    response = client.get("/data")
    assert response.status_code == 405


def test_not_found():
    response = client.get("/nonexistent")
    assert response.status_code == 404
