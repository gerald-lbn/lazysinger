package file

// Copied from https://github.com/navidrome/navidrome (GPL 3.0 License)
// Copyright (c) 2025 Navidrome

import (
	"os"
)

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
