package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ErrNoMessage = errors.New("no message in queue")

// PopMessage retrieves a message from the Redis queue. If the queue is empty, it returns ErrNoMessage.
// Parameters:
//   - ctx: The context for managing the request lifecycle.
//
// Returns:
//   - []byte: The message retrieved from the queue, or nil if the queue is empty.
//   - error: An error if there was an issue retrieving the message, or ErrNoMessage if the queue is empty.
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
