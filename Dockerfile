# Build Stage
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy source
COPY . .
# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# Runtime Stage
FROM alpine:3.18
WORKDIR /app
# Copy built binary
COPY --from=builder /app/server ./server
# Expose port
EXPOSE 8080
# Environment variables (can override in CI/production)
# Entrypoint
ENTRYPOINT ["./server"]
