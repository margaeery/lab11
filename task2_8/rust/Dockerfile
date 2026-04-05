FROM rust:1.85-alpine AS builder

RUN apk --no-cache add musl-dev

WORKDIR /app

COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -rf src

COPY . .
RUN cargo build --release

FROM alpine:3.19

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /app/target/release/rust_server .

EXPOSE 8080

ENV PORT=8080

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./rust_server"]
