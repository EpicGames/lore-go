// Copyright Epic Games, Inc. All Rights Reserved.

package lore

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/EpicGames/lore-go/internal/testutil"
	"github.com/EpicGames/lore-go/native"
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

// setupTestRepository creates a temporary Lore repository for testing
func setupTestRepository(t *testing.T) (types.LoreGlobalArgsFFI, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "lore-fluent-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: tempDir,
		Offline:        true,
	})

	cleanup := func() {
		flushArgs, cleanupFlushArgs := types.NewLoreRepositoryFlushArgs(types.LoreRepositoryFlushArgs{})
		defer cleanupFlushArgs()

		native.RepositoryFlush(&globals, &flushArgs, &types.LoreEventCallbackConfig{
			Callback: func(event *types.LoreEventFFI, userContext uint64) {},
		})

		cleanupGlobals()
		os.RemoveAll(tempDir)
	}

	// Create repository
	repoUrl := fmt.Sprintf("test-repo-%d", time.Now().UnixNano())
	createArgs, cleanupCreateArgs := types.NewLoreRepositoryCreateArgs(types.LoreRepositoryCreateArgs{
		RepositoryUrl: repoUrl,
	})
	defer cleanupCreateArgs()

	result, err := native.RepositoryCreate(&globals, &createArgs, &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {},
	})
	if err != nil || result != 0 {
		cleanup()
		t.Fatalf("Failed to create repository: err=%v, result=%d", err, result)
	}

	return globals, cleanup
}

func TestLoreBranchList_Wait_Success(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	returnCode, err := BranchList(&globals, &args).Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
}

func TestLoreBranchList_Callback_ReceivesEvents(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	var receivedTags []types.LoreEventTag

	returnCode, err := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			receivedTags = append(receivedTags, event.Tag)
		}).
		Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	// Should have received at least some events
	if len(receivedTags) == 0 {
		t.Error("Expected to receive events, got none")
	}

	// Should have received END event
	hasEnd := false
	for _, tag := range receivedTags {
		if tag == types.LoreEventTag_END {
			hasEnd = true
			break
		}
	}
	if !hasEnd {
		t.Error("Expected to receive END event")
	}

	t.Logf("Received %d events: %v", len(receivedTags), receivedTags)
}

func TestLoreBranchList_FilterByType(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	var receivedTags []types.LoreEventTag

	// Only filter for BRANCH_LIST_BEGIN and BRANCH_LIST_END events
	returnCode, err := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN, types.LoreEventTag_BRANCH_LIST_END).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			receivedTags = append(receivedTags, event.Tag)
		}).
		Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	// Should only have received filtered events
	for _, tag := range receivedTags {
		if tag != types.LoreEventTag_BRANCH_LIST_BEGIN && tag != types.LoreEventTag_BRANCH_LIST_END {
			t.Errorf("Received unexpected event tag %d, expected only BRANCH_LIST_BEGIN or BRANCH_LIST_END", tag)
		}
	}

	t.Logf("Filtered to %d events: %v", len(receivedTags), receivedTags)
}

func TestLoreBranchList_UserContext(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	expectedContext := uint64(12345)
	var receivedContexts []uint64

	returnCode, err := BranchList(&globals, &args).
		UserContext(expectedContext).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			receivedContexts = append(receivedContexts, userContext)
		}).
		Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	// All callbacks should have received the same user context
	for i, ctx := range receivedContexts {
		if ctx != expectedContext {
			t.Errorf("Event %d: expected user context %d, got %d", i, expectedContext, ctx)
		}
	}
}

