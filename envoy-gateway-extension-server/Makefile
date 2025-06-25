help:
	@echo "  build         - Build the extension server binary for Linux (CGO disabled)"
	@echo "  build-local   - Build the extension server binary for local machine (macOS)"
	@echo "  image         - Build container image"
	@echo "  help          - Display this help message"

IMAGE_NAME ?= envoy-gateway-extension-server:latest

build:
	@echo "=== Building extension-server binary for Linux ==="
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/extension-server ./cmd

build-local:
	@echo "=== Building extension-server binary for macOS ==="
	go build -o bin/extension-server-local ./cmd


image:
	@echo "=== Building Docker image $(IMAGE_NAME) with Linux binary ==="
	docker build --platform linux/amd64 -t $(IMAGE_NAME) -f Dockerfile . --no-cache

image-local:
	@echo "=== Building Docker image $(IMAGE_NAME) with local OS binary ==="
	docker build -t $(IMAGE_NAME) -f Dockerfile . --no-cache