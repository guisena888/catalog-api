# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.
   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
   - `api/`: HTTP response helpers
   - `database/`: Database connection
   - `handler/`: HTTP handlers
   - `middleware/`: HTTP middleware (logging, recovery)
   - `mocks/`: Generated mocks for testing
   - `repository/`: Data access layer
   - `service/`: Business logic layer

3. **pkg/model/**: Contains the domain models.
4. **errors/**: Contains application error types.
5. **sql/**: Contains a very simple database migration scripts setup.
6. `.env`: Environment variables file for configuration.

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run unit tests (skips integration tests).
  - `make test.integration`: Will run all tests including integration tests (requires docker).
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.
  - `make mocks`: Will generate mocks for testing.

## API Documentation

The API is documented using OpenAPI 3.0. See [openapi.yaml](openapi.yaml) for the full specification.

### View API Documentation

Start Swagger UI to explore the API interactively:

```bash
make swagger
```

Then open http://localhost:8081 in your browser.
