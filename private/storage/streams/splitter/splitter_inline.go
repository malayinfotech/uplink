// Copyright (C) 2023 Storx Labs, Inc.
// See LICENSE for copying information.

package splitter

import (
	"bytes"
	"io"

	"common/storx"
	"uplink/private/metaclient"
)

type splitterInline struct {
	position   metaclient.SegmentPosition
	encryption metaclient.SegmentEncryption
	encParams  storx.EncryptionParameters
	contentKey *storx.Key

	encData   []byte
	plainSize int64
}

func (s *splitterInline) Begin() metaclient.BatchItem {
	return &metaclient.MakeInlineSegmentParams{
		StreamID:            nil, // set by the stream batcher
		Position:            s.position,
		Encryption:          s.encryption,
		EncryptedInlineData: s.encData,
		PlainSize:           s.plainSize,
		EncryptedTag:        nil, // set by the segment tracker
	}
}

func (s *splitterInline) Position() metaclient.SegmentPosition { return s.position }
func (s *splitterInline) Inline() bool                         { return true }
func (s *splitterInline) Reader() io.Reader                    { return bytes.NewReader(s.encData) }
func (s *splitterInline) DoneReading(err error)                {}

func (s *splitterInline) EncryptETag(eTag []byte) ([]byte, error) {
	return encryptETag(eTag, s.encParams.CipherSuite, s.contentKey)
}

func (s *splitterInline) Finalize() *SegmentInfo {
	return &SegmentInfo{
		Encryption:    s.encryption,
		PlainSize:     s.plainSize,
		EncryptedSize: int64(len(s.encData)),
	}
}
