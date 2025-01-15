package load

import (
	"context"
	"fmt"

	"github.com/mattjh1/goetl/config"
	"github.com/mattjh1/goetl/config/logger"
	"github.com/tmc/langchaingo/schema"
)

func Load(ch <-chan schema.Document, cfg *config.Config) {
	ctx := context.Background()

	var db VectorDB
	var err error

	switch cfg.Database.Type {
	case "redis":
		redisConfig, _ := cfg.Database.GetRedisConfig()
		_, e := getEmbedding(cfg.EmbModelID, cfg.EmbAPIBase)

		if e != nil {
			logger.Log.Errorf("Failed to get embeddings: %v", e)
		}

		db, err = NewRedisDB(redisConfig, ctx, e)
		if err != nil {
			logger.Log.Fatalf("Failed to initialize RedisDB: %v", err)
			return
		}

	case "postgres":
		// TODO: Initialize Postgres connection with postgresConfig

	default:
		logger.Log.Fatalf("Unsupported database type: %s", cfg.Database.Type)
		return
	}

		// Process and load the data
	for chunk := range ch {
		// Generate document ID using the checksum
		docID := fmt.Sprintf("%s", chunk.Metadata["content_checksum"])

		// Insert the chunk's embeddings into the database
		err = db.InsertEmbedding([]schema.Document{chunk}, docID)
		if err != nil {
			logger.Log.Printf("Failed to insert embedding for chunk: %v", err)
			continue
		}
	}
}
