// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package metainfo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"common/memory"
	"common/testcontext"
	"common/testrand"
	"storx/private/testplanet"
	"uplink/private/metaclient"
)

func TestGetObject_RedundancySchemePerSegment(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 4, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		err := planet.Uplinks[0].Upload(ctx, planet.Satellites[0], "super-bucket", "super-object", testrand.Bytes(10*memory.KiB))
		require.NoError(t, err)

		objects, err := planet.Satellites[0].Metabase.DB.TestingAllCommittedObjects(ctx, planet.Uplinks[0].Projects[0].ID, "super-bucket")
		require.NoError(t, err)
		require.Len(t, objects, 1)

		apiKey := planet.Uplinks[0].APIKey[planet.Satellites[0].ID()]
		metainfoClient, err := planet.Uplinks[0].DialMetainfo(ctx, planet.Satellites[0], apiKey)
		require.NoError(t, err)
		defer ctx.Check(metainfoClient.Close)

		// RedundancySchemePerSegment == false means that GetObject SHOULD
		// return redundancy scheme
		object, err := metainfoClient.GetObject(ctx, metaclient.GetObjectParams{
			Bucket:                     []byte("super-bucket"),
			EncryptedObjectKey:         []byte(objects[0].ObjectKey),
			RedundancySchemePerSegment: false,
		})
		require.NoError(t, err)
		require.False(t, object.RedundancyScheme.IsZero())

		// RedundancySchemePerSegment == true means that GetObject SHOULDN'T
		// return redundancy scheme
		object, err = metainfoClient.GetObject(ctx, metaclient.GetObjectParams{
			Bucket:                     []byte("super-bucket"),
			EncryptedObjectKey:         []byte(objects[0].ObjectKey),
			RedundancySchemePerSegment: true,
		})
		require.NoError(t, err)
		require.True(t, object.RedundancyScheme.IsZero())

	})
}
