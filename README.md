# Insider Case - Message Scheduler Service

Message scheduler service with automatic webhook delivery, retry mechanism, and Redis caching.

## Installation

```bash
go mod download
```

## Usage

Run locally:
```bash
make run
```

Build:
```bash
make build
```

## Configuration

Create `.env` file for local development:

```bash
ENV=local
DB_TYPE=postgres                    # or sqlite
DB_PATH=insider_case.db          # for sqlite
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=insider_case
WEBHOOK_URL=https://webhook.site/your-id
WEBHOOK_AUTH_KEY=your-secret-key
REDIS_HOST=localhost
REDIS_PORT=6379
SCHEDULER_INTERVAL=2m
SCHEDULER_AUTO_START=true
ACCESS_TOKEN=your-access-token
```

## API

Health check:
```
GET /health
```

Swagger documentation:
```
GET /swagger/index.html
```

API v1 endpoints (requires `x-access-token` header):
```
POST /api/v1/sender/startScheduler
POST /api/v1/sender/stopScheduler
GET  /api/v1/sender/statusScheduler
GET  /api/v1/messages/sent?limit=10&offset=0
```

## Makefile

- `make build` - Build the application
- `make run` - Run locally
- `make test` - Run tests
- `make lint` - Run linter
- `make swagger` - Generate Swagger docs
- `make docker` - Build Docker image
- `make clean` - Clean build artifacts

## Docker

Build:
```bash
make docker
```

Run with SQLite:
```bash
docker run -d -p 8080:8080 \
  -e DB_TYPE=sqlite \
  -e DB_PATH=/app/data.db \
  insider-case
```

Run with PostgreSQL:
```bash
docker run -d -p 8080:8080 \
  -e DB_TYPE=postgres \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=insider_case \
  insider-case
```

## Tech Stack

- Go 1.21
- Gin
- GORM (PostgreSQL/SQLite)
- Redis (optional)
- Docker

## Code Quality

SonarQube Cloud: https://sonarcloud.io/project/overview?id=Subutay1352_insider-case