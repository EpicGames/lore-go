// Copyright Epic Games, Inc. All Rights Reserved.

package native

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/EpicGames/lore-go/types"
	"github.com/ebitengine/purego"
)

// Library handle and function pointers
var (
	libHandle                              uintptr
	libErr                                 error
	libOnce                                sync.Once
	loreLogConfigureFunc                   func(logConfigPtr uintptr) int32
	loreShutdownFunc                       func() int32
	loreSetThreadLimitFunc                 func(count uintptr) int32
	loreVersionFunc                        func() uintptr
	loreAuthUserInfoFunc                   loreFuncWithCallback
	loreAuthLoginWithTokenFunc             loreFuncWithCallback
	loreAuthListFunc                       loreFuncWithCallback
	loreAuthLogoutFunc                     loreFuncWithCallback
	loreAuthClearFunc                      loreFuncWithCallback
	loreAuthLocalUserInfoFunc              loreFuncWithCallback
	loreAuthLoginInteractiveFunc           loreFuncWithCallback
	loreBranchCreateFunc                   loreFuncWithCallback
	loreBranchInfoFunc                     loreFuncWithCallback
	loreBranchDiffFunc                     loreFuncWithCallback
	loreBranchProtectFunc                  loreFuncWithCallback
	loreBranchUnprotectFunc                loreFuncWithCallback
	loreBranchArchiveFunc                  loreFuncWithCallback
	loreBranchListFunc                     loreFuncWithCallback
	loreBranchMergeAbortFunc               loreFuncWithCallback
	loreBranchMergeUnresolveFunc           loreFuncWithCallback
	loreBranchMergeIntoFunc                loreFuncWithCallback
	loreBranchMergeResolveFunc             loreFuncWithCallback
	loreBranchMergeResolveMineFunc         loreFuncWithCallback
	loreBranchMergeResolveTheirsFunc       loreFuncWithCallback
	loreBranchMergeRestartFunc             loreFuncWithCallback
	loreBranchMergeStartFunc               loreFuncWithCallback
	loreBranchSwitchFunc                   loreFuncWithCallback
	loreBranchResetFunc                    loreFuncWithCallback
	loreBranchPushFunc                     loreFuncWithCallback
	loreBranchMetadataGetFunc              loreFuncWithCallback
	loreBranchMetadataSetFunc              loreFuncWithCallback
	loreBranchMetadataClearFunc            loreFuncWithCallback
	loreFileInfoFunc                       loreFuncWithCallback
	loreFileDiffFunc                       loreFuncWithCallback
	loreFileHashFunc                       loreFuncWithCallback
	loreFileHistoryFunc                    loreFuncWithCallback
	loreFileMetadataClearFunc              loreFuncWithCallback
	loreFileMetadataGetFunc                loreFuncWithCallback
	loreFileMetadataListFunc               loreFuncWithCallback
	loreFileMetadataSetFunc                loreFuncWithCallback
	loreFileResetFunc                      loreFuncWithCallback
	loreFileResetToLastMergedFunc          loreFuncWithCallback
	loreFileStageFunc                      loreFuncWithCallback
	loreFileStageMergeFunc                 loreFuncWithCallback
	loreFileStageMoveFunc                  loreFuncWithCallback
	loreFileDirtyFunc                      loreFuncWithCallback
	loreFileDirtyMoveFunc                  loreFuncWithCallback
	loreFileDirtyCopyFunc                  loreFuncWithCallback
	loreFileUnstageFunc                    loreFuncWithCallback
	loreFileWriteFunc                      loreFuncWithCallback
	loreFileObliterateFunc                 loreFuncWithCallback
	loreFileDumpFunc                       loreFuncWithCallback
	loreFileDependencyAddFunc              loreFuncWithCallback
	loreFileDependencyRemoveFunc           loreFuncWithCallback
	loreFileDependencyListFunc             loreFuncWithCallback
	loreLockFileAcquireFunc                loreFuncWithCallback
	loreLockFileStatusFunc                 loreFuncWithCallback
	loreLockFileQueryFunc                  loreFuncWithCallback
	loreLockFileReleaseFunc                loreFuncWithCallback
	loreLinkAddFunc                        loreFuncWithCallback
	loreLinkRemoveFunc                     loreFuncWithCallback
	loreLinkInfoFunc                       loreFuncWithCallback
	loreLinkListFunc                       loreFuncWithCallback
	loreLinkUpdateFunc                     loreFuncWithCallback
	loreRepositoryCloneFunc                loreFuncWithCallback
	loreRepositoryInfoFunc                 loreFuncWithCallback
	loreRepositoryDumpFunc                 loreFuncWithCallback
	loreRepositoryCreateFunc               loreFuncWithCallback
	loreRepositoryFlushFunc                loreFuncWithCallback
	loreRepositoryGcFunc                   loreFuncWithCallback
	loreRepositoryReleaseFunc              loreFuncWithCallback
	loreLayerAddFunc                       loreFuncWithCallback
	loreLayerRemoveFunc                    loreFuncWithCallback
	loreLayerListFunc                      loreFuncWithCallback
	loreRepositoryListFunc                 loreFuncWithCallback
	loreRepositoryStatusFunc               loreFuncWithCallback
	loreRepositoryStoreImmutableQueryFunc  loreFuncWithCallback
	loreRepositoryVerifyStateFunc          loreFuncWithCallback
	loreRevisionCommitFunc                 loreFuncWithCallback
	loreRevisionAmendFunc                  loreFuncWithCallback
	loreRevisionInfoFunc                   loreFuncWithCallback
	loreRevisionDiffFunc                   loreFuncWithCallback
	loreRevisionFindFunc                   loreFuncWithCallback
	loreRevisionHistoryFunc                loreFuncWithCallback
	loreRevisionRestoreFunc                loreFuncWithCallback
	loreRevisionMetadataClearFunc          loreFuncWithCallback
	loreRevisionMetadataGetFunc            loreFuncWithCallback
	loreRevisionMetadataListFunc           loreFuncWithCallback
	loreRevisionMetadataSetFunc            loreFuncWithCallback
	loreRevisionSyncFunc                   loreFuncWithCallback
	loreRevisionRevertFunc                 loreFuncWithCallback
	loreRevisionRevertAbortFunc            loreFuncWithCallback
	loreRevisionRevertUnresolveFunc        loreFuncWithCallback
	loreRevisionRevertRestartFunc          loreFuncWithCallback
	loreRevisionRevertResolveFunc          loreFuncWithCallback
	loreRevisionRevertResolveMineFunc      loreFuncWithCallback
	loreRevisionRevertResolveTheirsFunc    loreFuncWithCallback
	loreSharedStoreCreateFunc              loreFuncWithCallback
	loreSharedStoreInfoFunc                loreFuncWithCallback
	loreSharedStoreSetUseAutomaticallyFunc loreFuncWithCallback
	loreStorageOpenFunc                    loreFuncWithCallback
	loreStoragePutFunc                     loreFuncWithCallback
	loreStorageGetFunc                     loreFuncWithCallback
	loreStorageGetResolvedFunc             loreFuncWithCallback
	loreStoragePutResolvedFunc             loreFuncWithCallback
	loreStorageCloseFunc                   loreFuncWithCallback
	loreStorageFlushFunc                   loreFuncWithCallback
	loreStorageGetMetadataFunc             loreFuncWithCallback
	loreStorageObliterateFunc              loreFuncWithCallback
	loreStorageMutableLoadFunc             loreFuncWithCallback
	loreStorageMutableStoreFunc            loreFuncWithCallback
	loreStorageMutableCompareAndSwapFunc   loreFuncWithCallback
	loreStorageMutableListFunc             loreFuncWithCallback
	loreStorageCopyFunc                    loreFuncWithCallback
	loreStoragePutFileFunc                 loreFuncWithCallback
	loreStorageGetFileFunc                 loreFuncWithCallback
	loreStorageUploadFunc                  loreFuncWithCallback
	loreServiceStartFunc                   loreFuncWithCallback
	loreServiceStopFunc                    loreFuncWithCallback
	loreNotificationSubscribeFunc          loreFuncWithCallback
	loreNotificationUnsubscribeFunc        loreFuncWithCallback
	loreRepositoryMetadataGetFunc          loreFuncWithCallback
	loreRepositoryMetadataSetFunc          loreFuncWithCallback
	loreRepositoryMetadataClearFunc        loreFuncWithCallback
	loreRepositoryInstanceListFunc         loreFuncWithCallback
	loreRepositoryInstancePruneFunc        loreFuncWithCallback
	loreRepositoryUpdatePathFunc           loreFuncWithCallback
	loreRepositoryConfigGetFunc            loreFuncWithCallback
	loreRevisionTreeLoadFunc               loreFuncWithCallback
	loreRevisionTreeCloseFunc              loreFuncWithCallback
	loreRevisionTreeResolvePathFunc        loreFuncWithCallback
	loreRevisionTreeListChildrenFunc       loreFuncWithCallback
	loreRevisionTreeNodeInfoFunc           loreFuncWithCallback
	loreRevisionTreeInfoFunc               loreFuncWithCallback
	loreRevisionTreeNodePathFunc           loreFuncWithCallback
	loreRevisionTreeAddFunc                loreFuncWithCallback
	loreRevisionTreeDeleteFunc             loreFuncWithCallback
	loreRevisionTreeModifyFunc             loreFuncWithCallback
	loreRevisionTreeMoveFunc               loreFuncWithCallback
	loreRevisionTreeMetadataSetFunc        loreFuncWithCallback
	loreRevisionTreeMetadataGetFunc        loreFuncWithCallback
	loreRevisionTreeMetadataClearFunc      loreFuncWithCallback
	loreRevisionTreeCommitFunc             loreFuncWithCallback
)

// ensureLibrary loads the native library on first call. Concurrent callers
// block until loading completes; subsequent calls are a single atomic read.
func ensureLibrary() error {
	libOnce.Do(func() { libErr = initLibrary() })
	return libErr
}

