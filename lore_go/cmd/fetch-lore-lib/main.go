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
//     constructed as <base>/<dirVersion>/<artifactName>.
//
// How the artifact is named and retrieved is controlled by two build-time
// constants in version.go (see fetch_native_version.ji):
//   - loreArtifactNaming: "short" (liblore-<tag>-linux-x86_64.so, lore-<tag>.dll)
//     or "triple" (rust target triples, e.g. -x86_64-unknown-linux-gnu).
//   - loreArtifactFormat: "direct" (download the raw library file) or "archive"
//     (download a .tar.gz/.zip bundle and extract the library member).

package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
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
	artifactName string // remote filename to download, e.g. liblore-v0.8.1-macos-arm64.dylib
	localName    string // platform default name, e.g. liblore.dylib
	archiveKind  string // "tar.gz", "zip", or "" for a direct (unbundled) download
	member       string // when archived, the file within the archive to extract
}

// sanitizeBranch normalizes a branch name for use in a release path component:
// '/' and any char outside [A-Za-z0-9._-] become '-', then runs of '-' collapse.
func sanitizeBranch(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := b.String()
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	return out
}

// buildConfig is the set of build-time values baked into version.go that
// determine which artifact the fetcher retrieves and from where. Captured as a
// struct so the resolution logic is unit-testable independent of the consts.
type buildConfig struct {
	version         string
	revision        string
	siblingRevision string
	name            string
	branch          string
	siblingBranch   string
	format          string // "direct" | "archive"
	naming          string // "short" | "triple"
}

func bakedConfig() buildConfig {
	return buildConfig{
		version:         loreVersion,
		revision:        loreRevision,
		siblingRevision: loreSiblingRevision,
		name:            loreName,
		branch:          loreBranch,
		siblingBranch:   loreSiblingBranch,
		format:          loreArtifactFormat,
		naming:          loreArtifactNaming,
	}
}

// fileVersion returns the version embedded in artifact filenames, matching the
// layout the release pipeline uses when it uploads them:
//
//	stable release            -> vX.Y.Z
//	prerelease (composite)    -> vX.Y.Z-<sibling>_<primary>
//	prerelease with name      -> vX.Y.Z-<sibling>_<primary>-NAME
//
// version is baked from the full, un-normalized build version (e.g.
// "0.8.2-nightly"); a leading 'v' is added if missing.
func (c buildConfig) fileVersion() string {
	version := c.version
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	revision := c.revision
	if c.siblingRevision != "" {
		revision = c.siblingRevision + "_" + c.revision
	}

	fileVersion := version
	if revision != "" {
		fileVersion += "-" + revision
	}
	if c.name != "" {
		fileVersion += "-" + c.name
	}
	return fileVersion
}

// dirVersion returns the directory component of the artifact URL. For feature
// prerelease builds (no name, and a non-'main' primary or sibling branch) the
// release pipeline disambiguates the directory with a branch suffix
// ("-lore_<sibling-branch>__urc_<primary-branch>"); otherwise it equals the
// file version. Stable releases live directly under the version tag, so this
// returns fileVersion for them too.
func (c buildConfig) dirVersion() string {
	fileVersion := c.fileVersion()
	if c.name != "" {
		return fileVersion
	}
	primary := c.branch
	if primary == "" {
		primary = "main"
	}
	overlay := c.siblingBranch
	if overlay == "" {
		overlay = "main"
	}
	if primary != "main" || overlay != "main" {
		return fmt.Sprintf("%s-lore_%s__urc_%s",
			fileVersion, sanitizeBranch(overlay), sanitizeBranch(primary))
	}
	return fileVersion
}

