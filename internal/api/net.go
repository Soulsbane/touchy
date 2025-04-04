package api

import (
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads a file from a given URL and saves it to the specified file path.
// Returns a true indicating if the download was successful and false if not. The is returned in the second
// parameter.
func DownloadFile(destinationPath string, url string) (bool, string) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return false, err.Error()
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err.Error()
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("failed to download file: %s", url)
	}

	out, err := os.Create(destinationPath)

	if err != nil {
		return false, err.Error()
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return false, err.Error()
	}

	return true, ""
}

// DownloadFileWithProgress downloads a file from a given URL along with a progress bar displayed in the terminal
// and saves it to the specified file path. Returns a true indicating if the download was successful and false if not.
// The error string is returned in the second parameter.
func DownloadFileWithProgress(destinationPath string, url string) (bool, string) {
	tempDestinationPath := destinationPath + ".tmp"
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return false, err.Error()
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err.Error()
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Sprintf("failed to download file: %s", url)
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	defer resp.Body.Close()

	f, err := os.OpenFile(tempDestinationPath, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return false, fmt.Sprintf("failed create file: %s", err.Error())
	}

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)

	if err != nil {
		return false, fmt.Sprintf("failed create file: %s", err.Error())
	}

	err = os.Rename(tempDestinationPath, destinationPath)

	if err != nil {
		return false, fmt.Sprintf("failed create file: %s", err.Error())
	}

	return true, ""
}
