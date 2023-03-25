// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package testuplink_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"common/memory"
	"common/testcontext"
	"common/testrand"
	"storx/private/testplanet"
	"storx/satellite/metabase"
	"uplink/private/testuplink"
)

func TestWithMaxSegmentSize(t *testing.T) {
	ctx := context.Background()
	_, ok := testuplink.GetMaxSegmentSize(ctx)
	require.False(t, ok)

	newCtx := testuplink.WithMaxSegmentSize(ctx, memory.KiB)
	segmentSize, ok := testuplink.GetMaxSegmentSize(newCtx)
	require.True(t, ok)
	require.EqualValues(t, memory.KiB, segmentSize)

}

func TestWithMaxSegmentSize_Upload(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 4, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		newCtx := testuplink.WithMaxSegmentSize(ctx, 10*memory.KiB)

		expectedData := testrand.Bytes(19 * memory.KiB)
		err := planet.Uplinks[0].Upload(newCtx, planet.Satellites[0], "super-bucket", "super-object", expectedData)
		require.NoError(t, err)

		data, err := planet.Uplinks[0].Download(newCtx, planet.Satellites[0], "super-bucket", "super-object")
		require.NoError(t, err)
		require.Equal(t, expectedData, data)

		// verify we have two segments instead of one
		objects, err := planet.Satellites[0].Metabase.DB.TestingAllCommittedObjects(ctx, planet.Uplinks[0].Projects[0].ID, "super-bucket")
		require.NoError(t, err)
		require.Len(t, objects, 1)

		segments, err := planet.Satellites[0].Metabase.DB.TestingAllObjectSegments(ctx, metabase.ObjectLocation{
			ProjectID:  planet.Uplinks[0].Projects[0].ID,
			BucketName: "super-bucket",
			ObjectKey:  objects[0].ObjectKey,
		})
		require.NoError(t, err)
		require.Equal(t, 2, len(segments))
	})
}
