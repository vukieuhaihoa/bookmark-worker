package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark"
	mockBookmarkService "github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark/mocks"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupService func(ctx context.Context) *mockBookmarkService.Service

		inputMsg *bookmark.ImportMessage

		expectedError error
	}{
		{
			name: "handle import message successfully",

			setupService: func(ctx context.Context) *mockBookmarkService.Service {
				serviceMock := mockBookmarkService.NewService(t)
				serviceMock.On("CreateBatchBookmarks", ctx, &bookmark.ImportMessage{
					UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					Bookmarks: []*bookmark.ImportBookmarkInput{
						{
							URL:         "https://example.com/newbookmark1",
							Description: "New bookmark 1 for Test User 1",
						},
						{
							URL:         "https://example.com/newbookmark2",
							Description: "New bookmark 2 for Test User 1",
						},
					},
				}).Return(nil)
				return serviceMock
			},

			inputMsg: &bookmark.ImportMessage{
				UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				Bookmarks: []*bookmark.ImportBookmarkInput{
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
			name: "fail to handle import message due to service error",

			setupService: func(ctx context.Context) *mockBookmarkService.Service {
				serviceMock := mockBookmarkService.NewService(t)
				serviceMock.On("CreateBatchBookmarks", ctx, &bookmark.ImportMessage{
					UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
					Bookmarks: []*bookmark.ImportBookmarkInput{
						{
							URL:         "https://example.com/newbookmark1",
							Description: "New bookmark 1 for Test User 1",
						},
						{
							URL:         "https://example.com/newbookmark2",
							Description: "New bookmark 2 for Test User 1",
						},
					},
				}).Return(errors.New("service error"))
				return serviceMock
			},

			inputMsg: &bookmark.ImportMessage{
				UID: "4d9326d6-980c-4c62-9709-dbc70a82cbfe",
				Bookmarks: []*bookmark.ImportBookmarkInput{
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

			expectedError: errors.New("service error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			mockService := tc.setupService(ctx)
			h := NewHandler(mockService)

			messageBytes, err := json.Marshal(tc.inputMsg)
			assert.NoError(t, err)

			err = h.Handle(ctx, messageBytes)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
