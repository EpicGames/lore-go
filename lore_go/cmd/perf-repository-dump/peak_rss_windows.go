// Copyright Epic Games, Inc. All Rights Reserved.

//go:build windows

package main

// processPeakRssBytes is a no-op stub on Windows. The perf test is targeted
// at macOS for the actual measurement runs; Windows just needs the package
// to compile so the unit-test job can build the whole module.
func processPeakRssBytes() uint64 {
	return 0
}
