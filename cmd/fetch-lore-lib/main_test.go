// Copyright Epic Games, Inc. All Rights Reserved.

package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

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

	// A named branch build of the archive layout keeps the plain version.
	archiveNightly := nightly
	archiveNightly.format, archiveNightly.naming = "archive", "triple"
	archiveNightly.name, archiveNightly.branch = "probe", "feature/x"

	ossReleaseGeneric := ossRelease
	ossReleaseGeneric.linuxArm64Build = linuxArm64Generic
	archiveNightlyGeneric := archiveNightly
	archiveNightlyGeneric.linuxArm64Build = linuxArm64Generic

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

		// --- archive prerelease (named, on a branch): triple archive under the plain version ---
		{"archive-nightly/linux-amd64", archiveNightly, "linux", "amd64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-x86_64-unknown-linux-gnu.tar.gz", "tar.gz", "liblore.so"},

		// --- generic linux/arm64: the untuned triple; other platforms unaffected ---
		{"oss-release-generic/linux-arm64", ossReleaseGeneric, "linux", "arm64", "v0.8.1", "liblore-v0.8.1-aarch64-unknown-linux-gnu.tar.gz", "tar.gz", "liblore.so"},
		{"oss-release-generic/darwin", ossReleaseGeneric, "darwin", "arm64", "v0.8.1", "liblore-v0.8.1-aarch64-apple-darwin.tar.gz", "tar.gz", "liblore.dylib"},
		{"archive-nightly-generic/linux-arm64", archiveNightlyGeneric, "linux", "arm64", "v0.8.2-nightly-129_3612", "liblore-v0.8.2-nightly-129_3612-aarch64-unknown-linux-gnu.tar.gz", "tar.gz", "liblore.so"},
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

// The third-party-notices file is named after the remote library (without its
// extension) and placed next to the local library as
// <local-lib-basename>.THIRD-PARTY-NOTICES.txt — mirroring the build-time logic.
func TestNoticesNames(t *testing.T) {
	cases := []struct {
		name              string
		artifact          string
		local             string
		wantRemoteNotices string
		wantLocalNotices  string
	}{
		{"linux-x86_64", "liblore-v0.8.1-linux-x86_64.so", "liblore.so",
			"liblore-v0.8.1-linux-x86_64.THIRD-PARTY-NOTICES.txt", "liblore.THIRD-PARTY-NOTICES.txt"},
		{"windows", "lore-v0.8.1.dll", "lore.dll",
			"lore-v0.8.1.THIRD-PARTY-NOTICES.txt", "lore.THIRD-PARTY-NOTICES.txt"},
		{"oss-nightly-darwin", "liblore-v0.8.2-nightly-129_3612-aarch64-apple-darwin.dylib", "liblore.dylib",
			"liblore-v0.8.2-nightly-129_3612-aarch64-apple-darwin.THIRD-PARTY-NOTICES.txt", "liblore.THIRD-PARTY-NOTICES.txt"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := noExt(tc.artifact) + ".THIRD-PARTY-NOTICES.txt"; got != tc.wantRemoteNotices {
				t.Errorf("remote notices = %q, want %q", got, tc.wantRemoteNotices)
			}
			if got := noExt(tc.local) + ".THIRD-PARTY-NOTICES.txt"; got != tc.wantLocalNotices {
				t.Errorf("local notices = %q, want %q", got, tc.wantLocalNotices)
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

// One pass extracts every wanted member present and reports which were found.
func TestExtractMembers(t *testing.T) {
	dir := t.TempDir()
	members := map[string]string{"./liblore.so": "lib", "./Lore_Licenses.txt": "licenses", "./lore.h": "header"}

	tgz := filepath.Join(dir, "a.tar.gz")
	f, _ := os.Create(tgz)
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	for name, body := range members {
		_ = tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body))})
		_, _ = tw.Write([]byte(body))
	}
	tw.Close()
	gw.Close()
	f.Close()

	zp := filepath.Join(dir, "a.zip")
	f, _ = os.Create(zp)
	zw := zip.NewWriter(f)
	for name, body := range members {
		w, _ := zw.Create(name[2:])
		_, _ = w.Write([]byte(body))
	}
	zw.Close()
	f.Close()

	for _, tc := range []struct{ kind, path string }{{"tar.gz", tgz}, {"zip", zp}} {
		t.Run(tc.kind, func(t *testing.T) {
			out := t.TempDir()
			dests := map[string]string{
				"liblore.so":        filepath.Join(out, "liblore.so"),
				"Lore_Licenses.txt": filepath.Join(out, "Lore_Licenses.txt"),
				"LICENSE.txt":       filepath.Join(out, "LICENSE.txt"),
			}
			found, err := extractMembers(tc.path, tc.kind, dests)
			if err != nil {
				t.Fatal(err)
			}
			want := map[string]string{"liblore.so": "lib", "Lore_Licenses.txt": "licenses"}
			if len(found) != len(want) {
				t.Errorf("found = %v, want %v", found, want)
			}
			for name, body := range want {
				got, err := os.ReadFile(dests[name])
				if !found[name] || err != nil || string(got) != body {
					t.Errorf("%s: found=%v content=%q err=%v, want %q", name, found[name], got, err, body)
				}
			}
		})
	}
}

