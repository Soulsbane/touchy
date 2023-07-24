package api

import "os"

type Path struct {
}

func GetOutputDir() string {
	dir, _ := os.Getwd()
	return dir
}

func (path *Path) GetOutputDir() string {
	return GetOutputDir()
}
