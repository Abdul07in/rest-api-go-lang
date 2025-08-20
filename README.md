# Student API

A RESTful API for managing student data with full CRUD operations.

## Technologies Used

- Go 1.24
- MySQL 8.0
- Docker & Docker Compose
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
   ```

4. Run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

## API Endpoints

- `POST /api/students` - Create a new student
- `GET /api/students` - Get all students
- `GET /api/students/{id}` - Get a student by ID
- `PUT /api/students/{id}` - Update a student
- `DELETE /api/students/{id}` - Delete a student

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

## Logging

The application includes comprehensive request logging with:

- Unique trace ID for each request
- Request/response details
- Operation logging
- Client IP tracking
- Response times

Each log entry includes a trace ID that can be used to follow a request through all its operations.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request