// initLibrary loads the Lore library and function pointers.
// Search order:
//  1. LORE_LIB_PATH env var (use this for go run / go test)
//  2. Next to the compiled executable (for deployed binaries)
//  3. Next to the resolved executable (follows symlinks, for deployed binaries invoked via symlink)
func initLibrary() error {
	var searchPaths []string

	// 1. Check environment variable first
	if envPath := os.Getenv("LORE_LIB_PATH"); envPath != "" {
		searchPaths = append(searchPaths, envPath)
	}

	// 2. Check next to the compiled executable
	exePath, exeErr := os.Executable()
	if exeErr == nil {
		searchPaths = append(searchPaths, filepath.Join(filepath.Dir(exePath), nativeLibraryFileName()))
	}

	// 3. Check next to the resolved executable (follows symlinks)
	if exeErr == nil {
		if resolved, err := filepath.EvalSymlinks(exePath); err == nil && resolved != exePath {
			searchPaths = append(searchPaths, filepath.Join(filepath.Dir(resolved), nativeLibraryFileName()))
		}
	}

	if len(searchPaths) == 0 {
		return fmt.Errorf("no library path found for platform %s", runtime.GOOS)
	}

	var lastErr error
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		libHandle, lastErr = loadLibrary(path)
		if lastErr == nil {
			break
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to load library from any search path: %w", lastErr)
	}
	if libHandle == 0 {
		return fmt.Errorf("native library %s not found in any search path: %v", nativeLibraryFileName(), searchPaths)
	}

	purego.RegisterLibFunc(&loreAuthUserInfoFunc, libHandle, "lore_auth_user_info")
	purego.RegisterLibFunc(&loreAuthLoginWithTokenFunc, libHandle, "lore_auth_login_with_token")
	purego.RegisterLibFunc(&loreAuthListFunc, libHandle, "lore_auth_list")
	purego.RegisterLibFunc(&loreAuthLogoutFunc, libHandle, "lore_auth_logout")
	purego.RegisterLibFunc(&loreAuthClearFunc, libHandle, "lore_auth_clear")
	purego.RegisterLibFunc(&loreAuthLocalUserInfoFunc, libHandle, "lore_auth_local_user_info")
	purego.RegisterLibFunc(&loreAuthLoginInteractiveFunc, libHandle, "lore_auth_login_interactive")
	purego.RegisterLibFunc(&loreBranchCreateFunc, libHandle, "lore_branch_create")
	purego.RegisterLibFunc(&loreBranchInfoFunc, libHandle, "lore_branch_info")
	purego.RegisterLibFunc(&loreBranchDiffFunc, libHandle, "lore_branch_diff")
	purego.RegisterLibFunc(&loreBranchProtectFunc, libHandle, "lore_branch_protect")
	purego.RegisterLibFunc(&loreBranchUnprotectFunc, libHandle, "lore_branch_unprotect")
	purego.RegisterLibFunc(&loreBranchArchiveFunc, libHandle, "lore_branch_archive")
	purego.RegisterLibFunc(&loreBranchListFunc, libHandle, "lore_branch_list")
	purego.RegisterLibFunc(&loreBranchMergeAbortFunc, libHandle, "lore_branch_merge_abort")
	purego.RegisterLibFunc(&loreBranchMergeUnresolveFunc, libHandle, "lore_branch_merge_unresolve")
	purego.RegisterLibFunc(&loreBranchMergeIntoFunc, libHandle, "lore_branch_merge_into")
	purego.RegisterLibFunc(&loreBranchMergeResolveFunc, libHandle, "lore_branch_merge_resolve")
	purego.RegisterLibFunc(&loreBranchMergeResolveMineFunc, libHandle, "lore_branch_merge_resolve_mine")
	purego.RegisterLibFunc(&loreBranchMergeResolveTheirsFunc, libHandle, "lore_branch_merge_resolve_theirs")
	purego.RegisterLibFunc(&loreBranchMergeRestartFunc, libHandle, "lore_branch_merge_restart")
	purego.RegisterLibFunc(&loreBranchMergeStartFunc, libHandle, "lore_branch_merge_start")
	purego.RegisterLibFunc(&loreBranchSwitchFunc, libHandle, "lore_branch_switch")
	purego.RegisterLibFunc(&loreBranchResetFunc, libHandle, "lore_branch_reset")
	purego.RegisterLibFunc(&loreBranchPushFunc, libHandle, "lore_branch_push")
	purego.RegisterLibFunc(&loreBranchMetadataGetFunc, libHandle, "lore_branch_metadata_get")
	purego.RegisterLibFunc(&loreBranchMetadataSetFunc, libHandle, "lore_branch_metadata_set")
	purego.RegisterLibFunc(&loreBranchMetadataClearFunc, libHandle, "lore_branch_metadata_clear")
	purego.RegisterLibFunc(&loreFileInfoFunc, libHandle, "lore_file_info")
	purego.RegisterLibFunc(&loreFileDiffFunc, libHandle, "lore_file_diff")
	purego.RegisterLibFunc(&loreFileHashFunc, libHandle, "lore_file_hash")
	purego.RegisterLibFunc(&loreFileHistoryFunc, libHandle, "lore_file_history")
	purego.RegisterLibFunc(&loreFileMetadataClearFunc, libHandle, "lore_file_metadata_clear")
	purego.RegisterLibFunc(&loreFileMetadataGetFunc, libHandle, "lore_file_metadata_get")
	purego.RegisterLibFunc(&loreFileMetadataListFunc, libHandle, "lore_file_metadata_list")
	purego.RegisterLibFunc(&loreFileMetadataSetFunc, libHandle, "lore_file_metadata_set")
	purego.RegisterLibFunc(&loreFileResetFunc, libHandle, "lore_file_reset")
	purego.RegisterLibFunc(&loreFileResetToLastMergedFunc, libHandle, "lore_file_reset_to_last_merged")
	purego.RegisterLibFunc(&loreFileStageFunc, libHandle, "lore_file_stage")
	purego.RegisterLibFunc(&loreFileStageMergeFunc, libHandle, "lore_file_stage_merge")
	purego.RegisterLibFunc(&loreFileStageMoveFunc, libHandle, "lore_file_stage_move")
	purego.RegisterLibFunc(&loreFileDirtyFunc, libHandle, "lore_file_dirty")
	purego.RegisterLibFunc(&loreFileDirtyMoveFunc, libHandle, "lore_file_dirty_move")
	purego.RegisterLibFunc(&loreFileDirtyCopyFunc, libHandle, "lore_file_dirty_copy")
	purego.RegisterLibFunc(&loreFileUnstageFunc, libHandle, "lore_file_unstage")
	purego.RegisterLibFunc(&loreFileWriteFunc, libHandle, "lore_file_write")
	purego.RegisterLibFunc(&loreFileObliterateFunc, libHandle, "lore_file_obliterate")
	purego.RegisterLibFunc(&loreFileDumpFunc, libHandle, "lore_file_dump")
	purego.RegisterLibFunc(&loreFileDependencyAddFunc, libHandle, "lore_file_dependency_add")
	purego.RegisterLibFunc(&loreFileDependencyRemoveFunc, libHandle, "lore_file_dependency_remove")
	purego.RegisterLibFunc(&loreFileDependencyListFunc, libHandle, "lore_file_dependency_list")
	purego.RegisterLibFunc(&loreLockFileAcquireFunc, libHandle, "lore_lock_file_acquire")
	purego.RegisterLibFunc(&loreLockFileStatusFunc, libHandle, "lore_lock_file_status")
	purego.RegisterLibFunc(&loreLockFileQueryFunc, libHandle, "lore_lock_file_query")
	purego.RegisterLibFunc(&loreLockFileReleaseFunc, libHandle, "lore_lock_file_release")
	purego.RegisterLibFunc(&loreLinkAddFunc, libHandle, "lore_link_add")
	purego.RegisterLibFunc(&loreLinkRemoveFunc, libHandle, "lore_link_remove")
	purego.RegisterLibFunc(&loreLinkInfoFunc, libHandle, "lore_link_info")
	purego.RegisterLibFunc(&loreLinkListFunc, libHandle, "lore_link_list")
	purego.RegisterLibFunc(&loreLinkUpdateFunc, libHandle, "lore_link_update")
	purego.RegisterLibFunc(&loreRepositoryCloneFunc, libHandle, "lore_repository_clone")
	purego.RegisterLibFunc(&loreRepositoryInfoFunc, libHandle, "lore_repository_info")
	purego.RegisterLibFunc(&loreRepositoryDumpFunc, libHandle, "lore_repository_dump")
	purego.RegisterLibFunc(&loreRepositoryCreateFunc, libHandle, "lore_repository_create")
	purego.RegisterLibFunc(&loreRepositoryFlushFunc, libHandle, "lore_repository_flush")
	purego.RegisterLibFunc(&loreRepositoryGcFunc, libHandle, "lore_repository_gc")
	purego.RegisterLibFunc(&loreRepositoryReleaseFunc, libHandle, "lore_repository_release")
	purego.RegisterLibFunc(&loreLayerAddFunc, libHandle, "lore_layer_add")
	purego.RegisterLibFunc(&loreLayerRemoveFunc, libHandle, "lore_layer_remove")
	purego.RegisterLibFunc(&loreLayerListFunc, libHandle, "lore_layer_list")
	purego.RegisterLibFunc(&loreRepositoryListFunc, libHandle, "lore_repository_list")
	purego.RegisterLibFunc(&loreRepositoryStatusFunc, libHandle, "lore_repository_status")
	purego.RegisterLibFunc(&loreRepositoryStoreImmutableQueryFunc, libHandle, "lore_repository_store_immutable_query")
	purego.RegisterLibFunc(&loreRepositoryVerifyStateFunc, libHandle, "lore_repository_verify_state")
	purego.RegisterLibFunc(&loreRevisionCommitFunc, libHandle, "lore_revision_commit")
	purego.RegisterLibFunc(&loreRevisionAmendFunc, libHandle, "lore_revision_amend")
	purego.RegisterLibFunc(&loreRevisionInfoFunc, libHandle, "lore_revision_info")
	purego.RegisterLibFunc(&loreRevisionDiffFunc, libHandle, "lore_revision_diff")
	purego.RegisterLibFunc(&loreRevisionFindFunc, libHandle, "lore_revision_find")
	purego.RegisterLibFunc(&loreRevisionHistoryFunc, libHandle, "lore_revision_history")
	purego.RegisterLibFunc(&loreRevisionRestoreFunc, libHandle, "lore_revision_restore")
	purego.RegisterLibFunc(&loreRevisionMetadataClearFunc, libHandle, "lore_revision_metadata_clear")
	purego.RegisterLibFunc(&loreRevisionMetadataGetFunc, libHandle, "lore_revision_metadata_get")
	purego.RegisterLibFunc(&loreRevisionMetadataListFunc, libHandle, "lore_revision_metadata_list")
	purego.RegisterLibFunc(&loreRevisionMetadataSetFunc, libHandle, "lore_revision_metadata_set")
	purego.RegisterLibFunc(&loreRevisionSyncFunc, libHandle, "lore_revision_sync")
	purego.RegisterLibFunc(&loreRevisionRevertFunc, libHandle, "lore_revision_revert")
	purego.RegisterLibFunc(&loreRevisionRevertAbortFunc, libHandle, "lore_revision_revert_abort")
	purego.RegisterLibFunc(&loreRevisionRevertUnresolveFunc, libHandle, "lore_revision_revert_unresolve")
	purego.RegisterLibFunc(&loreRevisionRevertRestartFunc, libHandle, "lore_revision_revert_restart")
	purego.RegisterLibFunc(&loreRevisionRevertResolveFunc, libHandle, "lore_revision_revert_resolve")
	purego.RegisterLibFunc(&loreRevisionRevertResolveMineFunc, libHandle, "lore_revision_revert_resolve_mine")
	purego.RegisterLibFunc(&loreRevisionRevertResolveTheirsFunc, libHandle, "lore_revision_revert_resolve_theirs")
	purego.RegisterLibFunc(&loreSharedStoreCreateFunc, libHandle, "lore_shared_store_create")
	purego.RegisterLibFunc(&loreSharedStoreInfoFunc, libHandle, "lore_shared_store_info")
	purego.RegisterLibFunc(&loreSharedStoreSetUseAutomaticallyFunc, libHandle, "lore_shared_store_set_use_automatically")
	purego.RegisterLibFunc(&loreStorageOpenFunc, libHandle, "lore_storage_open")
	purego.RegisterLibFunc(&loreStoragePutFunc, libHandle, "lore_storage_put")
	purego.RegisterLibFunc(&loreStorageGetFunc, libHandle, "lore_storage_get")
	purego.RegisterLibFunc(&loreStorageGetResolvedFunc, libHandle, "lore_storage_get_resolved")
	purego.RegisterLibFunc(&loreStoragePutResolvedFunc, libHandle, "lore_storage_put_resolved")
	purego.RegisterLibFunc(&loreStorageCloseFunc, libHandle, "lore_storage_close")
	purego.RegisterLibFunc(&loreStorageFlushFunc, libHandle, "lore_storage_flush")
	purego.RegisterLibFunc(&loreStorageGetMetadataFunc, libHandle, "lore_storage_get_metadata")
	purego.RegisterLibFunc(&loreStorageObliterateFunc, libHandle, "lore_storage_obliterate")
	purego.RegisterLibFunc(&loreStorageMutableLoadFunc, libHandle, "lore_storage_mutable_load")
	purego.RegisterLibFunc(&loreStorageMutableStoreFunc, libHandle, "lore_storage_mutable_store")
	purego.RegisterLibFunc(&loreStorageMutableCompareAndSwapFunc, libHandle, "lore_storage_mutable_compare_and_swap")
	purego.RegisterLibFunc(&loreStorageMutableListFunc, libHandle, "lore_storage_mutable_list")
	purego.RegisterLibFunc(&loreStorageCopyFunc, libHandle, "lore_storage_copy")
	purego.RegisterLibFunc(&loreStoragePutFileFunc, libHandle, "lore_storage_put_file")
	purego.RegisterLibFunc(&loreStorageGetFileFunc, libHandle, "lore_storage_get_file")
	purego.RegisterLibFunc(&loreStorageUploadFunc, libHandle, "lore_storage_upload")
	purego.RegisterLibFunc(&loreServiceStartFunc, libHandle, "lore_service_start")
	purego.RegisterLibFunc(&loreServiceStopFunc, libHandle, "lore_service_stop")
	purego.RegisterLibFunc(&loreNotificationSubscribeFunc, libHandle, "lore_notification_subscribe")
	purego.RegisterLibFunc(&loreNotificationUnsubscribeFunc, libHandle, "lore_notification_unsubscribe")
	purego.RegisterLibFunc(&loreRepositoryMetadataGetFunc, libHandle, "lore_repository_metadata_get")
	purego.RegisterLibFunc(&loreRepositoryMetadataSetFunc, libHandle, "lore_repository_metadata_set")
	purego.RegisterLibFunc(&loreRepositoryMetadataClearFunc, libHandle, "lore_repository_metadata_clear")
	purego.RegisterLibFunc(&loreRepositoryInstanceListFunc, libHandle, "lore_repository_instance_list")
	purego.RegisterLibFunc(&loreRepositoryInstancePruneFunc, libHandle, "lore_repository_instance_prune")
	purego.RegisterLibFunc(&loreRepositoryUpdatePathFunc, libHandle, "lore_repository_update_path")
	purego.RegisterLibFunc(&loreRepositoryConfigGetFunc, libHandle, "lore_repository_config_get")
	purego.RegisterLibFunc(&loreRevisionTreeLoadFunc, libHandle, "lore_revision_tree_load")
	purego.RegisterLibFunc(&loreRevisionTreeCloseFunc, libHandle, "lore_revision_tree_close")
	purego.RegisterLibFunc(&loreRevisionTreeResolvePathFunc, libHandle, "lore_revision_tree_resolve_path")
	purego.RegisterLibFunc(&loreRevisionTreeListChildrenFunc, libHandle, "lore_revision_tree_list_children")
	purego.RegisterLibFunc(&loreRevisionTreeNodeInfoFunc, libHandle, "lore_revision_tree_node_info")
	purego.RegisterLibFunc(&loreRevisionTreeInfoFunc, libHandle, "lore_revision_tree_info")
	purego.RegisterLibFunc(&loreRevisionTreeNodePathFunc, libHandle, "lore_revision_tree_node_path")
	purego.RegisterLibFunc(&loreRevisionTreeAddFunc, libHandle, "lore_revision_tree_add")
	purego.RegisterLibFunc(&loreRevisionTreeDeleteFunc, libHandle, "lore_revision_tree_delete")
	purego.RegisterLibFunc(&loreRevisionTreeModifyFunc, libHandle, "lore_revision_tree_modify")
	purego.RegisterLibFunc(&loreRevisionTreeMoveFunc, libHandle, "lore_revision_tree_move")
	purego.RegisterLibFunc(&loreRevisionTreeMetadataSetFunc, libHandle, "lore_revision_tree_metadata_set")
	purego.RegisterLibFunc(&loreRevisionTreeMetadataGetFunc, libHandle, "lore_revision_tree_metadata_get")
	purego.RegisterLibFunc(&loreRevisionTreeMetadataClearFunc, libHandle, "lore_revision_tree_metadata_clear")
	purego.RegisterLibFunc(&loreRevisionTreeCommitFunc, libHandle, "lore_revision_tree_commit")

	purego.RegisterLibFunc(&loreLogConfigureFunc, libHandle, "lore_log_configure")
	purego.RegisterLibFunc(&loreShutdownFunc, libHandle, "lore_shutdown")
	purego.RegisterLibFunc(&loreSetThreadLimitFunc, libHandle, "lore_set_thread_limit")
	purego.RegisterLibFunc(&loreVersionFunc, libHandle, "lore_version")

	return nil
}

