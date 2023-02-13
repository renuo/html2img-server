# API for Generating Screenshots from HTML

[![Go](https://github.com/renuo/html2img-server/actions/workflows/go.yml/badge.svg)](https://github.com/renuo/html2img-server/actions/workflows/go.yml)

This is a simple API that takes an HTML file as input, generates a screenshot using Google Chrome, and returns the screenshot as a PNG image.

## Requirements

- Google Chrome
- Go

## Usage

Clone the repository
```bash
git clone git@github.com:renuo/html2img-server.git
cd html2img-server
```

Set the API token as an environment variable
```bash
export API_TOKEN=secret-token
```

Set the APP_PORT as an environment variable, for example `localhost:3001`
```bash
export APP_PORT=:3001
```

Set the path to Google Chrome as an environment variable
```bash
export CHROME_BIN=/usr/bin/google-chrome

# macOS
export CHROME_BIN=/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome
```

Build and run the API
```bash
go build
./html2img-server
```

Make a request to the API using curl or any other HTTP client
```bash
curl -X POST -d @sample.html http://localhost:3001/?token=secret-token --output screenshot.png 
```
Replace sample.html with the path to your HTML file, and secret-token with your API token.

