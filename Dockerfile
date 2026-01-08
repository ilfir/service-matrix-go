
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies if needed
# RUN apk add --no-cache git

COPY go.mod ./
# COPY go.sum ./ 
# If you have dependencies, uncomment above and run:
# RUN go mod download

COPY . .

# Build the application
RUN go build -o server cmd/server/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy resources and data
COPY --from=builder /app/resources ./resources
COPY --from=builder /app/data ./data

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
