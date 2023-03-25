// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package metaclient

import (
	"time"

	"common/storx"
)

// RawObjectItem represents raw object item from get request.
type RawObjectItem struct {
	Version            uint32
	Bucket             string
	EncryptedObjectKey []byte

	StreamID storx.StreamID

	Created  time.Time
	Modified time.Time
	Expires  time.Time

	PlainSize int64

	EncryptedMetadataNonce        storx.Nonce
	EncryptedMetadataEncryptedKey []byte
	EncryptedMetadata             []byte

	EncryptionParameters storx.EncryptionParameters
	RedundancyScheme     storx.RedundancyScheme
}

// RawObjectListItem represents raw object item from list objects request.
type RawObjectListItem struct {
	EncryptedObjectKey            []byte
	Version                       int32
	Status                        int32
	CreatedAt                     time.Time
	StatusAt                      time.Time
	ExpiresAt                     time.Time
	PlainSize                     int64
	EncryptedMetadataNonce        storx.Nonce
	EncryptedMetadataEncryptedKey []byte
	EncryptedMetadata             []byte
	StreamID                      storx.StreamID
	IsPrefix                      bool
}

// SegmentPosition the segment position within its parent object.
// It is an identifier for the segment.
type SegmentPosition struct {
	// PartNumber indicates the ordinal of the part within an object.
	// A part contains one or more segments.
	// PartNumber is defined by the user.
	// This is only relevant for multipart objects.
	// A non-multipart object only has one Part, and its number is 0.
	PartNumber int32
	// Index indicates the ordinal of this segment within a part.
	// Index is managed by uplink.
	// It is zero-indexed within each part.
	Index int32
}

// SegmentDownloadResponseInfo represents segment download information inline/remote.
type SegmentDownloadResponseInfo struct {
	SegmentID           storx.SegmentID
	EncryptedSize       int64
	EncryptedInlineData []byte
	Next                SegmentPosition
	Position            SegmentPosition
	PiecePrivateKey     storx.PiecePrivateKey

	SegmentEncryption SegmentEncryption
}

// SegmentEncryption represents segment encryption key and nonce.
type SegmentEncryption struct {
	EncryptedKeyNonce storx.Nonce
	EncryptedKey      storx.EncryptedPrivateKey
}

var (
	// ErrNoPath is an error class for using empty path.
	ErrNoPath = storx.ErrNoPath

	// ErrObjectNotFound is an error class for non-existing object.
	ErrObjectNotFound = storx.ErrObjectNotFound
)

// Object contains information about a specific object.
type Object struct {
	Version  uint32
	Bucket   Bucket
	Path     string
	IsPrefix bool

	Metadata map[string]string

	ContentType string
	Created     time.Time
	Modified    time.Time
	Expires     time.Time

	Stream
}

// Stream is information about an object stream.
type Stream struct {
	ID storx.StreamID

	// Size is the total size of the stream in bytes
	Size int64

	// SegmentCount is the number of segments
	SegmentCount int64
	// FixedSegmentSize is the size of each segment,
	// when all segments have the same size. It is -1 otherwise.
	FixedSegmentSize int64

	// RedundancyScheme specifies redundancy strategy used for this stream
	storx.RedundancyScheme
	// EncryptionParameters specifies encryption strategy used for this stream
	storx.EncryptionParameters

	LastSegment LastSegment // TODO: remove
}

// LastSegment contains info about last segment.
type LastSegment struct {
	Size              int64
	EncryptedKeyNonce storx.Nonce
	EncryptedKey      storx.EncryptedPrivateKey
}

var (
	// ErrBucket is an error class for general bucket errors.
	ErrBucket = storx.ErrBucket

	// ErrNoBucket is an error class for using empty bucket name.
	ErrNoBucket = storx.ErrNoBucket

	// ErrBucketNotFound is an error class for non-existing bucket.
	ErrBucketNotFound = storx.ErrBucketNotFound
)

// Bucket contains information about a specific bucket.
type Bucket struct {
	Name        string
	Created     time.Time
	Attribution string
}

// ListDirection specifies listing direction.
type ListDirection = storx.ListDirection

const (
	// Forward lists forwards from cursor, including cursor.
	Forward = storx.Forward
	// After lists forwards from cursor, without cursor.
	After = storx.After
)

// ListOptions lists objects.
type ListOptions struct {
	Prefix                storx.Path
	Cursor                storx.Path // Cursor is relative to Prefix, full path is Prefix + Cursor
	CursorEnc             []byte
	Delimiter             rune
	Recursive             bool
	Direction             ListDirection
	Limit                 int
	IncludeCustomMetadata bool
	IncludeSystemMetadata bool
	Status                int32
}

// NextPage returns options for listing the next page.
func (opts ListOptions) NextPage(list ObjectList) ListOptions {
	if !list.More || len(list.Items) == 0 {
		return ListOptions{}
	}

	return ListOptions{
		Prefix:                opts.Prefix,
		CursorEnc:             list.Cursor,
		Delimiter:             opts.Delimiter,
		Recursive:             opts.Recursive,
		IncludeSystemMetadata: opts.IncludeSystemMetadata,
		IncludeCustomMetadata: opts.IncludeCustomMetadata,
		Direction:             After,
		Limit:                 opts.Limit,
		Status:                opts.Status,
	}
}

// ObjectList is a list of objects.
type ObjectList struct {
	Bucket string
	Prefix string
	More   bool
	Cursor []byte

	// Items paths are relative to Prefix
	// To get the full path use list.Prefix + list.Items[0].Path
	Items []Object
}

// BucketList is a list of buckets.
type BucketList struct {
	More  bool
	Items []Bucket
}

// BucketListOptions lists objects.
type BucketListOptions struct {
	Cursor    string
	Direction ListDirection
	Limit     int
}

// NextPage returns options for listing the next page.
func (opts BucketListOptions) NextPage(list BucketList) BucketListOptions {
	if !list.More || len(list.Items) == 0 {
		return BucketListOptions{}
	}

	return BucketListOptions{
		Cursor:    list.Items[len(list.Items)-1].Name,
		Direction: After,
		Limit:     opts.Limit,
	}
}
