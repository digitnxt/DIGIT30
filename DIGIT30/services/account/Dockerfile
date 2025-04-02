# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Install git (required for fetching dependencies)
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker cache for dependency downloads.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the binary. CGO is disabled for a fully static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

# Final stage: create a minimal container image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .

# Expose port 8080.
EXPOSE 8080

# Command to run the application.
CMD ["./app"]
