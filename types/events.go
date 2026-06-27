// Copyright Epic Games, Inc. All Rights Reserved.

package types

import (
	"fmt"
	"unsafe"
)

// LoreEventFFI is a C-compatible representation of lore_event_t
// Only valid during the callback
type LoreEventFFI struct {
	Tag     LoreEventTag
	padding [4]byte // Ensure union starts at 8-byte boundary
	// Union data follows (we access it via unsafe pointer arithmetic)
}

// Precomputed offset to the union data in LoreEventFFI
const loreEventUnionOffset = unsafe.Sizeof(LoreEventTag(0)) + unsafe.Sizeof([4]byte{})

// LoreEvent is a Go-native event (safe to keep after callback)
type LoreEvent struct {
	Tag  LoreEventTag
	Data any
}

type LoreRepositoryVerifyFragmentMatchEventDataArrayFFI struct {
	Ptr   uintptr
	Count uint64
}

type LoreRepositoryVerifyFragmentMatchEventDataArray = []LoreRepositoryVerifyFragmentMatchEventData

func (arr LoreRepositoryVerifyFragmentMatchEventDataArrayFFI) Len() int {
	return int(arr.Count)
}

func (arr LoreRepositoryVerifyFragmentMatchEventDataArrayFFI) Get(index int) LoreRepositoryVerifyFragmentMatchEventData {
	if index < 0 || index >= int(arr.Count) {
		panic(fmt.Sprintf("index out of bounds: %d (len=%d)", index, arr.Count))
	}
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	slice := unsafe.Slice((*LoreRepositoryVerifyFragmentMatchEventData)(unsafe.Pointer(arr.Ptr)), arr.Count)
	return slice[index]
}

func (arr LoreRepositoryVerifyFragmentMatchEventDataArrayFFI) Clone() []LoreRepositoryVerifyFragmentMatchEventData {
	if arr.Ptr == 0 {
		panic("cannot access FFI data outside the callback function")
	}
	if arr.Count == 0 {
		return nil
	}
	cDataSlice := unsafe.Slice((*LoreRepositoryVerifyFragmentMatchEventData)(unsafe.Pointer(arr.Ptr)), arr.Count)
	result := make([]LoreRepositoryVerifyFragmentMatchEventData, arr.Count)
	copy(result, cDataSlice)
	return result
}

type LoreProgressEventDataFFI struct {
	/* Placeholder field; carries no meaningful value. */
	Unused uint32
}

type LoreProgressEventData struct {
	/* Placeholder field; carries no meaningful value. */
	Unused uint32
}
type LoreErrorEventDataFFI struct {
	/* The error code, matching one of the error codes. */
	ErrorType uint32
	/* The underlying error message. */
	ErrorInner LoreString
}

type LoreErrorEventData struct {
	/* The error code, matching one of the error codes. */
	ErrorType uint32
	/* The underlying error message. */
	ErrorInner string
}
type LoreCompleteEventDataFFI struct {
	/* The completion status code of the operation. */
	Status int32
	/* The error detail for the operation. The empty default detail on
	success; the populated detail on failure. `#[serde(default)]` lets an
	older payload that lacks this field deserialize: the detail then reads
	back as the empty default with an empty trace list. */
	Error LoreErrorDetailFFI
}

type LoreCompleteEventData struct {
	/* The completion status code of the operation. */
	Status int32
	/* The error detail for the operation. The empty default detail on
	success; the populated detail on failure. `#[serde(default)]` lets an
	older payload that lacks this field deserialize: the detail then reads
	back as the empty default with an empty trace list. */
	Error LoreErrorDetail
}
type LoreMetadataEventDataFFI struct {
	/* The metadata key. */
	Key LoreString
	/* The metadata value. */
	Value LoreMetadataFFI
}

type LoreMetadataEventData struct {
	/* The metadata key. */
	Key string
	/* The metadata value. */
	Value LoreMetadata
}
type LoreLogEventDataFFI struct {
	/* The severity level of the log message. */
	Level LoreLogLevel
	/* The category of the log message. */
	Category uint32
	/* The time the message was produced. */
	Timestamp uint64
	/* The source location that produced the message. */
	Location LoreString
	/* The log message text. */
	Message LoreString
}

type LoreLogEventData struct {
	/* The severity level of the log message. */
	Level LoreLogLevel
	/* The category of the log message. */
	Category uint32
	/* The time the message was produced. */
	Timestamp uint64
	/* The source location that produced the message. */
	Location string
	/* The log message text. */
	Message string
}
type LoreEndEventDataFFI struct {
	/* Placeholder field; carries no meaningful value. */
	Unused uint32
}

type LoreEndEventData struct {
	/* Placeholder field; carries no meaningful value. */
	Unused uint32
}
type LoreMaintenanceEventDataFFI struct {
	/* The maintenance message text. */
	Message LoreString
}

type LoreMaintenanceEventData struct {
	/* The maintenance message text. */
	Message string
}
type LoreAuthUrlEventDataFFI struct {
	/* Authentication URL */
	Url LoreString
}

type LoreAuthUrlEventData struct {
	/* Authentication URL */
	Url string
}
type LoreAuthUserInfoEventDataFFI struct {
	/* User identity */
	Id LoreString
	/* Display name for the user */
	Name LoreString
}

type LoreAuthUserInfoEventData struct {
	/* User identity */
	Id string
	/* Display name for the user */
	Name string
}
type LoreAuthUserTokenEventDataFFI struct {
	/* User identity */
	Id LoreString
	/* Display name for the user */
	Name LoreString
	/* The token string */
	Token LoreString
	/* Preferred username from the token */
	PreferredUsername LoreString
	/* Non-zero if the identity is a service account */
	FlagServiceAccount uint8
	/* Expiry time in milliseconds since UNIX epoch, or 0 if unavailable */
	Expires uint64
}

type LoreAuthUserTokenEventData struct {
	/* User identity */
	Id string
	/* Display name for the user */
	Name string
	/* The token string */
	Token string
	/* Preferred username from the token */
	PreferredUsername string
	/* Non-zero if the identity is a service account */
	FlagServiceAccount bool
	/* Expiry time in milliseconds since UNIX epoch, or 0 if unavailable */
	Expires uint64
}
type LoreAuthIdentityEventDataFFI struct {
	/* Auth service URL */
	AuthUrl LoreString
	/* Resource ID (empty for authentication tokens) */
	Resource LoreString
	/* User identity */
	UserId LoreString
	/* Comma-separated list of authorized root domains */
	AuthorizedDomains LoreString
	/* Expiry time in milliseconds since UNIX epoch, or 0 if unavailable */
	Expires uint64
	/* Cached token (only populated when requested) */
	Token LoreString
}

type LoreAuthIdentityEventData struct {
	/* Auth service URL */
	AuthUrl string
	/* Resource ID (empty for authentication tokens) */
	Resource string
	/* User identity */
	UserId string
	/* Comma-separated list of authorized root domains */
	AuthorizedDomains string
	/* Expiry time in milliseconds since UNIX epoch, or 0 if unavailable */
	Expires uint64
	/* Cached token (only populated when requested) */
	Token string
}
type LoreBranchCreateEventDataFFI struct {
	/* Name of the created branch. */
	Name LoreString
	/* Latest revision the new branch points at. */
	Latest LoreHash
	/* Set when creating the branch also produced a new commit. */
	IsCommit uint8
}

type LoreBranchCreateEventData struct {
	/* Name of the created branch. */
	Name string
	/* Latest revision the new branch points at. */
	Latest LoreHash
	/* Set when creating the branch also produced a new commit. */
	IsCommit bool
}
type LoreBranchMultipleInstanceEventDataFFI struct {
	/* The branch checked out by more than one instance */
	Branch LoreBranchId
	/* Identifiers of the other instances on the branch */
	InstanceIds LoreInstanceIdArrayFFI
	/* Filesystem paths of the other instances on the branch */
	InstancePaths LoreStringArrayFFI
}

type LoreBranchMultipleInstanceEventData struct {
	/* The branch checked out by more than one instance */
	Branch LoreBranchId
	/* Identifiers of the other instances on the branch */
	InstanceIds LoreInstanceIdArray
	/* Filesystem paths of the other instances on the branch */
	InstancePaths []string
}
type LoreBranchArchiveEventDataFFI struct {
	/* Name of the archived branch. */
	Name LoreString
}

type LoreBranchArchiveEventData struct {
	/* Name of the archived branch. */
	Name string
}
type LoreBranchListBeginEventDataFFI struct {
	/* Location the listed branches come from. */
	Location LoreBranchLocation
}

type LoreBranchListBeginEventData struct {
	/* Location the listed branches come from. */
	Location LoreBranchLocation
}
type LoreBranchListEntryEventDataFFI struct {
	/* Location this branch comes from. */
	Location LoreBranchLocation
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name LoreString
	/* Branch category. */
	Category LoreString
	/* Latest revision the branch points at. */
	Latest LoreHash
	/* Stack of branch points this branch was created from. */
	Stack LoreBranchPointArrayFFI
	/* Identifier of the user who created the branch. */
	Creator LoreString
	/* Creation time of the branch as a timestamp. */
	Created uint64
	/* Set when this branch is the current branch. */
	IsCurrent uint8
	/* Set when this branch has been archived. */
	Archived uint8
}

type LoreBranchListEntryEventData struct {
	/* Location this branch comes from. */
	Location LoreBranchLocation
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name string
	/* Branch category. */
	Category string
	/* Latest revision the branch points at. */
	Latest LoreHash
	/* Stack of branch points this branch was created from. */
	Stack LoreBranchPointArray
	/* Identifier of the user who created the branch. */
	Creator string
	/* Creation time of the branch as a timestamp. */
	Created uint64
	/* Set when this branch is the current branch. */
	IsCurrent bool
	/* Set when this branch has been archived. */
	Archived bool
}
type LoreBranchListEndEventDataFFI struct {
	/* Location the listed branches came from. */
	Location LoreBranchLocation
	/* Number of branches that were listed. */
	Count uint64
}

type LoreBranchListEndEventData struct {
	/* Location the listed branches came from. */
	Location LoreBranchLocation
	/* Number of branches that were listed. */
	Count uint64
}
type LoreBranchMergeAbortBeginEventDataFFI struct {
	/* The staged revision being discarded. */
	StateStagedRevision LoreHash
	/* The current revision the working state returns to. */
	StateCurrentRevision LoreHash
}

type LoreBranchMergeAbortBeginEventData struct {
	/* The staged revision being discarded. */
	StateStagedRevision LoreHash
	/* The current revision the working state returns to. */
	StateCurrentRevision LoreHash
}
type LoreBranchMergeAbortEndEventDataFFI struct {
	/* Placeholder field. The event carries no payload. */
	Unused uint32
}

type LoreBranchMergeAbortEndEventData struct {
	/* Placeholder field. The event carries no payload. */
	Unused uint32
}
type LoreBranchInfoEventDataFFI struct {
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name LoreString
	/* Branch category. */
	Category LoreString
	/* Latest revision known locally for the branch. */
	Latest LoreHash
	/* Latest revision known on the remote for the branch. */
	LatestRemote LoreHash
	/* Identifier of the parent branch. */
	Parent LoreBranchId
	/* Revision on the parent branch where this branch was created. */
	BranchPoint LoreHash
	/* Identifier of the user who created the branch. */
	Creator LoreString
	/* Creation time of the branch as a timestamp. */
	Created uint64
	/* Stack of branch points this branch was created from. */
	Stack LoreBranchPointArrayFFI
	/* Set when the branch has been archived. */
	Archived uint8
}

type LoreBranchInfoEventData struct {
	/* Branch identifier. */
	Id LoreBranchId
	/* Branch name. */
	Name string
	/* Branch category. */
	Category string
	/* Latest revision known locally for the branch. */
	Latest LoreHash
	/* Latest revision known on the remote for the branch. */
	LatestRemote LoreHash
	/* Identifier of the parent branch. */
	Parent LoreBranchId
	/* Revision on the parent branch where this branch was created. */
	BranchPoint LoreHash
	/* Identifier of the user who created the branch. */
	Creator string
	/* Creation time of the branch as a timestamp. */
	Created uint64
	/* Stack of branch points this branch was created from. */
	Stack LoreBranchPointArray
	/* Set when the branch has been archived. */
	Archived bool
}
type LoreBranchDiffBeginEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreBranchDiffBeginEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreBranchDiffChangeBeginEventDataFFI struct {
	/* Number of changes that follow. */
	ChangesCount uintptr
}

type LoreBranchDiffChangeBeginEventData struct {
	/* Number of changes that follow. */
	ChangesCount uintptr
}
type LoreBranchDiffChangeEventDataFFI struct {
	/* The changed node. */
	Change LoreBranchDiffNodeDataFFI
}

type LoreBranchDiffChangeEventData struct {
	/* The changed node. */
	Change LoreBranchDiffNodeData
}
type LoreBranchDiffChangeEndEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreBranchDiffChangeEndEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreBranchDiffConflictBeginEventDataFFI struct {
	/* Number of conflicts that follow. */
	ConflictsCount uintptr
}

type LoreBranchDiffConflictBeginEventData struct {
	/* Number of conflicts that follow. */
	ConflictsCount uintptr
}
type LoreBranchDiffConflictEventDataFFI struct {
	/* The change on the source side of the conflict. */
	SourceChange LoreBranchDiffNodeDataFFI
	/* The change on the target side of the conflict. */
	TargetChange LoreBranchDiffNodeDataFFI
}

type LoreBranchDiffConflictEventData struct {
	/* The change on the source side of the conflict. */
	SourceChange LoreBranchDiffNodeData
	/* The change on the target side of the conflict. */
	TargetChange LoreBranchDiffNodeData
}
type LoreBranchDiffConflictEndEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreBranchDiffConflictEndEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreBranchDiffEndEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreBranchDiffEndEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreBranchLatestListEntryEventDataFFI struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Revision recorded in the history entry. */
	Revision LoreHash
}

type LoreBranchLatestListEntryEventData struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Revision recorded in the history entry. */
	Revision LoreHash
}
type LoreBranchMergeConflictFileEventDataFFI struct {
	/* The path of the conflicted file. */
	Path LoreString
}

type LoreBranchMergeConflictFileEventData struct {
	/* The path of the conflicted file. */
	Path string
}
type LoreBranchMergeLinkSkippedEventDataFFI struct {
	/* The mount path of the skipped link. */
	LinkPath LoreString
	/* The repository of the skipped link. */
	Repository LoreRepositoryId
	/* The reason the link was skipped. */
	Reason uint8
}

type LoreBranchMergeLinkSkippedEventData struct {
	/* The mount path of the skipped link. */
	LinkPath string
	/* The repository of the skipped link. */
	Repository LoreRepositoryId
	/* The reason the link was skipped. */
	Reason bool
}
type LoreBranchMergeUnresolveFileEventDataFFI struct {
	/* The path of the file marked unresolved. */
	Path LoreString
}

type LoreBranchMergeUnresolveFileEventData struct {
	/* The path of the file marked unresolved. */
	Path string
}
type LoreBranchMergeUnresolveRevisionEventDataFFI struct {
	/* The repository of the revision marked unresolved. */
	Repository LoreRepositoryId
	/* The revision marked unresolved. */
	Revision LoreHash
}

type LoreBranchMergeUnresolveRevisionEventData struct {
	/* The repository of the revision marked unresolved. */
	Repository LoreRepositoryId
	/* The revision marked unresolved. */
	Revision LoreHash
}
type LoreBranchMergeIntoFileBeginEventDataFFI struct {
	/* The number of files to merge. */
	Count uintptr
}

type LoreBranchMergeIntoFileBeginEventData struct {
	/* The number of files to merge. */
	Count uintptr
}
type LoreBranchMergeIntoFileEventDataFFI struct {
	/* The path of the file. */
	Path LoreString
	/* The action applied to the file. */
	Action LoreFileAction
	/* The size of the file in bytes. */
	Size uint64
	/* Set when the entry is a regular file. */
	IsFile uint8
	/* Set when the entry is a directory. */
	IsDirectory uint8
	/* Set when the entry is a link. */
	IsLink uint8
}

type LoreBranchMergeIntoFileEventData struct {
	/* The path of the file. */
	Path string
	/* The action applied to the file. */
	Action LoreFileAction
	/* The size of the file in bytes. */
	Size uint64
	/* Set when the entry is a regular file. */
	IsFile bool
	/* Set when the entry is a directory. */
	IsDirectory bool
	/* Set when the entry is a link. */
	IsLink bool
}
type LoreBranchMergeIntoFileEndEventDataFFI struct {
	/* The number of files merged. */
	Count uintptr
}

type LoreBranchMergeIntoFileEndEventData struct {
	/* The number of files merged. */
	Count uintptr
}
type LoreBranchMergeIntoFragmentBeginEventDataFFI struct {
	/* The number of fragments to transfer. */
	Fragments uint64
}

type LoreBranchMergeIntoFragmentBeginEventData struct {
	/* The number of fragments to transfer. */
	Fragments uint64
}
type LoreBranchMergeIntoFragmentProgressEventDataFFI struct {
	/* The number of fragments transferred so far. */
	Complete uint64
	/* The total number of fragments to transfer. */
	Count uint64
}

type LoreBranchMergeIntoFragmentProgressEventData struct {
	/* The number of fragments transferred so far. */
	Complete uint64
	/* The total number of fragments to transfer. */
	Count uint64
}
type LoreBranchMergeIntoFragmentEndEventDataFFI struct {
	/* The number of fragments transferred. */
	Fragments uint64
}

type LoreBranchMergeIntoFragmentEndEventData struct {
	/* The number of fragments transferred. */
	Fragments uint64
}
type LoreBranchMergeIntoRevisionEventDataFFI struct {
	/* The revision merged. */
	Revision LoreHash
	/* The sequential number of the revision. */
	RevisionNumber uint64
}

type LoreBranchMergeIntoRevisionEventData struct {
	/* The revision merged. */
	Revision LoreHash
	/* The sequential number of the revision. */
	RevisionNumber uint64
}
type LoreBranchMergeIntoSyncBeginEventDataFFI struct {
	/* The number of revisions to synchronize. */
	Count uintptr
}

type LoreBranchMergeIntoSyncBeginEventData struct {
	/* The number of revisions to synchronize. */
	Count uintptr
}
type LoreBranchMergeIntoSyncEndEventDataFFI struct {
	/* The number of revisions synchronized. */
	Count uintptr
}

type LoreBranchMergeIntoSyncEndEventData struct {
	/* The number of revisions synchronized. */
	Count uintptr
}
type LoreBranchMergeResolveFileEventDataFFI struct {
	/* The path of the file marked resolved. */
	Path LoreString
}

type LoreBranchMergeResolveFileEventData struct {
	/* The path of the file marked resolved. */
	Path string
}
type LoreBranchMergeResolveRevisionEventDataFFI struct {
	/* The repository of the revision marked resolved. */
	Repository LoreRepositoryId
	/* The revision marked resolved. */
	Revision LoreHash
}

type LoreBranchMergeResolveRevisionEventData struct {
	/* The repository of the revision marked resolved. */
	Repository LoreRepositoryId
	/* The revision marked resolved. */
	Revision LoreHash
}
type LoreBranchMergeStartBeginEventDataFFI struct {
	/* The source branch being merged. */
	Branch LoreBranchId
	/* The source revision being merged. */
	Revision LoreHash
	/* The sequential number of the source revision. */
	RevisionNumber uint64
}

type LoreBranchMergeStartBeginEventData struct {
	/* The source branch being merged. */
	Branch LoreBranchId
	/* The source revision being merged. */
	Revision LoreHash
	/* The sequential number of the source revision. */
	RevisionNumber uint64
}
type LoreRevisionSyncProgressEventDataFFI struct {
	/* Number of files updated so far. */
	FileUpdate uintptr
	/* Total number of files to update. */
	FileUpdateTotal uintptr
	/* Number of files deleted so far. */
	FileDelete uintptr
	/* Total number of files to delete. */
	FileDeleteTotal uintptr
	/* Number of files merged automatically so far. */
	FileAutomerge uintptr
	/* Number of files with conflicts so far. */
	FileConflict uintptr
	/* Number of bytes updated so far. */
	BytesUpdate uint64
	/* Total number of bytes to update. */
	BytesUpdateTotal uint64
	/* Flag indicating discovery of the work to do has finished. */
	DiscoveryComplete uint8
}

type LoreRevisionSyncProgressEventData struct {
	/* Number of files updated so far. */
	FileUpdate uintptr
	/* Total number of files to update. */
	FileUpdateTotal uintptr
	/* Number of files deleted so far. */
	FileDelete uintptr
	/* Total number of files to delete. */
	FileDeleteTotal uintptr
	/* Number of files merged automatically so far. */
	FileAutomerge uintptr
	/* Number of files with conflicts so far. */
	FileConflict uintptr
	/* Number of bytes updated so far. */
	BytesUpdate uint64
	/* Total number of bytes to update. */
	BytesUpdateTotal uint64
	/* Flag indicating discovery of the work to do has finished. */
	DiscoveryComplete bool
}
type LoreBranchMergeStartEndEventDataFFI struct {
	/* Progress totals collected while applying the merge. */
	Stats LoreRevisionSyncProgressEventDataFFI
	/* The revision produced by the merge. */
	Signature LoreHash
	/* Set when the merge produced file conflicts. */
	HasConflicts uint8
}

type LoreBranchMergeStartEndEventData struct {
	/* Progress totals collected while applying the merge. */
	Stats LoreRevisionSyncProgressEventData
	/* The revision produced by the merge. */
	Signature LoreHash
	/* Set when the merge produced file conflicts. */
	HasConflicts bool
}
type LoreCherryPickStartBeginEventDataFFI struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Identifier of the revision being cherry-picked. */
	Revision LoreHash
	/* Number of the revision being cherry-picked. */
	RevisionNumber uint64
}

type LoreCherryPickStartBeginEventData struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Identifier of the revision being cherry-picked. */
	Revision LoreHash
	/* Number of the revision being cherry-picked. */
	RevisionNumber uint64
}
type LoreCherryPickStartEndEventDataFFI struct {
	/* Progress statistics for the applied changes. */
	Stats LoreRevisionSyncProgressEventDataFFI
	/* Resulting revision hash signature. */
	Signature LoreHash
	/* Flag indicating the cherry-pick produced conflicts. */
	HasConflicts uint8
}

type LoreCherryPickStartEndEventData struct {
	/* Progress statistics for the applied changes. */
	Stats LoreRevisionSyncProgressEventData
	/* Resulting revision hash signature. */
	Signature LoreHash
	/* Flag indicating the cherry-pick produced conflicts. */
	HasConflicts bool
}
type LoreCherryPickAbortBeginEventDataFFI struct {
	/* Identifier of the staged revision being discarded. */
	StateStagedRevision LoreHash
	/* Identifier of the current revision being restored. */
	StateCurrentRevision LoreHash
}

type LoreCherryPickAbortBeginEventData struct {
	/* Identifier of the staged revision being discarded. */
	StateStagedRevision LoreHash
	/* Identifier of the current revision being restored. */
	StateCurrentRevision LoreHash
}
type LoreCherryPickAbortEndEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreCherryPickAbortEndEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreCherryPickConflictFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreCherryPickConflictFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreCherryPickUnresolveFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreCherryPickUnresolveFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreCherryPickUnresolveRevisionEventDataFFI struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}

type LoreCherryPickUnresolveRevisionEventData struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}
type LoreCherryPickResolveFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreCherryPickResolveFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreCherryPickResolveRevisionEventDataFFI struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}

type LoreCherryPickResolveRevisionEventData struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}
type LoreRevertStartBeginEventDataFFI struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Identifier of the revision being reverted. */
	Revision LoreHash
	/* Number of the revision being reverted. */
	RevisionNumber uint64
}

type LoreRevertStartBeginEventData struct {
	/* Branch identifier. */
	Branch LoreBranchId
	/* Identifier of the revision being reverted. */
	Revision LoreHash
	/* Number of the revision being reverted. */
	RevisionNumber uint64
}
type LoreRevertStartEndEventDataFFI struct {
	/* Progress statistics for the applied changes. */
	Stats LoreRevisionSyncProgressEventDataFFI
	/* Resulting revision hash signature. */
	Signature LoreHash
	/* Flag indicating the revert produced conflicts. */
	HasConflicts uint8
}

type LoreRevertStartEndEventData struct {
	/* Progress statistics for the applied changes. */
	Stats LoreRevisionSyncProgressEventData
	/* Resulting revision hash signature. */
	Signature LoreHash
	/* Flag indicating the revert produced conflicts. */
	HasConflicts bool
}
type LoreRevertAbortBeginEventDataFFI struct {
	/* Identifier of the staged revision being discarded. */
	StateStagedRevision LoreHash
	/* Identifier of the current revision being restored. */
	StateCurrentRevision LoreHash
}

type LoreRevertAbortBeginEventData struct {
	/* Identifier of the staged revision being discarded. */
	StateStagedRevision LoreHash
	/* Identifier of the current revision being restored. */
	StateCurrentRevision LoreHash
}
type LoreRevertAbortEndEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreRevertAbortEndEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreRevertResolveFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreRevertResolveFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreRevertResolveRevisionEventDataFFI struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}

type LoreRevertResolveRevisionEventData struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}
type LoreRevertConflictFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreRevertConflictFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreRevertUnresolveFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
}

type LoreRevertUnresolveFileEventData struct {
	/* Path of the file. */
	Path string
}
type LoreRevertUnresolveRevisionEventDataFFI struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}

type LoreRevertUnresolveRevisionEventData struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Identifier of the revision. */
	Revision LoreHash
}
type LoreBranchProtectEventDataFFI struct {
	/* Name of the protected branch. */
	Name LoreString
}

type LoreBranchProtectEventData struct {
	/* Name of the protected branch. */
	Name string
}
type LoreBranchPushEventDataFFI struct {
	/* The remote being pushed to. */
	Remote LoreString
	/* The repository being pushed. */
	Repository LoreRepositoryId
	/* The branch being pushed. */
	Branch LoreBranchId
	/* The name of the branch being pushed. */
	BranchName LoreString
	/* The latest revision of the branch on the remote. */
	RemoteRevision LoreHash
	/* The latest revision of the branch in the local repository. */
	LocalRevision LoreHash
	/* The number of revisions on the remote that are not present locally. */
	RemoteHistory uint64
	/* The number of local revisions to push. */
	LocalHistory uint64
	/* Set when the local revision is already present on the remote. */
	FlagAlreadyPushed uint8
	/* Set when the branch is the repository's default branch. */
	FlagDefault uint8
	/* Set when the repository is a linked repository. */
	FlagLink uint8
	/* Set when the repository is a layer. */
	FlagLayer uint8
}

