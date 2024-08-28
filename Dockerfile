# Use the official Golang image as the base image
FROM golang:1.22.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./go.mod ./go.sum ./

# Copy the source from the current directory to the Working Directory inside the container
COPY ./hikhello ./hikhello
COPY ./hikaxprogo ./hikaxprogo
COPY ./templates ./templates

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main  ./hikhello/

# Start a new stage from scratch
FROM scratch

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates


# Expose port 8080 to the outside world
EXPOSE 8080

ENTRYPOINT ["/main"]
