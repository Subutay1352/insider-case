FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest


RUN apk --no-cache add ca-certificates curl && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser


WORKDIR /app

COPY --from=builder --chmod=755 /app/main /app/main
COPY --from=builder /app/migrations /app/migrations

# Change ownership of /app directory and binary
RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["/app/main"]

