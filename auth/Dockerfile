# Use an official Go runtime as the base image
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build

# Run the Go application when the container starts
CMD ["./auth"]