// platformTokens returns, for the running platform, the library's local name,
// raw file extension, the "short" platform token, the rust target triple, and
// the archive extension used when the lib is bundled.
func platformTokens(goos, goarch string) (localName, rawExt, short, triple, archiveExt string, err error) {
	switch {
	case goos == "windows" && goarch == "amd64":
		// Windows short names carry no platform token (lore-<tag>.dll).
		return "lore.dll", "dll", "", "x86_64-pc-windows-msvc", "zip", nil
	case goos == "darwin" && goarch == "arm64":
		return "liblore.dylib", "dylib", "macos-arm64", "aarch64-apple-darwin", "tar.gz", nil
	case goos == "linux" && goarch == "arm64":
		return "liblore.so", "so", "linux-arm64", "aarch64-unknown-linux-gnu-neoverse-512tvb", "tar.gz", nil
	case goos == "linux" && goarch == "amd64":
		return "liblore.so", "so", "linux-x86_64", "x86_64-unknown-linux-gnu", "tar.gz", nil
	default:
		return "", "", "", "", "", fmt.Errorf("unsupported platform: %s/%s", goos, goarch)
	}
}

func (c buildConfig) resolvePlatform(goos, goarch string) (*platformInfo, error) {
	localName, rawExt, short, triple, archiveExt, err := platformTokens(goos, goarch)
	if err != nil {
		return nil, err
	}
	fileVersion := c.fileVersion()

	// Archive distribution (e.g. open-source releases): a triple-named bundle
	// (always 'liblore-' prefixed, even on Windows) whose payload is the
	// generically-named library.
	if c.format == "archive" {
		return &platformInfo{
			artifactName: fmt.Sprintf("liblore-%s-%s.%s", fileVersion, triple, archiveExt),
			localName:    localName,
			archiveKind:  archiveExt,
			member:       localName,
		}, nil
	}

	// Direct distribution: a single raw library file. The Windows DLL uses the
	// 'lore-' prefix; the Unix shared libs use 'liblore-'.
	token := short
	if c.naming == "triple" {
		token = triple
	}
	var artifactName string
	if goos == "windows" {
		if token == "" {
			artifactName = fmt.Sprintf("lore-%s.%s", fileVersion, rawExt)
		} else {
			artifactName = fmt.Sprintf("lore-%s-%s.%s", fileVersion, token, rawExt)
		}
	} else {
		artifactName = fmt.Sprintf("liblore-%s-%s.%s", fileVersion, token, rawExt)
	}
	return &platformInfo{artifactName: artifactName, localName: localName}, nil
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

// downloadToTemp downloads url to a temp file and returns its path. The caller
// is responsible for removing it.
func downloadToTemp(url, pattern string) (string, error) {
	tmp, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmp.Close()
	if err := download(url, tmp.Name()); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

// extractMember writes the named member (matched by base name) from the archive
// to destPath. Supports gzipped tar and zip.
func extractMember(archivePath, kind, member, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	switch kind {
	case "tar.gz":
		return extractTarGzMember(archivePath, member, destPath)
	case "zip":
		return extractZipMember(archivePath, member, destPath)
	default:
		return fmt.Errorf("unsupported archive kind: %s", kind)
	}
}

func extractTarGzMember(archivePath, member, destPath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to open gzip: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar: %w", err)
		}
		if filepath.Base(hdr.Name) != member {
			continue
		}
		out, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination: %w", err)
		}
		defer out.Close()
		written, err := io.Copy(out, tr) //nolint:gosec // trusted release archive
		if err != nil {
			return fmt.Errorf("failed to extract %s: %w", member, err)
		}
		fmt.Printf("  Extracted %s (%d bytes)\n", member, written)
		return nil
	}
	return fmt.Errorf("member %s not found in archive", member)
}

func extractZipMember(archivePath, member, destPath string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer zr.Close()
	for _, zf := range zr.File {
		if filepath.Base(zf.Name) != member {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", member, err)
		}
		defer rc.Close()
		out, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination: %w", err)
		}
		defer out.Close()
		written, err := io.Copy(out, rc) //nolint:gosec // trusted release archive
		if err != nil {
			return fmt.Errorf("failed to extract %s: %w", member, err)
		}
		fmt.Printf("  Extracted %s (%d bytes)\n", member, written)
		return nil
	}
	return fmt.Errorf("member %s not found in archive", member)
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

