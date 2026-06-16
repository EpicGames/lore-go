// Copyright Epic Games, Inc. All Rights Reserved.

// Performance test for native vs. fluent API event handling in the Go SDK.
//
// Mirrors the JS SDK perf test at urc-js-sdk/src/perf/repository-dump.perf.ts
// so we can compare per-event JS-FFI vs Go-FFI overhead head to head.
//
// Run with:
//
//	go run -C lore_go ./cmd/perf-repository-dump
//
// To stabilize numbers on macOS / Linux (analog of the JS test:perf:*lowpower
// scripts):
//
//	taskpolicy -c utility go run -C lore_go ./cmd/perf-repository-dump    # macOS
//	nice -n -19           go run -C lore_go ./cmd/perf-repository-dump    # Linux
//
// Measures the cost of consuming LORE_EVENT_REPOSITORY_STATE_DUMP_NODE events
// across four SDK access patterns:
//  1. raw native callback (urc-go-sdk/native)
//  2. fluent Callback + Wait
//  3. fluent AsyncIter (channel)
//  4. fluent Collect (slice)
//
// Two variants per mode:
//
//	A. accumulate event.Size only
//	B. accumulate len(Name), len(TypeData), and every numeric field, forcing
//	   the LoreString -> Go string decode on the FFI paths (native + Callback)
//	   to match what the JS variant B does. Collect/AsyncIter pay the same
//	   decode cost up-front during Clone, so for those paths variant B is
//	   just plain `len()` reads.
//
// Each (mode, variant) pair runs in its OWN child process so peak RSS can be
// attributed cleanly per access pattern. Within one child: warmup + N_RUNS
// measured rounds. The parent orchestrates 4 modes × 2 variants = 8 children
// sequentially. The shared setup (create repo + stage 100k files + commit) is
// done once in the parent; children re-open the existing repo via the
// -child-repo flag.
//
// Trade-off: we lose per-round cross-mode interleaving (a system blip during
// one child only affects that mode's numbers). In exchange the per-mode peak
// RSS is no longer polluted by previous modes' allocations.
//
// Each child also tightens Go's GC pacing (debug.SetGCPercent(50)) before
// running. The default GOGC=100 lets the heap grow to ~2× live size before
// collecting, which makes peakRSS non-deterministic enough that the same
// mode can peak at 140 MB on one run and 320 MB on the next. GOGC=50 collects
// at 1.5× live size — smaller, more stable peaks at a small CPU cost. This
// keeps peakRSS reflective of working-set size (not cumulative allocation)
// while damping the run-to-run noise.
//
// To eliminate disk-cache variance, point the repo at a ramdisk by exporting
// LORE_PERF_REPO_PARENT before running. Defaults to os.TempDir() otherwise.
//
//	# Linux — /dev/shm is already tmpfs, no setup needed:
//	LORE_PERF_REPO_PARENT=/dev/shm go run -C lore_go ./cmd/perf-repository-dump
//
//	# macOS — create a 4 GB ramdisk once, reuse across runs, then eject:
//	diskutil erasevolume APFS perfdisk $(hdiutil attach -nomount ram://8388608)
//	LORE_PERF_REPO_PARENT=/Volumes/perfdisk go run -C lore_go ./cmd/perf-repository-dump
//	diskutil eject /Volumes/perfdisk
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/EpicGames/lore-go"
	"github.com/EpicGames/lore-go/internal/testutil"
	"github.com/EpicGames/lore-go/native"
	"github.com/EpicGames/lore-go/types"
)

const (
	fileCount        = 100_000
	filesPerLeafDir  = 100
	topDirs          = 10
	subDirs          = 100
	nRuns            = 10
	cooldownDuration = 500 * time.Millisecond
	// childGcPercent is passed to debug.SetGCPercent at the start of each
	// child. Default Go is 100 (heap can grow to 2× live before collecting),
	// which produces large run-to-run variance in peakRSS. Lower = more
	// frequent GC, smaller and more stable peaks, slightly higher CPU. Set
	// to -1 to disable GC entirely (peakRSS = cumulative allocation, even
	// more stable but harder to interpret).
	childGcPercent = 300
)

const nodeTag = types.LoreEventTag_REPOSITORY_STATE_DUMP_NODE

type mode string

const (
	modeNative     mode = "native"
	modeFluentCB   mode = "fluent-callback"
	modeFluentIter mode = "fluent-AsyncIter"
	modeFluentColl mode = "fluent-Collect"
)

