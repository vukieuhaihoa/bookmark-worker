package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNoMessage = errors.New("no message in queue")

func (r *redisQueue) PopMessage(ctx context.Context) ([]byte, error) {
	val, err := r.c.RPop(ctx, r.queueName).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNoMessage
		}

		return nil, err
	}

	return val, nil
}
