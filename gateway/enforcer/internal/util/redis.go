package util

import (
	"crypto/tls"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// CreateRedisClient creates a new Redis client with the given parameters.
func CreateRedisClient(address string, user string, password string, tlsConfig *tls.Config) *redis.Client {
	fmt.Println("Creating Redis client")
	// Create a Redis client
	options := &redis.Options{
		Addr: address,
	}
	options.Password = password
	if user != "" {
		options.Username = user
	}
	if tlsConfig != nil {
		options.TLSConfig = tlsConfig
	}
	return redis.NewClient(options)

}
