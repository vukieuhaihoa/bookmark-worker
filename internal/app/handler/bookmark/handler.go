package bookmark

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark"
)

// Handler defines the interface for handling bookmark-related messages.
type Handler interface {
	Handle(ctx context.Context, message []byte) error
}

// handler implements the Handler interface for processing bookmark-related messages.
type handler struct {
	svc bookmark.Service
}

// NewHandler creates a new instance of handler with the provided bookmark service.
// Parameters:
//   - svc: The bookmark service used to perform bookmark-related operations.
//
// Returns:
func NewHandler(svc bookmark.Service) Handler {
	return &handler{svc: svc}
}

var ErrUnmarshalMessage = errors.New("failed to unmarshal message")

// Handle processes the incoming message to create bookmarks in batch.
// It expects the message to be a JSON-encoded ImportMessage. The function first attempts to unmarshal the message into an ImportMessage struct. If unmarshaling fails, it returns an ErrUnmarshalMessage error. If unmarshaling is successful, it calls the CreateBatchBookmarks method of the bookmark service with the unmarshaled ImportMessage. If any error occurs during the creation of batch bookmarks, it returns that error; otherwise, it returns nil indicating successful handling of the message.
func (h *handler) Handle(ctx context.Context, message []byte) error {
	txn := newrelic.FromContext(ctx)
	s := txn.StartSegment("Handler_Handle")
	defer s.End()

	input := &bookmark.ImportMessage{}
	err := json.Unmarshal(message, input)
	if err != nil {
		return ErrUnmarshalMessage
	}

	// Custom attributes — visible in NR transaction details
	txn.AddAttribute("user_id", input.UID)
	txn.AddAttribute("bookmark_count", len(input.Bookmarks))

	err = h.svc.CreateBatchBookmarks(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