func TestLoreBranchInfo_NonZeroReturnCode_ReturnsLoreError(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	// Try to get info for a non-existent branch
	args, cleanupArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{
		Branch: "non-existent-branch-that-does-not-exist",
	})
	defer cleanupArgs()

	returnCode, err := BranchInfo(&globals, &args).Wait()

	// Should return non-zero and an error
	if returnCode == 0 {
		t.Error("Expected non-zero return code for non-existent branch")
	}

	if err == nil {
		t.Error("Expected error for non-zero return code")
	}

	var loreErr *LoreError
	if !errors.As(err, &loreErr) {
		t.Errorf("Expected LoreError, got %T: %v", err, err)
	} else {
		if loreErr.ReturnCode != returnCode {
			t.Errorf("LoreError.ReturnCode %d does not match return code %d", loreErr.ReturnCode, returnCode)
		}
		t.Logf("LoreError: code=%d, messages=%v", loreErr.ReturnCode, loreErr.Messages)
	}
}

func TestLoreCall_ColdHandle_NoExecutionUntilWait(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	callbackCalled := false

	// Create cold handle and configure it
	call := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			callbackCalled = true
		}).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN).
		UserContext(999)

	// Callback should NOT have been called yet (operation hasn't started)
	if callbackCalled {
		t.Error("Callback was called before Wait() - cold handle should not execute until Wait()")
	}

	// Now call Wait()
	_, _ = call.Wait()

	// Now callback should have been called
	if !callbackCalled {
		t.Error("Callback was not called after Wait()")
	}
}

func TestLoreCall_MethodChaining(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	var callbackCalled bool
	var receivedContext uint64
	var receivedTags []types.LoreEventTag

	// Test that all methods can be chained in a single expression
	returnCode, err := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			callbackCalled = true
			receivedContext = userContext
			receivedTags = append(receivedTags, event.Tag)
		}).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN).
		UserContext(42).
		Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
	if !callbackCalled {
		t.Error("Expected callback to be called")
	}
	if receivedContext != 42 {
		t.Errorf("Expected user context 42, got %d", receivedContext)
	}

	// Should only have received BRANCH_LIST_BEGIN due to filter
	for _, tag := range receivedTags {
		if tag != types.LoreEventTag_BRANCH_LIST_BEGIN {
			t.Errorf("Expected only BRANCH_LIST_BEGIN events, got %d", tag)
		}
	}
}

func TestLoreError_Error_WithMessages(t *testing.T) {
	err := &LoreError{
		ReturnCode: 42,
		Messages:   []string{"error 1", "error 2"},
	}

	errorString := err.Error()

	if errorString != "Lore operation failed with code 42: error 1; error 2" {
		t.Errorf("Unexpected error string: %s", errorString)
	}
}

func TestLoreError_Error_WithoutMessages(t *testing.T) {
	err := &LoreError{
		ReturnCode: 42,
		Messages:   nil,
	}

	errorString := err.Error()

	if errorString != "Lore operation failed with code 42" {
		t.Errorf("Unexpected error string: %s", errorString)
	}
}

func TestLoreBranchList_Collect_Success(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	events, err := BranchList(&globals, &args).Collect()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should have collected events
	if len(events) == 0 {
		t.Error("Expected to collect events, got none")
	}

	// Events should be valid LoreEvent (cloned to Go memory)
	hasEnd := false
	for _, event := range events {
		if event.Tag == types.LoreEventTag_END {
			hasEnd = true
		}
	}
	if !hasEnd {
		t.Error("Expected to collect END event")
	}

	t.Logf("Collected %d events", len(events))
}

func TestLoreBranchList_Collect_WithFilter(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Only collect BRANCH_LIST_BEGIN and BRANCH_LIST_END events
	events, err := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN, types.LoreEventTag_BRANCH_LIST_END).
		Collect()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should only have collected filtered events
	for _, event := range events {
		if event.Tag != types.LoreEventTag_BRANCH_LIST_BEGIN && event.Tag != types.LoreEventTag_BRANCH_LIST_END {
			t.Errorf("Collected unexpected event tag %d, expected only BRANCH_LIST_BEGIN or BRANCH_LIST_END", event.Tag)
		}
	}

	t.Logf("Filtered to %d events", len(events))
}

