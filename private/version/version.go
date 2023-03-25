// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

package version

import (
	"github.com/zeebo/errs"

	"common/useragent"
	"common/version"
)

// AppendVersionToUserAgent appends uplink product and version to user agent string.
//
// This doesn't work in test environment.
func AppendVersionToUserAgent(useragentStr string) (string, error) {
	version, err := version.FromBuild("uplink")
	if err != nil {
		return useragentStr, nil //nolint: nilerr // passthrough
	}
	entries := []useragent.Entry{}
	if len(useragentStr) > 0 {
		entries, err = useragent.ParseEntries([]byte(useragentStr))
		if err != nil {
			return "", errs.New("invalid user agent: %w", err)
		}
	}
	entries = append(entries, useragent.Entry{
		Product: "uplink",
		Version: version,
	})
	newUseragent, err := useragent.EncodeEntries(entries)
	if err != nil {
		return "", errs.New("unable to encode user agent entries: %w", err)
	}
	return string(newUseragent), nil
}
