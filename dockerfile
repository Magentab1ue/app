# Builder stage
FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./app/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/

# Copy necessary directories
# COPY assets /app/assets
COPY configs /app/configs
COPY modules /app/modules
COPY pkg /app/pkg

# Assuming .env file is at the root of your project
# Copy the .env file for environment variables
COPY config.env /app/

CMD ["/app/main"]

