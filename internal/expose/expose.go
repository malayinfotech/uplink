// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package expose

import (
	"context"
	_ "unsafe" // for go:linkname

	"common/grant"
	"common/macaroon"
	"common/rpc"
	"common/rpc/rpcpool"
	"uplink"
)

// ConfigSetConnectionPool exposes Config.setConnectionPool.
//
//go:linkname ConfigSetConnectionPool uplink.config_setConnectionPool
func ConfigSetConnectionPool(*uplink.Config, *rpcpool.Pool)

// ConfigSetSatelliteConnectionPool exposes Config.setSatelliteConnectionPool.
//
//go:linkname ConfigSetSatelliteConnectionPool uplink.config_setSatelliteConnectionPool
func ConfigSetSatelliteConnectionPool(*uplink.Config, *rpcpool.Pool)

// ConfigGetDialer exposes Config.getDialer.
//
//go:linkname ConfigGetDialer uplink.config_getDialer
//nolint:revive
func ConfigGetDialer(uplink.Config, context.Context) (rpc.Dialer, error)

// ConfigSetMaximumBufferSize exposes Config.setMaximumBufferSize.
//
//go:linkname ConfigSetMaximumBufferSize uplink.config_setMaximumBufferSize
func ConfigSetMaximumBufferSize(*uplink.Config, int)

// ConfigDisableObjectKeyEncryption exposes Config.disableObjectKeyEncryption.
//
//go:linkname ConfigDisableObjectKeyEncryption uplink.config_disableObjectKeyEncryption
func ConfigDisableObjectKeyEncryption(config *uplink.Config)

// AccessGetAPIKey exposes Access.getAPIKey.
//
//go:linkname AccessGetAPIKey uplink.access_getAPIKey
func AccessGetAPIKey(*uplink.Access) *macaroon.APIKey

// AccessGetEncAccess exposes Access.getEncAccess.
//
//go:linkname AccessGetEncAccess uplink.access_getEncAccess
func AccessGetEncAccess(*uplink.Access) *grant.EncryptionAccess

// ConfigRequestAccessWithPassphraseAndConcurrency exposes Config.requestAccessWithPassphraseAndConcurrency.
//
//nolint:revive
//go:linkname ConfigRequestAccessWithPassphraseAndConcurrency uplink.config_requestAccessWithPassphraseAndConcurrency
func ConfigRequestAccessWithPassphraseAndConcurrency(config uplink.Config, ctx context.Context, satelliteAddress, apiKey, passphrase string, concurrency uint8) (*uplink.Access, error)