// Callback registry to keep Go functions alive while the C library is making
// callbacks for them. sync.Map is used (rather than map+mutex) because the
// trampoline reads from this on every event — 100k+ events per
// repository_dump call — and we want lock-free reads on the hot path. Writes
// (Store on registration, Delete on END) are rare and don't contend.
var (
	callbackRegistry sync.Map      // map[uint64]*callbackData
	nextCallbackID   atomic.Uint64 // monotonic, starts at 1 after first Add
)

// eventCallbackFuncPtr is the C-callable function pointer for our trampoline.
// purego.NewCallback consumes a slot in a fixed-size table (max 2000 across
// the whole process) and the slot is *never released*, so we must register
// exactly once at init — registering per-call would leak slots and panic
// after ~2000 lore_* invocations.
var eventCallbackFuncPtr uintptr

func init() {
	eventCallbackFuncPtr = purego.NewCallback(eventCallbackTrampoline)
}

// callbackData stores the user's callback and context
type callbackData struct {
	callback    types.LoreEventCallback
	userContext uint64
}

// eventCallbackTrampoline is the C-callable callback function
// This must match the C signature: void (*func)(const struct lore_event_t *event, uint64_t user_context)
// Note: Returns uintptr to satisfy purego.NewCallback requirements on Windows (C void functions return 0)
func eventCallbackTrampoline(eventPtr uintptr, callbackID uint64) uintptr {
	v, exists := callbackRegistry.Load(callbackID)
	if !exists {
		return 0
	}
	data := v.(*callbackData)

	// Cast the C event pointer to our Go type
	event := (*types.LoreEventFFI)(unsafe.Pointer(eventPtr))

	// Call the user's callback
	data.callback(event, data.userContext)

	// END signals callback termination — drop the registry entry so the
	// user's closure (and anything it captures) becomes GC-eligible.
	if event.Tag == types.LoreEventTag_END {
		callbackRegistry.Delete(callbackID)
	}

	return 0 // C void function returns 0
}

// callLoreFunction is a generic helper that handles the common boilerplate for all Lore function calls.
// It takes a pointer to the function variable (not the value) so the variable is read AFTER
// ensureLibrary() has loaded the library and registered all function pointers.
func callLoreFunction[TArgs any](
	loreFuncPtr *loreFuncWithCallback,
	globals *types.LoreGlobalArgsFFI,
	args *TArgs,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	if err := ensureLibrary(); err != nil {
		return -1, err
	}

	// Register the callback. atomic.Uint64.Add returns the post-increment
	// value, so the first call here gets ID=1 (matching the original
	// semantics where nextCallbackID started at 1).
	callbackID := nextCallbackID.Add(1)
	callbackRegistry.Store(callbackID, &callbackData{
		callback:    config.Callback,
		userContext: config.UserContext,
	})

	// Reuse the package-level trampoline pointer — see eventCallbackFuncPtr
	// comment for why this MUST be created once, not per call.
	cCallbackConfig := types.LoreEventCallbackConfigFFI{
		UserContext: callbackID,
		FuncPtr:     eventCallbackFuncPtr,
	}

	// Pin every struct the native side will read.
	var pinner runtime.Pinner
	defer pinner.Unpin()
	pinner.Pin(globals)
	pinner.Pin(args)
	pinner.Pin(&cCallbackConfig)

	// Convert all pointers to uintptr for FFI call
	globalsPtr := uintptr(unsafe.Pointer(globals))
	argsPtr := uintptr(unsafe.Pointer(args))
	callbackConfigPtr := uintptr(unsafe.Pointer(&cCallbackConfig))
	callbackUserContext := uintptr(cCallbackConfig.UserContext)
	callbackFuncPtrValue := cCallbackConfig.FuncPtr

	result := callLoreFuncWithCallback(
		*loreFuncPtr,
		globalsPtr,
		argsPtr,
		callbackConfigPtr,
		callbackUserContext,
		callbackFuncPtrValue,
	)

	return result, nil
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name for each resolved user | */
func AuthUserInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthUserInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthUserInfoFunc, globals, args, config)
}

