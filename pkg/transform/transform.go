package transform

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"github.com/google/go-tika/tika"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/mattjh1/goetl/config"
	"github.com/mattjh1/goetl/config/logger"
)

func calculateChecksum(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func parse(in <-chan string, out chan<- schema.Document, client *tika.Client, cfg *config.Config) {
	ctx := context.Background()

	for filePath := range in {
		file, err := os.Open(filePath)
		if err != nil {
			logger.Log.Printf("Failed to open file %s: %v", filePath, err)
			continue
		}

		func() {
			defer file.Close()
			text, err := client.Parse(ctx, file)

			if err != nil {
				logger.Log.Printf("Failed to parse file %s: %v", filePath, err)
			}

			// Split the parsed text into chunks
			splitter := textsplitter.NewRecursiveCharacter()
			splitter.ChunkSize = cfg.ChunkSize
			splitter.ChunkOverlap = cfg.ChunkOverlap

			chunks, err := textsplitter.SplitDocuments(splitter, []schema.Document{
				{
					PageContent: text,
					Metadata:    map[string]any{"file_path": filePath},
				},
			})
			if err != nil {
				logger.Log.Printf("Failed to split document: %v", err)
				return
			}

			// Process each chunk separately
			for _, chunk := range chunks {
				checksum := calculateChecksum(chunk.PageContent)
				chunk.Metadata["content_checksum"] = checksum
				out <- chunk
			}
		}()
	}
	close(out)
}

func Transform(in <-chan string, out chan<- schema.Document, cfg *config.Config) {
	client := tika.NewClient(nil, cfg.TikaServerURL)
	parse(in, out, client, cfg)
}
