# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build binary for both architectures
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main-amd64 ./cmd/api
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main-arm64 ./cmd/api

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy binaries and web directory
COPY --from=builder /build/main-* ./
COPY web ./web

# Create entrypoint script
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'ARCH=$(uname -m)' >> /entrypoint.sh && \
    echo 'if [ "$ARCH" = "x86_64" ]; then' >> /entrypoint.sh && \
    echo '    cp main-amd64 main' >> /entrypoint.sh && \
    echo 'elif [ "$ARCH" = "aarch64" ]; then' >> /entrypoint.sh && \
    echo '    cp main-arm64 main' >> /entrypoint.sh && \
    echo 'fi' >> /entrypoint.sh && \
    echo 'chmod +x main' >> /entrypoint.sh && \
    echo 'exec ./main' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]