package api

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads a file from a given URL and saves it to the specified file path.
// Returns a true indicating if the download was successful and false if not. The is returned in the second
// parameter.
func DownloadFile(filePath string, url string) (bool, string) {
	client := req.C()
	resp, err := client.R().SetOutputFile(filePath).Get(url)

	if err != nil || resp.IsErrorState() {
		return true, fmt.Errorf("failed to download file: %s", url).Error()
	} else {
		return false, ""
	}
}

func DownloadFileWithProgress(destinationPath string, url string) (bool, string) {
	tempDestinationPath := destinationPath + ".tmp"
	req, err := http.NewRequest("GET", url, nil)

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
		return false, fmt.Sprintf("failed create file: ", err.Error())
	}

	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)

	if err != nil {
		return false, fmt.Sprintf("failed create file: ", err.Error())
	}

	err = os.Rename(tempDestinationPath, destinationPath)

	if err != nil {
		return false, fmt.Sprintf("failed create file: ", err.Error())
	}

	return true, ""
}
