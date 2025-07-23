// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package worker

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// WorkFunc ...
type WorkFunc func(ctx context.Context) error

// Worker ...
type Worker struct {
	interval time.Duration
	workFunc WorkFunc

	logger *zap.Logger
}

// New returns new Worker.
func New(name string, interval time.Duration, workFunc WorkFunc, logger *zap.Logger) *Worker {
	if logger == nil {
		logger = zap.L()
	}

	logger = logger.Named("async-worker")
	logger = logger.With(zap.String("name", name))

	return &Worker{
		interval: interval,
		workFunc: workFunc,

		logger: logger,
	}
}

// Run ...
func (w *Worker) Run(ctx context.Context) error {
	timer := time.NewTimer(w.nextRunAfter())
	defer timer.Stop()

	var wg sync.WaitGroup

	w.logger.Info("starting worker")

	for {
		select {
		case <-timer.C:
			wg.Add(1)

			if err := w.workFunc(ctx); err != nil {
				w.logger.Error("work func call was failed", zap.Error(err))
			}

			wg.Done()
			timer.Reset(w.nextRunAfter())
		case <-ctx.Done():
			w.logger.Info("stopping worker")
			wg.Wait() // TODO: add timeout?
			w.logger.Info("worker has been stopped")
			return nil
		}
	}
}

func (w *Worker) nextRunAfter() time.Duration {
	now := time.Now()
	return now.Truncate(w.interval).Add(w.interval).Sub(now)
}
