# Dockerfile for the Go Backend

# --- Stage 1: Build Stage ---
# Use an official Go image as the builder environment
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go application for Linux.
# The -o flag places the output in the specified directory/file.
# CGO_ENABLED=0 is important for creating a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/engine ./cmd/server

# --- Stage 2: Final Stage ---
# Use a minimal base image for the final container to keep it small
FROM alpine:latest

# We need ca-certificates for any potential HTTPS calls in the future
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /go/bin/engine .

# Copy the configuration file
COPY config.json .

# Expose the port the backend runs on
EXPOSE 8080

# The command to run when the container starts
CMD ["./engine"]
