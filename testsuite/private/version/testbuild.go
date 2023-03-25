// Copyright (C) 2021 Storx Labs, Inc.
// See LICENSE for copying information.

//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
)

func main() {
	useragent := os.Args[1]
	// useragentStr, err := version.AppendVersionToUserAgent(useragent)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }

	fmt.Printf("%#v", useragent)
}
