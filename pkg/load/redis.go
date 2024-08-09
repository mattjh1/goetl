package load

import (
	"context"
	"log"
  "github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/mattjh1/goetl/config"

)

type RedisDB struct {
    client *redisvector.RueidisClient
    ctx    context.Context
}

func NewRedisDB(client *redisvector.RueidisClient, ctx context.Context) *RedisDB {
	return &RedisDB{client: client, ctx: ctx}
}

func (r *RedisDB) InsertEmbedding(docs []schema.Document, id string, cfg *config.Config) error {
	// TODO replace with data from config 
	redisURL := "redis://127.0.0.1:6379"
	index := "test_redis_vectorstore"

	_, e := getEmbedding(cfg.EmbModelID, cfg.OllamaAPIBase)
	store, err := redisvector.New(r.ctx,
	redisvector.WithConnectionURL(redisURL),
	redisvector.WithIndexName(index, true),
	redisvector.WithEmbedder(e),
)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = store.AddDocuments(r.ctx, docs)

  return err
}

func getEmbedding(model string, connectionStr ...string) (llms.Model, *embeddings.EmbedderImpl) {
	opts := []ollama.Option{ollama.WithModel(model)}
	if len(connectionStr) > 0 {
		opts = append(opts, ollama.WithServerURL(connectionStr[0]))
	}
	llm, err := ollama.New(opts...)
	if err != nil {
		log.Fatal(err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}
	return llms.Model(llm), e
}
