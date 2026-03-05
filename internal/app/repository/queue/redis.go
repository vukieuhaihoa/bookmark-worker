package queue

import (
	"github.com/redis/go-redis/v9"
)

// redisQueue is a struct that implements the Queue interface using Redis as the underlying data store.
type redisQueue struct {
	c         *redis.Client
	queueName string
}

// NewRedisQueue creates a new instance of Redis queue implementing the Queue interface.
func NewRedisQueue(c *redis.Client, queueName string) Queue {
	return &redisQueue{
		c:         c,
		queueName: queueName,
	}
}
