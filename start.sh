#!/bin/bash

set -e

echo "Building application..."
make build || exit 1

echo "Building Docker image..."
make docker || exit 1

echo "Stopping existing container (if any)..."
docker stop insider-case 2>/dev/null || true
docker rm insider-case 2>/dev/null || true

echo "Starting Docker container..."
docker run -d \
  --name insider-case \
  -p 8080:8080 \
  -e ENV=local \
  -e DB_TYPE=sqlite \
  -e DB_PATH=/app/data.db \
  -e WEBHOOK_URL=${WEBHOOK_URL:-https://webhook.site/your-unique-id} \
  -e WEBHOOK_AUTH_KEY=${WEBHOOK_AUTH_KEY:-your-secret-key} \
  -e ACCESS_TOKEN=${ACCESS_TOKEN:-your-access-token} \
  -e SCHEDULER_INTERVAL=${SCHEDULER_INTERVAL:-2m} \
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

