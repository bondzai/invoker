# Stage 1: Build the application
FROM golang:1.21 AS builder

WORKDIR /go/src/app

# Copy only the necessary files for the Go modules
COPY go.mod .
COPY go.sum .

# Download and install Go dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o invoker ./cmd

# Stage 2: Create a minimal production image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /go/src/app/invoker /app/invoker

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./invoker"]