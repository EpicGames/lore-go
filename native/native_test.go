// Copyright Epic Games, Inc. All Rights Reserved.

package native

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/EpicGames/lore-go/internal/testutil"
	"github.com/EpicGames/lore-go/types"
)

func TestMain(m *testing.M) {
	if err := testutil.SetLibraryPath(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set library path: %v\n", err)
		os.Exit(1)
	}
	LogConfigure(&types.LoreLogConfigFFI{Level: types.LoreLogLevel_DEBUG})
	os.Exit(m.Run())
}

// createFileWithContents creates a file in the given directory with the specified contents
// Returns the full absolute path to the created file
func createFileWithContents(t *testing.T, repoDir, fileName, fileContents string) string {
	t.Helper()

	filePath := filepath.Join(repoDir, fileName)
	if err := os.WriteFile(filePath, []byte(fileContents), 0644); err != nil {
		t.Fatalf("Failed to create file %s: %v", fileName, err)
	}
	return filePath
}

// createRepository creates a new Lore repository with the given URL
func createRepository(t *testing.T, globals *types.LoreGlobalArgsFFI, url string) {
	t.Helper()

	createArgs, cleanupCreateArgs := types.NewLoreRepositoryCreateArgs(types.LoreRepositoryCreateArgs{
		RepositoryUrl: url,
	})
	defer cleanupCreateArgs()

	result, err := RepositoryCreate(globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil || result != 0 {
		t.Fatalf("Failed to create repository: err=%v, result=%d", err, result)
	}
}

// stageFiles stages the specified files in the repository
func stageFiles(t *testing.T, globals *types.LoreGlobalArgsFFI, paths []string) {
	t.Helper()

	stageArgs, cleanupStageArgs := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
		Paths: paths,
	})
	defer cleanupStageArgs()

	result, err := FileStage(globals, &stageArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil || result != 0 {
		t.Fatalf("Failed to stage files: err=%v, result=%d", err, result)
	}
}

// commitRevision commits the staged changes with the given message
func commitRevision(t *testing.T, globals *types.LoreGlobalArgsFFI, message string) {
	t.Helper()

	commitArgs, cleanupCommitArgs := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
		Message: message,
	})
	defer cleanupCommitArgs()

	result, err := RevisionCommit(globals, &commitArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil || result != 0 {
		t.Fatalf("Failed to commit revision: err=%v, result=%d", err, result)
	}
}

// setupTestRepository creates a temporary Lore repository for testing
func setupTestRepository(t *testing.T) types.LoreGlobalArgsFFI {
	t.Helper()

	if libErr != nil {
		t.Skipf("Lore library not loaded: %v", libErr)
	}

	tempDir, err := os.MkdirTemp("", "lore-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: tempDir,
		Offline:        true,
	})

	// Register cleanup: flush repository, then cleanup globals, then remove temp directory
	t.Cleanup(func() {
		flushArgs, cleanupFlushArgs := types.NewLoreRepositoryFlushArgs(types.LoreRepositoryFlushArgs{})
		defer cleanupFlushArgs()

		RepositoryFlush(&globals, &flushArgs, &types.LoreEventCallbackConfig{
			Callback: func(event *types.LoreEventFFI, userContext uint64) {},
		})

		cleanupGlobals()
		os.RemoveAll(tempDir)
	})

	// Create repository
	repoUrl := fmt.Sprintf("test-repo-%d", time.Now().UnixNano())
	createRepository(t, &globals, repoUrl)

	// Create README.md file
	readmePath := createFileWithContents(t, tempDir, "README.md", "Test")

	// Stage the file
	stageFiles(t, &globals, []string{readmePath})

	// Commit the file
	commitRevision(t, &globals, "initial commit")

	return globals
}

func TestLibraryLoaded(t *testing.T) {
	if libErr != nil {
		t.Fatalf("Lore library failed to load: %v", libErr)
	}
}

func TestLoreVersion(t *testing.T) {
	version, err := Version()

	if err != nil {
		t.Fatalf("LoreVersion failed: %v", err)
	}

	if version == "" {
		t.Fatal("LoreVersion returned empty string")
	}

	t.Logf("Lore library version: %s", version)
}

func TestLoreRepositoryStatus(t *testing.T) {
	globals := setupTestRepository(t)

	// Set up repository status args
	args, cleanupArgs := types.NewLoreRepositoryStatusArgs(types.LoreRepositoryStatusArgs{
		Staged:    true,
		Scan:      true,
		SyncPoint: true,
	})
	defer cleanupArgs()

	createFileWithContents(t, globals.RepositoryPath.String(), "test.txt", "content")
	stagedPath := createFileWithContents(t, globals.RepositoryPath.String(), "staged.txt", "staged content")

	stageFiles(t, &globals, []string{stagedPath})

	// Modify README.md
	readmePath := filepath.Join(globals.RepositoryPath.String(), "README.md")
	if err := os.WriteFile(readmePath, []byte("Modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify README.md: %v", err)
	}

	// Track callback events
	callbackCalled := false
	revisionEventReceived := false
	fileEvents := make(map[string]types.LoreRepositoryStatusFileEventData)

	// Define callback that verifies events
	callback := func(event *types.LoreEventFFI, userContext uint64) {
		callbackCalled = true

		if event.Tag == types.LoreEventTag_REPOSITORY_STATUS_REVISION {
			if revisionEvent, ok := event.GetData().(*types.LoreRepositoryStatusRevisionEventDataFFI); ok {
				revisionEventReceived = true
				branchName := revisionEvent.BranchName.String()
				if branchName != "main" {
					t.Errorf("Expected branch name 'main', got '%s'", branchName)
				}
			}
		}
		if event.Tag == types.LoreEventTag_REPOSITORY_STATUS_FILE {
			if fileEvent, ok := event.GetData().(*types.LoreRepositoryStatusFileEventDataFFI); ok {
				clonedEvent := fileEvent.Clone()
				fileEvents[clonedEvent.Path] = clonedEvent
			}
		}

		if event.Tag == types.LoreEventTag_END {
			t.Logf("Received END event")
		}
	}

	// Call the repository status function
	result, err := RepositoryStatus(&globals, &args, &types.LoreEventCallbackConfig{
		Callback:    callback,
		UserContext: 1,
	})

	if err != nil {
		t.Fatalf("LoreRepositoryStatus failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRepositoryStatus returned non-zero result: %d", result)
	}

	if !callbackCalled {
		t.Error("Callback was not called")
	}

	if !revisionEventReceived {
		t.Error("REPOSITORY_STATUS_REVISION event was not received")
	}

	// Verify staged.txt event
	if fileEvent, ok := fileEvents["staged.txt"]; ok {
		if !fileEvent.FlagStaged {
			t.Errorf("staged.txt: expected FlagStaged=true, got false")
		}
		if fileEvent.Action != types.LoreFileAction_ADD {
			t.Errorf("staged.txt: expected Action=ADD, got %d", fileEvent.Action)
		}
	} else {
		t.Error("staged.txt file event was not received")
	}

	// Verify test.txt event
	if fileEvent, ok := fileEvents["test.txt"]; ok {
		if fileEvent.FlagStaged {
			t.Errorf("test.txt: expected FlagStaged=false, got true")
		}
		if fileEvent.Action != types.LoreFileAction_ADD {
			t.Errorf("test.txt: expected Action=ADD, got %d", fileEvent.Action)
		}
	} else {
		t.Error("test.txt file event was not received")
	}

	// Verify README.md event
	if fileEvent, ok := fileEvents["README.md"]; ok {
		if fileEvent.FlagStaged {
			t.Errorf("README.md: expected FlagStaged=false, got true")
		}
		if fileEvent.Action != types.LoreFileAction_KEEP {
			t.Errorf("README.md: expected Action=KEEP, got %d", fileEvent.Action)
		}
	} else {
		t.Error("README.md file event was not received")
	}
}

func TestLoreRepositoryStatusBranchAndRepositoryIds(t *testing.T) {
	globals := setupTestRepository(t)

	args, cleanupArgs := types.NewLoreRepositoryStatusArgs(types.LoreRepositoryStatusArgs{
		Staged: true,
		Scan:   true,
	})
	defer cleanupArgs()

	// Track the FFI data accessed inside the callback and the cloned data accessed outside
	var ffiRepositoryStr, ffiBranchStr string
	var clonedRevisionEvent types.LoreRepositoryStatusRevisionEventData
	revisionEventReceived := false

	callback := func(event *types.LoreEventFFI, userContext uint64) {
		if event.Tag == types.LoreEventTag_REPOSITORY_STATUS_REVISION {
			if revisionEvent, ok := event.GetData().(*types.LoreRepositoryStatusRevisionEventDataFFI); ok {
				revisionEventReceived = true

				// Access FFI data directly inside the callback
				ffiRepositoryStr = revisionEvent.Repository.String()
				ffiBranchStr = revisionEvent.Branch.String()

				// Verify the binary Data fields are non-zero
				allZeroRepo := true
				for _, b := range revisionEvent.Repository.Data {
					if b != 0 {
						allZeroRepo = false
						break
					}
				}
				if allZeroRepo {
					t.Error("FFI Repository.Data is all zeros")
				}

				allZeroBranch := true
				for _, b := range revisionEvent.Branch.Data {
					if b != 0 {
						allZeroBranch = false
						break
					}
				}
				if allZeroBranch {
					t.Error("FFI Branch.Data is all zeros")
				}

				// Clone the event for use outside the callback
				clonedRevisionEvent = revisionEvent.Clone()
			}
		}
	}

	result, err := RepositoryStatus(&globals, &args, &types.LoreEventCallbackConfig{
		Callback:    callback,
		UserContext: 1,
	})

	if err != nil {
		t.Fatalf("LoreRepositoryStatus failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreRepositoryStatus returned non-zero result: %d", result)
	}
	if !revisionEventReceived {
		t.Fatal("REPOSITORY_STATUS_REVISION event was not received")
	}

	// Verify the FFI String() output is a valid 32-char hex string (16 bytes)
	if len(ffiRepositoryStr) != 32 {
		t.Errorf("FFI Repository.String() length: expected 32, got %d (%q)", len(ffiRepositoryStr), ffiRepositoryStr)
	}
	if len(ffiBranchStr) != 32 {
		t.Errorf("FFI Branch.String() length: expected 32, got %d (%q)", len(ffiBranchStr), ffiBranchStr)
	}

	// Verify cloned data matches the FFI data observed during the callback
	clonedRepoStr := clonedRevisionEvent.Repository.String()
	clonedBranchStr := clonedRevisionEvent.Branch.String()

	if clonedRepoStr != ffiRepositoryStr {
		t.Errorf("Cloned Repository.String() = %q, FFI had %q", clonedRepoStr, ffiRepositoryStr)
	}
	if clonedBranchStr != ffiBranchStr {
		t.Errorf("Cloned Branch.String() = %q, FFI had %q", clonedBranchStr, ffiBranchStr)
	}

	// Verify the cloned binary Data is non-zero (survived outside the callback)
	allZeroRepo := true
	for _, b := range clonedRevisionEvent.Repository.Data {
		if b != 0 {
			allZeroRepo = false
			break
		}
	}
	if allZeroRepo {
		t.Error("Cloned Repository.Data is all zeros")
	}

	allZeroBranch := true
	for _, b := range clonedRevisionEvent.Branch.Data {
		if b != 0 {
			allZeroBranch = false
			break
		}
	}
	if allZeroBranch {
		t.Error("Cloned Branch.Data is all zeros")
	}

	t.Logf("Repository ID: %s", clonedRepoStr)
	t.Logf("Branch ID:     %s", clonedBranchStr)
}

