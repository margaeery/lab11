use axum::{
    routing::get,
    Router,
    Json,
};
use serde::Serialize;
use std::net::SocketAddr;
use tokio::signal;

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

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/health", get(health))
        .route("/info", get(info))
        .route("/hello", get(hello));

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
