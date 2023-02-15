package utils

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func GetBookTitleFromPath(path string) string {
	_, filename := filepath.Split(path)
	return strings.Split(filename, ".")[0]
}

func GetFileNameFromPath(path string) string {
	_, filename := filepath.Split(path)
	return strings.Split(filename, ".")[0]
}

func GetFilesWithSameSuffix(dir, suffix string) []string {
	data := []string{}

	files, err := os.ReadDir(dir)
	if err == nil {
		for _, fileInfo := range files {
			filepath := path.Join(dir, fileInfo.Name())
			if fileInfo.IsDir() {
				continue
			}

			if strings.HasSuffix(filepath, suffix) {
				data = append(data, filepath)
			}
		}
	}
	sort.Strings(data)
	return data
}
