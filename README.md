# Student API

A high-performance, scalable RESTful API for managing student data with full CRUD operations. Features async processing, comprehensive monitoring, and distributed tracing.

## Technologies Used

- Go 1.24
- MySQL 8.0
- Docker & Docker Compose
- NGINX Load Balancer
- Prometheus & Grafana
- GitHub Actions for CI/CD

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── domain/
│   ├── handler/
│   ├── logging/
│   ├── repository/
│   └── service/
├── migrations/
├── .env
├── docker-compose.yml
└── Dockerfile
```

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- MySQL 8.0
- Prometheus (optional, for metrics)
- Grafana (optional, for dashboards)

## Local Development

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd rest-api-go
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up environment variables in `.env` file:

   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=student_db
   SERVER_PORT=8080
   WORKER_POOL_SIZE=50      # Number of worker goroutines
   MAX_JOB_QUEUE_SIZE=100   # Size of job queue buffer
   ```

4. Run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

## API Endpoints and Usage Examples

### Create a New Student (POST)

```bash
curl -X POST http://localhost:8080/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "age": 20,
    "grade": 85.5
  }'
```

Response:

```json
{
  "id": 1,
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "age": 20,
  "grade": 85.5,
  "createdAt": "2025-08-20T16:45:00Z",
  "updatedAt": "2025-08-20T16:45:00Z"
}
```

### Get All Students (GET)

```bash
curl http://localhost:8080/api/students
```

Response:

```json
[
  {
    "id": 1,
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "age": 20,
    "grade": 85.5,
    "createdAt": "2025-08-20T16:45:00Z",
    "updatedAt": "2025-08-20T16:45:00Z"
  },
  {
    "id": 2,
    "firstName": "Jane",
    "lastName": "Smith",
    "email": "jane.smith@example.com",
    "age": 22,
    "grade": 92.0,
    "createdAt": "2025-08-20T16:46:00Z",
    "updatedAt": "2025-08-20T16:46:00Z"
  }
]
```

### Get Student by ID (GET)

```bash
curl http://localhost:8080/api/students/1
```

Response:

```json
{
  "id": 1,
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "age": 20,
  "grade": 85.5,
  "createdAt": "2025-08-20T16:45:00Z",
  "updatedAt": "2025-08-20T16:45:00Z"
}
```

### Update Student (PUT)

```bash
curl -X PUT http://localhost:8080/api/students/1 \
  -H "Content-Type: application/json" \
  -d '{
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@example.com",
    "age": 21,
    "grade": 88.5
  }'
```

Response:

```json
{
  "id": 1,
  "firstName": "John",
  "lastName": "Doe",
  "email": "john.doe@example.com",
  "age": 21,
  "grade": 88.5,
  "createdAt": "2025-08-20T16:45:00Z",
  "updatedAt": "2025-08-20T16:50:00Z"
}
```

### Delete Student (DELETE)

```bash
curl -X DELETE http://localhost:8080/api/students/1
```

Response: Empty with status code 204 (No Content)

### PowerShell Examples

For Windows PowerShell users, here are the equivalent commands:

Create Student:

```powershell
$body = @{
    firstName = "John"
    lastName = "Doe"
    email = "john.doe@example.com"
    age = 20
    grade = 85.5
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "http://localhost:8080/api/students" `
  -ContentType "application/json" -Body $body
```

Update Student:

```powershell
$body = @{
    firstName = "John"
    lastName = "Doe"
    email = "john.doe@example.com"
    age = 21
    grade = 88.5
} | ConvertTo-Json

Invoke-RestMethod -Method Put -Uri "http://localhost:8080/api/students/1" `
  -ContentType "application/json" -Body $body
```

Get All Students:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/students"
```

Get Student by ID:

```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/api/students/1"
```

Delete Student:

```powershell
Invoke-RestMethod -Method Delete -Uri "http://localhost:8080/api/students/1"
```

## Deployment

The application can be deployed using Docker. The CI/CD pipeline will:

1. Run tests
2. Build Docker image
3. Push to Docker Hub
4. Deploy to production server

### Required Secrets for GitHub Actions

- `DOCKER_USERNAME`: Docker Hub username
- `DOCKER_PASSWORD`: Docker Hub password
- `DEPLOY_HOST`: Production server hostname
- `DEPLOY_USERNAME`: SSH username for production server
- `DEPLOY_SSH_KEY`: SSH private key for production server

## Architecture Features

### Async Processing

The API uses a worker pool pattern for handling requests:

- Configurable number of worker goroutines (default: 50)
- Job queue with buffer (default: 100)
- Context-aware operations with timeouts
- Graceful error handling and recovery

### Load Balancing

NGINX load balancer provides:

- Request distribution across API instances
- Health checks
- Connection limiting
- WebSocket support
- Custom timeout configurations

### Monitoring & Metrics

1. **Prometheus Integration**:

   - HTTP request metrics
   - Database connection stats
   - Worker pool utilization
   - Custom business metrics

2. **Grafana Dashboards**:
   - Real-time request monitoring
   - Error rate tracking
   - Performance metrics
   - System resource usage

## Logging and Tracing

The application includes comprehensive request logging with:

- Unique trace ID for each request
- Request/response details
- Operation logging
- Client IP tracking
- Response times
- Worker pool metrics

Each log entry includes a trace ID that can be used to follow a request through all its operations, including async processing.

## Infrastructure

The application is designed to be highly available and scalable:

1. **API Servers**:

   - Multiple replicas for high availability
   - Resource limits and reservations
   - Health checks
   - Rolling updates

2. **Database**:

   - MySQL 8.0 with proper timezone config
   - Automated schema creation
   - Connection pooling
   - Index optimization

3. **Load Balancer**:
   - Least connections algorithm
   - Health monitoring
   - Timeout configurations
   - SSL termination (optional)

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
