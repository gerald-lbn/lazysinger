package file_test

// Copied from https://github.com/navidrome/navidrome (GPL 3.0 License)
// Copyright (c) 2025 Navidrome

import (
	"os"
	"path/filepath"

	"github.com/gerald-lbn/refrain/utils/file"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Exists", func() {
	var tempFile *os.File
	var tempDir string

	BeforeEach(func() {
		var err error
		tempFile, err = os.CreateTemp("", "fileexists-test-*.txt")
		Expect(err).NotTo(HaveOccurred())

		tempDir, err = os.MkdirTemp("", "fileexists-test-dir-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if tempFile != nil {
			os.Remove(tempFile.Name())
			tempFile.Close()
		}
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	It("returns true for existing file", func() {
		Expect(file.Exists(tempFile.Name())).To(BeTrue())
	})

	It("returns true for existing directory", func() {
		Expect(file.Exists(tempDir)).To(BeTrue())
	})

	It("returns false for non-existing file", func() {
		nonExistentPath := filepath.Join(tempDir, "does-not-exist.txt")
		Expect(file.Exists(nonExistentPath)).To(BeFalse())
	})

	It("returns false for empty path", func() {
		Expect(file.Exists("")).To(BeFalse())
	})

	It("handles nested non-existing path", func() {
		nonExistentPath := "/this/path/definitely/does/not/exist/file.txt"
		Expect(file.Exists(nonExistentPath)).To(BeFalse())
	})

	Context("when file is deleted after creation", func() {
		It("returns false after file deletion", func() {
			filePath := tempFile.Name()
			Expect(file.Exists(filePath)).To(BeTrue())

			err := os.Remove(filePath)
			Expect(err).NotTo(HaveOccurred())
			tempFile = nil // Prevent cleanup attempt

			Expect(file.Exists(filePath)).To(BeFalse())
		})
	})

	Context("when directory is deleted after creation", func() {
		It("returns false after directory deletion", func() {
			dirPath := tempDir
			Expect(file.Exists(dirPath)).To(BeTrue())

			err := os.RemoveAll(dirPath)
			Expect(err).NotTo(HaveOccurred())
			tempDir = "" // Prevent cleanup attempt

			Expect(file.Exists(dirPath)).To(BeFalse())
		})
	})

	It("handles permission denied scenarios gracefully", func() {
		// This test might be platform specific, but we test the general case
		result := file.Exists("/root/.ssh/id_rsa") // Likely to not exist or be inaccessible
		Expect(result).To(Or(BeTrue(), BeFalse())) // Should not panic
	})
})
