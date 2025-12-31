# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Install swag-openapi3 for OpenAPI 3.0 documentation generation
RUN go install github.com/dimasdanz/swag-openapi3@latest

# Generate OpenAPI 3.0 docs
RUN swag-openapi3 init -g cmd/api/main.go -o ./docs --requiredByDefault

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/api .

# Copy swagger documentation files
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 3002

# Run the application
CMD ["./api"]

