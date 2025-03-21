# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod  ./

# Download the dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o blockchain-client .

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/blockchain-client .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./blockchain-client"]