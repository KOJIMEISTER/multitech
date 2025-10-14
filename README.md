[![CI](https://github.com/kojimeister/multitech/actions/workflows/main.yml/badge.svg)](https://github.com/kojimeister/multitech/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/kojimeister/multitech/branch/main/graph/badge.svg)](https://codecov.io/gh/kojimeister/multitech)

# Multitech - Go REST API Tutorial Project

## Introduction

A tutorial project demonstrating a REST API implementation in Go with:

- Gin web framework
- JWT authentication with Redis storage
- PostgreSQL user storage
- Swagger documentation
- Comprehensive testing (unit & integration)
- Docker-compose setup
- CI/CD with GitHub Actions

## Features

- JWT-based authentication with Redis session storage
- User management with PostgreSQL
- Swagger API documentation
- Healthcheck endpoint
- Unit and integration tests
- Docker-compose for local development
- GitHub CI/CD pipeline with test coverage

## Quick Start

1. Clone the repository
2. Create `.env` file (see Environment Variables section)
3. Run: `docker-compose up`
4. Access API documentation at: `http://localhost:8080/swagger/index.html`

Example API commands:

```bash
# Register new user
curl -X POST "http://localhost:8080/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass","email":"test@example.com"}'

# Login
curl -X POST "http://localhost:8080/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'

# Access protected endpoint
curl -X GET "http://localhost:8080/protected" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Environment Variables

Required `.env` variables:

- `JWT_SECRET`: Secret key for JWT token signing
- `REDIS_URL`: Redis connection URL (e.g. `redis://redis:6379`)
- `POSTGRES_USER`: PostgreSQL username
- `POSTGRES_PASSWORD`: PostgreSQL password
- `POSTGRES_DB`: PostgreSQL database name

Example `.env` file:

```
JWT_SECRET=your-secret-key-here
REDIS_URL=redis://redis:6379
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=multitech
```

## Testing

Run tests:

```bash
# Unit tests
go test ./... -tags=unit

# Integration tests
go test ./... -tags=integration
```

## CI/CD

The project includes GitHub Actions workflows that:

1. Run unit tests with coverage
2. Run integration tests
3. Generate test coverage reports
4. Upload coverage to Codecov

## Technology Stack

- Go 1.23.4
- Gin web framework
- Gorm ORM
- Redis
- PostgreSQL
- Swagger
- Testcontainers for integration tests

## Credits

Created by Nikita @ Kojimeister
