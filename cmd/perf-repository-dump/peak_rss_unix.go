// Copyright Epic Games, Inc. All Rights Reserved.

//go:build !windows

package main

import "syscall"

// processPeakRssBytes returns the high-water-mark RSS for this process via
// getrusage. On macOS Maxrss is in bytes (which is what we want — we're
// macOS-targeted). On Linux it's in KB, but we don't run there.
func processPeakRssBytes() uint64 {
	var r syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &r); err != nil {
		return 0
	}
	return uint64(r.Maxrss)
}
