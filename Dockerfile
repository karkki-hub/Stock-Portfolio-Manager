# STAGE 1: Build the binary
FROM golang:1.26-alpine AS builder

# Set the working directory
WORKDIR /app

RUN apk add --no-cache git
# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build the application
COPY . .
RUN go build -o app ./cmd/server

# STAGE 2: Create the minimal runtime image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose the application port (e.g., 8080)
EXPOSE 8080

# Command to run the application
CMD ["./app"]