type LoreBranchPushEventData struct {
	/* The remote being pushed to. */
	Remote string
	/* The repository being pushed. */
	Repository LoreRepositoryId
	/* The branch being pushed. */
	Branch LoreBranchId
	/* The name of the branch being pushed. */
	BranchName string
	/* The latest revision of the branch on the remote. */
	RemoteRevision LoreHash
	/* The latest revision of the branch in the local repository. */
	LocalRevision LoreHash
	/* The number of revisions on the remote that are not present locally. */
	RemoteHistory uint64
	/* The number of local revisions to push. */
	LocalHistory uint64
	/* Set when the local revision is already present on the remote. */
	FlagAlreadyPushed bool
	/* Set when the branch is the repository's default branch. */
	FlagDefault bool
	/* Set when the repository is a linked repository. */
	FlagLink bool
	/* Set when the repository is a layer. */
	FlagLayer bool
}
type LoreBranchPushRevisionUpdateBeginEventDataFFI struct {
	/* The revision being updated. */
	Revision LoreHash
	/* The previous parent revision. */
	OldParent LoreHash
	/* The new parent revision. */
	NewParent LoreHash
}

type LoreBranchPushRevisionUpdateBeginEventData struct {
	/* The revision being updated. */
	Revision LoreHash
	/* The previous parent revision. */
	OldParent LoreHash
	/* The new parent revision. */
	NewParent LoreHash
}
type LoreBranchPushRevisionUpdateEndEventDataFFI struct {
	/* The updated revision. */
	Revision LoreHash
}

type LoreBranchPushRevisionUpdateEndEventData struct {
	/* The updated revision. */
	Revision LoreHash
}
type LoreBranchPushFragmentBeginEventDataFFI struct {
	/* The number of fragments to transfer. */
	Fragments uint64
	/* The total number of bytes to transfer. */
	BytesTotal uint64
}

type LoreBranchPushFragmentBeginEventData struct {
	/* The number of fragments to transfer. */
	Fragments uint64
	/* The total number of bytes to transfer. */
	BytesTotal uint64
}
type LoreBranchPushFragmentProgressEventDataFFI struct {
	/* The number of fragments transferred so far. */
	Complete uint64
	/* The total number of fragments to transfer. */
	Count uint64
	/* The number of bytes transferred so far. */
	BytesTransferred uint64
	/* The total number of bytes to transfer. */
	BytesTotal uint64
}

type LoreBranchPushFragmentProgressEventData struct {
	/* The number of fragments transferred so far. */
	Complete uint64
	/* The total number of fragments to transfer. */
	Count uint64
	/* The number of bytes transferred so far. */
	BytesTransferred uint64
	/* The total number of bytes to transfer. */
	BytesTotal uint64
}
type LoreBranchPushFragmentEndEventDataFFI struct {
	/* The number of fragments transferred. */
	Fragments uint64
	/* The number of bytes transferred. */
	BytesTransferred uint64
}

type LoreBranchPushFragmentEndEventData struct {
	/* The number of fragments transferred. */
	Fragments uint64
	/* The number of bytes transferred. */
	BytesTransferred uint64
}
type LoreBranchPushBranchCreateBeginEventDataFFI struct {
	/* The local revision the branch starts from. */
	LocalRevision LoreHash
}

type LoreBranchPushBranchCreateBeginEventData struct {
	/* The local revision the branch starts from. */
	LocalRevision LoreHash
}
type LoreBranchPushBranchCreateEndEventDataFFI struct {
	/* The revision the branch points to on the remote. */
	RemoteRevision LoreHash
}

type LoreBranchPushBranchCreateEndEventData struct {
	/* The revision the branch points to on the remote. */
	RemoteRevision LoreHash
}
type LoreBranchPushRevisionPushBeginEventDataFFI struct {
	/* The latest revision of the branch on the remote. */
	RemoteRevision LoreHash
	/* The local revision being pushed. */
	LocalRevision LoreHash
}

type LoreBranchPushRevisionPushBeginEventData struct {
	/* The latest revision of the branch on the remote. */
	RemoteRevision LoreHash
	/* The local revision being pushed. */
	LocalRevision LoreHash
}
type LoreBranchPushRevisionPushUpdateEventDataFFI struct {
	/* The revision before the remote reassigned it. */
	OldRevision LoreHash
	/* The revision the remote assigned. */
	NewRevision LoreHash
	/* The sequential number of the new revision. */
	NewRevisionNumber uint64
}

type LoreBranchPushRevisionPushUpdateEventData struct {
	/* The revision before the remote reassigned it. */
	OldRevision LoreHash
	/* The revision the remote assigned. */
	NewRevision LoreHash
	/* The sequential number of the new revision. */
	NewRevisionNumber uint64
}
type LoreBranchPushRevisionPushEndEventDataFFI struct {
	/* The branch revision on the remote before the push. */
	OldRemoteRevision LoreHash
	/* The branch revision on the remote after the push. */
	NewRemoteRevision LoreHash
	/* The sequential number of the new remote revision. */
	NewRemoteRevisionNumber uint64
	/* A message returned by the remote for the push. */
	Message LoreString
	/* Set when the remote performed a fast-forward merge. */
	FastForwardMerged uint8
}

type LoreBranchPushRevisionPushEndEventData struct {
	/* The branch revision on the remote before the push. */
	OldRemoteRevision LoreHash
	/* The branch revision on the remote after the push. */
	NewRemoteRevision LoreHash
	/* The sequential number of the new remote revision. */
	NewRemoteRevisionNumber uint64
	/* A message returned by the remote for the push. */
	Message string
	/* Set when the remote performed a fast-forward merge. */
	FastForwardMerged bool
}
type LoreBranchResetEventDataFFI struct {
	/* Branch identifier. */
	Id LoreContext
	/* Branch name. */
	Name LoreString
	/* Revision the branch was reset to. */
	Revision LoreHash
}

type LoreBranchResetEventData struct {
	/* Branch identifier. */
	Id LoreContext
	/* Branch name. */
	Name string
	/* Revision the branch was reset to. */
	Revision LoreHash
}
type LoreBranchSwitchBeginEventDataFFI struct {
	/* Details of the branch being switched to. */
	Branch LoreBranchSwitchDataFFI
}

type LoreBranchSwitchBeginEventData struct {
	/* Details of the branch being switched to. */
	Branch LoreBranchSwitchData
}
type LoreBranchSwitchEndEventDataFFI struct {
	/* Details of the branch that was switched to. */
	Branch LoreBranchSwitchDataFFI
}

type LoreBranchSwitchEndEventData struct {
	/* Details of the branch that was switched to. */
	Branch LoreBranchSwitchData
}
type LoreBranchUnprotectEventDataFFI struct {
	/* Name of the unprotected branch. */
	Name LoreString
}

type LoreBranchUnprotectEventData struct {
	/* Name of the unprotected branch. */
	Name string
}
type LoreFileInfoEventDataFFI struct {
	/* Path of the file or directory. */
	Path LoreString
	/* Context identifying the file or directory. */
	Context LoreContext
	/* Content hash of the file or directory. */
	Hash LoreHash
	/* Set when the entry is a file. */
	IsFile uint8
	/* Set when the entry is a directory. */
	IsDir uint8
	/* Set when the entry has been modified. */
	FlagModified uint8
	/* Set when the entry has been deleted. */
	FlagDeleted uint8
	/* Set when the entry has been added. */
	FlagAdded uint8
	/* Set when the entry is in conflict. */
	FlagConflict uint8
	/* File mode bits. */
	Mode uint16
	/* Size of the entry in the repository, in bytes. */
	Size uint64
	/* Size of the entry on the local filesystem, in bytes. */
	LocalSize uint64
	/* Content hash of the entry on the local filesystem. */
	LocalHash LoreHash
	/* Size of the entry after filters are applied, in bytes. */
	FilterSize uint64
}

type LoreFileInfoEventData struct {
	/* Path of the file or directory. */
	Path string
	/* Context identifying the file or directory. */
	Context LoreContext
	/* Content hash of the file or directory. */
	Hash LoreHash
	/* Set when the entry is a file. */
	IsFile bool
	/* Set when the entry is a directory. */
	IsDir bool
	/* Set when the entry has been modified. */
	FlagModified bool
	/* Set when the entry has been deleted. */
	FlagDeleted bool
	/* Set when the entry has been added. */
	FlagAdded bool
	/* Set when the entry is in conflict. */
	FlagConflict bool
	/* File mode bits. */
	Mode uint16
	/* Size of the entry in the repository, in bytes. */
	Size uint64
	/* Size of the entry on the local filesystem, in bytes. */
	LocalSize uint64
	/* Content hash of the entry on the local filesystem. */
	LocalHash LoreHash
	/* Size of the entry after filters are applied, in bytes. */
	FilterSize uint64
}
type LoreFileDiffEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Unified-diff text describing the change. */
	Patch LoreString
	/* Action applied to the file. */
	Action LoreFileAction
}

type LoreFileDiffEventData struct {
	/* Path of the file. */
	Path string
	/* Unified-diff text describing the change. */
	Patch string
	/* Action applied to the file. */
	Action LoreFileAction
}
type LoreFileHashEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Size of the file in bytes. */
	Size uint64
	/* Content hash of the file. */
	Hash LoreHash
}

type LoreFileHashEventData struct {
	/* Path of the file. */
	Path string
	/* Size of the file in bytes. */
	Size uint64
	/* Content hash of the file. */
	Hash LoreHash
}
type LoreFileHistoryEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision this entry belongs to. */
	Revision LoreHash
	/* Sequential number of the revision. */
	RevisionNumber uint64
	/* Parent revisions of this revision. */
	Parent [2]LoreHash
	/* Address of the file content at this revision. */
	Address LoreAddress
	/* Size of the file in bytes at this revision. */
	Size uint64
	/* Action applied to the file at this revision. */
	Action LoreFileAction
}

type LoreFileHistoryEventData struct {
	/* Path of the file. */
	Path string
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision this entry belongs to. */
	Revision LoreHash
	/* Sequential number of the revision. */
	RevisionNumber uint64
	/* Parent revisions of this revision. */
	Parent [2]LoreHash
	/* Address of the file content at this revision. */
	Address LoreAddress
	/* Size of the file in bytes at this revision. */
	Size uint64
	/* Action applied to the file at this revision. */
	Action LoreFileAction
}
type LoreFileWriteEventDataFFI struct {
	/* Path that was written. */
	Path LoreString
}

type LoreFileWriteEventData struct {
	/* Path that was written. */
	Path string
}
type LoreFileObliterateEventDataFFI struct {
	/* Address of the obliterated content. */
	Address LoreAddress
	/* Number of fragments removed. */
	NumFragments uintptr
	/* Number of payloads removed. */
	NumPayloads uintptr
}

type LoreFileObliterateEventData struct {
	/* Address of the obliterated content. */
	Address LoreAddress
	/* Number of fragments removed. */
	NumFragments uintptr
	/* Number of payloads removed. */
	NumPayloads uintptr
}
type LoreFileDumpEventDataFFI struct {
	/* Address of the content. */
	Address LoreAddress
	/* Flags describing the stored content. */
	Flags uint32
	/* Size of the stored payload in bytes. */
	SizePayload uint32
	/* Size of the content in bytes. */
	SizeContent uint64
	/* Set when a matching stored object was found. */
	MatchMade uint8
}

type LoreFileDumpEventData struct {
	/* Address of the content. */
	Address LoreAddress
	/* Flags describing the stored content. */
	Flags uint32
	/* Size of the stored payload in bytes. */
	SizePayload uint32
	/* Size of the content in bytes. */
	SizeContent uint64
	/* Set when a matching stored object was found. */
	MatchMade bool
}
type LoreFileDependencyAddBeginEventDataFFI struct {
	/* Number of source files being processed. */
	PathCount uint64
	/* Number of dependency edges being added. */
	DependencyCount uint64
}

type LoreFileDependencyAddBeginEventData struct {
	/* Number of source files being processed. */
	PathCount uint64
	/* Number of dependency edges being added. */
	DependencyCount uint64
}
type LoreFileDependencyAddEntryEventDataFFI struct {
	/* Path of the source file that gains the dependency. */
	Path LoreString
	/* Path of the file being depended on. */
	Dependency LoreString
	/* Tags applied to this dependency edge. */
	Tags LoreStringArrayFFI
}

type LoreFileDependencyAddEntryEventData struct {
	/* Path of the source file that gains the dependency. */
	Path string
	/* Path of the file being depended on. */
	Dependency string
	/* Tags applied to this dependency edge. */
	Tags []string
}
type LoreFileDependencyAddEndEventDataFFI struct {
	/* Number of dependency edges that were added. */
	AddedCount uint64
}

type LoreFileDependencyAddEndEventData struct {
	/* Number of dependency edges that were added. */
	AddedCount uint64
}
type LoreFileDependencyRemoveBeginEventDataFFI struct {
	/* Number of source files being processed. */
	PathCount uint64
	/* Number of dependency edges being removed. */
	DependencyCount uint64
}

type LoreFileDependencyRemoveBeginEventData struct {
	/* Number of source files being processed. */
	PathCount uint64
	/* Number of dependency edges being removed. */
	DependencyCount uint64
}
type LoreFileDependencyRemoveEntryEventDataFFI struct {
	/* Path of the source file that loses the dependency. */
	Path LoreString
	/* Path of the file that was depended on. */
	Dependency LoreString
	/* Tags on the dependency edge being removed. */
	Tags LoreStringArrayFFI
}

type LoreFileDependencyRemoveEntryEventData struct {
	/* Path of the source file that loses the dependency. */
	Path string
	/* Path of the file that was depended on. */
	Dependency string
	/* Tags on the dependency edge being removed. */
	Tags []string
}
type LoreFileDependencyRemoveEndEventDataFFI struct {
	/* Number of dependency edges that were removed. */
	RemovedCount uint64
}

type LoreFileDependencyRemoveEndEventData struct {
	/* Number of dependency edges that were removed. */
	RemovedCount uint64
}
type LoreFileDependencyListBeginEventDataFFI struct {
	/* Number of files being listed. */
	FileCount uint64
}

type LoreFileDependencyListBeginEventData struct {
	/* Number of files being listed. */
	FileCount uint64
}
type LoreFileDependencyListFileEventDataFFI struct {
	/* Path of the file whose dependencies are being listed. */
	Path LoreString
	/* Number of dependency entries for this file. */
	EntryCount uint64
}

type LoreFileDependencyListFileEventData struct {
	/* Path of the file whose dependencies are being listed. */
	Path string
	/* Number of dependency entries for this file. */
	EntryCount uint64
}
type LoreFileDependencyListEntryEventDataFFI struct {
	/* Path of the dependency. */
	Path LoreString
	/* Tags on this dependency edge. */
	Tags LoreStringArrayFFI
	/* Traversal depth, zero for a direct dependency. */
	Depth uint32
}

type LoreFileDependencyListEntryEventData struct {
	/* Path of the dependency. */
	Path string
	/* Tags on this dependency edge. */
	Tags []string
	/* Traversal depth, zero for a direct dependency. */
	Depth uint32
}
type LoreFileDependencyListFileEndEventDataFFI struct {
	/* Path of the file whose dependencies were listed. */
	Path LoreString
}

type LoreFileDependencyListFileEndEventData struct {
	/* Path of the file whose dependencies were listed. */
	Path string
}
type LoreFileDependencyListEndEventDataFFI struct {
	/* Total number of dependency entries that were listed. */
	TotalEntryCount uint64
}

type LoreFileDependencyListEndEventData struct {
	/* Total number of dependency entries that were listed. */
	TotalEntryCount uint64
}
type LoreFileResetBeginEventDataFFI struct {
	/* Number of paths requested for reset. */
	PathCount uintptr
}

type LoreFileResetBeginEventData struct {
	/* Number of paths requested for reset. */
	PathCount uintptr
}
type LoreFileResetProgressEventDataFFI struct {
	/* Current counts of items processed. */
	Count LoreFileResetCountDataFFI
}

type LoreFileResetProgressEventData struct {
	/* Current counts of items processed. */
	Count LoreFileResetCountData
}
type LoreFileResetEndEventDataFFI struct {
	/* Final counts of items processed. */
	Count LoreFileResetCountDataFFI
}

type LoreFileResetEndEventData struct {
	/* Final counts of items processed. */
	Count LoreFileResetCountData
}
type LoreFileResetFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Action applied to the file. */
	Action LoreFileAction
	/* Previous path of the file, when it was moved. */
	FromPath LoreString
}

type LoreFileResetFileEventData struct {
	/* Path of the file. */
	Path string
	/* Action applied to the file. */
	Action LoreFileAction
	/* Previous path of the file, when it was moved. */
	FromPath string
}
type LoreFilterExcludeEventDataFFI struct {
	/* Reason the path was excluded. */
	Reason uint8
	/* Path that was excluded. */
	Path LoreString
}

type LoreFilterExcludeEventData struct {
	/* Reason the path was excluded. */
	Reason bool
	/* Path that was excluded. */
	Path string
}
type LoreFileStageBeginEventDataFFI struct {
	/* Number of paths requested for staging. */
	PathCount uintptr
}

type LoreFileStageBeginEventData struct {
	/* Number of paths requested for staging. */
	PathCount uintptr
}
type LoreFileStageProgressEventDataFFI struct {
	/* Current counts of items processed. */
	Count LoreFileStageCountDataFFI
}

type LoreFileStageProgressEventData struct {
	/* Current counts of items processed. */
	Count LoreFileStageCountData
}
type LoreFileStageEndEventDataFFI struct {
	/* Final counts of items processed. */
	Count LoreFileStageCountDataFFI
}

type LoreFileStageEndEventData struct {
	/* Final counts of items processed. */
	Count LoreFileStageCountData
}
type LoreFileStageRevisionEventDataFFI struct {
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision the files are staged against. */
	Revision LoreHash
}

type LoreFileStageRevisionEventData struct {
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision the files are staged against. */
	Revision LoreHash
}
type LoreFileStageFileEventDataFFI struct {
	/* Previous path of the file, when it was moved. */
	FromPath LoreString
	/* Path of the file. */
	Path LoreString
	/* Action applied to the file. */
	Action LoreFileAction
}

type LoreFileStageFileEventData struct {
	/* Previous path of the file, when it was moved. */
	FromPath string
	/* Path of the file. */
	Path string
	/* Action applied to the file. */
	Action LoreFileAction
}
type LoreFileUnstageBeginEventDataFFI struct {
	/* Number of paths requested for unstaging. */
	PathCount uintptr
}

type LoreFileUnstageBeginEventData struct {
	/* Number of paths requested for unstaging. */
	PathCount uintptr
}
type LoreFileUnstageProgressEventDataFFI struct {
	/* Current counts of items processed. */
	Count LoreFileUnstageCountDataFFI
}

type LoreFileUnstageProgressEventData struct {
	/* Current counts of items processed. */
	Count LoreFileUnstageCountData
}
type LoreFileUnstageEndEventDataFFI struct {
	/* Final counts of items processed. */
	Count LoreFileUnstageCountDataFFI
}

type LoreFileUnstageEndEventData struct {
	/* Final counts of items processed. */
	Count LoreFileUnstageCountData
}
type LoreFileUnstageRevisionEventDataFFI struct {
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision the files are unstaged against. */
	Revision LoreHash
}

type LoreFileUnstageRevisionEventData struct {
	/* Identifier of the repository. */
	Repository LoreRepositoryId
	/* Revision the files are unstaged against. */
	Revision LoreHash
}
type LoreFileUnstageFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Action applied to the file. */
	Action LoreFileAction
}

type LoreFileUnstageFileEventData struct {
	/* Path of the file. */
	Path string
	/* Action applied to the file. */
	Action LoreFileAction
}
type LoreFragmentWriteEventDataFFI struct {
	/* The fragment that was written */
	Fragment LoreFragmentFFI
	/* Non-zero if the fragment already existed and was deduplicated */
	Deduplicated uint8
}

type LoreFragmentWriteEventData struct {
	/* The fragment that was written */
	Fragment LoreFragment
	/* Non-zero if the fragment already existed and was deduplicated */
	Deduplicated bool
}
type LoreLayerAddEventDataFFI struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath LoreString
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer starts. */
	SourcePath LoreString
	/* Metadata used to match revisions between the repositories. */
	Metadata LoreString
	/* Revision of the source repository. */
	Revision LoreHash
}

type LoreLayerAddEventData struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath string
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer starts. */
	SourcePath string
	/* Metadata used to match revisions between the repositories. */
	Metadata string
	/* Revision of the source repository. */
	Revision LoreHash
}
type LoreLayerEntryEventDataFFI struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath LoreString
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer starts. */
	SourcePath LoreString
	/* Metadata used to match revisions between the repositories. */
	Metadata LoreString
	/* Revision of the source repository. */
	Revision LoreHash
}

type LoreLayerEntryEventData struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath string
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer starts. */
	SourcePath string
	/* Metadata used to match revisions between the repositories. */
	Metadata string
	/* Revision of the source repository. */
	Revision LoreHash
}
type LoreLayerRemoveEventDataFFI struct {
	/* Path in the outer repository where the layer was placed. */
	TargetPath LoreString
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer started. */
	SourcePath LoreString
	/* Revision of the source repository. */
	Revision LoreHash
	/* Set when removal was forced. */
	Forced uint8
	/* Set when the layer files were purged from disk. */
	Purged uint8
	/* Number of files removed. */
	FileCount uint64
	/* Number of directories removed. */
	DirectoryCount uint64
	/* Number of modified files encountered. */
	ModifiedCount uint64
}

type LoreLayerRemoveEventData struct {
	/* Path in the outer repository where the layer was placed. */
	TargetPath string
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Path inside the source repository where the layer started. */
	SourcePath string
	/* Revision of the source repository. */
	Revision LoreHash
	/* Set when removal was forced. */
	Forced bool
	/* Set when the layer files were purged from disk. */
	Purged bool
	/* Number of files removed. */
	FileCount uint64
	/* Number of directories removed. */
	DirectoryCount uint64
	/* Number of modified files encountered. */
	ModifiedCount uint64
}
type LoreLayerStagedEntryEventDataFFI struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath LoreString
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Number of staged files in the layer. */
	StagedFileCount uint64
}

type LoreLayerStagedEntryEventData struct {
	/* Path in the outer repository where the layer is placed. */
	TargetPath string
	/* Identifier of the source repository. */
	SourceRepository LoreRepositoryId
	/* Number of staged files in the layer. */
	StagedFileCount uint64
}
type LoreLinkChangeEventDataFFI struct {
	/* Path of the link within the parent repository. */
	LinkPath LoreString
	/* Identifier of the repository the link points to. */
	LinkRepository LoreRepositoryId
	/* Identifier of the branch the link is pinned to. */
	Branch LoreBranchId
	/* Hash of the revision the link is pinned to. */
	Revision LoreHash
	/* Kind of change applied to the link. */
	Action LoreFileAction
}

type LoreLinkChangeEventData struct {
	/* Path of the link within the parent repository. */
	LinkPath string
	/* Identifier of the repository the link points to. */
	LinkRepository LoreRepositoryId
	/* Identifier of the branch the link is pinned to. */
	Branch LoreBranchId
	/* Hash of the revision the link is pinned to. */
	Revision LoreHash
	/* Kind of change applied to the link. */
	Action LoreFileAction
}
type LoreLinkEntryEventDataFFI struct {
	/* Identifier of the repository the link points to. */
	Link LoreRepositoryId
	/* Identifier of the link node in the parent repository. */
	LinkNode uint32
	/* Path of the link within the parent repository. */
	LinkPath LoreString
	/* Identifier of the source node in the linked repository. */
	SourceNode uint32
	/* Path of the source within the linked repository. */
	SourcePath LoreString
	/* Identifier of the branch the link is pinned to. */
	Branch LoreBranchId
	/* Name of the branch the link is pinned to. */
	BranchName LoreString
	/* Hash of the revision the link is pinned to. */
	Revision LoreHash
	/* Link flags. */
	Flags uint32
}

type LoreLinkEntryEventData struct {
	/* Identifier of the repository the link points to. */
	Link LoreRepositoryId
	/* Identifier of the link node in the parent repository. */
	LinkNode uint32
	/* Path of the link within the parent repository. */
	LinkPath string
	/* Identifier of the source node in the linked repository. */
	SourceNode uint32
	/* Path of the source within the linked repository. */
	SourcePath string
	/* Identifier of the branch the link is pinned to. */
	Branch LoreBranchId
	/* Name of the branch the link is pinned to. */
	BranchName string
	/* Hash of the revision the link is pinned to. */
	Revision LoreHash
	/* Link flags. */
	Flags uint32
}
type LoreLockFileAcquireBeginEventDataFFI struct {
	/* Number of acquire entries that follow. */
	Count uint64
	/* Whether the entries that follow were already owned. */
	Ignored uint8
}

type LoreLockFileAcquireBeginEventData struct {
	/* Number of acquire entries that follow. */
	Count uint64
	/* Whether the entries that follow were already owned. */
	Ignored bool
}
type LoreLockFileAcquireEventDataFFI struct {
	/* The path whose lock is being acquired. */
	Path LoreString
}

type LoreLockFileAcquireEventData struct {
	/* The path whose lock is being acquired. */
	Path string
}
type LoreLockFileStatusBeginEventDataFFI struct {
	/* Number of status entries that follow. */
	Count uint64
}

type LoreLockFileStatusBeginEventData struct {
	/* Number of status entries that follow. */
	Count uint64
}
type LoreLockFileStatusEventDataFFI struct {
	/* Path the status applies to. */
	Path LoreString
	/* Identifier of the user that holds the lock. */
	Owner LoreString
	/* Timestamp recorded when the lock was acquired. */
	LockedAt uint64
}

type LoreLockFileStatusEventData struct {
	/* Path the status applies to. */
	Path string
	/* Identifier of the user that holds the lock. */
	Owner string
	/* Timestamp recorded when the lock was acquired. */
	LockedAt uint64
}
type LoreLockFileQueryBeginEventDataFFI struct {
	/* Number of query entries that follow. */
	Count uint64
}

type LoreLockFileQueryBeginEventData struct {
	/* Number of query entries that follow. */
	Count uint64
}
type LoreLockFileQueryEventDataFFI struct {
	/* Identifier of the branch the lock belongs to. */
	Branch LoreBranchId
	/* Path the lock applies to. */
	Path LoreString
	/* Identifier of the user that holds the lock. */
	Owner LoreString
	/* Timestamp recorded when the lock was acquired. */
	LockedAt uint64
}

type LoreLockFileQueryEventData struct {
	/* Identifier of the branch the lock belongs to. */
	Branch LoreBranchId
	/* Path the lock applies to. */
	Path string
	/* Identifier of the user that holds the lock. */
	Owner string
	/* Timestamp recorded when the lock was acquired. */
	LockedAt uint64
}
type LoreLockFileReleaseBeginEventDataFFI struct {
	/* Number of release entries that follow. */
	Count uint64
	/* Whether no matching lock was found to release. */
	NotFound uint8
}

