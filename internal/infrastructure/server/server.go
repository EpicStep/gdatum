// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Server is an HTTP server.
type Server struct {
	srv    *http.Server
	logger *zap.Logger
}

// New returns new Server.
func New(address string, handler http.Handler, logger *zap.Logger) *Server {
	if logger == nil {
		logger = zap.L()
	}

	return &Server{
		srv: &http.Server{
			Addr:    address,
			Handler: handler,
		},
		logger: logger,
	}
}

// Run HTTP server.
func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("stopping http server")

		sCtx, sCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer sCancel()

		if err := s.srv.Shutdown(sCtx); err != nil { //nolint:contextcheck
			s.logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	s.logger.Info("starting http server", zap.String("addr", s.srv.Addr))

	if err := s.srv.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}