var modes = []mode{modeNative, modeFluentCB, modeFluentIter, modeFluentColl}

type variant int

const (
	variantA variant = iota
	variantB
)

var variants = []variant{variantA, variantB}

func (v variant) name() string {
	if v == variantA {
		return "A"
	}
	return "B"
}

func (v variant) label() string {
	if v == variantA {
		return "(event.Size only)"
	}
	return "(len(Name) + len(TypeData) + numeric fields)"
}

func variantFromString(s string) (variant, error) {
	switch s {
	case "A":
		return variantA, nil
	case "B":
		return variantB, nil
	}
	return 0, fmt.Errorf("unknown variant %q", s)
}

// Exported field names so encoding/json can marshal them across the parent/
// child boundary.
type pass struct {
	Events           uint64  `json:"events"`
	AccumulatedSize  uint64  `json:"accumulatedSize"`
	Ms               float64 `json:"ms"`
	RssBytes         uint64  `json:"rssBytes"`
	NameLenTotal     uint64  `json:"nameLenTotal"`
	TypeDataLenTotal uint64  `json:"typeDataLenTotal"`
	NumericTotal     uint64  `json:"numericTotal"`
}

type childResult struct {
	Mode         mode   `json:"mode"`
	Variant      string `json:"variant"`
	Passes       []pass `json:"passes"`
	PeakRssBytes uint64 `json:"peakRssBytes"`
}

type variantResult struct {
	v       variant
	perMode map[mode]childResult
}

// CLI flags. Empty -child-mode means parent role.
var (
	flagChildMode    = flag.String("child-mode", "", "internal: run as child for this mode")
	flagChildVariant = flag.String("child-variant", "", "internal: child variant (A|B)")
	flagChildRepo    = flag.String("child-repo", "", "internal: child repository path")
)