type LoreLockFileReleaseBeginEventData struct {
	/* Number of release entries that follow. */
	Count uint64
	/* Whether no matching lock was found to release. */
	NotFound bool
}
type LoreLockFileReleaseEventDataFFI struct {
	/* The path whose lock is being released. */
	Path LoreString
}

type LoreLockFileReleaseEventData struct {
	/* The path whose lock is being released. */
	Path string
}
type LoreMetadataClearFileEventDataFFI struct {
	/* Path of the file whose metadata was cleared. */
	Path LoreString
}

type LoreMetadataClearFileEventData struct {
	/* Path of the file whose metadata was cleared. */
	Path string
}
type LoreMetadataClearRevisionEventDataFFI struct {
	/* Hash of the revision whose metadata was cleared. */
	Revision LoreHash
}

type LoreMetadataClearRevisionEventData struct {
	/* Hash of the revision whose metadata was cleared. */
	Revision LoreHash
}
type LorePathIgnoreEventDataFFI struct {
	/* The ignored path */
	Path LoreString
}

type LorePathIgnoreEventData struct {
	/* The ignored path */
	Path string
}
type LoreRepositoryCreateEventDataFFI struct {
	/* Identifier of the created repository. */
	Id LoreRepositoryId
	/* Name of the created repository. */
	Name LoreString
	/* Local path of the created repository. */
	Path LoreString
}

type LoreRepositoryCreateEventData struct {
	/* Identifier of the created repository. */
	Id LoreRepositoryId
	/* Name of the created repository. */
	Name string
	/* Local path of the created repository. */
	Path string
}
type LoreRepositoryCloneBeginEventDataFFI struct {
	/* Identifier of the repository being cloned. */
	Repository LoreRepositoryId
	/* Name of the branch being cloned. */
	Branch LoreString
	/* Revision being cloned. */
	Revision LoreHash
	/* Local path the clone is written to. */
	Path LoreString
}

type LoreRepositoryCloneBeginEventData struct {
	/* Identifier of the repository being cloned. */
	Repository LoreRepositoryId
	/* Name of the branch being cloned. */
	Branch string
	/* Revision being cloned. */
	Revision LoreHash
	/* Local path the clone is written to. */
	Path string
}
type LoreRepositoryCloneProgressEventDataFFI struct {
	/* Current progress counts. */
	Count LoreRepositoryCloneCountDataFFI
}

type LoreRepositoryCloneProgressEventData struct {
	/* Current progress counts. */
	Count LoreRepositoryCloneCountData
}
type LoreRepositoryCloneEndEventDataFFI struct {
	/* Name of the branch that was cloned. */
	Branch LoreString
	/* Revision that was cloned. */
	Revision LoreHash
	/* Final progress counts. */
	Count LoreRepositoryCloneCountDataFFI
}

type LoreRepositoryCloneEndEventData struct {
	/* Name of the branch that was cloned. */
	Branch string
	/* Revision that was cloned. */
	Revision LoreHash
	/* Final progress counts. */
	Count LoreRepositoryCloneCountData
}
type LoreDependencyResolveBeginEventDataFFI struct {
	/* Number of root files resolution starts from. */
	RootCount uint64
}

type LoreDependencyResolveBeginEventData struct {
	/* Number of root files resolution starts from. */
	RootCount uint64
}
type LoreDependencyResolveItemEventDataFFI struct {
	/* Path of the file the dependency comes from. */
	Source LoreString
	/* Path of the file the dependency points to. */
	Target LoreString
	/* Tags on this dependency edge. */
	Tags LoreStringArrayFFI
}

type LoreDependencyResolveItemEventData struct {
	/* Path of the file the dependency comes from. */
	Source string
	/* Path of the file the dependency points to. */
	Target string
	/* Tags on this dependency edge. */
	Tags []string
}
type LoreDependencyResolveEndEventDataFFI struct {
	/* Number of dependency edges that were resolved. */
	ResolvedCount uint64
}

type LoreDependencyResolveEndEventData struct {
	/* Number of dependency edges that were resolved. */
	ResolvedCount uint64
}
type LoreRepositoryDataEventDataFFI struct {
	/* Remote URL of the repository. */
	RemoteUrl LoreString
	/* Repository identifier. */
	Id LoreRepositoryId
	/* Repository name. */
	Name LoreString
	/* Repository description. */
	Description LoreString
	/* Identifier of the default branch. */
	DefaultBranch LoreBranchId
	/* Name of the default branch. */
	DefaultBranchName LoreString
	/* Name of the user who created the repository. */
	Creator LoreString
	/* Creation time of the repository, in seconds since the Unix epoch. */
	Created uint64
}

type LoreRepositoryDataEventData struct {
	/* Remote URL of the repository. */
	RemoteUrl string
	/* Repository identifier. */
	Id LoreRepositoryId
	/* Repository name. */
	Name string
	/* Repository description. */
	Description string
	/* Identifier of the default branch. */
	DefaultBranch LoreBranchId
	/* Name of the default branch. */
	DefaultBranchName string
	/* Name of the user who created the repository. */
	Creator string
	/* Creation time of the repository, in seconds since the Unix epoch. */
	Created uint64
}
type LoreRepositoryConfigGetEventDataFFI struct {
	/* Configuration key. */
	Key LoreString
	/* Configuration value for the key. */
	Value LoreString
}

type LoreRepositoryConfigGetEventData struct {
	/* Configuration key. */
	Key string
	/* Configuration value for the key. */
	Value string
}
type LoreRepositoryDumpBeginEventDataFFI struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Revision being dumped. */
	Revision LoreHash
}

type LoreRepositoryDumpBeginEventData struct {
	/* Repository identifier. */
	Repository LoreRepositoryId
	/* Revision being dumped. */
	Revision LoreHash
}
type LoreRepositoryDumpEndEventDataFFI struct {
	/* Placeholder field. The event carries no data. */
	Unused uint32
}

type LoreRepositoryDumpEndEventData struct {
	/* Placeholder field. The event carries no data. */
	Unused uint32
}
type LoreRepositoryListEntryEventDataFFI struct {
	/* Repository identifier. */
	Id LoreRepositoryId
	/* Repository name. */
	Name LoreString
}

type LoreRepositoryListEntryEventData struct {
	/* Repository identifier. */
	Id LoreRepositoryId
	/* Repository name. */
	Name string
}
type LoreRepositoryInstanceEventDataFFI struct {
	/* Identifier of the instance */
	InstanceId LoreInstanceId
	/* Filesystem path of the instance */
	Path LoreString
	/* Name of the branch the instance has checked out */
	BranchName LoreString
	/* Identifier of the branch the instance has checked out */
	Branch LoreBranchId
	/* Current revision hash for the instance */
	Revision LoreHash
	/* Non-zero if the instance path no longer exists on disk */
	Stale uint8
}

type LoreRepositoryInstanceEventData struct {
	/* Identifier of the instance */
	InstanceId LoreInstanceId
	/* Filesystem path of the instance */
	Path string
	/* Name of the branch the instance has checked out */
	BranchName string
	/* Identifier of the branch the instance has checked out */
	Branch LoreBranchId
	/* Current revision hash for the instance */
	Revision LoreHash
	/* Non-zero if the instance path no longer exists on disk */
	Stale bool
}
type LoreRepositoryVerifyStateBeginEventDataFFI struct {
	/* Placeholder field. The event carries no data. */
	Unused uint32
}

type LoreRepositoryVerifyStateBeginEventData struct {
	/* Placeholder field. The event carries no data. */
	Unused uint32
}
type LoreRepositoryVerifyStateEndEventDataFFI struct {
	/* Identifier of the staged state after healing. Zero when nothing was healed. */
	HealedStagedState LoreHash
}

type LoreRepositoryVerifyStateEndEventData struct {
	/* Identifier of the staged state after healing. Zero when nothing was healed. */
	HealedStagedState LoreHash
}
type LoreRepositoryVerifyFragmentMatchEventDataFFI struct {
	/* Slot the match was found in. */
	Slot uint32
	/* Index of the match within the slot. */
	Index uint32
	/* Identifier of the repository the match belongs to. */
	Repository LoreRepositoryId
	/* Hash part of the fragment address. */
	AddressHash LoreHash
	/* Context part of the fragment address. */
	AddressContext LoreContext
	/* Storage flags recorded for the fragment. */
	Flags uint32
	/* Stored size of the fragment payload in bytes. */
	SizePayload uint32
	/* Size of the fragment content in bytes. */
	SizeContent uint64
	/* Offset of the fragment within its pack file. */
	PackOffset uint32
	/* Index of the pack file holding the fragment. */
	PackFile uint32
	/* Time the fragment was last accessed, in seconds since the Unix epoch. */
	LastAccess uint64
}

type LoreRepositoryVerifyFragmentMatchEventData struct {
	/* Slot the match was found in. */
	Slot uint32
	/* Index of the match within the slot. */
	Index uint32
	/* Identifier of the repository the match belongs to. */
	Repository LoreRepositoryId
	/* Hash part of the fragment address. */
	AddressHash LoreHash
	/* Context part of the fragment address. */
	AddressContext LoreContext
	/* Storage flags recorded for the fragment. */
	Flags uint32
	/* Stored size of the fragment payload in bytes. */
	SizePayload uint32
	/* Size of the fragment content in bytes. */
	SizeContent uint64
	/* Offset of the fragment within its pack file. */
	PackOffset uint32
	/* Index of the pack file holding the fragment. */
	PackFile uint32
	/* Time the fragment was last accessed, in seconds since the Unix epoch. */
	LastAccess uint64
}
type LoreRepositoryVerifyFragmentEventDataFFI struct {
	/* Hash of the fragment that was verified. */
	Hash LoreHash
	/* Index of the group the fragment belongs to. */
	GroupIndex uint32
	/* Index of the bucket the fragment belongs to. */
	BucketIndex uint32
	/* Path of the index file examined for the fragment. */
	IndexPath LoreString
	/* Number of entries in the index. */
	EntryCount uint32
	/* Number of entries in the pack file. */
	PackfileEntryCount uint32
	/* Number of stored copies found for the fragment. */
	MatchCount uint32
	/* The stored copies found for the fragment. */
	Matches LoreRepositoryVerifyFragmentMatchEventDataArrayFFI
	/* Error message produced during verification. Empty on success. */
	Error LoreString
}

type LoreRepositoryVerifyFragmentEventData struct {
	/* Hash of the fragment that was verified. */
	Hash LoreHash
	/* Index of the group the fragment belongs to. */
	GroupIndex uint32
	/* Index of the bucket the fragment belongs to. */
	BucketIndex uint32
	/* Path of the index file examined for the fragment. */
	IndexPath string
	/* Number of entries in the index. */
	EntryCount uint32
	/* Number of entries in the pack file. */
	PackfileEntryCount uint32
	/* Number of stored copies found for the fragment. */
	MatchCount uint32
	/* The stored copies found for the fragment. */
	Matches LoreRepositoryVerifyFragmentMatchEventDataArray
	/* Error message produced during verification. Empty on success. */
	Error string
}
type LoreRepositoryVerifyFragmentRemoteEventDataFFI struct {
	/* Hash part of the fragment address. */
	AddressHash LoreHash
	/* Context part of the fragment address. */
	AddressContext LoreContext
	/* Non-zero when the fragment was found to be corrupted. */
	Corrupted uint8
	/* Non-zero when the fragment was healed. */
	Healed uint8
	/* Error message produced during verification. Empty on success. */
	Error LoreString
}

type LoreRepositoryVerifyFragmentRemoteEventData struct {
	/* Hash part of the fragment address. */
	AddressHash LoreHash
	/* Context part of the fragment address. */
	AddressContext LoreContext
	/* Non-zero when the fragment was found to be corrupted. */
	Corrupted bool
	/* Non-zero when the fragment was healed. */
	Healed bool
	/* Error message produced during verification. Empty on success. */
	Error string
}
type LoreRepositoryStateDumpEventDataFFI struct {
	/* Sequence number of the revision. */
	RevisionNumber uint64
	/* Hash of the revision. */
	Revision LoreHash
	/* Hash of the state's node tree. */
	TreeHash LoreHash
	/* Size of the node tree in bytes. */
	TreeSize uint64
}

type LoreRepositoryStateDumpEventData struct {
	/* Sequence number of the revision. */
	RevisionNumber uint64
	/* Hash of the revision. */
	Revision LoreHash
	/* Hash of the state's node tree. */
	TreeHash LoreHash
	/* Size of the node tree in bytes. */
	TreeSize uint64
}
type LoreRepositoryStateDumpNodeEventDataFFI struct {
	/* Name of the node. */
	Name LoreString
	/* Identifier of the node. */
	Id uint32
	/* Identifier of the parent node. */
	Parent uint32
	/* Identifier of the next sibling node. */
	Sibling uint32
	/* File mode of the node. */
	Mode uint16
	/* Size of the node's content in bytes. */
	Size uint64
	/* Node flags. */
	Flags uint16
	/* Type-specific detail for the node. */
	TypeData LoreString
}

type LoreRepositoryStateDumpNodeEventData struct {
	/* Name of the node. */
	Name string
	/* Identifier of the node. */
	Id uint32
	/* Identifier of the parent node. */
	Parent uint32
	/* Identifier of the next sibling node. */
	Sibling uint32
	/* File mode of the node. */
	Mode uint16
	/* Size of the node's content in bytes. */
	Size uint64
	/* Node flags. */
	Flags uint16
	/* Type-specific detail for the node. */
	TypeData string
}
type LoreRepositoryStatusRevisionEventDataFFI struct {
	/* Repository identifier */
	Repository LoreRepositoryId
	/* Current branch identifier */
	Branch LoreBranchId
	/* Current branch name */
	BranchName LoreString
	/* Current revision identifier */
	Revision LoreHash
	/* Current revision number */
	RevisionNumber uint64
	/* Staged revision identifier (zero when nothing is staged) */
	RevisionStaged LoreHash
	/* Incoming revision identifier of a pending merge (zero when none) */
	RevisionMerged LoreHash
	/* Last revision merged in from the parent branch (calculated and reported if sync point option is set). */
	RevisionMergedParentBranch LoreHash
	/* Local branch latest revision identifier */
	RevisionLocal LoreHash
	/* Local branch latest revision number */
	RevisionLocalNumber uint64
	/* Remote branch latest revision identifier (zero if unknown, branch not existing on remote or remote not available) */
	RevisionRemote LoreHash
	/* Remote branch latest revision number (zero if corresponding identifier is zero) */
	RevisionRemoteNumber uint64
	/* Local holds revisions not on the remote history line */
	IsLocalAhead uint8
	/* Remote holds revisions not present locally */
	IsRemoteAhead uint8
	/* Remote configured and reachable with a local identity; connectivity only, not authorization */
	RemoteAvailable uint8
	/* Remote revision query returned an authoritative answer, identity is authorized to access the repository */
	RemoteAuthorized uint8
	/* Branch exists on the remote and the query returned a latest revisoin (possibly zero if branch does not exist on remote) */
	RemoteBranchExist uint8
}

type LoreRepositoryStatusRevisionEventData struct {
	/* Repository identifier */
	Repository LoreRepositoryId
	/* Current branch identifier */
	Branch LoreBranchId
	/* Current branch name */
	BranchName string
	/* Current revision identifier */
	Revision LoreHash
	/* Current revision number */
	RevisionNumber uint64
	/* Staged revision identifier (zero when nothing is staged) */
	RevisionStaged LoreHash
	/* Incoming revision identifier of a pending merge (zero when none) */
	RevisionMerged LoreHash
	/* Last revision merged in from the parent branch (calculated and reported if sync point option is set). */
	RevisionMergedParentBranch LoreHash
	/* Local branch latest revision identifier */
	RevisionLocal LoreHash
	/* Local branch latest revision number */
	RevisionLocalNumber uint64
	/* Remote branch latest revision identifier (zero if unknown, branch not existing on remote or remote not available) */
	RevisionRemote LoreHash
	/* Remote branch latest revision number (zero if corresponding identifier is zero) */
	RevisionRemoteNumber uint64
	/* Local holds revisions not on the remote history line */
	IsLocalAhead bool
	/* Remote holds revisions not present locally */
	IsRemoteAhead bool
	/* Remote configured and reachable with a local identity; connectivity only, not authorization */
	RemoteAvailable bool
	/* Remote revision query returned an authoritative answer, identity is authorized to access the repository */
	RemoteAuthorized bool
	/* Branch exists on the remote and the query returned a latest revisoin (possibly zero if branch does not exist on remote) */
	RemoteBranchExist bool
}
type LoreRepositoryStatusFileEventDataFFI struct {
	/* Path of the file relative to the repository root. */
	Path LoreString
	/* Size of the file in bytes. */
	Size uint64
	/* Change applied to the file, such as add, modify, delete, or move. */
	Action LoreFileAction
	/* Kind of node: file, directory, or link. */
	Type LoreNodeType
	/* Non-zero when the change is staged. */
	FlagStaged uint8
	/* Non-zero when the change comes from a merge. */
	FlagMerged uint8
	/* Non-zero when the file is in conflict. */
	FlagConflict uint8
	/* Non-zero when the conflict is not yet resolved. */
	FlagConflictUnresolved uint8
	/* Non-zero when the conflict was resolved automatically. */
	FlagConflictAutomerged uint8
	/* Non-zero when the local side was chosen to resolve the conflict. */
	FlagConflictMine uint8
	/* Non-zero when the incoming side was chosen to resolve the conflict. */
	FlagConflictTheirs uint8
	/* Non-zero when the file differs from the recorded state. */
	FlagDirty uint8
	/* Previous path of the file when it was moved or copied. Empty otherwise. */
	FromPath LoreString
}

type LoreRepositoryStatusFileEventData struct {
	/* Path of the file relative to the repository root. */
	Path string
	/* Size of the file in bytes. */
	Size uint64
	/* Change applied to the file, such as add, modify, delete, or move. */
	Action LoreFileAction
	/* Kind of node: file, directory, or link. */
	Type LoreNodeType
	/* Non-zero when the change is staged. */
	FlagStaged bool
	/* Non-zero when the change comes from a merge. */
	FlagMerged bool
	/* Non-zero when the file is in conflict. */
	FlagConflict bool
	/* Non-zero when the conflict is not yet resolved. */
	FlagConflictUnresolved bool
	/* Non-zero when the conflict was resolved automatically. */
	FlagConflictAutomerged bool
	/* Non-zero when the local side was chosen to resolve the conflict. */
	FlagConflictMine bool
	/* Non-zero when the incoming side was chosen to resolve the conflict. */
	FlagConflictTheirs bool
	/* Non-zero when the file differs from the recorded state. */
	FlagDirty bool
	/* Previous path of the file when it was moved or copied. Empty otherwise. */
	FromPath string
}
type LoreRepositoryStatusCountEventDataFFI struct {
	/* Number of directories in the tree, view-filtered (staged state if
	present, otherwise the current revision) */
	Directories uint64
	/* Number of files in the tree, view-filtered (staged state if present,
	otherwise the current revision) */
	Files uint64
}

type LoreRepositoryStatusCountEventData struct {
	/* Number of directories in the tree, view-filtered (staged state if
	present, otherwise the current revision) */
	Directories uint64
	/* Number of files in the tree, view-filtered (staged state if present,
	otherwise the current revision) */
	Files uint64
}
type LoreRepositoryStatusSummaryEventDataFFI struct {
	/* Number of files added. */
	Adds uint64
	/* Number of files deleted. */
	Deletes uint64
	/* Number of files modified. */
	Modifies uint64
	/* Number of files moved. */
	Moves uint64
	/* Number of files copied. */
	Copies uint64
}

type LoreRepositoryStatusSummaryEventData struct {
	/* Number of files added. */
	Adds uint64
	/* Number of files deleted. */
	Deletes uint64
	/* Number of files modified. */
	Modifies uint64
	/* Number of files moved. */
	Moves uint64
	/* Number of files copied. */
	Copies uint64
}
type LoreRepositoryStoreImmutableQueryEventDataFFI struct {
	/* Address of fragment */
	Address LoreAddress
	/* Remote flag, true if results are from remote store, false if local store */
	Remote uint8
	/* Status, where
	0 = exact address exist
	1 = hash exist in repository
	2 = hash exist in other repository
	3 = hash does not exist */
	Status uint32
	/* Payload flag, true if payload data is present in the store, false if not */
	Payload uint8
	/* Subfragment flag, true if this fragment was a subfragment of the original query, false if not */
	Subfragment uint8
	/* Internal flags */
	Flags uint32
	/* Payload size */
	PayloadSize uint32
	/* Content size */
	ContentSize uint64
}

type LoreRepositoryStoreImmutableQueryEventData struct {
	/* Address of fragment */
	Address LoreAddress
	/* Remote flag, true if results are from remote store, false if local store */
	Remote bool
	/* Status, where
	0 = exact address exist
	1 = hash exist in repository
	2 = hash exist in other repository
	3 = hash does not exist */
	Status uint32
	/* Payload flag, true if payload data is present in the store, false if not */
	Payload bool
	/* Subfragment flag, true if this fragment was a subfragment of the original query, false if not */
	Subfragment bool
	/* Internal flags */
	Flags uint32
	/* Payload size */
	PayloadSize uint32
	/* Content size */
	ContentSize uint64
}
type LoreRevisionCommitBeginEventDataFFI struct {
	/* Unused placeholder field. */
	Unused uint32
}

type LoreRevisionCommitBeginEventData struct {
	/* Unused placeholder field. */
	Unused uint32
}
type LoreRevisionCommitProgressEventDataFFI struct {
	/* Current progress counters. */
	Count LoreRevisionCommitCountDataFFI
}

type LoreRevisionCommitProgressEventData struct {
	/* Current progress counters. */
	Count LoreRevisionCommitCountData
}
type LoreRevisionCommitEndEventDataFFI struct {
	/* Final progress counters. */
	Count LoreRevisionCommitCountDataFFI
}

type LoreRevisionCommitEndEventData struct {
	/* Final progress counters. */
	Count LoreRevisionCommitCountData
}
type LoreRevisionCommitRevisionEventDataFFI struct {
	/* Identifier of the repository the revision belongs to. */
	Repository LoreRepositoryId
	/* Identifier of the branch the revision was committed on. */
	Branch LoreBranchId
	/* Signature of the committed revision. */
	Revision LoreHash
	/* Sequential number of the revision. */
	RevisionNumber uint64
	/* Signature of the first parent revision. */
	Parent LoreHash
	/* Signature of the second parent revision, set for a merge. */
	ParentOther LoreHash
}

type LoreRevisionCommitRevisionEventData struct {
	/* Identifier of the repository the revision belongs to. */
	Repository LoreRepositoryId
	/* Identifier of the branch the revision was committed on. */
	Branch LoreBranchId
	/* Signature of the committed revision. */
	Revision LoreHash
	/* Sequential number of the revision. */
	RevisionNumber uint64
	/* Signature of the first parent revision. */
	Parent LoreHash
	/* Signature of the second parent revision, set for a merge. */
	ParentOther LoreHash
}
type LoreRevisionInfoEventDataFFI struct {
	/* Repository identifier the revision belongs to. */
	Repository LoreRepositoryId
	/* Revision hash signature. */
	Revision LoreHash
	/* Revision number. */
	RevisionNumber uint64
	/* Parent revision hashes; the first is the direct parent and the second
	is the other parent of a merge, or zero when there is none. */
	Parent [2]LoreHash
}

type LoreRevisionInfoEventData struct {
	/* Repository identifier the revision belongs to. */
	Repository LoreRepositoryId
	/* Revision hash signature. */
	Revision LoreHash
	/* Revision number. */
	RevisionNumber uint64
	/* Parent revision hashes; the first is the direct parent and the second
	is the other parent of a merge, or zero when there is none. */
	Parent [2]LoreHash
}
type LoreRevisionInfoDeltaEventDataFFI struct {
	/* Path of the file relative to the repository root. */
	Path LoreString
	/* Size of the file in bytes. */
	Size uint64
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the file content was modified. */
	FlagModify uint8
	/* Flag indicating the change came from a merge. */
	FlagMerged uint8
	/* Flag indicating the entry is a file rather than a directory. */
	FlagFile uint8
}

type LoreRevisionInfoDeltaEventData struct {
	/* Path of the file relative to the repository root. */
	Path string
	/* Size of the file in bytes. */
	Size uint64
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the file content was modified. */
	FlagModify bool
	/* Flag indicating the change came from a merge. */
	FlagMerged bool
	/* Flag indicating the entry is a file rather than a directory. */
	FlagFile bool
}
type LoreRevisionDiffFileEventDataFFI struct {
	/* Path of the file relative to the repository root. */
	Path LoreString
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the entry on the source side is a file rather than a directory. */
	OldIsFile uint8
	/* Flag indicating the entry on the target side is a file rather than a directory. */
	NewIsFile uint8
	/* Address of the file content on the source side. */
	OldAddress LoreAddress
	/* Address of the file content on the target side. */
	NewAddress LoreAddress
}

type LoreRevisionDiffFileEventData struct {
	/* Path of the file relative to the repository root. */
	Path string
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the entry on the source side is a file rather than a directory. */
	OldIsFile bool
	/* Flag indicating the entry on the target side is a file rather than a directory. */
	NewIsFile bool
	/* Address of the file content on the source side. */
	OldAddress LoreAddress
	/* Address of the file content on the target side. */
	NewAddress LoreAddress
}
type LoreRevisionFindEventDataFFI struct {
	/* Signature of the revision that was found. */
	Signature LoreHash
}

type LoreRevisionFindEventData struct {
	/* Signature of the revision that was found. */
	Signature LoreHash
}
type LoreRevisionHistoryEventDataFFI struct {
	/* Repository identifier the history belongs to. */
	Repository LoreRepositoryId
	/* Branch identifier the history is listed for. */
	Branch LoreBranchId
}

type LoreRevisionHistoryEventData struct {
	/* Repository identifier the history belongs to. */
	Repository LoreRepositoryId
	/* Branch identifier the history is listed for. */
	Branch LoreBranchId
}
type LoreRevisionHistoryEntryEventDataFFI struct {
	/* Revision hash signature. */
	Revision LoreHash
	/* Revision number. */
	RevisionNumber uint64
	/* Parent revision hashes; the first is the direct parent and the second
	is the other parent of a merge, or zero when there is none. */
	Parent [2]LoreHash
}

type LoreRevisionHistoryEntryEventData struct {
	/* Revision hash signature. */
	Revision LoreHash
	/* Revision number. */
	RevisionNumber uint64
	/* Parent revision hashes; the first is the direct parent and the second
	is the other parent of a merge, or zero when there is none. */
	Parent [2]LoreHash
}
type LoreRevisionRestoreFileBeginEventDataFFI struct {
	/* Number of files to process. */
	Count uintptr
}

