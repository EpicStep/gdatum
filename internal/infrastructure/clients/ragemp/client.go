// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package ragemp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	serverListURL = "https://cdn.rage.mp/master/"
)

// Client ...
type Client struct {
	client *http.Client
}

// NewOpts ...
type NewOpts struct {
	HTTPClient *http.Client
}

func (o *NewOpts) setDefaults() {
	if o.HTTPClient == nil {
		o.HTTPClient = http.DefaultClient
	}
}

// New returns new Client.
func New(opts NewOpts) *Client {
	opts.setDefaults()

	return &Client{
		client: opts.HTTPClient,
	}
}

// Servers returns ragemp servers.
func (c *Client) Servers(ctx context.Context) (Servers, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverListURL, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %w", err)
	}

	defer resp.Body.Close() //nolint:errcheck

	var respServers Servers

	err = json.NewDecoder(resp.Body).Decode(&respServers)
	if err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return respServers, nil
}