func main() {
	flag.Parse()

	if err := testutil.SetLibraryPath(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to locate native library: %v\n", err)
		os.Exit(1)
	}

	if *flagChildMode != "" {
		if err := runChild(); err != nil {
			fmt.Fprintf(os.Stderr, "child failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := runParent(); err != nil {
		fmt.Fprintf(os.Stderr, "perf failed: %v\n", err)
		os.Exit(1)
	}
}

// --- child path -----------------------------------------------------------

func runChild() error {
	// Tighter-than-default GC pacing for stable per-mode RSS comparisons.
	// See childGcPercent for tuning notes.
	debug.SetGCPercent(childGcPercent)

	m := mode(*flagChildMode)
	if !isKnownMode(m) {
		return fmt.Errorf("unknown -child-mode %q", *flagChildMode)
	}

	// Pin the main goroutine to one OS thread for modes that run the dump
	// synchronously on this goroutine (native, fluent-callback,
	// fluent-Collect). Why: the lib's malloc/free calls inside lore_* go
	// through macOS libmalloc, which keeps a per-thread arena per M that
	// holds freed pages indefinitely. Without pinning, the Go scheduler
	// is free to bounce the goroutine between M's between iterations, so
	// some runs touch 1 arena and others touch 2 — producing the run-to-run
	// peakRSS spread we saw (e.g. native: 91-206 MB across runs). Pinning
	// makes "1 arena" deterministic, dropping each mode's peakRSS to the
	// low end of its previous range and producing stable cross-mode
	// comparisons.
	//
	// We deliberately SKIP fluent-AsyncIter: it streams events via a
	// channel between two goroutines (one in the lib callback feeding the
	// channel, one consuming via `for range`). Locking the consumer pins
	// only half the work; the unlocked feeder still bounces across M's,
	// growing additional arenas. Worse, the locked consumer can't be
	// rescheduled efficiently, so synchronization stalls and the run gets
	// 3× slower with HIGHER peakRSS than the unpinned baseline. AsyncIter
	// is inherently multi-goroutine, so its arena count is a property of
	// the access pattern, not noise we can pin away.
	if m != modeFluentIter {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	v, err := variantFromString(*flagChildVariant)
	if err != nil {
		return err
	}
	if *flagChildRepo == "" {
		return fmt.Errorf("-child-repo required")
	}

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: *flagChildRepo,
		Offline:        true,
		CorrelationId:  fmt.Sprintf("perf-child-%s-%s", m, v.name()),
	})
	defer cleanupGlobals()

	tag := fmt.Sprintf("[mode=%-22s variant=%s]", string(m), v.name())

	// Warmup
	time.Sleep(cooldownDuration)
	warm, err := runMode(m, v, &globals)
	if err != nil {
		return fmt.Errorf("warmup: %w", err)
	}
	logChild("%s warmup    time=%s  events=%d  rss=%s",
		tag, fmtMs(warm.Ms), warm.Events, fmtMb(float64(warm.RssBytes)))

	passes := make([]pass, 0, nRuns)
	for round := 1; round <= nRuns; round++ {
		time.Sleep(cooldownDuration)
		p, err := runMode(m, v, &globals)
		if err != nil {
			return fmt.Errorf("round %d: %w", round, err)
		}
		passes = append(passes, p)
		logChild("%s round=%2d time=%s  events=%d  rss=%s",
			tag, round, fmtMs(p.Ms), p.Events, fmtMb(float64(p.RssBytes)))
	}

	result := childResult{
		Mode:         m,
		Variant:      v.name(),
		Passes:       passes,
		PeakRssBytes: processPeakRssBytes(),
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	// Single JSON line on stdout; progress already went to stderr.
	if _, err := os.Stdout.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("write stdout: %w", err)
	}
	return nil
}

func isKnownMode(m mode) bool {
	for _, known := range modes {
		if m == known {
			return true
		}
	}
	return false
}

// --- parent path ----------------------------------------------------------

func runParent() error {
	globals, repositoryPath, cleanupGlobals, err := setupRepo()
	if err != nil {
		return err
	}
	defer func() {
		flushArgs, cleanupFlush := types.NewLoreRepositoryFlushArgs(types.LoreRepositoryFlushArgs{})
		defer cleanupFlush()
		_, _ = lore.RepositoryFlush(globals, &flushArgs).Wait()
		cleanupGlobals()
		if err := os.RemoveAll(repositoryPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to remove temp dir %s: %v\n", repositoryPath, err)
		}
	}()

	results := make([]variantResult, 0, len(variants))
	for _, v := range variants {
		logParent("\n--- Variant %s %s: running each mode in its own child process ---",
			v.name(), v.label())
		perMode := make(map[mode]childResult, len(modes))
		for _, m := range modes {
			// Cooldown between children so the previous one's tail GC etc.
			// has time to settle before the next child starts measuring.
			time.Sleep(cooldownDuration)
			r, err := spawnChild(m, v, repositoryPath)
			if err != nil {
				return fmt.Errorf("spawn child %s/%s: %w", m, v.name(), err)
			}
			perMode[m] = r
		}
		results = append(results, variantResult{v: v, perMode: perMode})
	}

	for _, r := range results {
		checkConsistency(r)
	}
	for _, r := range results {
		printSummary(r)
	}
	return nil
}

func spawnChild(m mode, v variant, repoPath string) (childResult, error) {
	cmd := exec.Command(os.Args[0],
		"-child-mode", string(m),
		"-child-variant", v.name(),
		"-child-repo", repoPath,
	)
	// stderr passes through so per-round progress is visible in real time.
	cmd.Stderr = os.Stderr
	var stdoutBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf

	if err := cmd.Run(); err != nil {
		return childResult{}, fmt.Errorf("run: %w", err)
	}

	out := strings.TrimSpace(stdoutBuf.String())
	var result childResult
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		preview := out
		if len(preview) > 500 {
			preview = preview[:500]
		}
		return childResult{}, fmt.Errorf("decode child JSON: %w (raw=%q)", err, preview)
	}
	return result, nil
}

// --- shared: setup / runners / accumulators -------------------------------

