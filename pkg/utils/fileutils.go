package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CountFilesInPath(path string, globPattern string) (int, error) {
	matches := 0

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", p, err)
		}

		// Check if the current file matches the glob pattern
		matched, matchErr := filepath.Match(globPattern, filepath.Base(p))
		if matchErr != nil {
			return fmt.Errorf("error matching file %s: %w", p, matchErr)
		}

		if matched && !info.IsDir() {
			matches++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return matches, nil
}
