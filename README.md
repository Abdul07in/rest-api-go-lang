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

````markdown
# Docker Setup Instructions for Student API Project

## Prerequisites

- Docker installed
- Docker Swarm initialized
- Git repository cloned

## 1. Initialize Docker Swarm (if not already done)

```bash
# If you have multiple network interfaces, specify the advertise address
docker swarm init --advertise-addr <your-ip-address>

# If you have only one network interface, simply use:
docker swarm init
```
````

## 2. Build the API Image

```bash
# Navigate to project directory
cd rest-api-go-lang

# Build the API image
docker build -t localhost/student-api:latest .
```

## 3. Create Required Networks

```bash
# Create overlay network for services
docker network create --driver overlay student-api-network
```

## 4. Prepare Configuration Files

```bash
# Create necessary directories
mkdir -p mysql/config prometheus grafana nginx

# Ensure configuration files exist:
# - mysql/config/my.cnf
# - prometheus/prometheus.yml
# - nginx/nginx.conf
# These should already be in your repository
```

## 5. Deploy the Stack

```bash
# Deploy all services
docker stack deploy -c docker-compose.yml student-api

# Verify services are running
docker service ls

# Check specific service status
docker service ps student-api_api
docker service ps student-api_mysql
docker service ps student-api_nginx
docker service ps student-api_prometheus
docker service ps student-api_grafana
```

## 6. Monitor Deployment

```bash
# Watch service logs
docker service logs -f student-api_api

# Check MySQL logs
docker service logs -f student-api_mysql

# Check NGINX logs
docker service logs -f student-api_nginx
```

## 7. Access Services

- API: http://localhost:80/api/students
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (default credentials: admin/admin)

## 8. Scale Services

```bash
# Scale API service
docker service scale student-api_api=5

# Scale NGINX
docker service scale student-api_nginx=3
```

## 9. Useful Management Commands

```bash
# List all running services
docker service ls

# Check service details
docker service inspect student-api_api

# Check service logs
docker service logs student-api_api

# Update a service
docker service update --image localhost/student-api:new-version student-api_api

# Remove the entire stack
docker stack rm student-api

# List all stacks
docker stack ls

# List stack services
docker stack services student-api

# View stack tasks
docker stack ps student-api
```

## 10. Troubleshooting Commands

```bash
# Check service tasks (containers)
docker service ps --no-trunc student-api_api

# View container logs
docker service logs -f student-api_api

# Inspect service configuration
docker service inspect student-api_api

# Check network connectivity
docker network inspect student-api_student-api-network

# View resource usage
docker stats
```

## 11. Health Checks

```bash
# API health check
curl http://localhost/health

# MySQL health check
docker exec $(docker ps -q -f name=student-api_mysql) mysqladmin ping -h localhost -u root -proot

# Check all service health status
docker stack ps student-api
```

## 12. Cleanup

```bash
# Remove the stack
docker stack rm student-api

# Remove the network
docker network rm student-api_student-api-network

# Clean up unused resources
docker system prune -f

# Remove volumes (caution: this deletes data)
docker volume prune -f
```

## Environment Variables

Make sure these environment variables are set correctly in docker-compose.yml:

```yaml
environment:
  - DB_HOST=mysql
  - DB_PORT=3306
  - DB_USER=root
  - DB_PASSWORD=root
  - DB_NAME=student_db
  - SERVER_PORT=8080
  - WORKER_POOL_SIZE=50
  - MAX_JOB_QUEUE_SIZE=100
```

## Service URLs and Ports

- API: http://localhost:80
- MySQL: localhost:3306
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000

## Common Issues and Solutions

1. If services fail to start:

   ```bash
   # Check service logs
   docker service logs student-api_api
   ```

2. If MySQL connection fails:

   ```bash
   # Check if MySQL is healthy
   docker service ps student-api_mysql
   ```

3. If network issues occur:

   ```bash
   # Recreate the network
   docker network rm student-api_student-api-network
   docker network create --driver overlay student-api-network
   ```

4. If you need to rebuild and redeploy:

   ```bash
   # Rebuild image
   docker build -t localhost/student-api:latest .

   # Update service
   docker service update --force student-api_api
   ```

Remember to always check the logs when troubleshooting issues. The stack is designed to be resilient with health checks and automatic restarts, but monitoring the logs will help identify any persistent problems.
