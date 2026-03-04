package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark"
)

const DEFAULT_CODE_LENGTH = 10

// Service defines the interface for bookmark-related business logic.
//
//go:generate mockery --name=Service --filename=service.go --output=./mocks
type Service interface {
	// CreateBookmark creates a new bookmark with the provided information.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - url: The URL of the bookmark.
	//   - description: A description of the bookmark.
	//   - userID: The ID of the user who owns the bookmark.
	//
	// Returns:
	//   - *model.Bookmark: The created bookmark model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateBookmark(ctx context.Context, url, description, userID string) (*model.Bookmark, error)

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
	repo bookmark.Repository
}

// NewBookmarkService creates a new instance of bookmarkService.
// It takes a bookmark.Repository as its dependency
// and returns a Service interface.
func NewBookmarkService(repo bookmark.Repository) Service {
	return &bookmarkService{
		repo: repo,
	}
}