func TestLoreBranchList_Collect_WithCallback_ReturnsError(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Set a callback and then try to call Collect()
	events, err := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {}).
		Collect()

	if err != ErrCallbackSet {
		t.Errorf("Expected ErrCallbackSet, got %v", err)
	}

	if events != nil {
		t.Error("Expected nil events when error is returned")
	}
}

func TestLoreBranchInfo_Collect_NonZeroReturnCode(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	// Try to get info for a non-existent branch
	args, cleanupArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{
		Branch: "non-existent-branch-that-does-not-exist",
	})
	defer cleanupArgs()

	_, err := BranchInfo(&globals, &args).Collect()

	// Should return error
	if err == nil {
		t.Error("Expected error for non-existent branch")
	}

	var loreErr *LoreError
	if !errors.As(err, &loreErr) {
		t.Errorf("Expected LoreError, got %T: %v", err, err)
	}
}

func TestLoreBranchList_Collect_EventDataAccessibleOutsideCallback(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Collect only BRANCH_LIST_ENTRY events which contain branch data
	events, err := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_ENTRY).
		Collect()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least one BRANCH_LIST_ENTRY event")
	}

	// Access the event data outside the FFI callback - this verifies Clone() worked
	for i, event := range events {
		if event.Tag != types.LoreEventTag_BRANCH_LIST_ENTRY {
			t.Errorf("Event %d: expected BRANCH_LIST_ENTRY tag, got %d", i, event.Tag)
			continue
		}

		// Cast to the Go-native event data type (value type, not pointer)
		data, ok := event.Data.(types.LoreBranchListEntryEventData)
		if !ok {
			t.Errorf("Event %d: expected LoreBranchListEntryEventData, got %T", i, event.Data)
			continue
		}

		// Read string fields - these should be valid Go strings copied from FFI memory
		t.Logf("Branch entry %d: Name=%q, Category=%q, Creator=%q, IsCurrent=%v",
			i, data.Name, data.Category, data.Creator, data.IsCurrent)

		// Verify the name is not empty (repository should have at least main branch)
		if data.Name == "" {
			t.Errorf("Event %d: branch name should not be empty", i)
		}

		// Verify we can read the hash (32 bytes)
		hashStr := data.Latest.String()
		if len(hashStr) != 64 { // 32 bytes = 64 hex chars
			t.Errorf("Event %d: expected 64-char hash string, got %d chars: %s", i, len(hashStr), hashStr)
		}

		// Verify we can read the context/ID (16 bytes)
		idStr := data.Id.String()
		if len(idStr) != 32 { // 16 bytes = 32 hex chars
			t.Errorf("Event %d: expected 32-char ID string, got %d chars: %s", i, len(idStr), idStr)
		}
	}
}

func TestLoreCall_Wait_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call should succeed
	_, err := call.Wait()
	if err != nil {
		t.Fatalf("First Wait() failed: %v", err)
	}

	// Second call should return ErrAlreadyStarted
	_, err = call.Wait()
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreCall_Collect_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call should succeed
	_, err := call.Collect()
	if err != nil {
		t.Fatalf("First Collect() failed: %v", err)
	}

	// Second call should return ErrAlreadyStarted
	_, err = call.Collect()
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreCall_Wait_ThenCollect_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call Wait()
	_, err := call.Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Then Collect() should return ErrAlreadyStarted
	_, err = call.Collect()
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreBranchList_AsyncIter_Success(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	eventCh, errCh := BranchList(&globals, &args).AsyncIter()

	var events []types.LoreEvent
	for event := range eventCh {
		events = append(events, event)
	}

	err := <-errCh
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should have collected events
	if len(events) == 0 {
		t.Error("Expected to receive events, got none")
	}

	// Events should include END event
	hasEnd := false
	for _, event := range events {
		if event.Tag == types.LoreEventTag_END {
			hasEnd = true
		}
	}
	if !hasEnd {
		t.Error("Expected to receive END event")
	}

	t.Logf("Received %d events via AsyncIter", len(events))
}

