# Frontend build stage
FROM node:18-alpine AS frontend-builder
WORKDIR /frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

# Go build stage
FROM golang:1.22-alpine AS backend-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main-amd64 ./cmd/api
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main-arm64 ./cmd/api

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /build/main-* ./
COPY --from=frontend-builder /frontend/build ./frontend/build

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