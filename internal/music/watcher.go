package music

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type LibraryWatcher struct {
	Directory string
	watcher   *fsnotify.Watcher

	HandleFileOnInitialScan func(pathToFile string) error
	HandleCreate            func(pathToFile string) error
	HandleRename            func(pathToFile string) error
	HandleMove              func(pathToFile string) error
	HandleChmod             func(pathToFile string) error
	HandleDelete            func(pathToFile string) error
}

func NewLibraryWatcher(directory string) (*LibraryWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &LibraryWatcher{
		Directory: directory,
		watcher:   watcher,
	}, nil
}

func (lw *LibraryWatcher) InitialScan() error {
	log.Printf("Performing initial scan of %s\n", lw.Directory)
	return filepath.WalkDir(lw.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("error walking %s: %v", lw.Directory, err)
			return err
		}

		if d.IsDir() {
			return lw.watcher.Add(path)
		}

		if d.Type().IsRegular() && IsMusicFile(path) {
			log.Printf("Found %s while performing library scan", path)
			if lw.HandleFileOnInitialScan != nil {
				if err := lw.HandleFileOnInitialScan(path); err != nil {
					log.Printf("Unable to handle file %s when performing library scan. %v", path, err)
				}
			}
		}

		return nil
	})
}

func (lw *LibraryWatcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-lw.watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Printf("Error stating %s: %v\n", event.Name, err)
						continue
					}
					if fileInfo.IsDir() {
						if err := lw.watcher.Add(event.Name); err != nil {
							log.Printf("Error adding directory %s to watcher: %v\n", event.Name, err)
						}
					} else if fileInfo.Mode().IsRegular() && IsMusicFile(event.Name) && lw.HandleCreate != nil {
						if err := lw.HandleCreate(event.Name); err != nil {
							log.Printf("Error handling create event for %s: %v\n", event.Name, err)
						}
					}
				}

				if event.Op&fsnotify.Rename == fsnotify.Rename {
					if lw.HandleRename != nil {
						if err := lw.HandleRename(event.Name); err != nil {
							log.Printf("Error handling rename event for %s: %v\n", event.Name, err)
						}
					}
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if lw.HandleDelete != nil {
						if err := lw.HandleDelete(event.Name); err != nil {
							log.Printf("Error handling delete event for %s: %v\n", event.Name, err)
						}
					}
				}

				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					if lw.HandleChmod != nil {
						if err := lw.HandleChmod(event.Name); err != nil {
							log.Printf("Error handling chmod event for %s: %v\n", event.Name, err)
						}
					}
				}

			case err, ok := <-lw.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v\n", err)
			}
		}
	}()
}

func (lw *LibraryWatcher) Wait() {
	// select {}
	<-make(chan struct{})
}

func (lw *LibraryWatcher) Close() error {
	return lw.watcher.Close()
}