/* Authenticate using an existing bearer token.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name after successful token authentication | */
func AuthLoginWithToken(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLoginWithTokenArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthLoginWithTokenFunc, globals, args, config)
}

/* List all stored authentication identities.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_IDENTITY` | `lore_auth_identity_event_data_t` | Emitted once per stored identity | */
func AuthList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthListFunc, globals, args, config)
}

/* Remove stored authentication and authorization tokens.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func AuthLogout(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLogoutArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthLogoutFunc, globals, args, config)
}

/* Clear all stored authentication data.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func AuthClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthClearFunc, globals, args, config)
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with the resolved user id and display name | */
func AuthLocalUserInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLocalUserInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthLocalUserInfoFunc, globals, args, config)
}

/* Authenticate interactively via a browser-based login flow.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Auth Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_AUTH_URL` | `lore_auth_url_event_data_t` | Emitted with the login URL when no_browser mode is requested |
| `LORE_EVENT_AUTH_USER_INFO` | `lore_auth_user_info_event_data_t` | Emitted with user id and display name after successful interactive authentication | */
func AuthLoginInteractive(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreAuthLoginInteractiveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreAuthLoginInteractiveFunc, globals, args, config)
}

/* Create a new branch in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_CREATE` | `lore_branch_create_event_data_t` | Emitted when the branch has been successfully created, includes branch name and id |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_BRANCH_CREATE` | `lore_link_branch_create_event_data_t` | Emitted once per linked repository mount, reporting whether its branch was created or an existing one reused | */
func BranchCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchCreateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchCreateFunc, globals, args, config)
}

/* Retrieve metadata about a specific branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_INFO` | `lore_branch_info_event_data_t` | Emitted with branch metadata (name, id, category, protection status, etc.) | */
func BranchInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchInfoFunc, globals, args, config)
}

/* Show the changes and conflicts between two branches.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchDiffFunc, globals, args, config)
}

/* Enable write protection on a branch to prevent direct commits.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_PROTECT` | `lore_branch_protect_event_data_t` | Emitted when the branch has been successfully protected | */
func BranchProtect(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchProtectArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchProtectFunc, globals, args, config)
}

/* Remove write protection from a branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_UNPROTECT` | `lore_branch_unprotect_event_data_t` | Emitted when the branch has been successfully unprotected | */
func BranchUnprotect(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchUnprotectArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchUnprotectFunc, globals, args, config)
}

/* Archive a branch in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_ARCHIVE` | `lore_branch_archive_event_data_t` | Emitted when the branch has been successfully archived | */
func BranchArchive(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchArchiveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchArchiveFunc, globals, args, config)
}

/* List all branches in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchListFunc, globals, args, config)
}

/* Abort an in-progress branch merge and restore the pre-merge state.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeAbortFunc, globals, args, config)
}

/* Mark conflicting files in a merge as unresolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_UNRESOLVE_FILE` | `lore_branch_merge_unresolve_file_event_data_t` | Emitted for each file that was marked as unresolved |
| `LORE_EVENT_BRANCH_MERGE_UNRESOLVE_REVISION` | `lore_branch_merge_unresolve_revision_event_data_t` | Emitted with the updated staged revision after unresolve completes | */
func BranchMergeUnresolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeUnresolveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeUnresolveFunc, globals, args, config)
}

/* Merge the current branch into a target branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeIntoFunc, globals, args, config)
}

/* Mark conflicting files in a merge as resolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file that was marked as resolved |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision after resolve completes | */
func BranchMergeResolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeResolveFunc, globals, args, config)
}

/* Resolve a merge conflict by accepting the "mine" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file resolved by keeping "mine" |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func BranchMergeResolveMine(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveMineArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeResolveMineFunc, globals, args, config)
}

/* Resolve a merge conflict by accepting the "theirs" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_FILE` | `lore_branch_merge_resolve_file_event_data_t` | Emitted for each file resolved by keeping "theirs" |
| `LORE_EVENT_BRANCH_MERGE_RESOLVE_REVISION` | `lore_branch_merge_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func BranchMergeResolveTheirs(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMergeResolveTheirsArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeResolveTheirsFunc, globals, args, config)
}

/* Restart an in-progress merge, re-materializing conflicted files.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeRestartFunc, globals, args, config)
}

/* Start a merge from another branch into the current branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMergeStartFunc, globals, args, config)
}

/* Switch to a different branch and update the working directory.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchSwitchFunc, globals, args, config)
}

/* Reset the current branch to a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Branch Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_BRANCH_RESET` | `lore_branch_reset_event_data_t` | Emitted when the branch has been reset to the target revision | */
func BranchReset(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchResetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchResetFunc, globals, args, config)
}

/* Push local branch commits to the remote repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchPushFunc, globals, args, config)
}

/* Retrieve branch metadata. */
func BranchMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMetadataGetFunc, globals, args, config)
}

/* Set branch metadata key-value pairs. */
func BranchMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataSetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMetadataSetFunc, globals, args, config)
}

/* Clear branch metadata keys. */
func BranchMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreBranchMetadataClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreBranchMetadataClearFunc, globals, args, config)
}

/* Retrieve metadata for one or more files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_INFO` | `lore_file_info_event_data_t` | Emitted for each file with its metadata (size, hash, staged status, etc.) | */
func FileInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileInfoFunc, globals, args, config)
}

/* Show which files differ between two revisions.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DIFF` | `lore_file_diff_event_data_t` | Emitted for each file that differs between the two revisions | */
func FileDiff(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDiffArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDiffFunc, globals, args, config)
}

/* Compute the hash of a local file for comparison with repository content.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_HASH` | `lore_file_hash_event_data_t` | Emitted with the computed hash and size of the specified file | */
func FileHash(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileHashArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileHashFunc, globals, args, config)
}

/* Retrieve the revision history for a specific file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_HISTORY` | `lore_file_history_event_data_t` | Emitted for each revision in which the file was modified | */
func FileHistory(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileHistoryArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileHistoryFunc, globals, args, config)
}

/* Clear all metadata from a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA_CLEAR_FILE` | `lore_metadata_clear_file_event_data_t` | Emitted when metadata has been cleared for the file | */
func FileMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileMetadataClearFunc, globals, args, config)
}

/* Get a specific metadata key/value pair from a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for the requested metadata key/value pair | */
func FileMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileMetadataGetFunc, globals, args, config)
}

/* List all metadata key/value pairs associated with a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata key/value pair associated with the file | */
func FileMetadataList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileMetadataListFunc, globals, args, config)
}

/* Set a metadata key/value pair on a file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileMetadataSetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileMetadataSetFunc, globals, args, config)
}

/* Reset files to the state recorded in the current or target revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileResetFunc, globals, args, config)
}

/* Reset files to their state at the last merged revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileResetToLastMergedFunc, globals, args, config)
}

/* Stage files for the next commit.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileStageFunc, globals, args, config)
}

/* Stage files for a merge commit, recording resolved merge content.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileStageMergeFunc, globals, args, config)
}

/* Stage a file move (rename) operation for commit.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileStageMoveFunc, globals, args, config)
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_PATH_IGNORE` | `lore_path_ignore_event_data_t` | Emitted for each input path that could not be resolved to a repository-relative path |
| `LORE_EVENT_FILTER_EXCLUDE` | `lore_filter_exclude_event_data_t` | Emitted for each path excluded by view or ignore filters | */
func FileDirty(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDirtyFunc, globals, args, config)
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileDirtyMove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyMoveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDirtyMoveFunc, globals, args, config)
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func FileDirtyCopy(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDirtyCopyArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDirtyCopyFunc, globals, args, config)
}

/* Remove files from the staging area without discarding local changes.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileUnstageFunc, globals, args, config)
}

/* Write binary content to a file in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_WRITE` | `lore_file_write_event_data_t` | Emitted when the file has been successfully written to the repository | */
func FileWrite(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileWriteArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileWriteFunc, globals, args, config)
}

/* Permanently remove a file and all its history from the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_OBLITERATE` | `lore_file_obliterate_event_data_t` | Emitted for each file permanently removed from repository history | */
func FileObliterate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileObliterateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileObliterateFunc, globals, args, config)
}

/* Retrieve the binary content of a file at a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## File Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_FILE_DUMP` | `lore_file_dump_event_data_t` | Emitted with binary content of the requested file | */
func FileDump(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreFileDumpArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDumpFunc, globals, args, config)
}

/* Adds dependency relationships between files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDependencyAddFunc, globals, args, config)
}

/* Removes dependency relationships between files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDependencyRemoveFunc, globals, args, config)
}

/* Queries dependency information for files.

# Events

## Standard Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreFileDependencyListFunc, globals, args, config)
}

/* Acquire exclusive locks on one or more files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_ACQUIRE` | `lore_lock_file_acquire_event_data_t` | Emitted for each file for which a lock was successfully acquired |
| `LORE_EVENT_LOCK_FILE_ACQUIRE_IGNORE` | `lore_lock_file_acquire_ignore_event_data_t` | Emitted for each file for which a lock was ignored (already owned) | */
func LockFileAcquire(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileAcquireArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLockFileAcquireFunc, globals, args, config)
}

/* Get the lock status of files in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_STATUS_BEGIN` | `lore_lock_file_status_begin_event_data_t` | Emitted before lock status results begin streaming |
| `LORE_EVENT_LOCK_FILE_STATUS` | `lore_lock_file_status_event_data_t` | Emitted for each locked file with owner and lock details | */
func LockFileStatus(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileStatusArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLockFileStatusFunc, globals, args, config)
}