type LoreRevisionRestoreFileBeginEventData struct {
	/* Number of files to process. */
	Count uintptr
}
type LoreRevisionRestoreFileEventDataFFI struct {
	/* Path of the file. */
	Path LoreString
	/* Action applied to the file. */
	Action LoreFileAction
	/* Size of the file in bytes. */
	Size uint64
	/* Flag indicating the entry is a file. */
	IsFile uint8
	/* Flag indicating the entry is a directory. */
	IsDirectory uint8
	/* Flag indicating the entry is a module. */
	IsModule uint8
}

type LoreRevisionRestoreFileEventData struct {
	/* Path of the file. */
	Path string
	/* Action applied to the file. */
	Action LoreFileAction
	/* Size of the file in bytes. */
	Size uint64
	/* Flag indicating the entry is a file. */
	IsFile bool
	/* Flag indicating the entry is a directory. */
	IsDirectory bool
	/* Flag indicating the entry is a module. */
	IsModule bool
}
type LoreRevisionRestoreFileEndEventDataFFI struct {
	/* Number of files processed. */
	Count uintptr
}

type LoreRevisionRestoreFileEndEventData struct {
	/* Number of files processed. */
	Count uintptr
}
type LoreRevisionRestoreFragmentBeginEventDataFFI struct {
	/* Number of fragments to transfer. */
	Fragments uint64
}

type LoreRevisionRestoreFragmentBeginEventData struct {
	/* Number of fragments to transfer. */
	Fragments uint64
}
type LoreRevisionRestoreFragmentProgressEventDataFFI struct {
	/* Number of fragments completed. */
	Complete uint64
	/* Total number of fragments. */
	Count uint64
}

type LoreRevisionRestoreFragmentProgressEventData struct {
	/* Number of fragments completed. */
	Complete uint64
	/* Total number of fragments. */
	Count uint64
}
type LoreRevisionRestoreFragmentEndEventDataFFI struct {
	/* Number of fragments transferred. */
	Fragments uint64
}

type LoreRevisionRestoreFragmentEndEventData struct {
	/* Number of fragments transferred. */
	Fragments uint64
}
type LoreRevisionRestoreRevisionEventDataFFI struct {
	/* Resulting revision hash signature. */
	Revision LoreHash
	/* Resulting revision number. */
	RevisionNumber uint64
}

type LoreRevisionRestoreRevisionEventData struct {
	/* Resulting revision hash signature. */
	Revision LoreHash
	/* Resulting revision number. */
	RevisionNumber uint64
}
type LoreRevisionRestoreSyncBeginEventDataFFI struct {
	/* Number of changes to apply. */
	Count uintptr
}

type LoreRevisionRestoreSyncBeginEventData struct {
	/* Number of changes to apply. */
	Count uintptr
}
type LoreRevisionRestoreSyncEndEventDataFFI struct {
	/* Number of changes applied. */
	Count uintptr
}

type LoreRevisionRestoreSyncEndEventData struct {
	/* Number of changes applied. */
	Count uintptr
}
type LoreRevisionResolveEventDataFFI struct {
	/* Repository identifier in which repository */
	Repository LoreRepositoryId
	/* Identifier of the branch on which resolution is being done */
	Branch LoreBranchId
	/* If set to non-empty, the partial hash being resolved */
	Revision LoreString
	/* If set to non-zero, the revision number being resolved */
	RevisionNumber uint64
	/* Resolving using remote data */
	Remote uint8
	/* Resolving using local data */
	Local uint8
}

type LoreRevisionResolveEventData struct {
	/* Repository identifier in which repository */
	Repository LoreRepositoryId
	/* Identifier of the branch on which resolution is being done */
	Branch LoreBranchId
	/* If set to non-empty, the partial hash being resolved */
	Revision string
	/* If set to non-zero, the revision number being resolved */
	RevisionNumber uint64
	/* Resolving using remote data */
	Remote bool
	/* Resolving using local data */
	Local bool
}
type LoreRevisionSyncTargetEventDataFFI struct {
	/* Remote URL */
	Remote LoreString
	/* Repository identifier */
	Repository LoreRepositoryId
	/* Branch identifier (if any) */
	Branch LoreBranchId
	/* Branch name (if any) */
	BranchName LoreString
	/* Current (source) revision identifier */
	SourceRevision LoreHash
	/* Current (source) revision number */
	SourceRevisionNumber uint64
	/* Target revision identifier */
	TargetRevision LoreHash
	/* Target revision number */
	TargetRevisionNumber uint64
	/* Flag indicating revision is the latest revision of the branch */
	IsLatest uint8
	/* Flag indicating revision was from local revision history, not remote */
	Local uint8
}

type LoreRevisionSyncTargetEventData struct {
	/* Remote URL */
	Remote string
	/* Repository identifier */
	Repository LoreRepositoryId
	/* Branch identifier (if any) */
	Branch LoreBranchId
	/* Branch name (if any) */
	BranchName string
	/* Current (source) revision identifier */
	SourceRevision LoreHash
	/* Current (source) revision number */
	SourceRevisionNumber uint64
	/* Target revision identifier */
	TargetRevision LoreHash
	/* Target revision number */
	TargetRevisionNumber uint64
	/* Flag indicating revision is the latest revision of the branch */
	IsLatest bool
	/* Flag indicating revision was from local revision history, not remote */
	Local bool
}
type LoreRevisionSyncFileEventDataFFI struct {
	/* Path of the file relative to the repository root. */
	Path LoreString
	/* Size of the file in bytes. */
	Size uint64
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the entry is a file rather than a directory. */
	FlagFile uint8
}

type LoreRevisionSyncFileEventData struct {
	/* Path of the file relative to the repository root. */
	Path string
	/* Size of the file in bytes. */
	Size uint64
	/* Action applied to the file. */
	Action LoreFileAction
	/* Flag indicating the entry is a file rather than a directory. */
	FlagFile bool
}
type LoreRevisionSyncRevisionEventDataFFI struct {
	/* Branch (if any) */
	Branch LoreBranchId
	/* Resulting revision hash signature */
	Revision LoreHash
	/* Resulting revision number, or 0 if sync resulted in a merge */
	RevisionNumber uint64
	/* Sync resulted in a staged merge revision */
	FlagMerge uint8
	/* Sync resulted in a staged merged revision with conflicts */
	FlagConflict uint8
}

type LoreRevisionSyncRevisionEventData struct {
	/* Branch (if any) */
	Branch LoreBranchId
	/* Resulting revision hash signature */
	Revision LoreHash
	/* Resulting revision number, or 0 if sync resulted in a merge */
	RevisionNumber uint64
	/* Sync resulted in a staged merge revision */
	FlagMerge bool
	/* Sync resulted in a staged merged revision with conflicts */
	FlagConflict bool
}
type LoreRevisionBisectEventDataFFI struct {
	/* Revision number at the start of the search range. */
	StartRevisionNumber uint64
	/* Revision number selected to test next. */
	TargetRevisionNumber uint64
	/* Revision number at the end of the search range. */
	EndRevisionNumber uint64
	/* Flag indicating the search has finished. */
	Done uint8
}

type LoreRevisionBisectEventData struct {
	/* Revision number at the start of the search range. */
	StartRevisionNumber uint64
	/* Revision number selected to test next. */
	TargetRevisionNumber uint64
	/* Revision number at the end of the search range. */
	EndRevisionNumber uint64
	/* Flag indicating the search has finished. */
	Done bool
}
type LoreNotificationBranchCreatedEventDataFFI struct {
	/* Identifier of the created branch. */
	Branch LoreBranchId
}

type LoreNotificationBranchCreatedEventData struct {
	/* Identifier of the created branch. */
	Branch LoreBranchId
}
type LoreNotificationBranchDeletedEventDataFFI struct {
	/* Identifier of the deleted branch. */
	Branch LoreBranchId
}

type LoreNotificationBranchDeletedEventData struct {
	/* Identifier of the deleted branch. */
	Branch LoreBranchId
}
type LoreNotificationBranchPushedEventDataFFI struct {
	/* Hash of the pushed revision. */
	Revision LoreHash
	/* Sequence number of the pushed revision. */
	RevisionNumber uint64
	/* Identifier of the branch that received the revision. */
	Branch LoreBranchId
	/* Identifier of the user that pushed the revision. */
	UserId LoreString
}

type LoreNotificationBranchPushedEventData struct {
	/* Hash of the pushed revision. */
	Revision LoreHash
	/* Sequence number of the pushed revision. */
	RevisionNumber uint64
	/* Identifier of the branch that received the revision. */
	Branch LoreBranchId
	/* Identifier of the user that pushed the revision. */
	UserId string
}
type LoreNotificationResourceLockedEventDataFFI struct {
	/* Identifier of the user that locked the resources. */
	UserId LoreString
	/* Identifier of the branch the resources belong to. */
	Branch LoreBranchId
	/* Paths of the locked resources. */
	Paths LoreStringArrayFFI
}

type LoreNotificationResourceLockedEventData struct {
	/* Identifier of the user that locked the resources. */
	UserId string
	/* Identifier of the branch the resources belong to. */
	Branch LoreBranchId
	/* Paths of the locked resources. */
	Paths []string
}
type LoreNotificationResourceUnlockedEventDataFFI struct {
	/* Identifier of the user that unlocked the resources. */
	UserId LoreString
	/* Identifier of the branch the resources belong to. */
	Branch LoreBranchId
	/* Paths of the unlocked resources. */
	Paths LoreStringArrayFFI
}

type LoreNotificationResourceUnlockedEventData struct {
	/* Identifier of the user that unlocked the resources. */
	UserId string
	/* Identifier of the branch the resources belong to. */
	Branch LoreBranchId
	/* Paths of the unlocked resources. */
	Paths []string
}
type LoreNotificationSubscribedEventDataFFI struct {
	/* Identifier of the subscribed repository. */
	Repository LoreRepositoryId
}

type LoreNotificationSubscribedEventData struct {
	/* Identifier of the subscribed repository. */
	Repository LoreRepositoryId
}
type LoreNotificationUnsubscribedEventDataFFI struct {
	/* Identifier of the unsubscribed repository. */
	Repository LoreRepositoryId
}

type LoreNotificationUnsubscribedEventData struct {
	/* Identifier of the unsubscribed repository. */
	Repository LoreRepositoryId
}
type LoreSharedStoreCreateEventDataFFI struct {
	/* Filesystem path of the created shared store. */
	Path LoreString
}

type LoreSharedStoreCreateEventData struct {
	/* Filesystem path of the created shared store. */
	Path string
}
type LoreSharedStoreInfoEventDataFFI struct {
	/* Nonzero when a shared store is used automatically for the repository. */
	UseAutomatically uint8
	/* Remote URLs of the shared stores. */
	RemoteUrls LoreStringArrayFFI
	/* Filesystem paths of the shared stores. */
	Paths LoreStringArrayFFI
	/* Per-store flag, nonzero when the store exists on disk. */
	Exists LoreUint8ArrayFFI
}

type LoreSharedStoreInfoEventData struct {
	/* Nonzero when a shared store is used automatically for the repository. */
	UseAutomatically bool
	/* Remote URLs of the shared stores. */
	RemoteUrls []string
	/* Filesystem paths of the shared stores. */
	Paths []string
	/* Per-store flag, nonzero when the store exists on disk. */
	Exists []bool
}
type LoreLinkStagedEntryEventDataFFI struct {
	/* Path of the link within the parent repository. */
	Path LoreString
	/* Identifier of the repository the link points to. */
	Repository LoreRepositoryId
	/* Number of staged files inside the link. */
	StagedFileCount uint64
}

type LoreLinkStagedEntryEventData struct {
	/* Path of the link within the parent repository. */
	Path string
	/* Identifier of the repository the link points to. */
	Repository LoreRepositoryId
	/* Number of staged files inside the link. */
	StagedFileCount uint64
}
type LoreStorageOpenedEventDataFFI struct {
	/* Handle id for the opened store. */
	HandleId uint64
}

type LoreStorageOpenedEventData struct {
	/* Handle id for the opened store. */
	HandleId uint64
}
type LoreStoragePutItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The computed content address of the stored item. */
	Address LoreAddress
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStoragePutItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The computed content address of the stored item. */
	Address LoreAddress
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageGetHeaderEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The total reassembled content size in bytes. */
	SizeContent uint64
}

type LoreStorageGetHeaderEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The total reassembled content size in bytes. */
	SizeContent uint64
}
type LoreStorageGetDataEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The byte offset of this payload within the item's content. */
	Offset uint64
	/* The payload bytes for this part of the item. */
	Bytes LoreBytesFFI
}

type LoreStorageGetDataEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The byte offset of this payload within the item's content. */
	Offset uint64
	/* The payload bytes for this part of the item. */
	Bytes LoreBytes
}
type LoreStorageGetItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageGetItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageGetMetadataItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The metadata fragment for the item. */
	Fragment LoreFragmentFFI
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageGetMetadataItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* The metadata fragment for the item. */
	Fragment LoreFragment
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageCopyItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The partition the item was copied from. */
	SourcePartition LorePartition
	/* The partition the item was copied to. */
	TargetPartition LorePartition
	/* The address of the item in the source. */
	SourceAddress LoreAddress
	/* The context of the item in the target. */
	TargetContext LoreContext
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageCopyItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The partition the item was copied from. */
	SourcePartition LorePartition
	/* The partition the item was copied to. */
	TargetPartition LorePartition
	/* The address of the item in the source. */
	SourceAddress LoreAddress
	/* The context of the item in the target. */
	TargetContext LoreContext
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageObliterateItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* 1 when the local side completed without error. */
	LocalSuccess uint8
	/* 1 when the remote side completed without error. */
	RemoteSuccess uint8
	/* 1 when the local side was skipped. */
	LocalSkipped uint8
	/* 1 when the remote side was skipped. */
	RemoteSkipped uint8
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageObliterateItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* 1 when the local side completed without error. */
	LocalSuccess bool
	/* 1 when the remote side completed without error. */
	RemoteSuccess bool
	/* 1 when the local side was skipped. */
	LocalSkipped bool
	/* 1 when the remote side was skipped. */
	RemoteSkipped bool
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageUploadItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* 1 when the item was already durable and no upload was performed. */
	AlreadyDurable uint8
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageUploadItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The content address of the item. */
	Address LoreAddress
	/* 1 when the item was already durable and no upload was performed. */
	AlreadyDurable bool
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeLoadedEventDataFFI struct {
	/* Registry id for the loaded revision tree. */
	HandleId uint64
}

type LoreRevisionTreeLoadedEventData struct {
	/* Registry id for the loaded revision tree. */
	HandleId uint64
}
type LoreRevisionTreeResolvePathCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The resolved node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeResolvePathCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The resolved node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeChildEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The child node. */
	NodeId uint32
	/* The name of the child node. */
	Name LoreString
	/* The parent node. */
	ParentId uint32
	/* The kind of node. */
	Kind uint32
	/* The file mode bits. */
	Mode uint16
	/* The size of the node's content in bytes. */
	Size uint64
	/* The address of the node's content. */
	Address LoreAddress
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeChildEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The child node. */
	NodeId uint32
	/* The name of the child node. */
	Name string
	/* The parent node. */
	ParentId uint32
	/* The kind of node. */
	Kind uint32
	/* The file mode bits. */
	Mode uint16
	/* The size of the node's content in bytes. */
	Size uint64
	/* The address of the node's content. */
	Address LoreAddress
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeNodeInfoEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The queried node. */
	NodeId uint32
	/* The name of the node. */
	Name LoreString
	/* The parent node. */
	ParentId uint32
	/* The kind of node. */
	Kind uint32
	/* The file mode bits. */
	Mode uint16
	/* The size of the node's content in bytes. */
	Size uint64
	/* The address of the node's content. */
	Address LoreAddress
	/* The preserved file id of the node. */
	FileId LoreContext
	/* Root metadata, valid only when the node is the revision root. */
	RootInfo LoreRevisionTreeRootInfoDataFFI
}

type LoreRevisionTreeNodeInfoEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The queried node. */
	NodeId uint32
	/* The name of the node. */
	Name string
	/* The parent node. */
	ParentId uint32
	/* The kind of node. */
	Kind uint32
	/* The file mode bits. */
	Mode uint16
	/* The size of the node's content in bytes. */
	Size uint64
	/* The address of the node's content. */
	Address LoreAddress
	/* The preserved file id of the node. */
	FileId LoreContext
	/* Root metadata, valid only when the node is the revision root. */
	RootInfo LoreRevisionTreeRootInfoData
}
type LoreRevisionTreeNodePathEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The reconstructed path from the root to the queried node. */
	Path LoreString
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeNodePathEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The reconstructed path from the root to the queried node. */
	Path string
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeAddCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The newly-added node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeAddCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The newly-added node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeDeleteCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeDeleteCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeModifyCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The modified node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeModifyCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The modified node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeMoveCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The moved node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeMoveCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The moved node. */
	NodeId uint32
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeMetadataSetCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeMetadataSetCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeMetadataGetCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The metadata key. */
	Key LoreString
	/* The metadata value. */
	Value LoreMetadataFFI
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeMetadataGetCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The metadata key. */
	Key string
	/* The metadata value. */
	Value LoreMetadata
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeCommitCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The newly-committed revision. */
	RevisionHash LoreHash
	/* The observed branch tip when the branch had advanced. */
	NewTipHash LoreHash
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeCommitCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The newly-committed revision. */
	RevisionHash LoreHash
	/* The observed branch tip when the branch had advanced. */
	NewTipHash LoreHash
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreRevisionTreeCloseCompleteEventDataFFI struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}

