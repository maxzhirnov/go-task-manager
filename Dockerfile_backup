# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with specific platform target
RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/main .

# Copy web directory
COPY web ./web

EXPOSE 8080

CMD ["./main"]