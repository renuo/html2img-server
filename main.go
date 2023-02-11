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

func runCommand(html []byte) ([]byte, error) {
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
	html, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	result, err := runCommand(html)
	if err != nil {
		http.Error(w, "Error running command", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(result)
	if err != nil {
		http.Error(w, "Error writing result", http.StatusInternalServerError)
		return
	}
}

func chromeExecutable() string {
	_, err := exec.LookPath("google-chrome-stable")
	if err == nil {
		return "google-chrome-stable"
	}

	log.Fatal("Google Chrome not found")
	return ""
}

func main() {
	// Pre-boot google-chrome-stable
	err := exec.Command(chromeExecutable(), "--headless", "--disable-gpu", "--no-sandbox", "--disable-dev-shm-usage").Run()
	if err != nil {
		log.Fatalf("Error pre-booting google-chrome-stable: %v", err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
