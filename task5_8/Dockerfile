FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server .

FROM alpine:latest

COPY --from=builder /server /server

EXPOSE 8080

LABEL com.centurylinklabs.watchtower.enable="true"
LABEL com.centurylinklabs.watchtower.lifecycle.pre-check="/bin/sh -c 'echo Updating'"
LABEL com.centurylinklabs.watchtower.lifecycle.post-check="/bin/sh -c 'echo Updated'"

ENTRYPOINT ["/server"]
