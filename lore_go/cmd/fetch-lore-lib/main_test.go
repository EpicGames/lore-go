// Copyright Epic Games, Inc. All Rights Reserved.

package main

import "testing"

// Covers the four distribution variants:
//   open-source release -> bundled triple-named archives under v0.8.1
//   open-source prerelease -> raw triple-named files under v0.8.2-nightly-129_3612
//   private release     -> raw short-named files under v0.8.1
//   private prerelease  -> raw short-named files under v0.8.2-nightly-129_3612
func TestResolveArtifact(t *testing.T) {
	// Composite prerelease: sibling revision 129, primary revision 3612 -> 129_3612.
	nightly := buildConfig{version: "0.8.2-nightly", revision: "3612", siblingRevision: "129"}
	release := buildConfig{version: "v0.8.1"}

	ossRelease := release
	ossRelease.format, ossRelease.naming = "archive", "triple"
	privRelease := release
	privRelease.format, privRelease.naming = "direct", "short"

	ossNightly := nightly
	ossNightly.format, ossNightly.naming = "direct", "triple"
	privNightly := nightly
	privNightly.format, privNightly.naming = "direct", "short"

	cases := []struct {
		name         string
		cfg          buildConfig
		goos, goarch string
		wantDir      string
		wantArtifact string
		wantArchive  string
		wantMember   string
	}{
		// --- OSS release: bundled triple archives on GitHub ---
		{"oss-release/darwin", ossRelease, "darwin", "arm64", "v0.8.1", "liblore-v0.8.1-aarch64-apple-darwin.tar.gz", "tar.gz", "liblore.dylib"},
		{"oss-release/linux-amd64", ossRelease, "linux", "amd64", "v0.8.1", "liblore-v0.8.1-x86_64-unknown-linux-gnu.tar.gz", "tar.gz", "liblore.so"},
		{"oss-release/linux-arm64", ossRelease, "linux", "arm64", "v0.8.1", "liblore-v0.8.1-aarch64-unknown-linux-gnu-neoverse-512tvb.tar.gz", "tar.gz", "liblore.so"},
		{"oss-release/windows", ossRelease, "windows", "amd64", "v0.8.1", "liblore-v0.8.1-x86_64-pc-windows-msvc.zip", "zip", "lore.dll"},

		// --- private release: raw short-named files ---
		{"priv-release/darwin", privRelease, "darwin", "arm64", "v0.8.1", "liblore-v0.8.1-macos-arm64.dylib", "", ""},
		{"priv-release/linux-amd64", privRelease, "linux", "amd64", "v0.8.1", "liblore-v0.8.1-linux-x86_64.so", "", ""},
		{"priv-release/windows", privRelease, "windows", "amd64", "v0.8.1", "lore-v0.8.1.dll", "", ""},

		// --- open-source prerelease: raw triple-named files ---
		{"oss-nightly/darwin", ossNightly, "darwin", "arm64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-aarch64-apple-darwin.dylib", "", ""},
		{"oss-nightly/linux-amd64", ossNightly, "linux", "amd64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-x86_64-unknown-linux-gnu.so", "", ""},
		{"oss-nightly/linux-arm64", ossNightly, "linux", "arm64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-aarch64-unknown-linux-gnu-neoverse-512tvb.so", "", ""},
		{"oss-nightly/windows", ossNightly, "windows", "amd64", "v0.8.2-nightly-129_3612", "lore-v0.8.2-nightly-129_3612-x86_64-pc-windows-msvc.dll", "", ""},

		// --- private prerelease: raw short-named files ---
		{"priv-nightly/darwin", privNightly, "darwin", "arm64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-macos-arm64.dylib", "", ""},
		{"priv-nightly/linux-arm64", privNightly, "linux", "arm64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-linux-arm64.so", "", ""},
		{"priv-nightly/windows", privNightly, "windows", "amd64", "v0.8.2-nightly-129_3612", "lore-v0.8.2-nightly-129_3612.dll", "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cfg.dirVersion(); got != tc.wantDir {
				t.Errorf("dirVersion = %q, want %q", got, tc.wantDir)
			}
			p, err := tc.cfg.resolvePlatform(tc.goos, tc.goarch)
			if err != nil {
				t.Fatalf("resolvePlatform: %v", err)
			}
			if p.artifactName != tc.wantArtifact {
				t.Errorf("artifactName = %q, want %q", p.artifactName, tc.wantArtifact)
			}
			if p.archiveKind != tc.wantArchive {
				t.Errorf("archiveKind = %q, want %q", p.archiveKind, tc.wantArchive)
			}
			if p.member != tc.wantMember {
				t.Errorf("member = %q, want %q", p.member, tc.wantMember)
			}
		})
	}
}

// A feature prerelease build (non-main branch) gets the branch dir suffix on
// the URL path, while the filename keeps the bare file version.
func TestDirVersionFeatureBranch(t *testing.T) {
	cfg := buildConfig{
		version: "0.8.2-nightly", revision: "3612", siblingRevision: "129",
		branch: "feature/x", siblingBranch: "main",
	}
	const wantFile = "v0.8.2-nightly-129_3612"
	const wantDir = "v0.8.2-nightly-129_3612-lore_main__urc_feature-x"
	if got := cfg.fileVersion(); got != wantFile {
		t.Errorf("fileVersion = %q, want %q", got, wantFile)
	}
	if got := cfg.dirVersion(); got != wantDir {
		t.Errorf("dirVersion = %q, want %q", got, wantDir)
	}
}
