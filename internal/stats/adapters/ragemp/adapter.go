// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package ragemp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/EpicStep/gdatum/internal/external/clients/ragemp"
	"github.com/EpicStep/gdatum/internal/stats/domain"
)

// Adapter ...
type Adapter struct {
	client *ragemp.Client
}

// New ...
func New(client *ragemp.Client) *Adapter {
	return &Adapter{
		client: client,
	}
}

// Servers ...
func (a *Adapter) Servers(ctx context.Context, collectedAt time.Time) ([]domain.Server, error) {
	servers, err := a.client.Servers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ragemp servers: %w", err)
	}

	return lo.MapToSlice(servers, func(ip string, server ragemp.Server) domain.Server {
		if index := strings.Index(ip, ":"); index != -1 {
			ip = ip[:index] // remove port
		}

		return domain.Server{
			Multiplayer: domain.MultiplayerRagemp,
			Identifier:  ip,
			Name:        server.Name,
			Gamemode:    server.Gamemode,
			Lang:        server.Lang,
			Players:     server.Players,
			CollectedAt: collectedAt,
		}
	}), nil
}
