# AI Agent Instructions for rest-api-go

This document provides essential context for AI coding agents working in this codebase.

## Architecture Overview

This is a RESTful API built with Go, following clean architecture principles:

- `cmd/api/main.go` - Application entry point, wires up dependencies
- `internal/` - Core application code
  - `domain/` - Business entities and interfaces
  - `handler/` - HTTP handlers with async processing
  - `service/` - Business logic layer
  - `repository/` - Data access layer
  - `config/` - Configuration management
  - `logging/` - Request tracing and logging
  - `database/` - Database initialization

### Key Design Patterns

1. **Async Processing**: All handler operations use a worker pool pattern:

   ```go
   // See internal/handler/student_handler.go
   type StudentHandler struct {
       workers int     // Number of worker goroutines (default: 50)
       jobs    chan func() // Job queue (buffer: 100)
   }
   ```

2. **Request Tracing**: Every request gets a unique trace ID:
   ```go
   // See internal/logging/request_logger.go
   traceID := uuid.New().String()
   ctx = AddTraceIDToContext(ctx, traceID)
   ```

### Infrastructure Setup

- **Docker Compose Services**:
  - API (x3 replicas)
  - MySQL 8.0
  - NGINX Load Balancer
  - Prometheus
  - Grafana

## Development Workflow

1. **Local Development**:

   ```bash
   # Start all services
   docker-compose up --build

   # Development-only DB
   docker-compose up mysql
   go run cmd/api/main.go
   ```

2. **Adding New Features**:

   - Add domain types to `internal/domain/`
   - Implement repository interface in `internal/repository/`
   - Add service logic in `internal/service/`
   - Create handler in `internal/handler/`

3. **Database Changes**:
   - Add migrations to `migrations/`
   - Follow naming: `NNN_description.sql`

## Testing Guidelines

1. **Integration Tests**:

   - Require running MySQL instance
   - Set test environment in `.env.test`

2. **Unit Tests**:
   - Mock interfaces from `domain` package
   - Use table-driven tests

## Key Integration Points

1. **Database**:

   - MySQL 8.0 with Asia/Kolkata timezone
   - See `mysql/config/my.cnf` for configuration
   - Automated schema creation

2. **Monitoring**:

   - Prometheus metrics at `/metrics`
   - Grafana dashboards in `grafana/`
   - Health check at `/health`

3. **Load Balancing**:
   - NGINX handles request distribution
   - Least connections algorithm
   - Sticky sessions disabled

## Common Tasks

1. **Adding Prometheus Metrics**:

   - Define metric in handler
   - Add to `prometheus/prometheus.yml`

2. **Scaling Services**:

   ```yaml
   # In docker-compose.yml
   deploy:
     replicas: N
     resources:
       limits:
         cpus: 'X'
         memory: YM
   ```

3. **Updating Worker Pool**:
   - Environment variables:
     - WORKER_POOL_SIZE (default: 50)
     - MAX_JOB_QUEUE_SIZE (default: 100)

## Production Considerations

1. **Error Handling**:

   - All errors are logged with trace ID
   - Client errors return 4xx
   - System errors return 5xx

2. **Rate Limiting**:

   - Implemented in NGINX
   - See `nginx/nginx.conf`

3. **Observability**:
   - Structured logging format
   - Trace IDs in all logs
   - Request timing metrics

## API Documentation

See README.md for detailed API documentation, including PowerShell and curl examples for all CRUD operations.
