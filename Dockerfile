# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations
COPY scripts/wait-for-it.sh /wait-for-it.sh

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates netcat-openbsd && \
    chmod +x /wait-for-it.sh

# Expose the application port
EXPOSE 8080

# Command to run the application will be specified in docker-compose.yml
ENTRYPOINT ["/wait-for-it.sh"]