func TestLoreBranchList_AsyncIter_WithFilter(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Only iterate BRANCH_LIST_BEGIN and BRANCH_LIST_END events
	eventCh, errCh := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN, types.LoreEventTag_BRANCH_LIST_END).
		AsyncIter()

	var events []types.LoreEvent
	for event := range eventCh {
		events = append(events, event)
	}

	err := <-errCh
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should only have received filtered events
	for _, event := range events {
		if event.Tag != types.LoreEventTag_BRANCH_LIST_BEGIN && event.Tag != types.LoreEventTag_BRANCH_LIST_END {
			t.Errorf("Received unexpected event tag %d, expected only BRANCH_LIST_BEGIN or BRANCH_LIST_END", event.Tag)
		}
	}

	t.Logf("Filtered to %d events via AsyncIter", len(events))
}

func TestLoreBranchList_AsyncIter_WithCallback_ReturnsError(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Set a callback and then try to call AsyncIter()
	eventCh, errCh := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {}).
		AsyncIter()

	// Event channel should be closed immediately
	events := []types.LoreEvent{}
	for event := range eventCh {
		events = append(events, event)
	}

	if len(events) != 0 {
		t.Error("Expected no events when callback is set")
	}

	// Error channel should have ErrCallbackSet
	err := <-errCh
	if err != ErrCallbackSet {
		t.Errorf("Expected ErrCallbackSet, got %v", err)
	}
}

func TestLoreBranchInfo_AsyncIter_NonZeroReturnCode(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	// Try to get info for a non-existent branch
	args, cleanupArgs := types.NewLoreBranchInfoArgs(types.LoreBranchInfoArgs{
		Branch: "non-existent-branch-that-does-not-exist",
	})
	defer cleanupArgs()

	eventCh, errCh := BranchInfo(&globals, &args).AsyncIter()

	// Drain the event channel
	for range eventCh {
	}

	// Should return error
	err := <-errCh
	if err == nil {
		t.Error("Expected error for non-existent branch")
	}

	var loreErr *LoreError
	if !errors.As(err, &loreErr) {
		t.Errorf("Expected LoreError, got %T: %v", err, err)
	}
}

func TestLoreBranchList_AsyncIter_EventDataAccessibleOutsideCallback(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Iterate only BRANCH_LIST_ENTRY events which contain branch data
	eventCh, errCh := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_ENTRY).
		AsyncIter()

	var events []types.LoreEvent
	for event := range eventCh {
		events = append(events, event)
	}

	err := <-errCh
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(events) == 0 {
		t.Fatal("Expected at least one BRANCH_LIST_ENTRY event")
	}

	// Access the event data outside the channel receive - this verifies Clone() worked
	for i, event := range events {
		if event.Tag != types.LoreEventTag_BRANCH_LIST_ENTRY {
			t.Errorf("Event %d: expected BRANCH_LIST_ENTRY tag, got %d", i, event.Tag)
			continue
		}

		// Cast to the Go-native event data type (value type, not pointer)
		data, ok := event.Data.(types.LoreBranchListEntryEventData)
		if !ok {
			t.Errorf("Event %d: expected LoreBranchListEntryEventData, got %T", i, event.Data)
			continue
		}

		// Read string fields - these should be valid Go strings copied from FFI memory
		t.Logf("Branch entry %d: Name=%q, Category=%q, Creator=%q, IsCurrent=%v",
			i, data.Name, data.Category, data.Creator, data.IsCurrent)

		// Verify the name is not empty (repository should have at least main branch)
		if data.Name == "" {
			t.Errorf("Event %d: branch name should not be empty", i)
		}
	}
}