// noExt strips the final extension from a filename (e.g. liblore.dylib ->
// liblore, liblore-v0.8.1-linux-x86_64.so -> liblore-v0.8.1-linux-x86_64).
func noExt(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

// downloadOptional fetches url to destPath best-effort. License files may not
// exist for every build, so a missing file (or any download error) is logged
// and ignored rather than failing the run.
func downloadOptional(url, destPath string) {
	if err := download(url, destPath); err != nil {
		fmt.Printf("  (optional) %s not fetched: %v\n", filepath.Base(destPath), err)
	}
}

// extractMemberOptional extracts an archive member best-effort, ignoring a
// missing member or extraction error (used for bundled license files).
func extractMemberOptional(archivePath, kind, member, destPath string) {
	if err := extractMember(archivePath, kind, member, destPath); err != nil {
		fmt.Printf("  (optional) %s not extracted: %v\n", member, err)
	}
}

// fetchDirectLicenses downloads the optional license files that sit next to a
// directly-distributed library in the same remote directory:
//   - Lore_Licenses.txt              (shared; name kept as-is)
//   - <remote-lib-basename>.THIRD-PARTY-NOTICES.txt, placed next to the library
//     as <local-lib-basename>.THIRD-PARTY-NOTICES.txt (e.g. liblore.THIRD-PARTY-NOTICES.txt)
//
// Both are optional — a build typically ships one or the other.
func fetchDirectLicenses(dirURL string, platform *platformInfo, outputDir string) {
	downloadOptional(dirURL+"/Lore_Licenses.txt", filepath.Join(outputDir, "Lore_Licenses.txt"))

	remoteNotices := noExt(platform.artifactName) + ".THIRD-PARTY-NOTICES.txt"
	localNotices := noExt(platform.localName) + ".THIRD-PARTY-NOTICES.txt"
	downloadOptional(dirURL+"/"+remoteNotices, filepath.Join(outputDir, localNotices))
}

func main() {
	outputDir := flag.String("o", ".", "output directory for the native library")
	targetOS := flag.String("os", runtime.GOOS, "target OS (e.g. linux, darwin, windows)")
	targetArch := flag.String("arch", runtime.GOARCH, "target architecture (e.g. amd64, arm64)")
	flag.Parse()

	cfg := bakedConfig()
	fileVersion := cfg.fileVersion()
	dirVersion := cfg.dirVersion()

	platform, err := cfg.resolvePlatform(*targetOS, *targetArch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	destPath := filepath.Join(*outputDir, platform.localName)

	// Skip if the destination already has the right version
	if _, err := os.Stat(destPath); err == nil {
		if checkVersion(*outputDir, fileVersion) {
			fmt.Printf("Lore library %s already present at %s\n", fileVersion, destPath)
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
		dirURL := fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), dirVersion)
		url := fmt.Sprintf("%s/%s", dirURL, platform.artifactName)

		fmt.Printf("Downloading Lore %s for %s/%s...\n", fileVersion, *targetOS, *targetArch)
		fmt.Printf("  URL: %s\n", url)
		fmt.Printf("  Destination: %s\n", destPath)

		if platform.archiveKind == "" {
			if err := download(url, destPath); err != nil {
				fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
				os.Exit(1)
			}
			// License files ship as separate files alongside the library.
			fetchDirectLicenses(dirURL, platform, *outputDir)
		} else {
			tmp, err := downloadToTemp(url, "lore-lib-*."+platform.archiveKind)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
				os.Exit(1)
			}
			defer os.Remove(tmp)
			if err := extractMember(tmp, platform.archiveKind, platform.member, destPath); err != nil {
				fmt.Fprintf(os.Stderr, "Extract failed: %v\n", err)
				os.Exit(1)
			}
			// License files are bundled inside the release archive.
			extractMemberOptional(tmp, platform.archiveKind, "LICENSE.txt", filepath.Join(*outputDir, "LICENSE.txt"))
			extractMemberOptional(tmp, platform.archiveKind, "THIRD-PARTY-NOTICES.txt", filepath.Join(*outputDir, "THIRD-PARTY-NOTICES.txt"))
		}
	}

	if err := writeVersion(*outputDir, fileVersion); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to write version marker: %v\n", err)
	}

	fmt.Printf("Successfully obtained Lore library %s\n", fileVersion)
}
