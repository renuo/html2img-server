# Build stage
FROM golang:1.19 AS builder

WORKDIR /app
COPY go.mod ./
COPY . .

RUN go build -o main .

# Final stage
FROM debian:bullseye-slim

# Install Chromium and required dependencies
RUN apt-get update && apt-get install -y \
    chromium \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/main .

# Set environment variables
ENV APP_PORT=:8080
ENV CHROME_BIN=/usr/bin/chromium

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"] 