// Copyright Epic Games, Inc. All Rights Reserved.

package lore

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/EpicGames/lore-go/native"
	"github.com/EpicGames/lore-go/types"
)

// LoreError represents an error from a Lore operation with a non-zero return code.
type LoreError struct {
	ReturnCode int32
	Messages   []string
}

func (e *LoreError) Error() string {
	if len(e.Messages) > 0 {
		return fmt.Sprintf("Lore operation failed with code %d: %s", e.ReturnCode, strings.Join(e.Messages, "; "))
	}
	return fmt.Sprintf("Lore operation failed with code %d", e.ReturnCode)
}

// ErrCallbackSet is returned when Collect() or AsyncIter() is called on a LoreCall that has a callback set.
var ErrCallbackSet = fmt.Errorf("Collect() or AsyncIter() cannot be used with Callback(); use Wait() instead")

// ErrAlreadyStarted is returned when Wait() or Collect() is called on a LoreCall that has already started.
var ErrAlreadyStarted = fmt.Errorf("operation has already started")

// globalCallbackEntry represents a registered global callback with a unique ID.
type globalCallbackEntry struct {
	id       uint64
	callback types.LoreEventCallback
}

// globalCallbacksSnapshot is an immutable snapshot of callbacks per event type.
// This is used for copy-on-write: writes create a new snapshot, reads use the current one.
type globalCallbacksSnapshot map[types.LoreEventTag][]globalCallbackEntry

// Global callback registry using copy-on-write for efficient reads.
// Writes (register/unregister) are rare and protected by mutex.
// Reads (invoke) are frequent and lock-free using atomic snapshot.
var (
	globalCallbacksPtr   atomic.Pointer[globalCallbacksSnapshot]
	globalCallbacksMutex sync.Mutex // protects writes only
	nextGlobalCallbackID uint64     = 1
)

func init() {
	// Initialize with empty snapshot
	empty := make(globalCallbacksSnapshot)
	globalCallbacksPtr.Store(&empty)
}

// GlobalCallback registers a global callback handler for events with the given type.
// The callback will be invoked for all Lore operations that emit events of this type,
// regardless of which terminating method (Wait, Collect, AsyncIter) is used.
// Multiple callbacks can be registered for the same event type.
// Returns a cleanup function that unregisters the callback when called.
//
// Example usage:
//
//	cleanup := lore.GlobalCallback(types.LoreEventTag_LOG, func(event *types.LoreEventFFI, userContext uint64) {
//	    if data, ok := event.GetData().(*types.LoreLogEventDataFFI); ok {
//	        log.Printf("[Lore] %s", data.Message.String())
//	    }
//	})
//	defer cleanup()
func GlobalCallback(eventType types.LoreEventTag, callback types.LoreEventCallback) func() {
	globalCallbacksMutex.Lock()
	defer globalCallbacksMutex.Unlock()

	id := nextGlobalCallbackID
	nextGlobalCallbackID++

	// Create new snapshot with the added callback
	oldSnapshot := *globalCallbacksPtr.Load()
	newSnapshot := make(globalCallbacksSnapshot, len(oldSnapshot))
	for k, v := range oldSnapshot {
		newSnapshot[k] = v
	}
	newSnapshot[eventType] = append(newSnapshot[eventType], globalCallbackEntry{
		id:       id,
		callback: callback,
	})
	globalCallbacksPtr.Store(&newSnapshot)

	// Return cleanup function
	return func() {
		globalCallbacksMutex.Lock()
		defer globalCallbacksMutex.Unlock()

		oldSnapshot := *globalCallbacksPtr.Load()
		entries := oldSnapshot[eventType]

		// Find the entry index
		idx := -1
		for i, entry := range entries {
			if entry.id == id {
				idx = i
				break
			}
		}
		if idx == -1 {
			return // Already removed
		}

		// Create new snapshot without this entry
		newSnapshot := make(globalCallbacksSnapshot, len(oldSnapshot))
		for k, v := range oldSnapshot {
			newSnapshot[k] = v
		}
		newEntries := make([]globalCallbackEntry, 0, len(entries)-1)
		newEntries = append(newEntries, entries[:idx]...)
		newEntries = append(newEntries, entries[idx+1:]...)
		if len(newEntries) > 0 {
			newSnapshot[eventType] = newEntries
		} else {
			delete(newSnapshot, eventType)
		}
		globalCallbacksPtr.Store(&newSnapshot)
	}
}

// invokeGlobalCallbacks calls all registered global callbacks for the given event.
// This is lock-free and allocation-free on the read path.
func invokeGlobalCallbacks(event *types.LoreEventFFI, userContext uint64) {
	snapshot := globalCallbacksPtr.Load()
	entries := (*snapshot)[event.Tag]
	for _, entry := range entries {
		entry.callback(event, userContext)
	}
}

// LoreCallExecutor is the function type that executes a Lore operation.
type LoreCallExecutor[TArgs any] func(
	globals *types.LoreGlobalArgsFFI,
	args *TArgs,
	config *types.LoreEventCallbackConfig,
) (int32, error)

// LoreCall is a cold handle for a Lore operation.
// It does not execute until a terminating method like Wait() is called.
type LoreCall[TArgs any] struct {
	globals     *types.LoreGlobalArgsFFI
	args        *TArgs
	callback    types.LoreEventCallback
	filterTypes map[types.LoreEventTag]struct{}
	userContext uint64
	execFunc    LoreCallExecutor[TArgs]
	started     bool
}

// Callback sets the event handler for this operation.
// Returns the same handle for method chaining.
func (c *LoreCall[TArgs]) Callback(fn types.LoreEventCallback) *LoreCall[TArgs] {
	c.callback = fn
	return c
}

// FilterByType sets which event types should be passed to the callback.
// Only events with tags in the provided list will trigger the callback.
// If no filter is set, all events are passed to the callback.
// Returns the same handle for method chaining.
func (c *LoreCall[TArgs]) FilterByType(tags ...types.LoreEventTag) *LoreCall[TArgs] {
	c.filterTypes = make(map[types.LoreEventTag]struct{}, len(tags))
	for _, tag := range tags {
		c.filterTypes[tag] = struct{}{}
	}
	return c
}

// UserContext sets a user-provided context value that will be passed to the callback.
// Returns the same handle for method chaining.
func (c *LoreCall[TArgs]) UserContext(ctx uint64) *LoreCall[TArgs] {
	c.userContext = ctx
	return c
}

// Wait executes the Lore operation and blocks until completion.
// Returns the return code and an error if the return code was non-zero.
func (c *LoreCall[TArgs]) Wait() (int32, error) {
	if c.started {
		return -1, ErrAlreadyStarted
	}
	c.started = true

	var errorMessages []string

	callbackConfig := &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			// Invoke global callbacks first
			invokeGlobalCallbacks(event, userContext)

			// Collect error messages from ERROR events
			if event.Tag == types.LoreEventTag_ERROR {
				if data, ok := event.GetData().(*types.LoreErrorEventDataFFI); ok {
					errorMessages = append(errorMessages, data.ErrorInner.String())
				}
			}

			// Call user's callback if set and event passes filter
			if c.callback != nil {
				if len(c.filterTypes) == 0 {
					c.callback(event, userContext)
				} else if _, ok := c.filterTypes[event.Tag]; ok {
					c.callback(event, userContext)
				}
			}
		},
		UserContext: c.userContext,
	}

	returnCode, err := c.execFunc(c.globals, c.args, callbackConfig)
	if err != nil {
		return -1, err
	}

	if returnCode != 0 {
		return returnCode, &LoreError{
			ReturnCode: returnCode,
			Messages:   errorMessages,
		}
	}

	return returnCode, nil
}

// Collect executes the Lore operation and collects all events into a slice.
// Returns the collected events and an error if the operation failed.
// Events are cloned to Go memory and remain valid after the call completes.
// If FilterByType was called, only matching events are collected.
// Cannot be used together with Callback() - returns ErrCallbackSet if a callback is set.
func (c *LoreCall[TArgs]) Collect() ([]types.LoreEvent, error) {
	if c.callback != nil {
		return nil, ErrCallbackSet
	}
	if c.started {
		return nil, ErrAlreadyStarted
	}
	c.started = true

	var events []types.LoreEvent
	var errorMessages []string

	callbackConfig := &types.LoreEventCallbackConfig{
		Callback: func(event *types.LoreEventFFI, userContext uint64) {
			// Invoke global callbacks first
			invokeGlobalCallbacks(event, userContext)

			// Collect error messages from ERROR events
			if event.Tag == types.LoreEventTag_ERROR {
				if data, ok := event.GetData().(*types.LoreErrorEventDataFFI); ok {
					errorMessages = append(errorMessages, data.ErrorInner.String())
				}
			}

			// Collect event if it passes filter
			if len(c.filterTypes) == 0 {
				events = append(events, event.Clone())
			} else if _, ok := c.filterTypes[event.Tag]; ok {
				events = append(events, event.Clone())
			}
		},
		UserContext: c.userContext,
	}

	returnCode, err := c.execFunc(c.globals, c.args, callbackConfig)
	if err != nil {
		return nil, err
	}

	if returnCode != 0 {
		return events, &LoreError{
			ReturnCode: returnCode,
			Messages:   errorMessages,
		}
	}

	return events, nil
}

