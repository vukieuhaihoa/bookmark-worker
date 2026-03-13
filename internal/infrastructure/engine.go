package infrastructure

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/logger"

	bookmarkHandler "github.com/vukieuhaihoa/bookmark-worker/internal/app/handler/bookmark"
	bookmarkRepository "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/cache"
	bookmarkService "github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark"

	queueRepository "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/queue"
	"github.com/vukieuhaihoa/bookmark-worker/internal/worker"
)

func CreateEngineConfig() *worker.Config {
	cfg, err := worker.NewConfig()
	common.HandlerError(err)

	return cfg
}

func CreateEngine() {
	ctx := context.Background()

	logger.SetLogLevel()

	cfg := CreateEngineConfig()

	redisClient := CreateRedisCon()

	dbClient := CreateSQLDB()

	queueRepo := queueRepository.NewRedisQueue(redisClient, cfg.QueueName)
	cacheRepo := cache.NewRedisCache(redisClient)

	bookmarkRepo := bookmarkRepository.NewBookmarkRepository(dbClient)
	bookmarkSvc := bookmarkService.NewBookmarkService(bookmarkRepo, cacheRepo)
	handler := bookmarkHandler.NewHandler(bookmarkSvc)

	nrApp := CreateNewRelicClient()

	engine := worker.NewEngine(queueRepo, handler, nrApp)

	engine.Start(ctx)
}
