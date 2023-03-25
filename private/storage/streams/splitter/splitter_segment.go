// Copyright (C) 2023 Storx Labs, Inc.
// See LICENSE for copying information.

package splitter

import (
	"io"

	"common/encryption"
	"common/storx"
	"uplink/private/metaclient"
)

type splitterSegment struct {
	position   metaclient.SegmentPosition
	encryption metaclient.SegmentEncryption
	encParams  storx.EncryptionParameters
	contentKey *storx.Key

	maxSegmentSize int64
	encTransformer encryption.Transformer
	encBuf         *encryptedBuffer
}

func (s *splitterSegment) Begin() metaclient.BatchItem {
	return &metaclient.BeginSegmentParams{
		StreamID:      nil, // set by the stream batcher
		Position:      s.position,
		MaxOrderLimit: s.maxSegmentSize,
	}
}

func (s *splitterSegment) Position() metaclient.SegmentPosition { return s.position }
func (s *splitterSegment) Inline() bool                         { return false }
func (s *splitterSegment) Reader() io.Reader                    { return s.encBuf.Reader() }
func (s *splitterSegment) DoneReading(err error)                { s.encBuf.DoneReading(err) }

func (s *splitterSegment) EncryptETag(eTag []byte) ([]byte, error) {
	return encryptETag(eTag, s.encParams.CipherSuite, s.contentKey)
}

func (s *splitterSegment) Finalize() *SegmentInfo {
	plainSize := s.encBuf.PlainSize()
	return &SegmentInfo{
		Encryption:    s.encryption,
		PlainSize:     plainSize,
		EncryptedSize: encryption.CalcTransformerEncryptedSize(plainSize, s.encTransformer),
	}
}