// AsyncIter executes the Lore operation and returns channels for async iteration.
// Returns an event channel and an error channel.
// Events are cloned to Go memory and remain valid after the call completes.
// If FilterByType was called, only matching events are sent to the channel.
// The event channel is closed when the operation completes.
// The error channel receives the final error (or nil) and is then closed.
// Cannot be used together with Callback() - returns ErrCallbackSet if a callback is set.
//
// Example usage:
//
//	events, errCh := lore.RevisionHistory(globals, args).
//	    FilterByType(types.LoreEventTag_REVISION_HISTORY_ENTRY).
//	    AsyncIter()
//	for event := range events {
//	    // process event
//	}
//	if err := <-errCh; err != nil {
//	    // handle error
//	}
func (c *LoreCall[TArgs]) AsyncIter() (<-chan types.LoreEvent, <-chan error) {
	eventCh := make(chan types.LoreEvent)
	errCh := make(chan error, 1)

	if c.callback != nil {
		close(eventCh)
		errCh <- ErrCallbackSet
		close(errCh)
		return eventCh, errCh
	}
	if c.started {
		close(eventCh)
		errCh <- ErrAlreadyStarted
		close(errCh)
		return eventCh, errCh
	}
	c.started = true

	go func() {
		defer close(eventCh)
		defer close(errCh)

		var errorMessages []string

		callbackConfig := &types.LoreEventCallbackConfig{
			Callback: func(event *types.LoreEventFFI, userContext uint64) {
				// Invoke global callbacks first
				invokeGlobalCallbacks(event, userContext)

				// Collect error messages from ERROR events
				if event.Tag == types.LoreEventTag_ERROR {
					if data, ok := event.GetData().(*types.LoreErrorEventDataFFI); ok {
						errorMessages = append(errorMessages, data.ErrorInner.String())
					}
				}

				// Send event if it passes filter
				if len(c.filterTypes) == 0 {
					eventCh <- event.Clone()
				} else if _, ok := c.filterTypes[event.Tag]; ok {
					eventCh <- event.Clone()
				}
			},
			UserContext: c.userContext,
		}

		returnCode, err := c.execFunc(c.globals, c.args, callbackConfig)
		if err != nil {
			errCh <- err
			return
		}

		if returnCode != 0 {
			errCh <- &LoreError{
				ReturnCode: returnCode,
				Messages:   errorMessages,
			}
			return
		}

		errCh <- nil
	}()

	return eventCh, errCh
}

/* Resolve user IDs to display names using the remote authentication service.
Requires an authenticated connection.

When no user IDs are provided, returns the current user's identity using
locally cached tokens (equivalent to `lore_auth_local_user_info`).

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name for each resolved user | */
func AuthUserInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthUserInfoArgsFFI,
) *LoreCall[types.LoreAuthUserInfoArgsFFI] {
	return &LoreCall[types.LoreAuthUserInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthUserInfo,
	}
}

/* Authenticate using an existing bearer token.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name after successful token authentication | */
func AuthLoginWithToken(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLoginWithTokenArgsFFI,
) *LoreCall[types.LoreAuthLoginWithTokenArgsFFI] {
	return &LoreCall[types.LoreAuthLoginWithTokenArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthLoginWithToken,
	}
}

/* List all stored authentication identities.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_IDENTITY` | `lore_auth_identity_event_data_t` | Emitted once per stored identity | */
func AuthList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthListArgsFFI,
) *LoreCall[types.LoreAuthListArgsFFI] {
	return &LoreCall[types.LoreAuthListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthList,
	}
}

/* Remove stored authentication and authorization tokens.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func AuthLogout(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLogoutArgsFFI,
) *LoreCall[types.LoreAuthLogoutArgsFFI] {
	return &LoreCall[types.LoreAuthLogoutArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthLogout,
	}
}

/* Clear all stored authentication data.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func AuthClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthClearArgsFFI,
) *LoreCall[types.LoreAuthClearArgsFFI] {
	return &LoreCall[types.LoreAuthClearArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthClear,
	}
}

/* Resolve user identities to display names from locally stored JWT tokens.

Does not contact the auth service. Decodes cached JWT tokens to extract
display names. For user IDs without a local token, returns the raw user
ID. For remote resolution with proper authorization, use
`lore_auth_user_info` which queries the remote authentication service.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with the resolved user id and display name | */
func AuthLocalUserInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLocalUserInfoArgsFFI,
) *LoreCall[types.LoreAuthLocalUserInfoArgsFFI] {
	return &LoreCall[types.LoreAuthLocalUserInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthLocalUserInfo,
	}
}

/* Authenticate interactively via a browser-based login flow.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_URL` | `lore_auth_url_event_data_t` | Emitted with the login URL when no_browser mode is requested |
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name after successful interactive authentication | */
func AuthLoginInteractive(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLoginInteractiveArgsFFI,
) *LoreCall[types.LoreAuthLoginInteractiveArgsFFI] {
	return &LoreCall[types.LoreAuthLoginInteractiveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.AuthLoginInteractive,
	}
}

/* Create a new branch in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_CREATE` | `lore_branch_create_event_data_t` | Emitted when the branch has been successfully created, includes branch name and id | */
func BranchCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchCreateArgsFFI,
) *LoreCall[types.LoreBranchCreateArgsFFI] {
	return &LoreCall[types.LoreBranchCreateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchCreate,
	}
}

/* Retrieve metadata about a specific branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_INFO` | `lore_branch_info_event_data_t` | Emitted with branch metadata (name, id, category, protection status, etc.) | */
func BranchInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchInfoArgsFFI,
) *LoreCall[types.LoreBranchInfoArgsFFI] {
	return &LoreCall[types.LoreBranchInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchInfo,
	}
}

/* Show the changes and conflicts between two branches.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_DIFF_BEGIN` | `lore_branch_diff_begin_event_data_t` | Emitted before diff results begin streaming |
| `LORE_EVENT_BRANCH_DIFF_CHANGE_BEGIN` | `lore_branch_diff_change_begin_event_data_t` | Emitted before the list of changed files begins |
| `LORE_EVENT_BRANCH_DIFF_CHANGE` | `lore_branch_diff_change_event_data_t` | Emitted for each changed file between the two branches |
| `LORE_EVENT_BRANCH_DIFF_CHANGE_END` | `lore_branch_diff_change_end_event_data_t` | Emitted after all changed files have been reported |
| `LORE_EVENT_BRANCH_DIFF_CONFLICT_BEGIN` | `lore_branch_diff_conflict_begin_event_data_t` | Emitted before the list of conflicting files begins |
| `LORE_EVENT_BRANCH_DIFF_CONFLICT` | `lore_branch_diff_conflict_event_data_t` | Emitted for each file that has a conflict between the two branches |
| `LORE_EVENT_BRANCH_DIFF_CONFLICT_END` | `lore_branch_diff_conflict_end_event_data_t` | Emitted after all conflict files have been reported |
| `LORE_EVENT_BRANCH_DIFF_END` | `lore_branch_diff_end_event_data_t` | Emitted after all diff results have been streamed | */
func BranchDiff(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchDiffArgsFFI,
) *LoreCall[types.LoreBranchDiffArgsFFI] {
	return &LoreCall[types.LoreBranchDiffArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchDiff,
	}
}

/* Enable write protection on a branch to prevent direct commits.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_PROTECT` | `lore_branch_protect_event_data_t` | Emitted when the branch has been successfully protected | */
func BranchProtect(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchProtectArgsFFI,
) *LoreCall[types.LoreBranchProtectArgsFFI] {
	return &LoreCall[types.LoreBranchProtectArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchProtect,
	}
}

/* Remove write protection from a branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_UNPROTECT` | `lore_branch_unprotect_event_data_t` | Emitted when the branch has been successfully unprotected | */
func BranchUnprotect(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchUnprotectArgsFFI,
) *LoreCall[types.LoreBranchUnprotectArgsFFI] {
	return &LoreCall[types.LoreBranchUnprotectArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchUnprotect,
	}
}

/* Archive a branch in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_ARCHIVE` | `lore_branch_archive_event_data_t` | Emitted when the branch has been successfully archived | */
func BranchArchive(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchArchiveArgsFFI,
) *LoreCall[types.LoreBranchArchiveArgsFFI] {
	return &LoreCall[types.LoreBranchArchiveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchArchive,
	}
}

/* List all branches in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_LIST_BEGIN` | `lore_branch_list_begin_event_data_t` | Emitted before branch list entries begin streaming |
| `LORE_EVENT_BRANCH_LIST_ENTRY` | `lore_branch_list_entry_event_data_t` | Emitted for each branch in the repository |
| `LORE_EVENT_BRANCH_LIST_END` | `lore_branch_list_end_event_data_t` | Emitted after all branch entries have been streamed | */
func BranchList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchListArgsFFI,
) *LoreCall[types.LoreBranchListArgsFFI] {
	return &LoreCall[types.LoreBranchListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchList,
	}
}

/* Abort an in-progress branch merge and restore the pre-merge state.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_ABORT_BEGIN` | `lore_branch_merge_abort_begin_event_data_t` | Emitted when aborting a branch merge, includes staged and current revision hashes |
| `LORE_EVENT_BRANCH_MERGE_ABORT_END` | `lore_branch_merge_abort_end_event_data_t` | Emitted after the merge abort has been completed |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization while reverting merge changes | */
func BranchMergeAbort(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeAbortArgsFFI,
) *LoreCall[types.LoreBranchMergeAbortArgsFFI] {
	return &LoreCall[types.LoreBranchMergeAbortArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeAbort,
	}
}

