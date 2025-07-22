package semanticcache

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wso2/apk/gateway/enforcer/internal/logging"
)

const (
	embeddingField = "embedding"
	responseField  = "response"
	keyPrefix      = "doc:"
)

// RedisVectorDBProvider implements the VectorDBProvider interface for Redis
type RedisVectorDBProvider struct {
	redisURL  string
	database  int
	username  string
	password  string
	indexID   string
	dimension int
	ttl       int
	client    *redis.Client
}

// Init initializes the Redis vector DB provider with configuration
func (r *RedisVectorDBProvider) Init(logger *logging.Logger, config VectorDBProviderConfig) error {
	err := ValidateVectorStoreConfigProps(config)
	if err != nil {
		logger.Sugar().Errorf("Invalid vector store config properties: %v", err)
		return err
	}

	r.redisURL = config.DBHost + ":" + strconv.Itoa(config.DBPort)
	r.username = config.Username
	r.password = config.Password
	r.database, err = strconv.Atoi(config.DatabaseName)
	if err != nil {
		r.database = 0
		logger.Info("Switching to default due to invalid database number: %v", err)
	}

	embeddingDimension := config.EmbeddingDimention
	r.indexID = VectorIndexPrefix + embeddingDimension
	r.dimension, err = strconv.Atoi(embeddingDimension)
	if err != nil {
		logger.Sugar().Errorf("unable to parse and convert the embedding dimension to Int: %v", err)
		return err
	}

	r.ttl = DefaultTTL
	if config.TTL != "" {
		parsedTTL, err := strconv.Atoi(config.TTL)
		if err != nil {
			logger.Sugar().Errorf("invalid TTL value: %v", err)
			return err
		}
		r.ttl = parsedTTL
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     r.redisURL,
		Username: r.username,
		Password: r.password,
		DB:       r.database,
		Protocol: 2,
	})
	return nil
}

// GetType returns the type of the provider
func (r *RedisVectorDBProvider) GetType() string {
	return "REDIS"
}

// CreateIndex creates the Redis index for vector search
func (r *RedisVectorDBProvider) CreateIndex(logger *logging.Logger) error {
	// Check if collection/index exists
	_, err := r.client.Do(context.Background(), "FT.INFO", r.indexID).Result()
	if err == nil {
		// Index already exists
		logger.Sugar().Infof("Index %s already exists, skipping creation", r.indexID)
		return nil
	}

	_, err = r.client.FTCreate(context.Background(),
		r.indexID,
		&redis.FTCreateOptions{
			OnHash: true,
			Prefix: []any{"doc:"},
		},
		&redis.FieldSchema{
			FieldName: "api_id",
			FieldType: redis.SearchFieldTypeTag,
		},
		&redis.FieldSchema{
			FieldName: embeddingField,
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				HNSWOptions: &redis.FTHNSWOptions{
					Dim:            r.dimension,
					DistanceMetric: "L2",
					Type:           "FLOAT32",
				},
			},
		},
	).Result()

	if err != nil {
		return err
	}
	logger.Sugar().Info("Index successfully created with the given parameters")
	return nil
}

// Store stores an embedding in Redis along with the response
func (r *RedisVectorDBProvider) Store(logger *logging.Logger, embeddings []float32, response CacheResponse, filter map[string]interface{}) error {
	ctx := filter["ctx"].(context.Context)
	embeddingBytes := FloatsToBytes(embeddings)
	responseBytes, err := SerializeObject(response)
	if err != nil {
		logger.Sugar().Debugf("Unable to serialize the response object: %v\n", err.Error())
		return err
	}

	docID := uuid.New().String()
	redisKey := keyPrefix + docID

	_, err = r.client.HSet(ctx, redisKey, map[string]any{
		responseField:  responseBytes,
		"api_id":       filter["api_id"].(string),
		embeddingField: embeddingBytes,
	}).Result()

	if err != nil {
		logger.Sugar().Debugf("Failed to store the redis entry: %v\n", err.Error())
		return err
	}

	if r.ttl > 0 {
		_, err = r.client.Expire(ctx, redisKey, time.Duration(r.ttl)*time.Second).Result()
		if err != nil {
			logger.Sugar().Debugf("Failed to set the ttl for the specified redis entry: %v\n", err.Error())
			return err
		}
	}

	return nil
}

// Retrieve retrieves the most similar embedding from Redis
func (r *RedisVectorDBProvider) Retrieve(logger *logging.Logger, embeddings []float32, filter map[string]interface{}) (CacheResponse, error) {
	ctx := filter["ctx"].(context.Context)
	embeddingBytes := FloatsToBytes(embeddings)
	apiID := filter["api_id"].(string)
	if apiID == "" {
		logger.Sugar().Debugf("Given API ID: %s", apiID)
		logger.Sugar().Debug("Error: api_id is required in filter")
		return CacheResponse{}, errors.New("api_id is required in filter")
	}

	knnQuery := fmt.Sprintf(
		"@api_id:{\"%s\"}=>[KNN $K @%s $vec AS score]",
		apiID, embeddingField,
	)
	logger.Sugar().Debugf("KNN Query: %s", knnQuery)
	results, err := r.client.FTSearchWithArgs(ctx,
		r.indexID,
		knnQuery,
		&redis.FTSearchOptions{
			Return: []redis.FTSearchReturn{
				{FieldName: responseField},
				{FieldName: "score"},
			},
			DialectVersion: 2,
			Params: map[string]any{
				"K":   1,
				"vec": embeddingBytes,
			},
		},
	).Result()

	if err != nil {
		logger.Sugar().Errorf("Error during FTSearch: %v\n", err)
		return CacheResponse{}, err
	}

	if results.Total == 0 {
		logger.Sugar().Errorf("No results found: %v\n", err)
		return CacheResponse{}, errors.New("no results found")
	}

	// Take the top‐hit document
	doc := results.Docs[0]
	score, err := strconv.ParseFloat(doc.Fields["score"], 64)
	if err != nil {
		logger.Sugar().Errorf("Invalid doc score found: %s", err)
	}
	
	thresholdStr, ok := filter["threshold"].(string)
	if !ok {
		return CacheResponse{}, fmt.Errorf("missing threshold in filter")
	}
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		return CacheResponse{}, fmt.Errorf("invalid threshold: %w", err)
	}
	logger.Sugar().Debugf("Match Score: %f | Threshold: %f", score, threshold)
	if score > threshold {
		return CacheResponse{}, nil
	}

	// Fetch the serialized response blob from Redis
	respBytes, err := r.client.HGet(ctx, doc.ID, responseField).Bytes()
	if err != nil {
		return CacheResponse{}, err
	}

	var resp CacheResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return CacheResponse{}, err
	}

	return resp, nil
}

// Helper functions

// FloatsToBytes convert float[] to byte[] for storing in Redis(FROM DOCS)
func FloatsToBytes(fs []float32) []byte {
	buf := make([]byte, len(fs)*4)

	for i, f := range fs {
		u := math.Float32bits(f)
		binary.NativeEndian.PutUint32(buf[i*4:], u)
	}

	return buf
}

// SerializeObject Serialize an object to byte array
func SerializeObject(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize a byte array to object
func deserializeObject(data []byte, obj interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	return dec.Decode(obj)
}
