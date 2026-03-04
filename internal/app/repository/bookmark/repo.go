package bookmark

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	"gorm.io/gorm"
)

// Repository defines the interface for bookmark data operations.
//
//go:generate mockery --name=Repository --output=./mocks --filename=repo.go
type Repository interface {
	// CreateBookmark creates a new bookmark in the database.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - bookmark: The bookmark model to be created.
	//
	// Returns:
	//   - *model.Bookmark: The created bookmark model.
	//   - error: An error if the creation fails, otherwise nil.
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)

	// CreateBatchBookmarks creates multiple bookmarks in the database in a single transaction.
	//
	// Parameters:
	//   - ctx: The context for managing request-scoped values and cancellation.
	//   - bookmarks: A slice of bookmark models to be created.
	//
	// Returns:
	//   - error: An error if the creation fails, otherwise nil.
	CreateBatchBookmarks(ctx context.Context, bookmarks []*model.Bookmark) error
}

// bookmarkRepository implements the Repository interface for bookmark data operations.
type bookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new instance of bookmarkRepository.
// It takes a gorm.DB instance as a parameter and returns a Repository interface.
func NewBookmarkRepository(db *gorm.DB) Repository {
	return &bookmarkRepository{
		db: db,
	}
}
