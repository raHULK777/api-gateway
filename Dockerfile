# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api-gateway ./cmd/server

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install ca-certificates and tzdata
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=UTC

# Copy the binary from builder
COPY --from=builder /app/api-gateway .

# Create non-root user
RUN addgroup -g 1001 -S nonroot && \
    adduser -u 1001 -S nonroot -G nonroot && \
    chown -R nonroot:nonroot /app

USER nonroot

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["./api-gateway"]
