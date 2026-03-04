package infrastructure

import (
	"context"

	"github.com/vukieuhaihoa/bookmark-libs/pkg/common"
	"github.com/vukieuhaihoa/bookmark-libs/pkg/logger"
	bookmarkHandler "github.com/vukieuhaihoa/bookmark-worker/internal/app/handler/bookmark"
	bookmarkRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/bookmark"
	queueRepo "github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/queue"
	bookmarkService "github.com/vukieuhaihoa/bookmark-worker/internal/app/service/bookmark"
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

	queueRepo := queueRepo.NewRedisQueue(redisClient, cfg.QueueName)

	bookmarkRepo := bookmarkRepo.NewBookmarkRepository(dbClient)
	bookmarkSvc := bookmarkService.NewBookmarkService(bookmarkRepo)
	handler := bookmarkHandler.NewHandler(bookmarkSvc)

	engine := worker.NewEngine(queueRepo, handler)

	engine.Start(ctx)
}
