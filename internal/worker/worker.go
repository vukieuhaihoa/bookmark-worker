package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vukieuhaihoa/bookmark-worker/internal/app/handler/bookmark"
)

type worker struct {
	id           int
	err          error
	handler      bookmark.Handler
	wg           *sync.WaitGroup
	messagesChan <-chan []byte
	errChan      chan<- *worker
}

func (w *worker) Work(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				w.err = err
			} else {
				w.err = fmt.Errorf("worker panicked: %v", r)
			}
			w.errChan <- w
		} else {
			log.Debug().Int("worker_id", w.id).Msg("worker exiting normally")
			w.wg.Done()
		}
	}()

	for {
		msg, ok := <-w.messagesChan
		if !ok {
			log.Debug().Int("worker_id", w.id).Msg("message channel closed, worker exiting")
			return
		}

		log.Debug().Int("worker_id", w.id).Bytes("message", msg).Msg("worker received message")

		err := w.handler.Handle(ctx, msg)
		if err != nil {
			log.Error().Err(err).Bytes("message", msg).Int("worker_id", w.id).Msg("failed to handle message")
		} else {
			log.Debug().Int("worker_id", w.id).Bytes("message", msg).Msg("worker handled message successfully")
		}
	}
}

type Pool struct {
	handler      bookmark.Handler
	numWorkers   int
	wg           *sync.WaitGroup
	messagesChan chan []byte
	errChan      chan *worker
}

func NewPool(ctx context.Context, handler bookmark.Handler, numWorkers int) *Pool {
	log.Info().Int("num_workers", numWorkers).Msg("creating worker pool")
	messagesChan := make(chan []byte, numWorkers)
	errChan := make(chan *worker, numWorkers)

	pool := &Pool{
		handler:      handler,
		numWorkers:   numWorkers,
		wg:           &sync.WaitGroup{},
		messagesChan: messagesChan,
		errChan:      errChan,
	}

	pool.init(ctx)

	return pool
}

func (p *Pool) init(ctx context.Context) {
	for i := 1; i <= p.numWorkers; i++ {
		log.Debug().Int("worker_id", i).Msg("starting worker")
		worker := &worker{
			id:           i,
			handler:      p.handler,
			wg:           p.wg,
			messagesChan: p.messagesChan,
			errChan:      p.errChan,
		}

		p.wg.Add(1)
		go worker.Work(ctx)
	}

	go func() {
		for w := range p.errChan {
			log.Error().Err(w.err).Int("worker_id", w.id).Msg("worker encountered an error")

			w.err = nil // reset error to prevent repeated logging

			// restart the worker
			time.Sleep(2 * time.Second) // add a small delay before restarting
			log.Info().Int("worker_id", w.id).Msg("restarting worker")
			go w.Work(ctx)
		}
	}()
}

func (p *Pool) Consume(msg []byte) {
	p.messagesChan <- msg
}

func (p *Pool) Close() {
	close(p.messagesChan)
	p.wg.Wait()
	close(p.errChan)
	log.Debug().Msg("all workers have exited")
}