func TestLoreRevisionHistory(t *testing.T) {
	globals := setupTestRepository(t)

	// Create another file and commit it
	secondFilePath := createFileWithContents(t, globals.RepositoryPath.String(), "second.txt", "second commit content")
	stageFiles(t, &globals, []string{secondFilePath})
	commitRevision(t, &globals, "second commit")

	// Set up revision history args (empty revision and branch to get all history)
	args, cleanupArgs := types.NewLoreRevisionHistoryArgs(types.LoreRevisionHistoryArgs{})
	defer cleanupArgs()

	// Track events
	historyEntries := []types.LoreEvent{}
	commitMessages := []string{}

	callback := func(event *types.LoreEventFFI, userContext uint64) {
		if event.Tag == types.LoreEventTag_REVISION_HISTORY_ENTRY {
			clonedEvent := event.Clone()
			historyEntries = append(historyEntries, clonedEvent)
		}

		if event.Tag == types.LoreEventTag_METADATA {
			if metadataEvent, ok := event.GetData().(*types.LoreMetadataEventDataFFI); ok {
				key := metadataEvent.Key.String()
				if key == "message" && metadataEvent.Value.Tag == types.LoreMetadataTag_STRING {
					messageStr := metadataEvent.Value.AsLoreString().String()
					commitMessages = append(commitMessages, messageStr)
				}
			}
		}
	}

	// Call revision history
	result, err := RevisionHistory(&globals, &args, &types.LoreEventCallbackConfig{
		Callback:    callback,
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionHistory failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionHistory returned non-zero result: %d", result)
	}

	// Verify we received 2 revision history entries
	if len(historyEntries) != 2 {
		t.Errorf("Expected 2 revision history entries, got %d", len(historyEntries))
	}

	// Verify we received commit messages
	if len(commitMessages) < 2 {
		t.Errorf("Expected at least 2 commit messages, got %d", len(commitMessages))
	}

	// Verify commit messages content
	expectedMessages := []string{"initial commit", "second commit"}
	for _, expected := range expectedMessages {
		found := false
		for _, actual := range commitMessages {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected commit message '%s' not found in: %v", expected, commitMessages)
		}
	}
}

func TestLoreRevisionAmend(t *testing.T) {
	globals := setupTestRepository(t)

	// Create another file and commit with original message
	testFilePath := createFileWithContents(t, globals.RepositoryPath.String(), "amend-test.txt", "test content")
	stageFiles(t, &globals, []string{testFilePath})

	originalMessage := "original commit message"
	commitRevision(t, &globals, originalMessage)

	// Amend the commit with a new message
	amendedMessage := "amended commit message"
	amendArgs, cleanupAmendArgs := types.NewLoreRevisionAmendArgs(types.LoreRevisionAmendArgs{
		Message: amendedMessage,
	})
	defer cleanupAmendArgs()

	result, err := RevisionAmend(&globals, &amendArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionAmend failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionAmend returned non-zero result: %d", result)
	}

	// Get revision history to verify the amended message
	historyArgs, cleanupHistoryArgs := types.NewLoreRevisionHistoryArgs(types.LoreRevisionHistoryArgs{})
	defer cleanupHistoryArgs()

	commitMessages := []string{}

	result, err = RevisionHistory(&globals, &historyArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_METADATA {
				if metadataEvent, ok := event.GetData().(*types.LoreMetadataEventDataFFI); ok {
					key := metadataEvent.Key.String()
					if key == "message" && metadataEvent.Value.Tag == types.LoreMetadataTag_STRING {
						messageStr := metadataEvent.Value.AsLoreString().String()
						commitMessages = append(commitMessages, messageStr)
					}
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionHistory failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionHistory returned non-zero result: %d", result)
	}

	// Verify the amended message is present
	amendedFound := false
	for _, msg := range commitMessages {
		if msg == amendedMessage {
			amendedFound = true
			t.Logf("Found amended commit message: %s", amendedMessage)
		}
	}

	if !amendedFound {
		t.Errorf("Amended message '%s' not found in commit messages: %v", amendedMessage, commitMessages)
	}

	// Verify the original message is NOT present
	originalFound := false
	for _, msg := range commitMessages {
		if msg == originalMessage {
			originalFound = true
		}
	}

	if originalFound {
		t.Errorf("Original message '%s' should not be present after amend, but was found in: %v", originalMessage, commitMessages)
	}

	t.Logf("Successfully amended commit message from '%s' to '%s'", originalMessage, amendedMessage)
}

func TestLoreBranchMerge(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a feature branch
	featureBranch := "feature-branch"
	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: featureBranch,
	})
	defer cleanupCreateArgs()

	result, err := BranchCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate returned non-zero result: %d", result)
	}

	// Create and commit a file on the feature branch
	conflictFile := "conflict-file.txt"
	conflictFilePath := filepath.Join(globals.RepositoryPath.String(), conflictFile)
	featureContent := "Feature branch content"
	if err := os.WriteFile(conflictFilePath, []byte(featureContent), 0644); err != nil {
		t.Fatalf("Failed to create conflict file on feature branch: %v", err)
	}
	stageFiles(t, &globals, []string{conflictFilePath})
	commitRevision(t, &globals, "commit on feature branch")

	// Switch back to main branch
	switchArgs, cleanupSwitchArgs := types.NewLoreBranchSwitchArgs(types.LoreBranchSwitchArgs{
		Branch: "main",
	})
	defer cleanupSwitchArgs()

	result, err = BranchSwitch(&globals, &switchArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchSwitch to main failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchSwitch returned non-zero result: %d", result)
	}

	// Commit the same file with conflicting content on main
	mainContent := "Main branch content - conflicts with feature"
	if err := os.WriteFile(conflictFilePath, []byte(mainContent), 0644); err != nil {
		t.Fatalf("Failed to create conflict file on main branch: %v", err)
	}
	stageFiles(t, &globals, []string{conflictFilePath})
	commitRevision(t, &globals, "commit on main branch")

	// Start merge from feature branch to main
	mergeStartArgs, cleanupMergeStartArgs := types.NewLoreBranchMergeStartArgs(types.LoreBranchMergeStartArgs{
		Branch:   featureBranch,
		Message:  "merge feature branch",
		NoCommit: false,
	})
	defer cleanupMergeStartArgs()

	conflictFound := false
	conflictFilePaths := []string{}

	result, err = BranchMergeStart(&globals, &mergeStartArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_MERGE_CONFLICT_FILE {
				conflictFound = true
				if conflictEvent, ok := event.GetData().(*types.LoreBranchMergeConflictFileEventDataFFI); ok {
					path := conflictEvent.Path.String()
					conflictFilePaths = append(conflictFilePaths, path)
					t.Logf("BRANCH_MERGE_CONFLICT_FILE: %s", path)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchMergeStart failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchMergeStart returned non-zero result: %d", result)
	}

	// Verify conflict was detected
	if !conflictFound {
		t.Fatal("Expected BRANCH_MERGE_CONFLICT_FILE event, but none was received")
	}

	if len(conflictFilePaths) == 0 {
		t.Fatal("No conflict files detected")
	}

	t.Logf("Detected %d conflict file(s): %v", len(conflictFilePaths), conflictFilePaths)

	// Resolve conflict using "mine" (main branch version)
	resolveMineArgs, cleanupResolveMineArgs := types.NewLoreBranchMergeResolveMineArgs(types.LoreBranchMergeResolveMineArgs{
		Paths: []string{conflictFilePath},
	})
	defer cleanupResolveMineArgs()

	stageFileEventFound := false

	result, err = BranchMergeResolveMine(&globals, &resolveMineArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_STAGE_FILE {
				stageFileEventFound = true
				if stageEvent, ok := event.GetData().(*types.LoreFileStageFileEventDataFFI); ok {
					path := stageEvent.Path.String()
					t.Logf("FILE_STAGE_FILE: %s", path)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchMergeResolveMine failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchMergeResolveMine returned non-zero result: %d", result)
	}

	// Verify FILE_STAGE_FILE event was received
	if !stageFileEventFound {
		t.Fatal("Expected FILE_STAGE_FILE event after conflict resolution, but none was received")
	}

	// Commit the merge
	commitRevision(t, &globals, "merge commit after resolving conflicts")

	t.Logf("Successfully merged %s into main with conflict resolution", featureBranch)
}

func TestLoreBranchCreate(t *testing.T) {
	globals := setupTestRepository(t)

	// Get main branch info before creating new branch
	infoArgs, cleanupInfoArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{})
	defer cleanupInfoArgs()

	var mainBranchInfo *types.LoreBranchInfoEventData

	result, err := BranchInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreBranchInfoEventDataFFI); ok {
					clonedInfo := infoEvent.Clone()
					mainBranchInfo = &clonedInfo
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchInfo (main) failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchInfo (main) returned non-zero result: %d", result)
	}

	if mainBranchInfo == nil {
		t.Fatal("No BRANCH_INFO event received for main branch")
	}

	if mainBranchInfo.Name != "main" {
		t.Errorf("Expected main branch name 'main', got '%s'", mainBranchInfo.Name)
	}

	// Create a new branch
	branchName := "feature-branch"
	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: branchName,
	})
	defer cleanupCreateArgs()

	result, err = BranchCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate returned non-zero result: %d", result)
	}

	// Get feature branch info (empty string means current branch)
	var featureBranchInfo *types.LoreBranchInfoEventData

	result, err = BranchInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreBranchInfoEventDataFFI); ok {
					clonedInfo := infoEvent.Clone()
					featureBranchInfo = &clonedInfo
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchInfo (feature) failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchInfo (feature) returned non-zero result: %d", result)
	}

	// Verify branch info was received
	if featureBranchInfo == nil {
		t.Fatal("No BRANCH_INFO event received for feature branch")
	}

	// Verify branch name matches
	if featureBranchInfo.Name != branchName {
		t.Errorf("Expected branch name '%s', got '%s'", branchName, featureBranchInfo.Name)
	}

	// Verify parent matches main branch ID
	if featureBranchInfo.Parent != mainBranchInfo.Id {
		t.Errorf("Expected feature branch parent to be main branch ID %s, got %s",
			mainBranchInfo.Id, featureBranchInfo.Parent)
	}
}

func TestLoreBranchArchive(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a new branch
	branchName := "branch-to-archive"
	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: branchName,
	})
	defer cleanupCreateArgs()

	result, err := BranchCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate returned non-zero result: %d", result)
	}

	// Switch back to main branch
	switchArgs, cleanupSwitchArgs := types.NewLoreBranchSwitchArgs(types.LoreBranchSwitchArgs{
		Branch: "main",
	})
	defer cleanupSwitchArgs()

	result, err = BranchSwitch(&globals, &switchArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchSwitch failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchSwitch returned non-zero result: %d", result)
	}

	// Archive the branch
	archiveArgs, cleanupArchiveArgs := types.NewLoreBranchArchiveArgs(types.LoreBranchArchiveArgs{
		Branch: branchName,
	})
	defer cleanupArchiveArgs()

	result, err = BranchArchive(&globals, &archiveArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchArchive failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchArchive returned non-zero result: %d", result)
	}

	// List branches to verify the archived branch is gone
	listArgs, cleanupListArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupListArgs()

	branchNames := []string{}

	result, err = BranchList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_LIST_ENTRY {
				if entryEvent, ok := event.GetData().(*types.LoreBranchListEntryEventDataFFI); ok {
					clonedEntry := entryEvent.Clone()
					branchNames = append(branchNames, clonedEntry.Name)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchList failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchList returned non-zero result: %d", result)
	}

	// Verify the archived branch is not in the list
	for _, name := range branchNames {
		if name == branchName {
			t.Errorf("Archived branch '%s' still exists in branch list", branchName)
		}
	}

	// Verify main branch still exists
	mainFound := false
	for _, name := range branchNames {
		if name == "main" {
			mainFound = true
			break
		}
	}

	if !mainFound {
		t.Error("Main branch not found in branch list")
	}

	t.Logf("Successfully deleted branch '%s'. Remaining branches: %v", branchName, branchNames)
}

