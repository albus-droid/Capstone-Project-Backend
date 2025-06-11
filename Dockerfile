# Build Stage
FROM golang:1.20-alpine AS builder
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
ENV DATABASE_URL="postgres://user:pass@db:5432/homecooked?sslmode=disable"
ENV JWT_SECRET="replace-with-secure-secret"
# Entrypoint
ENTRYPOINT ["./server"]
