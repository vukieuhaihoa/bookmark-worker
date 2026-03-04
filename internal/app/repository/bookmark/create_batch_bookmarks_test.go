package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	"github.com/vukieuhaihoa/bookmark-worker/test/data"
	"gorm.io/gorm"
)

func TestRepository_CreateBatchBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupDB func(t *testing.T) *gorm.DB

		inputBookmarks []*model.Bookmark

		expectedError error
		verifyFunc    func(t *testing.T, db *gorm.DB, ctx context.Context, inputBookmarks []*model.Bookmark)
	}{
		{
			name: "create batch bookmarks successfully",

			setupDB: func(t *testing.T) *gorm.DB {
				return data.NewFixture(t, &data.BookmarkCommonTestDB{})
			},

			inputBookmarks: []*model.Bookmark{
				{
					Base: model.Base{
						ID: "a1b2c3d4-e5f6-7890-abcd-ef0000000077",
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/newbookmark1",
					Description: "New bookmark 1 for Test User 1",
				},
				{
					Base: model.Base{
						ID: "b1c2d3e4-f5g6-7890-abcd-ef0000000088",
					},
					UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					URL:         "https://example.com/newbookmark2",
					Description: "New bookmark 2 for Test User 1",
				},
			},

			expectedError: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB, ctx context.Context, inputBookmarks []*model.Bookmark) {
				for _, bookmark := range inputBookmarks {
					var foundBookmark model.Bookmark
					err := db.WithContext(ctx).First(&foundBookmark, "id = ?", bookmark.ID).Error
					assert.NoError(t, err)
					assert.Equal(t, bookmark.URL, foundBookmark.URL)
					assert.Equal(t, bookmark.Description, foundBookmark.Description)
					assert.Equal(t, bookmark.UserID, foundBookmark.UserID)
					assert.NotEmpty(t, foundBookmark.CodeShortenEncoded)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewBookmarkRepository(db)

			err := repo.CreateBatchBookmarks(ctx, tc.inputBookmarks)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			tc.verifyFunc(t, db, ctx, tc.inputBookmarks)
		})
	}
}
