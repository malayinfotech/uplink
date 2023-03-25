// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package testsuite_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"common/testcontext"
	"storx/private/testplanet"
	"storx/satellite"
	"uplink"
)

func TestListBuckets_EmptyProject(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		project := openProject(t, ctx, planet)
		defer ctx.Check(project.Close)

		list := listBuckets(ctx, t, project, nil)
		assertNoNextBucket(t, list)
	})
}

func TestListBuckets_SingleBucket(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		project := openProject(t, ctx, planet)
		defer ctx.Check(project.Close)

		createBucket(t, ctx, project, "testbucket")
		defer func() {
			_, err := project.DeleteBucket(ctx, "testbucket")
			require.NoError(t, err)
		}()

		list := listBuckets(ctx, t, project, nil)

		assert.True(t, list.Next())
		require.NoError(t, list.Err())
		require.NotNil(t, list.Item())
		require.Equal(t, "testbucket", list.Item().Name)
		require.WithinDuration(t, time.Now(), list.Item().Created, 10*time.Second)

		assertNoNextBucket(t, list)
	})
}

func TestListBuckets_TwoBuckets(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		project := openProject(t, ctx, planet)
		defer ctx.Check(project.Close)

		expectedBuckets := map[string]bool{
			"testbucket1": true,
			"testbucket2": true,
		}

		for bucket := range expectedBuckets {
			bucket := bucket
			createBucket(t, ctx, project, bucket)
			defer func() {
				_, err := project.DeleteBucket(ctx, bucket)
				require.NoError(t, err)
			}()
		}

		list := listBuckets(ctx, t, project, nil)

		more := list.Next()
		require.True(t, more)
		require.NoError(t, list.Err())
		require.NotNil(t, list.Item())
		require.WithinDuration(t, time.Now(), list.Item().Created, 10*time.Second)
		delete(expectedBuckets, list.Item().Name)

		more = list.Next()
		require.True(t, more)
		require.NoError(t, list.Err())
		require.NotNil(t, list.Item())
		require.WithinDuration(t, time.Now(), list.Item().Created, 10*time.Second)
		delete(expectedBuckets, list.Item().Name)

		require.Empty(t, expectedBuckets)
		assertNoNextBucket(t, list)
	})
}

func TestListBuckets_Cursor(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		project := openProject(t, ctx, planet)
		defer ctx.Check(project.Close)

		expectedBuckets := map[string]bool{
			"testbucket1": true,
			"testbucket2": true,
		}

		for bucket := range expectedBuckets {
			bucket := bucket
			createBucket(t, ctx, project, bucket)
			defer func() {
				_, err := project.DeleteBucket(ctx, bucket)
				require.NoError(t, err)
			}()
		}

		list := listBuckets(ctx, t, project, nil)

		// get the first list item and make it a cursor for the next list request
		more := list.Next()
		require.True(t, more)
		require.NoError(t, list.Err())
		delete(expectedBuckets, list.Item().Name)
		cursor := list.Item().Name

		// list again with cursor set to the first item from previous list request
		list = listBuckets(ctx, t, project, &uplink.ListBucketsOptions{Cursor: cursor})

		// expect the second item as the first item in this new list request
		more = list.Next()
		require.True(t, more)
		require.NoError(t, list.Err())
		require.NotNil(t, list.Item())
		require.WithinDuration(t, time.Now(), list.Item().Created, 10*time.Second)
		delete(expectedBuckets, list.Item().Name)

		require.Empty(t, expectedBuckets)
		assertNoNextBucket(t, list)
	})
}

type satelliteDBWithBucketsListLimit struct {
	limit int
	satellite.DB
}

func TestListBuckets_AutoPaging(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
		Reconfigure: testplanet.Reconfigure{
			SatelliteDB: func(log *zap.Logger, index int, satellitedb satellite.DB) (satellite.DB, error) {
				return &satelliteDBWithBucketsListLimit{2, satellitedb}, nil
			},
		},
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		project := openProject(t, ctx, planet)
		defer ctx.Check(project.Close)

		totalBuckets := 5
		expectedBuckets := map[string]bool{}

		for i := 0; i < totalBuckets; i++ {
			bucketName := fmt.Sprintf("bucket%d", i)
			expectedBuckets[bucketName] = true
			createBucket(t, ctx, project, bucketName)

			defer func() {
				_, err := project.DeleteBucket(ctx, bucketName)
				require.NoError(t, err)
			}()
		}

		list := listBuckets(ctx, t, project, nil)

		var ok bool
		for list.Next() {
			bucket := list.Item()

			_, ok = expectedBuckets[bucket.Name]
			require.True(t, ok)

			delete(expectedBuckets, bucket.Name)
		}

		require.NoError(t, list.Err())
		require.Equal(t, 0, len(expectedBuckets))
	})
}

func listBuckets(ctx context.Context, t *testing.T, project *uplink.Project, options *uplink.ListBucketsOptions) *uplink.BucketIterator {
	list := project.ListBuckets(ctx, options)
	require.NoError(t, list.Err())
	require.Nil(t, list.Item())
	return list
}

func assertNoNextBucket(t *testing.T, list *uplink.BucketIterator) {
	require.False(t, list.Next())
	require.NoError(t, list.Err())
	require.Nil(t, list.Item())
}
