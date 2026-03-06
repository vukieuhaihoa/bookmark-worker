package cache

import "github.com/redis/go-redis/v9"

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new instance of Redis cache implementing the Repository interface.
//
// Parameters:
//   - client: The Redis client to be used for cache operations.
//
// Returns:
//   - Repository: An instance of the cache repository interface.
func NewRedisCache(client *redis.Client) Repository {
	return &redisCache{
		client: client,
	}
}
