package bookmark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	mockBookmarkRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark/mocks"
	mockCacheRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/cache/mocks"
)

func TestService_CreateBatchBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockCacheRepo    func(ctx context.Context) *mockCacheRepo.Repository
		setupMockBookmarkRepo func(ctx context.Context) *mockBookmarkRepo.Repository

		inputImportMsg *ImportMessage

		expectedError error
	}{
		{
			name: "create batch bookmarks successfully",

			setupMockCacheRepo: func(ctx context.Context) *mockCacheRepo.Repository {
				cacheMock := mockCacheRepo.NewRepository(t)
				cacheMock.On("DelCacheData", ctx, "list_bookmarks_4d9326d6-980c-4c62-9709-dbc70a82cbfe").Return(nil)
				return cacheMock
			},

			setupMockBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
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
			name: "fail to create batch bookmarks due to cache repository error",

			setupMockCacheRepo: func(ctx context.Context) *mockCacheRepo.Repository {
				cacheMock := mockCacheRepo.NewRepository(t)
				cacheMock.On("DelCacheData", ctx, "list_bookmarks_4d9326d6-980c-4c62-9709-dbc70a82cbfe").Return(assert.AnError)
				return cacheMock
			},

			setupMockBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
				repoMock := mockBookmarkRepo.NewRepository(t)
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
		{
			name: "fail to create batch bookmarks due to bookmark repository error",

			setupMockCacheRepo: func(ctx context.Context) *mockCacheRepo.Repository {
				cacheMock := mockCacheRepo.NewRepository(t)
				cacheMock.On("DelCacheData", ctx, "list_bookmarks_4d9326d6-980c-4c62-9709-dbc70a82cbfe").Return(nil)
				return cacheMock
			},

			setupMockBookmarkRepo: func(ctx context.Context) *mockBookmarkRepo.Repository {
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
			mockCacheRepo := tc.setupMockCacheRepo(ctx)
			mockBookmarkRepo := tc.setupMockBookmarkRepo(ctx)
			service := NewBookmarkService(mockBookmarkRepo, mockCacheRepo)

			err := service.CreateBatchBookmarks(ctx, tc.inputImportMsg)
			if err != nil {
				assert.Equal(t, tc.expectedError, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
