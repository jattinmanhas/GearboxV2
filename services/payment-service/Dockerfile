# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go work files
COPY go.work go.work.sum ./
COPY services/shared ./services/shared
COPY services/payment-service ./services/payment-service

# Download dependencies
WORKDIR /build/services/payment-service
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/services/payment-service/api .

# Copy migrations
COPY --from=builder /build/services/payment-service/migrations ./migrations

# Expose port
EXPOSE 8083

# Run the application
CMD ["./api"]

