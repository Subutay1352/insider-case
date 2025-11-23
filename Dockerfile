FROM golang:1.21-alpine AS builder

# SQLite dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# SQLite dependencies
RUN CGO_ENABLED=1 CGO_CFLAGS="-D_LARGEFILE64_SOURCE" GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest

# SQLite dependencies
RUN apk --no-cache add ca-certificates sqlite curl

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Create app directory
WORKDIR /app

# Copy the binary from builder with execute permission
COPY --from=builder --chmod=755 /app/main /app/main

# Change ownership of /app directory and binary
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["/app/main"]