/* Mark conflicting files in a merge as unresolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_UNRESOLVE_FILE` | `lore_branch_merge_unresolve_file_event_data_t` | Emitted for each file that was marked as unresolved |
| `LORE_EVENT_BRANCH_MERGE_UNRESOLVE_REVISION` | `lore_branch_merge_unresolve_revision_event_data_t` | Emitted with the updated staged revision after unresolve completes | */
func BranchMergeUnresolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeUnresolveArgsFFI,
) *LoreCall[types.LoreBranchMergeUnresolveArgsFFI] {
	return &LoreCall[types.LoreBranchMergeUnresolveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeUnresolve,
	}
}

/* Merge the current branch into a target branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_INTO_FILE_BEGIN` | `lore_branch_merge_into_file_begin_event_data_t` | Emitted when starting to merge files into the target branch |
| `LORE_EVENT_BRANCH_MERGE_INTO_FILE` | `lore_branch_merge_into_file_event_data_t` | Emitted for each file being merged into the target branch |
| `LORE_EVENT_BRANCH_MERGE_INTO_FILE_END` | `lore_branch_merge_into_file_end_event_data_t` | Emitted after all files have been merged |
| `LORE_EVENT_BRANCH_MERGE_INTO_FRAGMENT_BEGIN` | `lore_branch_merge_into_fragment_begin_event_data_t` | Emitted when starting fragment transfer for a file |
| `LORE_EVENT_BRANCH_MERGE_INTO_FRAGMENT_PROGRESS` | `lore_branch_merge_into_fragment_progress_event_data_t` | Emitted periodically during fragment transfer |
| `LORE_EVENT_BRANCH_MERGE_INTO_FRAGMENT_END` | `lore_branch_merge_into_fragment_end_event_data_t` | Emitted when fragment transfer for a file completes |
| `LORE_EVENT_BRANCH_MERGE_INTO_REVISION` | `lore_branch_merge_into_revision_event_data_t` | Emitted with the resulting revision after the merge into is complete |
| `LORE_EVENT_BRANCH_MERGE_INTO_SYNC_BEGIN` | `lore_branch_merge_into_sync_begin_event_data_t` | Emitted when starting to apply the changes on the target state |
| `LORE_EVENT_BRANCH_MERGE_INTO_SYNC_END` | `lore_branch_merge_into_sync_end_event_data_t` | Emitted after applying the changes on the target state is complete |
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when auto-commit starts (if no conflicts) |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted periodically during auto-commit file processing |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when auto-commit file processing completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed revision details |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during changes realization |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata entry of the committed revision |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for each file fragment written or deduplicated during commit | */
func BranchMergeInto(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeIntoArgsFFI,
) *LoreCall[types.LoreBranchMergeIntoArgsFFI] {
	return &LoreCall[types.LoreBranchMergeIntoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeInto,
	}
}

/* Mark conflicting files in a merge as resolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file that was marked as resolved |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision after resolve completes | */
func BranchMergeResolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveArgsFFI,
) *LoreCall[types.LoreBranchMergeResolveArgsFFI] {
	return &LoreCall[types.LoreBranchMergeResolveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeResolve,
	}
}

/* Resolve a merge conflict by accepting the "mine" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file resolved by keeping "mine" |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func BranchMergeResolveMine(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveMineArgsFFI,
) *LoreCall[types.LoreBranchMergeResolveMineArgsFFI] {
	return &LoreCall[types.LoreBranchMergeResolveMineArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeResolveMine,
	}
}

/* Resolve a merge conflict by accepting the "theirs" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file resolved by keeping "theirs" |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func BranchMergeResolveTheirs(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveTheirsArgsFFI,
) *LoreCall[types.LoreBranchMergeResolveTheirsArgsFFI] {
	return &LoreCall[types.LoreBranchMergeResolveTheirsArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeResolveTheirs,
	}
}

/* Restart an in-progress merge, re-materializing conflicted files.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_CONFLICT_FILE` | `lore_branch_merge_conflict_file_event_data_t` | Emitted for each file with a remaining merge conflict |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization during restart |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file re-materialized during restart | */
func BranchMergeRestart(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeRestartArgsFFI,
) *LoreCall[types.LoreBranchMergeRestartArgsFFI] {
	return &LoreCall[types.LoreBranchMergeRestartArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeRestart,
	}
}

/* Start a merge from another branch into the current branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_START_BEGIN` | `lore_branch_merge_start_begin_event_data_t` | Emitted when merge begins, includes source branch and revision info |
| `LORE_EVENT_BRANCH_MERGE_START_END` | `lore_branch_merge_start_end_event_data_t` | Emitted when merge operation completes, includes sync stats and conflict flag |
| `LORE_EVENT_BRANCH_MERGE_CONFLICT_FILE` | `lore_branch_merge_conflict_file_event_data_t` | Emitted for each file with an unresolved merge conflict |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during the apply_diff phase of the merge |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file modified during merge realization |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged for deletion during merge realization |
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when auto-commit starts (no conflicts, no_commit=false) |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted periodically during auto-commit |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when auto-commit file processing completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed revision details |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata entry of the committed revision |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for each fragment written during auto-commit | */
func BranchMergeStart(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeStartArgsFFI,
) *LoreCall[types.LoreBranchMergeStartArgsFFI] {
	return &LoreCall[types.LoreBranchMergeStartArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMergeStart,
	}
}

/* Switch to a different branch and update the working directory.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_SWITCH_BEGIN` | `lore_branch_switch_begin_event_data_t` | Emitted when branch switch starts |
| `LORE_EVENT_BRANCH_SWITCH_END` | `lore_branch_switch_end_event_data_t` | Emitted when branch switch completes successfully |
| `LORE_EVENT_REVISION_SYNC_TARGET` | `lore_revision_sync_target_event_data_t` | Emitted with target revision info after resolving the switch target |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file modified/added/deleted during switch |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted periodically during file realization |
| `LORE_EVENT_REVISION_SYNC_REVISION` | `lore_revision_sync_revision_event_data_t` | Emitted with the resulting revision after switch |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by view or ignore filters |
| `LORE_EVENT_REVISION_RESOLVE` | `lore_revision_resolve_event_data_t` | Emitted when resolving a partial revision reference | */
func BranchSwitch(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchSwitchArgsFFI,
) *LoreCall[types.LoreBranchSwitchArgsFFI] {
	return &LoreCall[types.LoreBranchSwitchArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchSwitch,
	}
}

/* Reset the current branch to a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_RESET` | `lore_branch_reset_event_data_t` | Emitted when the branch has been reset to the target revision | */
func BranchReset(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchResetArgsFFI,
) *LoreCall[types.LoreBranchResetArgsFFI] {
	return &LoreCall[types.LoreBranchResetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchReset,
	}
}

/* Push local branch commits to the remote repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_PUSH` | `lore_branch_push_event_data_t` | Emitted when push begins, includes branch name and revision info |
| `LORE_EVENT_BRANCH_PUSH_BRANCH_CREATE_BEGIN` | `lore_branch_push_branch_create_begin_event_data_t` | Emitted when creating the remote branch (first push) |
| `LORE_EVENT_BRANCH_PUSH_BRANCH_CREATE_END` | `lore_branch_push_branch_create_end_event_data_t` | Emitted when remote branch creation completes |
| `LORE_EVENT_BRANCH_PUSH_REVISION_UPDATE_BEGIN` | `lore_branch_push_revision_update_begin_event_data_t` | Emitted when updating a revision on the remote |
| `LORE_EVENT_BRANCH_PUSH_REVISION_UPDATE_END` | `lore_branch_push_revision_update_end_event_data_t` | Emitted when a revision update completes |
| `LORE_EVENT_BRANCH_PUSH_FRAGMENT_BEGIN` | `lore_branch_push_fragment_begin_event_data_t` | Emitted when uploading fragment data begins |
| `LORE_EVENT_BRANCH_PUSH_FRAGMENT_PROGRESS` | `lore_branch_push_fragment_progress_event_data_t` | Emitted periodically during fragment upload |
| `LORE_EVENT_BRANCH_PUSH_FRAGMENT_END` | `lore_branch_push_fragment_end_event_data_t` | Emitted when fragment upload completes |
| `LORE_EVENT_BRANCH_PUSH_REVISION_PUSH_BEGIN` | `lore_branch_push_revision_push_begin_event_data_t` | Emitted when pushing a revision to the remote begins |
| `LORE_EVENT_BRANCH_PUSH_REVISION_PUSH_UPDATE` | `lore_branch_push_revision_push_update_event_data_t` | Emitted with progress updates during revision push |
| `LORE_EVENT_BRANCH_PUSH_REVISION_PUSH_END` | `lore_branch_push_revision_push_end_event_data_t` | Emitted when revision push completes | */
func BranchPush(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchPushArgsFFI,
) *LoreCall[types.LoreBranchPushArgsFFI] {
	return &LoreCall[types.LoreBranchPushArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchPush,
	}
}

/* Retrieve branch metadata. */
func BranchMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataGetArgsFFI,
) *LoreCall[types.LoreBranchMetadataGetArgsFFI] {
	return &LoreCall[types.LoreBranchMetadataGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMetadataGet,
	}
}

/* Set branch metadata key-value pairs. */
func BranchMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataSetArgsFFI,
) *LoreCall[types.LoreBranchMetadataSetArgsFFI] {
	return &LoreCall[types.LoreBranchMetadataSetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMetadataSet,
	}
}

