use axum::{
    body::Body,
    http::{Request, StatusCode},
    routing::get,
    Router,
    Json,
};
use serde::Serialize;
use serde_json::Value;
use std::net::SocketAddr;
use tokio::signal;
use tower::ServiceExt;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    service: String,
}

#[derive(Serialize)]
struct InfoResponse {
    name: String,
    version: String,
    description: String,
}

#[derive(Serialize)]
struct MessageResponse {
    message: String,
    timestamp: u64,
}

async fn health() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        service: "rust-web-server".to_string(),
    })
}

async fn info() -> Json<InfoResponse> {
    Json(InfoResponse {
        name: "rust-web-server".to_string(),
        version: env!("CARGO_PKG_VERSION").to_string(),
        description: "Simple Rust web server with musl support".to_string(),
    })
}

async fn hello() -> Json<MessageResponse> {
    Json(MessageResponse {
        message: "Hello from Rust!".to_string(),
        timestamp: std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap()
            .as_secs(),
    })
}

fn app() -> Router {
    Router::new()
        .route("/health", get(health))
        .route("/info", get(info))
        .route("/hello", get(hello))
}

#[tokio::main]
async fn main() {
    let app = app();

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("Listening on {}", addr);
    
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    
    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .unwrap();
    
    println!("Server stopped gracefully");
}

async fn shutdown_signal() {
    signal::ctrl_c()
        .await
        .expect("failed to install Ctrl+C handler");
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::body::to_bytes;

    #[tokio::test]
    async fn test_health_positive() {
        let response = app()
            .oneshot(
                Request::builder()
                    .uri("/health")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
        
        let body = to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        
        assert_eq!(json["status"], "ok");
        assert_eq!(json["service"], "rust-web-server");
    }

    #[tokio::test]
    async fn test_info_positive() {
        let response = app()
            .oneshot(
                Request::builder()
                    .uri("/info")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
        
        let body = to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        
        assert_eq!(json["name"], "rust-web-server");
        assert_eq!(json["version"], "0.1.0");
    }

    #[tokio::test]
    async fn test_hello_positive() {
        let response = app()
            .oneshot(
                Request::builder()
                    .uri("/hello")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::OK);
        
        let body = to_bytes(response.into_body(), usize::MAX)
            .await
            .unwrap();
        let json: Value = serde_json::from_slice(&body).unwrap();
        
        assert_eq!(json["message"], "Hello from Rust!");
        assert!(json["timestamp"].is_number());
    }

    #[tokio::test]
    async fn test_health_negative_wrong_method() {
        let response = app()
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
    async fn test_info_negative_wrong_method() {
        let response = app()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/info")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    }

    #[tokio::test]
    async fn test_hello_negative_wrong_method() {
        let response = app()
            .oneshot(
                Request::builder()
                    .method("POST")
                    .uri("/hello")
                    .body(Body::empty())
                    .unwrap(),
            )
            .await
            .unwrap();

        assert_eq!(response.status(), StatusCode::METHOD_NOT_ALLOWED);
    }

    #[tokio::test]
    async fn test_not_found_negative() {
        let response = app()
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
