package load

import (
	"log"
	"context"

	"github.com/mattjh1/goetl/config"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
)

// Load processes transformed data and stores it into the VectorDB.
func Load(ch <-chan string, cfg *config.Config) {
	ctx := context.Background()
	client, err := redisvector.NewRueidisClient("redis://127.0.0.1:6379")
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}
	redisDB := NewRedisDB(client, ctx)

	for transformedData := range ch {
		log.Println("Loading transformed data:", transformedData)
		
		// Wrap transformed data into a schema.Document
		doc := schema.Document{
			PageContent: transformedData,
		}

		// Insert document embedding into Redis
		err := redisDB.InsertEmbedding([]schema.Document{doc}, doc.PageContent, cfg)
		if err != nil {
			log.Printf("Failed to insert embedding for document %s: %v", doc.PageContent, err)
			continue
		}

		log.Printf("Successfully loaded document with ID: %s", doc.PageContent)
	}
}