type LoreRevisionTreeCloseCompleteEventData struct {
	/* Correlation id of the originating call. */
	Id uint64
	/* The outcome of the call. */
	ErrorCode LoreErrorCode
}
type LoreStorageMutableLoadItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The value stored for the key. */
	Value LoreHash
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageMutableLoadItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The value stored for the key. */
	Value LoreHash
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageMutableStoreItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageMutableStoreItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageMutableCompareAndSwapItemCompleteEventDataFFI struct {
	/* Correlation id of the item. */
	Id uint64
	/* The value the key held before the swap. */
	Previous LoreHash
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageMutableCompareAndSwapItemCompleteEventData struct {
	/* Correlation id of the item. */
	Id uint64
	/* The value the key held before the swap. */
	Previous LoreHash
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreStorageMutableListEntryEventDataFFI struct {
	/* Correlation id of the listing item. */
	Id uint64
	/* The key of this entry. */
	Key LoreHash
	/* The value stored for the key. */
	Value LoreHash
}

type LoreStorageMutableListEntryEventData struct {
	/* Correlation id of the listing item. */
	Id uint64
	/* The key of this entry. */
	Key LoreHash
	/* The value stored for the key. */
	Value LoreHash
}
type LoreStorageMutableListItemCompleteEventDataFFI struct {
	/* Correlation id of the listing item. */
	Id uint64
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}

type LoreStorageMutableListItemCompleteEventData struct {
	/* Correlation id of the listing item. */
	Id uint64
	/* The outcome for the item. */
	ErrorCode LoreErrorCode
}
type LoreEvictionBeginEventDataFFI struct {
	/* Fragment capacity the pass is reducing the store toward. */
	TargetFragments uint64
}

type LoreEvictionBeginEventData struct {
	/* Fragment capacity the pass is reducing the store toward. */
	TargetFragments uint64
}
type LoreEvictionProgressEventDataFFI struct {
	/* Fragments evicted from this bucket. */
	Evicted uint64
}

type LoreEvictionProgressEventData struct {
	/* Fragments evicted from this bucket. */
	Evicted uint64
}
type LoreEvictionEndEventDataFFI struct {
	/* Total fragments evicted across the pass. */
	TotalEvicted uint64
}

type LoreEvictionEndEventData struct {
	/* Total fragments evicted across the pass. */
	TotalEvicted uint64
}
type LoreCompactionBeginEventDataFFI struct {
	/* Store size in bytes the pass is reducing the store toward. */
	TargetBytes uint64
}

type LoreCompactionBeginEventData struct {
	/* Store size in bytes the pass is reducing the store toward. */
	TargetBytes uint64
}
type LoreCompactionProgressEventDataFFI struct {
	/* Bytes reclaimed from this group. */
	CompactedBytes uint64
}

type LoreCompactionProgressEventData struct {
	/* Bytes reclaimed from this group. */
	CompactedBytes uint64
}
type LoreCompactionEndEventDataFFI struct {
	/* Total bytes reclaimed across the pass. */
	TotalCompactedBytes uint64
}

type LoreCompactionEndEventData struct {
	/* Total bytes reclaimed across the pass. */
	TotalCompactedBytes uint64
}

func (e *LoreEventFFI) asProgressEventDataFFI() *LoreProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asErrorEventDataFFI() *LoreErrorEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreErrorEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCompleteEventDataFFI() *LoreCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asMetadataEventDataFFI() *LoreMetadataEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreMetadataEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLogEventDataFFI() *LoreLogEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLogEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asEndEventDataFFI() *LoreEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asMaintenanceEventDataFFI() *LoreMaintenanceEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreMaintenanceEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asAuthUrlEventDataFFI() *LoreAuthUrlEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreAuthUrlEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asAuthUserInfoEventDataFFI() *LoreAuthUserInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreAuthUserInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asAuthUserTokenEventDataFFI() *LoreAuthUserTokenEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreAuthUserTokenEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asAuthIdentityEventDataFFI() *LoreAuthIdentityEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreAuthIdentityEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchCreateEventDataFFI() *LoreBranchCreateEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchCreateEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMultipleInstanceEventDataFFI() *LoreBranchMultipleInstanceEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMultipleInstanceEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchArchiveEventDataFFI() *LoreBranchArchiveEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchArchiveEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchListBeginEventDataFFI() *LoreBranchListBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchListBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchListEntryEventDataFFI() *LoreBranchListEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchListEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchListEndEventDataFFI() *LoreBranchListEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchListEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeAbortBeginEventDataFFI() *LoreBranchMergeAbortBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeAbortBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeAbortEndEventDataFFI() *LoreBranchMergeAbortEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeAbortEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchInfoEventDataFFI() *LoreBranchInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffBeginEventDataFFI() *LoreBranchDiffBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffChangeBeginEventDataFFI() *LoreBranchDiffChangeBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffChangeBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffChangeEventDataFFI() *LoreBranchDiffChangeEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffChangeEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffChangeEndEventDataFFI() *LoreBranchDiffChangeEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffChangeEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffConflictBeginEventDataFFI() *LoreBranchDiffConflictBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffConflictBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffConflictEventDataFFI() *LoreBranchDiffConflictEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffConflictEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffConflictEndEventDataFFI() *LoreBranchDiffConflictEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffConflictEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchDiffEndEventDataFFI() *LoreBranchDiffEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchDiffEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchLatestListEntryEventDataFFI() *LoreBranchLatestListEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchLatestListEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeConflictFileEventDataFFI() *LoreBranchMergeConflictFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeConflictFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeLinkSkippedEventDataFFI() *LoreBranchMergeLinkSkippedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeLinkSkippedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeUnresolveFileEventDataFFI() *LoreBranchMergeUnresolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeUnresolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeUnresolveRevisionEventDataFFI() *LoreBranchMergeUnresolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeUnresolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFileBeginEventDataFFI() *LoreBranchMergeIntoFileBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFileBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFileEventDataFFI() *LoreBranchMergeIntoFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFileEndEventDataFFI() *LoreBranchMergeIntoFileEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFileEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFragmentBeginEventDataFFI() *LoreBranchMergeIntoFragmentBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFragmentBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFragmentProgressEventDataFFI() *LoreBranchMergeIntoFragmentProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFragmentProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoFragmentEndEventDataFFI() *LoreBranchMergeIntoFragmentEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoFragmentEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoRevisionEventDataFFI() *LoreBranchMergeIntoRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoSyncBeginEventDataFFI() *LoreBranchMergeIntoSyncBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoSyncBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeIntoSyncEndEventDataFFI() *LoreBranchMergeIntoSyncEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeIntoSyncEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeResolveFileEventDataFFI() *LoreBranchMergeResolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeResolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeResolveRevisionEventDataFFI() *LoreBranchMergeResolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeResolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeStartBeginEventDataFFI() *LoreBranchMergeStartBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeStartBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchMergeStartEndEventDataFFI() *LoreBranchMergeStartEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchMergeStartEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickStartBeginEventDataFFI() *LoreCherryPickStartBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickStartBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickStartEndEventDataFFI() *LoreCherryPickStartEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickStartEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickAbortBeginEventDataFFI() *LoreCherryPickAbortBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickAbortBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickAbortEndEventDataFFI() *LoreCherryPickAbortEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickAbortEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickConflictFileEventDataFFI() *LoreCherryPickConflictFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickConflictFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickUnresolveFileEventDataFFI() *LoreCherryPickUnresolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickUnresolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickUnresolveRevisionEventDataFFI() *LoreCherryPickUnresolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickUnresolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickResolveFileEventDataFFI() *LoreCherryPickResolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickResolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCherryPickResolveRevisionEventDataFFI() *LoreCherryPickResolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCherryPickResolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertStartBeginEventDataFFI() *LoreRevertStartBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertStartBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertStartEndEventDataFFI() *LoreRevertStartEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertStartEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertAbortBeginEventDataFFI() *LoreRevertAbortBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertAbortBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertAbortEndEventDataFFI() *LoreRevertAbortEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertAbortEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertResolveFileEventDataFFI() *LoreRevertResolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertResolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertResolveRevisionEventDataFFI() *LoreRevertResolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertResolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertConflictFileEventDataFFI() *LoreRevertConflictFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertConflictFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertUnresolveFileEventDataFFI() *LoreRevertUnresolveFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertUnresolveFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevertUnresolveRevisionEventDataFFI() *LoreRevertUnresolveRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevertUnresolveRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchProtectEventDataFFI() *LoreBranchProtectEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchProtectEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushEventDataFFI() *LoreBranchPushEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushRevisionUpdateBeginEventDataFFI() *LoreBranchPushRevisionUpdateBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushRevisionUpdateBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushRevisionUpdateEndEventDataFFI() *LoreBranchPushRevisionUpdateEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushRevisionUpdateEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushFragmentBeginEventDataFFI() *LoreBranchPushFragmentBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushFragmentBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushFragmentProgressEventDataFFI() *LoreBranchPushFragmentProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushFragmentProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushFragmentEndEventDataFFI() *LoreBranchPushFragmentEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushFragmentEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushBranchCreateBeginEventDataFFI() *LoreBranchPushBranchCreateBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushBranchCreateBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushBranchCreateEndEventDataFFI() *LoreBranchPushBranchCreateEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushBranchCreateEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushRevisionPushBeginEventDataFFI() *LoreBranchPushRevisionPushBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushRevisionPushBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushRevisionPushUpdateEventDataFFI() *LoreBranchPushRevisionPushUpdateEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushRevisionPushUpdateEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchPushRevisionPushEndEventDataFFI() *LoreBranchPushRevisionPushEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchPushRevisionPushEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchResetEventDataFFI() *LoreBranchResetEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchResetEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchSwitchBeginEventDataFFI() *LoreBranchSwitchBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchSwitchBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchSwitchEndEventDataFFI() *LoreBranchSwitchEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchSwitchEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asBranchUnprotectEventDataFFI() *LoreBranchUnprotectEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreBranchUnprotectEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileInfoEventDataFFI() *LoreFileInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDiffEventDataFFI() *LoreFileDiffEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDiffEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileHashEventDataFFI() *LoreFileHashEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileHashEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileHistoryEventDataFFI() *LoreFileHistoryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileHistoryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileWriteEventDataFFI() *LoreFileWriteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileWriteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileObliterateEventDataFFI() *LoreFileObliterateEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileObliterateEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDumpEventDataFFI() *LoreFileDumpEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDumpEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyAddBeginEventDataFFI() *LoreFileDependencyAddBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyAddBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyAddEntryEventDataFFI() *LoreFileDependencyAddEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyAddEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyAddEndEventDataFFI() *LoreFileDependencyAddEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyAddEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyRemoveBeginEventDataFFI() *LoreFileDependencyRemoveBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyRemoveBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyRemoveEntryEventDataFFI() *LoreFileDependencyRemoveEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyRemoveEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyRemoveEndEventDataFFI() *LoreFileDependencyRemoveEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyRemoveEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyListBeginEventDataFFI() *LoreFileDependencyListBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyListBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyListFileEventDataFFI() *LoreFileDependencyListFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyListFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyListEntryEventDataFFI() *LoreFileDependencyListEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyListEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyListFileEndEventDataFFI() *LoreFileDependencyListFileEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyListFileEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileDependencyListEndEventDataFFI() *LoreFileDependencyListEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileDependencyListEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileResetBeginEventDataFFI() *LoreFileResetBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileResetBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileResetProgressEventDataFFI() *LoreFileResetProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileResetProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileResetEndEventDataFFI() *LoreFileResetEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileResetEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileResetFileEventDataFFI() *LoreFileResetFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileResetFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFilterExcludeEventDataFFI() *LoreFilterExcludeEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFilterExcludeEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileStageBeginEventDataFFI() *LoreFileStageBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileStageBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileStageProgressEventDataFFI() *LoreFileStageProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileStageProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileStageEndEventDataFFI() *LoreFileStageEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileStageEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileStageRevisionEventDataFFI() *LoreFileStageRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileStageRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileStageFileEventDataFFI() *LoreFileStageFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileStageFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileUnstageBeginEventDataFFI() *LoreFileUnstageBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileUnstageBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileUnstageProgressEventDataFFI() *LoreFileUnstageProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileUnstageProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileUnstageEndEventDataFFI() *LoreFileUnstageEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileUnstageEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileUnstageRevisionEventDataFFI() *LoreFileUnstageRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileUnstageRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFileUnstageFileEventDataFFI() *LoreFileUnstageFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFileUnstageFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asFragmentWriteEventDataFFI() *LoreFragmentWriteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreFragmentWriteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLayerAddEventDataFFI() *LoreLayerAddEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLayerAddEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLayerEntryEventDataFFI() *LoreLayerEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLayerEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLayerRemoveEventDataFFI() *LoreLayerRemoveEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLayerRemoveEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLayerStagedEntryEventDataFFI() *LoreLayerStagedEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLayerStagedEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLinkChangeEventDataFFI() *LoreLinkChangeEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLinkChangeEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLinkEntryEventDataFFI() *LoreLinkEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLinkEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileAcquireBeginEventDataFFI() *LoreLockFileAcquireBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileAcquireBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileAcquireEventDataFFI() *LoreLockFileAcquireEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileAcquireEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileStatusBeginEventDataFFI() *LoreLockFileStatusBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileStatusBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileStatusEventDataFFI() *LoreLockFileStatusEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileStatusEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileQueryBeginEventDataFFI() *LoreLockFileQueryBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileQueryBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileQueryEventDataFFI() *LoreLockFileQueryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileQueryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileReleaseBeginEventDataFFI() *LoreLockFileReleaseBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileReleaseBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLockFileReleaseEventDataFFI() *LoreLockFileReleaseEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLockFileReleaseEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asMetadataClearFileEventDataFFI() *LoreMetadataClearFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreMetadataClearFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asMetadataClearRevisionEventDataFFI() *LoreMetadataClearRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreMetadataClearRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asPathIgnoreEventDataFFI() *LorePathIgnoreEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LorePathIgnoreEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryCreateEventDataFFI() *LoreRepositoryCreateEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryCreateEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryCloneBeginEventDataFFI() *LoreRepositoryCloneBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryCloneBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryCloneProgressEventDataFFI() *LoreRepositoryCloneProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryCloneProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryCloneEndEventDataFFI() *LoreRepositoryCloneEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryCloneEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asDependencyResolveBeginEventDataFFI() *LoreDependencyResolveBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreDependencyResolveBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asDependencyResolveItemEventDataFFI() *LoreDependencyResolveItemEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreDependencyResolveItemEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asDependencyResolveEndEventDataFFI() *LoreDependencyResolveEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreDependencyResolveEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryDataEventDataFFI() *LoreRepositoryDataEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryDataEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryConfigGetEventDataFFI() *LoreRepositoryConfigGetEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryConfigGetEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryDumpBeginEventDataFFI() *LoreRepositoryDumpBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryDumpBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryDumpEndEventDataFFI() *LoreRepositoryDumpEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryDumpEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryListEntryEventDataFFI() *LoreRepositoryListEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryListEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryInstanceEventDataFFI() *LoreRepositoryInstanceEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryInstanceEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryVerifyStateBeginEventDataFFI() *LoreRepositoryVerifyStateBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryVerifyStateBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryVerifyStateEndEventDataFFI() *LoreRepositoryVerifyStateEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryVerifyStateEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryVerifyFragmentEventDataFFI() *LoreRepositoryVerifyFragmentEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryVerifyFragmentEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryVerifyFragmentMatchEventDataFFI() *LoreRepositoryVerifyFragmentMatchEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryVerifyFragmentMatchEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryVerifyFragmentRemoteEventDataFFI() *LoreRepositoryVerifyFragmentRemoteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryVerifyFragmentRemoteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStateDumpEventDataFFI() *LoreRepositoryStateDumpEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStateDumpEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStateDumpNodeEventDataFFI() *LoreRepositoryStateDumpNodeEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStateDumpNodeEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStatusRevisionEventDataFFI() *LoreRepositoryStatusRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStatusRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStatusFileEventDataFFI() *LoreRepositoryStatusFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStatusFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStatusCountEventDataFFI() *LoreRepositoryStatusCountEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStatusCountEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStatusSummaryEventDataFFI() *LoreRepositoryStatusSummaryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStatusSummaryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRepositoryStoreImmutableQueryEventDataFFI() *LoreRepositoryStoreImmutableQueryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRepositoryStoreImmutableQueryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionCommitBeginEventDataFFI() *LoreRevisionCommitBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionCommitBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionCommitProgressEventDataFFI() *LoreRevisionCommitProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionCommitProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionCommitEndEventDataFFI() *LoreRevisionCommitEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionCommitEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionCommitRevisionEventDataFFI() *LoreRevisionCommitRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionCommitRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionInfoEventDataFFI() *LoreRevisionInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionInfoDeltaEventDataFFI() *LoreRevisionInfoDeltaEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionInfoDeltaEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionDiffFileEventDataFFI() *LoreRevisionDiffFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionDiffFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionFindEventDataFFI() *LoreRevisionFindEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionFindEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionHistoryEventDataFFI() *LoreRevisionHistoryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionHistoryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionHistoryEntryEventDataFFI() *LoreRevisionHistoryEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionHistoryEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFileBeginEventDataFFI() *LoreRevisionRestoreFileBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFileBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFileEventDataFFI() *LoreRevisionRestoreFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFileEndEventDataFFI() *LoreRevisionRestoreFileEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFileEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFragmentBeginEventDataFFI() *LoreRevisionRestoreFragmentBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFragmentBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFragmentProgressEventDataFFI() *LoreRevisionRestoreFragmentProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFragmentProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreFragmentEndEventDataFFI() *LoreRevisionRestoreFragmentEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreFragmentEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreRevisionEventDataFFI() *LoreRevisionRestoreRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreSyncBeginEventDataFFI() *LoreRevisionRestoreSyncBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreSyncBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionRestoreSyncEndEventDataFFI() *LoreRevisionRestoreSyncEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionRestoreSyncEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionResolveEventDataFFI() *LoreRevisionResolveEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionResolveEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionSyncTargetEventDataFFI() *LoreRevisionSyncTargetEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionSyncTargetEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionSyncFileEventDataFFI() *LoreRevisionSyncFileEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionSyncFileEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionSyncProgressEventDataFFI() *LoreRevisionSyncProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionSyncProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionSyncRevisionEventDataFFI() *LoreRevisionSyncRevisionEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionSyncRevisionEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionBisectEventDataFFI() *LoreRevisionBisectEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionBisectEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationBranchCreatedEventDataFFI() *LoreNotificationBranchCreatedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationBranchCreatedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationBranchDeletedEventDataFFI() *LoreNotificationBranchDeletedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationBranchDeletedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationBranchPushedEventDataFFI() *LoreNotificationBranchPushedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationBranchPushedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationResourceLockedEventDataFFI() *LoreNotificationResourceLockedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationResourceLockedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationResourceUnlockedEventDataFFI() *LoreNotificationResourceUnlockedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationResourceUnlockedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationSubscribedEventDataFFI() *LoreNotificationSubscribedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationSubscribedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asNotificationUnsubscribedEventDataFFI() *LoreNotificationUnsubscribedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreNotificationUnsubscribedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asSharedStoreCreateEventDataFFI() *LoreSharedStoreCreateEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreSharedStoreCreateEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asSharedStoreInfoEventDataFFI() *LoreSharedStoreInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreSharedStoreInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asLinkStagedEntryEventDataFFI() *LoreLinkStagedEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreLinkStagedEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageOpenedEventDataFFI() *LoreStorageOpenedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageOpenedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStoragePutItemCompleteEventDataFFI() *LoreStoragePutItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStoragePutItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageGetHeaderEventDataFFI() *LoreStorageGetHeaderEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageGetHeaderEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageGetDataEventDataFFI() *LoreStorageGetDataEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageGetDataEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageGetItemCompleteEventDataFFI() *LoreStorageGetItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageGetItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageGetMetadataItemCompleteEventDataFFI() *LoreStorageGetMetadataItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageGetMetadataItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageCopyItemCompleteEventDataFFI() *LoreStorageCopyItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageCopyItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageObliterateItemCompleteEventDataFFI() *LoreStorageObliterateItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageObliterateItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageUploadItemCompleteEventDataFFI() *LoreStorageUploadItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageUploadItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeLoadedEventDataFFI() *LoreRevisionTreeLoadedEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeLoadedEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeResolvePathCompleteEventDataFFI() *LoreRevisionTreeResolvePathCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeResolvePathCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeChildEventDataFFI() *LoreRevisionTreeChildEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeChildEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeNodeInfoEventDataFFI() *LoreRevisionTreeNodeInfoEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeNodeInfoEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeNodePathEventDataFFI() *LoreRevisionTreeNodePathEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeNodePathEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeAddCompleteEventDataFFI() *LoreRevisionTreeAddCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeAddCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeDeleteCompleteEventDataFFI() *LoreRevisionTreeDeleteCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeDeleteCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeModifyCompleteEventDataFFI() *LoreRevisionTreeModifyCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeModifyCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeMoveCompleteEventDataFFI() *LoreRevisionTreeMoveCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeMoveCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeMetadataSetCompleteEventDataFFI() *LoreRevisionTreeMetadataSetCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeMetadataSetCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeMetadataGetCompleteEventDataFFI() *LoreRevisionTreeMetadataGetCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeMetadataGetCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeCommitCompleteEventDataFFI() *LoreRevisionTreeCommitCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeCommitCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asRevisionTreeCloseCompleteEventDataFFI() *LoreRevisionTreeCloseCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreRevisionTreeCloseCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageMutableLoadItemCompleteEventDataFFI() *LoreStorageMutableLoadItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageMutableLoadItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageMutableStoreItemCompleteEventDataFFI() *LoreStorageMutableStoreItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageMutableStoreItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageMutableCompareAndSwapItemCompleteEventDataFFI() *LoreStorageMutableCompareAndSwapItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageMutableCompareAndSwapItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageMutableListEntryEventDataFFI() *LoreStorageMutableListEntryEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageMutableListEntryEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asStorageMutableListItemCompleteEventDataFFI() *LoreStorageMutableListItemCompleteEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreStorageMutableListItemCompleteEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asEvictionBeginEventDataFFI() *LoreEvictionBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreEvictionBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asEvictionProgressEventDataFFI() *LoreEvictionProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreEvictionProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asEvictionEndEventDataFFI() *LoreEvictionEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreEvictionEndEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCompactionBeginEventDataFFI() *LoreCompactionBeginEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCompactionBeginEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCompactionProgressEventDataFFI() *LoreCompactionProgressEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCompactionProgressEventDataFFI)(unionPtr)
}
func (e *LoreEventFFI) asCompactionEndEventDataFFI() *LoreCompactionEndEventDataFFI {
	unionPtr := unsafe.Add(unsafe.Pointer(e), loreEventUnionOffset)
	return (*LoreCompactionEndEventDataFFI)(unionPtr)
}

func (e *LoreEventFFI) GetData() any {
	switch e.Tag {
	case LoreEventTag_PROGRESS:
		return e.asProgressEventDataFFI()
	case LoreEventTag_ERROR:
		return e.asErrorEventDataFFI()
	case LoreEventTag_COMPLETE:
		return e.asCompleteEventDataFFI()
	case LoreEventTag_METADATA:
		return e.asMetadataEventDataFFI()
	case LoreEventTag_LOG:
		return e.asLogEventDataFFI()
	case LoreEventTag_END:
		return e.asEndEventDataFFI()
	case LoreEventTag_MAINTENANCE:
		return e.asMaintenanceEventDataFFI()
	case LoreEventTag_AUTH_URL:
		return e.asAuthUrlEventDataFFI()
	case LoreEventTag_AUTH_USER_INFO:
		return e.asAuthUserInfoEventDataFFI()
	case LoreEventTag_AUTH_USER_TOKEN:
		return e.asAuthUserTokenEventDataFFI()
	case LoreEventTag_AUTH_IDENTITY:
		return e.asAuthIdentityEventDataFFI()
	case LoreEventTag_BRANCH_CREATE:
		return e.asBranchCreateEventDataFFI()
	case LoreEventTag_BRANCH_MULTIPLE_INSTANCE:
		return e.asBranchMultipleInstanceEventDataFFI()
	case LoreEventTag_BRANCH_ARCHIVE:
		return e.asBranchArchiveEventDataFFI()
	case LoreEventTag_BRANCH_LIST_BEGIN:
		return e.asBranchListBeginEventDataFFI()
	case LoreEventTag_BRANCH_LIST_ENTRY:
		return e.asBranchListEntryEventDataFFI()
	case LoreEventTag_BRANCH_LIST_END:
		return e.asBranchListEndEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_ABORT_BEGIN:
		return e.asBranchMergeAbortBeginEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_ABORT_END:
		return e.asBranchMergeAbortEndEventDataFFI()
	case LoreEventTag_BRANCH_INFO:
		return e.asBranchInfoEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_BEGIN:
		return e.asBranchDiffBeginEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CHANGE_BEGIN:
		return e.asBranchDiffChangeBeginEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CHANGE:
		return e.asBranchDiffChangeEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CHANGE_END:
		return e.asBranchDiffChangeEndEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CONFLICT_BEGIN:
		return e.asBranchDiffConflictBeginEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CONFLICT:
		return e.asBranchDiffConflictEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_CONFLICT_END:
		return e.asBranchDiffConflictEndEventDataFFI()
	case LoreEventTag_BRANCH_DIFF_END:
		return e.asBranchDiffEndEventDataFFI()
	case LoreEventTag_BRANCH_LATEST_LIST_ENTRY:
		return e.asBranchLatestListEntryEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_CONFLICT_FILE:
		return e.asBranchMergeConflictFileEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_LINK_SKIPPED:
		return e.asBranchMergeLinkSkippedEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_UNRESOLVE_FILE:
		return e.asBranchMergeUnresolveFileEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_UNRESOLVE_REVISION:
		return e.asBranchMergeUnresolveRevisionEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FILE_BEGIN:
		return e.asBranchMergeIntoFileBeginEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FILE:
		return e.asBranchMergeIntoFileEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FILE_END:
		return e.asBranchMergeIntoFileEndEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_BEGIN:
		return e.asBranchMergeIntoFragmentBeginEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_PROGRESS:
		return e.asBranchMergeIntoFragmentProgressEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_END:
		return e.asBranchMergeIntoFragmentEndEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_REVISION:
		return e.asBranchMergeIntoRevisionEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_SYNC_BEGIN:
		return e.asBranchMergeIntoSyncBeginEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_INTO_SYNC_END:
		return e.asBranchMergeIntoSyncEndEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_RESOLVE_FILE:
		return e.asBranchMergeResolveFileEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_RESOLVE_REVISION:
		return e.asBranchMergeResolveRevisionEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_START_BEGIN:
		return e.asBranchMergeStartBeginEventDataFFI()
	case LoreEventTag_BRANCH_MERGE_START_END:
		return e.asBranchMergeStartEndEventDataFFI()
	case LoreEventTag_CHERRY_PICK_START_BEGIN:
		return e.asCherryPickStartBeginEventDataFFI()
	case LoreEventTag_CHERRY_PICK_START_END:
		return e.asCherryPickStartEndEventDataFFI()
	case LoreEventTag_CHERRY_PICK_ABORT_BEGIN:
		return e.asCherryPickAbortBeginEventDataFFI()
	case LoreEventTag_CHERRY_PICK_ABORT_END:
		return e.asCherryPickAbortEndEventDataFFI()
	case LoreEventTag_CHERRY_PICK_CONFLICT_FILE:
		return e.asCherryPickConflictFileEventDataFFI()
	case LoreEventTag_CHERRY_PICK_UNRESOLVE_FILE:
		return e.asCherryPickUnresolveFileEventDataFFI()
	case LoreEventTag_CHERRY_PICK_UNRESOLVE_REVISION:
		return e.asCherryPickUnresolveRevisionEventDataFFI()
	case LoreEventTag_CHERRY_PICK_RESOLVE_FILE:
		return e.asCherryPickResolveFileEventDataFFI()
	case LoreEventTag_CHERRY_PICK_RESOLVE_REVISION:
		return e.asCherryPickResolveRevisionEventDataFFI()
	case LoreEventTag_REVERT_START_BEGIN:
		return e.asRevertStartBeginEventDataFFI()
	case LoreEventTag_REVERT_START_END:
		return e.asRevertStartEndEventDataFFI()
	case LoreEventTag_REVERT_ABORT_BEGIN:
		return e.asRevertAbortBeginEventDataFFI()
	case LoreEventTag_REVERT_ABORT_END:
		return e.asRevertAbortEndEventDataFFI()
	case LoreEventTag_REVERT_RESOLVE_FILE:
		return e.asRevertResolveFileEventDataFFI()
	case LoreEventTag_REVERT_RESOLVE_REVISION:
		return e.asRevertResolveRevisionEventDataFFI()
	case LoreEventTag_REVERT_CONFLICT_FILE:
		return e.asRevertConflictFileEventDataFFI()
	case LoreEventTag_REVERT_UNRESOLVE_FILE:
		return e.asRevertUnresolveFileEventDataFFI()
	case LoreEventTag_REVERT_UNRESOLVE_REVISION:
		return e.asRevertUnresolveRevisionEventDataFFI()
	case LoreEventTag_BRANCH_PROTECT:
		return e.asBranchProtectEventDataFFI()
	case LoreEventTag_BRANCH_PUSH:
		return e.asBranchPushEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_BEGIN:
		return e.asBranchPushRevisionUpdateBeginEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_END:
		return e.asBranchPushRevisionUpdateEndEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_BEGIN:
		return e.asBranchPushFragmentBeginEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_PROGRESS:
		return e.asBranchPushFragmentProgressEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_END:
		return e.asBranchPushFragmentEndEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_BEGIN:
		return e.asBranchPushBranchCreateBeginEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_END:
		return e.asBranchPushBranchCreateEndEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_BEGIN:
		return e.asBranchPushRevisionPushBeginEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_UPDATE:
		return e.asBranchPushRevisionPushUpdateEventDataFFI()
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_END:
		return e.asBranchPushRevisionPushEndEventDataFFI()
	case LoreEventTag_BRANCH_RESET:
		return e.asBranchResetEventDataFFI()
	case LoreEventTag_BRANCH_SWITCH_BEGIN:
		return e.asBranchSwitchBeginEventDataFFI()
	case LoreEventTag_BRANCH_SWITCH_END:
		return e.asBranchSwitchEndEventDataFFI()
	case LoreEventTag_BRANCH_UNPROTECT:
		return e.asBranchUnprotectEventDataFFI()
	case LoreEventTag_FILE_INFO:
		return e.asFileInfoEventDataFFI()
	case LoreEventTag_FILE_DIFF:
		return e.asFileDiffEventDataFFI()
	case LoreEventTag_FILE_HASH:
		return e.asFileHashEventDataFFI()
	case LoreEventTag_FILE_HISTORY:
		return e.asFileHistoryEventDataFFI()
	case LoreEventTag_FILE_WRITE:
		return e.asFileWriteEventDataFFI()
	case LoreEventTag_FILE_OBLITERATE:
		return e.asFileObliterateEventDataFFI()
	case LoreEventTag_FILE_DUMP:
		return e.asFileDumpEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_ADD_BEGIN:
		return e.asFileDependencyAddBeginEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_ADD_ENTRY:
		return e.asFileDependencyAddEntryEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_ADD_END:
		return e.asFileDependencyAddEndEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_BEGIN:
		return e.asFileDependencyRemoveBeginEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_ENTRY:
		return e.asFileDependencyRemoveEntryEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_END:
		return e.asFileDependencyRemoveEndEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_LIST_BEGIN:
		return e.asFileDependencyListBeginEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_LIST_FILE:
		return e.asFileDependencyListFileEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_LIST_ENTRY:
		return e.asFileDependencyListEntryEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_LIST_FILE_END:
		return e.asFileDependencyListFileEndEventDataFFI()
	case LoreEventTag_FILE_DEPENDENCY_LIST_END:
		return e.asFileDependencyListEndEventDataFFI()
	case LoreEventTag_FILE_RESET_BEGIN:
		return e.asFileResetBeginEventDataFFI()
	case LoreEventTag_FILE_RESET_PROGRESS:
		return e.asFileResetProgressEventDataFFI()
	case LoreEventTag_FILE_RESET_END:
		return e.asFileResetEndEventDataFFI()
	case LoreEventTag_FILE_RESET_FILE:
		return e.asFileResetFileEventDataFFI()
	case LoreEventTag_FILTER_EXCLUDE:
		return e.asFilterExcludeEventDataFFI()
	case LoreEventTag_FILE_STAGE_BEGIN:
		return e.asFileStageBeginEventDataFFI()
	case LoreEventTag_FILE_STAGE_PROGRESS:
		return e.asFileStageProgressEventDataFFI()
	case LoreEventTag_FILE_STAGE_END:
		return e.asFileStageEndEventDataFFI()
	case LoreEventTag_FILE_STAGE_REVISION:
		return e.asFileStageRevisionEventDataFFI()
	case LoreEventTag_FILE_STAGE_FILE:
		return e.asFileStageFileEventDataFFI()
	case LoreEventTag_FILE_UNSTAGE_BEGIN:
		return e.asFileUnstageBeginEventDataFFI()
	case LoreEventTag_FILE_UNSTAGE_PROGRESS:
		return e.asFileUnstageProgressEventDataFFI()
	case LoreEventTag_FILE_UNSTAGE_END:
		return e.asFileUnstageEndEventDataFFI()
	case LoreEventTag_FILE_UNSTAGE_REVISION:
		return e.asFileUnstageRevisionEventDataFFI()
	case LoreEventTag_FILE_UNSTAGE_FILE:
		return e.asFileUnstageFileEventDataFFI()
	case LoreEventTag_FRAGMENT_WRITE:
		return e.asFragmentWriteEventDataFFI()
	case LoreEventTag_LAYER_ADD:
		return e.asLayerAddEventDataFFI()
	case LoreEventTag_LAYER_ENTRY:
		return e.asLayerEntryEventDataFFI()
	case LoreEventTag_LAYER_REMOVE:
		return e.asLayerRemoveEventDataFFI()
	case LoreEventTag_LAYER_STAGED_ENTRY:
		return e.asLayerStagedEntryEventDataFFI()
	case LoreEventTag_LINK_CHANGE:
		return e.asLinkChangeEventDataFFI()
	case LoreEventTag_LINK_ENTRY:
		return e.asLinkEntryEventDataFFI()
	case LoreEventTag_LOCK_FILE_ACQUIRE_BEGIN:
		return e.asLockFileAcquireBeginEventDataFFI()
	case LoreEventTag_LOCK_FILE_ACQUIRE:
		return e.asLockFileAcquireEventDataFFI()
	case LoreEventTag_LOCK_FILE_STATUS_BEGIN:
		return e.asLockFileStatusBeginEventDataFFI()
	case LoreEventTag_LOCK_FILE_STATUS:
		return e.asLockFileStatusEventDataFFI()
	case LoreEventTag_LOCK_FILE_QUERY_BEGIN:
		return e.asLockFileQueryBeginEventDataFFI()
	case LoreEventTag_LOCK_FILE_QUERY:
		return e.asLockFileQueryEventDataFFI()
	case LoreEventTag_LOCK_FILE_RELEASE_BEGIN:
		return e.asLockFileReleaseBeginEventDataFFI()
	case LoreEventTag_LOCK_FILE_RELEASE:
		return e.asLockFileReleaseEventDataFFI()
	case LoreEventTag_METADATA_CLEAR_FILE:
		return e.asMetadataClearFileEventDataFFI()
	case LoreEventTag_METADATA_CLEAR_REVISION:
		return e.asMetadataClearRevisionEventDataFFI()
	case LoreEventTag_PATH_IGNORE:
		return e.asPathIgnoreEventDataFFI()
	case LoreEventTag_REPOSITORY_CREATE:
		return e.asRepositoryCreateEventDataFFI()
	case LoreEventTag_REPOSITORY_CLONE_BEGIN:
		return e.asRepositoryCloneBeginEventDataFFI()
	case LoreEventTag_REPOSITORY_CLONE_PROGRESS:
		return e.asRepositoryCloneProgressEventDataFFI()
	case LoreEventTag_REPOSITORY_CLONE_END:
		return e.asRepositoryCloneEndEventDataFFI()
	case LoreEventTag_DEPENDENCY_RESOLVE_BEGIN:
		return e.asDependencyResolveBeginEventDataFFI()
	case LoreEventTag_DEPENDENCY_RESOLVE_ITEM:
		return e.asDependencyResolveItemEventDataFFI()
	case LoreEventTag_DEPENDENCY_RESOLVE_END:
		return e.asDependencyResolveEndEventDataFFI()
	case LoreEventTag_REPOSITORY_DATA:
		return e.asRepositoryDataEventDataFFI()
	case LoreEventTag_REPOSITORY_CONFIG_GET:
		return e.asRepositoryConfigGetEventDataFFI()
	case LoreEventTag_REPOSITORY_DUMP_BEGIN:
		return e.asRepositoryDumpBeginEventDataFFI()
	case LoreEventTag_REPOSITORY_DUMP_END:
		return e.asRepositoryDumpEndEventDataFFI()
	case LoreEventTag_REPOSITORY_LIST_ENTRY:
		return e.asRepositoryListEntryEventDataFFI()
	case LoreEventTag_REPOSITORY_INSTANCE:
		return e.asRepositoryInstanceEventDataFFI()
	case LoreEventTag_REPOSITORY_VERIFY_STATE_BEGIN:
		return e.asRepositoryVerifyStateBeginEventDataFFI()
	case LoreEventTag_REPOSITORY_VERIFY_STATE_END:
		return e.asRepositoryVerifyStateEndEventDataFFI()
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT:
		return e.asRepositoryVerifyFragmentEventDataFFI()
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_MATCH:
		return e.asRepositoryVerifyFragmentMatchEventDataFFI()
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_REMOTE:
		return e.asRepositoryVerifyFragmentRemoteEventDataFFI()
	case LoreEventTag_REPOSITORY_STATE_DUMP:
		return e.asRepositoryStateDumpEventDataFFI()
	case LoreEventTag_REPOSITORY_STATE_DUMP_NODE:
		return e.asRepositoryStateDumpNodeEventDataFFI()
	case LoreEventTag_REPOSITORY_STATUS_REVISION:
		return e.asRepositoryStatusRevisionEventDataFFI()
	case LoreEventTag_REPOSITORY_STATUS_FILE:
		return e.asRepositoryStatusFileEventDataFFI()
	case LoreEventTag_REPOSITORY_STATUS_COUNT:
		return e.asRepositoryStatusCountEventDataFFI()
	case LoreEventTag_REPOSITORY_STATUS_SUMMARY:
		return e.asRepositoryStatusSummaryEventDataFFI()
	case LoreEventTag_REPOSITORY_STORE_IMMUTABLE_QUERY:
		return e.asRepositoryStoreImmutableQueryEventDataFFI()
	case LoreEventTag_REVISION_COMMIT_BEGIN:
		return e.asRevisionCommitBeginEventDataFFI()
	case LoreEventTag_REVISION_COMMIT_PROGRESS:
		return e.asRevisionCommitProgressEventDataFFI()
	case LoreEventTag_REVISION_COMMIT_END:
		return e.asRevisionCommitEndEventDataFFI()
	case LoreEventTag_REVISION_COMMIT_REVISION:
		return e.asRevisionCommitRevisionEventDataFFI()
	case LoreEventTag_REVISION_INFO:
		return e.asRevisionInfoEventDataFFI()
	case LoreEventTag_REVISION_INFO_DELTA:
		return e.asRevisionInfoDeltaEventDataFFI()
	case LoreEventTag_REVISION_DIFF_FILE:
		return e.asRevisionDiffFileEventDataFFI()
	case LoreEventTag_REVISION_FIND:
		return e.asRevisionFindEventDataFFI()
	case LoreEventTag_REVISION_HISTORY:
		return e.asRevisionHistoryEventDataFFI()
	case LoreEventTag_REVISION_HISTORY_ENTRY:
		return e.asRevisionHistoryEntryEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FILE_BEGIN:
		return e.asRevisionRestoreFileBeginEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FILE:
		return e.asRevisionRestoreFileEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FILE_END:
		return e.asRevisionRestoreFileEndEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_BEGIN:
		return e.asRevisionRestoreFragmentBeginEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_PROGRESS:
		return e.asRevisionRestoreFragmentProgressEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_END:
		return e.asRevisionRestoreFragmentEndEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_REVISION:
		return e.asRevisionRestoreRevisionEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_SYNC_BEGIN:
		return e.asRevisionRestoreSyncBeginEventDataFFI()
	case LoreEventTag_REVISION_RESTORE_SYNC_END:
		return e.asRevisionRestoreSyncEndEventDataFFI()
	case LoreEventTag_REVISION_RESOLVE:
		return e.asRevisionResolveEventDataFFI()
	case LoreEventTag_REVISION_SYNC_TARGET:
		return e.asRevisionSyncTargetEventDataFFI()
	case LoreEventTag_REVISION_SYNC_FILE:
		return e.asRevisionSyncFileEventDataFFI()
	case LoreEventTag_REVISION_SYNC_PROGRESS:
		return e.asRevisionSyncProgressEventDataFFI()
	case LoreEventTag_REVISION_SYNC_REVISION:
		return e.asRevisionSyncRevisionEventDataFFI()
	case LoreEventTag_REVISION_BISECT:
		return e.asRevisionBisectEventDataFFI()
	case LoreEventTag_NOTIFICATION_BRANCH_CREATED:
		return e.asNotificationBranchCreatedEventDataFFI()
	case LoreEventTag_NOTIFICATION_BRANCH_DELETED:
		return e.asNotificationBranchDeletedEventDataFFI()
	case LoreEventTag_NOTIFICATION_BRANCH_PUSHED:
		return e.asNotificationBranchPushedEventDataFFI()
	case LoreEventTag_NOTIFICATION_RESOURCE_LOCKED:
		return e.asNotificationResourceLockedEventDataFFI()
	case LoreEventTag_NOTIFICATION_RESOURCE_UNLOCKED:
		return e.asNotificationResourceUnlockedEventDataFFI()
	case LoreEventTag_NOTIFICATION_SUBSCRIBED:
		return e.asNotificationSubscribedEventDataFFI()
	case LoreEventTag_NOTIFICATION_UNSUBSCRIBED:
		return e.asNotificationUnsubscribedEventDataFFI()
	case LoreEventTag_SHARED_STORE_CREATE:
		return e.asSharedStoreCreateEventDataFFI()
	case LoreEventTag_SHARED_STORE_INFO:
		return e.asSharedStoreInfoEventDataFFI()
	case LoreEventTag_LINK_STAGED_ENTRY:
		return e.asLinkStagedEntryEventDataFFI()
	case LoreEventTag_STORAGE_OPENED:
		return e.asStorageOpenedEventDataFFI()
	case LoreEventTag_STORAGE_PUT_ITEM_COMPLETE:
		return e.asStoragePutItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_GET_HEADER:
		return e.asStorageGetHeaderEventDataFFI()
	case LoreEventTag_STORAGE_GET_DATA:
		return e.asStorageGetDataEventDataFFI()
	case LoreEventTag_STORAGE_GET_ITEM_COMPLETE:
		return e.asStorageGetItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_GET_METADATA_ITEM_COMPLETE:
		return e.asStorageGetMetadataItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_COPY_ITEM_COMPLETE:
		return e.asStorageCopyItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_OBLITERATE_ITEM_COMPLETE:
		return e.asStorageObliterateItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_UPLOAD_ITEM_COMPLETE:
		return e.asStorageUploadItemCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_LOADED:
		return e.asRevisionTreeLoadedEventDataFFI()
	case LoreEventTag_REVISION_TREE_RESOLVE_PATH_COMPLETE:
		return e.asRevisionTreeResolvePathCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_CHILD:
		return e.asRevisionTreeChildEventDataFFI()
	case LoreEventTag_REVISION_TREE_NODE_INFO:
		return e.asRevisionTreeNodeInfoEventDataFFI()
	case LoreEventTag_REVISION_TREE_NODE_PATH:
		return e.asRevisionTreeNodePathEventDataFFI()
	case LoreEventTag_REVISION_TREE_ADD_COMPLETE:
		return e.asRevisionTreeAddCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_DELETE_COMPLETE:
		return e.asRevisionTreeDeleteCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_MODIFY_COMPLETE:
		return e.asRevisionTreeModifyCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_MOVE_COMPLETE:
		return e.asRevisionTreeMoveCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_METADATA_SET_COMPLETE:
		return e.asRevisionTreeMetadataSetCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_METADATA_GET_COMPLETE:
		return e.asRevisionTreeMetadataGetCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_COMMIT_COMPLETE:
		return e.asRevisionTreeCommitCompleteEventDataFFI()
	case LoreEventTag_REVISION_TREE_CLOSE_COMPLETE:
		return e.asRevisionTreeCloseCompleteEventDataFFI()
	case LoreEventTag_STORAGE_MUTABLE_LOAD_ITEM_COMPLETE:
		return e.asStorageMutableLoadItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_MUTABLE_STORE_ITEM_COMPLETE:
		return e.asStorageMutableStoreItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE:
		return e.asStorageMutableCompareAndSwapItemCompleteEventDataFFI()
	case LoreEventTag_STORAGE_MUTABLE_LIST_ENTRY:
		return e.asStorageMutableListEntryEventDataFFI()
	case LoreEventTag_STORAGE_MUTABLE_LIST_ITEM_COMPLETE:
		return e.asStorageMutableListItemCompleteEventDataFFI()
	case LoreEventTag_EVICTION_BEGIN:
		return e.asEvictionBeginEventDataFFI()
	case LoreEventTag_EVICTION_PROGRESS:
		return e.asEvictionProgressEventDataFFI()
	case LoreEventTag_EVICTION_END:
		return e.asEvictionEndEventDataFFI()
	case LoreEventTag_COMPACTION_BEGIN:
		return e.asCompactionBeginEventDataFFI()
	case LoreEventTag_COMPACTION_PROGRESS:
		return e.asCompactionProgressEventDataFFI()
	case LoreEventTag_COMPACTION_END:
		return e.asCompactionEndEventDataFFI()
	default:
		return nil
	}
}

