package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/log"
)

func ScanAll(ctx context.Context) error {
	return filepath.WalkDir(config.Server.Scanner.Directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error().Err(err).Str("path", path).Msg("Unable to walk down the file tree")
		}

		fmt.Println(path)

		return nil
	})
}
