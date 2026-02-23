# Build stage
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
# Assuming go.sum exists. If not, remove go.sum from this line.
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o midtrans-gateway ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests (useful for webhooks)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/midtrans-gateway .

# Expose port
EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=40s --retries=5 \
  CMD curl -f http://localhost:8000/up || exit 1

# Run the binary
CMD ["./midtrans-gateway"]