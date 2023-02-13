package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var CHROME_ARGS = []string{
	"--headless",
	"--hide-scrollbars",
	"--window-size=600,310",
	"--disable-extensions",
	"--no-sandbox",
	"--disable-background-networking",
	"--disable-cache",
	"--force-device-scale-factor=2",
}

func takeScreenshot(html []byte) ([]byte, error) {
	htmlFile, err := ioutil.TempFile("", "*.html")
	if err != nil {
		return nil, err
	}
	defer os.Remove(htmlFile.Name())

	_, err = htmlFile.Write(html)
	if err != nil {
		return nil, err
	}

	screenshotFile, err := ioutil.TempFile("", "*.png")
	if err != nil {
		return nil, err
	}
	defer os.Remove(screenshotFile.Name())

	params := append(CHROME_ARGS, fmt.Sprintf("file://%s", htmlFile.Name()), fmt.Sprintf("--screenshot=%s", screenshotFile.Name()))
	cmd := exec.Command(chromeExecutable(), params...)
	log.Println("Running command: ", cmd.String())

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)
	log.Printf("Time taken: %s", elapsed)

	if err != nil {
		return nil, err
	}

	screenshot, err := ioutil.ReadFile(screenshotFile.Name())
	if err != nil {
		return nil, err
	}

	return screenshot, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	apiToken := os.Getenv("API_TOKEN")
	if token != apiToken {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Body == nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	result, err := takeScreenshot(html)
	if err != nil {
		http.Error(w, "Error while taking screenshot", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func chromeExecutable() string {
	chrome := os.Getenv("CHROME_BIN")
	_, err := exec.LookPath(chrome)
	if err == nil {
		return chrome
	}

	log.Fatal("Google Chrome not found")
	return ""
}

func main() {
	// check if chrome is installed
	chromeExecutable()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