/* Clear branch metadata keys. */
func BranchMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataClearArgsFFI,
) *LoreCall[types.LoreBranchMetadataClearArgsFFI] {
	return &LoreCall[types.LoreBranchMetadataClearArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.BranchMetadataClear,
	}
}

/* Retrieve metadata for one or more files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_INFO` | `lore_file_info_event_data_t` | Emitted for each file with its metadata (size, hash, staged status, etc.) | */
func FileInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileInfoArgsFFI,
) *LoreCall[types.LoreFileInfoArgsFFI] {
	return &LoreCall[types.LoreFileInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileInfo,
	}
}

/* Show which files differ between two revisions.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DIFF` | `lore_file_diff_event_data_t` | Emitted for each file that differs between the two revisions | */
func FileDiff(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDiffArgsFFI,
) *LoreCall[types.LoreFileDiffArgsFFI] {
	return &LoreCall[types.LoreFileDiffArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDiff,
	}
}

/* Compute the hash of a local file for comparison with repository content.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_HASH` | `lore_file_hash_event_data_t` | Emitted with the computed hash and size of the specified file | */
func FileHash(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileHashArgsFFI,
) *LoreCall[types.LoreFileHashArgsFFI] {
	return &LoreCall[types.LoreFileHashArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileHash,
	}
}

/* Retrieve the revision history for a specific file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_HISTORY` | `lore_file_history_event_data_t` | Emitted for each revision in which the file was modified | */
func FileHistory(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileHistoryArgsFFI,
) *LoreCall[types.LoreFileHistoryArgsFFI] {
	return &LoreCall[types.LoreFileHistoryArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileHistory,
	}
}

/* Clear all metadata from a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA_CLEAR_FILE` | `lore_metadata_clear_file_event_data_t` | Emitted when metadata has been cleared for the file | */
func FileMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataClearArgsFFI,
) *LoreCall[types.LoreFileMetadataClearArgsFFI] {
	return &LoreCall[types.LoreFileMetadataClearArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileMetadataClear,
	}
}

/* Get a specific metadata key/value pair from a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for the requested metadata key/value pair | */
func FileMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataGetArgsFFI,
) *LoreCall[types.LoreFileMetadataGetArgsFFI] {
	return &LoreCall[types.LoreFileMetadataGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileMetadataGet,
	}
}

/* List all metadata key/value pairs associated with a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata key/value pair associated with the file | */
func FileMetadataList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataListArgsFFI,
) *LoreCall[types.LoreFileMetadataListArgsFFI] {
	return &LoreCall[types.LoreFileMetadataListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileMetadataList,
	}
}

/* Set a metadata key/value pair on a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataSetArgsFFI,
) *LoreCall[types.LoreFileMetadataSetArgsFFI] {
	return &LoreCall[types.LoreFileMetadataSetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileMetadataSet,
	}
}

/* Reset files to the state recorded in the current or target revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_RESET_BEGIN` | `lore_file_reset_begin_event_data_t` | Emitted when reset starts, includes path count |
| `LORE_EVENT_FILE_RESET_PROGRESS` | `lore_file_reset_progress_event_data_t` | Emitted periodically during file reset with progress counts |
| `LORE_EVENT_FILE_RESET_END` | `lore_file_reset_end_event_data_t` | Emitted when reset completes |
| `LORE_EVENT_FILE_RESET_FILE` | `lore_file_reset_file_event_data_t` | Emitted for each file that was reset |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file materialized |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by filters | */
func FileReset(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileResetArgsFFI,
) *LoreCall[types.LoreFileResetArgsFFI] {
	return &LoreCall[types.LoreFileResetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileReset,
	}
}

/* Reset files to their state at the last merged revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_RESET_BEGIN` | `lore_file_reset_begin_event_data_t` | Emitted when reset starts |
| `LORE_EVENT_FILE_RESET_PROGRESS` | `lore_file_reset_progress_event_data_t` | Emitted periodically during file reset |
| `LORE_EVENT_FILE_RESET_END` | `lore_file_reset_end_event_data_t` | Emitted when reset completes |
| `LORE_EVENT_FILE_RESET_FILE` | `lore_file_reset_file_event_data_t` | Emitted for each file that was reset |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file materialized | */
func FileResetToLastMerged(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileResetToLastMergedArgsFFI,
) *LoreCall[types.LoreFileResetToLastMergedArgsFFI] {
	return &LoreCall[types.LoreFileResetToLastMergedArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileResetToLastMerged,
	}
}

/* Stage files for the next commit.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_STAGE_BEGIN` | `lore_file_stage_begin_event_data_t` | Emitted when staging begins, includes path count |
| `LORE_EVENT_FILE_STAGE_PROGRESS` | `lore_file_stage_progress_event_data_t` | Emitted periodically during staging with file counts |
| `LORE_EVENT_FILE_STAGE_END` | `lore_file_stage_end_event_data_t` | Emitted when staging completes |
| `LORE_EVENT_FILE_STAGE_REVISION` | `lore_file_stage_revision_event_data_t` | Emitted with the resulting staged revision |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged or staged for deletion |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by filters | */
func FileStage(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileStageArgsFFI,
) *LoreCall[types.LoreFileStageArgsFFI] {
	return &LoreCall[types.LoreFileStageArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileStage,
	}
}

/* Stage files for a merge commit, recording resolved merge content.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_STAGE_BEGIN` | `lore_file_stage_begin_event_data_t` | Emitted when merge-staging begins |
| `LORE_EVENT_FILE_STAGE_PROGRESS` | `lore_file_stage_progress_event_data_t` | Emitted periodically during merge-staging |
| `LORE_EVENT_FILE_STAGE_REVISION` | `lore_file_stage_revision_event_data_t` | Emitted with the resulting staged revision |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged | */
func FileStageMerge(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileStageMergeArgsFFI,
) *LoreCall[types.LoreFileStageMergeArgsFFI] {
	return &LoreCall[types.LoreFileStageMergeArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileStageMerge,
	}
}

/* Stage a file move (rename) operation for commit.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_STAGE_BEGIN` | `lore_file_stage_begin_event_data_t` | Emitted when move staging begins |
| `LORE_EVENT_FILE_STAGE_END` | `lore_file_stage_end_event_data_t` | Emitted when move staging completes |
| `LORE_EVENT_FILE_STAGE_REVISION` | `lore_file_stage_revision_event_data_t` | Emitted with the resulting staged revision |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged (deletion of original and new path) | */
func FileStageMove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileStageMoveArgsFFI,
) *LoreCall[types.LoreFileStageMoveArgsFFI] {
	return &LoreCall[types.LoreFileStageMoveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileStageMove,
	}
}

/* Mark files as dirty in the staged state without staging their content.

Action is determined by checking filesystem existence and current revision state
(modify, add, delete, or revert-add). Respects ignore and view filters.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_PATH_IGNORE` | `lore_path_ignore_event_data_t` | Emitted for each input path that could not be resolved to a repository-relative path |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by view or ignore filters | */
func FileDirty(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyArgsFFI,
) *LoreCall[types.LoreFileDirtyArgsFFI] {
	return &LoreCall[types.LoreFileDirtyArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDirty,
	}
}

/* Mark a file as dirty-moved from one path to another in the staged state.

Updates the source node's parent/name and flags it with `DirtyMove`, propagating
`Dirty` to both the old and new parent directories. For directories, the move
is propagated recursively to children. No filesystem access is performed.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileDirtyMove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyMoveArgsFFI,
) *LoreCall[types.LoreFileDirtyMoveArgsFFI] {
	return &LoreCall[types.LoreFileDirtyMoveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDirtyMove,
	}
}

/* Mark a file as dirty-copied from one path to another in the staged state.

Creates a new destination node flagged `DirtyCopy`; the source node is unchanged.
No filesystem access is performed.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileDirtyCopy(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyCopyArgsFFI,
) *LoreCall[types.LoreFileDirtyCopyArgsFFI] {
	return &LoreCall[types.LoreFileDirtyCopyArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDirtyCopy,
	}
}

/* Remove files from the staging area without discarding local changes.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_UNSTAGE_BEGIN` | `lore_file_unstage_begin_event_data_t` | Emitted when unstage begins, includes path count |
| `LORE_EVENT_FILE_UNSTAGE_PROGRESS` | `lore_file_unstage_progress_event_data_t` | Emitted periodically during unstaging |
| `LORE_EVENT_FILE_UNSTAGE_END` | `lore_file_unstage_end_event_data_t` | Emitted when unstaging completes |
| `LORE_EVENT_FILE_UNSTAGE_REVISION` | `lore_file_unstage_revision_event_data_t` | Emitted with the resulting staged revision |
| `LORE_EVENT_FILE_UNSTAGE_FILE` | `lore_file_unstage_file_event_data_t` | Emitted for each file that was unstaged | */
func FileUnstage(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileUnstageArgsFFI,
) *LoreCall[types.LoreFileUnstageArgsFFI] {
	return &LoreCall[types.LoreFileUnstageArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileUnstage,
	}
}

/* Write binary content to a file in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_WRITE` | `lore_file_write_event_data_t` | Emitted when the file has been successfully written to the repository | */
func FileWrite(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileWriteArgsFFI,
) *LoreCall[types.LoreFileWriteArgsFFI] {
	return &LoreCall[types.LoreFileWriteArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileWrite,
	}
}

