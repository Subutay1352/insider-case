# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies for CGO (required for SQLite, optional for PostgreSQL)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 for database independence (supports both PostgreSQL and SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Final stage
FROM alpine:latest

# Install runtime dependencies (SQLite support + health check)
RUN apk --no-cache add ca-certificates sqlite curl

WORKDIR /root/

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Copy the binary from builder
COPY --from=builder /app/main .

# Change ownership
RUN chown appuser:appuser /root/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]

