// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package transport

import (
	"context"

	"common/rpc/rpcpool"
	"uplink"
	"uplink/internal/expose"
)

// SetConnectionPool configures connection pool on the passed in config. If
// argument pool is nil, it will clear the pool on the config.
func SetConnectionPool(ctx context.Context, config *uplink.Config, pool *rpcpool.Pool) error {
	expose.ConfigSetConnectionPool(config, pool)
	return nil
}

// SetSatelliteConnectionPool configures connection pool (for satellite) on the passed in config. If
// argument pool is nil, it will clear the pool on the config.
func SetSatelliteConnectionPool(ctx context.Context, config *uplink.Config, pool *rpcpool.Pool) error {
	expose.ConfigSetSatelliteConnectionPool(config, pool)
	return nil
}

// SetMaximumBufferSize sets maximumBufferSize in config.
func SetMaximumBufferSize(config *uplink.Config, maximumBufferSize int) {
	expose.ConfigSetMaximumBufferSize(config, maximumBufferSize)
}