/* Permanently remove a file and all its history from the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_OBLITERATE` | `lore_file_obliterate_event_data_t` | Emitted for each file permanently removed from repository history | */
func FileObliterate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileObliterateArgsFFI,
) *LoreCall[types.LoreFileObliterateArgsFFI] {
	return &LoreCall[types.LoreFileObliterateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileObliterate,
	}
}

/* Retrieve the binary content of a file at a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DUMP` | `lore_file_dump_event_data_t` | Emitted with binary content of the requested file | */
func FileDump(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDumpArgsFFI,
) *LoreCall[types.LoreFileDumpArgsFFI] {
	return &LoreCall[types.LoreFileDumpArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDump,
	}
}

/* Adds dependency relationships between files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Dependency Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DEPENDENCY_ADD_BEGIN` | `lore_file_dependency_add_begin_event_data_t` | Start of operation |
| `LORE_EVENT_FILE_DEPENDENCY_ADD_ENTRY` | `lore_file_dependency_add_entry_event_data_t` | Each dependency added |
| `LORE_EVENT_FILE_DEPENDENCY_ADD_END` | `lore_file_dependency_add_end_event_data_t` | Operation complete | */
func FileDependencyAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDependencyAddArgsFFI,
) *LoreCall[types.LoreFileDependencyAddArgsFFI] {
	return &LoreCall[types.LoreFileDependencyAddArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDependencyAdd,
	}
}

/* Removes dependency relationships between files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Dependency Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DEPENDENCY_REMOVE_BEGIN` | `lore_file_dependency_remove_begin_event_data_t` | Start of operation |
| `LORE_EVENT_FILE_DEPENDENCY_REMOVE_ENTRY` | `lore_file_dependency_remove_entry_event_data_t` | Each dependency removed |
| `LORE_EVENT_FILE_DEPENDENCY_REMOVE_END` | `lore_file_dependency_remove_end_event_data_t` | Operation complete | */
func FileDependencyRemove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDependencyRemoveArgsFFI,
) *LoreCall[types.LoreFileDependencyRemoveArgsFFI] {
	return &LoreCall[types.LoreFileDependencyRemoveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDependencyRemove,
	}
}

/* Queries dependency information for files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Dependency Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DEPENDENCY_LIST_BEGIN` | `lore_file_dependency_list_begin_event_data_t` | Start of listing |
| `LORE_EVENT_FILE_DEPENDENCY_LIST_FILE` | `lore_file_dependency_list_file_event_data_t` | Start of entries for one file |
| `LORE_EVENT_FILE_DEPENDENCY_LIST_ENTRY` | `lore_file_dependency_list_entry_event_data_t` | One dependency entry |
| `LORE_EVENT_FILE_DEPENDENCY_LIST_FILE_END` | `lore_file_dependency_list_file_end_event_data_t` | End of entries for one file |
| `LORE_EVENT_FILE_DEPENDENCY_LIST_END` | `lore_file_dependency_list_end_event_data_t` | End of listing | */
func FileDependencyList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDependencyListArgsFFI,
) *LoreCall[types.LoreFileDependencyListArgsFFI] {
	return &LoreCall[types.LoreFileDependencyListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.FileDependencyList,
	}
}

/* Acquire exclusive locks on one or more files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_ACQUIRE` | `lore_lock_file_acquire_event_data_t` | Emitted for each file for which a lock was successfully acquired |
| `LORE_EVENT_LOCK_FILE_ACQUIRE_IGNORE` | `lore_lock_file_acquire_ignore_event_data_t` | Emitted for each file for which a lock was ignored (already owned) | */
func LockFileAcquire(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileAcquireArgsFFI,
) *LoreCall[types.LoreLockFileAcquireArgsFFI] {
	return &LoreCall[types.LoreLockFileAcquireArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LockFileAcquire,
	}
}

/* Get the lock status of files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_STATUS_BEGIN` | `lore_lock_file_status_begin_event_data_t` | Emitted before lock status results begin streaming |
| `LORE_EVENT_LOCK_FILE_STATUS` | `lore_lock_file_status_event_data_t` | Emitted for each locked file with owner and lock details | */
func LockFileStatus(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileStatusArgsFFI,
) *LoreCall[types.LoreLockFileStatusArgsFFI] {
	return &LoreCall[types.LoreLockFileStatusArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LockFileStatus,
	}
}

/* Query which files are currently locked, optionally filtered by user or path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_QUERY_BEGIN` | `lore_lock_file_query_begin_event_data_t` | Emitted before query results begin streaming |
| `LORE_EVENT_LOCK_FILE_QUERY` | `lore_lock_file_query_event_data_t` | Emitted for each file matching the query | */
func LockFileQuery(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileQueryArgsFFI,
) *LoreCall[types.LoreLockFileQueryArgsFFI] {
	return &LoreCall[types.LoreLockFileQueryArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LockFileQuery,
	}
}

/* Release file locks previously acquired by this client.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_RELEASE` | `lore_lock_file_release_event_data_t` | Emitted for each file lock successfully released |
| `LORE_EVENT_LOCK_FILE_RELEASE_NOT_FOUND` | `lore_lock_file_release_not_found_event_data_t` | Emitted for each file whose lock was not found | */
func LockFileRelease(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileReleaseArgsFFI,
) *LoreCall[types.LoreLockFileReleaseArgsFFI] {
	return &LoreCall[types.LoreLockFileReleaseArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LockFileRelease,
	}
}

/* Add a link to another repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_CLONE_BEGIN` | `lore_repository_clone_begin_event_data_t` | Emitted when cloning a linked repository begins |
| `LORE_EVENT_REPOSITORY_CLONE_END` | `lore_repository_clone_end_event_data_t` | Emitted when cloning a linked repository completes |
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when the link has been added and saved | */
func LinkAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkAddArgsFFI,
) *LoreCall[types.LoreLinkAddArgsFFI] {
	return &LoreCall[types.LoreLinkAddArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LinkAdd,
	}
}

/* Remove a link to another repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when the link has been removed | */
func LinkRemove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkRemoveArgsFFI,
) *LoreCall[types.LoreLinkRemoveArgsFFI] {
	return &LoreCall[types.LoreLinkRemoveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LinkRemove,
	}
}

/* List all repository links configured in the current repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_ENTRY` | `lore_link_entry_event_data_t` | Emitted for each linked repository | */
func LinkList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkListArgsFFI,
) *LoreCall[types.LoreLinkListArgsFFI] {
	return &LoreCall[types.LoreLinkListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LinkList,
	}
}

/* Update properties of an existing repository link.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when a link property is updated or finalized | */
func LinkUpdate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkUpdateArgsFFI,
) *LoreCall[types.LoreLinkUpdateArgsFFI] {
	return &LoreCall[types.LoreLinkUpdateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LinkUpdate,
	}
}

/* Clone a remote repository to a local path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_CLONE_BEGIN` | `lore_repository_clone_begin_event_data_t` | Emitted when clone begins, includes remote URL and target path |
| `LORE_EVENT_REPOSITORY_CLONE_PROGRESS` | `lore_repository_clone_progress_event_data_t` | Emitted periodically during clone with progress data |
| `LORE_EVENT_REPOSITORY_CLONE_END` | `lore_repository_clone_end_event_data_t` | Emitted when clone completes successfully |
| `LORE_EVENT_REVISION_SYNC_TARGET` | `lore_revision_sync_target_event_data_t` | Emitted after resolving the target revision to sync during clone |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file written during initial sync |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted periodically during initial file sync |
| `LORE_EVENT_REVISION_SYNC_REVISION` | `lore_revision_sync_revision_event_data_t` | Emitted with the resulting revision |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by view filters |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for each fragment written to the local store | */
func RepositoryClone(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryCloneArgsFFI,
) *LoreCall[types.LoreRepositoryCloneArgsFFI] {
	return &LoreCall[types.LoreRepositoryCloneArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryClone,
	}
}

/* Retrieve metadata about the current repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_DATA` | `lore_repository_data_event_data_t` | Emitted with repository metadata (name, URL, branch info, etc.) | */
func RepositoryInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInfoArgsFFI,
) *LoreCall[types.LoreRepositoryInfoArgsFFI] {
	return &LoreCall[types.LoreRepositoryInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryInfo,
	}
}

/* Dump the internal state of the repository for diagnostic purposes.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_DUMP_BEGIN` | `lore_repository_dump_begin_event_data_t` | Emitted before dump output begins |
| `LORE_EVENT_REPOSITORY_DUMP_END` | `lore_repository_dump_end_event_data_t` | Emitted when dump completes |
| `LORE_EVENT_REPOSITORY_STATE_DUMP` | `lore_repository_state_dump_event_data_t` | Emitted with repository state summary |
| `LORE_EVENT_REPOSITORY_STATE_DUMP_NODE` | `lore_repository_state_dump_node_event_data_t` | Emitted for each node in the state tree | */
func RepositoryDump(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryDumpArgsFFI,
) *LoreCall[types.LoreRepositoryDumpArgsFFI] {
	return &LoreCall[types.LoreRepositoryDumpArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryDump,
	}
}

/* Create a new Lore repository on the remote server.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_CREATE` | `lore_repository_create_event_data_t` | Emitted when the repository has been successfully created | */
func RepositoryCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryCreateArgsFFI,
) *LoreCall[types.LoreRepositoryCreateArgsFFI] {
	return &LoreCall[types.LoreRepositoryCreateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryCreate,
	}
}

