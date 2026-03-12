package worker

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/handler/bookmark"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/repository/queue"
)

const (
	DefaultNumWorkers = 5
)

type Engine interface {
	Start(ctx context.Context)
}

type engine struct {
	queue   queue.Queue
	handler bookmark.Handler
	run     bool
	sigChan chan os.Signal
	nrApp   *newrelic.Application
}

func NewEngine(queue queue.Queue, handler bookmark.Handler, nrApp *newrelic.Application) Engine {
	return &engine{
		queue:   queue,
		handler: handler,
		nrApp:   nrApp,
		sigChan: make(chan os.Signal, 1),
	}
}

func (e *engine) Start(ctx context.Context) {
	log.Info().Msg("starting worker engine")
	workerPool := NewPool(ctx, e.handler, DefaultNumWorkers, e.nrApp)
	signal.Notify(e.sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	e.run = true
	for e.run {
		select {
		case sig := <-e.sigChan:
			log.Info().Str("signal", sig.String()).Msg("received shutdown signal, stopping worker engine")
			e.run = false
		default:
			// get task from queue
			msg, err := e.queue.PopMessage(ctx)
			if err != nil {
				if !errors.Is(err, queue.ErrNoMessage) {
					log.Error().Err(err).Msg("failed to pop message from queue")
				}

				time.Sleep(1 * time.Second)
				continue
			}

			// dispatch task to worker pool message channel
			workerPool.Consume(msg)
		}
	}
	// wait for all workers to finish processing
	workerPool.Close()
}
