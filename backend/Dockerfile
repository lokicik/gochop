# Build stage
FROM golang:1.22.1-alpine AS builder

# Install git and ca-certificates (needed for go modules with private repos and SSL)
RUN apk add --no-cache git ca-certificates tzdata

# Create and change to the app directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/gochop-server

# Final stage
FROM alpine:3.19

# Install ca-certificates for SSL connections
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1001 -S app && \
    adduser -S app -u 1001

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

# Change ownership to app user
RUN chown -R app:app /root/

# Switch to non-root user
USER app

# Expose port
EXPOSE 3001

# Command to run
CMD ["./main"] 