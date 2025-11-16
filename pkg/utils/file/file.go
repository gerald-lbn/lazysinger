package file

// Copied from https://github.com/navidrome/navidrome (GPL 3.0 License)
// Copyright (c) 2025 Navidrome

import (
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// Exists checks if a file or directory exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func IsDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// IsAudioFile reads the mimetype of a file to check if it's an audio file
func IsAudioFile(path string) (bool, error) {
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return false, err
	}

	return strings.Contains(mtype.String(), "audio"), nil
}
