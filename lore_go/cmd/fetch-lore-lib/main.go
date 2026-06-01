// Copyright Epic Games, Inc. All Rights Reserved.
//
// This tool installs the Lore native library next to your application binary.
// It is typically invoked via "go generate". By default the library is placed
// in the current working directory (next to the file containing the
// //go:generate directive). Use -o to specify a different directory.
//
// Source priority:
//  1. If LORE_LIB_PATH is set, copy the library from that file path. This
//     matches the runtime semantic in lore_go/native (set the env var to a
//     specific .so/.dylib/.dll file and both fetch-time and runtime use it).
//  2. Otherwise download from LORE_RELEASE_BASE_URL, falling back to the
//     build-time loreReleaseBaseURL constant in version.go (set by
//     generator/generate.py from $LORE_RELEASE_BASE_URL). The URL is
//     constructed as <base>/<versionTag>/<artifactName>.

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type platformInfo struct {
	artifactName string // versioned filename, e.g. liblore-v0.8.1-macos-arm64.dylib
	localName    string // platform default name, e.g. liblore.dylib
}

// resolveVersionTag returns the full version tag used as the release
// directory and embedded in the artifact filenames. Mirrors the layout
// used by release.yml / release-nightly.yml when uploading binaries:
//
//	stable release             -> vX.Y.Z
//	nightly main, no name      -> vX.Y.Z-nightly-REV
//	nightly with name          -> vX.Y.Z-nightly-REV-NAME
//	nightly non-main, no name  -> vX.Y.Z-nightly-REV-BRANCH
//
// loreVersion may or may not already include a leading 'v' or a trailing
// '-nightly' segment (Cargo.toml on main is e.g. "0.8.2-nightly"); both
// shapes are normalized to a single canonical form.
func resolveVersionTag() string {
	version := strings.TrimPrefix(loreVersion, "v")

	if loreRevision == "" {
		return "v" + version
	}

	if !strings.Contains(version, "-nightly") {
		version += "-nightly"
	}
	tag := fmt.Sprintf("v%s-%s", version, loreRevision)
	if loreName != "" {
		tag += "-" + loreName
	} else if loreBranch != "" && loreBranch != "main" {
		tag += "-" + loreBranch
	}
	return tag
}

func resolvePlatform(goos, goarch, versionTag string) (*platformInfo, error) {
	switch {
	case goos == "windows" && goarch == "amd64":
		return &platformInfo{
			artifactName: fmt.Sprintf("lore-%s.dll", versionTag),
			localName:    "lore.dll",
		}, nil
	case goos == "darwin" && goarch == "arm64":
		return &platformInfo{
			artifactName: fmt.Sprintf("liblore-%s-macos-arm64.dylib", versionTag),
			localName:    "liblore.dylib",
		}, nil
	case goos == "linux" && goarch == "arm64":
		return &platformInfo{
			artifactName: fmt.Sprintf("liblore-%s-linux-arm64.so", versionTag),
			localName:    "liblore.so",
		}, nil
	case goos == "linux" && goarch == "amd64":
		return &platformInfo{
			artifactName: fmt.Sprintf("liblore-%s-linux-x86_64.so", versionTag),
			localName:    "liblore.so",
		}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer out.Close()
	written, err := io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	fmt.Printf("  Copied %d bytes\n", written)
	return nil
}

func download(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("  Written %d bytes\n", written)
	return nil
}

func versionFilePath(dir string) string {
	return filepath.Join(dir, ".lore-version")
}

func checkVersion(dir, expectedVersionTag string) bool {
	data, err := os.ReadFile(versionFilePath(dir))
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == expectedVersionTag
}

func writeVersion(dir, versionTag string) error {
	return os.WriteFile(versionFilePath(dir), []byte(versionTag+"\n"), 0o644)
}

func main() {
	outputDir := flag.String("o", ".", "output directory for the native library")
	targetOS := flag.String("os", runtime.GOOS, "target OS (e.g. linux, darwin, windows)")
	targetArch := flag.String("arch", runtime.GOARCH, "target architecture (e.g. amd64, arm64)")
	flag.Parse()

	versionTag := resolveVersionTag()

	platform, err := resolvePlatform(*targetOS, *targetArch, versionTag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	destPath := filepath.Join(*outputDir, platform.localName)

	// Skip if the destination already has the right version
	if _, err := os.Stat(destPath); err == nil {
		if checkVersion(*outputDir, versionTag) {
			fmt.Printf("Lore library %s already present at %s\n", versionTag, destPath)
			return
		}
		fmt.Printf("Lore library version mismatch or missing marker, refreshing\n")
		os.Remove(destPath)
	}

	if libPath := os.Getenv("LORE_LIB_PATH"); libPath != "" {
		if _, err := os.Stat(libPath); err != nil {
			fmt.Fprintf(os.Stderr,
				"Error: LORE_LIB_PATH=%s: %v\n", libPath, err)
			os.Exit(1)
		}
		fmt.Printf("Copying Lore library from %s to %s\n", libPath, destPath)
		if err := copyFile(libPath, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Copy failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		baseURL := os.Getenv("LORE_RELEASE_BASE_URL")
		if baseURL == "" {
			baseURL = loreReleaseBaseURL
		}
		if baseURL == "" {
			fmt.Fprintln(os.Stderr,
				"Error: no Lore release base URL configured. Set LORE_RELEASE_BASE_URL or rebuild the SDK with $LORE_RELEASE_BASE_URL set.")
			os.Exit(1)
		}
		url := fmt.Sprintf("%s/%s/%s",
			strings.TrimRight(baseURL, "/"),
			versionTag, platform.artifactName)

		fmt.Printf("Downloading Lore %s for %s/%s...\n", versionTag, *targetOS, *targetArch)
		fmt.Printf("  URL: %s\n", url)
		fmt.Printf("  Destination: %s\n", destPath)

		if err := download(url, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}
	}

	if err := writeVersion(*outputDir, versionTag); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write version marker: %v\n", err)
	}

	fmt.Printf("Successfully obtained Lore library %s\n", versionTag)
}
