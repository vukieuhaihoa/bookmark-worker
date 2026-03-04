package queue

import "context"

type Queue interface {
	PopMessage(ctx context.Context) ([]byte, error)
}
