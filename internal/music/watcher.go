package music

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/fsnotify/fsnotify"
	"github.com/gerald-lbn/lazysinger/internal/log"
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
	log.Debug().Str("directory", directory).Msg("Starting library watcher for direction")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("An error occured when creating watcher")
		return nil, err
	}

	return &LibraryWatcher{
		Directory: directory,
		watcher:   watcher,
	}, nil
}

func (lw *LibraryWatcher) InitialScan() error {
	log.Info().Str("directory", lw.Directory).Msg("Performing initial scan")
	return filepath.WalkDir(lw.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Str("directory", lw.Directory).Err(err).Msg("An errror occured when running initial scan")
			return err
		}

		if d.IsDir() {
			log.Debug().Str("path", path).Msg("Adding path to the watchlist")
			return lw.watcher.Add(path)
		}

		if d.Type().IsRegular() && IsMusicFile(path) {
			log.Debug().Str("path", path).Msg("Found file while performing initial scan")
			if lw.HandleFileOnInitialScan != nil {
				if err := lw.HandleFileOnInitialScan(path); err != nil {
					log.Error().Str("path", path).Err(err).Msg("Unable to handle file when performing library scan")
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

				log.Debug().Str("event", event.Op.String()).Str("path", event.Name).Str("directory", lw.Directory).Msg("An event occured while watching directory")
				if event.Has(fsnotify.Create) {
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Error().Err(err).Str("path", event.Name).Msg("An error occured when stating filepath")
					}
					if fileInfo.IsDir() {
						if err := lw.watcher.Add(event.Name); err != nil {
							log.Error().Str("path", event.Name).Err(err).Msg("An error occured when adding directory to watchlist")
						}
					} else if fileInfo.Mode().IsRegular() && IsMusicFile(event.Name) && lw.HandleCreate != nil {
						if err := lw.HandleCreate(event.Name); err != nil {
							log.Error().Err(err).Str("path", event.Name).Msg("An error occured when handling file creation")
						}
					}
				}

				if event.Has(fsnotify.Rename) {
					if slices.Contains(lw.watcher.WatchList(), event.Name) {
						err := lw.watcher.Remove(event.Name)
						if err != nil {
							log.Error().Err(err).Str("path", event.Name).Msg("An error occured while removing filepath from watchlist")
						} else {
							log.Info().Str("path", event.Name).Msg("Successfully removed filepath from watchlist")
						}
					}
					if lw.HandleRename != nil {
						if err := lw.HandleRename(event.Name); err != nil {
							log.Error().Err(err).Str("path", event.Name).Msg("An error occured when handling renaming")
						}
					}
				}

				if event.Has(fsnotify.Remove) {
					if lw.HandleDelete != nil {
						if err := lw.HandleDelete(event.Name); err != nil {
							log.Error().Err(err).Str("path", event.Name).Msg("An error occurred when handling deletion")
						}
					}
				}

				if event.Has(fsnotify.Chmod) {
					if lw.HandleChmod != nil {
						if err := lw.HandleChmod(event.Name); err != nil {
							log.Error().Err(err).Str("path", event.Name).Msg("An error occured when handling chmod event")
						}
					}
				}

			case err, ok := <-lw.watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msg("An error occured with the watcher")
			}
		}
	}()

	err := lw.watcher.Add(lw.Directory)
	if err != nil {
		log.Fatal().Err(err).Str("path", lw.Directory).Msg("An error occurred while directory to watcher")
	}
}

func (lw *LibraryWatcher) Wait() {
	log.Info().Msg("Waiting for FS event")
	<-make(chan struct{})
}

func (lw *LibraryWatcher) Close() error {
	log.Info().Msg("Closing watcher")
	return lw.watcher.Close()
}
