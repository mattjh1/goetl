package extract

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mattjh1/goetl/config"
	"github.com/mattjh1/goetl/config/logger"
)

func directory(path string, globPattern string, since time.Time) ([]string, error) {
	var files []string
	logger.Log.Infof("Starting directory walk for path: %s with pattern: %s", path, globPattern)
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		logger.Log.Infof("Checking file: %s", info.Name())

		// Check if the file matches the glob pattern
		if matched, err := filepath.Match(globPattern, info.Name()); err != nil {
			return err
		} else if matched {
			// Check if the file's modification time is after or equal to the 'since' time
			if info.ModTime().After(since) || info.ModTime().Equal(since) {
				files = append(files, filePath)
				logger.Log.Infof("File matched: %s", filePath)
			} else {
				logger.Log.Infof("File skipped (not modified since %v): %s", since, filePath)
			}
		} else {
			logger.Log.Infof("File did not match pattern: %s", info.Name())
		}
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return files, nil
}

// Extract function that reads files from a directory and sends their paths to a channel
func Extract(ch chan<- string, cfg *config.Config, path string, globPattern string, since time.Time) {
	logger.Log.Info("Starting file extraction...")
	files, err := directory(path, globPattern, since)
	if err != nil {
		logger.Log.Errorf("Error extracting files: %v", err)
		close(ch)
		return
	}
	for _, file := range files {
		logger.Log.Infof("Sending file to channel: %s", file)
		ch <- file
	}
	close(ch)
	logger.Log.Info("File extraction completed.")
}
