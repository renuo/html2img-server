package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var CHROME_ARGS = []string{
	"--headless",
	"--hide-scrollbars",
	"--window-size=600,314",
	"--disable-extensions",
	"--no-sandbox",
	"--disable-background-networking",
	"--disable-cache",
	"--force-device-scale-factor=2",
}

// A well-formed screenshot is always several KB; anything under this is
// almost certainly an empty/truncated capture.
const minScreenshotBytes = 1024

func chromeTimeout() time.Duration {
	if raw := os.Getenv("CHROME_TIMEOUT_SECONDS"); raw != "" {
		if seconds, err := strconv.Atoi(raw); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return 10 * time.Second
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

	timeout := chromeTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, chromeExecutable(), params...)
	log.Println("Running command: ", cmd.String())

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)
	log.Printf("Time taken: %s", elapsed)

	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("chrome timed out after %s", timeout)
	}
	if err != nil {
		return nil, err
	}

	screenshot, err := ioutil.ReadFile(screenshotFile.Name())
	if err != nil {
		return nil, err
	}

	if err := validateScreenshot(screenshot); err != nil {
		return nil, err
	}

	return screenshot, nil
}

// validateScreenshot rejects captures that are empty, not a decodable
// image, or effectively blank (chrome's --screenshot flag can take the
// shot before webfonts/remote images have finished painting).
func validateScreenshot(data []byte) error {
	if len(data) < minScreenshotBytes {
		return fmt.Errorf("screenshot is too small (%d bytes)", len(data))
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("screenshot is not a decodable image: %w", err)
	}

	if isBlank(img) {
		return fmt.Errorf("screenshot appears blank")
	}

	return nil
}

// isBlank samples a grid of pixels across the image and reports whether
// they are all (near) identical, the signature of a screenshot taken
// before the page finished rendering.
func isBlank(img image.Image) bool {
	const step = 10
	const tolerance = 1024 // out of 65535 per channel

	bounds := img.Bounds()
	firstR, firstG, firstB, _ := img.At(bounds.Min.X, bounds.Min.Y).RGBA()

	for y := bounds.Min.Y; y < bounds.Max.Y; y += step {
		for x := bounds.Min.X; x < bounds.Max.X; x += step {
			r, g, b, _ := img.At(x, y).RGBA()
			if channelDiff(r, firstR) > tolerance || channelDiff(g, firstG) > tolerance || channelDiff(b, firstB) > tolerance {
				return false
			}
		}
	}

	return true
}

func channelDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
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
		log.Println("takeScreenshot failed:", err)
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
	log.Fatal(http.ListenAndServe(os.Getenv("APP_PORT"), nil))
}
