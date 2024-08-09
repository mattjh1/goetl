package transform

import (
	"context"
	"fmt"
	"os"
	"github.com/google/go-tika/tika"
	"github.com/mattjh1/goetl/config/logger"
	"github.com/mattjh1/goetl/config"
)

func parse(in <-chan string, out chan<- string, client *tika.Client) {
	ctx := context.Background()

	for filePath := range in {
			file, err := os.Open(filePath)
			if err != nil {
					logger.Log.Printf("Failed to open file %s: %v", filePath, err)
					continue
			}
			defer file.Close()

			// Extract text from the file
			text, err := client.Parse(ctx, file)
			if err != nil {
					logger.Log.Printf("Failed to parse file %s: %v", filePath, err)
					continue
			}

			fmt.Printf("Extracted Text from %s:\n%s\n", filePath, text)
			out <- text
	}
	close(out)
}

func Transform(in <-chan string, out chan<- string, cfg *config.Config) {
	client := tika.NewClient(nil, cfg.TikaServerURL)
	parse(in, out, client)
}
