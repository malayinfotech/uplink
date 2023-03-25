// Copyright (C) 2023 Storx Labs, Inc.
// See LICENSE for copying information.

package splitter

import (
	"common/encryption"
	"common/storx"
	"uplink/private/metaclient"
)

// TODO: move it to separate package?
func encryptETag(etag []byte, cipherSuite storx.CipherSuite, contentKey *storx.Key) ([]byte, error) {
	etagKey, err := encryption.DeriveKey(contentKey, "storx-etag-v1")
	if err != nil {
		return nil, err
	}

	encryptedETag, err := encryption.Encrypt(etag, cipherSuite, etagKey, &storx.Nonce{})
	if err != nil {
		return nil, err
	}

	return encryptedETag, nil
}

func nonceForPosition(position metaclient.SegmentPosition) (storx.Nonce, error) {
	var nonce storx.Nonce
	inc := (int64(position.PartNumber) << 32) | (int64(position.Index) + 1)
	_, err := encryption.Increment(&nonce, inc)
	return nonce, err
}
