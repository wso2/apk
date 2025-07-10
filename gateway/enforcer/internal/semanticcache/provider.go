package semanticcache

import "github.com/wso2/apk/gateway/enforcer/internal/logging"

// EmbeddingProvider defines the interface for services that provide text embedding
type EmbeddingProvider interface {
	Init(logger *logging.Logger, config EmbeddingProviderConfig) error
	GetType() string
	GetEmbedding(logger *logging.Logger, input string) ([]float32, error)
	GetEmbeddings(logger *logging.Logger, inputs []string) ([][]float32, error)
}

// VectorDBProvider defines the interface for vector database providers
type VectorDBProvider interface {
	Init(logger *logging.Logger, config VectorDBProviderConfig) error
	GetType() string
	CreateIndex(logger *logging.Logger) error
	Store(logger *logging.Logger, embeddings []float32, response CacheResponse, filter map[string]interface{}) error
	Retrieve(logger *logging.Logger, embeddings []float32, filter map[string]interface{}) (CacheResponse, error)
}

// VectorDBProviderConfig defines the properties required for initializing a vector DB provider
type VectorDBProviderConfig struct {
	VectorStoreProvider string
	EmbeddingDimention  string
	DistanceMetric      string
	Threshold           string
	DBHost              string
	DBPort              int
	Username            string
	Password            string
	DatabaseName        string
	TTL                 string
}

// EmbeddingProviderConfig defines the properties required for initializing an embedding provider
type EmbeddingProviderConfig struct {
	AuthHeaderName    string
	EmbeddingProvider string
	EmbeddingEndpoint string
	APIKey            string
	EmbeddingModel    string
	TimeOut           string
}
