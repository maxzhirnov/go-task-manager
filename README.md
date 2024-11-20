# Go Task Manager

A robust task management system built with Go, demonstrating modern web development practices and Docker containerization. This project serves as a practical example of building a full-stack application using Go and PostgreSQL.

## Features

- RESTful API for task management
- PostgreSQL database integration
- Docker and Docker Compose setup for easy deployment
- Simple and responsive web interface
- Database initialization with sample data
- Containerized development environment

## Tech Stack

- **Backend**: Go 1.20
- **Database**: PostgreSQL 13
- **Containerization**: Docker & Docker Compose
- **API**: RESTful with Gorilla Mux
- **Frontend**: HTML, CSS, JavaScript

## Prerequisites

- Docker and Docker Compose
- Go 1.20 (for local development)
- PostgreSQL (for local development)

## Quick Start

1. Clone the repository:

```bash
git clone https://github.com/maxzhirnov/go-task-manager.git
cd go-task-manager
```

2. Start the application using Docker Compose:

```bash
docker-compose up --build
```

1. Access the application:

- Web Interface: <http://localhost:8080>
- API Endpoint: <http://localhost:8080/api/tasks>

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/tasks` | Retrieve all tasks |
| GET | `/api/tasks/{id}` | Retrieve a specific task |
| POST | `/api/tasks` | Create a new task |
| PUT | `/api/tasks/{id}` | Update an existing task |
| DELETE | `/api/tasks/{id}` | Delete a task |

### Example API Request

```bash
# Create a new task
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"New Task","description":"Task description","status":"pending"}'
```

## Project Structure

```
go-task-manager/
├── cmd/
│   └── api/            # Application entrypoint
├── internal/
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   └── models/         # Data models
├── pkg/
│   └── database/       # Database utilities
├── scripts/
│   └── db/            # Database initialization scripts
├── web/               # Frontend assets
└── docker-compose.yml # Docker composition
```

## Local Development

1. Set up the environment variables:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=mysecretpassword
export DB_NAME=taskmanager
```

2. Run the database:

```bash
docker-compose up db
```

3. Run the application:

```bash
go run cmd/api/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux) for HTTP routing
- [PostgreSQL](https://www.postgresql.org/) for database
- [Docker](https://www.docker.com/) for containerization