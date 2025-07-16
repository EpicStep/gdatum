// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package sql

import "github.com/huandu/go-sqlbuilder"

// Builder ...
type Builder interface {
	BuildWithFlavor(flavor sqlbuilder.Flavor, initialArg ...any) (string, []any)
}

// Build query with ClickHouse flavor.
func Build(b Builder) (string, []any) {
	return b.BuildWithFlavor(sqlbuilder.ClickHouse)
}