/* Query which files are currently locked, optionally filtered by user or path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_QUERY_BEGIN` | `lore_lock_file_query_begin_event_data_t` | Emitted before query results begin streaming |
| `LORE_EVENT_LOCK_FILE_QUERY` | `lore_lock_file_query_event_data_t` | Emitted for each file matching the query | */
func LockFileQuery(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileQueryArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLockFileQueryFunc, globals, args, config)
}

/* Release file locks previously acquired by this client.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Lock Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOCK_FILE_RELEASE` | `lore_lock_file_release_event_data_t` | Emitted for each file lock successfully released |
| `LORE_EVENT_LOCK_FILE_RELEASE_NOT_FOUND` | `lore_lock_file_release_not_found_event_data_t` | Emitted for each file whose lock was not found | */
func LockFileRelease(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLockFileReleaseArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLockFileReleaseFunc, globals, args, config)
}

/* Add a link to another repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_CLONE_BEGIN` | `lore_repository_clone_begin_event_data_t` | Emitted when cloning a linked repository begins |
| `LORE_EVENT_REPOSITORY_CLONE_END` | `lore_repository_clone_end_event_data_t` | Emitted when cloning a linked repository completes |
| `LORE_EVENT_LINK_BRANCH_CREATE` | `lore_link_branch_create_event_data_t` | Emitted when branching is enabled, reporting whether the link's branch was created or an existing one reused |
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when the link has been added and saved | */
func LinkAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkAddArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLinkAddFunc, globals, args, config)
}

/* Remove a link to another repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when the link has been removed | */
func LinkRemove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkRemoveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLinkRemoveFunc, globals, args, config)
}

/* Report detailed information about a single repository link.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_INFO` | `lore_link_info_event_data_t` | Emitted once for the described link | */
func LinkInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLinkInfoFunc, globals, args, config)
}

/* List all repository links configured in the current repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_ENTRY` | `lore_link_entry_event_data_t` | Emitted for each linked repository | */
func LinkList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLinkListFunc, globals, args, config)
}

/* Update properties of an existing repository link.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Link Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LINK_CHANGE` | `lore_link_change_event_data_t` | Emitted when a link property is updated or finalized | */
func LinkUpdate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLinkUpdateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLinkUpdateFunc, globals, args, config)
}

/* Clone a remote repository to a local path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryCloneFunc, globals, args, config)
}

/* Retrieve metadata about the current repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_DATA` | `lore_repository_data_event_data_t` | Emitted with repository metadata (name, URL, branch info, etc.) | */
func RepositoryInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryInfoFunc, globals, args, config)
}

/* Dump the internal state of the repository for diagnostic purposes.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryDumpFunc, globals, args, config)
}

/* Create a new Lore repository on the remote server.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_CREATE` | `lore_repository_create_event_data_t` | Emitted when the repository has been successfully created | */
func RepositoryCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryCreateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryCreateFunc, globals, args, config)
}

/* Flush pending repository state to persistent storage.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryFlush(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryFlushArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryFlushFunc, globals, args, config)
}

/* Run garbage collection to reclaim unreferenced storage in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryGc(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryGcArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryGcFunc, globals, args, config)
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
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RepositoryRelease(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryReleaseArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryReleaseFunc, globals, args, config)
}

/* Add a new layer to the repository configuration.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Layer Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LAYER_ADD` | `lore_layer_add_event_data_t` | Emitted when a layer has been successfully added | */
func LayerAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerAddArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLayerAddFunc, globals, args, config)
}

/* Remove a layer from the repository configuration.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func LayerRemove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerRemoveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLayerRemoveFunc, globals, args, config)
}

/* List all layers configured in the repository.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Layer Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LAYER_ENTRY` | `lore_layer_entry_event_data_t` | Emitted for each layer configured in the repository | */
func LayerList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreLayerListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreLayerListFunc, globals, args, config)
}

/* List all repositories available on the remote server.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_LIST_ENTRY` | `lore_repository_list_entry_event_data_t` | Emitted for each repository found | */
func RepositoryList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryListFunc, globals, args, config)
}

/* Show the working directory status, including staged, dirty, and conflicted files.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryStatusFunc, globals, args, config)
}

/* Query the repository's immutable fragment store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Repository Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REPOSITORY_STORE_IMMUTABLE_QUERY` | `lore_repository_store_immutable_query_event_data_t` | Emitted for each fragment entry found in the immutable store | */
func RepositoryStoreImmutableQuery(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryStoreImmutableQueryArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryStoreImmutableQueryFunc, globals, args, config)
}

/* Verify the integrity of the repository's stored fragments.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryVerifyStateFunc, globals, args, config)
}

/* Commit staged files to create a new revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionCommitFunc, globals, args, config)
}

/* Amend the most recent revision with updated metadata.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_COMMIT_REVISION` | `lore_revision_commit_revision_event_data_t` | Emitted with the amended revision details |
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata entry of the amended revision | */
func RevisionAmend(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionAmendArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionAmendFunc, globals, args, config)
}

/* Retrieve metadata and delta information about a specific revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionInfoFunc, globals, args, config)
}

/* Show files that differ between two revisions.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_DIFF_FILE` | `lore_revision_diff_file_event_data_t` | Emitted for each file that differs between the two revisions |
| `LORE_EVENT_REVISION_RESOLVE` | `lore_revision_resolve_event_data_t` | Emitted when resolving a partial or numbered revision reference | */
func RevisionDiff(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionDiffArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionDiffFunc, globals, args, config)
}

/* Find a revision by metadata or revision number.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_FIND` | `lore_revision_find_event_data_t` | Emitted when a matching revision is found (exact or partial match) | */
func RevisionFind(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionFindArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionFindFunc, globals, args, config)
}

/* Retrieve the commit history of the current branch.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVISION_HISTORY` | `lore_revision_history_event_data_t` | Emitted once with summary info before entries stream |
| `LORE_EVENT_REVISION_HISTORY_ENTRY` | `lore_revision_history_entry_event_data_t` | Emitted for each revision in the history | */
func RevisionHistory(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionHistoryArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionHistoryFunc, globals, args, config)
}

/* Restore the working directory to a previously committed revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRestoreFunc, globals, args, config)
}

/* Clear all metadata from the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA_CLEAR_REVISION` | `lore_metadata_clear_revision_event_data_t` | Emitted when metadata has been cleared for the current revision | */
func RevisionMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionMetadataClearFunc, globals, args, config)
}

/* Get a specific metadata key/value pair from the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted with the requested key/value for the revision | */
func RevisionMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionMetadataGetFunc, globals, args, config)
}

/* List all metadata key/value pairs associated with the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revision Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_METADATA` | `lore_metadata_event_data_t` | Emitted for each metadata key/value associated with the revision | */
func RevisionMetadataList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionMetadataListFunc, globals, args, config)
}

/* Set a metadata key/value pair on the current revision.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func RevisionMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionMetadataSetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionMetadataSetFunc, globals, args, config)
}

/* Synchronize the working directory to a target revision, optionally merging divergent branches.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionSyncFunc, globals, args, config)
}

/* Revert a revision, applying its inverse changes to the working tree.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertFunc, globals, args, config)
}

/* Abort an in-progress revert operation and restore the previous state.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertAbortFunc, globals, args, config)
}

/* Mark conflicting files in a revert operation as unresolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_UNRESOLVE_FILE` | `lore_revert_unresolve_file_event_data_t` | Emitted for each file marked as unresolved |
| `LORE_EVENT_REVERT_UNRESOLVE_REVISION` | `lore_revert_unresolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertUnresolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertUnresolveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertUnresolveFunc, globals, args, config)
}

/* Restart a revert operation, re-materializing files with conflicts.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertRestartFunc, globals, args, config)
}

/* Resolve a revert conflict by marking conflicting files as resolved.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file marked as resolved |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolve(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertResolveFunc, globals, args, config)
}

/* Resolve a revert conflict by accepting the "mine" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file resolved by keeping "mine" |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolveMine(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveMineArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertResolveMineFunc, globals, args, config)
}

/* Resolve a revert conflict by accepting the "theirs" version of each conflicting file.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Revert Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_REVERT_RESOLVE_FILE` | `lore_revert_resolve_file_event_data_t` | Emitted for each file resolved by keeping "theirs" |
| `LORE_EVENT_REVERT_RESOLVE_REVISION` | `lore_revert_resolve_revision_event_data_t` | Emitted with the updated staged revision | */
func RevisionRevertResolveTheirs(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionRevertResolveTheirsArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionRevertResolveTheirsFunc, globals, args, config)
}

/* Create a new shared store at the specified path.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Shared Store Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_SHARED_STORE_CREATE` | `lore_shared_store_create_event_data_t` | Emitted on success after the shared store is created, carrying the path of the newly created store | */
func SharedStoreCreate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreCreateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreSharedStoreCreateFunc, globals, args, config)
}

/* Retrieve the path of the configured default shared store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Shared Store Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_SHARED_STORE_INFO` | `lore_shared_store_info_event_data_t` | Emitted on success carrying the path of the configured default shared store | */
func SharedStoreInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreSharedStoreInfoFunc, globals, args, config)
}

/* Set whether to automatically use the shared store.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func SharedStoreSetUseAutomatically(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreSharedStoreSetUseAutomaticallyArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreSharedStoreSetUseAutomaticallyFunc, globals, args, config)
}

/* Open a content-addressed storage handle.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_OPENED` | `lore_storage_opened_event_data_t` | Emitted on success carrying the opened handle id |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` on success or the error code on failure | */
func StorageOpen(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageOpenArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageOpenFunc, globals, args, config)
}

/* Store one or more content-addressed buffers.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_PUT_ITEM_COMPLETE` | `lore_storage_put_item_complete_event_data_t` | Emitted once per input item — success or failure; `stored_local`/`stored_remote` report where the content landed |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` iff every item succeeded, else the error code | */
func StoragePut(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStoragePutArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStoragePutFunc, globals, args, config)
}

