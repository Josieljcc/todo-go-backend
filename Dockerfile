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

# Install swag for Swagger documentation (use specific version to avoid compatibility issues)
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

# Generate Swagger docs
RUN swag init -g cmd/api/main.go

# Fix Swagger docs compatibility issue (remove unsupported LeftDelim and RightDelim fields)
# Using grep to filter out problematic lines (more reliable than sed in Alpine)
RUN grep -v "LeftDelim:" docs/docs.go > docs/docs.go.tmp && \
    grep -v "RightDelim:" docs/docs.go.tmp > docs/docs.go && \
    rm -f docs/docs.go.tmp || true

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api"]

