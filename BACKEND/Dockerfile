# syntax=docker/dockerfile:1
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
