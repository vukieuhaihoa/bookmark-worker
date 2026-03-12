package worker

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	redisPkg "github.com/vukieuhaihoa/bookmark-libs/pkg/redis"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/sqldb"
	bookmarkHandler "github.com/vukieuhaihoa/bookmark-worker/internal/app/handler/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/model"
	bookmarkRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/cache"
	queueRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/queue"
	bookmarkService "github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/worker"
	"gorm.io/gorm"
)

func makeTestWorkerEngine(t *testing.T, ctx context.Context, message []string) (worker.Engine, *gorm.DB) {
	mockRedisClient := redisPkg.InitMockRedis(t)
	for _, msg := range message {
		err := mockRedisClient.LPush(ctx, "test_queue", msg).Err()
		assert.NoError(t, err)
	}

	db := sqldb.InitMockDB(t)

	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	db.AutoMigrate(&model.Bookmark{})

	testQueueRepo := queueRepo.NewRedisQueue(mockRedisClient, "test_queue")
	testCacheRepo := cache.NewRedisCache(mockRedisClient)

	bookmarkRepo := bookmarkRepo.NewBookmarkRepository(db)
	bookmarkSvc := bookmarkService.NewBookmarkService(bookmarkRepo, testCacheRepo)
	handler := bookmarkHandler.NewHandler(bookmarkSvc)
	testEngine := worker.NewEngine(testQueueRepo, handler, nil)

	return testEngine, db
}

func TestWorkerEngine_Start(t *testing.T) {
	ctx := t.Context()

	messages := []string{
		`{"user_id":"4d9326d6-980c-4c62-9709-dbc70a82cbfe","bookmarks":[{"url":"https://example.com/newbookmark1","description":"New bookmark 1 for Test User 1"},{"url":"https://example.com/newbookmark2","description":"New bookmark 2 for Test User 1"}]}`,
		`{"user_id":"4d9326d6-980c-4c62-9709-dbc70a82cbfe","bookmarks":[{"url":"https://example.com/newbookmark3","description":"New bookmark 3 for Test User 1"},{"url":"https://example.com/newbookmark4","description":"New bookmark 4 for Test User 1"}]}`,
		`{"user_id":"4d9326d6-980c-4c62-9709-dbc70a82cbfe","bookmarks":[{"url":"https://example.com/newbookmark5","description":"New bookmark 5 for Test User 1"},{"url":"https://example.com/newbookmark6","description":"New bookmark 6 for Test User 1"}]}`,
	}

	testEngine, db := makeTestWorkerEngine(t, ctx, messages)
	go testEngine.Start(ctx)
	time.Sleep(1 * time.Second)
	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)
	err = process.Signal(os.Interrupt)
	assert.NoError(t, err)

	time.Sleep(10 * time.Second)

	var count int64
	err = db.WithContext(ctx).Model(&model.Bookmark{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(6), count)
}
