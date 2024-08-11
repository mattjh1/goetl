package load
import (
	"github.com/tmc/langchaingo/schema"
)

type VectorDB interface {
    InsertEmbedding(docs []schema.Document, id string) error
}
