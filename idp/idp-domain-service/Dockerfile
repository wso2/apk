# Use a minimal Go image
FROM golang:1.23.4 as builder

# Set working directory
WORKDIR /app

# Copy go modules and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o idp ./cmd

# Use a small base image
FROM gcr.io/distroless/static

COPY --from=builder /app/idp /idp

EXPOSE 8080

ENTRYPOINT ["/idp"]