func setupRepo() (*types.LoreGlobalArgsFFI, string, func(), error) {
	parentDir := os.Getenv("LORE_PERF_REPO_PARENT")
	parentLabel := "(parent from LORE_PERF_REPO_PARENT)"
	if parentDir == "" {
		parentDir = os.TempDir()
		parentLabel = "(parent from os.TempDir())"
	}

	repositoryPath, err := os.MkdirTemp(parentDir, "lore-go-sdk-perf-")
	if err != nil {
		return nil, "", nil, fmt.Errorf("MkdirTemp: %w", err)
	}

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: repositoryPath,
		Offline:        true,
		CorrelationId:  "perf-repository-dump",
	})

	logParent("setup: repo at %s %s", repositoryPath, parentLabel)

	t := time.Now()
	repoUrl := fmt.Sprintf("perf-%d", time.Now().UnixNano())
	createArgs, cleanupCreate := types.NewLoreRepositoryCreateArgs(types.LoreRepositoryCreateArgs{
		RepositoryUrl: repoUrl,
	})
	if _, err := lore.RepositoryCreate(&globals, &createArgs).Wait(); err != nil {
		cleanupCreate()
		cleanupGlobals()
		os.RemoveAll(repositoryPath)
		return nil, "", nil, fmt.Errorf("RepositoryCreate: %w", err)
	}
	cleanupCreate()
	logParent("setup: RepositoryCreate done (%s)", fmtSince(t))

	t = time.Now()
	if err := createFiles(repositoryPath); err != nil {
		cleanupGlobals()
		os.RemoveAll(repositoryPath)
		return nil, "", nil, fmt.Errorf("createFiles: %w", err)
	}
	logParent("setup: created %d files in %d leaf dirs (%s)", fileCount, topDirs*subDirs, fmtSince(t))

	t = time.Now()
	stageArgs, cleanupStage := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
		Paths: []string{repositoryPath},
	})
	if _, err := lore.FileStage(&globals, &stageArgs).Wait(); err != nil {
		cleanupStage()
		cleanupGlobals()
		os.RemoveAll(repositoryPath)
		return nil, "", nil, fmt.Errorf("FileStage: %w", err)
	}
	cleanupStage()
	logParent("setup: FileStage done (%s)", fmtSince(t))

	t = time.Now()
	commitArgs, cleanupCommit := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
		Message: "perf setup",
	})
	if _, err := lore.RevisionCommit(&globals, &commitArgs).Wait(); err != nil {
		cleanupCommit()
		cleanupGlobals()
		os.RemoveAll(repositoryPath)
		return nil, "", nil, fmt.Errorf("RevisionCommit: %w", err)
	}
	cleanupCommit()
	logParent("setup: RevisionCommit done (%s)", fmtSince(t))

	t = time.Now()
	flushArgs, cleanupFlush := types.NewLoreRepositoryFlushArgs(types.LoreRepositoryFlushArgs{})
	if _, err := lore.RepositoryFlush(&globals, &flushArgs).Wait(); err != nil {
		cleanupFlush()
		cleanupGlobals()
		os.RemoveAll(repositoryPath)
		return nil, "", nil, fmt.Errorf("RepositoryFlush: %w", err)
	}
	cleanupFlush()
	logParent("setup: RepositoryFlush done (%s)", fmtSince(t))

	return &globals, repositoryPath, cleanupGlobals, nil
}

func createFiles(repoPath string) error {
	for top := 0; top < topDirs; top++ {
		for sub := 0; sub < subDirs; sub++ {
			if err := os.MkdirAll(filepath.Join(repoPath, pad2(top), pad2(sub)), 0o755); err != nil {
				return err
			}
		}
	}
	for n := 0; n < fileCount; n++ {
		top := n / 10_000
		sub := (n / filesPerLeafDir) % subDirs
		name := pad6(n)
		path := filepath.Join(repoPath, pad2(top), pad2(sub), name+".txt")
		if err := os.WriteFile(path, []byte(name), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func pad2(n int) string { return fmt.Sprintf("%02d", n) }
func pad6(n int) string { return fmt.Sprintf("%06d", n) }

func runMode(m mode, v variant, globals *types.LoreGlobalArgsFFI) (pass, error) {
	switch m {
	case modeNative:
		return runNative(globals, v)
	case modeFluentCB:
		return runFluentCallback(globals, v)
	case modeFluentIter:
		return runFluentAsyncIter(globals, v)
	case modeFluentColl:
		return runFluentCollect(globals, v)
	default:
		return pass{}, fmt.Errorf("unknown mode %s", m)
	}
}

// currentRssBytes returns the process's current resident set size in bytes.
// macOS has no Go stdlib API for this, so we shell out to `ps -o rss=`, which
// returns kilobytes. ~5ms overhead per call; called outside the timed window.
func currentRssBytes() uint64 {
	out, err := exec.Command("ps", "-o", "rss=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		return 0
	}
	kb, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return 0
	}
	return kb * 1024
}

func runNative(globals *types.LoreGlobalArgsFFI, v variant) (pass, error) {
	args, cleanup := types.NewLoreRepositoryDumpArgs(types.LoreRepositoryDumpArgs{})
	defer cleanup()

	var p pass
	cb := types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag != nodeTag {
				return
			}
			data, ok := event.GetData().(*types.LoreRepositoryStateDumpNodeEventDataFFI)
			if !ok {
				return
			}
			consumeFFI(&p, v, data)
		},
	}

	t0 := time.Now()
	rc, err := native.RepositoryDump(globals, &args, &cb)
	p.Ms = msSince(t0)
	p.RssBytes = currentRssBytes()
	if err != nil {
		return p, err
	}
	if rc != 0 {
		return p, fmt.Errorf("native LoreRepositoryDump returned %d", rc)
	}
	return p, nil
}

