package queue

import "context"

// Queue is an interface that defines the methods for interacting with a message queue.
//
//go:generate mockery --name=Queue --output=./mocks --filename=queue.go --outpkg=mock_queue
type Queue interface {
	// PopMessage retrieves a message from the queue. If the queue is empty, it returns ErrNoMessage.
	// Parameters:
	//   - ctx: The context for managing the request lifecycle.
	//
	// Returns:
	//   - []byte: The message retrieved from the queue, or nil if the queue is empty.
	//   - error: An error if there was an issue retrieving the message, or ErrNoMessage if the queue is empty.
	PopMessage(ctx context.Context) ([]byte, error)
}
