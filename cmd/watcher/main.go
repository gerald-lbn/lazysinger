package main

import (
	"log"

	"github.com/gerald-lbn/lyrisync/internal/config"
	"github.com/gerald-lbn/lyrisync/internal/music"
)

func main() {
	cfg := config.LoadConfig()

	watcher, err := music.NewLibraryWatcher(cfg.MusicLibraryPath)
	if err != nil {
		log.Fatalf("error creating scanner: %v", err)
	}
	defer watcher.Close()

	watcher.HandleCreate = func(pathToFile string) error {
		log.Printf("New file created: %s", pathToFile)
		return nil
	}

	watcher.HandleDelete = func(pathToFile string) error {
		log.Printf("File deleted: %s.", pathToFile)
		return nil
	}

	watcher.HandleMove = func(pathToFile string) error {
		log.Printf("Filed moved: %s.", pathToFile)
		return nil
	}

	watcher.HandleRename = func(pathToFile string) error {
		log.Printf("Filed renamed: %s", pathToFile)
		return nil
	}

	if err := watcher.InitialScan(); err != nil {
		log.Fatalf("Initial scan failed: %v", err)
	} else {
		log.Printf("Initial scan completed")
	}

	watcher.Start()
	watcher.Wait()
}