/* Flush pending repository state to persistent storage.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryFlush(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryFlushArgsFFI,
) *LoreCall[types.LoreRepositoryFlushArgsFFI] {
	return &LoreCall[types.LoreRepositoryFlushArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryFlush,
	}
}

/* Run garbage collection to reclaim unreferenced storage in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryGc(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryGcArgsFFI,
) *LoreCall[types.LoreRepositoryGcArgsFFI] {
	return &LoreCall[types.LoreRepositoryGcArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryGc,
	}
}

/* Release all cached store references for the given repository path.

Frees in-memory store data and releases file-backed store cache entries.
Any active repository contexts for this path remain valid, but once they
are dropped the stores will be freed. Subsequent opens will create fresh stores.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryRelease(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryReleaseArgsFFI,
) *LoreCall[types.LoreRepositoryReleaseArgsFFI] {
	return &LoreCall[types.LoreRepositoryReleaseArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryRelease,
	}
}

/* Add a new layer to the repository configuration.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Layer Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LAYER_ADD` | `lore_layer_add_event_data_t` | Emitted when a layer has been successfully added | */
func LayerAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerAddArgsFFI,
) *LoreCall[types.LoreLayerAddArgsFFI] {
	return &LoreCall[types.LoreLayerAddArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LayerAdd,
	}
}

/* Remove a layer from the repository configuration.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func LayerRemove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerRemoveArgsFFI,
) *LoreCall[types.LoreLayerRemoveArgsFFI] {
	return &LoreCall[types.LoreLayerRemoveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LayerRemove,
	}
}

/* List all layers configured in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Layer Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LAYER_ENTRY` | `lore_layer_entry_event_data_t` | Emitted for each layer configured in the repository | */
func LayerList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerListArgsFFI,
) *LoreCall[types.LoreLayerListArgsFFI] {
	return &LoreCall[types.LoreLayerListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.LayerList,
	}
}

/* List all repositories available on the remote server.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_LIST_ENTRY` | `lore_repository_list_entry_event_data_t` | Emitted for each repository found | */
func RepositoryList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryListArgsFFI,
) *LoreCall[types.LoreRepositoryListArgsFFI] {
	return &LoreCall[types.LoreRepositoryListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryList,
	}
}

/* Show the working directory status, including staged, dirty, and conflicted files.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_STATUS_REVISION` | `lore_repository_status_revision_event_data_t` | Emitted with current and staged revision info |
| `LORE_EVENT_REPOSITORY_STATUS_FILE` | `lore_repository_status_file_event_data_t` | Emitted for each file with pending changes, conflict status, or untracked status |
| `LORE_EVENT_PATH_IGNORE` | `lore_path_ignore_event_data_t` | Emitted for each path excluded by ignore rules | */
func RepositoryStatus(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryStatusArgsFFI,
) *LoreCall[types.LoreRepositoryStatusArgsFFI] {
	return &LoreCall[types.LoreRepositoryStatusArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryStatus,
	}
}

/* Query the repository's immutable fragment store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_STORE_IMMUTABLE_QUERY` | `lore_repository_store_immutable_query_event_data_t` | Emitted for each fragment entry found in the immutable store | */
func RepositoryStoreImmutableQuery(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryStoreImmutableQueryArgsFFI,
) *LoreCall[types.LoreRepositoryStoreImmutableQueryArgsFFI] {
	return &LoreCall[types.LoreRepositoryStoreImmutableQueryArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryStoreImmutableQuery,
	}
}

/* Verify the integrity of the repository's stored fragments.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_VERIFY_STATE_BEGIN` | `lore_repository_verify_state_begin_event_data_t` | Emitted when verify begins |
| `LORE_EVENT_REPOSITORY_VERIFY_STATE_END` | `lore_repository_verify_state_end_event_data_t` | Emitted when verify completes (success or with errors) |
| `LORE_EVENT_REPOSITORY_VERIFY_FRAGMENT` | `lore_repository_verify_fragment_event_data_t` | Emitted for each fragment verified in the local store |
| `LORE_EVENT_REPOSITORY_VERIFY_FRAGMENT_REMOTE` | `lore_repository_verify_fragment_remote_event_data_t` | Emitted for each fragment verified against the remote store | */
func RepositoryVerifyState(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryVerifyStateArgsFFI,
) *LoreCall[types.LoreRepositoryVerifyStateArgsFFI] {
	return &LoreCall[types.LoreRepositoryVerifyStateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryVerifyState,
	}
}

/* Commit staged files to create a new revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when commit begins fragmenting files |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted periodically during commit with file processing counts |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when commit file processing completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed revision details (hash, branch, parents) |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata entry of the committed revision |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for each fragment written or deduplicated | */
func RevisionCommit(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionCommitArgsFFI,
) *LoreCall[types.LoreRevisionCommitArgsFFI] {
	return &LoreCall[types.LoreRevisionCommitArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionCommit,
	}
}

/* Amend the most recent revision with updated metadata.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the amended revision details |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata entry of the amended revision | */
func RevisionAmend(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionAmendArgsFFI,
) *LoreCall[types.LoreRevisionAmendArgsFFI] {
	return &LoreCall[types.LoreRevisionAmendArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionAmend,
	}
}

/* Retrieve metadata and delta information about a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_INFO` | `lore_revision_info_event_data_t` | Emitted with revision metadata (hash, branch, parents, file count, etc.) |
| `LORE_EVENT_REVISION_INFO_DELTA` | `lore_revision_info_delta_event_data_t` | Emitted with delta information between revision and its parent (when delta=true) |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata key/value of the revision (when metadata=true) | */
func RevisionInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionInfoArgsFFI,
) *LoreCall[types.LoreRevisionInfoArgsFFI] {
	return &LoreCall[types.LoreRevisionInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionInfo,
	}
}

/* Show files that differ between two revisions.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_DIFF_FILE` | `lore_revision_diff_file_event_data_t` | Emitted for each file that differs between the two revisions |
| `LORE_EVENT_REVISION_RESOLVE` | `lore_revision_resolve_event_data_t` | Emitted when resolving a partial or numbered revision reference | */
func RevisionDiff(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionDiffArgsFFI,
) *LoreCall[types.LoreRevisionDiffArgsFFI] {
	return &LoreCall[types.LoreRevisionDiffArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionDiff,
	}
}

/* Find a revision by metadata or revision number.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_FIND` | `lore_revision_find_event_data_t` | Emitted when a matching revision is found (exact or partial match) | */
func RevisionFind(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionFindArgsFFI,
) *LoreCall[types.LoreRevisionFindArgsFFI] {
	return &LoreCall[types.LoreRevisionFindArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionFind,
	}
}

/* Retrieve the commit history of the current branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_HISTORY` | `lore_revision_history_event_data_t` | Emitted once with summary info before entries stream |
| `LORE_EVENT_REVISION_HISTORY_ENTRY` | `lore_revision_history_entry_event_data_t` | Emitted for each revision in the history | */
func RevisionHistory(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionHistoryArgsFFI,
) *LoreCall[types.LoreRevisionHistoryArgsFFI] {
	return &LoreCall[types.LoreRevisionHistoryArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionHistory,
	}
}

/* Restore the working directory to a previously committed revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_RESTORE_FILE_BEGIN` | `lore_revision_restore_file_begin_event_data_t` | Emitted when restore starts processing files |
| `LORE_EVENT_REVISION_RESTORE_FILE` | `lore_revision_restore_file_event_data_t` | Emitted for each file being restored |
| `LORE_EVENT_REVISION_RESTORE_FILE_END` | `lore_revision_restore_file_end_event_data_t` | Emitted when file processing completes |
| `LORE_EVENT_REVISION_RESTORE_FRAGMENT_BEGIN` | `lore_revision_restore_fragment_begin_event_data_t` | Emitted when fragment download begins for a file |
| `LORE_EVENT_REVISION_RESTORE_FRAGMENT_PROGRESS` | `lore_revision_restore_fragment_progress_event_data_t` | Emitted periodically during fragment download |
| `LORE_EVENT_REVISION_RESTORE_FRAGMENT_END` | `lore_revision_restore_fragment_end_event_data_t` | Emitted when fragment download completes |
| `LORE_EVENT_REVISION_RESTORE_REVISION` | `lore_revision_restore_revision_event_data_t` | Emitted with the restored revision details |
| `LORE_EVENT_REVISION_RESTORE_SYNC_BEGIN` | `lore_revision_restore_sync_begin_event_data_t` | Emitted when starting to apply the changes on the target state |
| `LORE_EVENT_REVISION_RESTORE_SYNC_END` | `lore_revision_restore_sync_end_event_data_t` | Emitted after applying the changes on the target state is complete |
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when auto-commit of restored revision starts |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted during auto-commit |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when auto-commit completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed restored revision |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during changes realization |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for metadata of the restored revision |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for fragments written during restore commit | */
func RevisionRestore(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRestoreArgsFFI,
) *LoreCall[types.LoreRevisionRestoreArgsFFI] {
	return &LoreCall[types.LoreRevisionRestoreArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRestore,
	}
}

/* Clear all metadata from the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA_CLEAR_REVISION` | `lore_metadata_clear_revision_event_data_t` | Emitted when metadata has been cleared for the current revision | */
func RevisionMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataClearArgsFFI,
) *LoreCall[types.LoreRevisionMetadataClearArgsFFI] {
	return &LoreCall[types.LoreRevisionMetadataClearArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionMetadataClear,
	}
}

