use axum::{
    extract::State,
    http::StatusCode,
    response::IntoResponse,
    routing::{get, post},
    Json, Router,
};
use serde_json::Value;
use std::sync::atomic::{AtomicU64, Ordering};
use std::sync::Arc;

pub struct AppState {
    pub request_count: AtomicU64,
}

pub fn create_app() -> Router {
    let state = Arc::new(AppState {
        request_count: AtomicU64::new(0),
    });

    Router::new()
        .route("/health", get(health))
        .route("/", get(root))
        .route("/data", post(data))
        .with_state(state)
}

pub async fn health() -> impl IntoResponse {
    (StatusCode::OK, Json(serde_json::json!({"status": "ok"})))
}

pub async fn root() -> impl IntoResponse {
    (StatusCode::OK, Json(serde_json::json!({"message": "Hello, World!"})))
}

pub async fn data(
    State(state): State<Arc<AppState>>,
    Json(body): Json<Value>,
) -> impl IntoResponse {
    state.request_count.fetch_add(1, Ordering::SeqCst);
    (StatusCode::OK, Json(body))
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::Body;
    use axum::http::Request;
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_health() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .uri("/health")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = axum::body::to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["status"], "ok");
    }

    #[tokio::test]
    async fn test_root() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .uri("/")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = axum::body::to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["message"], "Hello, World!");
    }

    #[tokio::test]
    async fn test_data_valid() {
        let app = create_app();
        let body = Body::from(r#"{"key":"value"}"#);
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = axum::body::to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["key"], "value");
    }

    #[tokio::test]
    async fn test_data_empty() {
        let app = create_app();
        let body = Body::from(r#"{}"#);
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
    }

    #[tokio::test]
    async fn test_data_nested() {
        let app = create_app();
        let body = Body::from(r#"{"user":{"name":"Alice","age":25}}"#);
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);

        let body = axum::body::to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["user"]["name"], "Alice");
        assert_eq!(json["user"]["age"], 25);
    }

    #[tokio::test]
    async fn test_data_invalid_json() {
        let app = create_app();
        let body = Body::from(r#"{bad"#);
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_data_not_json() {
        let app = create_app();
        let body = Body::from("hello");
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_data_empty_body() {
        let app = create_app();
        let body = Body::empty();
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/data")
                    .header("content-type", "application/json")
                    .body(body)
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::BAD_REQUEST);
    }

    #[tokio::test]
    async fn test_health_post() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/health")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    }

    #[tokio::test]
    async fn test_root_post() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    }

    #[tokio::test]
    async fn test_data_get() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .uri("/data")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    }

    #[tokio::test]
    async fn test_not_found() {
        let app = create_app();
        let response = app
            .oneshot(
                Request::builder()
                    .uri("/nonexistent")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::NOT_FOUND);
    }
}
