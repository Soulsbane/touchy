package api

import (
	"fmt"
	"github.com/imroc/req/v3"
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
