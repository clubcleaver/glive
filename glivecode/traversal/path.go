package traversal

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type PathAndMod map[string]time.Time

func GetPathAndMod(root string, skip []string) (pathAndMod PathAndMod, err error) {
	pathAndMod = make(PathAndMod)
	// root is already checked for validity in conf package

	// Root Traversal
	err = filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		stats, err := os.Stat(path) // Get Stats on path
		if err != nil {
			return err
		}

		for _, sfv := range skip { // else check against skip files
			if sfv == path || sfv == "./"+path || sfv == "/"+path || sfv == "."+path {
				if stats.IsDir() {
					return filepath.SkipDir // skip directory

				}
				return nil // skip file
			}
		}
		pathAndMod[path] = stats.ModTime()

		return nil
	})
	return
}
