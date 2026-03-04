package queue

import (
	"github.com/redis/go-redis/v9"
)

type redisQueue struct {
	c         *redis.Client
	queueName string
}

func NewRedisQueue(c *redis.Client, queueName string) Queue {
	return &redisQueue{
		c:         c,
		queueName: queueName,
	}
}