func TestLoreCall_AsyncIter_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call should succeed
	eventCh, errCh := call.AsyncIter()
	for range eventCh {
	}
	err := <-errCh
	if err != nil {
		t.Fatalf("First AsyncIter() failed: %v", err)
	}

	// Second call should return ErrAlreadyStarted
	eventCh2, errCh2 := call.AsyncIter()
	for range eventCh2 {
	}
	err = <-errCh2
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreCall_Wait_ThenAsyncIter_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call Wait()
	_, err := call.Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Then AsyncIter() should return ErrAlreadyStarted
	eventCh, errCh := call.AsyncIter()
	for range eventCh {
	}
	err = <-errCh
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreCall_AsyncIter_ThenWait_AlreadyStarted(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	call := BranchList(&globals, &args)

	// First call AsyncIter()
	eventCh, errCh := call.AsyncIter()
	for range eventCh {
	}
	err := <-errCh
	if err != nil {
		t.Fatalf("AsyncIter() failed: %v", err)
	}

	// Then Wait() should return ErrAlreadyStarted
	_, err = call.Wait()
	if err != ErrAlreadyStarted {
		t.Errorf("Expected ErrAlreadyStarted, got %v", err)
	}
}

func TestLoreGlobalCallback_Wait(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	var globalLogCount int
	var perCallLogCount int

	// Register global callback for LOG events (common use case for global logging)
	cleanupCallback := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		globalLogCount++
	})
	defer cleanupCallback()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Call Wait() with a per-call callback to count LOG events
	_, err := BranchList(&globals, &args).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_LOG {
				perCallLogCount++
			}
		}).
		Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Global and per-call callbacks should have received the same count
	if globalLogCount != perCallLogCount {
		t.Errorf("Global callback received %d LOG events, but per-call received %d", globalLogCount, perCallLogCount)
	}
	t.Logf("Both callbacks received %d LOG events", globalLogCount)
}

func TestLoreGlobalCallback_Collect(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	var globalLogCount int

	// Register global callback for LOG events
	cleanupCallback := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		globalLogCount++
	})
	defer cleanupCallback()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Call Collect()
	events, err := BranchList(&globals, &args).Collect()
	if err != nil {
		t.Fatalf("Collect() failed: %v", err)
	}

	// Count LOG events in collected events
	collectedLogCount := 0
	for _, event := range events {
		if event.Tag == types.LoreEventTag_LOG {
			collectedLogCount++
		}
	}

	// Global callback should have been called for each LOG event
	if globalLogCount != collectedLogCount {
		t.Errorf("Global callback received %d LOG events, but collected %d", globalLogCount, collectedLogCount)
	}
	t.Logf("Global callback and Collect both received %d LOG events", globalLogCount)
}

func TestLoreGlobalCallback_AsyncIter(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	var globalLogCount int

	// Register global callback for LOG events
	cleanupCallback := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		globalLogCount++
	})
	defer cleanupCallback()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Call AsyncIter()
	eventCh, errCh := BranchList(&globals, &args).AsyncIter()
	iterLogCount := 0
	for event := range eventCh {
		if event.Tag == types.LoreEventTag_LOG {
			iterLogCount++
		}
	}
	if err := <-errCh; err != nil {
		t.Fatalf("AsyncIter() failed: %v", err)
	}

	// Global callback should have been called for each LOG event
	if globalLogCount != iterLogCount {
		t.Errorf("Global callback received %d LOG events, but AsyncIter received %d", globalLogCount, iterLogCount)
	}
	t.Logf("Global callback and AsyncIter both received %d LOG events", globalLogCount)
}

func TestLoreGlobalCallback_Cleanup(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	callCount := 0

	// Register global callback for LOG events
	cleanupCallback := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		callCount++
	})

	args1, cleanupArgs1 := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs1()

	// First call - callback should be called
	_, err := BranchList(&globals, &args1).Wait()
	if err != nil {
		t.Fatalf("First Wait() failed: %v", err)
	}
	firstCallCount := callCount
	t.Logf("Before cleanup: received %d LOG events", firstCallCount)

	// Cleanup the callback
	cleanupCallback()

	args2, cleanupArgs2 := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs2()

	// Second call - callback should NOT be called
	_, err = BranchList(&globals, &args2).Wait()
	if err != nil {
		t.Fatalf("Second Wait() failed: %v", err)
	}

	if callCount != firstCallCount {
		t.Errorf("Global callback was called after cleanup: before=%d, after=%d", firstCallCount, callCount)
	}
}