func (e *LoreProgressEventDataFFI) Clone() LoreProgressEventData {
	return LoreProgressEventData{
		Unused: e.Unused,
	}
}
func (e *LoreErrorEventDataFFI) Clone() LoreErrorEventData {
	return LoreErrorEventData{
		ErrorType:  e.ErrorType,
		ErrorInner: e.ErrorInner.Clone(),
	}
}
func (e *LoreCompleteEventDataFFI) Clone() LoreCompleteEventData {
	return LoreCompleteEventData{
		Status: e.Status,
		Error:  e.Error.Clone(),
	}
}
func (e *LoreMetadataEventDataFFI) Clone() LoreMetadataEventData {
	return LoreMetadataEventData{
		Key:   e.Key.Clone(),
		Value: e.Value.Clone(),
	}
}
func (e *LoreLogEventDataFFI) Clone() LoreLogEventData {
	return LoreLogEventData{
		Level:     e.Level,
		Category:  e.Category,
		Timestamp: e.Timestamp,
		Location:  e.Location.Clone(),
		Message:   e.Message.Clone(),
	}
}
func (e *LoreEndEventDataFFI) Clone() LoreEndEventData {
	return LoreEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreMaintenanceEventDataFFI) Clone() LoreMaintenanceEventData {
	return LoreMaintenanceEventData{
		Message: e.Message.Clone(),
	}
}
func (e *LoreAuthUrlEventDataFFI) Clone() LoreAuthUrlEventData {
	return LoreAuthUrlEventData{
		Url: e.Url.Clone(),
	}
}
func (e *LoreAuthUserInfoEventDataFFI) Clone() LoreAuthUserInfoEventData {
	return LoreAuthUserInfoEventData{
		Id:   e.Id.Clone(),
		Name: e.Name.Clone(),
	}
}
func (e *LoreAuthUserTokenEventDataFFI) Clone() LoreAuthUserTokenEventData {
	return LoreAuthUserTokenEventData{
		Id:                 e.Id.Clone(),
		Name:               e.Name.Clone(),
		Token:              e.Token.Clone(),
		PreferredUsername:  e.PreferredUsername.Clone(),
		FlagServiceAccount: e.FlagServiceAccount != 0,
		Expires:            e.Expires,
	}
}
func (e *LoreAuthIdentityEventDataFFI) Clone() LoreAuthIdentityEventData {
	return LoreAuthIdentityEventData{
		AuthUrl:           e.AuthUrl.Clone(),
		Resource:          e.Resource.Clone(),
		UserId:            e.UserId.Clone(),
		AuthorizedDomains: e.AuthorizedDomains.Clone(),
		Expires:           e.Expires,
		Token:             e.Token.Clone(),
	}
}
func (e *LoreBranchCreateEventDataFFI) Clone() LoreBranchCreateEventData {
	return LoreBranchCreateEventData{
		Name:     e.Name.Clone(),
		Latest:   e.Latest,
		IsCommit: e.IsCommit != 0,
	}
}
func (e *LoreBranchMultipleInstanceEventDataFFI) Clone() LoreBranchMultipleInstanceEventData {
	return LoreBranchMultipleInstanceEventData{
		Branch:        e.Branch,
		InstanceIds:   e.InstanceIds.Clone(),
		InstancePaths: e.InstancePaths.Clone(),
	}
}
func (e *LoreBranchArchiveEventDataFFI) Clone() LoreBranchArchiveEventData {
	return LoreBranchArchiveEventData{
		Name: e.Name.Clone(),
	}
}
func (e *LoreBranchListBeginEventDataFFI) Clone() LoreBranchListBeginEventData {
	return LoreBranchListBeginEventData{
		Location: e.Location,
	}
}
func (e *LoreBranchListEntryEventDataFFI) Clone() LoreBranchListEntryEventData {
	return LoreBranchListEntryEventData{
		Location:  e.Location,
		Id:        e.Id,
		Name:      e.Name.Clone(),
		Category:  e.Category.Clone(),
		Latest:    e.Latest,
		Stack:     e.Stack.Clone(),
		Creator:   e.Creator.Clone(),
		Created:   e.Created,
		IsCurrent: e.IsCurrent != 0,
		Archived:  e.Archived != 0,
	}
}
func (e *LoreBranchListEndEventDataFFI) Clone() LoreBranchListEndEventData {
	return LoreBranchListEndEventData{
		Location: e.Location,
		Count:    e.Count,
	}
}
func (e *LoreBranchMergeAbortBeginEventDataFFI) Clone() LoreBranchMergeAbortBeginEventData {
	return LoreBranchMergeAbortBeginEventData{
		StateStagedRevision:  e.StateStagedRevision,
		StateCurrentRevision: e.StateCurrentRevision,
	}
}
func (e *LoreBranchMergeAbortEndEventDataFFI) Clone() LoreBranchMergeAbortEndEventData {
	return LoreBranchMergeAbortEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreBranchInfoEventDataFFI) Clone() LoreBranchInfoEventData {
	return LoreBranchInfoEventData{
		Id:           e.Id,
		Name:         e.Name.Clone(),
		Category:     e.Category.Clone(),
		Latest:       e.Latest,
		LatestRemote: e.LatestRemote,
		Parent:       e.Parent,
		BranchPoint:  e.BranchPoint,
		Creator:      e.Creator.Clone(),
		Created:      e.Created,
		Stack:        e.Stack.Clone(),
		Archived:     e.Archived != 0,
	}
}
func (e *LoreBranchDiffBeginEventDataFFI) Clone() LoreBranchDiffBeginEventData {
	return LoreBranchDiffBeginEventData{
		Unused: e.Unused,
	}
}
func (e *LoreBranchDiffChangeBeginEventDataFFI) Clone() LoreBranchDiffChangeBeginEventData {
	return LoreBranchDiffChangeBeginEventData{
		ChangesCount: e.ChangesCount,
	}
}
func (e *LoreBranchDiffChangeEventDataFFI) Clone() LoreBranchDiffChangeEventData {
	return LoreBranchDiffChangeEventData{
		Change: e.Change.Clone(),
	}
}
func (e *LoreBranchDiffChangeEndEventDataFFI) Clone() LoreBranchDiffChangeEndEventData {
	return LoreBranchDiffChangeEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreBranchDiffConflictBeginEventDataFFI) Clone() LoreBranchDiffConflictBeginEventData {
	return LoreBranchDiffConflictBeginEventData{
		ConflictsCount: e.ConflictsCount,
	}
}
func (e *LoreBranchDiffConflictEventDataFFI) Clone() LoreBranchDiffConflictEventData {
	return LoreBranchDiffConflictEventData{
		SourceChange: e.SourceChange.Clone(),
		TargetChange: e.TargetChange.Clone(),
	}
}
func (e *LoreBranchDiffConflictEndEventDataFFI) Clone() LoreBranchDiffConflictEndEventData {
	return LoreBranchDiffConflictEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreBranchDiffEndEventDataFFI) Clone() LoreBranchDiffEndEventData {
	return LoreBranchDiffEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreBranchLatestListEntryEventDataFFI) Clone() LoreBranchLatestListEntryEventData {
	return LoreBranchLatestListEntryEventData{
		Branch:   e.Branch,
		Revision: e.Revision,
	}
}
func (e *LoreBranchMergeConflictFileEventDataFFI) Clone() LoreBranchMergeConflictFileEventData {
	return LoreBranchMergeConflictFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreBranchMergeLinkSkippedEventDataFFI) Clone() LoreBranchMergeLinkSkippedEventData {
	return LoreBranchMergeLinkSkippedEventData{
		LinkPath:   e.LinkPath.Clone(),
		Repository: e.Repository,
		Reason:     e.Reason != 0,
	}
}
func (e *LoreBranchMergeUnresolveFileEventDataFFI) Clone() LoreBranchMergeUnresolveFileEventData {
	return LoreBranchMergeUnresolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreBranchMergeUnresolveRevisionEventDataFFI) Clone() LoreBranchMergeUnresolveRevisionEventData {
	return LoreBranchMergeUnresolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreBranchMergeIntoFileBeginEventDataFFI) Clone() LoreBranchMergeIntoFileBeginEventData {
	return LoreBranchMergeIntoFileBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreBranchMergeIntoFileEventDataFFI) Clone() LoreBranchMergeIntoFileEventData {
	return LoreBranchMergeIntoFileEventData{
		Path:        e.Path.Clone(),
		Action:      e.Action,
		Size:        e.Size,
		IsFile:      e.IsFile != 0,
		IsDirectory: e.IsDirectory != 0,
		IsLink:      e.IsLink != 0,
	}
}
func (e *LoreBranchMergeIntoFileEndEventDataFFI) Clone() LoreBranchMergeIntoFileEndEventData {
	return LoreBranchMergeIntoFileEndEventData{
		Count: e.Count,
	}
}
func (e *LoreBranchMergeIntoFragmentBeginEventDataFFI) Clone() LoreBranchMergeIntoFragmentBeginEventData {
	return LoreBranchMergeIntoFragmentBeginEventData{
		Fragments: e.Fragments,
	}
}
func (e *LoreBranchMergeIntoFragmentProgressEventDataFFI) Clone() LoreBranchMergeIntoFragmentProgressEventData {
	return LoreBranchMergeIntoFragmentProgressEventData{
		Complete: e.Complete,
		Count:    e.Count,
	}
}
func (e *LoreBranchMergeIntoFragmentEndEventDataFFI) Clone() LoreBranchMergeIntoFragmentEndEventData {
	return LoreBranchMergeIntoFragmentEndEventData{
		Fragments: e.Fragments,
	}
}
func (e *LoreBranchMergeIntoRevisionEventDataFFI) Clone() LoreBranchMergeIntoRevisionEventData {
	return LoreBranchMergeIntoRevisionEventData{
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
	}
}
func (e *LoreBranchMergeIntoSyncBeginEventDataFFI) Clone() LoreBranchMergeIntoSyncBeginEventData {
	return LoreBranchMergeIntoSyncBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreBranchMergeIntoSyncEndEventDataFFI) Clone() LoreBranchMergeIntoSyncEndEventData {
	return LoreBranchMergeIntoSyncEndEventData{
		Count: e.Count,
	}
}
func (e *LoreBranchMergeResolveFileEventDataFFI) Clone() LoreBranchMergeResolveFileEventData {
	return LoreBranchMergeResolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreBranchMergeResolveRevisionEventDataFFI) Clone() LoreBranchMergeResolveRevisionEventData {
	return LoreBranchMergeResolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreBranchMergeStartBeginEventDataFFI) Clone() LoreBranchMergeStartBeginEventData {
	return LoreBranchMergeStartBeginEventData{
		Branch:         e.Branch,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
	}
}
func (e *LoreRevisionSyncProgressEventDataFFI) Clone() LoreRevisionSyncProgressEventData {
	return LoreRevisionSyncProgressEventData{
		FileUpdate:        e.FileUpdate,
		FileUpdateTotal:   e.FileUpdateTotal,
		FileDelete:        e.FileDelete,
		FileDeleteTotal:   e.FileDeleteTotal,
		FileAutomerge:     e.FileAutomerge,
		FileConflict:      e.FileConflict,
		BytesUpdate:       e.BytesUpdate,
		BytesUpdateTotal:  e.BytesUpdateTotal,
		DiscoveryComplete: e.DiscoveryComplete != 0,
	}
}
func (e *LoreBranchMergeStartEndEventDataFFI) Clone() LoreBranchMergeStartEndEventData {
	return LoreBranchMergeStartEndEventData{
		Stats:        e.Stats.Clone(),
		Signature:    e.Signature,
		HasConflicts: e.HasConflicts != 0,
	}
}
func (e *LoreCherryPickStartBeginEventDataFFI) Clone() LoreCherryPickStartBeginEventData {
	return LoreCherryPickStartBeginEventData{
		Branch:         e.Branch,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
	}
}
func (e *LoreCherryPickStartEndEventDataFFI) Clone() LoreCherryPickStartEndEventData {
	return LoreCherryPickStartEndEventData{
		Stats:        e.Stats.Clone(),
		Signature:    e.Signature,
		HasConflicts: e.HasConflicts != 0,
	}
}
func (e *LoreCherryPickAbortBeginEventDataFFI) Clone() LoreCherryPickAbortBeginEventData {
	return LoreCherryPickAbortBeginEventData{
		StateStagedRevision:  e.StateStagedRevision,
		StateCurrentRevision: e.StateCurrentRevision,
	}
}
func (e *LoreCherryPickAbortEndEventDataFFI) Clone() LoreCherryPickAbortEndEventData {
	return LoreCherryPickAbortEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreCherryPickConflictFileEventDataFFI) Clone() LoreCherryPickConflictFileEventData {
	return LoreCherryPickConflictFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreCherryPickUnresolveFileEventDataFFI) Clone() LoreCherryPickUnresolveFileEventData {
	return LoreCherryPickUnresolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreCherryPickUnresolveRevisionEventDataFFI) Clone() LoreCherryPickUnresolveRevisionEventData {
	return LoreCherryPickUnresolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreCherryPickResolveFileEventDataFFI) Clone() LoreCherryPickResolveFileEventData {
	return LoreCherryPickResolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreCherryPickResolveRevisionEventDataFFI) Clone() LoreCherryPickResolveRevisionEventData {
	return LoreCherryPickResolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreRevertStartBeginEventDataFFI) Clone() LoreRevertStartBeginEventData {
	return LoreRevertStartBeginEventData{
		Branch:         e.Branch,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
	}
}
func (e *LoreRevertStartEndEventDataFFI) Clone() LoreRevertStartEndEventData {
	return LoreRevertStartEndEventData{
		Stats:        e.Stats.Clone(),
		Signature:    e.Signature,
		HasConflicts: e.HasConflicts != 0,
	}
}
func (e *LoreRevertAbortBeginEventDataFFI) Clone() LoreRevertAbortBeginEventData {
	return LoreRevertAbortBeginEventData{
		StateStagedRevision:  e.StateStagedRevision,
		StateCurrentRevision: e.StateCurrentRevision,
	}
}
func (e *LoreRevertAbortEndEventDataFFI) Clone() LoreRevertAbortEndEventData {
	return LoreRevertAbortEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreRevertResolveFileEventDataFFI) Clone() LoreRevertResolveFileEventData {
	return LoreRevertResolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreRevertResolveRevisionEventDataFFI) Clone() LoreRevertResolveRevisionEventData {
	return LoreRevertResolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreRevertConflictFileEventDataFFI) Clone() LoreRevertConflictFileEventData {
	return LoreRevertConflictFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreRevertUnresolveFileEventDataFFI) Clone() LoreRevertUnresolveFileEventData {
	return LoreRevertUnresolveFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreRevertUnresolveRevisionEventDataFFI) Clone() LoreRevertUnresolveRevisionEventData {
	return LoreRevertUnresolveRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreBranchProtectEventDataFFI) Clone() LoreBranchProtectEventData {
	return LoreBranchProtectEventData{
		Name: e.Name.Clone(),
	}
}
func (e *LoreBranchPushEventDataFFI) Clone() LoreBranchPushEventData {
	return LoreBranchPushEventData{
		Remote:            e.Remote.Clone(),
		Repository:        e.Repository,
		Branch:            e.Branch,
		BranchName:        e.BranchName.Clone(),
		RemoteRevision:    e.RemoteRevision,
		LocalRevision:     e.LocalRevision,
		RemoteHistory:     e.RemoteHistory,
		LocalHistory:      e.LocalHistory,
		FlagAlreadyPushed: e.FlagAlreadyPushed != 0,
		FlagDefault:       e.FlagDefault != 0,
		FlagLink:          e.FlagLink != 0,
		FlagLayer:         e.FlagLayer != 0,
	}
}
func (e *LoreBranchPushRevisionUpdateBeginEventDataFFI) Clone() LoreBranchPushRevisionUpdateBeginEventData {
	return LoreBranchPushRevisionUpdateBeginEventData{
		Revision:  e.Revision,
		OldParent: e.OldParent,
		NewParent: e.NewParent,
	}
}
func (e *LoreBranchPushRevisionUpdateEndEventDataFFI) Clone() LoreBranchPushRevisionUpdateEndEventData {
	return LoreBranchPushRevisionUpdateEndEventData{
		Revision: e.Revision,
	}
}
func (e *LoreBranchPushFragmentBeginEventDataFFI) Clone() LoreBranchPushFragmentBeginEventData {
	return LoreBranchPushFragmentBeginEventData{
		Fragments:  e.Fragments,
		BytesTotal: e.BytesTotal,
	}
}
func (e *LoreBranchPushFragmentProgressEventDataFFI) Clone() LoreBranchPushFragmentProgressEventData {
	return LoreBranchPushFragmentProgressEventData{
		Complete:         e.Complete,
		Count:            e.Count,
		BytesTransferred: e.BytesTransferred,
		BytesTotal:       e.BytesTotal,
	}
}
func (e *LoreBranchPushFragmentEndEventDataFFI) Clone() LoreBranchPushFragmentEndEventData {
	return LoreBranchPushFragmentEndEventData{
		Fragments:        e.Fragments,
		BytesTransferred: e.BytesTransferred,
	}
}
func (e *LoreBranchPushBranchCreateBeginEventDataFFI) Clone() LoreBranchPushBranchCreateBeginEventData {
	return LoreBranchPushBranchCreateBeginEventData{
		LocalRevision: e.LocalRevision,
	}
}
func (e *LoreBranchPushBranchCreateEndEventDataFFI) Clone() LoreBranchPushBranchCreateEndEventData {
	return LoreBranchPushBranchCreateEndEventData{
		RemoteRevision: e.RemoteRevision,
	}
}
func (e *LoreBranchPushRevisionPushBeginEventDataFFI) Clone() LoreBranchPushRevisionPushBeginEventData {
	return LoreBranchPushRevisionPushBeginEventData{
		RemoteRevision: e.RemoteRevision,
		LocalRevision:  e.LocalRevision,
	}
}
func (e *LoreBranchPushRevisionPushUpdateEventDataFFI) Clone() LoreBranchPushRevisionPushUpdateEventData {
	return LoreBranchPushRevisionPushUpdateEventData{
		OldRevision:       e.OldRevision,
		NewRevision:       e.NewRevision,
		NewRevisionNumber: e.NewRevisionNumber,
	}
}
func (e *LoreBranchPushRevisionPushEndEventDataFFI) Clone() LoreBranchPushRevisionPushEndEventData {
	return LoreBranchPushRevisionPushEndEventData{
		OldRemoteRevision:       e.OldRemoteRevision,
		NewRemoteRevision:       e.NewRemoteRevision,
		NewRemoteRevisionNumber: e.NewRemoteRevisionNumber,
		Message:                 e.Message.Clone(),
		FastForwardMerged:       e.FastForwardMerged != 0,
	}
}
func (e *LoreBranchResetEventDataFFI) Clone() LoreBranchResetEventData {
	return LoreBranchResetEventData{
		Id:       e.Id,
		Name:     e.Name.Clone(),
		Revision: e.Revision,
	}
}
func (e *LoreBranchSwitchBeginEventDataFFI) Clone() LoreBranchSwitchBeginEventData {
	return LoreBranchSwitchBeginEventData{
		Branch: e.Branch.Clone(),
	}
}
func (e *LoreBranchSwitchEndEventDataFFI) Clone() LoreBranchSwitchEndEventData {
	return LoreBranchSwitchEndEventData{
		Branch: e.Branch.Clone(),
	}
}
func (e *LoreBranchUnprotectEventDataFFI) Clone() LoreBranchUnprotectEventData {
	return LoreBranchUnprotectEventData{
		Name: e.Name.Clone(),
	}
}
func (e *LoreFileInfoEventDataFFI) Clone() LoreFileInfoEventData {
	return LoreFileInfoEventData{
		Path:         e.Path.Clone(),
		Context:      e.Context,
		Hash:         e.Hash,
		IsFile:       e.IsFile != 0,
		IsDir:        e.IsDir != 0,
		FlagModified: e.FlagModified != 0,
		FlagDeleted:  e.FlagDeleted != 0,
		FlagAdded:    e.FlagAdded != 0,
		FlagConflict: e.FlagConflict != 0,
		Mode:         e.Mode,
		Size:         e.Size,
		LocalSize:    e.LocalSize,
		LocalHash:    e.LocalHash,
		FilterSize:   e.FilterSize,
	}
}
func (e *LoreFileDiffEventDataFFI) Clone() LoreFileDiffEventData {
	return LoreFileDiffEventData{
		Path:   e.Path.Clone(),
		Patch:  e.Patch.Clone(),
		Action: e.Action,
	}
}
func (e *LoreFileHashEventDataFFI) Clone() LoreFileHashEventData {
	return LoreFileHashEventData{
		Path: e.Path.Clone(),
		Size: e.Size,
		Hash: e.Hash,
	}
}
func (e *LoreFileHistoryEventDataFFI) Clone() LoreFileHistoryEventData {
	return LoreFileHistoryEventData{
		Path:           e.Path.Clone(),
		Repository:     e.Repository,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		Parent: [2]LoreHash{
			e.Parent[0].Clone(),
			e.Parent[1].Clone(),
		},
		Address: e.Address,
		Size:    e.Size,
		Action:  e.Action,
	}
}
func (e *LoreFileWriteEventDataFFI) Clone() LoreFileWriteEventData {
	return LoreFileWriteEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreFileObliterateEventDataFFI) Clone() LoreFileObliterateEventData {
	return LoreFileObliterateEventData{
		Address:      e.Address,
		NumFragments: e.NumFragments,
		NumPayloads:  e.NumPayloads,
	}
}
func (e *LoreFileDumpEventDataFFI) Clone() LoreFileDumpEventData {
	return LoreFileDumpEventData{
		Address:     e.Address,
		Flags:       e.Flags,
		SizePayload: e.SizePayload,
		SizeContent: e.SizeContent,
		MatchMade:   e.MatchMade != 0,
	}
}
func (e *LoreFileDependencyAddBeginEventDataFFI) Clone() LoreFileDependencyAddBeginEventData {
	return LoreFileDependencyAddBeginEventData{
		PathCount:       e.PathCount,
		DependencyCount: e.DependencyCount,
	}
}
func (e *LoreFileDependencyAddEntryEventDataFFI) Clone() LoreFileDependencyAddEntryEventData {
	return LoreFileDependencyAddEntryEventData{
		Path:       e.Path.Clone(),
		Dependency: e.Dependency.Clone(),
		Tags:       e.Tags.Clone(),
	}
}
func (e *LoreFileDependencyAddEndEventDataFFI) Clone() LoreFileDependencyAddEndEventData {
	return LoreFileDependencyAddEndEventData{
		AddedCount: e.AddedCount,
	}
}
func (e *LoreFileDependencyRemoveBeginEventDataFFI) Clone() LoreFileDependencyRemoveBeginEventData {
	return LoreFileDependencyRemoveBeginEventData{
		PathCount:       e.PathCount,
		DependencyCount: e.DependencyCount,
	}
}
func (e *LoreFileDependencyRemoveEntryEventDataFFI) Clone() LoreFileDependencyRemoveEntryEventData {
	return LoreFileDependencyRemoveEntryEventData{
		Path:       e.Path.Clone(),
		Dependency: e.Dependency.Clone(),
		Tags:       e.Tags.Clone(),
	}
}
func (e *LoreFileDependencyRemoveEndEventDataFFI) Clone() LoreFileDependencyRemoveEndEventData {
	return LoreFileDependencyRemoveEndEventData{
		RemovedCount: e.RemovedCount,
	}
}
func (e *LoreFileDependencyListBeginEventDataFFI) Clone() LoreFileDependencyListBeginEventData {
	return LoreFileDependencyListBeginEventData{
		FileCount: e.FileCount,
	}
}
func (e *LoreFileDependencyListFileEventDataFFI) Clone() LoreFileDependencyListFileEventData {
	return LoreFileDependencyListFileEventData{
		Path:       e.Path.Clone(),
		EntryCount: e.EntryCount,
	}
}
func (e *LoreFileDependencyListEntryEventDataFFI) Clone() LoreFileDependencyListEntryEventData {
	return LoreFileDependencyListEntryEventData{
		Path:  e.Path.Clone(),
		Tags:  e.Tags.Clone(),
		Depth: e.Depth,
	}
}
func (e *LoreFileDependencyListFileEndEventDataFFI) Clone() LoreFileDependencyListFileEndEventData {
	return LoreFileDependencyListFileEndEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreFileDependencyListEndEventDataFFI) Clone() LoreFileDependencyListEndEventData {
	return LoreFileDependencyListEndEventData{
		TotalEntryCount: e.TotalEntryCount,
	}
}
func (e *LoreFileResetBeginEventDataFFI) Clone() LoreFileResetBeginEventData {
	return LoreFileResetBeginEventData{
		PathCount: e.PathCount,
	}
}
func (e *LoreFileResetProgressEventDataFFI) Clone() LoreFileResetProgressEventData {
	return LoreFileResetProgressEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileResetEndEventDataFFI) Clone() LoreFileResetEndEventData {
	return LoreFileResetEndEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileResetFileEventDataFFI) Clone() LoreFileResetFileEventData {
	return LoreFileResetFileEventData{
		Path:     e.Path.Clone(),
		Action:   e.Action,
		FromPath: e.FromPath.Clone(),
	}
}
func (e *LoreFilterExcludeEventDataFFI) Clone() LoreFilterExcludeEventData {
	return LoreFilterExcludeEventData{
		Reason: e.Reason != 0,
		Path:   e.Path.Clone(),
	}
}
func (e *LoreFileStageBeginEventDataFFI) Clone() LoreFileStageBeginEventData {
	return LoreFileStageBeginEventData{
		PathCount: e.PathCount,
	}
}
func (e *LoreFileStageProgressEventDataFFI) Clone() LoreFileStageProgressEventData {
	return LoreFileStageProgressEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileStageEndEventDataFFI) Clone() LoreFileStageEndEventData {
	return LoreFileStageEndEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileStageRevisionEventDataFFI) Clone() LoreFileStageRevisionEventData {
	return LoreFileStageRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreFileStageFileEventDataFFI) Clone() LoreFileStageFileEventData {
	return LoreFileStageFileEventData{
		FromPath: e.FromPath.Clone(),
		Path:     e.Path.Clone(),
		Action:   e.Action,
	}
}
func (e *LoreFileUnstageBeginEventDataFFI) Clone() LoreFileUnstageBeginEventData {
	return LoreFileUnstageBeginEventData{
		PathCount: e.PathCount,
	}
}
func (e *LoreFileUnstageProgressEventDataFFI) Clone() LoreFileUnstageProgressEventData {
	return LoreFileUnstageProgressEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileUnstageEndEventDataFFI) Clone() LoreFileUnstageEndEventData {
	return LoreFileUnstageEndEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreFileUnstageRevisionEventDataFFI) Clone() LoreFileUnstageRevisionEventData {
	return LoreFileUnstageRevisionEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreFileUnstageFileEventDataFFI) Clone() LoreFileUnstageFileEventData {
	return LoreFileUnstageFileEventData{
		Path:   e.Path.Clone(),
		Action: e.Action,
	}
}
func (e *LoreFragmentWriteEventDataFFI) Clone() LoreFragmentWriteEventData {
	return LoreFragmentWriteEventData{
		Fragment:     e.Fragment.Clone(),
		Deduplicated: e.Deduplicated != 0,
	}
}
func (e *LoreLayerAddEventDataFFI) Clone() LoreLayerAddEventData {
	return LoreLayerAddEventData{
		TargetPath:       e.TargetPath.Clone(),
		SourceRepository: e.SourceRepository,
		SourcePath:       e.SourcePath.Clone(),
		Metadata:         e.Metadata.Clone(),
		Revision:         e.Revision,
	}
}
func (e *LoreLayerEntryEventDataFFI) Clone() LoreLayerEntryEventData {
	return LoreLayerEntryEventData{
		TargetPath:       e.TargetPath.Clone(),
		SourceRepository: e.SourceRepository,
		SourcePath:       e.SourcePath.Clone(),
		Metadata:         e.Metadata.Clone(),
		Revision:         e.Revision,
	}
}
func (e *LoreLayerRemoveEventDataFFI) Clone() LoreLayerRemoveEventData {
	return LoreLayerRemoveEventData{
		TargetPath:       e.TargetPath.Clone(),
		SourceRepository: e.SourceRepository,
		SourcePath:       e.SourcePath.Clone(),
		Revision:         e.Revision,
		Forced:           e.Forced != 0,
		Purged:           e.Purged != 0,
		FileCount:        e.FileCount,
		DirectoryCount:   e.DirectoryCount,
		ModifiedCount:    e.ModifiedCount,
	}
}
func (e *LoreLayerStagedEntryEventDataFFI) Clone() LoreLayerStagedEntryEventData {
	return LoreLayerStagedEntryEventData{
		TargetPath:       e.TargetPath.Clone(),
		SourceRepository: e.SourceRepository,
		StagedFileCount:  e.StagedFileCount,
	}
}
func (e *LoreLinkChangeEventDataFFI) Clone() LoreLinkChangeEventData {
	return LoreLinkChangeEventData{
		LinkPath:       e.LinkPath.Clone(),
		LinkRepository: e.LinkRepository,
		Branch:         e.Branch,
		Revision:       e.Revision,
		Action:         e.Action,
	}
}
func (e *LoreLinkEntryEventDataFFI) Clone() LoreLinkEntryEventData {
	return LoreLinkEntryEventData{
		Link:       e.Link,
		LinkNode:   e.LinkNode,
		LinkPath:   e.LinkPath.Clone(),
		SourceNode: e.SourceNode,
		SourcePath: e.SourcePath.Clone(),
		Branch:     e.Branch,
		BranchName: e.BranchName.Clone(),
		Revision:   e.Revision,
		Flags:      e.Flags,
	}
}
func (e *LoreLockFileAcquireBeginEventDataFFI) Clone() LoreLockFileAcquireBeginEventData {
	return LoreLockFileAcquireBeginEventData{
		Count:   e.Count,
		Ignored: e.Ignored != 0,
	}
}
func (e *LoreLockFileAcquireEventDataFFI) Clone() LoreLockFileAcquireEventData {
	return LoreLockFileAcquireEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreLockFileStatusBeginEventDataFFI) Clone() LoreLockFileStatusBeginEventData {
	return LoreLockFileStatusBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreLockFileStatusEventDataFFI) Clone() LoreLockFileStatusEventData {
	return LoreLockFileStatusEventData{
		Path:     e.Path.Clone(),
		Owner:    e.Owner.Clone(),
		LockedAt: e.LockedAt,
	}
}
func (e *LoreLockFileQueryBeginEventDataFFI) Clone() LoreLockFileQueryBeginEventData {
	return LoreLockFileQueryBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreLockFileQueryEventDataFFI) Clone() LoreLockFileQueryEventData {
	return LoreLockFileQueryEventData{
		Branch:   e.Branch,
		Path:     e.Path.Clone(),
		Owner:    e.Owner.Clone(),
		LockedAt: e.LockedAt,
	}
}
func (e *LoreLockFileReleaseBeginEventDataFFI) Clone() LoreLockFileReleaseBeginEventData {
	return LoreLockFileReleaseBeginEventData{
		Count:    e.Count,
		NotFound: e.NotFound != 0,
	}
}
func (e *LoreLockFileReleaseEventDataFFI) Clone() LoreLockFileReleaseEventData {
	return LoreLockFileReleaseEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreMetadataClearFileEventDataFFI) Clone() LoreMetadataClearFileEventData {
	return LoreMetadataClearFileEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreMetadataClearRevisionEventDataFFI) Clone() LoreMetadataClearRevisionEventData {
	return LoreMetadataClearRevisionEventData{
		Revision: e.Revision,
	}
}
func (e *LorePathIgnoreEventDataFFI) Clone() LorePathIgnoreEventData {
	return LorePathIgnoreEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreRepositoryCreateEventDataFFI) Clone() LoreRepositoryCreateEventData {
	return LoreRepositoryCreateEventData{
		Id:   e.Id,
		Name: e.Name.Clone(),
		Path: e.Path.Clone(),
	}
}
func (e *LoreRepositoryCloneBeginEventDataFFI) Clone() LoreRepositoryCloneBeginEventData {
	return LoreRepositoryCloneBeginEventData{
		Repository: e.Repository,
		Branch:     e.Branch.Clone(),
		Revision:   e.Revision,
		Path:       e.Path.Clone(),
	}
}
func (e *LoreRepositoryCloneProgressEventDataFFI) Clone() LoreRepositoryCloneProgressEventData {
	return LoreRepositoryCloneProgressEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreRepositoryCloneEndEventDataFFI) Clone() LoreRepositoryCloneEndEventData {
	return LoreRepositoryCloneEndEventData{
		Branch:   e.Branch.Clone(),
		Revision: e.Revision,
		Count:    e.Count.Clone(),
	}
}
func (e *LoreDependencyResolveBeginEventDataFFI) Clone() LoreDependencyResolveBeginEventData {
	return LoreDependencyResolveBeginEventData{
		RootCount: e.RootCount,
	}
}
func (e *LoreDependencyResolveItemEventDataFFI) Clone() LoreDependencyResolveItemEventData {
	return LoreDependencyResolveItemEventData{
		Source: e.Source.Clone(),
		Target: e.Target.Clone(),
		Tags:   e.Tags.Clone(),
	}
}
func (e *LoreDependencyResolveEndEventDataFFI) Clone() LoreDependencyResolveEndEventData {
	return LoreDependencyResolveEndEventData{
		ResolvedCount: e.ResolvedCount,
	}
}
func (e *LoreRepositoryDataEventDataFFI) Clone() LoreRepositoryDataEventData {
	return LoreRepositoryDataEventData{
		RemoteUrl:         e.RemoteUrl.Clone(),
		Id:                e.Id,
		Name:              e.Name.Clone(),
		Description:       e.Description.Clone(),
		DefaultBranch:     e.DefaultBranch,
		DefaultBranchName: e.DefaultBranchName.Clone(),
		Creator:           e.Creator.Clone(),
		Created:           e.Created,
	}
}
func (e *LoreRepositoryConfigGetEventDataFFI) Clone() LoreRepositoryConfigGetEventData {
	return LoreRepositoryConfigGetEventData{
		Key:   e.Key.Clone(),
		Value: e.Value.Clone(),
	}
}
func (e *LoreRepositoryDumpBeginEventDataFFI) Clone() LoreRepositoryDumpBeginEventData {
	return LoreRepositoryDumpBeginEventData{
		Repository: e.Repository,
		Revision:   e.Revision,
	}
}
func (e *LoreRepositoryDumpEndEventDataFFI) Clone() LoreRepositoryDumpEndEventData {
	return LoreRepositoryDumpEndEventData{
		Unused: e.Unused,
	}
}
func (e *LoreRepositoryListEntryEventDataFFI) Clone() LoreRepositoryListEntryEventData {
	return LoreRepositoryListEntryEventData{
		Id:   e.Id,
		Name: e.Name.Clone(),
	}
}
func (e *LoreRepositoryInstanceEventDataFFI) Clone() LoreRepositoryInstanceEventData {
	return LoreRepositoryInstanceEventData{
		InstanceId: e.InstanceId,
		Path:       e.Path.Clone(),
		BranchName: e.BranchName.Clone(),
		Branch:     e.Branch,
		Revision:   e.Revision,
		Stale:      e.Stale != 0,
	}
}
func (e *LoreRepositoryVerifyStateBeginEventDataFFI) Clone() LoreRepositoryVerifyStateBeginEventData {
	return LoreRepositoryVerifyStateBeginEventData{
		Unused: e.Unused,
	}
}
func (e *LoreRepositoryVerifyStateEndEventDataFFI) Clone() LoreRepositoryVerifyStateEndEventData {
	return LoreRepositoryVerifyStateEndEventData{
		HealedStagedState: e.HealedStagedState,
	}
}
func (e *LoreRepositoryVerifyFragmentMatchEventDataFFI) Clone() LoreRepositoryVerifyFragmentMatchEventData {
	return LoreRepositoryVerifyFragmentMatchEventData{
		Slot:           e.Slot,
		Index:          e.Index,
		Repository:     e.Repository,
		AddressHash:    e.AddressHash,
		AddressContext: e.AddressContext,
		Flags:          e.Flags,
		SizePayload:    e.SizePayload,
		SizeContent:    e.SizeContent,
		PackOffset:     e.PackOffset,
		PackFile:       e.PackFile,
		LastAccess:     e.LastAccess,
	}
}
func (e *LoreRepositoryVerifyFragmentEventDataFFI) Clone() LoreRepositoryVerifyFragmentEventData {
	return LoreRepositoryVerifyFragmentEventData{
		Hash:               e.Hash,
		GroupIndex:         e.GroupIndex,
		BucketIndex:        e.BucketIndex,
		IndexPath:          e.IndexPath.Clone(),
		EntryCount:         e.EntryCount,
		PackfileEntryCount: e.PackfileEntryCount,
		MatchCount:         e.MatchCount,
		Matches:            e.Matches.Clone(),
		Error:              e.Error.Clone(),
	}
}
func (e *LoreRepositoryVerifyFragmentRemoteEventDataFFI) Clone() LoreRepositoryVerifyFragmentRemoteEventData {
	return LoreRepositoryVerifyFragmentRemoteEventData{
		AddressHash:    e.AddressHash,
		AddressContext: e.AddressContext,
		Corrupted:      e.Corrupted != 0,
		Healed:         e.Healed != 0,
		Error:          e.Error.Clone(),
	}
}
func (e *LoreRepositoryStateDumpEventDataFFI) Clone() LoreRepositoryStateDumpEventData {
	return LoreRepositoryStateDumpEventData{
		RevisionNumber: e.RevisionNumber,
		Revision:       e.Revision,
		TreeHash:       e.TreeHash,
		TreeSize:       e.TreeSize,
	}
}
func (e *LoreRepositoryStateDumpNodeEventDataFFI) Clone() LoreRepositoryStateDumpNodeEventData {
	return LoreRepositoryStateDumpNodeEventData{
		Name:     e.Name.Clone(),
		Id:       e.Id,
		Parent:   e.Parent,
		Sibling:  e.Sibling,
		Mode:     e.Mode,
		Size:     e.Size,
		Flags:    e.Flags,
		TypeData: e.TypeData.Clone(),
	}
}
func (e *LoreRepositoryStatusRevisionEventDataFFI) Clone() LoreRepositoryStatusRevisionEventData {
	return LoreRepositoryStatusRevisionEventData{
		Repository:                 e.Repository,
		Branch:                     e.Branch,
		BranchName:                 e.BranchName.Clone(),
		Revision:                   e.Revision,
		RevisionNumber:             e.RevisionNumber,
		RevisionStaged:             e.RevisionStaged,
		RevisionMerged:             e.RevisionMerged,
		RevisionMergedParentBranch: e.RevisionMergedParentBranch,
		RevisionLocal:              e.RevisionLocal,
		RevisionLocalNumber:        e.RevisionLocalNumber,
		RevisionRemote:             e.RevisionRemote,
		RevisionRemoteNumber:       e.RevisionRemoteNumber,
		IsLocalAhead:               e.IsLocalAhead != 0,
		IsRemoteAhead:              e.IsRemoteAhead != 0,
		RemoteAvailable:            e.RemoteAvailable != 0,
		RemoteAuthorized:           e.RemoteAuthorized != 0,
		RemoteBranchExist:          e.RemoteBranchExist != 0,
	}
}
func (e *LoreRepositoryStatusFileEventDataFFI) Clone() LoreRepositoryStatusFileEventData {
	return LoreRepositoryStatusFileEventData{
		Path:                   e.Path.Clone(),
		Size:                   e.Size,
		Action:                 e.Action,
		Type:                   e.Type,
		FlagStaged:             e.FlagStaged != 0,
		FlagMerged:             e.FlagMerged != 0,
		FlagConflict:           e.FlagConflict != 0,
		FlagConflictUnresolved: e.FlagConflictUnresolved != 0,
		FlagConflictAutomerged: e.FlagConflictAutomerged != 0,
		FlagConflictMine:       e.FlagConflictMine != 0,
		FlagConflictTheirs:     e.FlagConflictTheirs != 0,
		FlagDirty:              e.FlagDirty != 0,
		FromPath:               e.FromPath.Clone(),
	}
}
func (e *LoreRepositoryStatusCountEventDataFFI) Clone() LoreRepositoryStatusCountEventData {
	return LoreRepositoryStatusCountEventData{
		Directories: e.Directories,
		Files:       e.Files,
	}
}
func (e *LoreRepositoryStatusSummaryEventDataFFI) Clone() LoreRepositoryStatusSummaryEventData {
	return LoreRepositoryStatusSummaryEventData{
		Adds:     e.Adds,
		Deletes:  e.Deletes,
		Modifies: e.Modifies,
		Moves:    e.Moves,
		Copies:   e.Copies,
	}
}
func (e *LoreRepositoryStoreImmutableQueryEventDataFFI) Clone() LoreRepositoryStoreImmutableQueryEventData {
	return LoreRepositoryStoreImmutableQueryEventData{
		Address:     e.Address,
		Remote:      e.Remote != 0,
		Status:      e.Status,
		Payload:     e.Payload != 0,
		Subfragment: e.Subfragment != 0,
		Flags:       e.Flags,
		PayloadSize: e.PayloadSize,
		ContentSize: e.ContentSize,
	}
}
func (e *LoreRevisionCommitBeginEventDataFFI) Clone() LoreRevisionCommitBeginEventData {
	return LoreRevisionCommitBeginEventData{
		Unused: e.Unused,
	}
}
func (e *LoreRevisionCommitProgressEventDataFFI) Clone() LoreRevisionCommitProgressEventData {
	return LoreRevisionCommitProgressEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreRevisionCommitEndEventDataFFI) Clone() LoreRevisionCommitEndEventData {
	return LoreRevisionCommitEndEventData{
		Count: e.Count.Clone(),
	}
}
func (e *LoreRevisionCommitRevisionEventDataFFI) Clone() LoreRevisionCommitRevisionEventData {
	return LoreRevisionCommitRevisionEventData{
		Repository:     e.Repository,
		Branch:         e.Branch,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		Parent:         e.Parent,
		ParentOther:    e.ParentOther,
	}
}
func (e *LoreRevisionInfoEventDataFFI) Clone() LoreRevisionInfoEventData {
	return LoreRevisionInfoEventData{
		Repository:     e.Repository,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		Parent: [2]LoreHash{
			e.Parent[0].Clone(),
			e.Parent[1].Clone(),
		},
	}
}
func (e *LoreRevisionInfoDeltaEventDataFFI) Clone() LoreRevisionInfoDeltaEventData {
	return LoreRevisionInfoDeltaEventData{
		Path:       e.Path.Clone(),
		Size:       e.Size,
		Action:     e.Action,
		FlagModify: e.FlagModify != 0,
		FlagMerged: e.FlagMerged != 0,
		FlagFile:   e.FlagFile != 0,
	}
}
func (e *LoreRevisionDiffFileEventDataFFI) Clone() LoreRevisionDiffFileEventData {
	return LoreRevisionDiffFileEventData{
		Path:       e.Path.Clone(),
		Action:     e.Action,
		OldIsFile:  e.OldIsFile != 0,
		NewIsFile:  e.NewIsFile != 0,
		OldAddress: e.OldAddress,
		NewAddress: e.NewAddress,
	}
}
func (e *LoreRevisionFindEventDataFFI) Clone() LoreRevisionFindEventData {
	return LoreRevisionFindEventData{
		Signature: e.Signature,
	}
}
func (e *LoreRevisionHistoryEventDataFFI) Clone() LoreRevisionHistoryEventData {
	return LoreRevisionHistoryEventData{
		Repository: e.Repository,
		Branch:     e.Branch,
	}
}
func (e *LoreRevisionHistoryEntryEventDataFFI) Clone() LoreRevisionHistoryEntryEventData {
	return LoreRevisionHistoryEntryEventData{
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		Parent: [2]LoreHash{
			e.Parent[0].Clone(),
			e.Parent[1].Clone(),
		},
	}
}
func (e *LoreRevisionRestoreFileBeginEventDataFFI) Clone() LoreRevisionRestoreFileBeginEventData {
	return LoreRevisionRestoreFileBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreRevisionRestoreFileEventDataFFI) Clone() LoreRevisionRestoreFileEventData {
	return LoreRevisionRestoreFileEventData{
		Path:        e.Path.Clone(),
		Action:      e.Action,
		Size:        e.Size,
		IsFile:      e.IsFile != 0,
		IsDirectory: e.IsDirectory != 0,
		IsModule:    e.IsModule != 0,
	}
}
func (e *LoreRevisionRestoreFileEndEventDataFFI) Clone() LoreRevisionRestoreFileEndEventData {
	return LoreRevisionRestoreFileEndEventData{
		Count: e.Count,
	}
}
func (e *LoreRevisionRestoreFragmentBeginEventDataFFI) Clone() LoreRevisionRestoreFragmentBeginEventData {
	return LoreRevisionRestoreFragmentBeginEventData{
		Fragments: e.Fragments,
	}
}
func (e *LoreRevisionRestoreFragmentProgressEventDataFFI) Clone() LoreRevisionRestoreFragmentProgressEventData {
	return LoreRevisionRestoreFragmentProgressEventData{
		Complete: e.Complete,
		Count:    e.Count,
	}
}
func (e *LoreRevisionRestoreFragmentEndEventDataFFI) Clone() LoreRevisionRestoreFragmentEndEventData {
	return LoreRevisionRestoreFragmentEndEventData{
		Fragments: e.Fragments,
	}
}
func (e *LoreRevisionRestoreRevisionEventDataFFI) Clone() LoreRevisionRestoreRevisionEventData {
	return LoreRevisionRestoreRevisionEventData{
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
	}
}
func (e *LoreRevisionRestoreSyncBeginEventDataFFI) Clone() LoreRevisionRestoreSyncBeginEventData {
	return LoreRevisionRestoreSyncBeginEventData{
		Count: e.Count,
	}
}
func (e *LoreRevisionRestoreSyncEndEventDataFFI) Clone() LoreRevisionRestoreSyncEndEventData {
	return LoreRevisionRestoreSyncEndEventData{
		Count: e.Count,
	}
}
func (e *LoreRevisionResolveEventDataFFI) Clone() LoreRevisionResolveEventData {
	return LoreRevisionResolveEventData{
		Repository:     e.Repository,
		Branch:         e.Branch,
		Revision:       e.Revision.Clone(),
		RevisionNumber: e.RevisionNumber,
		Remote:         e.Remote != 0,
		Local:          e.Local != 0,
	}
}
func (e *LoreRevisionSyncTargetEventDataFFI) Clone() LoreRevisionSyncTargetEventData {
	return LoreRevisionSyncTargetEventData{
		Remote:               e.Remote.Clone(),
		Repository:           e.Repository,
		Branch:               e.Branch,
		BranchName:           e.BranchName.Clone(),
		SourceRevision:       e.SourceRevision,
		SourceRevisionNumber: e.SourceRevisionNumber,
		TargetRevision:       e.TargetRevision,
		TargetRevisionNumber: e.TargetRevisionNumber,
		IsLatest:             e.IsLatest != 0,
		Local:                e.Local != 0,
	}
}
func (e *LoreRevisionSyncFileEventDataFFI) Clone() LoreRevisionSyncFileEventData {
	return LoreRevisionSyncFileEventData{
		Path:     e.Path.Clone(),
		Size:     e.Size,
		Action:   e.Action,
		FlagFile: e.FlagFile != 0,
	}
}
func (e *LoreRevisionSyncRevisionEventDataFFI) Clone() LoreRevisionSyncRevisionEventData {
	return LoreRevisionSyncRevisionEventData{
		Branch:         e.Branch,
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		FlagMerge:      e.FlagMerge != 0,
		FlagConflict:   e.FlagConflict != 0,
	}
}
func (e *LoreRevisionBisectEventDataFFI) Clone() LoreRevisionBisectEventData {
	return LoreRevisionBisectEventData{
		StartRevisionNumber:  e.StartRevisionNumber,
		TargetRevisionNumber: e.TargetRevisionNumber,
		EndRevisionNumber:    e.EndRevisionNumber,
		Done:                 e.Done != 0,
	}
}
func (e *LoreNotificationBranchCreatedEventDataFFI) Clone() LoreNotificationBranchCreatedEventData {
	return LoreNotificationBranchCreatedEventData{
		Branch: e.Branch,
	}
}
func (e *LoreNotificationBranchDeletedEventDataFFI) Clone() LoreNotificationBranchDeletedEventData {
	return LoreNotificationBranchDeletedEventData{
		Branch: e.Branch,
	}
}
func (e *LoreNotificationBranchPushedEventDataFFI) Clone() LoreNotificationBranchPushedEventData {
	return LoreNotificationBranchPushedEventData{
		Revision:       e.Revision,
		RevisionNumber: e.RevisionNumber,
		Branch:         e.Branch,
		UserId:         e.UserId.Clone(),
	}
}
func (e *LoreNotificationResourceLockedEventDataFFI) Clone() LoreNotificationResourceLockedEventData {
	return LoreNotificationResourceLockedEventData{
		UserId: e.UserId.Clone(),
		Branch: e.Branch,
		Paths:  e.Paths.Clone(),
	}
}
func (e *LoreNotificationResourceUnlockedEventDataFFI) Clone() LoreNotificationResourceUnlockedEventData {
	return LoreNotificationResourceUnlockedEventData{
		UserId: e.UserId.Clone(),
		Branch: e.Branch,
		Paths:  e.Paths.Clone(),
	}
}
func (e *LoreNotificationSubscribedEventDataFFI) Clone() LoreNotificationSubscribedEventData {
	return LoreNotificationSubscribedEventData{
		Repository: e.Repository,
	}
}
func (e *LoreNotificationUnsubscribedEventDataFFI) Clone() LoreNotificationUnsubscribedEventData {
	return LoreNotificationUnsubscribedEventData{
		Repository: e.Repository,
	}
}
func (e *LoreSharedStoreCreateEventDataFFI) Clone() LoreSharedStoreCreateEventData {
	return LoreSharedStoreCreateEventData{
		Path: e.Path.Clone(),
	}
}
func (e *LoreSharedStoreInfoEventDataFFI) Clone() LoreSharedStoreInfoEventData {
	return LoreSharedStoreInfoEventData{
		UseAutomatically: e.UseAutomatically != 0,
		RemoteUrls:       e.RemoteUrls.Clone(),
		Paths:            e.Paths.Clone(),
		Exists:           e.Exists.Clone(),
	}
}
func (e *LoreLinkStagedEntryEventDataFFI) Clone() LoreLinkStagedEntryEventData {
	return LoreLinkStagedEntryEventData{
		Path:            e.Path.Clone(),
		Repository:      e.Repository,
		StagedFileCount: e.StagedFileCount,
	}
}
func (e *LoreStorageOpenedEventDataFFI) Clone() LoreStorageOpenedEventData {
	return LoreStorageOpenedEventData{
		HandleId: e.HandleId,
	}
}
func (e *LoreStoragePutItemCompleteEventDataFFI) Clone() LoreStoragePutItemCompleteEventData {
	return LoreStoragePutItemCompleteEventData{
		Id:        e.Id,
		Address:   e.Address,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageGetHeaderEventDataFFI) Clone() LoreStorageGetHeaderEventData {
	return LoreStorageGetHeaderEventData{
		Id:          e.Id,
		Address:     e.Address,
		SizeContent: e.SizeContent,
	}
}
func (e *LoreStorageGetDataEventDataFFI) Clone() LoreStorageGetDataEventData {
	return LoreStorageGetDataEventData{
		Id:      e.Id,
		Address: e.Address,
		Offset:  e.Offset,
		Bytes:   e.Bytes.Clone(),
	}
}
func (e *LoreStorageGetItemCompleteEventDataFFI) Clone() LoreStorageGetItemCompleteEventData {
	return LoreStorageGetItemCompleteEventData{
		Id:        e.Id,
		Address:   e.Address,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageGetMetadataItemCompleteEventDataFFI) Clone() LoreStorageGetMetadataItemCompleteEventData {
	return LoreStorageGetMetadataItemCompleteEventData{
		Id:        e.Id,
		Address:   e.Address,
		Fragment:  e.Fragment.Clone(),
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageCopyItemCompleteEventDataFFI) Clone() LoreStorageCopyItemCompleteEventData {
	return LoreStorageCopyItemCompleteEventData{
		Id:              e.Id,
		SourcePartition: e.SourcePartition,
		TargetPartition: e.TargetPartition,
		SourceAddress:   e.SourceAddress,
		TargetContext:   e.TargetContext,
		ErrorCode:       e.ErrorCode,
	}
}
func (e *LoreStorageObliterateItemCompleteEventDataFFI) Clone() LoreStorageObliterateItemCompleteEventData {
	return LoreStorageObliterateItemCompleteEventData{
		Id:            e.Id,
		Address:       e.Address,
		LocalSuccess:  e.LocalSuccess != 0,
		RemoteSuccess: e.RemoteSuccess != 0,
		LocalSkipped:  e.LocalSkipped != 0,
		RemoteSkipped: e.RemoteSkipped != 0,
		ErrorCode:     e.ErrorCode,
	}
}
func (e *LoreStorageUploadItemCompleteEventDataFFI) Clone() LoreStorageUploadItemCompleteEventData {
	return LoreStorageUploadItemCompleteEventData{
		Id:             e.Id,
		Address:        e.Address,
		AlreadyDurable: e.AlreadyDurable != 0,
		ErrorCode:      e.ErrorCode,
	}
}
func (e *LoreRevisionTreeLoadedEventDataFFI) Clone() LoreRevisionTreeLoadedEventData {
	return LoreRevisionTreeLoadedEventData{
		HandleId: e.HandleId,
	}
}
func (e *LoreRevisionTreeResolvePathCompleteEventDataFFI) Clone() LoreRevisionTreeResolvePathCompleteEventData {
	return LoreRevisionTreeResolvePathCompleteEventData{
		Id:        e.Id,
		NodeId:    e.NodeId,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeChildEventDataFFI) Clone() LoreRevisionTreeChildEventData {
	return LoreRevisionTreeChildEventData{
		Id:        e.Id,
		NodeId:    e.NodeId,
		Name:      e.Name.Clone(),
		ParentId:  e.ParentId,
		Kind:      e.Kind,
		Mode:      e.Mode,
		Size:      e.Size,
		Address:   e.Address,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeNodeInfoEventDataFFI) Clone() LoreRevisionTreeNodeInfoEventData {
	return LoreRevisionTreeNodeInfoEventData{
		Id:       e.Id,
		NodeId:   e.NodeId,
		Name:     e.Name.Clone(),
		ParentId: e.ParentId,
		Kind:     e.Kind,
		Mode:     e.Mode,
		Size:     e.Size,
		Address:  e.Address,
		FileId:   e.FileId,
		RootInfo: e.RootInfo.Clone(),
	}
}
func (e *LoreRevisionTreeNodePathEventDataFFI) Clone() LoreRevisionTreeNodePathEventData {
	return LoreRevisionTreeNodePathEventData{
		Id:        e.Id,
		Path:      e.Path.Clone(),
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeAddCompleteEventDataFFI) Clone() LoreRevisionTreeAddCompleteEventData {
	return LoreRevisionTreeAddCompleteEventData{
		Id:        e.Id,
		NodeId:    e.NodeId,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeDeleteCompleteEventDataFFI) Clone() LoreRevisionTreeDeleteCompleteEventData {
	return LoreRevisionTreeDeleteCompleteEventData{
		Id:        e.Id,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeModifyCompleteEventDataFFI) Clone() LoreRevisionTreeModifyCompleteEventData {
	return LoreRevisionTreeModifyCompleteEventData{
		Id:        e.Id,
		NodeId:    e.NodeId,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeMoveCompleteEventDataFFI) Clone() LoreRevisionTreeMoveCompleteEventData {
	return LoreRevisionTreeMoveCompleteEventData{
		Id:        e.Id,
		NodeId:    e.NodeId,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeMetadataSetCompleteEventDataFFI) Clone() LoreRevisionTreeMetadataSetCompleteEventData {
	return LoreRevisionTreeMetadataSetCompleteEventData{
		Id:        e.Id,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeMetadataGetCompleteEventDataFFI) Clone() LoreRevisionTreeMetadataGetCompleteEventData {
	return LoreRevisionTreeMetadataGetCompleteEventData{
		Id:        e.Id,
		Key:       e.Key.Clone(),
		Value:     e.Value.Clone(),
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreRevisionTreeCommitCompleteEventDataFFI) Clone() LoreRevisionTreeCommitCompleteEventData {
	return LoreRevisionTreeCommitCompleteEventData{
		Id:           e.Id,
		RevisionHash: e.RevisionHash,
		NewTipHash:   e.NewTipHash,
		ErrorCode:    e.ErrorCode,
	}
}
func (e *LoreRevisionTreeCloseCompleteEventDataFFI) Clone() LoreRevisionTreeCloseCompleteEventData {
	return LoreRevisionTreeCloseCompleteEventData{
		Id:        e.Id,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageMutableLoadItemCompleteEventDataFFI) Clone() LoreStorageMutableLoadItemCompleteEventData {
	return LoreStorageMutableLoadItemCompleteEventData{
		Id:        e.Id,
		Value:     e.Value,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageMutableStoreItemCompleteEventDataFFI) Clone() LoreStorageMutableStoreItemCompleteEventData {
	return LoreStorageMutableStoreItemCompleteEventData{
		Id:        e.Id,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageMutableCompareAndSwapItemCompleteEventDataFFI) Clone() LoreStorageMutableCompareAndSwapItemCompleteEventData {
	return LoreStorageMutableCompareAndSwapItemCompleteEventData{
		Id:        e.Id,
		Previous:  e.Previous,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreStorageMutableListEntryEventDataFFI) Clone() LoreStorageMutableListEntryEventData {
	return LoreStorageMutableListEntryEventData{
		Id:    e.Id,
		Key:   e.Key,
		Value: e.Value,
	}
}
func (e *LoreStorageMutableListItemCompleteEventDataFFI) Clone() LoreStorageMutableListItemCompleteEventData {
	return LoreStorageMutableListItemCompleteEventData{
		Id:        e.Id,
		ErrorCode: e.ErrorCode,
	}
}
func (e *LoreEvictionBeginEventDataFFI) Clone() LoreEvictionBeginEventData {
	return LoreEvictionBeginEventData{
		TargetFragments: e.TargetFragments,
	}
}
func (e *LoreEvictionProgressEventDataFFI) Clone() LoreEvictionProgressEventData {
	return LoreEvictionProgressEventData{
		Evicted: e.Evicted,
	}
}
func (e *LoreEvictionEndEventDataFFI) Clone() LoreEvictionEndEventData {
	return LoreEvictionEndEventData{
		TotalEvicted: e.TotalEvicted,
	}
}
func (e *LoreCompactionBeginEventDataFFI) Clone() LoreCompactionBeginEventData {
	return LoreCompactionBeginEventData{
		TargetBytes: e.TargetBytes,
	}
}
func (e *LoreCompactionProgressEventDataFFI) Clone() LoreCompactionProgressEventData {
	return LoreCompactionProgressEventData{
		CompactedBytes: e.CompactedBytes,
	}
}
func (e *LoreCompactionEndEventDataFFI) Clone() LoreCompactionEndEventData {
	return LoreCompactionEndEventData{
		TotalCompactedBytes: e.TotalCompactedBytes,
	}
}

// Clone creates a Go-native copy of the event data
// This copy is safe to keep after the callback returns
func (e *LoreEventFFI) Clone() LoreEvent {
	switch e.Tag {
	case LoreEventTag_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_ERROR:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asErrorEventDataFFI().Clone(),
		}
	case LoreEventTag_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_METADATA:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asMetadataEventDataFFI().Clone(),
		}
	case LoreEventTag_LOG:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLogEventDataFFI().Clone(),
		}
	case LoreEventTag_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asEndEventDataFFI().Clone(),
		}
	case LoreEventTag_MAINTENANCE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asMaintenanceEventDataFFI().Clone(),
		}
	case LoreEventTag_AUTH_URL:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asAuthUrlEventDataFFI().Clone(),
		}
	case LoreEventTag_AUTH_USER_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asAuthUserInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_AUTH_USER_TOKEN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asAuthUserTokenEventDataFFI().Clone(),
		}
	case LoreEventTag_AUTH_IDENTITY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asAuthIdentityEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_CREATE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchCreateEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MULTIPLE_INSTANCE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMultipleInstanceEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_ARCHIVE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchArchiveEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_LIST_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchListBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_LIST_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchListEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_LIST_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchListEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_ABORT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeAbortBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_ABORT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeAbortEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CHANGE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffChangeBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CHANGE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffChangeEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CHANGE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffChangeEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CONFLICT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffConflictBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CONFLICT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffConflictEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_CONFLICT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffConflictEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_DIFF_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchDiffEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_LATEST_LIST_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchLatestListEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_CONFLICT_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeConflictFileEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_LINK_SKIPPED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeLinkSkippedEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_UNRESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeUnresolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_UNRESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeUnresolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FILE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFileBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFileEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FILE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFileEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFragmentBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFragmentProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoFragmentEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_SYNC_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoSyncBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_INTO_SYNC_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeIntoSyncEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_RESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeResolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_RESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeResolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_START_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeStartBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_MERGE_START_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchMergeStartEndEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_START_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickStartBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_START_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickStartEndEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_ABORT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickAbortBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_ABORT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickAbortEndEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_CONFLICT_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickConflictFileEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_UNRESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickUnresolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_UNRESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickUnresolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_RESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickResolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_CHERRY_PICK_RESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCherryPickResolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_START_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertStartBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_START_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertStartEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_ABORT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertAbortBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_ABORT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertAbortEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_RESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertResolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_RESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertResolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_CONFLICT_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertConflictFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_UNRESOLVE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertUnresolveFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVERT_UNRESOLVE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevertUnresolveRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PROTECT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchProtectEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushRevisionUpdateBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushRevisionUpdateEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushFragmentBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushFragmentProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_FRAGMENT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushFragmentEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushBranchCreateBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushBranchCreateEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushRevisionPushBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_UPDATE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushRevisionPushUpdateEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_PUSH_REVISION_PUSH_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchPushRevisionPushEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_RESET:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchResetEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_SWITCH_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchSwitchBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_SWITCH_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchSwitchEndEventDataFFI().Clone(),
		}
	case LoreEventTag_BRANCH_UNPROTECT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asBranchUnprotectEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DIFF:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDiffEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_HASH:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileHashEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_HISTORY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileHistoryEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_WRITE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileWriteEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_OBLITERATE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileObliterateEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DUMP:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDumpEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_ADD_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyAddBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_ADD_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyAddEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_ADD_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyAddEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyRemoveBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyRemoveEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_REMOVE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyRemoveEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_LIST_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyListBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_LIST_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyListFileEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_LIST_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyListEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_LIST_FILE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyListFileEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_DEPENDENCY_LIST_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileDependencyListEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_RESET_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileResetBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_RESET_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileResetProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_RESET_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileResetEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_RESET_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileResetFileEventDataFFI().Clone(),
		}
	case LoreEventTag_FILTER_EXCLUDE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFilterExcludeEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_STAGE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileStageBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_STAGE_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileStageProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_STAGE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileStageEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_STAGE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileStageRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_STAGE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileStageFileEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_UNSTAGE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileUnstageBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_UNSTAGE_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileUnstageProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_UNSTAGE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileUnstageEndEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_UNSTAGE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileUnstageRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_FILE_UNSTAGE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFileUnstageFileEventDataFFI().Clone(),
		}
	case LoreEventTag_FRAGMENT_WRITE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asFragmentWriteEventDataFFI().Clone(),
		}
	case LoreEventTag_LAYER_ADD:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLayerAddEventDataFFI().Clone(),
		}
	case LoreEventTag_LAYER_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLayerEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_LAYER_REMOVE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLayerRemoveEventDataFFI().Clone(),
		}
	case LoreEventTag_LAYER_STAGED_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLayerStagedEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_LINK_CHANGE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLinkChangeEventDataFFI().Clone(),
		}
	case LoreEventTag_LINK_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLinkEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_ACQUIRE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileAcquireBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_ACQUIRE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileAcquireEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_STATUS_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileStatusBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_STATUS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileStatusEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_QUERY_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileQueryBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_QUERY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileQueryEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_RELEASE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileReleaseBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_LOCK_FILE_RELEASE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLockFileReleaseEventDataFFI().Clone(),
		}
	case LoreEventTag_METADATA_CLEAR_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asMetadataClearFileEventDataFFI().Clone(),
		}
	case LoreEventTag_METADATA_CLEAR_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asMetadataClearRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_PATH_IGNORE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asPathIgnoreEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_CREATE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryCreateEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_CLONE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryCloneBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_CLONE_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryCloneProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_CLONE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryCloneEndEventDataFFI().Clone(),
		}
	case LoreEventTag_DEPENDENCY_RESOLVE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asDependencyResolveBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_DEPENDENCY_RESOLVE_ITEM:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asDependencyResolveItemEventDataFFI().Clone(),
		}
	case LoreEventTag_DEPENDENCY_RESOLVE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asDependencyResolveEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_DATA:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryDataEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_CONFIG_GET:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryConfigGetEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_DUMP_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryDumpBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_DUMP_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryDumpEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_LIST_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryListEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_INSTANCE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryInstanceEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_VERIFY_STATE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryVerifyStateBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_VERIFY_STATE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryVerifyStateEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryVerifyFragmentEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_MATCH:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryVerifyFragmentMatchEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_REMOTE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryVerifyFragmentRemoteEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATE_DUMP:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStateDumpEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATE_DUMP_NODE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStateDumpNodeEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATUS_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStatusRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATUS_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStatusFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATUS_COUNT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStatusCountEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STATUS_SUMMARY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStatusSummaryEventDataFFI().Clone(),
		}
	case LoreEventTag_REPOSITORY_STORE_IMMUTABLE_QUERY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRepositoryStoreImmutableQueryEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_COMMIT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionCommitBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_COMMIT_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionCommitProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_COMMIT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionCommitEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_COMMIT_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionCommitRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_INFO_DELTA:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionInfoDeltaEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_DIFF_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionDiffFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_FIND:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionFindEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_HISTORY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionHistoryEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_HISTORY_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionHistoryEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FILE_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFileBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FILE_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFileEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFragmentBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFragmentProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_FRAGMENT_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreFragmentEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_SYNC_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreSyncBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESTORE_SYNC_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionRestoreSyncEndEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_RESOLVE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionResolveEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_SYNC_TARGET:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionSyncTargetEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_SYNC_FILE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionSyncFileEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_SYNC_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionSyncProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_SYNC_REVISION:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionSyncRevisionEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_BISECT:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionBisectEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_BRANCH_CREATED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationBranchCreatedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_BRANCH_DELETED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationBranchDeletedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_BRANCH_PUSHED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationBranchPushedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_RESOURCE_LOCKED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationResourceLockedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_RESOURCE_UNLOCKED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationResourceUnlockedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_SUBSCRIBED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationSubscribedEventDataFFI().Clone(),
		}
	case LoreEventTag_NOTIFICATION_UNSUBSCRIBED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asNotificationUnsubscribedEventDataFFI().Clone(),
		}
	case LoreEventTag_SHARED_STORE_CREATE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asSharedStoreCreateEventDataFFI().Clone(),
		}
	case LoreEventTag_SHARED_STORE_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asSharedStoreInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_LINK_STAGED_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asLinkStagedEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_OPENED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageOpenedEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_PUT_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStoragePutItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_GET_HEADER:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageGetHeaderEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_GET_DATA:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageGetDataEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_GET_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageGetItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_GET_METADATA_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageGetMetadataItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_COPY_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageCopyItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_OBLITERATE_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageObliterateItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_UPLOAD_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageUploadItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_LOADED:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeLoadedEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_RESOLVE_PATH_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeResolvePathCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_CHILD:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeChildEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_NODE_INFO:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeNodeInfoEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_NODE_PATH:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeNodePathEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_ADD_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeAddCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_DELETE_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeDeleteCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_MODIFY_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeModifyCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_MOVE_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeMoveCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_METADATA_SET_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeMetadataSetCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_METADATA_GET_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeMetadataGetCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_COMMIT_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeCommitCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_REVISION_TREE_CLOSE_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asRevisionTreeCloseCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_MUTABLE_LOAD_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageMutableLoadItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_MUTABLE_STORE_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageMutableStoreItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageMutableCompareAndSwapItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_MUTABLE_LIST_ENTRY:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageMutableListEntryEventDataFFI().Clone(),
		}
	case LoreEventTag_STORAGE_MUTABLE_LIST_ITEM_COMPLETE:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asStorageMutableListItemCompleteEventDataFFI().Clone(),
		}
	case LoreEventTag_EVICTION_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asEvictionBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_EVICTION_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asEvictionProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_EVICTION_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asEvictionEndEventDataFFI().Clone(),
		}
	case LoreEventTag_COMPACTION_BEGIN:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCompactionBeginEventDataFFI().Clone(),
		}
	case LoreEventTag_COMPACTION_PROGRESS:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCompactionProgressEventDataFFI().Clone(),
		}
	case LoreEventTag_COMPACTION_END:
		return LoreEvent{
			Tag:  e.Tag,
			Data: e.asCompactionEndEventDataFFI().Clone(),
		}
	default:
		return LoreEvent{
			Tag: e.Tag,
		}
	}
}
