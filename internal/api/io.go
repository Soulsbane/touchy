package api

import "os"

func GetOutputDir() string {
	dir, _ := os.Getwd()
	return dir
}
