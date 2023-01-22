package utils

import (
	"path/filepath"
	"strings"
)

func GetBookTitleFromPath(path string) string {
	_, filename := filepath.Split(path)
	return strings.Split(filename, ".")[0]
}
