// Copyright (C) 2020 Storx Labs, Inc.
// See LICENSE for copying information.

package testuplink

import (
	"context"

	"common/memory"
	"uplink/private/eestream/scheduler"
)

type segmentSizeKey struct{}

type plainSizeKey struct{}

type listLimitKey struct{}

type concurrentSegmentUploadsConfigKey struct{}

// WithMaxSegmentSize creates context with max segment size for testing purposes.
//
// Created context needs to be used with uplink.OpenProject to manipulate default
// segment size.
func WithMaxSegmentSize(ctx context.Context, segmentSize memory.Size) context.Context {
	return context.WithValue(ctx, segmentSizeKey{}, segmentSize)
}

// GetMaxSegmentSize returns max segment size from context if exists.
func GetMaxSegmentSize(ctx context.Context) (memory.Size, bool) {
	segmentSize, ok := ctx.Value(segmentSizeKey{}).(memory.Size)
	return segmentSize, ok
}

// WithoutPlainSize creates context with information that segment plain size shouldn't be sent.
// Only for testing purposes.
func WithoutPlainSize(ctx context.Context) context.Context {
	return context.WithValue(ctx, plainSizeKey{}, true)
}

// IsWithoutPlainSize returns true if information about not sending segment plain size exists in context.
// Only for testing purposes.
func IsWithoutPlainSize(ctx context.Context) bool {
	withoutPlainSize, _ := ctx.Value(plainSizeKey{}).(bool)
	return withoutPlainSize
}

// WithListLimit creates context with information about list limit that will be used with request.
// Only for testing purposes.
func WithListLimit(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, listLimitKey{}, limit)
}

// GetListLimit returns value for list limit if exists in context.
// Only for testing purposes.
func GetListLimit(ctx context.Context) int {
	limit, _ := ctx.Value(listLimitKey{}).(int)
	return limit
}

// ConcurrentSegmentUploadsConfig is the configuration for concurrent
// segment uploads using the new upload codepath.
type ConcurrentSegmentUploadsConfig struct {
	// SchedulerOptions are the options for the scheduler used to place limits
	// on the amount of concurrent piece limits per-upload, across all
	// segments.
	SchedulerOptions scheduler.Options

	// LongTailMargin represents the maximum number of piece uploads beyond the
	// optimal threshold that will be uploaded for a given segment. Once an
	// upload has reached the optimal threshold, the remaining piece uploads
	// are cancelled.
	LongTailMargin int
}

// DefaultConcurrentSegmentUploadsConfig returns the default ConcurrentSegmentUploadsConfig.
func DefaultConcurrentSegmentUploadsConfig() ConcurrentSegmentUploadsConfig {
	return ConcurrentSegmentUploadsConfig{
		SchedulerOptions: scheduler.Options{
			MaximumConcurrent: 200,
		},
		LongTailMargin: 15,
	}
}

// WithConcurrentSegmentUploadsDefaultConfig creates a context that enables the
// new concurrent segment upload codepath for testing purposes using the
// default configuration.
//
// The context needs to be used with uplink.OpenProject to have effect.
func WithConcurrentSegmentUploadsDefaultConfig(ctx context.Context) context.Context {
	return WithConcurrentSegmentUploadsConfig(ctx, DefaultConcurrentSegmentUploadsConfig())
}

// WithConcurrentSegmentUploadsConfig creates a context that enables the
// new concurrent segment upload codepath for testing purposes using the
// given scheduler options.
//
// The context needs to be used with uplink.OpenProject to have effect.
func WithConcurrentSegmentUploadsConfig(ctx context.Context, config ConcurrentSegmentUploadsConfig) context.Context {
	return context.WithValue(ctx, concurrentSegmentUploadsConfigKey{}, config)
}

// GetConcurrentSegmentUploadsConfig returns the scheduler options to
// use with the new concurrent segment upload codepath, or nil if no scheduler
// options have been set.
func GetConcurrentSegmentUploadsConfig(ctx context.Context) *ConcurrentSegmentUploadsConfig {
	if config, ok := ctx.Value(concurrentSegmentUploadsConfigKey{}).(ConcurrentSegmentUploadsConfig); ok {
		return &config
	}
	return nil
}
