// Copyright (C) 2022 Storx Labs, Inc.
// See LICENSE for copying information

package piecestore_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/memory"
	"common/pb"
	"common/pkcrypto"
	"common/signing"
	"common/storx"
	"common/testcontext"
	"common/testrand"
	"storx/private/testplanet"
)

func TestUploadStreamClosing(t *testing.T) {
	testplanet.Run(t, testplanet.Config{
		SatelliteCount: 1, StorageNodeCount: 1, UplinkCount: 1,
	}, func(t *testing.T, ctx *testcontext.Context, planet *testplanet.Planet) {
		client, err := planet.Uplinks[0].DialPiecestore(ctx, planet.StorageNodes[0])
		require.NoError(t, err)
		defer ctx.Check(client.Close)

		tts := []struct {
			pieceID       storx.PieceID
			contentLength memory.Size
			action        pb.PieceAction
			err           string
		}{
			{ // should successfully store data
				pieceID:       storx.PieceID{1},
				contentLength: 32,
				action:        pb.PieceAction_PUT,
				err:           "",
			},
			{ // should err with piece ID not specified
				pieceID:       storx.PieceID{},
				contentLength: 32,
				action:        pb.PieceAction_PUT,
				err:           "missing piece id",
			},
		}

		for i := 0; i < 1000; i++ {
			tt := tts[i%len(tts)]
			tt.pieceID = storx.PieceID{byte(i / 256), byte(i % 256), 1}
			if i%2 == 1 {
				tt.pieceID = storx.PieceID{}
			}

			data := testrand.Bytes(tt.contentLength)
			serialNumber := testrand.SerialNumber()

			orderLimit, piecePrivateKey := generateOrderLimit(
				t,
				planet.Satellites[0].ID(),
				planet.StorageNodes[0].ID(),
				tt.pieceID,
				tt.action,
				serialNumber,
				24*time.Hour,
				24*time.Hour,
				int64(len(data)),
			)
			signer := signing.SignerFromFullIdentity(planet.Satellites[0].Identity)
			orderLimit, err = signing.SignOrderLimit(ctx, signer, orderLimit)
			require.NoError(t, err)

			pieceHash, err := client.UploadReader(ctx, orderLimit, piecePrivateKey, bytes.NewReader(data))
			if tt.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			} else {
				require.NoError(t, err)

				expectedHash := pkcrypto.SHA256Hash(data)
				assert.Equal(t, expectedHash, pieceHash.Hash)

				signee := signing.SignerFromFullIdentity(planet.StorageNodes[0].Identity)
				require.NoError(t, signing.VerifyPieceHashSignature(ctx, signee, pieceHash))
			}
		}
	})
}

func generateOrderLimit(t *testing.T, satellite storx.NodeID, storageNode storx.NodeID, pieceID storx.PieceID, action pb.PieceAction, serialNumber storx.SerialNumber, pieceExpiration, orderExpiration time.Duration, limit int64) (*pb.OrderLimit, storx.PiecePrivateKey) {
	piecePublicKey, piecePrivateKey, err := storx.NewPieceKey()
	require.NoError(t, err)

	now := time.Now()
	return &pb.OrderLimit{
		SatelliteId:     satellite,
		UplinkPublicKey: piecePublicKey,
		StorageNodeId:   storageNode,
		PieceId:         pieceID,
		Action:          action,
		SerialNumber:    serialNumber,
		OrderCreation:   time.Now(),
		OrderExpiration: now.Add(orderExpiration),
		PieceExpiration: now.Add(pieceExpiration),
		Limit:           limit,
	}, piecePrivateKey
}
