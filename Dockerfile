# Build stage
FROM golang:1.19 AS builder

WORKDIR /app
COPY go.mod ./
COPY . .

RUN go build -o main .

# Final stage
FROM debian:bullseye-slim

# Install Chromium, base fonts, and tools to fetch custom fonts
RUN apt-get update && apt-get install -y \
    chromium \
    ca-certificates \
    curl \
    unzip \
    fontconfig \
    fonts-liberation \
    fonts-dejavu \
    fonts-noto \
    fonts-noto-cjk \
    && rm -rf /var/lib/apt/lists/*

# Install Inter font (https://rsms.me/inter/)
RUN mkdir -p /usr/share/fonts/truetype/inter \
    && curl -sSL -o /tmp/inter.zip https://github.com/rsms/inter/releases/download/v4.0/Inter-4.0.zip \
    && unzip -j /tmp/inter.zip 'extras/ttf/*.ttf' -d /usr/share/fonts/truetype/inter \
    && rm /tmp/inter.zip

# Install Twitter Color Emoji SVGinOT font
RUN mkdir -p /usr/share/fonts/truetype/twemoji \
    && curl -sSL -o /usr/share/fonts/truetype/twemoji/TwitterColorEmoji-SVGinOT.ttf \
        https://github.com/eosrei/twemoji-color-font/releases/download/v14.0.2/TwitterColorEmoji-SVGinOT.ttf

RUN fc-cache -f -v

WORKDIR /app
COPY --from=builder /app/main .

# Set environment variables
ENV APP_PORT=:8080
ENV CHROME_BIN=/usr/bin/chromium

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]
