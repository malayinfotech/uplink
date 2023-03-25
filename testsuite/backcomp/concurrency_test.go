// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package backcomp_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"common/memory"
	"common/testcontext"
	"common/testrand"
	"storx/private/testplanet"
	"uplink"
	"uplink/backcomp"
)

func TestRequestAccessWithPassphraseAndConcurrency(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount:   1,
		StorageNodeCount: 0,
		UplinkCount:      1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		projectInfo := planet.Uplinks[0].Projects[0]

		uplinkConfig := uplink.Config{}

		// create access with custom concurrency, create a bucket, upload a file, check listing
		customAccess, err := backcomp.RequestAccessWithPassphraseAndConcurrency(ctx, uplinkConfig, projectInfo.Satellite.URL(), projectInfo.APIKey, "mypassphrase", 4)
		require.NoError(t, err)

		project, err := uplinkConfig.OpenProject(ctx, customAccess)
		require.NoError(t, err)
		defer ctx.Check(project.Close)

		bucket, err := project.CreateBucket(ctx, "bucket-abcd")
		require.NoError(t, err)

		upload, err := project.UploadObject(ctx, bucket.Name, "test.dat", nil)
		require.NoError(t, err)

		source := bytes.NewBuffer(testrand.Bytes(1 * memory.KiB))
		_, err = io.Copy(upload, source)
		require.NoError(t, err)

		err = upload.Commit()
		require.NoError(t, err)

		objects := project.ListObjects(ctx, bucket.Name, nil)
		require.True(t, objects.Next())
		require.NoError(t, objects.Err())

		download, err := project.DownloadObject(ctx, bucket.Name, "test.dat", nil)
		require.NoError(t, err)
		defer ctx.Check(download.Close)

		// try to use access with default concurrency to download object in a bucket, should fail
		standardAccess, err := uplinkConfig.RequestAccessWithPassphrase(ctx, projectInfo.Satellite.URL(), projectInfo.APIKey, "mypassphrase")
		require.NoError(t, err)

		differentProject, err := uplinkConfig.OpenProject(ctx, standardAccess)
		require.NoError(t, err)
		defer ctx.Check(differentProject.Close)

		_, err = differentProject.DownloadObject(ctx, bucket.Name, "test.dat", nil)
		require.Error(t, err)
	})
}