/* Read one or more content-addressed buffers.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_GET_HEADER` | `lore_storage_get_header_event_data_t` | Size of the item's reassembled content, emitted before any DATA events |
| `LORE_EVENT_STORAGE_GET_DATA` | `lore_storage_get_data_event_data_t` | Payload bytes — valid only during the callback invocation |
| `LORE_EVENT_STORAGE_GET_ITEM_COMPLETE` | `lore_storage_get_item_complete_event_data_t` | Terminal per-item event |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` iff every item succeeded, else the error code | */
func StorageGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageGetFunc, globals, args, config)
}

/* Resolve one or more mutable keys and read the content they name, in one round trip.

`lore_storage_mutable_load` followed by `lore_storage_get`, performed by the server. The keys
are read under the `LORE_KEY_TYPE_RESOLVE` key type, which is what `lore_storage_put_resolved`
publishes; no other key type is resolvable this way.

The `address` in every event is the *resolved* address, so a caller can learn the key-to-hash
mapping from the event stream. A key with no mapping, or one naming absent content, reports
`error_code = ADDRESS_NOT_FOUND`, and the terminal event then carries a zero address.

Set `streaming` to receive one `LORE_EVENT_STORAGE_GET_DATA` per leaf fragment instead of a
single reassembled buffer, exactly as `lore_storage_get` does. Without it the whole content is
materialised in memory before the first byte reaches the callback, so a key naming something
large should set it.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_GET_HEADER` | `lore_storage_get_header_event_data_t` | Size of the item's reassembled content, emitted before any DATA events |
| `LORE_EVENT_STORAGE_GET_DATA` | `lore_storage_get_data_event_data_t` | Payload bytes — valid only during the callback invocation. One event per item, or one per leaf fragment when `streaming` is set |
| `LORE_EVENT_STORAGE_GET_ITEM_COMPLETE` | `lore_storage_get_item_complete_event_data_t` | Terminal per-item event |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` iff every item succeeded, else the error code | */
func StorageGetResolved(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetResolvedArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageGetResolvedFunc, globals, args, config)
}

/* Store one or more buffers and publish a mutable key naming each, in one round trip.

`lore_storage_put` followed by `lore_storage_mutable_store`, fused into one request when the
content fits a single fragment. The key is published under `LORE_KEY_TYPE_RESOLVE`, making it
readable by `lore_storage_get_resolved`, and the mapping is written only once the content is
stored — so a key published this way never resolves to content that is not there. Writing the
same key type directly with `lore_storage_mutable_store` carries no such guarantee.

The local store always receives both the content and the mapping. `remote_write = 1` also
publishes them remotely, matching `lore_storage_put`; there is no local-then-remote fallback.
A zero `key` or a zero `partition` rejects with `INVALID_ARGUMENTS`.

A zero-length `data` **removes** the key's mapping rather than publishing one: no content is
stored, the key is set to the zero hash, and `lore_storage_get_resolved` then reports
`ADDRESS_NOT_FOUND` for it. The terminal event carries the zero content hash and the caller's
context.

With `remote_write = 0` this evicts only the locally cached mapping. The local mutable store
is a cache, not an authority, so a key published remotely resolves again on the next call.
Deleting a published key requires `remote_write = 1`.

A remote content upload that fails still leaves a successful local write, so the key is not
published remotely and `stored_remote` is `0` while `error_code` stays `NONE`. Check
`stored_remote`, not `error_code`, to confirm the key is visible to other clients.

Publishing is last-writer-wins. Two callers publishing the same key concurrently both
succeed, and the key ends up naming whichever content was published second — the first
publisher is not told it was overwritten. Callers needing to detect a lost update should
store the content with `lore_storage_put` and publish with
`lore_storage_mutable_compare_and_swap`, which costs the second round trip this operation
exists to avoid.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_PUT_ITEM_COMPLETE` | `lore_storage_put_item_complete_event_data_t` | Emitted once per input item; `address` is the content the key now resolves to, and `stored_local`/`stored_remote` report where it landed |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` iff every item succeeded, else the error code | */
func StoragePutResolved(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStoragePutResolvedArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStoragePutResolvedFunc, globals, args, config)
}

/* Release a content-addressed storage handle.

Subsequent calls against the same handle return `InvalidArguments`.
Close does not block on the flush it spawns — `Complete` fires after
the in-flight counter drains, not after the flush finishes. */
func StorageClose(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageCloseArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageCloseFunc, globals, args, config)
}

/* Flush pending writes through the handle's stores.

On disk-backed stores this performs an fsync honoring `globals.sync_data`.
On in-memory stores the underlying flush is a no-op and the call still
completes with `status: 0`. */
func StorageFlush(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageFlushArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageFlushFunc, globals, args, config)
}

/* Fetch fragment metadata for one or more `(partition, address)` pairs without paying the
payload bytes. Each item's terminal event carries the resolved `Fragment` (`flags`,
`size_payload`, `size_content`); on miss `error_code == ADDRESS_NOT_FOUND`.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_GET_METADATA_ITEM_COMPLETE` | `lore_storage_get_metadata_item_complete_event_data_t` | Per-item terminal event |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status` is `0` iff every item succeeded, else the error code | */
func StorageGetMetadata(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetMetadataArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageGetMetadataFunc, globals, args, config)
}

/* Delete one or more `(partition, address)` entries from the store.

Idempotent on absent items; emits one `OBLITERATE_ITEM_COMPLETE` event
per item carrying `local_success` / `remote_success` / `error_code`. */
func StorageObliterate(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageObliterateArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageObliterateFunc, globals, args, config)
}

/* Read one or more mutable key values.

Each item acts on the local mutable store by default, or the remote mutable store when
`globals.remote` is set (or the handle was opened remote-bound), over the shared storage
session.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_MUTABLE_LOAD_ITEM_COMPLETE` | `lore_storage_mutable_load_item_complete_event_data_t` | Per-item terminal event carrying the value; `error_code == ADDRESS_NOT_FOUND` on a miss |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageMutableLoad(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageMutableLoadArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageMutableLoadFunc, globals, args, config)
}

/* Write one or more mutable key-value pairs. Storing the null value removes the key.

Each item acts on the local mutable store by default, or the remote mutable store when
`globals.remote` is set (or the handle was opened remote-bound), over the shared storage
session.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_MUTABLE_STORE_ITEM_COMPLETE` | `lore_storage_mutable_store_item_complete_event_data_t` | Per-item terminal event |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageMutableStore(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageMutableStoreArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageMutableStoreFunc, globals, args, config)
}

/* Conditionally swap one or more mutable key values. Each item updates the key to `value` when
its current value matches `expected`, and reports the value the key held before the swap.

Each item acts on the local mutable store by default, or the remote mutable store when
`globals.remote` is set (or the handle was opened remote-bound), over the shared storage
session.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE` | `lore_storage_mutable_compare_and_swap_item_complete_event_data_t` | Per-item terminal event carrying `previous`; the swap took effect when `previous == expected` |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageMutableCompareAndSwap(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageMutableCompareAndSwapArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageMutableCompareAndSwapFunc, globals, args, config)
}

/* List the mutable key-value pairs of a given type for one or more partitions.

Acts on the local mutable store only; a remote-targeted call (`globals.remote`, or a
remote-bound handle) is rejected with `INVALID_ARGUMENTS`. A zero/default partition lists
every partition the caller can access.

# Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_STORAGE_MUTABLE_LIST_ENTRY` | `lore_storage_mutable_list_entry_event_data_t` | One `(key, value)` pair, emitted before the item's terminal event |
| `LORE_EVENT_STORAGE_MUTABLE_LIST_ITEM_COMPLETE` | `lore_storage_mutable_list_item_complete_event_data_t` | Per-item terminal event |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | `status: 0` iff every item succeeded | */
func StorageMutableList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageMutableListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageMutableListFunc, globals, args, config)
}

/* Copy content from one partition to another within the same store.

Same-partition source/target rejects with `INVALID_ARGUMENTS`. The
item's content hash is preserved; only the source address is carried
in the per-item event. */
func StorageCopy(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageCopyArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageCopyFunc, globals, args, config)
}

/* Read one or more files into the content-addressed store.

Each item emits `LORE_EVENT_STORAGE_PUT_ITEM_COMPLETE` carrying the
computed address. Empty files short-circuit to the zero-hash address
without opening for read. */
func StoragePutFile(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStoragePutFileArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStoragePutFileFunc, globals, args, config)
}

/* Write content-addressed payloads to filesystem paths.

Each item emits `LORE_EVENT_STORAGE_GET_ITEM_COMPLETE`. No HEADER or
DATA events are produced — the payload is written straight to disk.
On partial-write failure the library leaves whatever state the
failure produced; cleanup is the caller's responsibility. */
func StorageGetFile(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreStorageGetFileArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageGetFileFunc, globals, args, config)
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreStorageUploadFunc, globals, args, config)
}

/* Start the Lore background service.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func ServiceStart(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreServiceStartArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreServiceStartFunc, globals, args, config)
}

/* Stop the Lore background service.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination | */
func ServiceStop(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreServiceStopArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreServiceStopFunc, globals, args, config)
}

/* Subscribe to repository notifications.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
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
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreNotificationSubscribeFunc, globals, args, config)
}

/* Unsubscribe from repository notifications.

# Events

Events are delivered via the callback as `lore_event_t`. Use the `tag` field to identify the event type.

## Standard Events

These events are emitted by all interface functions:

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_LOG` | `lore_log_event_data_t` | Diagnostic messages throughout execution |
| `LORE_EVENT_ERROR` | `lore_error_event_data_t` | Emitted for a non-fatal error during the operation |
| `LORE_EVENT_COMPLETE` | `lore_complete_event_data_t` | Always emitted at the end; `status` is `0` on success or the error code on failure |
| `LORE_EVENT_END` | `lore_end_event_data_t` | Always emitted after `COMPLETE` to signal callback termination |

## Notification Events

| Tag | Data Type | Description |
|-----|-----------|-------------|
| `LORE_EVENT_NOTIFICATION_UNSUBSCRIBED` | `lore_notification_unsubscribed_event_data_t` | Emitted when successfully unsubscribed from repository notifications | */
func NotificationUnsubscribe(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreNotificationUnsubscribeArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreNotificationUnsubscribeFunc, globals, args, config)
}

/* Retrieve repository metadata. Reads a single key, or all entries when no
key is given. */
func RepositoryMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryMetadataGetFunc, globals, args, config)
}

