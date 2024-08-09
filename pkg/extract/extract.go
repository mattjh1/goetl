package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mattjh1/goetl/config"
)

// directory function that returns a slice of file paths modified since the given date
func directory(path string, globPattern string, since time.Time) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, err := filepath.Match(globPattern, info.Name()); err != nil {
			return err
		} else if matched {
			if info.ModTime().After(since) || info.ModTime().Equal(since) {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

// Extract function that reads files from a directory and sends their paths to a channel
func Extract(ch chan<- string, cfg *config.Config, path string, globPattern string, since time.Time) {
	files, err := directory(path, globPattern, since)
	if err != nil {
		fmt.Println("Error extracting files:", err)
		close(ch)
		return
	}
	for _, file := range files {
		ch <- file
	}
	close(ch)
}
