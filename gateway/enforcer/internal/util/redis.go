package util

import (
	"crypto/tls"

	"github.com/redis/go-redis/v9"
)

func CreateRedisClient(address string, user string, password string, tlsConfig *tls.Config) *redis.Client {
	// Create a Redis client
	options := &redis.Options{
        Addr:     address,
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