/* Set repository metadata key-value pairs. */
func RepositoryMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataSetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryMetadataSetFunc, globals, args, config)
}

/* Clear repository metadata keys. Clears all user-defined keys when none are
given. */
func RepositoryMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryMetadataClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryMetadataClearFunc, globals, args, config)
}

/* List the tracked instances of the repository. */
func RepositoryInstanceList(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInstanceListArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryInstanceListFunc, globals, args, config)
}

/* Remove stale instances of the repository that are no longer present. */
func RepositoryInstancePrune(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryInstancePruneArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryInstancePruneFunc, globals, args, config)
}

/* Update the recorded path of the current repository instance to its present
location. */
func RepositoryUpdatePath(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryUpdatePathArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryUpdatePathFunc, globals, args, config)
}

/* Read a configuration value of the current repository by key. */
func RepositoryConfigGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRepositoryConfigGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRepositoryConfigGetFunc, globals, args, config)
}

/* Open a memory-based revision tree handle on the given
`(store, repository, revision_hash)` tuple. `revision_hash == 0` opens an
empty tree suitable for committing an initial revision.

| Terminal event                       | Payload                                | Notes                                              |
|--------------------------------------|----------------------------------------|----------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_LOADED`    | `lore_revision_tree_loaded_event_data_t` | Emitted on success carrying the opened handle id | */
func RevisionTreeLoad(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeLoadArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeLoadFunc, globals, args, config)
}

/* Release a memory-based revision tree handle.

Subsequent calls against the same handle return `InvalidArguments`. The
call blocks until every in-flight op on the handle has paired its
decrement.

| Terminal event                              | Payload                                       | Notes                                              |
|---------------------------------------------|-----------------------------------------------|----------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_CLOSE_COMPLETE`   | `lore_revision_tree_close_complete_event_data_t` | Emitted on success carrying the caller id       | */
func RevisionTreeClose(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeCloseArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeCloseFunc, globals, args, config)
}

/* Resolve a UTF-8 path against a loaded revision tree to a node id. An empty
path resolves to the root node.

| Terminal event                                       | Payload                                             | Notes                                                       |
|------------------------------------------------------|-----------------------------------------------------|-------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_RESOLVE_PATH_COMPLETE`     | `lore_revision_tree_resolve_path_complete_event_data_t` | Carries the resolved node id and the per-call outcome   | */
func RevisionTreeResolvePath(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeResolvePathArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeResolvePathFunc, globals, args, config)
}

/* Stream the children of a directory node in a loaded revision tree.

| Terminal event                       | Payload                                | Notes                                                          |
|--------------------------------------|----------------------------------------|----------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_CHILD`     | `lore_revision_tree_child_event_data_t` | One per child; an empty directory emits none before `Complete` | */
func RevisionTreeListChildren(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeListChildrenArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeListChildrenFunc, globals, args, config)
}

/* Fetch the per-node record for a single node id in a loaded revision tree.

| Terminal event                          | Payload                                     | Notes                                                          |
|-----------------------------------------|---------------------------------------------|----------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_NODE_INFO`    | `lore_revision_tree_node_info_event_data_t` | Carries the node record, uniform across every node id (revision metadata: `lore_revision_tree_info`) | */
func RevisionTreeNodeInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeNodeInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeNodeInfoFunc, globals, args, config)
}

/* Fetch the loaded revision's record-level metadata (parents, creation
timestamp, author identity, metadata key count). Revision-scoped — no node id.

| Terminal event                     | Payload                                | Notes                                                   |
|------------------------------------|----------------------------------------|---------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_INFO`    | `lore_revision_tree_info_event_data_t` | Carries the revision record metadata for the handle     | */
func RevisionTreeInfo(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeInfoArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeInfoFunc, globals, args, config)
}

/* Reconstruct the full UTF-8 path for a node id by walking parent pointers,
relative to the handle's own tree root.

| Terminal event                       | Payload                                     | Notes                                                  |
|--------------------------------------|---------------------------------------------|--------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_NODE_PATH` | `lore_revision_tree_node_path_event_data_t` | Carries the path; the root resolves to the empty path  | */
func RevisionTreeNodePath(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeNodePathArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeNodePathFunc, globals, args, config)
}

/* Add a batch of nodes to a loaded revision tree. An entry parents onto an
existing node or onto an earlier entry, so one call builds a subtree. Every
entry is checked before any node is created, so one bad entry rejects the
call and creates nothing; the reason names the offending entry's batch index,
which a caller leaving `entry_id` at zero has no other way to identify. A failure
after those checks pass is internal and may leave part of the batch created.

A link entry's target revision is not resolved here, so a link naming a
revision that cannot be read is accepted and fails only when something later
reads through it. Entries under separate parents are created concurrently,
but allocating a node slot is serialized per loaded tree.

| Terminal event                            | Payload                                          | Notes                                                    |
|-------------------------------------------|--------------------------------------------------|----------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_ADD_COMPLETE`   | `lore_revision_tree_add_complete_event_data_t`   | One per entry, carrying its `entry_id`                    |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE` | `lore_revision_tree_batch_complete_event_data_t` | Exactly one, carrying the `batch_id` and the call's outcome | */
func RevisionTreeAdd(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeAddArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeAddFunc, globals, args, config)
}

/* Remove a batch of subtrees from a loaded revision tree. Each entry names a
subtree root and removes it whole, transitive children included. Every entry
is checked before any node is touched, so one bad entry rejects the call and
leaves every subtree in place; the reason names the offending entry's batch
index, which a caller leaving `entry_id` at zero has no other way to
identify. A failure after those checks pass is internal and may leave part of
the batch removed.

A node the loaded revision holds is **staged** for deletion and stays in the
tree: it keeps its name and its place among its siblings, still lists through
`lore_revision_tree_list_children`, and reports
`LORE_NODE_STAGED_ACTION_DELETE` in its `staged_action`. The commit that
freezes the tree is what drops it. A node added through this handle is in no
revision yet, so it is discarded outright instead, freeing its name and its
node id. A link is removed as one node — its subtree belongs to the linked
repository's tree.

A staged deletion is reversible: `lore_revision_tree_add` of the same name
under the same parent with the same kind restores the node. A zero
`address.context` on that add preserves the node's `file_id`; supplying one
replaces it, as on `lore_revision_tree_modify`. Only the named node comes back
— restoring a directory leaves its children staged for deletion, so each has to
be added back in turn. A discarded node is not restorable, since its id is
gone.

Staging fans out one depth level at a time; discarding an added node rewrites
sibling pointers and so runs serially, deepest first. Memory while the call
runs is proportional to the widest level of the subtrees being removed rather
than to the entry count, since a level is collected before it is staged. The
root cannot be deleted, and an entry whose ancestor another entry deletes is
rejected rather than removed twice.

| Terminal event                              | Payload                                           | Notes                                                    |
|---------------------------------------------|---------------------------------------------------|----------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_DELETE_COMPLETE`  | `lore_revision_tree_delete_complete_event_data_t` | One per entry, carrying its `entry_id` and `node_count`   |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE`   | `lore_revision_tree_batch_complete_event_data_t`  | Exactly one, carrying the `batch_id` and the call's outcome | */
func RevisionTreeDelete(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeDeleteArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeDeleteFunc, globals, args, config)
}

/* Rewrite a batch of file nodes' `mode`, `size` and `address` in a loaded
revision tree. Every entry is checked before any node is rewritten, so one bad
entry rejects the call and leaves every target untouched; the reason names the
offending entry's batch index, which a caller leaving `entry_id` at zero has no
other way to identify. A failure after those checks pass is internal and may
leave part of the batch rewritten.

Only a file is modifiable: a directory's size and address are derived at
commit and a link's address is its target. A zero `address.context` preserves
the node's existing file id rather than generating one, which is the opposite
of `lore_revision_tree_add` — the node already has an identity, and replacing
it would record the edit as a move. Entries touch no parent or sibling chain,
so the whole batch applies concurrently.

| Terminal event                              | Payload                                           | Notes                                                    |
|---------------------------------------------|---------------------------------------------------|----------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_MODIFY_COMPLETE`  | `lore_revision_tree_modify_complete_event_data_t` | One per entry, carrying its `entry_id`                    |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE`   | `lore_revision_tree_batch_complete_event_data_t`  | Exactly one, carrying the `batch_id` and the call's outcome | */
func RevisionTreeModify(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeModifyArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeModifyFunc, globals, args, config)
}

/* Move a batch of nodes to new parents and/or new names in a loaded revision tree. An
entry naming the node's current parent renames it where it is. Every entry is checked
before any node is moved, so one bad entry rejects the call and leaves every node
where it was; the reason names the offending entry's batch index, which a caller
leaving `entry_id` at zero has no other way to identify. A failure after those checks
pass is internal and may leave earlier entries applied.

A move keeps the node: its node id, its `file_id` and its children come along, and
the change is recorded as a move rather than as a deletion and an addition, so the
revision graph carries the node's history across it. The node reports
`LORE_NODE_STAGED_ACTION_MOVE` until the commit that freezes the tree, and so does
every node under a moved directory — their records do not change, but their paths do.
Two exceptions: a node added through this handle stays staged as an addition wherever
it lands, since it is in no revision a move could be recorded against; and a node
under the moved directory that is staged for deletion keeps its deletion, since it is
leaving the revision at the commit either way.

Both batch-level rules read the tree the whole batch produces rather than the one in
front of them. A destination inside the moved node's own subtree is rejected, and so
is one that lands there once the batch is applied — moving A under B and B under A is
a loop neither entry shows on its own. A name a live child of the destination already
holds is rejected, but a name the batch itself vacates is not: moving `x` out of a
directory while moving another node to `x` in it succeeds, and two entries taking one
name under one destination reject even though neither collides with the tree.

Entries apply one at a time, in batch order, because a move rewrites the parent and
sibling pointers of two child chains where `lore_revision_tree_add` only prepends to
one. For the same reason concurrent calls have more to lose here than on `add` or
`modify`: two calls moving nodes that share a parent chain can interleave their
unlinks, which the pre-commit validator then refuses. Moves that may touch one parent
chain belong in one call.

| Terminal event                            | Payload                                          | Notes                                                       |
|-------------------------------------------|--------------------------------------------------|-------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_MOVE_COMPLETE`  | `lore_revision_tree_move_complete_event_data_t`  | One per entry, carrying its `entry_id` and the moved node    |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE` | `lore_revision_tree_batch_complete_event_data_t` | Exactly one, carrying the `batch_id` and the call's outcome  | */
func RevisionTreeMove(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeMoveArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeMoveFunc, globals, args, config)
}

