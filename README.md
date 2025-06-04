# API for Generating Screenshots from HTML

[![Go](https://github.com/renuo/html2img-server/actions/workflows/go.yml/badge.svg)](https://github.com/renuo/html2img-server/actions/workflows/go.yml)

This is a simple API that takes an HTML file as input, generates a screenshot using Google Chrome, and returns the screenshot as a PNG image.

## Requirements

- Docker

## Usage

Clone the repository
```bash
git clone git@github.com:renuo/html2img-server.git
```

Build the image
```bash
cd html2img-server
docker build -t html2img-server .
```

Run it and test it:

```bash
docker run -p 8080:8080 -e API_TOKEN=your-secret-token html2img-server
curl -X POST  -d @test.html "http://localhost:8080/?token=your-secret-token" --output screenshot.png
```


Build and run the API
```bash
go build
./html2img-server
```

Make a request to the API using curl or any other HTTP client
```bash
curl -X POST -d @test.html http://localhost:3001/?token=your-secret-token --output screenshot.png 
```

# systemd setup

To run the executable as a background service you can use a systemd configuration. It is based on the `html2img-server.service.template` file and needs to be placed in /etc/systemd/system/html2img-server.service with the correct token set.

After creating/changing the file you need to run `systemctl daemon-reload`

When setting it up on a new server you need to run `systemctl enable html2img-server.service` once.

Then you can manage the service like any other with these commands:

```
service html2img-server stop
service html2img-server start
service html2img-server status
```

If the service does crash it will log to `/var/log/syslog`
