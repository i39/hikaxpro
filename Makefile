# Define variables
BINARY_NAME=hikaxpro
DOCKER_IMAGE_NAME=hikaxpro-image
DOCKER_TAG=latest

# Build the binary for Mac Os
.PHONY: build-binary
build-binary:
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME) ./hikhello/

# Build the binary for Linux
.PHONY: build-binary-linux
build-binary-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux ./hikhello/


# Build Docker image
.PHONY: build-docker
build-docker: build-binary-linux
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Clean up
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