const (
	graviton3CPUInfo = "processor\t: 0\nBogoMIPS\t: 2100.00\nFeatures\t: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm jscvt fcma lrcpc dcpop sha3 sm3 sm4 asimddp sha512 sve asimdfhm dit uscat ilrcpc flagm ssbs paca pacg dcpodp svei8mm svebf16 i8mm bf16 dgh rng\nCPU implementer\t: 0x41\n"
	graviton2CPUInfo = "processor\t: 0\nFeatures\t: fp asimd evtstrm aes pmull sha1 sha2 crc32 atomics fphp asimdhp cpuid asimdrdm lrcpc dcpop asimddp ssbs\nCPU implementer\t: 0x41\n"
)

func TestResolveLinuxArm64Build(t *testing.T) {
	cpuinfo := func(s string) func() (string, error) {
		return func() (string, error) { return s, nil }
	}
	noCPUInfo := func() (string, error) { return "", errors.New("no /proc") }

	cases := []struct {
		name                 string
		requested, format    string
		targetOS, targetArch string
		hostOS, hostArch     string
		readCPUInfo          func() (string, error)
		want                 string
		wantErr              bool
	}{
		{"not-linux-arm64", "generic", "archive", "linux", "amd64", "linux", "arm64", cpuinfo(graviton3CPUInfo), "", false},
		{"auto/graviton3-host", "auto", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton3CPUInfo), linuxArm64Graviton, false},
		{"empty-means-auto", "", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton3CPUInfo), linuxArm64Graviton, false},
		{"auto/graviton2-host", "auto", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton2CPUInfo), linuxArm64Generic, false},
		{"auto/unreadable-cpuinfo", "auto", "archive", "linux", "arm64", "linux", "arm64", noCPUInfo, linuxArm64Generic, false},
		{"auto/cross-fetch", "auto", "archive", "linux", "arm64", "darwin", "arm64", cpuinfo(graviton3CPUInfo), linuxArm64Generic, false},
		{"auto/direct-layout", "auto", "direct", "linux", "arm64", "linux", "arm64", cpuinfo(graviton2CPUInfo), linuxArm64Graviton, false},
		{"graviton/overrides-cpu", "graviton", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton2CPUInfo), linuxArm64Graviton, false},
		{"generic/overrides-cpu", "generic", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton3CPUInfo), linuxArm64Generic, false},
		{"generic/direct-layout", "generic", "direct", "linux", "arm64", "linux", "arm64", cpuinfo(graviton3CPUInfo), "", true},
		{"unknown-value", "neoverse", "archive", "linux", "arm64", "linux", "arm64", cpuinfo(graviton3CPUInfo), "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _, err := resolveLinuxArm64Build(tc.requested, tc.format, tc.targetOS, tc.targetArch, tc.hostOS, tc.hostArch, tc.readCPUInfo)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("build = %q, want %q", got, tc.want)
			}
		})
	}
}
