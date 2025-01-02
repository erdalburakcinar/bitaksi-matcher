.PHONY: build run test docker-build docker-run clean

# Application settings
BINARY_NAME=matcher-service
MAIN=cmd/main.go

# Docker settings
DOCKER_IMAGE=matcher-service
DOCKER_TAG=latest

# Default port
PORT=8081

# Build the Go application
build:
	@echo "Building application..."
	go build -o $(BINARY_NAME) $(MAIN)

# Run the application locally
run: build
	@echo "Running application..."
	./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Build the Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run the Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -d -p $(PORT):$(PORT) --name $(BINARY_NAME) $(DOCKER_IMAGE):$(DOCKER_TAG)

# Clean up the build and Docker artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	docker rm -f $(BINARY_NAME) || true
	docker rmi -f $(DOCKER_IMAGE):$(DOCKER_TAG) || true