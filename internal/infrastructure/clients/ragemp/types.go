// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package ragemp

// Servers is a ragemp servers.
type Servers map[string]Server

// Server is a regemp server.
type Server struct {
	Name     string `json:"name"`
	Gamemode string `json:"gamemode"`
	URL      string `json:"url"`
	Language string `json:"lang"`
	Players  int32  `json:"players"`
}