func runFluentCallback(globals *types.LoreGlobalArgsFFI, v variant) (pass, error) {
	args, cleanup := types.NewLoreRepositoryDumpArgs(types.LoreRepositoryDumpArgs{})
	defer cleanup()

	var p pass
	t0 := time.Now()
	_, err := lore.RepositoryDump(globals, &args).
		FilterByType(nodeTag).
		Callback(func(event *types.LoreEventFFI, _ uint64) {
			data, ok := event.GetData().(*types.LoreRepositoryStateDumpNodeEventDataFFI)
			if !ok {
				return
			}
			consumeFFI(&p, v, data)
		}).
		Wait()
	p.Ms = msSince(t0)
	p.RssBytes = currentRssBytes()
	return p, err
}

func runFluentAsyncIter(globals *types.LoreGlobalArgsFFI, v variant) (pass, error) {
	args, cleanup := types.NewLoreRepositoryDumpArgs(types.LoreRepositoryDumpArgs{})
	defer cleanup()

	var p pass
	t0 := time.Now()
	eventCh, errCh := lore.RepositoryDump(globals, &args).
		FilterByType(nodeTag).
		AsyncIter()
	for event := range eventCh {
		data, ok := event.Data.(types.LoreRepositoryStateDumpNodeEventData)
		if !ok {
			continue
		}
		consumeCloned(&p, v, &data)
	}
	err := <-errCh
	p.Ms = msSince(t0)
	p.RssBytes = currentRssBytes()
	return p, err
}

func runFluentCollect(globals *types.LoreGlobalArgsFFI, v variant) (pass, error) {
	args, cleanup := types.NewLoreRepositoryDumpArgs(types.LoreRepositoryDumpArgs{})
	defer cleanup()

	var p pass
	t0 := time.Now()
	events, err := lore.RepositoryDump(globals, &args).
		FilterByType(nodeTag).
		Collect()
	if err != nil {
		p.Ms = msSince(t0)
		p.RssBytes = currentRssBytes()
		return p, err
	}
	for i := range events {
		data, ok := events[i].Data.(types.LoreRepositoryStateDumpNodeEventData)
		if !ok {
			continue
		}
		consumeCloned(&p, v, &data)
	}
	p.Ms = msSince(t0)
	p.RssBytes = currentRssBytes()
	return p, nil
}

// consumeFFI handles events that are still backed by FFI memory. Variant B
// calls .String() on the LoreString fields to force the C-to-Go decode,
// matching what the JS variant B is paying via the lazy koffi struct decode.
func consumeFFI(p *pass, v variant, data *types.LoreRepositoryStateDumpNodeEventDataFFI) {
	p.Events++
	p.AccumulatedSize += data.Size
	if v == variantB {
		p.NameLenTotal += uint64(len(data.Name.String()))
		p.TypeDataLenTotal += uint64(len(data.TypeData.String()))
		p.NumericTotal += uint64(data.Id) + uint64(data.Parent) + uint64(data.Sibling) +
			uint64(data.Mode) + data.Size + uint64(data.Flags)
	}
}

// consumeCloned handles already-Go-native events from Collect / AsyncIter.
// Strings are already decoded by the SDK during Clone, so variant B just
// reads len() with no allocation here — same as the JS Collect/AsyncIter
// paths reading already-materialized .data.name.
func consumeCloned(p *pass, v variant, data *types.LoreRepositoryStateDumpNodeEventData) {
	p.Events++
	p.AccumulatedSize += data.Size
	if v == variantB {
		p.NameLenTotal += uint64(len(data.Name))
		p.TypeDataLenTotal += uint64(len(data.TypeData))
		p.NumericTotal += uint64(data.Id) + uint64(data.Parent) + uint64(data.Sibling) +
			uint64(data.Mode) + data.Size + uint64(data.Flags)
	}
}

// --- reporting -------------------------------------------------------------

