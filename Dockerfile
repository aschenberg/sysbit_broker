# Stage 1: Build the Go application
FROM golang:1.22-alpine AS build

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.* ./
RUN go mod download

# Copy the source code and .env file
COPY . .

# Build the Go application
RUN GOOS=linux go build -o sysbitbroker main.go

# Stage 2: Create a minimal image for running the application
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary and .env file from the build stage
COPY --from=build /app/sysbitbroker .
COPY --from=build /app/.env .

# Set the entrypoint command
ENTRYPOINT ["./sysbitbroker"]
