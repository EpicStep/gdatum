// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package api

import (
	"github.com/EpicStep/gdatum/pkg/api"
)

// Handlers ...
type Handlers struct {
	api.UnimplementedHandler
}

// New returns new Handlers.
func New() *Handlers {
	return &Handlers{}
}
