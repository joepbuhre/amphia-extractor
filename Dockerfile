# Stage 1: Build the Go application
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all the dependencies. This is a separate step so that Docker caches it.
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o amphia-extractor

# Stage 2: Run the Go application
FROM alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/amphia-extractor .

# Expose the application port
EXPOSE 8080

RUN chmod +x /app/amphia-extractor

ENV GIN_MODE=release

# Run the binary program produced by the build step
CMD ["/app/amphia-extractor"]
