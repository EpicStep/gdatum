// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package collector

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v5"
	"go.uber.org/zap"

	"github.com/EpicStep/gdatum/internal/domain"
)

type collectFunc func(ctx context.Context, collectedAt time.Time) ([]domain.Server, error)

type collectInstance struct {
	Multiplayer domain.Multiplayer
	Collect     collectFunc
}

func (h *Handler) collect(ctx context.Context) []domain.Server {
	collectedAt := time.Now().Truncate(time.Hour)

	var wg sync.WaitGroup

	var result []domain.Server
	var resultMux sync.Mutex

	for _, collector := range h.collectors {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var attempt int

			collectedServers, err := backoff.Retry(
				ctx,
				func() ([]domain.Server, error) {
					collectedServers, err := collector.Collect(ctx, collectedAt)
					if err != nil {
						attempt++

						h.logger.Error("failed to collect servers",
							zap.String("multiplayer", string(collector.Multiplayer)),
							zap.Int("attempt", attempt),
							zap.Error(err),
						)
						return nil, err
					}

					return collectedServers, nil
				},
				backoff.WithBackOff(backoff.NewExponentialBackOff()),
				backoff.WithMaxTries(3),
			)
			if err != nil {
				h.logger.Error("failed to collect servers",
					zap.String("multiplayer", string(collector.Multiplayer)),
					zap.Error(err),
				)

				h.metrics.RecordCollectionError(collector.Multiplayer)

				return
			}

			h.logger.Debug("collected servers",
				zap.String("multiplayer", string(collector.Multiplayer)),
				zap.Int("count", len(collectedServers)),
			)

			h.metrics.RecordServersCollected(collector.Multiplayer, len(collectedServers))

			resultMux.Lock()
			defer resultMux.Unlock()

			result = append(result, collectedServers...)
		}()
	}

	wg.Wait()

	return result
}
