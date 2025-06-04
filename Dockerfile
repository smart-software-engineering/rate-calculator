# Multi-stage build for Go + Svelte Rate Calculator

# Stage 1: Build Svelte frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files first for better caching
COPY frontend/package*.json ./
RUN npm ci

# Copy frontend source and build
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend from previous stage
COPY --from=frontend-builder /app/internal/server/static ./internal/server/static

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Stage 3: Create minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=backend-builder /app/main .

# Copy any static assets and templates that might be needed
COPY --from=backend-builder /app/internal/template ./internal/template
COPY --from=backend-builder /app/internal/server/static ./internal/server/static
COPY --from=backend-builder /app/data ./data

# Expose port 8080
EXPOSE 8080

# Command to run
CMD ["./main"] 