package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/cache"
)

const DEFAULT_CODE_LENGTH = 10

// Service defines the interface for bookmark-related business logic.
//
//go:generate mockery --name=Service --filename=service.go --output=./mocks
type Service interface {
	// CreateBatchBookmarks creates multiple bookmarks based on the provided import message.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - importMsg: The message containing the batch of bookmarks to be created.
	//
	// Returns:
	//   - error: An error if the creation fails, otherwise nil.
	CreateBatchBookmarks(ctx context.Context, importMsg *ImportMessage) error
}

// bookmarkService is the concrete implementation of the Service interface.
type bookmarkService struct {
	repo  bookmark.Repository
	cache cache.Repository
}

// NewBookmarkService creates a new instance of bookmarkService.
// It takes a bookmark.Repository and a cache.Repository as its dependencies
// and returns a Service interface.
func NewBookmarkService(repo bookmark.Repository, cache cache.Repository) Service {
	return &bookmarkService{
		repo:  repo,
		cache: cache,
	}
}
