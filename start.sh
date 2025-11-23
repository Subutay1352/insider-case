#!/bin/bash

set -e

echo "Building application..."
make build || exit 1

echo "Building Docker image..."
make docker || exit 1

echo "Creating Docker network..."
docker network create insider-case-network 2>/dev/null || echo "Network already exists"

echo "Starting PostgreSQL container..."
docker stop insider-case-postgres 2>/dev/null || true
docker rm insider-case-postgres 2>/dev/null || true
docker run -d \
  --name insider-case-postgres \
  --network insider-case-network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=insider_case \
  -p 5432:5432 \
  postgres:15-alpine

echo "Starting Redis container..."
docker stop insider-case-redis 2>/dev/null || true
docker rm insider-case-redis 2>/dev/null || true
docker run -d \
  --name insider-case-redis \
  --network insider-case-network \
  -p 6379:6379 \
  redis:7-alpine

echo "Waiting for PostgreSQL to be ready..."
until docker exec insider-case-postgres pg_isready -U postgres > /dev/null 2>&1; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

echo "Waiting for Redis to be ready..."
until docker exec insider-case-redis redis-cli ping > /dev/null 2>&1; do
  echo "Waiting for Redis..."
  sleep 1
done

echo "Services are ready!"

echo "Stopping existing application container (if any)..."
docker stop insider-case 2>/dev/null || true
docker rm insider-case 2>/dev/null || true

echo "Starting application container..."
docker run -d \
  --name insider-case \
  --network insider-case-network \
  -p 8080:8080 \
  -e ENV=local \
  -e DB_TYPE=postgres \
  -e DB_HOST=insider-case-postgres \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  -e DB_NAME=insider_case \
  -e REDIS_HOST=insider-case-redis \
  -e REDIS_PORT=6379 \
  -e WEBHOOK_URL=${WEBHOOK_URL:-https://webhook.site/your-unique-id} \
  -e WEBHOOK_AUTH_KEY=${WEBHOOK_AUTH_KEY:-your-secret-key} \
  -e ACCESS_TOKEN=${ACCESS_TOKEN:-your-access-token} \
  -e SCHEDULER_INTERVAL=${SCHEDULER_INTERVAL:-20s} \
  -e SCHEDULER_AUTO_START=${SCHEDULER_AUTO_START:-true} \
  insider-case || exit 1

echo "Waiting for application to start..."
sleep 3

if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "Application is ready!"
    echo "Health: http://localhost:8080/health"
    echo "Swagger: http://localhost:8080/swagger/index.html"
else
    echo "Health check failed. Check logs: docker logs insider-case"
fi

