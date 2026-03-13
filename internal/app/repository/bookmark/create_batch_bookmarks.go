package bookmark

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/dbutils"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/encoding"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	"gorm.io/gorm"
)

func (r *bookmarkRepository) CreateBatchBookmarks(ctx context.Context, bookmarks []*model.Bookmark) error {
	txn := newrelic.FromContext(ctx)
	s := txn.StartSegment("Repository_CreateBatchBookmarks")
	defer s.End()
	// Datastore segment for Postgres batch insert
	dataSeg := newrelic.DatastoreSegment{
		StartTime:  newrelic.FromContext(ctx).StartSegmentNow(),
		Product:    newrelic.DatastorePostgres,
		Collection: "bookmarks",
		Operation:  "INSERT",
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, bookmark := range bookmarks {
			if err := tx.Omit("CodeShortenEncoded").Create(bookmark).Error; err != nil {
				return err
			}

			encoded, err := encoding.StdEncoding.EncodeInt64ToString(bookmark.CodeShorten)
			if err != nil {
				return err
			}

			bookmark.CodeShortenEncoded = model.BookmarkShortenPrefix + encoded

			if err := tx.Save(bookmark).Error; err != nil {
				return err
			}
		}

		return nil
	})

	dataSeg.End()
	if err != nil {
		return dbutils.CatchDBError(err)
	}

	return nil
}