/* Get a specific metadata key/value pair from the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted with the requested key/value for the revision | */
func RevisionMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataGetArgsFFI,
) *LoreCall[types.LoreRevisionMetadataGetArgsFFI] {
	return &LoreCall[types.LoreRevisionMetadataGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionMetadataGet,
	}
}

/* List all metadata key/value pairs associated with the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata key/value associated with the revision | */
func RevisionMetadataList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataListArgsFFI,
) *LoreCall[types.LoreRevisionMetadataListArgsFFI] {
	return &LoreCall[types.LoreRevisionMetadataListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionMetadataList,
	}
}

/* Set a metadata key/value pair on the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RevisionMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataSetArgsFFI,
) *LoreCall[types.LoreRevisionMetadataSetArgsFFI] {
	return &LoreCall[types.LoreRevisionMetadataSetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionMetadataSet,
	}
}

/* Synchronize the working directory to a target revision, optionally merging divergent branches.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Sync Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_SYNC_TARGET` | `lore_revision_sync_target_event_data_t` | Emitted once after resolving the target revision with source/target revision info, branch, and remote URL |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file deleted, modified, added, or merged during sync |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted periodically during file realization and once at completion with cumulative update/delete/automerge/conflict counts |
| `LORE_EVENT_REVISION_SYNC_REVISION` | `lore_revision_sync_revision_event_data_t` | Emitted once at the end with the resulting revision, branch, and merge/conflict flags |
| `LORE_EVENT_REVISION_RESOLVE` | `lore_revision_resolve_event_data_t` | Emitted when resolving a partial or numbered revision reference |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by view or ignore filters |
| `LORE_EVENT_BRANCH_MERGE_START_BEGIN` | `lore_branch_merge_start_begin_event_data_t` | Emitted when an auto-merge is initiated (diverged branches) |
| `LORE_EVENT_BRANCH_MERGE_START_END` | `lore_branch_merge_start_end_event_data_t` | Emitted when the auto-merge operation completes |
| `LORE_EVENT_BRANCH_MERGE_CONFLICT_FILE` | `lore_branch_merge_conflict_file_event_data_t` | Emitted for each file with an unresolved merge conflict |
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when auto-merge auto-commits (no conflicts) |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted during auto-commit |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when auto-commit completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed merge revision |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for metadata of the auto-merge commit |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for each fragment written during auto-merge commit |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged for deletion during merge realization | */
func RevisionSync(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionSyncArgsFFI,
) *LoreCall[types.LoreRevisionSyncArgsFFI] {
	return &LoreCall[types.LoreRevisionSyncArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionSync,
	}
}

/* Revert a revision, applying its inverse changes to the working tree.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_START_BEGIN` | `lore_revert_start_begin_event_data_t` | Emitted when revert begins, includes target revision info |
| `LORE_EVENT_REVERT_START_END` | `lore_revert_start_end_event_data_t` | Emitted when revert completes, includes conflict flag |
| `LORE_EVENT_REVERT_CONFLICT_FILE` | `lore_revert_conflict_file_event_data_t` | Emitted for each file with an unresolved revert conflict |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during apply_diff phase |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file modified during revert realization |
| `LORE_EVENT_FILE_STAGE_FILE` | `lore_file_stage_file_event_data_t` | Emitted for each file staged for deletion during revert |
| `LORE_EVENT_REVISION_COMMIT_BEGIN` | `lore_revision_commit_begin_event_data_t` | Emitted when auto-commit starts (no conflicts) |
| `LORE_EVENT_REVISION_COMMIT_PROGRESS` | `lore_revision_commit_progress_event_data_t` | Emitted during auto-commit |
| `LORE_EVENT_REVISION_COMMIT_END` | `lore_revision_commit_end_event_data_t` | Emitted when auto-commit completes |
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the committed revert revision |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for metadata of the auto-commit |
| `LORE_EVENT_FRAGMENT_WRITE` | `lore_fragment_write_event_data_t` | Emitted for fragments written during auto-commit | */
func RevisionRevert(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertArgsFFI,
) *LoreCall[types.LoreRevisionRevertArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevert,
	}
}

/* Abort an in-progress revert operation and restore the previous state.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_ABORT_BEGIN` | `lore_revert_abort_begin_event_data_t` | Emitted when revert abort begins |
| `LORE_EVENT_REVERT_ABORT_END` | `lore_revert_abort_end_event_data_t` | Emitted when revert abort completes |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization while reverting | */
func RevisionRevertAbort(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertAbortArgsFFI,
) *LoreCall[types.LoreRevisionRevertAbortArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertAbortArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertAbort,
	}
}

/* Mark conflicting files in a revert operation as unresolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_UNRESOLVE_FILE` | `lore_revert_unresolve_file_event_data_t` | Emitted for each file marked as unresolved |
| `LORE_EVENT_REVERT_UNRESOLVE_REVISION` | `lore_revert_unresolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertUnresolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertUnresolveArgsFFI,
) *LoreCall[types.LoreRevisionRevertUnresolveArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertUnresolveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertUnresolve,
	}
}

/* Restart a revert operation, re-materializing files with conflicts.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_CONFLICT_FILE` | `lore_revert_conflict_file_event_data_t` | Emitted for each file with a remaining revert conflict |
| `LORE_EVENT_REVISION_SYNC_PROGRESS` | `lore_revision_sync_progress_event_data_t` | Emitted during file realization |
| `LORE_EVENT_REVISION_SYNC_FILE` | `lore_revision_sync_file_event_data_t` | Emitted for each file re-materialized | */
func RevisionRevertRestart(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertRestartArgsFFI,
) *LoreCall[types.LoreRevisionRevertRestartArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertRestartArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertRestart,
	}
}

/* Resolve a revert conflict by marking conflicting files as resolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file marked as resolved |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveArgsFFI,
) *LoreCall[types.LoreRevisionRevertResolveArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertResolveArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertResolve,
	}
}

/* Resolve a revert conflict by accepting the "mine" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file resolved by keeping "mine" |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolveMine(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveMineArgsFFI,
) *LoreCall[types.LoreRevisionRevertResolveMineArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertResolveMineArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertResolveMine,
	}
}

/* Resolve a revert conflict by accepting the "theirs" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file resolved by keeping "theirs" |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolveTheirs(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveTheirsArgsFFI,
) *LoreCall[types.LoreRevisionRevertResolveTheirsArgsFFI] {
	return &LoreCall[types.LoreRevisionRevertResolveTheirsArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RevisionRevertResolveTheirs,
	}
}

/* Create a new shared store at the specified path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Shared Store Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_SHARED_STORE_CREATE` | `lore_shared_store_create_event_data_t` | Emitted on success after the shared store is created, carrying the path of the newly created store | */
func SharedStoreCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreCreateArgsFFI,
) *LoreCall[types.LoreSharedStoreCreateArgsFFI] {
	return &LoreCall[types.LoreSharedStoreCreateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.SharedStoreCreate,
	}
}

/* Retrieve the path of the configured default shared store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Shared Store Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_SHARED_STORE_INFO` | `lore_shared_store_info_event_data_t` | Emitted on success carrying the path of the configured default shared store | */
func SharedStoreInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreInfoArgsFFI,
) *LoreCall[types.LoreSharedStoreInfoArgsFFI] {
	return &LoreCall[types.LoreSharedStoreInfoArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.SharedStoreInfo,
	}
}

/* Set whether to automatically use the shared store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func SharedStoreSetUseAutomatically(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreSetUseAutomaticallyArgsFFI,
) *LoreCall[types.LoreSharedStoreSetUseAutomaticallyArgsFFI] {
	return &LoreCall[types.LoreSharedStoreSetUseAutomaticallyArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.SharedStoreSetUseAutomatically,
	}
}

/* Open a content-addressed storage handle.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_OPENED` | `lore_storage_opened_event_data_t` | Emitted on success carrying the opened handle id |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted on failure (invalid mode, invalid path, cache construction error) |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` on success, `status: 1` otherwise | */
func StorageOpen(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageOpenArgsFFI,
) *LoreCall[types.LoreStorageOpenArgsFFI] {
	return &LoreCall[types.LoreStorageOpenArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageOpen,
	}
}

/* Store one or more content-addressed buffers.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_PUT_ITEM_COMPLETE` | `lore_storage_put_item_complete_event_data_t` | Emitted once per input item — success or failure |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Aggregate error when any item failed |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StoragePut(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStoragePutArgsFFI,
) *LoreCall[types.LoreStoragePutArgsFFI] {
	return &LoreCall[types.LoreStoragePutArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StoragePut,
	}
}

/* Read one or more content-addressed buffers.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_GET_HEADER` | `lore_storage_get_header_event_data_t` | Size of the item's reassembled content, emitted before any DATA events |
| `LORE_EVENT_STORAGE_GET_DATA` | `lore_storage_get_data_event_data_t` | Payload bytes — valid only during the callback invocation |
| `LORE_EVENT_STORAGE_GET_ITEM_COMPLETE` | `lore_storage_get_item_complete_event_data_t` | Terminal per-item event |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetArgsFFI,
) *LoreCall[types.LoreStorageGetArgsFFI] {
	return &LoreCall[types.LoreStorageGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageGet,
	}
}

/* Release a content-addressed storage handle.

Subsequent calls against the same handle return `InvalidArguments`.
Close does not block on the flush it spawns — `Complete` fires after
the in-flight counter drains, not after the flush finishes. */
func StorageClose(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageCloseArgsFFI,
) *LoreCall[types.LoreStorageCloseArgsFFI] {
	return &LoreCall[types.LoreStorageCloseArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageClose,
	}
}