func TestLoreRevisionDiff(t *testing.T) {
	globals := setupTestRepository(t)

	// Get current revision info
	infoArgs, cleanupInfoArgs := types.NewLoreRevisionInfoArgs(types.LoreRevisionInfoArgs{
		Metadata: false,
	})
	defer cleanupInfoArgs()

	var originalRevisionHash string

	result, err := RevisionInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreRevisionInfoEventDataFFI); ok {
					originalRevisionHash = infoEvent.Revision.String()
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionInfo failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionInfo returned non-zero result: %d", result)
	}

	// Create and commit another file
	newFilePath := createFileWithContents(t, globals.RepositoryPath.String(), "newfile.txt", "new content")
	stageFiles(t, &globals, []string{newFilePath})
	commitRevision(t, &globals, "add new file")

	// Call RevisionDiff to compare original revision with current
	diffArgs, cleanupDiffArgs := types.NewLoreRevisionDiffArgs(types.LoreRevisionDiffArgs{
		RevisionSource: originalRevisionHash,
	})
	defer cleanupDiffArgs()

	diffFiles := []types.LoreRevisionDiffFileEventData{}

	result, err = RevisionDiff(&globals, &diffArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_DIFF_FILE {
				if diffEvent, ok := event.GetData().(*types.LoreRevisionDiffFileEventDataFFI); ok {
					clonedDiff := diffEvent.Clone()
					diffFiles = append(diffFiles, clonedDiff)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionDiff failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionDiff returned non-zero result: %d", result)
	}

	// Verify we received at least one diff file event
	if len(diffFiles) == 0 {
		t.Fatal("No REVISION_DIFF_FILE events received")
	}

	// Find the newfile.txt diff event
	var newFileDiff *types.LoreRevisionDiffFileEventData
	for i := range diffFiles {
		if diffFiles[i].Path == "newfile.txt" {
			newFileDiff = &diffFiles[i]
			break
		}
	}

	if newFileDiff == nil {
		t.Fatal("No REVISION_DIFF_FILE event found for newfile.txt")
	}

	// Verify the action is ADD
	if newFileDiff.Action != types.LoreFileAction_ADD {
		t.Errorf("Expected Action=ADD for newfile.txt, got %d", newFileDiff.Action)
	}
}

func TestLoreFileUnstage(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a new file
	newFilePath := createFileWithContents(t, globals.RepositoryPath.String(), "unstage-test.txt", "test content")

	// Stage the file
	stageArgs, cleanupStageArgs := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
		Paths: []string{newFilePath},
	})
	defer cleanupStageArgs()

	var stageCountData *types.LoreFileStageCountData

	result, err := FileStage(&globals, &stageArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_STAGE_END {
				if endEvent, ok := event.GetData().(*types.LoreFileStageEndEventDataFFI); ok {
					clone := endEvent.Clone()
					stageCountData = &clone.Count
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileStage failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileStage returned non-zero result: %d", result)
	}

	// Verify stage count data
	if stageCountData == nil {
		t.Fatal("No FILE_STAGE_END event received")
	}

	if stageCountData.FileAddCount != 1 {
		t.Errorf("Expected FileAddCount=1, got %d", stageCountData.FileAddCount)
	}

	if stageCountData.TotalCount != 1 {
		t.Errorf("Expected TotalCount=1, got %d", stageCountData.TotalCount)
	}

	// Unstage the file
	unstageArgs, cleanupUnstageArgs := types.NewLoreFileUnstageArgs(types.LoreFileUnstageArgs{
		Paths: []string{newFilePath},
	})
	defer cleanupUnstageArgs()

	var unstageCountData *types.LoreFileUnstageCountData

	result, err = FileUnstage(&globals, &unstageArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_UNSTAGE_END {
				if endEvent, ok := event.GetData().(*types.LoreFileUnstageEndEventDataFFI); ok {
					clone := endEvent.Clone()
					unstageCountData = &clone.Count
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileUnstage failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileUnstage returned non-zero result: %d", result)
	}

	// Verify unstage count data
	if unstageCountData == nil {
		t.Fatal("No FILE_UNSTAGE_END event received")
	}

	if unstageCountData.FileUnstagedCount != 1 {
		t.Errorf("Expected FileUnstagedCount=1, got %d", unstageCountData.FileUnstagedCount)
	}

	if unstageCountData.TotalCount != 1 {
		t.Errorf("Expected TotalCount=1, got %d", unstageCountData.TotalCount)
	}
}

func TestLoreUnicodeSupport(t *testing.T) {
	globals := setupTestRepository(t)

	// Get initial revision info
	infoArgs, cleanupInfoArgs := types.NewLoreRevisionInfoArgs(types.LoreRevisionInfoArgs{
		Metadata: false,
	})
	defer cleanupInfoArgs()

	var originalRevisionHash string

	result, err := RevisionInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreRevisionInfoEventDataFFI); ok {
					originalRevisionHash = infoEvent.Revision.String()
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionInfo failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionInfo returned non-zero result: %d", result)
	}

	// Create file with multibyte unicode characters in name and content
	unicodeFileName := "öäÄÅ的ЛЛЛµÐ𒁻𒁃𓉡𓉢‼️🌏文件名-こんにちは-🚀🇩🇪.txt"
	unicodeContent := "Hello 世界! Привет мир! 🎉 This is a test with emoji 😊 and various scripts: 日本語, Русский, العربية, öäÄÅ的ЛЛЛµÐ𒁻𒁃𓉡𓉢‼️🌏🇩🇪"
	unicodeFilePath := createFileWithContents(t, globals.RepositoryPath.String(), unicodeFileName, unicodeContent)

	// Stage the file and verify the path in FILE_STAGE_FILE event
	stageArgs, cleanupStageArgs := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
		Paths: []string{unicodeFilePath},
	})
	defer cleanupStageArgs()

	var stagedFilePath string

	result, err = FileStage(&globals, &stageArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_STAGE_FILE {
				if fileEvent, ok := event.GetData().(*types.LoreFileStageFileEventDataFFI); ok {
					stagedFilePath = fileEvent.Path.String()
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileStage failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileStage returned non-zero result: %d", result)
	}

	// Verify the staged file path contains the unicode filename
	if stagedFilePath != unicodeFileName {
		t.Errorf("Expected staged file path '%s', got '%s'", unicodeFileName, stagedFilePath)
	}

	// Commit the file
	commitMessage := "add unicode file with öäÄÅ的ЛЛЛµÐ𒁻𒁃𓉡𓉢‼️🌏🇩🇪 in commit message"
	commitRevision(t, &globals, commitMessage)

	// Call FileDiff to get the diff patch
	diffArgs, cleanupDiffArgs := types.NewLoreFileDiffArgs(types.LoreFileDiffArgs{
		Paths:          []string{unicodeFilePath},
		SourceRevision: originalRevisionHash,
		TargetRevision: "", // Current revision
	})
	defer cleanupDiffArgs()

	var fileDiffPatch string

	result, err = FileDiff(&globals, &diffArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_DIFF {
				if diffEvent, ok := event.GetData().(*types.LoreFileDiffEventDataFFI); ok {
					clonedDiff := diffEvent.Clone()
					fileDiffPatch = clonedDiff.Patch
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileDiff failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileDiff returned non-zero result: %d", result)
	}

	// Verify the patch contains the multibyte unicode content
	if fileDiffPatch == "" {
		t.Fatal("No FILE_DIFF patch received")
	}

	// Check that the patch contains the unicode content
	if !strings.Contains(fileDiffPatch, unicodeContent) {
		t.Errorf("Expected patch to contain '%s', but it was not found", unicodeContent)
	}

	// Call RevisionHistory to verify commit message with unicode
	historyArgs, cleanupHistoryArgs := types.NewLoreRevisionHistoryArgs(types.LoreRevisionHistoryArgs{})
	defer cleanupHistoryArgs()

	commitMessages := []string{}

	result, err = RevisionHistory(&globals, &historyArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_METADATA {
				if metadataEvent, ok := event.GetData().(*types.LoreMetadataEventDataFFI); ok {
					key := metadataEvent.Key.String()
					if key == "message" && metadataEvent.Value.Tag == types.LoreMetadataTag_STRING {
						messageStr := metadataEvent.Value.AsLoreString().String()
						commitMessages = append(commitMessages, messageStr)
					}
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionHistory failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionHistory returned non-zero result: %d", result)
	}

	// Verify the commit message with unicode is present
	found := false
	for _, msg := range commitMessages {
		if msg == commitMessage {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected commit message '%s' not found in revision history. Found messages: %v", commitMessage, commitMessages)
	}
}

func TestLoreBranchUnicodeNames(t *testing.T) {
	globals := setupTestRepository(t)

	// Create first branch with unicode characters in name
	unicodeBranch1 := "feature-🚀-世界-Привет-öäÄÅ-🌏🇩🇪"
	createArgs1, cleanupCreateArgs1 := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: unicodeBranch1,
	})
	defer cleanupCreateArgs1()

	result, err := BranchCreate(&globals, &createArgs1, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate (unicode1) failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate (unicode1) returned non-zero result: %d", result)
	}

	// Create second branch with different unicode characters
	unicodeBranch2 := "branch-日本語-🎉-𒁻𒁃"
	createArgs2, cleanupCreateArgs2 := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: unicodeBranch2,
	})
	defer cleanupCreateArgs2()

	result, err = BranchCreate(&globals, &createArgs2, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate (unicode2) failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate (unicode2) returned non-zero result: %d", result)
	}

	// Switch back to the first unicode branch
	switchArgs, cleanupSwitchArgs := types.NewLoreBranchSwitchArgs(types.LoreBranchSwitchArgs{
		Branch: unicodeBranch1,
	})
	defer cleanupSwitchArgs()

	result, err = BranchSwitch(&globals, &switchArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchSwitch failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchSwitch returned non-zero result: %d", result)
	}

	// List branches and verify unicode branches are present
	listArgs, cleanupListArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupListArgs()

	branchNames := []string{}

	result, err = BranchList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_LIST_ENTRY {
				if entryEvent, ok := event.GetData().(*types.LoreBranchListEntryEventDataFFI); ok {
					clonedEntry := entryEvent.Clone()
					branchNames = append(branchNames, clonedEntry.Name)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchList failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchList returned non-zero result: %d", result)
	}

	// Verify both unicode branches are in the list
	expectedBranches := []string{"main", unicodeBranch1, unicodeBranch2}
	for _, expected := range expectedBranches {
		found := false
		for _, actual := range branchNames {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected branch '%s' not found in branch list. Found branches: %v", expected, branchNames)
		}
	}

	// Verify first unicode branch is current using LoreBranchInfo
	infoArgs, cleanupInfoArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{})
	defer cleanupInfoArgs()

	var currentBranchInfo *types.LoreBranchInfoEventData

	result, err = BranchInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreBranchInfoEventDataFFI); ok {
					clonedInfo := infoEvent.Clone()
					currentBranchInfo = &clonedInfo
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchInfo failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchInfo returned non-zero result: %d", result)
	}

	if currentBranchInfo == nil {
		t.Fatal("No BRANCH_INFO event received")
	}

	if currentBranchInfo.Name != unicodeBranch1 {
		t.Errorf("Expected current branch to be '%s', got '%s'", unicodeBranch1, currentBranchInfo.Name)
	}
}

func TestLoreBranchListStackHandling(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a new branch
	branchName := "test-branch"
	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: branchName,
	})
	defer cleanupCreateArgs()

	result, err := BranchCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate returned non-zero result: %d", result)
	}

	// List branches and verify we can access the Stack field inside the callback
	listArgs, cleanupListArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupListArgs()

	branchEntries := make(map[string]types.LoreBranchListEntryEventData)
	callbackCalled := false

	result, err = BranchList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_LIST_ENTRY {
				callbackCalled = true
				if entryEvent, ok := event.GetData().(*types.LoreBranchListEntryEventDataFFI); ok {
					// Access the Stack field (branch_point_array) inside the callback
					// This verifies that the array data is accessible during the callback
					eventBranchName := entryEvent.Name.String()

					// 1. Get the length
					stackLen := entryEvent.Stack.Len()
					t.Logf("Branch %s: Stack has %d entries", eventBranchName, stackLen)

					// 2. Access individual elements by index (if any exist)
					if eventBranchName == branchName {
						if stackLen != 1 {
							t.Errorf("Expected stack length 1 for feature branch")
						}
					}
					if stackLen > 0 {
						firstEntry := entryEvent.Stack.Get(0)
						branchId := firstEntry.Branch.String()
						t.Logf("Branch %s: First stack entry - Parent=%s", eventBranchName, branchId)
					}

					// Clone the entire event to convert FFI data to Go data
					clonedEntry := entryEvent.Clone()

					// Store the cloned entry for verification outside the callback
					branchEntries[eventBranchName] = clonedEntry
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchList failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchList returned non-zero result: %d", result)
	}

	if !callbackCalled {
		t.Fatal("BRANCH_LIST_ENTRY callback was not called")
	}

	// Verify we received branch entries
	if len(branchEntries) == 0 {
		t.Fatal("No branch entries were collected")
	}

	// Verify main branch exists in the results and can access its data outside callback
	if mainEntry, ok := branchEntries["main"]; ok {
		// Access the Name field outside the callback - this will demonstrate if Clone() works properly
		t.Logf("Outside callback: main branch Name=%s, Creator=%s", mainEntry.Name, mainEntry.Creator)
		if mainEntry.Name != "main" {
			t.Errorf("Expected Name='main', got '%s'", mainEntry.Name)
		}
	} else {
		t.Errorf("Expected 'main' branch in results, got branches: %v", getBranchNames(branchEntries))
	}

	// Verify test branch exists in the results and can access its data outside callback
	if testEntry, ok := branchEntries[branchName]; ok {
		// Access fields outside the callback
		t.Logf("Outside callback: test-branch Name=%s, Creator=%s", testEntry.Name, testEntry.Creator)

		// The test branch should have a stack with at least one entry (the branch point from main)
		stackLen := len(testEntry.Stack)
		t.Logf("Outside callback: Stack length = %d", stackLen)

		if stackLen == 0 {
			t.Errorf("Expected test branch '%s' to have at least one branch point in stack, got 0", branchName)
		}
	} else {
		t.Errorf("Expected '%s' branch in results, got branches: %v", branchName, getBranchNames(branchEntries))
	}
}

// Helper function to get branch names from map
func getBranchNames(branchData map[string]types.LoreBranchListEntryEventData) []string {
	names := make([]string, 0, len(branchData))
	for name := range branchData {
		names = append(names, name)
	}
	return names
}

