package load

import (
	"context"
	"github.com/mattjh1/goetl/config/logger"
  "github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/mattjh1/goetl/config"

)

type RedisDB struct {
    ctx        context.Context
    config     *config.RedisConfig
    embedder   *redisvector.Store
}

// NewRedisDB initializes and returns a RedisDB instance
func NewRedisDB(redisConfig *config.RedisConfig, ctx context.Context, model *embeddings.EmbedderImpl) (*RedisDB, error) {
    embedder, err := redisvector.New(ctx,
        redisvector.WithConnectionURL(redisConfig.URL),
        redisvector.WithIndexName(redisConfig.Index, true),
        redisvector.WithEmbedder(model),
    )
    if err != nil {
        return nil, err
    }

    return &RedisDB{
        ctx:      ctx,
        config:   redisConfig,
        embedder: embedder,
    }, nil
}

func (r *RedisDB) InsertEmbedding(docs []schema.Document, id string) error {
    _, err := r.embedder.AddDocuments(r.ctx, docs)
    return err
}

func getEmbedding(model string, connectionStr ...string) (llms.Model, *embeddings.EmbedderImpl) {
	opts := []ollama.Option{ollama.WithModel(model)}
	if len(connectionStr) > 0 {
		opts = append(opts, ollama.WithServerURL(connectionStr[0]))
	}
	llm, err := ollama.New(opts...)
	if err != nil {
		logger.Log.Fatal(err)
	}

	e, err := embeddings.NewEmbedder(llm)
	if err != nil {
		logger.Log.Fatal(err)
	}
	return llms.Model(llm), e
}
