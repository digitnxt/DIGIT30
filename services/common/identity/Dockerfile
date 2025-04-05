# Use Go base image for building
FROM golang:1.24-alpine AS builder

# Install Git (required for some swag operations)
RUN apk update && apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Run go mod tidy to ensure dependencies are up to date
RUN go mod tidy

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate the Swagger docs (ensuring main.go is used as the entry point for annotations)
RUN swag init --generalInfo main.go --output ./docs

# Build the Go application binary
RUN go build -o server .

# Final runtime image based on distroless for a minimal footprint
FROM gcr.io/distroless/base-debian11

# Set working directory in the final image
WORKDIR /app

# Copy built binary and generated docs from the builder stage
COPY --from=builder /app/server /app/
COPY --from=builder /app/docs /app/docs

# Expose the service port
EXPOSE 8080

# Set the entrypoint to run the server binary
ENTRYPOINT ["/app/server"]