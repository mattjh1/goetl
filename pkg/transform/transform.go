package transform

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"github.com/google/go-tika/tika"
	"github.com/tmc/langchaingo/schema"
	"github.com/mattjh1/goetl/config/logger"
	"github.com/mattjh1/goetl/config"
)

func calculateChecksum(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
func parse(in <-chan string, out chan<- schema.Document, client *tika.Client) {
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

				checksum := calculateChecksum(text)
				doc := schema.Document{
					PageContent: text,
					Metadata: map[string]any{
						"content_checksum": checksum,
					},
				}
				out <- doc
				}()

	}
	close(out)
}

func Transform(in <-chan string, out chan<- schema.Document, cfg *config.Config) {
	client := tika.NewClient(nil, cfg.TikaServerURL)
	parse(in, out, client)
}
