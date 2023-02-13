package main

import (
	"bytes"
	"image"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	html := []byte(`<html>
		<head><style>body {background-color: yellow;width: 600px;height: 310px;}</style></head>
		<body><h1>Hello from Go!</h1></body>
		</html>`)

	apiToken := os.Getenv("API_TOKEN")
	req, err := http.NewRequest("POST", "http://localhost:8080/?token="+apiToken, bytes.NewBuffer(html))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	screenshot, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}

	if len(screenshot) == 0 {
		t.Errorf("handler returned an empty response")
	}

	img, _, err := image.Decode(bytes.NewBuffer(screenshot))
	if err != nil {
		t.Fatalf("Error decoding screenshot: %v", err)
	}

	zoomFactor := 2
	expectedWidth := 600 * zoomFactor
	expectedHeight := 310 * zoomFactor
	if img.Bounds().Dx() != expectedWidth || img.Bounds().Dy() != expectedHeight {
		t.Errorf("Unexpected screenshot dimensions: got %dx%d, want %dx%d", img.Bounds().Dx(), img.Bounds().Dy(), expectedWidth, expectedHeight)
	}

	if err := ioutil.WriteFile("result.png", screenshot, 0644); err != nil {
		t.Fatalf("Error writing result to file: %v", err)
	}
}

func TestUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/?token=invalid", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestEmptyRequest(t *testing.T) {
	apiToken := os.Getenv("API_TOKEN")
	req, _ := http.NewRequest("POST", "http://localhost:8080/?token="+apiToken, nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
