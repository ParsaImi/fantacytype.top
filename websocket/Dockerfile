# Start from the official Golang image
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and necessary build tools
RUN apk add --no-cache git gcc musl-dev

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o websocket-app ./cmd/main.go

# Create a minimal production image
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/websocket-app .

# Expose the application port
EXPOSE 9000

# Run the binary
CMD ["./websocket-app"]