func checkConsistency(r variantResult) {
	type sample struct {
		mode  mode
		round int
		p     pass
	}
	var samples []sample
	for _, m := range modes {
		child, ok := r.perMode[m]
		if !ok {
			continue
		}
		for i, p := range child.Passes {
			samples = append(samples, sample{mode: m, round: i + 1, p: p})
		}
	}
	if len(samples) == 0 {
		return
	}
	ref := samples[0].p
	for _, s := range samples {
		if s.p.Events != ref.Events {
			logParent("  WARN variant %s %s round%d: events=%d differs from reference %d",
				r.v.name(), s.mode, s.round, s.p.Events, ref.Events)
		}
		if s.p.AccumulatedSize != ref.AccumulatedSize {
			logParent("  WARN variant %s %s round%d: accumulatedSize=%d differs from reference %d",
				r.v.name(), s.mode, s.round, s.p.AccumulatedSize, ref.AccumulatedSize)
		}
		if r.v == variantB {
			if s.p.NameLenTotal != ref.NameLenTotal ||
				s.p.TypeDataLenTotal != ref.TypeDataLenTotal ||
				s.p.NumericTotal != ref.NumericTotal {
				logParent("  WARN variant B %s round%d: heavy-field accumulators differ from reference (nameLen=%d/%d typeDataLen=%d/%d numeric=%d/%d)",
					s.mode, s.round,
					s.p.NameLenTotal, ref.NameLenTotal,
					s.p.TypeDataLenTotal, ref.TypeDataLenTotal,
					s.p.NumericTotal, ref.NumericTotal)
			}
		}
	}
}

func printSummary(r variantResult) {
	logParent("\n=== Variant %s %s — summary over %d runs per mode (each mode in its own child process) ===",
		r.v.name(), r.v.label(), nRuns)

	type stats struct {
		m       mode
		min     float64
		mean    float64
		max     float64
		peakRss uint64
		eps     uint64
		events  uint64
	}
	var rows []stats
	for _, m := range modes {
		child, ok := r.perMode[m]
		if !ok || len(child.Passes) == 0 {
			continue
		}
		mn, mx, sum := math.Inf(1), math.Inf(-1), 0.0
		var totalEvents uint64
		var totalMs float64
		for _, p := range child.Passes {
			if p.Ms < mn {
				mn = p.Ms
			}
			if p.Ms > mx {
				mx = p.Ms
			}
			sum += p.Ms
			totalEvents += p.Events
			totalMs += p.Ms
		}
		mean := sum / float64(len(child.Passes))
		var eps uint64
		if totalMs > 0 {
			eps = uint64(float64(totalEvents) * 1000.0 / totalMs)
		}
		rows = append(rows, stats{
			m: m, min: mn, mean: mean, max: mx,
			peakRss: child.PeakRssBytes,
			eps:     eps,
			events:  child.Passes[0].Events,
		})
	}

	fastestMean := math.Inf(1)
	for _, s := range rows {
		if s.mean < fastestMean {
			fastestMean = s.mean
		}
	}

	for _, s := range rows {
		ratio := s.mean / fastestMean
		logParent("mode=%-24s events=%d  min=%s  mean=%s  max=%s  ev/s=%9s  peakRSS=%s  (mean %.2fx)",
			string(s.m), s.events, fmtMs(s.min), fmtMs(s.mean), fmtMs(s.max),
			withThousands(s.eps), fmtMb(float64(s.peakRss)), ratio)
	}
}

func fmtMb(bytes float64) string {
	return fmt.Sprintf("%6.1f MB", bytes/1024.0/1024.0)
}

// --- formatting helpers ----------------------------------------------------

// logParent writes to stdout — used by the parent process.
func logParent(format string, args ...any) {
	fmt.Fprintf(os.Stdout, format+"\n", args...)
}

// logChild writes to stderr — child processes write their per-round progress
// to stderr so the parent (which inherits stderr) can show it live, while
// stdout is reserved for the JSON result.
func logChild(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func fmtSince(t time.Time) string {
	return fmtMs(msSince(t))
}

func msSince(t time.Time) float64 {
	return float64(time.Since(t)) / float64(time.Millisecond)
}

func fmtMs(ms float64) string {
	return fmt.Sprintf("%7.1fms", ms)
}

// withThousands formats an unsigned int with comma thousands separators,
// e.g. 50000 -> "50,000".
func withThousands(n uint64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	out := make([]byte, 0, len(s)+len(s)/3)
	rem := len(s) % 3
	if rem > 0 {
		out = append(out, s[:rem]...)
		if len(s) > rem {
			out = append(out, ',')
		}
	}
	for i := rem; i < len(s); i += 3 {
		out = append(out, s[i:i+3]...)
		if i+3 < len(s) {
			out = append(out, ',')
		}
	}
	return string(out)
}
