package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark/mocks"
)

func TestService_CreateBatchBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockRepo func(ctx context.Context) *mockBookmarkRepo.Repository

		inputImportMsg *ImportMessage

		expectedError error
	}{
		{
			name: "create batch bookmarks successfully",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBatchBookmarks", ctx, []*model.Bookmark{
					{
						URL:         "https://example.com/newbookmark1",
						Description: "New bookmark 1 for Test User 1",
						UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					},
					{
						URL:         "https://example.com/newbookmark2",
						Description: "New bookmark 2 for Test User 1",
						UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					},
				}).Return(nil)
				return repoMock
			},

			inputImportMsg: &ImportMessage{
				UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				Bookmarks: []*ImportBookmarkInput{
					{
						URL:         "https://example.com/newbookmark1",
						Description: "New bookmark 1 for Test User 1",
					},
					{
						URL:         "https://example.com/newbookmark2",
						Description: "New bookmark 2 for Test User 1",
					},
				},
			},

			expectedError: nil,
		},
		{
			name: "fail to create batch bookmarks due to repository error",

			setupMockRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
				repoMock.On("CreateBatchBookmarks", ctx, []*model.Bookmark{
					{
						URL:         "https://example.com/newbookmark1",
						Description: "New bookmark 1 for Test User 1",
						UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					},
					{
						URL:         "https://example.com/newbookmark2",
						Description: "New bookmark 2 for Test User 1",
						UserID:      "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					},
				}).Return(assert.AnError)
				return repoMock
			},

			inputImportMsg: &ImportMessage{
				UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				Bookmarks: []*ImportBookmarkInput{
					{
						URL:         "https://example.com/newbookmark1",
						Description: "New bookmark 1 for Test User 1",
					},
					{
						URL:         "https://example.com/newbookmark2",
						Description: "New bookmark 2 for Test User 1",
					},
				},
			},

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)
			service := NewBookmarkService(mockRepo)

			err := service.CreateBatchBookmarks(ctx, tc.inputImportMsg)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
