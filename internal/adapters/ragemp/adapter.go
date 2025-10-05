// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package ragemp

import (
	"context"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/EpicStep/gdatum/internal/domain"
	"github.com/EpicStep/gdatum/internal/infrastructure/clients/ragemp"
)

type client interface {
	Servers(ctx context.Context) (ragemp.Servers, error)
}

// Adapter ...
type Adapter struct {
	client client
}

// New ...
func New(client client) *Adapter {
	return &Adapter{
		client: client,
	}
}

// Servers ...
func (a *Adapter) Servers(ctx context.Context, collectedAt time.Time) ([]domain.Server, error) {
	servers, err := a.client.Servers(ctx)
	if err != nil {
		return nil, err
	}

	return lo.MapToSlice(servers, func(ip string, server ragemp.Server) domain.Server {
		if index := strings.Index(ip, ":"); index != -1 {
			ip = ip[:index] // remove port
		}

		return domain.Server{
			Multiplayer: domain.MultiplayerRagemp,
			ID:          ip,
			Name:        server.Name,
			Gamemode:    server.Gamemode,
			Lang:        server.Lang,
			Players:     server.Players,
			CollectedAt: collectedAt,
		}
	}), nil
}