func TestLoreBranchDiffChangeEventAccess(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a new branch
	branchName := "feature-branch"
	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: branchName,
	})
	defer cleanupCreateArgs()

	result, err := BranchCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchCreate failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchCreate returned non-zero result: %d", result)
	}

	// Create and commit a new file on the feature branch
	newFilePath := createFileWithContents(t, globals.RepositoryPath.String(), "feature.txt", "feature content")
	stageFiles(t, &globals, []string{newFilePath})
	commitRevision(t, &globals, "add feature file")

	// Modify README.md on the feature branch
	readmePath := filepath.Join(globals.RepositoryPath.String(), "README.md")
	if err := os.WriteFile(readmePath, []byte("Modified on feature branch"), 0644); err != nil {
		t.Fatalf("Failed to modify README.md: %v", err)
	}
	stageFiles(t, &globals, []string{readmePath})
	commitRevision(t, &globals, "update README")

	// Call BranchDiff to compare main with feature branch
	diffArgs, cleanupDiffArgs := types.NewLoreBranchDiffArgs(types.LoreBranchDiffArgs{
		Source: branchName,
		Target: "main",
	})
	defer cleanupDiffArgs()

	// Track BRANCH_DIFF_CHANGE events
	changeEvents := []types.LoreBranchDiffChangeEventData{}
	callbackCalled := false

	result, err = BranchDiff(&globals, &diffArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_DIFF_CHANGE {
				callbackCalled = true
				if changeEvent, ok := event.GetData().(*types.LoreBranchDiffChangeEventDataFFI); ok {
					// Verify that we can access Change.Path and Change.Action inside the callback through FFI data
					path := changeEvent.Change.Path.String()
					action := changeEvent.Change.Action

					t.Logf("BRANCH_DIFF_CHANGE: Path=%s, Action=%d", path, action)

					// Clone the event for storage outside the callback
					clonedEvent := changeEvent.Clone()
					changeEvents = append(changeEvents, clonedEvent)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchDiff failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreBranchDiff returned non-zero result: %d", result)
	}

	if !callbackCalled {
		t.Error("BRANCH_DIFF_CHANGE callback was not called")
	}

	// Verify we received change events
	if len(changeEvents) == 0 {
		t.Fatal("No BRANCH_DIFF_CHANGE events received")
	}

	// Verify feature.txt change event
	var featureTxtEvent *types.LoreBranchDiffChangeEventData
	for i := range changeEvents {
		if changeEvents[i].Change.Path == "feature.txt" {
			featureTxtEvent = &changeEvents[i]
			break
		}
	}

	if featureTxtEvent == nil {
		t.Error("No BRANCH_DIFF_CHANGE event found for feature.txt")
	} else {
		// Verify we can access the Path field
		if featureTxtEvent.Change.Path != "feature.txt" {
			t.Errorf("Expected Path='feature.txt', got '%s'", featureTxtEvent.Change.Path)
		}

		// Verify we can access the Action field (should be ADD)
		if featureTxtEvent.Change.Action != types.LoreFileAction_ADD {
			t.Errorf("Expected Action=ADD for feature.txt, got %d", featureTxtEvent.Change.Action)
		}
	}

	// Verify README.md change event
	var readmeEvent *types.LoreBranchDiffChangeEventData
	for i := range changeEvents {
		if changeEvents[i].Change.Path == "README.md" {
			readmeEvent = &changeEvents[i]
			break
		}
	}

	if readmeEvent == nil {
		t.Error("No BRANCH_DIFF_CHANGE event found for README.md")
	} else {
		// Verify we can access the Path field
		if readmeEvent.Change.Path != "README.md" {
			t.Errorf("Expected Path='README.md', got '%s'", readmeEvent.Change.Path)
		}

		// Verify we can access the Action field (should be KEEP for modified file)
		if readmeEvent.Change.Action != types.LoreFileAction_KEEP {
			t.Errorf("Expected Action=KEEP for README.md, got %d", readmeEvent.Change.Action)
		}
	}
}

func TestLoreHashContextAddressClone(t *testing.T) {
	// Test LoreHash.Clone()
	t.Run("LoreHash", func(t *testing.T) {
		original := types.LoreHash{
			Data: [32]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
				17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		}

		cloned := original.Clone()

		// Verify data is copied
		for i := 0; i < 32; i++ {
			if cloned.Data[i] != original.Data[i] {
				t.Errorf("Index %d: expected %d, got %d", i, original.Data[i], cloned.Data[i])
			}
		}

		// Modify cloned to ensure it's independent
		cloned.Data[0] = 99
		if original.Data[0] == 99 {
			t.Error("Modifying cloned data affected original - not a proper copy")
		}

		// Verify String() works
		hashStr := original.String()
		if len(hashStr) != 64 { // 32 bytes = 64 hex chars
			t.Errorf("Expected hash string length 64, got %d", len(hashStr))
		}
	})

	// Test LoreContext.Clone()
	t.Run("LoreContext", func(t *testing.T) {
		original := types.LoreContext{
			Data: [16]uint8{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160},
		}

		cloned := original.Clone()

		// Verify data is copied
		for i := 0; i < 16; i++ {
			if cloned.Data[i] != original.Data[i] {
				t.Errorf("Index %d: expected %d, got %d", i, original.Data[i], cloned.Data[i])
			}
		}

		// Modify cloned to ensure it's independent
		cloned.Data[0] = 255
		if original.Data[0] == 255 {
			t.Error("Modifying cloned data affected original - not a proper copy")
		}

		// Verify String() works
		ctxStr := original.String()
		if len(ctxStr) != 32 { // 16 bytes = 32 hex chars
			t.Errorf("Expected context string length 32, got %d", len(ctxStr))
		}
	})

	// Test LoreAddress.Clone()
	t.Run("LoreAddress", func(t *testing.T) {
		original := types.LoreAddress{
			Hash: types.LoreHash{
				Data: [32]uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
			},
			Context: types.LoreContext{
				Data: [16]uint8{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160},
			},
		}

		cloned := original.Clone()

		// Verify Hash is copied
		for i := 0; i < 32; i++ {
			if cloned.Hash.Data[i] != original.Hash.Data[i] {
				t.Errorf("Hash Index %d: expected %d, got %d", i, original.Hash.Data[i], cloned.Hash.Data[i])
			}
		}

		// Verify Context is copied
		for i := 0; i < 16; i++ {
			if cloned.Context.Data[i] != original.Context.Data[i] {
				t.Errorf("Context Index %d: expected %d, got %d", i, original.Context.Data[i], cloned.Context.Data[i])
			}
		}

		// Modify cloned to ensure it's independent
		cloned.Hash.Data[0] = 99
		cloned.Context.Data[0] = 255
		if original.Hash.Data[0] == 99 || original.Context.Data[0] == 255 {
			t.Error("Modifying cloned data affected original - not a proper copy")
		}
	})
}

func TestLoreRevisionInfoWithHashArray(t *testing.T) {
	globals := setupTestRepository(t)

	// Get revision info which includes Parent [2]LoreHash field
	infoArgs, cleanupInfoArgs := types.NewLoreRevisionInfoArgs(types.LoreRevisionInfoArgs{
		Metadata: false,
	})
	defer cleanupInfoArgs()

	var revisionInfo *types.LoreRevisionInfoEventData

	result, err := RevisionInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_INFO {
				if infoEvent, ok := event.GetData().(*types.LoreRevisionInfoEventDataFFI); ok {
					// Clone the event - this should properly clone the [2]LoreHash Parent array
					clonedInfo := infoEvent.Clone()
					revisionInfo = &clonedInfo
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionInfo failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionInfo returned non-zero result: %d", result)
	}

	if revisionInfo == nil {
		t.Fatal("No REVISION_INFO event received")
	}

	// Verify we can access the Parent array outside the callback
	// The Parent field is [2]LoreHash - both parents should be accessible
	t.Logf("Parent[0]: %s", revisionInfo.Parent[0].String())
	t.Logf("Parent[1]: %s", revisionInfo.Parent[1].String())

	// Verify the Parent array is properly cloned (should be all zeros for initial commit)
	allZeros := true
	for i := 0; i < 2; i++ {
		for j := 0; j < 32; j++ {
			if revisionInfo.Parent[i].Data[j] != 0 {
				allZeros = false
				break
			}
		}
	}

	if !allZeros {
		t.Logf("Parent hashes are non-zero (this is expected for non-initial commits)")
	}

	// Verify Repository and Revision fields are also cloned
	t.Logf("Repository: %s", revisionInfo.Repository.String())
	t.Logf("Revision: %s", revisionInfo.Revision.String())
	t.Logf("RevisionNumber: %d", revisionInfo.RevisionNumber)
}

func TestLoreFileHistoryWithAddress(t *testing.T) {
	globals := setupTestRepository(t)

	// Create and modify a file across multiple commits
	testFile := "history-test.txt"
	testFilePath := filepath.Join(globals.RepositoryPath.String(), testFile)

	// First commit - create file
	if err := os.WriteFile(testFilePath, []byte("Version 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	stageFiles(t, &globals, []string{testFilePath})
	commitRevision(t, &globals, "add history-test.txt v1")

	// Second commit - modify file
	if err := os.WriteFile(testFilePath, []byte("Version 2 - updated"), 0644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}
	stageFiles(t, &globals, []string{testFilePath})
	commitRevision(t, &globals, "update history-test.txt v2")

	// Third commit - modify file again
	if err := os.WriteFile(testFilePath, []byte("Version 3 - final update"), 0644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}
	stageFiles(t, &globals, []string{testFilePath})
	commitRevision(t, &globals, "update history-test.txt v3")

	// Call FileHistory to get commits where this file was updated
	historyArgs, cleanupHistoryArgs := types.NewLoreFileHistoryArgs(types.LoreFileHistoryArgs{
		Path:   testFilePath, // Use absolute path
		Length: 10,           // Get up to 10 revisions
	})
	defer cleanupHistoryArgs()

	historyEvents := []types.LoreFileHistoryEventData{}

	result, err := FileHistory(&globals, &historyArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_HISTORY {
				if historyEvent, ok := event.GetData().(*types.LoreFileHistoryEventDataFFI); ok {
					// Access the Address field inside the callback
					hashStr := historyEvent.Address.Hash.String()
					contextStr := historyEvent.Address.Context.String()
					t.Logf("Inside callback - Address: Hash=%s, Context=%s", hashStr, contextStr)

					// Clone the event to convert FFI data to Go data
					clonedEvent := historyEvent.Clone()
					historyEvents = append(historyEvents, clonedEvent)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileHistory failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileHistory returned non-zero result: %d", result)
	}

	// Verify we received history events
	if len(historyEvents) == 0 {
		t.Fatal("No FILE_HISTORY events received")
	}

	t.Logf("Received %d file history events", len(historyEvents))

	// Verify each event's Address field is accessible outside the callback
	for i, event := range historyEvents {
		// Access Address fields outside the callback
		hashStr := event.Address.Hash.String()
		contextStr := event.Address.Context.String()

		t.Logf("Outside callback [%d] - Address: Hash=%s, Context=%s", i, hashStr, contextStr)

		// Verify Address is not zero
		allZeroHash := true
		for _, b := range event.Address.Hash.Data {
			if b != 0 {
				allZeroHash = false
				break
			}
		}

		allZeroContext := true
		for _, b := range event.Address.Context.Data {
			if b != 0 {
				allZeroContext = false
				break
			}
		}

		if allZeroHash {
			t.Errorf("Event %d: Address.Hash is all zeros (unexpected)", i)
		}

		if allZeroContext {
			t.Errorf("Event %d: Address.Context is all zeros (unexpected)", i)
		}

		// Verify other cloned fields are accessible
		if event.Path != testFile {
			t.Errorf("Event %d: Expected Path='%s', got '%s'", i, testFile, event.Path)
		}

		t.Logf("Event %d: Path=%s, Revision=%s, RevisionNumber=%d, Size=%d",
			i, event.Path, event.Revision.String(), event.RevisionNumber, event.Size)

		// Verify Parent [2]LoreHash array is also accessible
		t.Logf("Event %d: Parent[0]=%s", i, event.Parent[0].String())
		t.Logf("Event %d: Parent[1]=%s", i, event.Parent[1].String())
	}

	// Verify we got events for all 3 commits
	if len(historyEvents) < 3 {
		t.Errorf("Expected at least 3 history events, got %d", len(historyEvents))
	}
}

func TestLoreFileInfo(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a directory
	testDir := "testdir"
	testDirPath := filepath.Join(globals.RepositoryPath.String(), testDir)
	if err := os.MkdirAll(testDirPath, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a file inside the directory
	testFile := "testfile.txt"
	testFilePath := filepath.Join(testDirPath, testFile)
	if err := os.WriteFile(testFilePath, []byte("test file content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	stageFiles(t, &globals, []string{testFilePath})

	// Commit the changes
	commitRevision(t, &globals, "add test directory and file")

	// Call FileInfo to get information about both the directory and file
	fileInfoArgs, cleanupFileInfoArgs := types.NewLoreFileInfoArgs(types.LoreFileInfoArgs{
		Paths: []string{testDirPath, testFilePath},
	})
	defer cleanupFileInfoArgs()

	fileInfoEvents := []types.LoreEvent{}

	result, err := FileInfo(&globals, &fileInfoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_INFO {
				// Clone the entire event
				clonedEvent := event.Clone()
				fileInfoEvents = append(fileInfoEvents, clonedEvent)
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileInfo failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileInfo returned non-zero result: %d", result)
	}

	// Verify we received file info events
	if len(fileInfoEvents) == 0 {
		t.Fatal("No FILE_INFO events received")
	}

	// Find and verify file info
	var fileInfo *types.LoreFileInfoEventData
	var dirInfo *types.LoreFileInfoEventData

	for _, event := range fileInfoEvents {
		if data, ok := event.Data.(types.LoreFileInfoEventData); ok {
			if strings.HasSuffix(data.Path, testFile) {
				fileInfo = &data
			}
			if data.Path == testDir {
				dirInfo = &data
			}
		}
	}

	// Verify file info
	if fileInfo != nil {
		if !fileInfo.IsFile {
			t.Errorf("Expected IsFile=true for file, got false")
		}
		if fileInfo.IsDir {
			t.Errorf("Expected IsDir=false for file, got true")
		}
		t.Logf("File info: Path=%s, IsFile=%v, IsDir=%v, Size=%d",
			fileInfo.Path, fileInfo.IsFile, fileInfo.IsDir, fileInfo.Size)
	} else {
		t.Error("No file info received for the test file")
	}

	// Verify directory info if present
	if dirInfo != nil {
		if dirInfo.IsFile {
			t.Errorf("Expected IsFile=false for directory, got true")
		}
		if !dirInfo.IsDir {
			t.Errorf("Expected IsDir=true for directory, got false")
		}
		t.Logf("Directory info: Path=%s, IsFile=%v, IsDir=%v",
			dirInfo.Path, dirInfo.IsFile, dirInfo.IsDir)
	}
}

func TestLoreFileWrite(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a file with version 1 content
	testFile := "versioned-file.txt"
	testFilePath := filepath.Join(globals.RepositoryPath.String(), testFile)
	version1Content := "Version 1 content"
	if err := os.WriteFile(testFilePath, []byte(version1Content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the file
	stageFiles(t, &globals, []string{testFilePath})

	// Commit the first version and capture the revision hash
	commitArgs, cleanupCommitArgs := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
		Message: "commit version 1",
	})
	defer cleanupCommitArgs()

	var firstRevisionHash string

	result, err := RevisionCommit(&globals, &commitArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_COMMIT_REVISION {
				if revisionData, ok := event.GetData().(*types.LoreRevisionCommitRevisionEventDataFFI); ok {
					firstRevisionHash = revisionData.Revision.String()
					t.Logf("First revision hash: %s", firstRevisionHash)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("First commit failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("First commit returned non-zero result: %d", result)
	}

	if firstRevisionHash == "" {
		t.Fatal("Failed to capture first revision hash from REVISION_COMMIT_REVISION event")
	}

	// Modify the file to version 2
	version2Content := "Version 2 content - updated"
	if err := os.WriteFile(testFilePath, []byte(version2Content), 0644); err != nil {
		t.Fatalf("Failed to update test file: %v", err)
	}

	// Stage and commit the second version
	stageFiles(t, &globals, []string{testFilePath})
	commitRevision(t, &globals, "commit version 2")

	// Verify the file now has version 2 content
	currentContent, err := os.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("Failed to read current file: %v", err)
	}
	if string(currentContent) != version2Content {
		t.Errorf("Expected current content '%s', got '%s'", version2Content, string(currentContent))
	}

	// Call FileWrite to write version 1 to a .old file
	outputPath := testFilePath + ".old"
	fileWriteArgs, cleanupFileWriteArgs := types.NewLoreFileWriteArgs(types.LoreFileWriteArgs{
		Path:     testFilePath, // Use absolute path
		Revision: firstRevisionHash,
		Output:   outputPath,
	})
	defer cleanupFileWriteArgs()

	result, err = FileWrite(&globals, &fileWriteArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_WRITE {
				if writeData, ok := event.GetData().(*types.LoreFileWriteEventDataFFI); ok {
					t.Logf("FILE_WRITE event: Path=%s", writeData.Path.String())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileWrite failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreFileWrite returned non-zero result: %d", result)
	}

	// Read the .old file and verify it contains version 1 content
	oldContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read .old file: %v", err)
	}

	if string(oldContent) != version1Content {
		t.Errorf("Expected .old file content '%s', got '%s'", version1Content, string(oldContent))
	}

	t.Logf("Successfully verified .old file contains version 1 content: %s", string(oldContent))
}

func TestLoreHashAndContextStringHexEncoding(t *testing.T) {
	// Test LoreHash.String() with various byte values, especially 0x00-0x0F
	t.Run("LoreHash", func(t *testing.T) {
		hash := types.LoreHash{
			Data: [32]uint8{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, // Test leading zeros (00-07)
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, // Test leading zeros (08-0F)
				0x10, 0x1F, 0x20, 0x2F, 0xA0, 0xAF, 0xF0, 0xFF, // Test non-zero leading digits
				0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF1, // Mix of values
			},
		}

		result := hash.String()

		// Expected string: each byte as 2 hex characters (lowercase)
		expected := "000102030405060708090a0b0c0d0e0f" + // First 16 bytes
			"101f202fa0aff0ff" + // Next 8 bytes
			"123456789abcdef1" // Last 8 bytes

		if result != expected {
			t.Errorf("LoreHash.String() failed\nExpected: %s\nGot:      %s", expected, result)
		}

		// Verify length is always 64 characters (32 bytes * 2)
		if len(result) != 64 {
			t.Errorf("Expected hash string length 64, got %d", len(result))
		}

		// Verify that bytes 0x00-0x0F are printed with leading zero
		for i := 0; i < 16; i++ {
			expectedByte := fmt.Sprintf("%02x", i)
			actualByte := result[i*2 : i*2+2]
			if actualByte != expectedByte {
				t.Errorf("Byte %d (value 0x%02X): expected '%s', got '%s'", i, i, expectedByte, actualByte)
			}
		}
	})

	// Test LoreContext.String() with various byte values, especially 0x00-0x0F
	t.Run("LoreContext", func(t *testing.T) {
		context := types.LoreContext{
			Data: [16]uint8{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, // Test leading zeros (00-07)
				0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, // Test leading zeros (08-0F)
			},
		}

		result := context.String()

		// Expected string: each byte as 2 hex characters (lowercase)
		expected := "000102030405060708090a0b0c0d0e0f"

		if result != expected {
			t.Errorf("LoreContext.String() failed\nExpected: %s\nGot:      %s", expected, result)
		}

		// Verify length is always 32 characters (16 bytes * 2)
		if len(result) != 32 {
			t.Errorf("Expected context string length 32, got %d", len(result))
		}

		// Verify that bytes 0x00-0x0F are printed with leading zero
		for i := 0; i < 16; i++ {
			expectedByte := fmt.Sprintf("%02x", i)
			actualByte := result[i*2 : i*2+2]
			if actualByte != expectedByte {
				t.Errorf("Byte %d (value 0x%02X): expected '%s', got '%s'", i, i, expectedByte, actualByte)
			}
		}
	})

	// Test edge cases
	t.Run("EdgeCases", func(t *testing.T) {
		// All zeros
		allZeroHash := types.LoreHash{Data: [32]uint8{}}
		zeroHashStr := allZeroHash.String()
		expectedZero := strings.Repeat("0", 64)
		if zeroHashStr != expectedZero {
			t.Errorf("All-zero hash: expected %s, got %s", expectedZero, zeroHashStr)
		}

		// All 0xFF
		allFFHash := types.LoreHash{Data: [32]uint8{
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		}}
		allFFHashStr := allFFHash.String()
		expectedFF := strings.Repeat("ff", 32)
		if allFFHashStr != expectedFF {
			t.Errorf("All-0xFF hash: expected %s, got %s", expectedFF, allFFHashStr)
		}

		// All zeros context
		allZeroContext := types.LoreContext{Data: [16]uint8{}}
		zeroContextStr := allZeroContext.String()
		expectedZeroContext := strings.Repeat("0", 32)
		if zeroContextStr != expectedZeroContext {
			t.Errorf("All-zero context: expected %s, got %s", expectedZeroContext, zeroContextStr)
		}

		// All 0xFF context
		allFFContext := types.LoreContext{Data: [16]uint8{
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		}}
		allFFContextStr := allFFContext.String()
		expectedFFContext := strings.Repeat("ff", 16)
		if allFFContextStr != expectedFFContext {
			t.Errorf("All-0xFF context: expected %s, got %s", expectedFFContext, allFFContextStr)
		}
	})
}

func TestLoreRevisionMetadataSetAndList(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a metadata file for binary metadata
	metadataFile := "metadata-content.bin"
	metadataFilePath := filepath.Join(globals.RepositoryPath.String(), metadataFile)
	binaryContent := []byte("Binary metadata content from file")
	if err := os.WriteFile(metadataFilePath, binaryContent, 0644); err != nil {
		t.Fatalf("Failed to create metadata file: %v", err)
	}

	// Create a regular file to commit
	testFile := "test-with-metadata.txt"
	testFilePath := filepath.Join(globals.RepositoryPath.String(), testFile)
	if err := os.WriteFile(testFilePath, []byte("Test file content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Stage the test file (not the metadata file)
	stageFiles(t, &globals, []string{testFilePath})

	// Set metadata before commit
	metadataSetArgs, cleanupMetadataSetArgs := types.NewLoreRevisionMetadataSetArgs(types.LoreRevisionMetadataSetArgs{
		Keys:    []string{"meta-string", "meta-file", "empty-string"},
		Values:  []string{"string value", metadataFilePath, ""},
		Formats: []types.LoreMetadataType{types.LoreMetadataType_STRING, types.LoreMetadataType_BINARY, types.LoreMetadataType_STRING},
	})
	defer cleanupMetadataSetArgs()

	result, err := RevisionMetadataSet(&globals, &metadataSetArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionMetadataSet failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionMetadataSet returned non-zero result: %d", result)
	}

	// Commit the revision and capture the revision hash
	commitArgs, cleanupCommitArgs := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
		Message: "commit with metadata",
	})
	defer cleanupCommitArgs()

	var revisionHash string

	result, err = RevisionCommit(&globals, &commitArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_COMMIT_REVISION {
				if revisionData, ok := event.GetData().(*types.LoreRevisionCommitRevisionEventDataFFI); ok {
					revisionHash = revisionData.Revision.String()
					t.Logf("Committed revision hash: %s", revisionHash)
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("Commit failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("Commit returned non-zero result: %d", result)
	}

	if revisionHash == "" {
		t.Fatal("Failed to capture revision hash from commit")
	}

	// List metadata for the committed revision
	metadataListArgs, cleanupMetadataListArgs := types.NewLoreRevisionMetadataListArgs(types.LoreRevisionMetadataListArgs{
		Revision: revisionHash,
	})
	defer cleanupMetadataListArgs()

	metadataEvents := make(map[string]types.LoreMetadataEventData)

	result, err = RevisionMetadataList(&globals, &metadataListArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_METADATA {
				if metadataEvent, ok := event.GetData().(*types.LoreMetadataEventDataFFI); ok {
					key := metadataEvent.Key.String()
					// Clone the event to store it
					clonedEvent := event.Clone()
					if clonedData, ok := clonedEvent.Data.(types.LoreMetadataEventData); ok {
						metadataEvents[key] = clonedData
					}
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionMetadataList failed: %v", err)
	}

	if result != 0 {
		t.Fatalf("LoreRevisionMetadataList returned non-zero result: %d", result)
	}

	// Verify we received metadata events
	if len(metadataEvents) == 0 {
		t.Fatal("No METADATA events received")
	}

	t.Logf("Received %d metadata entries", len(metadataEvents))

	// Verify implicit metadata that's always present

	// Verify "branch" metadata (CONTEXT type)
	if branch, ok := metadataEvents["branch"]; ok {
		if branch.Value.Tag != types.LoreMetadataTag_CONTEXT {
			t.Errorf("Expected branch to have tag CONTEXT, got %d", branch.Value.Tag)
		} else {
			ctx := branch.Value.Context
			t.Logf("branch: key=%s, tag=CONTEXT, context=%s", branch.Key, ctx.String())
		}
	} else {
		t.Error("branch metadata not found (expected implicit metadata)")
	}

	// Verify "timestamp" metadata (NUMERIC type)
	if timestamp, ok := metadataEvents["timestamp"]; ok {
		if timestamp.Value.Tag != types.LoreMetadataTag_NUMERIC {
			t.Errorf("Expected timestamp to have tag NUMERIC, got %d", timestamp.Value.Tag)
		} else {
			ts := *timestamp.Value.Numeric
			t.Logf("timestamp: key=%s, tag=NUMERIC, value=%d", timestamp.Key, ts)
			// Verify it's a reasonable Unix timestamp (should be > 0 and not too far in the future)
			if ts == 0 {
				t.Error("timestamp value is 0 (unexpected)")
			}
		}
	} else {
		t.Error("timestamp metadata not found (expected implicit metadata)")
	}

	// Verify "message" metadata (STRING type, should contain commit message)
	if message, ok := metadataEvents["message"]; ok {
		if message.Value.Tag != types.LoreMetadataTag_STRING {
			t.Errorf("Expected message to have tag STRING, got %d", message.Value.Tag)
		} else {
			msg := *message.Value.String
			if msg != "commit with metadata" {
				t.Errorf("Expected message value 'commit with metadata', got '%s'", msg)
			}
			t.Logf("message: key=%s, tag=STRING, value=%s", message.Key, msg)
		}
	} else {
		t.Error("message metadata not found (expected implicit metadata)")
	}

	// Verify custom metadata

	// Verify meta-string metadata
	if metaString, ok := metadataEvents["meta-string"]; ok {
		if metaString.Value.Tag != types.LoreMetadataTag_STRING {
			t.Errorf("Expected meta-string to have tag STRING, got %d", metaString.Value.Tag)
		} else {
			value := *metaString.Value.String
			if value != "string value" {
				t.Errorf("Expected meta-string value 'string value', got '%s'", value)
			}
			t.Logf("meta-string: key=%s, tag=STRING, value=%s", metaString.Key, value)
		}
	} else {
		t.Error("meta-string metadata not found")
	}

	// Verify meta-file metadata (should be ADDRESS type for binary data)
	if metaFile, ok := metadataEvents["meta-file"]; ok {
		if metaFile.Value.Tag != types.LoreMetadataTag_ADDRESS {
			t.Errorf("Expected meta-file to have tag ADDRESS, got %d", metaFile.Value.Tag)
		} else {
			address := metaFile.Value.Address
			hashStr := address.Hash.String()
			contextStr := address.Context.String()
			t.Logf("meta-file: key=%s, tag=ADDRESS, hash=%s, context=%s", metaFile.Key, hashStr, contextStr)

			// Verify address is not all zeros
			allZero := true
			for _, b := range address.Hash.Data {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				t.Error("meta-file address hash is all zeros (unexpected)")
			}
		}
	} else {
		t.Error("meta-file metadata not found")
	}

	// Verify empty-string metadata
	if emptyString, ok := metadataEvents["empty-string"]; ok {
		if emptyString.Value.Tag != types.LoreMetadataTag_STRING {
			t.Errorf("Expected empty-string to have tag STRING, got %d", emptyString.Value.Tag)
		} else {
			value := *emptyString.Value.String
			if value != "" {
				t.Errorf("Expected empty-string value '', got '%s'", value)
			}
			t.Logf("empty-string: key=%s, tag=STRING, value='' (empty)", emptyString.Key)
		}
	} else {
		t.Error("empty-string metadata not found")
	}

	t.Logf("Successfully verified all metadata entries")
}

// createBranch creates a new branch with the given name
func createBranch(t *testing.T, globals *types.LoreGlobalArgsFFI, name string) {
	t.Helper()

	createArgs, cleanupCreateArgs := types.NewLoreBranchCreateArgs(types.LoreBranchCreateArgs{
		Branch: name,
	})
	defer cleanupCreateArgs()

	result, err := BranchCreate(globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil || result != 0 {
		t.Fatalf("Failed to create branch %s: err=%v, result=%d", name, err, result)
	}
}

// switchBranch switches to the given branch
func switchBranch(t *testing.T, globals *types.LoreGlobalArgsFFI, name string) {
	t.Helper()

	switchArgs, cleanupSwitchArgs := types.NewLoreBranchSwitchArgs(types.LoreBranchSwitchArgs{
		Branch: name,
	})
	defer cleanupSwitchArgs()

	result, err := BranchSwitch(globals, &switchArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})

	if err != nil || result != 0 {
		t.Fatalf("Failed to switch to branch %s: err=%v, result=%d", name, err, result)
	}
}

func TestLoreBranchInfo(t *testing.T) {
	globals := setupTestRepository(t)

	branchInfoEvents := []types.LoreBranchInfoEventData{}

	args, cleanupArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{
		Branch: "main",
	})
	defer cleanupArgs()

	result, err := BranchInfo(&globals, &args, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_INFO {
				if data, ok := event.GetData().(*types.LoreBranchInfoEventDataFFI); ok {
					branchInfoEvents = append(branchInfoEvents, data.Clone())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchInfo failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreBranchInfo returned non-zero result: %d", result)
	}

	if len(branchInfoEvents) != 1 {
		t.Fatalf("Expected 1 branch info event, got %d", len(branchInfoEvents))
	}

	if branchInfoEvents[0].Name != "main" {
		t.Errorf("Expected branch name 'main', got %q", branchInfoEvents[0].Name)
	}

	t.Logf("Branch info: Name=%q, Creator=%q", branchInfoEvents[0].Name, branchInfoEvents[0].Creator)
}

func TestLoreBranchSwitch(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a feature branch (we're on main after setupTestRepository)
	createBranch(t, &globals, "feature-branch")

	// Create a file on the feature branch and commit
	featureFile := createFileWithContents(t, globals.RepositoryPath.String(), "feature-only.txt", "feature content")
	stageFiles(t, &globals, []string{featureFile})
	commitRevision(t, &globals, "feature commit")

	// Verify file exists
	featureFilePath := filepath.Join(globals.RepositoryPath.String(), "feature-only.txt")
	if _, err := os.Stat(featureFilePath); os.IsNotExist(err) {
		t.Fatal("feature-only.txt should exist on feature branch")
	}

	// Switch back to main
	switchBranch(t, &globals, "main")

	// File should no longer exist on main
	if _, err := os.Stat(featureFilePath); !os.IsNotExist(err) {
		t.Error("feature-only.txt should not exist after switching to main")
	}
}

func TestLoreBranchDiff(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a feature branch with a new file
	createBranch(t, &globals, "diff-branch")
	diffFile := createFileWithContents(t, globals.RepositoryPath.String(), "diff-file.txt", "diff content")
	stageFiles(t, &globals, []string{diffFile})
	commitRevision(t, &globals, "diff commit")

	args, cleanupArgs := types.NewLoreBranchDiffArgs(types.LoreBranchDiffArgs{
		Target: "main",
	})
	defer cleanupArgs()

	diffChanges := []types.LoreBranchDiffChangeEventData{}

	result, err := BranchDiff(&globals, &args, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_BRANCH_DIFF_CHANGE {
				if data, ok := event.GetData().(*types.LoreBranchDiffChangeEventDataFFI); ok {
					diffChanges = append(diffChanges, data.Clone())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreBranchDiff failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreBranchDiff returned non-zero result: %d", result)
	}

	// Should have at least one change (the new file)
	if len(diffChanges) == 0 {
		t.Error("Expected at least one diff change event")
	}

	t.Logf("Received %d diff change events", len(diffChanges))
}

func TestLoreRevisionFindByMetadata(t *testing.T) {
	globals := setupTestRepository(t)

	// Set metadata on staged changes
	setArgs, cleanupSetArgs := types.NewLoreRevisionMetadataSetArgs(types.LoreRevisionMetadataSetArgs{
		Keys:    []string{"search-key"},
		Values:  []string{"search-value"},
		Formats: []types.LoreMetadataType{types.LoreMetadataType_STRING},
	})
	defer cleanupSetArgs()

	result, err := RevisionMetadataSet(&globals, &setArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreRevisionMetadataSet failed: err=%v, result=%d", err, result)
	}

	// Create a file and commit so the metadata is attached to a revision
	testFile := createFileWithContents(t, globals.RepositoryPath.String(), "find-test.txt", "find content")
	stageFiles(t, &globals, []string{testFile})
	commitRevision(t, &globals, "commit with metadata")

	// Find revision by metadata
	findResults := []types.LoreRevisionFindEventData{}

	findArgs, cleanupFindArgs := types.NewLoreRevisionFindArgs(types.LoreRevisionFindArgs{
		Key:   "search-key",
		Value: "search-value",
	})
	defer cleanupFindArgs()

	result, err = RevisionFind(&globals, &findArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_FIND {
				if data, ok := event.GetData().(*types.LoreRevisionFindEventDataFFI); ok {
					findResults = append(findResults, data.Clone())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionFind failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreRevisionFind returned non-zero result: %d", result)
	}

	if len(findResults) != 1 {
		t.Fatalf("Expected 1 find result, got %d", len(findResults))
	}

	// Verify the signature hash is 32 bytes (64 hex chars)
	sigStr := findResults[0].Signature.String()
	if len(sigStr) != 64 {
		t.Errorf("Expected 64-char signature hash, got %d chars: %s", len(sigStr), sigStr)
	}

	t.Logf("Found revision with signature: %s", sigStr)
}

func TestLoreRevisionFindByNumber(t *testing.T) {
	globals := setupTestRepository(t)

	// setupTestRepository already creates revision 1 (initial commit)
	// Create another file and commit for revision 2
	testFile := createFileWithContents(t, globals.RepositoryPath.String(), "find-by-number.txt", "number content")
	stageFiles(t, &globals, []string{testFile})
	commitRevision(t, &globals, "second commit")

	findResults := []types.LoreRevisionFindEventData{}

	findArgs, cleanupFindArgs := types.NewLoreRevisionFindArgs(types.LoreRevisionFindArgs{
		Number: 1,
	})
	defer cleanupFindArgs()

	result, err := RevisionFind(&globals, &findArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_REVISION_FIND {
				if data, ok := event.GetData().(*types.LoreRevisionFindEventDataFFI); ok {
					findResults = append(findResults, data.Clone())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreRevisionFind failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreRevisionFind returned non-zero result: %d", result)
	}

	if len(findResults) != 1 {
		t.Fatalf("Expected 1 find result, got %d", len(findResults))
	}

	sigStr := findResults[0].Signature.String()
	if len(sigStr) != 64 {
		t.Errorf("Expected 64-char signature hash, got %d chars: %s", len(sigStr), sigStr)
	}

	t.Logf("Found revision #1 with signature: %s", sigStr)
}

func TestLoreFileMetadataSetAndList(t *testing.T) {
	globals := setupTestRepository(t)

	// Create and stage a file
	testFile := createFileWithContents(t, globals.RepositoryPath.String(), "metadata-file.txt", "metadata content")
	stageFiles(t, &globals, []string{testFile})

	// Set metadata on the file
	setArgs, cleanupSetArgs := types.NewLoreFileMetadataSetArgs(types.LoreFileMetadataSetArgs{
		Paths:   []string{testFile},
		Keys:    []string{"test-key"},
		Values:  []string{"test-value"},
		Formats: []types.LoreMetadataType{types.LoreMetadataType_STRING},
		Entries: []uint32{1},
	})
	defer cleanupSetArgs()

	result, err := FileMetadataSet(&globals, &setArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, userContext uint64) {},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreFileMetadataSet failed: err=%v, result=%d", err, result)
	}

	// List metadata for the file
	metadataEvents := []types.LoreMetadataEventData{}

	listArgs, cleanupListArgs := types.NewLoreFileMetadataListArgs(types.LoreFileMetadataListArgs{
		Path: testFile,
	})
	defer cleanupListArgs()

	result, err = FileMetadataList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_METADATA {
				if data, ok := event.GetData().(*types.LoreMetadataEventDataFFI); ok {
					metadataEvents = append(metadataEvents, data.Clone())
				}
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileMetadataList failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreFileMetadataList returned non-zero result: %d", result)
	}

	// Verify we got the metadata we set
	found := false
	for _, m := range metadataEvents {
		if m.Key == "test-key" && m.Value.Tag == types.LoreMetadataTag_STRING && m.Value.String != nil && *m.Value.String == "test-value" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find metadata with key='test-key' and value='test-value'")
	}

	t.Logf("Received %d metadata events", len(metadataEvents))
}

func TestLoreSharedStoreCreateAndInfo(t *testing.T) {
	if libErr != nil {
		t.Skipf("Lore library not loaded: %v", libErr)
	}

	tempDir, err := os.MkdirTemp("", "lore-global-store-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	globalStorePath := filepath.Join(tempDir, "store")
	remoteUrl := "lore-go-unit-test-global-store"

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		Offline: true,
	})
	defer cleanupGlobals()

	// Create global store
	createArgs, cleanupCreateArgs := types.NewLoreSharedStoreCreateArgs(types.LoreSharedStoreCreateArgs{
		RemoteUrl:   remoteUrl,
		Path:        globalStorePath,
		MakeDefault: true,
	})
	defer cleanupCreateArgs()

	var storeCreateEvents []types.LoreSharedStoreCreateEventData

	result, err := SharedStoreCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_SHARED_STORE_CREATE {
				storeCreateEvents = append(storeCreateEvents, event.Clone().Data.(types.LoreSharedStoreCreateEventData))
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreSharedStoreCreate failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreSharedStoreCreate returned non-zero result: %d", result)
	}

	if len(storeCreateEvents) == 0 {
		t.Fatal("Expected at least one GLOBAL_STORE_CREATE event")
	}
	if !strings.Contains(storeCreateEvents[0].Path, globalStorePath) {
		t.Errorf("Expected store create path to contain %q, got %q", globalStorePath, storeCreateEvents[0].Path)
	}

	// Query global store info
	infoArgs, cleanupInfoArgs := types.NewLoreSharedStoreInfoArgs(types.LoreSharedStoreInfoArgs{})
	defer cleanupInfoArgs()

	var globalStoreInfoEvents []types.LoreSharedStoreInfoEventData

	result, err = SharedStoreInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_SHARED_STORE_INFO {
				globalStoreInfoEvents = append(globalStoreInfoEvents, event.Clone().Data.(types.LoreSharedStoreInfoEventData))
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreSharedStoreInfo failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreSharedStoreInfo returned non-zero result: %d", result)
	}

	if len(globalStoreInfoEvents) != 1 {
		t.Fatalf("Expected exactly 1 GLOBAL_STORE_INFO event, got %d", len(globalStoreInfoEvents))
	}

	info := globalStoreInfoEvents[0]

	// Find the store by remote URL
	storeIndex := -1
	for i, url := range info.RemoteUrls {
		if url == remoteUrl {
			storeIndex = i
			break
		}
	}

	if storeIndex < 0 {
		t.Fatalf("Expected to find remote URL %q in global store info, got %v", remoteUrl, info.RemoteUrls)
	}

	// Verify the Exists []bool field (tests LoreUint8Array -> []bool mapping)
	if storeIndex >= len(info.Exists) {
		t.Fatalf("Exists array too short: index %d, len %d", storeIndex, len(info.Exists))
	}
	if !info.Exists[storeIndex] {
		t.Error("Expected store to exist")
	}

	if !strings.Contains(info.Paths[storeIndex], globalStorePath) {
		t.Errorf("Expected path to contain %q, got %q", globalStorePath, info.Paths[storeIndex])
	}
	if info.RemoteUrls[storeIndex] != remoteUrl {
		t.Errorf("Expected remote URL %q, got %q", remoteUrl, info.RemoteUrls[storeIndex])
	}
}

func TestLoreFileStage(t *testing.T) {
	globals := setupTestRepository(t)

	// Create a new file
	testFile := createFileWithContents(t, globals.RepositoryPath.String(), "stage-test.txt", "stage content")

	// Stage the file and verify via events
	stageArgs, cleanupStageArgs := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
		Paths: []string{testFile},
	})
	defer cleanupStageArgs()

	stageEndReceived := false

	result, err := FileStage(&globals, &stageArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_FILE_STAGE_END {
				stageEndReceived = true
			}
		},
		UserContext: 0,
	})

	if err != nil {
		t.Fatalf("LoreFileStage failed: %v", err)
	}
	if result != 0 {
		t.Fatalf("LoreFileStage returned non-zero result: %d", result)
	}

	if !stageEndReceived {
		t.Error("Expected to receive FILE_STAGE_END event")
	}
}

// sortedCopy returns a sorted copy of the slice (does not mutate the input).
func sortedCopy(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}

// TestLoreFileDependencyAddListRemove exercises the parallel-array contract
// shared by LoreFileDependencyAdd / List / Remove. It is the only operation
// that drives both LoreUint32ArrayFFI (DepCounts/TagCounts) on input and
// LoreStringArrayFFI (Tags) on per-entry output, so it is the densest
// regression test for the autogenerated _array_t marshaling.
func TestLoreFileDependencyAddListRemove(t *testing.T) {
	globals := setupTestRepository(t)
	repoDir := globals.RepositoryPath.String()

	// Five files: a, b are sources; x, y, z are targets.
	aPath := createFileWithContents(t, repoDir, "a.txt", "a")
	bPath := createFileWithContents(t, repoDir, "b.txt", "b")
	xPath := createFileWithContents(t, repoDir, "x.txt", "x")
	yPath := createFileWithContents(t, repoDir, "y.txt", "y")
	zPath := createFileWithContents(t, repoDir, "z.txt", "z")
	stageFiles(t, &globals, []string{aPath, bPath, xPath, yPath, zPath})
	commitRevision(t, &globals, "files for dependency test")

	// Dependency layout:
	//   a.txt -> x.txt (tags: ["alpha"])
	//   a.txt -> y.txt (tags: ["alpha", "beta"])
	//   b.txt -> z.txt (tags: [])
	//
	// Parallel arrays:
	//   Paths        = [a, b]                               len = 2
	//   DepCounts    = [2, 1]                               len = 2 (matches Paths)
	//   Dependencies = [x, y, z]                            len = sum(DepCounts) = 3
	//   TagCounts    = [1, 2, 0]                            len = 3 (matches Dependencies)
	//   Tags         = ["alpha", "alpha", "beta"]           len = sum(TagCounts) = 3
	addArgs, cleanupAdd := types.NewLoreFileDependencyAddArgs(types.LoreFileDependencyAddArgs{
		Paths:        []string{aPath, bPath},
		DepCounts:    []uint32{2, 1},
		Dependencies: []string{xPath, yPath, zPath},
		TagCounts:    []uint32{1, 2, 0},
		Tags:         []string{"alpha", "alpha", "beta"},
	})
	defer cleanupAdd()

	type entry struct {
		Path       string
		Dependency string
		Tags       []string
	}
	var addEntries []entry

	result, err := FileDependencyAdd(&globals, &addArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag == types.LoreEventTag_FILE_DEPENDENCY_ADD_ENTRY {
				if d, ok := event.GetData().(*types.LoreFileDependencyAddEntryEventDataFFI); ok {
					// Touch the FFI array inside the callback to exercise Len/Get
					// (Clone() also covers it, but Len/Get is the no-allocation path).
					_ = d.Tags.Len()
					cloned := d.Clone()
					addEntries = append(addEntries, entry{
						Path:       cloned.Path,
						Dependency: cloned.Dependency,
						Tags:       cloned.Tags,
					})
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreFileDependencyAdd: err=%v result=%d", err, result)
	}
	if len(addEntries) != 3 {
		t.Fatalf("expected 3 ADD_ENTRY events, got %d", len(addEntries))
	}

	// Index by (Path, Dependency) for order-independent verification.
	addByKey := make(map[string][]string, len(addEntries))
	for _, e := range addEntries {
		addByKey[e.Path+"|"+e.Dependency] = e.Tags
	}
	wantAdd := map[string][]string{
		aPath + "|" + xPath: {"alpha"},
		aPath + "|" + yPath: {"alpha", "beta"},
		bPath + "|" + zPath: nil, // empty tag array (LoreStringArrayFFI with Count=0)
	}
	for key, want := range wantAdd {
		got, ok := addByKey[key]
		if !ok {
			t.Errorf("ADD_ENTRY missing for %s", key)
			continue
		}
		if !reflect.DeepEqual(sortedCopy(got), sortedCopy(want)) {
			t.Errorf("ADD_ENTRY %s tags: want %v got %v", key, want, got)
		}
	}

	// List dependencies starting from a.txt and b.txt.
	listArgs, cleanupList := types.NewLoreFileDependencyListArgs(types.LoreFileDependencyListArgs{
		Paths: []string{aPath, bPath},
	})
	defer cleanupList()

	listEntries := make(map[string][]string) // dependency path -> tags
	var listFileCount, listEntryCount uint64

	result, err = FileDependencyList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			switch event.Tag {
			case types.LoreEventTag_FILE_DEPENDENCY_LIST_BEGIN:
				if d, ok := event.GetData().(*types.LoreFileDependencyListBeginEventDataFFI); ok {
					listFileCount = d.FileCount
				}
			case types.LoreEventTag_FILE_DEPENDENCY_LIST_ENTRY:
				if d, ok := event.GetData().(*types.LoreFileDependencyListEntryEventDataFFI); ok {
					_ = d.Tags.Len()
					cloned := d.Clone()
					listEntries[cloned.Path] = cloned.Tags
				}
			case types.LoreEventTag_FILE_DEPENDENCY_LIST_END:
				if d, ok := event.GetData().(*types.LoreFileDependencyListEndEventDataFFI); ok {
					listEntryCount = d.TotalEntryCount
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreFileDependencyList: err=%v result=%d", err, result)
	}
	if listFileCount != 2 {
		t.Errorf("LIST_BEGIN.FileCount: want 2 got %d", listFileCount)
	}
	if listEntryCount != 3 {
		t.Errorf("LIST_END.TotalEntryCount: want 3 got %d", listEntryCount)
	}

	// LIST_ENTRY emits paths relative to the repo root, not the absolute
	// paths that were passed in. Compare by basename.
	wantList := map[string][]string{
		"x.txt": {"alpha"},
		"y.txt": {"alpha", "beta"},
		"z.txt": nil,
	}
	for path, want := range wantList {
		got, ok := listEntries[path]
		if !ok {
			t.Errorf("LIST_ENTRY missing for %s", path)
			continue
		}
		if !reflect.DeepEqual(sortedCopy(got), sortedCopy(want)) {
			t.Errorf("LIST_ENTRY %s tags: want %v got %v", path, want, got)
		}
	}

	// Remove a.txt -> y.txt only. This exercises a degenerate parallel-array
	// shape (Paths=1, DepCounts=[1], Dependencies=[y], TagCounts=[0], Tags=[]).
	removeArgs, cleanupRemove := types.NewLoreFileDependencyRemoveArgs(types.LoreFileDependencyRemoveArgs{
		Paths:        []string{aPath},
		DepCounts:    []uint32{1},
		Dependencies: []string{yPath},
		TagCounts:    []uint32{0},
		Tags:         []string{},
	})
	defer cleanupRemove()

	var removeEntries []entry
	result, err = FileDependencyRemove(&globals, &removeArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag == types.LoreEventTag_FILE_DEPENDENCY_REMOVE_ENTRY {
				if d, ok := event.GetData().(*types.LoreFileDependencyRemoveEntryEventDataFFI); ok {
					_ = d.Tags.Len()
					cloned := d.Clone()
					removeEntries = append(removeEntries, entry{
						Path:       cloned.Path,
						Dependency: cloned.Dependency,
						Tags:       cloned.Tags,
					})
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreFileDependencyRemove: err=%v result=%d", err, result)
	}
	if len(removeEntries) != 1 {
		t.Fatalf("expected 1 REMOVE_ENTRY event, got %d", len(removeEntries))
	}
	if removeEntries[0].Path != aPath || removeEntries[0].Dependency != yPath {
		t.Errorf("REMOVE_ENTRY: want a→y got %s→%s", removeEntries[0].Path, removeEntries[0].Dependency)
	}

	// Re-list and verify y.txt is gone, x.txt and z.txt remain.
	postEntries := make(map[string][]string)
	result, err = FileDependencyList(&globals, &listArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag == types.LoreEventTag_FILE_DEPENDENCY_LIST_ENTRY {
				if d, ok := event.GetData().(*types.LoreFileDependencyListEntryEventDataFFI); ok {
					cloned := d.Clone()
					postEntries[cloned.Path] = cloned.Tags
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreFileDependencyList (post-remove): err=%v result=%d", err, result)
	}
	if _, gone := postEntries["y.txt"]; gone {
		t.Errorf("y.txt should have been removed, but still appears in list")
	}
	if _, ok := postEntries["x.txt"]; !ok {
		t.Errorf("x.txt missing after remove")
	}
	if _, ok := postEntries["z.txt"]; !ok {
		t.Errorf("z.txt missing after remove")
	}
}

// TestLoreStoragePutGet round-trips multi-item content through an in-memory
// store. It is the only test that drives LoreStoragePutItemArrayFFI and
// LoreStorageGetItemArrayFFI (struct-of-array shapes). The payload set
// includes ASCII, a 2KB block, and multibyte UTF-8 to verify byte-exact
// transfer for non-trivial sizes. Modeled on the JS SDK test
// "should open in-memory storage, put and get data, then close".
func TestLoreStoragePutGet(t *testing.T) {
	if libErr != nil {
		t.Skipf("Lore library not loaded: %v", libErr)
	}

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		Offline: true,
	})
	defer cleanupGlobals()

	openArgs, cleanupOpen := types.NewLoreStorageOpenArgs(types.LoreStorageOpenArgs{
		InMemory: true,
	})
	defer cleanupOpen()

	var handleId uint64
	var handleSeen bool
	result, err := StorageOpen(&globals, &openArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag == types.LoreEventTag_STORAGE_OPENED {
				if d, ok := event.GetData().(*types.LoreStorageOpenedEventDataFFI); ok {
					handleId = d.HandleId
					handleSeen = true
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreStorageOpen: err=%v result=%d", err, result)
	}
	if !handleSeen {
		t.Fatal("STORAGE_OPENED event not received")
	}
	handle := types.LoreStore{HandleId: handleId}

	// Fixed 16-byte UUIDs (deterministic).
	var partition types.LorePartition
	var ctx types.LoreContext
	copy(partition.Data[:], "1234567890123456")
	copy(ctx.Data[:], "1234567890123456")

	// Three items: short ASCII, 2KB block, multibyte UTF-8.
	payloads := []string{
		"hello",
		strings.Repeat("ab", 1024),
		"with multibyte unicode chars -öäÄÅ𒂔𒀱的ЛЛЛµ𒅌𓉡𓉢‼️🌏🇩🇪",
	}
	items := make([]types.LoreStoragePutItem, len(payloads))
	for i, p := range payloads {
		items[i] = types.LoreStoragePutItem{
			Id:        uint64(i + 1),
			Partition: partition,
			Context:   ctx,
			Data:      []byte(p),
		}
	}

	putArgs, cleanupPut := types.NewLoreStoragePutArgs(types.LoreStoragePutArgs{
		Handle: handle,
		Items:  items,
	})
	defer cleanupPut()

	putAddresses := make(map[uint64]types.LoreAddress)
	result, err = StoragePut(&globals, &putArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag == types.LoreEventTag_STORAGE_PUT_ITEM_COMPLETE {
				if d, ok := event.GetData().(*types.LoreStoragePutItemCompleteEventDataFFI); ok {
					if d.ErrorCode != 0 {
						t.Errorf("PUT_ITEM_COMPLETE id=%d ErrorCode=%d", d.Id, d.ErrorCode)
						return
					}
					putAddresses[d.Id] = d.Address.Clone()
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreStoragePut: err=%v result=%d", err, result)
	}
	if len(putAddresses) != len(items) {
		t.Fatalf("expected %d PUT_ITEM_COMPLETE events, got %d", len(items), len(putAddresses))
	}

	// Build GET items pointing at the addresses we just wrote.
	getItems := make([]types.LoreStorageGetItem, 0, len(items))
	for _, it := range items {
		addr, ok := putAddresses[it.Id]
		if !ok {
			t.Fatalf("no put address for id %d", it.Id)
		}
		getItems = append(getItems, types.LoreStorageGetItem{
			Id:        it.Id,
			Partition: it.Partition,
			Address:   addr,
			Streaming: false,
		})
	}

	getArgs, cleanupGet := types.NewLoreStorageGetArgs(types.LoreStorageGetArgs{
		Handle: handle,
		Items:  getItems,
	})
	defer cleanupGet()

	headerSize := make(map[uint64]uint64)
	fetched := make(map[uint64][]byte)
	completeCount := 0

	result, err = StorageGet(&globals, &getArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			switch event.Tag {
			case types.LoreEventTag_STORAGE_GET_HEADER:
				if d, ok := event.GetData().(*types.LoreStorageGetHeaderEventDataFFI); ok {
					headerSize[d.Id] = d.SizeContent
				}
			case types.LoreEventTag_STORAGE_GET_DATA:
				if d, ok := event.GetData().(*types.LoreStorageGetDataEventDataFFI); ok {
					// d.Bytes is a LoreBytesFFI (pointer + length into C-owned
					// memory). Clone() copies it into a Go []byte that remains
					// valid after the callback returns.
					fetched[d.Id] = append(fetched[d.Id], d.Bytes.Clone()...)
				}
			case types.LoreEventTag_STORAGE_GET_ITEM_COMPLETE:
				if d, ok := event.GetData().(*types.LoreStorageGetItemCompleteEventDataFFI); ok {
					if d.ErrorCode != 0 {
						t.Errorf("GET_ITEM_COMPLETE id=%d ErrorCode=%d", d.Id, d.ErrorCode)
					}
					completeCount++
				}
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreStorageGet: err=%v result=%d", err, result)
	}
	if completeCount != len(items) {
		t.Errorf("expected %d GET_ITEM_COMPLETE events, got %d", len(items), completeCount)
	}
	for _, it := range items {
		want := []byte(payloads[it.Id-1])
		if got := headerSize[it.Id]; got != uint64(len(want)) {
			t.Errorf("id=%d: GET_HEADER.SizeContent want %d got %d", it.Id, len(want), got)
		}
		if got := fetched[it.Id]; !bytes.Equal(got, want) {
			t.Errorf("id=%d: payload mismatch: want %d bytes, got %d bytes", it.Id, len(want), len(got))
		}
	}

	closeArgs, cleanupClose := types.NewLoreStorageCloseArgs(types.LoreStorageCloseArgs{
		Handle: handle,
	})
	defer cleanupClose()
	result, err = StorageClose(&globals, &closeArgs, &types.LoreEventCallbackConfig{
		Callback:    func(event *types.LoreEventFFI, _ uint64) {},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreStorageClose: err=%v result=%d", err, result)
	}
}

// TestLoreSharedStoreInfoExistsArray verifies that LoreUint8ArrayFFI parses
// correctly when the C library populates the Exists field with multiple
// entries. The existing TestLoreSharedStoreCreateAndInfo only tests a single
// store; this test creates three to exercise multi-element parsing and
// confirm the array length matches Paths/RemoteUrls.
func TestLoreSharedStoreInfoExistsArray(t *testing.T) {
	if libErr != nil {
		t.Skipf("Lore library not loaded: %v", libErr)
	}

	tempDir, err := os.MkdirTemp("", "lore-shared-store-exists-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		Offline: true,
	})
	defer cleanupGlobals()

	// Create three shared stores, each with a unique remote URL so we can
	// pick them out of the info response.
	remoteUrls := []string{
		fmt.Sprintf("lore-go-exists-test-%d-a", time.Now().UnixNano()),
		fmt.Sprintf("lore-go-exists-test-%d-b", time.Now().UnixNano()),
		fmt.Sprintf("lore-go-exists-test-%d-c", time.Now().UnixNano()),
	}
	storePaths := []string{
		filepath.Join(tempDir, "store-a"),
		filepath.Join(tempDir, "store-b"),
		filepath.Join(tempDir, "store-c"),
	}

	for i := range remoteUrls {
		createArgs, cleanupCreateArgs := types.NewLoreSharedStoreCreateArgs(types.LoreSharedStoreCreateArgs{
			RemoteUrl:   remoteUrls[i],
			Path:        storePaths[i],
			MakeDefault: true,
		})
		result, err := SharedStoreCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
			Callback:    func(event *types.LoreEventFFI, _ uint64) {},
			UserContext: 0,
		})
		cleanupCreateArgs()
		if err != nil || result != 0 {
			t.Fatalf("LoreSharedStoreCreate[%d]: err=%v result=%d", i, err, result)
		}
	}

	// Query info — exactly one event with all stores listed in parallel arrays.
	infoArgs, cleanupInfo := types.NewLoreSharedStoreInfoArgs(types.LoreSharedStoreInfoArgs{})
	defer cleanupInfo()

	// Capture both the cloned event (Go-side []bool) AND the raw Len() seen
	// inside the callback (FFI-side). They must agree.
	var info types.LoreSharedStoreInfoEventData
	var ffiExistsLen, ffiPathsLen, ffiRemoteUrlsLen int
	var existsAtIndex []bool

	result, err := SharedStoreInfo(&globals, &infoArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag != types.LoreEventTag_SHARED_STORE_INFO {
				return
			}
			if d, ok := event.GetData().(*types.LoreSharedStoreInfoEventDataFFI); ok {
				ffiExistsLen = d.Exists.Len()
				ffiPathsLen = d.Paths.Len()
				ffiRemoteUrlsLen = d.RemoteUrls.Len()
				// Drive the per-index Get() path on the FFI array — independent
				// from the Clone() path used by event.Clone() below.
				existsAtIndex = make([]bool, ffiExistsLen)
				for i := 0; i < ffiExistsLen; i++ {
					existsAtIndex[i] = d.Exists.Get(i)
				}
			}
			if e, ok := event.Clone().Data.(types.LoreSharedStoreInfoEventData); ok {
				info = e
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreSharedStoreInfo: err=%v result=%d", err, result)
	}

	// All three parallel arrays must agree in length.
	if ffiPathsLen != ffiExistsLen || ffiRemoteUrlsLen != ffiExistsLen {
		t.Errorf("FFI parallel-array length mismatch: Paths=%d Exists=%d RemoteUrls=%d",
			ffiPathsLen, ffiExistsLen, ffiRemoteUrlsLen)
	}
	if len(info.Exists) != len(info.Paths) || len(info.Exists) != len(info.RemoteUrls) {
		t.Errorf("Cloned parallel-array length mismatch: Paths=%d Exists=%d RemoteUrls=%d",
			len(info.Paths), len(info.Exists), len(info.RemoteUrls))
	}
	if len(info.Exists) != ffiExistsLen {
		t.Errorf("Cloned Exists len %d does not match FFI Len() %d", len(info.Exists), ffiExistsLen)
	}
	if !reflect.DeepEqual(info.Exists, existsAtIndex) {
		t.Errorf("Clone() vs Get() disagree: clone=%v get=%v", info.Exists, existsAtIndex)
	}

	// Each of our three stores should be present and Exists==true.
	for i, url := range remoteUrls {
		idx := -1
		for j, u := range info.RemoteUrls {
			if u == url {
				idx = j
				break
			}
		}
		if idx < 0 {
			t.Errorf("remote URL %d (%q) not found in info", i, url)
			continue
		}
		if !info.Exists[idx] {
			t.Errorf("Exists[%d] = false for store %q (just created, should be true)", idx, url)
		}
		if !strings.Contains(info.Paths[idx], storePaths[i]) {
			t.Errorf("Paths[%d] = %q, expected to contain %q", idx, info.Paths[idx], storePaths[i])
		}
	}
}

// TestLoreBranchInfoStackMultiElement extends the existing single-element
// stack coverage in TestLoreBranchListStackHandling. It builds a 3-deep
// branch chain (main → b1 → b2 → b3) and verifies that
// LoreBranchPointArrayFFI.Len()/Get()/Clone() all agree and produce
// non-zero branch IDs at every depth.
func TestLoreBranchInfoStackMultiElement(t *testing.T) {
	globals := setupTestRepository(t)

	// LoreBranchCreate implicitly switches to the new branch (see
	// TestLoreBranchSwitch which relies on this behavior). So the chain is
	// built by creating, committing, and creating again — each new branch
	// is a child of the previous tip.
	for i, name := range []string{"b1", "b2", "b3"} {
		createBranch(t, &globals, name)
		// A commit per branch ensures each branch contributes a distinct
		// branch point to the stack.
		fname := fmt.Sprintf("file-%s.txt", name)
		fpath := createFileWithContents(t, globals.RepositoryPath.String(), fname, fmt.Sprintf("content %d", i))
		stageFiles(t, &globals, []string{fpath})
		commitRevision(t, &globals, fmt.Sprintf("commit on %s", name))
	}

	args, cleanupArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{
		Branch: "b3",
	})
	defer cleanupArgs()

	var (
		ffiStackLen int
		stackGet    []types.LoreBranchPoint
		cloned      types.LoreBranchInfoEventData
		infoSeen    bool
	)

	result, err := BranchInfo(&globals, &args, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, _ uint64) {
			if event.Tag != types.LoreEventTag_BRANCH_INFO {
				return
			}
			if d, ok := event.GetData().(*types.LoreBranchInfoEventDataFFI); ok {
				infoSeen = true
				ffiStackLen = d.Stack.Len()
				// Read each entry through Get() while still inside the
				// callback. This is the only safe time for FFI access.
				stackGet = make([]types.LoreBranchPoint, ffiStackLen)
				for i := 0; i < ffiStackLen; i++ {
					stackGet[i] = d.Stack.Get(i)
				}
				cloned = d.Clone()
			}
		},
		UserContext: 0,
	})
	if err != nil || result != 0 {
		t.Fatalf("LoreBranchInfo: err=%v result=%d", err, result)
	}
	if !infoSeen {
		t.Fatal("BRANCH_INFO event not received")
	}

	// b3 was created off b2 which was created off b1 which was created off main.
	// Expect at least 3 ancestor branch points.
	if ffiStackLen < 3 {
		t.Errorf("expected stack length >= 3 for chain main→b1→b2→b3, got %d", ffiStackLen)
	}
	if len(cloned.Stack) != ffiStackLen {
		t.Errorf("Clone() stack len %d != FFI Len() %d", len(cloned.Stack), ffiStackLen)
	}

	var zeroBranch types.LoreBranchId
	for i := 0; i < ffiStackLen; i++ {
		if !reflect.DeepEqual(stackGet[i].Branch, cloned.Stack[i].Branch) {
			t.Errorf("Stack[%d].Branch: Get() vs Clone() disagree", i)
		}
		if !reflect.DeepEqual(stackGet[i].Revision, cloned.Stack[i].Revision) {
			t.Errorf("Stack[%d].Revision: Get() vs Clone() disagree", i)
		}
		if stackGet[i].Branch == zeroBranch {
			t.Errorf("Stack[%d].Branch is zero-valued (uninitialized)", i)
		}
	}

	// All branch points in the stack should be unique — the chain has no
	// duplicate ancestors.
	seen := make(map[types.LoreBranchId]int, ffiStackLen)
	for i, bp := range cloned.Stack {
		if prev, dup := seen[bp.Branch]; dup {
			t.Errorf("duplicate branch in stack at indices %d and %d", prev, i)
		}
		seen[bp.Branch] = i
	}
}

// TestLoreRepositoryStatusParallel reproduces the JS SDK test
// "should support multiple parallel Lore calls" (lore-native.test.ts) against
// the Go FFI bindings.
func TestLoreRepositoryStatusParallel(t *testing.T) {
	globals := setupTestRepository(t)

	// Stage a random file so status has something to report, mirroring
	// stageRandomFile() in the JS test.
	staged := createFileWithContents(t, globals.RepositoryPath.String(), "parallel-staged.txt", "staged content")
	stageFiles(t, &globals, []string{staged})

	const calls = 200

	type callResult struct {
		result         int32
		err            error
		fileEventCount int
		errorMessages  []string
	}
	results := make([]callResult, calls)

	var wg sync.WaitGroup
	wg.Add(calls)

	for i := 0; i < calls; i++ {
		go func(idx int) {
			defer wg.Done()

			// Each goroutine gets its own args and callback, but they all
			// share the single `globals` (and therefore the single C-allocated
			// repository path string) — exactly as the JS test shares one
			// globalArgs across every parallel call.
			args, cleanupArgs := types.NewLoreRepositoryStatusArgs(types.LoreRepositoryStatusArgs{
				Staged: true,
				Scan:   true,
			})
			defer cleanupArgs()

			r := &results[idx]
			callback := func(event *types.LoreEventFFI, userContext uint64) {
				switch event.Tag {
				case types.LoreEventTag_REPOSITORY_STATUS_FILE:
					if _, ok := event.GetData().(*types.LoreRepositoryStatusFileEventDataFFI); ok {
						r.fileEventCount++
					}
				case types.LoreEventTag_ERROR:
					if errData, ok := event.GetData().(*types.LoreErrorEventDataFFI); ok {
						r.errorMessages = append(r.errorMessages, errData.ErrorInner.String())
					}
				}
			}

			result, err := RepositoryStatus(&globals, &args, &types.LoreEventCallbackConfig{
				Callback:    callback,
				UserContext: uint64(idx + 1),
			})
			r.result = result
			r.err = err
		}(i)
	}

	wg.Wait()

	succeeded := 0
	withFileEvents := 0
	for i := range results {
		r := &results[i]
		if r.err != nil {
			t.Errorf("call %d: RepositoryStatus returned error: %v (errorMessages=%v)", i, r.err, r.errorMessages)
			continue
		}
		if r.result != 0 {
			t.Errorf("call %d: RepositoryStatus returned non-zero result %d (errorMessages=%v)", i, r.result, r.errorMessages)
			continue
		}
		succeeded++
		if r.fileEventCount > 0 {
			withFileEvents++
		} else {
			t.Errorf("call %d: succeeded but produced no file events", i)
		}
	}

	if succeeded != calls {
		t.Errorf("expected all %d parallel calls to succeed, got %d", calls, succeeded)
	}
	if withFileEvents != calls {
		t.Errorf("expected all %d parallel calls to produce file events, got %d", calls, withFileEvents)
	}
}
