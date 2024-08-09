package load
import (
	"github.com/tmc/langchaingo/schema"
	"github.com/mattjh1/goetl/config"
)

type VectorDB interface {
    InsertEmbedding(docs []schema.Document, id string, cfg *config.Config) error
}
