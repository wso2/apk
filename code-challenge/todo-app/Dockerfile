# Build Stage
FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Deployment Stage
FROM alpine:3.19.1

WORKDIR /root/

COPY --from=build /app/main .

# Expose the application port
EXPOSE 8080

CMD ["./main"]
