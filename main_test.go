package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandler(t *testing.T) {
	html := []byte(`<html><body><h1>Hello, World!</h1></body></html>`)
	htmlFile, err := ioutil.TempFile("", "*.html")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(htmlFile.Name())
	if _, err := htmlFile.Write(html); err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}

	apiToken := os.Getenv("API_TOKEN")
	req, err := http.NewRequest("POST", "http://localhost:8080/?token="+apiToken, bytes.NewBuffer(html))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "text/html")

	rr := httptest.NewRecorder()

	handler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	result, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if len(result) == 0 {
		t.Errorf("handler returned an empty response")
	}

	if err := ioutil.WriteFile("result.png", result, 0644); err != nil {
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