func TestLoreGlobalCallback_MultipleCallbacksSameType(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	callback1Count := 0
	callback2Count := 0

	// Register two global callbacks for the same event type (LOG)
	cleanup1 := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		callback1Count++
	})
	defer cleanup1()

	cleanup2 := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		callback2Count++
	})
	defer cleanup2()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	_, err := BranchList(&globals, &args).Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Both callbacks should have been called the same number of times
	if callback1Count != callback2Count {
		t.Errorf("Callbacks received different counts: %d vs %d", callback1Count, callback2Count)
	}
	t.Logf("Both callbacks received %d LOG events each", callback1Count)
}

func TestLoreGlobalCallback_IgnoresPerCallFilter(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	globalLogCount := 0
	perCallLogCount := 0

	// Register global callback for LOG events
	cleanupGlobal := GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
		globalLogCount++
	})
	defer cleanupGlobal()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Call with per-call filter that EXCLUDES LOG events (only BRANCH_LIST_BEGIN)
	_, err := BranchList(&globals, &args).
		FilterByType(types.LoreEventTag_BRANCH_LIST_BEGIN).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_LOG {
				perCallLogCount++
			}
		}).
		Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	// Per-call callback should NOT have received any LOG events (filtered out)
	if perCallLogCount != 0 {
		t.Errorf("Per-call callback should not receive LOG events due to filter, but got %d", perCallLogCount)
	}

	// Global callback should STILL have received LOG events (ignores per-call filter)
	if globalLogCount == 0 {
		t.Error("Global callback should receive LOG events regardless of per-call filter")
	}
	t.Logf("Global callback received %d LOG events, per-call received %d (filtered)", globalLogCount, perCallLogCount)
}

func TestLoreCall_Wait_WithoutCallback_Succeeds(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	args, cleanupArgs := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs()

	// Call Wait() without setting any callback — should still succeed
	returnCode, err := BranchList(&globals, &args).Wait()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
}

func TestLoreBranchList_Collect_CompleteAndEndEvents(t *testing.T) {
	globals, cleanup := setupTestRepository(t)
	defer cleanup()

	// Test with Collect()
	args1, cleanupArgs1 := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs1()

	events, err := BranchList(&globals, &args1).Collect()
	if err != nil {
		t.Fatalf("Collect() failed: %v", err)
	}

	hasComplete := false
	hasEnd := false
	for _, event := range events {
		if event.Tag == types.LoreEventTag_COMPLETE {
			hasComplete = true
		}
		if event.Tag == types.LoreEventTag_END {
			hasEnd = true
		}
	}

	if !hasComplete {
		t.Error("Collect(): expected COMPLETE event")
	}
	if !hasEnd {
		t.Error("Collect(): expected END event")
	}

	// Test with Wait() + Callback
	args2, cleanupArgs2 := types.NewLoreBranchListArgs(types.LoreBranchListArgs{})
	defer cleanupArgs2()

	var waitHasComplete, waitHasEnd bool
	_, err = BranchList(&globals, &args2).
		Callback(func(event *types.LoreEventFFI, userContext uint64) {
			if event.Tag == types.LoreEventTag_COMPLETE {
				waitHasComplete = true
			}
			if event.Tag == types.LoreEventTag_END {
				waitHasEnd = true
			}
		}).
		Wait()
	if err != nil {
		t.Fatalf("Wait() failed: %v", err)
	}

	if !waitHasComplete {
		t.Error("Wait(): expected COMPLETE event")
	}
	if !waitHasEnd {
		t.Error("Wait(): expected END event")
	}
}
