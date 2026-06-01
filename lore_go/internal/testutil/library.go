// Copyright Epic Games, Inc. All Rights Reserved.

// Package testutil provides helpers for SDK tests to locate the native library.
package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// SetLibraryPath sets LORE_LIB_PATH to the native library in the SDK's
// lib/ directory. Call this from TestMain before running tests.
func SetLibraryPath() error {
	// Find the SDK root relative to this file:
	// this file is at lore_go/internal/testutil/library.go
	// native libs are at lore_go/lib/<platform>/<libname>
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to determine source file location")
	}
	sdkDir := filepath.Join(filepath.Dir(filename), "..", "..")

	var libPath string
	switch {
	case runtime.GOOS == "darwin":
		libPath = filepath.Join(sdkDir, "lib", "darwin", "liblore.dylib")
	case runtime.GOOS == "linux" && runtime.GOARCH == "arm64":
		libPath = filepath.Join(sdkDir, "lib", "linux_arm64", "liblore.so")
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		libPath = filepath.Join(sdkDir, "lib", "linux_amd64", "liblore.so")
	case runtime.GOOS == "windows":
		libPath = filepath.Join(sdkDir, "lib", "windows", "lore.dll")
	default:
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("native library not found at %s: %w", libPath, err)
	}

	os.Setenv("LORE_LIB_PATH", libPath)
	return nil
}