/* Record a batch of `(key, value)` pairs on a loaded revision tree's in-progress
metadata.

`value` is a `lore_metadata_t`, which carries its own kind — set the union's
`tag` and the matching member. There is no separate format field and nothing
is parsed, so a value cannot be stored under a kind it is not, and a binary
value can hold any bytes rather than only text. This differs from
`lore_revision_metadata_set`, which takes text plus a parallel format array.
`lore_revision_tree_metadata_get` returns the same union.

Every entry is checked before any pair is recorded, so one bad entry rejects
the call and records nothing; the reason names the offending entry's batch
index, which a caller leaving `entry_id` at zero has no other way to
identify.

A repeated key is **not** rejected, unlike the duplicate-target rules on the
node verbs: entries apply in index order, so the last entry naming a key wins
— the same result as sending those pairs as separate calls.

Nothing reaches storage here. The pairs live on the handle until
`lore_revision_tree_commit` serializes them, so only this handle's
`lore_revision_tree_metadata_get` sees them. The whole batch applies under one
write lock, which is what makes it atomic and why there is no concurrency to
gain: the work is buffer writes, not I/O.

A revision's whole metadata is capped at 1 MiB. The cap counts the metadata
itself — keys, values and per-entry overhead — and not what a value refers
to: a value holding a content address costs the address, not the content
behind it.

A single entry larger than the whole cap is rejected here, since no amount of
removing other keys could make it fit. The running total is not checked here,
because what a revision ends up carrying is only known once every set has
run: a batch of individually legal entries that together push past the limit
is reported as recorded and fails later, at `lore_revision_tree_commit`.

| Terminal event                                    | Payload                                                 | Notes                                                       |
|---------------------------------------------------|---------------------------------------------------------|-------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_METADATA_SET_COMPLETE`  | `lore_revision_tree_metadata_set_complete_event_data_t` | One per entry, carrying its `entry_id`                      |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE`         | `lore_revision_tree_batch_complete_event_data_t`        | Exactly one, carrying the `batch_id` and the call's outcome  | */
func RevisionTreeMetadataSet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeMetadataSetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeMetadataSetFunc, globals, args, config)
}

/* Read a batch of metadata values from a loaded revision tree.

By default this reads only the **revision being built** — what
`lore_revision_tree_metadata_set` recorded on this handle, which is exactly
what `lore_revision_tree_commit` will write. Nothing is inherited from the
revision the handle was loaded on. Set `include_revision` to `1` to also fall
back to that revision for a key the handle has no entry for; a pending entry
still wins, so the flag only adds answers.

A key present in neither emits **no event at all** and does not fail the call,
matching `lore_revision_metadata_get`: detect an absent key by tracking
whether a value event arrived for its `entry_id`. This verb is therefore **not
all-or-nothing**, unlike the other batch verbs — it mutates nothing, so one
unanswerable key costs the others nothing. Bad arguments still reject the
whole call.

A value comes back as the same `lore_metadata_t` that
`lore_revision_tree_metadata_set` takes, so it round-trips without either
side encoding it as text. Every kind is returned, raw binary included.

A value whose stored bytes do not match the kind they are tagged with, or
whose tag this build does not recognize, reports `LORE_ERROR_CODE_INTERNAL`
on that entry rather than staying silent, so it is never mistaken for an
absent key.

The revision's metadata is read once for the whole batch, which is what
batching buys here.

| Terminal event                                    | Payload                                                 | Notes                                                       |
|---------------------------------------------------|---------------------------------------------------------|-------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_METADATA_GET_COMPLETE`  | `lore_revision_tree_metadata_get_complete_event_data_t` | One per key that resolved, carrying its `entry_id`; none for an absent key |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE`         | `lore_revision_tree_batch_complete_event_data_t`        | Exactly one, carrying the `batch_id` and the call's outcome  | */
func RevisionTreeMetadataGet(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeMetadataGetArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeMetadataGetFunc, globals, args, config)
}

/* Remove a batch of keys from a loaded revision tree's in-progress metadata.
Every entry is checked before any key is removed, so one bad entry rejects the
call and removes nothing; the reason names the offending entry's batch index,
which a caller leaving `entry_id` at zero has no other way to identify.

**Clearing a key that is not set is a no-op success**, not a failure. The
terminal's `removed` field says which happened: `1` when the key was there and
is now gone, `0` when there was nothing to remove. A repeated key is likewise
not rejected — the second entry naming it reports `removed = 0`.

This clears the **pending** metadata that `lore_revision_tree_metadata_set`
records and `lore_revision_tree_metadata_get` reads first; a key frozen in the
loaded revision is not reachable from here, since clearing edits the revision
being built and the one it was loaded from is immutable. The whole batch
applies under one write lock, which is what makes it atomic.

| Terminal event                                      | Payload                                                   | Notes                                                       |
|-----------------------------------------------------|-----------------------------------------------------------|-------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_METADATA_CLEAR_COMPLETE`  | `lore_revision_tree_metadata_clear_complete_event_data_t` | One per entry, carrying its `entry_id` and `removed`        |
| `LORE_EVENT_REVISION_TREE_BATCH_COMPLETE`           | `lore_revision_tree_batch_complete_event_data_t`          | Exactly one, carrying the `batch_id` and the call's outcome  | */
func RevisionTreeMetadataClear(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeMetadataClearArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeMetadataClearFunc, globals, args, config)
}

/* Freeze a loaded revision tree into a new revision and advance its branch tip.

**The branch is not an argument.** It is the revision's own, read from the
`branch` metadata key: set it with `lore_revision_tree_metadata_set` to start a
branch's history, and leave it unset to continue the loaded revision's branch. A
key that is set must name either the loaded revision's branch or a branch whose
branch point is exactly the loaded revision. A handle loaded from the zero
revision has no parent to read a branch from and must set the key.

The revision records exactly the metadata set on the handle — nothing is
inherited from the revision it was loaded on — plus the three facts about the
commit the caller did not supply: the branch, the timestamp if unset, and
`created-by` / `committed-by` if unset. The commit message is caller metadata like
any other: set `"message"` before committing.

On success the handle stays usable and now *is* the new revision: node ids
captured before the commit still resolve, and further edits commit on top. The
pending metadata is emptied, so the next revision starts fresh.

**A commit is all-or-nothing against the handle.** Either it succeeds and the
handle is consistent on the new revision, or it fails and the handle is
consistent on the state it had before the call. A call rejected before any write
— nothing staged, an unusable branch, a tree the validator refuses, or a branch
tip that has already moved — writes nothing at all. A failure once the freeze has
begun leaves a part-frozen tree, which is discarded and rebuilt from a snapshot
taken before the freeze started, so the handle comes back on the revision it was
on with the edits still staged. Either way recovery is to fix what the terminal
reported and retry **on the same handle**: no close, no reload, no re-applying
edits. The one failure that still poisons the handle is a restore that itself
fails, which reports `INTERNAL` saying so. When the branch had advanced,
`new_tip_hash` on the terminal carries that tip, which is also how a caller tells
that failure apart: neither a tip collision nor an empty commit has a
`lore_error_code_t` of its own, so both report `INTERNAL` with the reason in the
completion detail — the same codes the file-system commit returns.

`options.remote_write = 1` uploads within the call. It is a request, not a
guarantee: a store bound offline or local-only, or a call passing
`globals.local`, silently commits local-only. So does a store opened without a
remote configuration — there is nothing to upload to, and the commit still
reports success. Per-call flags contradicting the store's bound flags reject the
call.

**The commit holds the handle for the length of the call.** No other call on the
same handle runs while the tree is frozen, so an edit issued concurrently lands
wholly before the commit reads the tree or wholly after it finishes, two commits
on one handle serialize, and `metadata_set` can no longer lose an edit to the
commit. A commit on a large tree therefore blocks reads on that handle for its
duration, and a callback that re-enters the API on the same handle deadlocks —
which the callback contract already forbids. Commits from *different* handles or
processes still race, and the branch tip compare-and-swap decides them.

`remote_write` is resolved onto the handle's shared repository context, so the
value outlives the call: the handle carries whatever the last commit resolved.

| Terminal event                                | Payload                                             | Notes                                                             |
|-----------------------------------------------|-----------------------------------------------------|-------------------------------------------------------------------|
| `LORE_EVENT_REVISION_TREE_COMMIT_COMPLETE`    | `lore_revision_tree_commit_complete_event_data_t`   | Exactly one; carries the new revision, or the new tip on collision |
| `LORE_EVENT_REVISION_COMMIT_REVISION`         | `lore_revision_commit_revision_event_data_t`        | On success, for continuity with file-system commit consumers       | */
func RevisionTreeCommit(
	globals *types.LoreGlobalArgsFFI,
	args *types.LoreRevisionTreeCommitArgsFFI,
	config *types.LoreEventCallbackConfig,
) (int32, error) {
	return callLoreFunction(&loreRevisionTreeCommitFunc, globals, args, config)
}

func LogConfigure(logConfig *types.LoreLogConfigFFI) (int32, error) {
	if err := ensureLibrary(); err != nil {
		return -1, err
	}
	var pinner runtime.Pinner
	defer pinner.Unpin()
	pinner.Pin(logConfig)
	result := loreLogConfigureFunc(uintptr(unsafe.Pointer(logConfig)))
	return result, nil
}

func Shutdown() (int32, error) {
	if err := ensureLibrary(); err != nil {
		return -1, err
	}
	result := loreShutdownFunc()
	return result, nil
}

func SetThreadLimit(count uintptr) (int32, error) {
	if err := ensureLibrary(); err != nil {
		return -1, err
	}
	result := loreSetThreadLimitFunc(count)
	return result, nil
}

func Version() (string, error) {
	if err := ensureLibrary(); err != nil {
		return "", err
	}

	cString := loreVersionFunc()

	// find the char* string null terminator
	p := (*byte)(unsafe.Pointer(cString))
	length := 0
	for {
		if *(*byte)(unsafe.Add(unsafe.Pointer(p), length)) == 0 {
			break
		}
		length++
	}
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(cString)), length)

	return string(bytes), nil
}
