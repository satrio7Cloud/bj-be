# Build Stage
FROM golang:1.25-alpine AS builder

# Set directory
WORKDIR /app

# Install dependencies (if any CGO needed, though alpine usually needs apk add gcc etc)
# CGO_ENABLED=0 creates a statically linked binary which is smaller and runs without C libraries.
ENV CGO_ENABLED=0 
ENV GOOS=linux

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application (binary will be named 'server')
RUN go build -o server ./cmd/server

# Final Stage (Minimal Image)
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
# Copy config files if needed (adjust the paths based on how you load .env or config)
# COPY --from=builder /app/.env .

# Set timezone if necessary
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

# Expose port (adjust to your config)
EXPOSE 3000

# Run the binary
CMD ["./server"]