/* Flush pending writes through the handle's stores.

On disk-backed stores this performs an fsync honoring `globals.sync_data`.
On in-memory stores the underlying flush is a no-op and the call still
completes with `status: 0`. */
func StorageFlush(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageFlushArgsFFI,
) *LoreCall[types.LoreStorageFlushArgsFFI] {
	return &LoreCall[types.LoreStorageFlushArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageFlush,
	}
}

/* Fetch fragment metadata for one or more `(partition, address)` pairs without paying the
payload bytes. Each item's terminal event carries the resolved `Fragment` (`flags`,
`size_payload`, `size_content`); on miss `error_code == ADDRESS_NOT_FOUND`.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_GET_METADATA_ITEM_COMPLETE` | `lore_storage_get_metadata_item_complete_event_data_t` | Per-item terminal event |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageGetMetadata(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetMetadataArgsFFI,
) *LoreCall[types.LoreStorageGetMetadataArgsFFI] {
	return &LoreCall[types.LoreStorageGetMetadataArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageGetMetadata,
	}
}

/* Delete one or more `(partition, address)` entries from the store.

Idempotent on absent items; emits one `OBLITERATE_ITEM_COMPLETE` event
per item carrying `local_success` / `remote_success` / `error_code`. */
func StorageObliterate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageObliterateArgsFFI,
) *LoreCall[types.LoreStorageObliterateArgsFFI] {
	return &LoreCall[types.LoreStorageObliterateArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageObliterate,
	}
}

/* Copy content from one partition to another within the same store.

Same-partition source/target rejects with `INVALID_ARGUMENTS`. The
item's content hash is preserved; only the source address is carried
in the per-item event. */
func StorageCopy(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageCopyArgsFFI,
) *LoreCall[types.LoreStorageCopyArgsFFI] {
	return &LoreCall[types.LoreStorageCopyArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageCopy,
	}
}

/* Read one or more files into the content-addressed store.

Each item emits `LORE_EVENT_STORAGE_PUT_ITEM_COMPLETE` carrying the
computed address. Empty files short-circuit to the zero-hash address
without opening for read. */
func StoragePutFile(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStoragePutFileArgsFFI,
) *LoreCall[types.LoreStoragePutFileArgsFFI] {
	return &LoreCall[types.LoreStoragePutFileArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StoragePutFile,
	}
}

/* Write content-addressed payloads to filesystem paths.

Each item emits `LORE_EVENT_STORAGE_GET_ITEM_COMPLETE`. No HEADER or
DATA events are produced — the payload is written straight to disk.
On partial-write failure the library leaves whatever state the
failure produced; cleanup is the caller's responsibility. */
func StorageGetFile(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetFileArgsFFI,
) *LoreCall[types.LoreStorageGetFileArgsFFI] {
	return &LoreCall[types.LoreStorageGetFileArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageGetFile,
	}
}

/* Push locally-stored, not-yet-durable content to the remote store.

Whole-call pre-dispatch fails when the handle has no remote, when `globals.offline=1`,
or when `globals.local=1`. Per-item: `partition == 0` → `INVALID_ARGUMENTS`; zero hash and
already-durable both succeed with `already_durable=1` and no remote call; missing local
payload → `ADDRESS_NOT_FOUND`. Otherwise the bytes are uploaded and the local entry is
updated with `PayloadStoredDurable` set. */
func StorageUpload(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageUploadArgsFFI,
) *LoreCall[types.LoreStorageUploadArgsFFI] {
	return &LoreCall[types.LoreStorageUploadArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.StorageUpload,
	}
}

/* Start the Lore background service.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func ServiceStart(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreServiceStartArgsFFI,
) *LoreCall[types.LoreServiceStartArgsFFI] {
	return &LoreCall[types.LoreServiceStartArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.ServiceStart,
	}
}

/* Stop the Lore background service.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func ServiceStop(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreServiceStopArgsFFI,
) *LoreCall[types.LoreServiceStopArgsFFI] {
	return &LoreCall[types.LoreServiceStopArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.ServiceStop,
	}
}

/* Subscribe to repository notifications.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Notification Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_NOTIFICATION_SUBSCRIBED` | `lore_notification_subscribed_event_data_t` | Emitted when successfully subscribed to repository notifications |
| `LORE_EVENT_NOTIFICATION_BRANCH_CREATED` | `lore_notification_branch_created_event_data_t` | Emitted when a branch is created in the repository (push notification) |
| `LORE_EVENT_NOTIFICATION_BRANCH_DELETED` | `lore_notification_branch_deleted_event_data_t` | Emitted when a branch is deleted in the repository (push notification) |
| `LORE_EVENT_NOTIFICATION_BRANCH_PUSHED` | `lore_notification_branch_pushed_event_data_t` | Emitted when a branch is pushed to (push notification) |
| `LORE_EVENT_NOTIFICATION_RESOURCE_LOCKED` | `lore_notification_resource_locked_event_data_t` | Emitted when a resource is locked (push notification) |
| `LORE_EVENT_NOTIFICATION_RESOURCE_UNLOCKED` | `lore_notification_resource_unlocked_event_data_t` | Emitted when a resource is unlocked (push notification) | */
func NotificationSubscribe(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreNotificationSubscribeArgsFFI,
) *LoreCall[types.LoreNotificationSubscribeArgsFFI] {
	return &LoreCall[types.LoreNotificationSubscribeArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.NotificationSubscribe,
	}
}

/* Unsubscribe from repository notifications.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted when an error occurs |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end (`status: 0` success, `status: 1` failure) |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Notification Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_NOTIFICATION_UNSUBSCRIBED` | `lore_notification_unsubscribed_event_data_t` | Emitted when successfully unsubscribed from repository notifications | */
func NotificationUnsubscribe(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreNotificationUnsubscribeArgsFFI,
) *LoreCall[types.LoreNotificationUnsubscribeArgsFFI] {
	return &LoreCall[types.LoreNotificationUnsubscribeArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.NotificationUnsubscribe,
	}
}

/* Retrieve repository metadata. Reads a single key, or all entries when no
key is given. */
func RepositoryMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataGetArgsFFI,
) *LoreCall[types.LoreRepositoryMetadataGetArgsFFI] {
	return &LoreCall[types.LoreRepositoryMetadataGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryMetadataGet,
	}
}

/* Set repository metadata key-value pairs. */
func RepositoryMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataSetArgsFFI,
) *LoreCall[types.LoreRepositoryMetadataSetArgsFFI] {
	return &LoreCall[types.LoreRepositoryMetadataSetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryMetadataSet,
	}
}

/* Clear repository metadata keys. Clears all user-defined keys when none are
given. */
func RepositoryMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataClearArgsFFI,
) *LoreCall[types.LoreRepositoryMetadataClearArgsFFI] {
	return &LoreCall[types.LoreRepositoryMetadataClearArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryMetadataClear,
	}
}

/* List the tracked instances of the repository. */
func RepositoryInstanceList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInstanceListArgsFFI,
) *LoreCall[types.LoreRepositoryInstanceListArgsFFI] {
	return &LoreCall[types.LoreRepositoryInstanceListArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryInstanceList,
	}
}

/* Remove stale instances of the repository that are no longer present. */
func RepositoryInstancePrune(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInstancePruneArgsFFI,
) *LoreCall[types.LoreRepositoryInstancePruneArgsFFI] {
	return &LoreCall[types.LoreRepositoryInstancePruneArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryInstancePrune,
	}
}

/* Update the recorded path of the current repository instance to its present
location. */
func RepositoryUpdatePath(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryUpdatePathArgsFFI,
) *LoreCall[types.LoreRepositoryUpdatePathArgsFFI] {
	return &LoreCall[types.LoreRepositoryUpdatePathArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryUpdatePath,
	}
}

/* Read a configuration value of the current repository by key. */
func RepositoryConfigGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryConfigGetArgsFFI,
) *LoreCall[types.LoreRepositoryConfigGetArgsFFI] {
	return &LoreCall[types.LoreRepositoryConfigGetArgsFFI]{
		globals:  globals,
		args:     args,
		execFunc: native.RepositoryConfigGet,
	}
}

// LogConfigure configures Lore logging.
func LogConfigure(logConfig *types.LoreLogConfigFFI) (int32, error) {
	returnCode, err := native.LogConfigure(logConfig)
	if err != nil {
		return returnCode, err
	}
	if returnCode != 0 {
		return returnCode, &LoreError{ReturnCode: returnCode}
	}
	return returnCode, nil
}

// Shutdown shuts down the Lore library.
func Shutdown() (int32, error) {
	returnCode, err := native.Shutdown()
	if err != nil {
		return returnCode, err
	}
	if returnCode != 0 {
		return returnCode, &LoreError{ReturnCode: returnCode}
	}
	return returnCode, nil
}

// SetThreadLimit limits the total number of threads Lore sizes its pools for.
// Must be called before the first Lore operation. Returns 0 if the limit was
// applied, or 1 if it had already been set (or the runtime was already
// running) — the latter is not treated as an error.
func SetThreadLimit(count uintptr) (int32, error) {
	return native.SetThreadLimit(count)
}

// Version returns the Lore library version string.
func Version() (string, error) {
	return native.Version()
}
