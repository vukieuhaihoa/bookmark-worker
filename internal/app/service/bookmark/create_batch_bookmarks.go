package bookmark

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
)

const (
	ListBookmarksCacheGroupKey = "list_bookmarks_%s"
)

// ImportBookmarkInput represents the structure of a bookmark to be imported. It includes a description and a URL, both of which are required fields with specific validation rules. The description must be a string with a maximum length of 255 characters, while the URL must be a valid URL string with a maximum length of 1024 characters. This struct is used to capture the details of each bookmark that is being imported through the API.
type ImportBookmarkInput struct {
	Description string `csv:"description" binding:"required,lte=255" json:"description"`
	URL         string `csv:"url" binding:"required,url,lte=1024" json:"url"`
}

// ImportMessage represents the structure of the message that will be sent to the queue for importing bookmarks. It contains the user ID (UID) and a slice of ImportBookmarkInput, which holds the details of each bookmark to be imported. This struct is used to encapsulate the data that will be processed by the worker responsible for handling bookmark imports from the queue.
type ImportMessage struct {
	UID       string                 `json:"user_id"`
	Bookmarks []*ImportBookmarkInput `json:"bookmarks"`
}

func (s *bookmarkService) CreateBatchBookmarks(ctx context.Context, importMsg *ImportMessage) error {
	// Segment for the entire service method
	seg := newrelic.FromContext(ctx).StartSegment("Service_CreateBatchBookmarks")
	defer seg.End()

	// Datastore segment for cache deletion (Redis)
	cacheSeg := newrelic.DatastoreSegment{
		StartTime:  newrelic.FromContext(ctx).StartSegmentNow(),
		Product:    newrelic.DatastoreRedis,
		Collection: "cache",
		Operation:  "DEL",
	}
	err := s.cache.DelCacheData(ctx, fmt.Sprintf(ListBookmarksCacheGroupKey, importMsg.UID))
	cacheSeg.End()
	if err != nil {
		return err
	}
	// Convert ImportBookmarkInput to model.Bookmark
	input := make([]*model.Bookmark, len(importMsg.Bookmarks))
	for i, bookmark := range importMsg.Bookmarks {
		input[i] = &model.Bookmark{
			URL:         bookmark.URL,
			Description: bookmark.Description,
			UserID:      importMsg.UID,
		}
	}

	// Call the repository method to create batch bookmarks
	err = s.repo.CreateBatchBookmarks(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
