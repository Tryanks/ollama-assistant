FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ollama-assistant .

# Second stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ollama-assistant .

# Expose the default port
EXPOSE 11434

# Run the application
CMD ["./ollama-assistant"]
