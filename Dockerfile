# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and required files from builder
COPY --from=builder /app/api .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations
COPY scripts/wait-for-it.sh /wait-for-it.sh

# Install necessary runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    netcat-openbsd \
    curl \
    tzdata \
    && cp /usr/share/zoneinfo/Asia/Kolkata /etc/localtime \
    && echo "Asia/Kolkata" > /etc/timezone \
    && chmod +x /wait-for-it.sh \
    && adduser -D -u 1000 appuser \
    && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Expose the application port
EXPOSE 8080

# Command to run the application will be specified in docker-compose.yml
ENTRYPOINT ["/wait-for-it.sh"]
