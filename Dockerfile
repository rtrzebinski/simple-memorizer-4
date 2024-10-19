# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Copy source code
COPY ./ ./

# Install dependencies required for building
RUN apk add --update --no-cache git make

# Build the application
RUN make build

# Install the migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy files
COPY --from=builder /app/bin/simple-memorizer ./
COPY --from=builder /app/web ./web
COPY --from=builder /go/bin/migrate ./migrate
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/version ./version
COPY --from=builder /app/ssl/localhost-cert.pem ./ssl/localhost-cert.pem
COPY --from=builder /app/ssl/localhost-key.pem ./ssl/localhost-key.pem

# Ensure migrate is executable
RUN chmod +x ./migrate

# Add /app to PATH
ENV PATH="/app:${PATH}"

# Expose port 8000
EXPOSE 8000

# Start the application
CMD ["./simple-memorizer"]
