# Start from golang base image.
FROM golang:alpine as builder

# Add Maintainer info.
LABEL maintainer="Vladimir Andrianov"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container.
COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main srv/cmd/main.go

# Start a new stage from scratch.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create directory to place the configuration file.
RUN mkdir -p /app/srv/cmd/config

# Copy the Pre-built binary file from the previous stage and certificates and the configuration file.
COPY --from=builder /app/main .
COPY --from=builder /app/srv/cmd/config/configuration.yaml /app/srv/cmd/config
COPY --from=builder /app/srv/rsa.crt /app/srv
COPY --from=builder /app/srv/rsa.key /app/srv

# Command to run the executable
CMD ["./main"]