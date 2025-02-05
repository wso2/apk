# Stage 1: Build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and source code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application statically linked
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Use Alpine as the base image
FROM alpine:latest

# Install minimal dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the statically built binary
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Run the app
CMD ["./main"]
