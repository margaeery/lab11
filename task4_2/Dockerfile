FROM --platform=$BUILDPLATFORM rust:alpine AS builder

RUN apk add --no-cache musl-dev build-base

WORKDIR /app

COPY Cargo.toml ./
COPY src ./src

RUN cargo build --release --target x86_64-unknown-linux-musl

RUN cp /app/target/x86_64-unknown-linux-musl/release/rust-web-server /server

FROM scratch

COPY --from=builder /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
