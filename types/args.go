// Copyright Epic Games, Inc. All Rights Reserved.

package types

type LoreGlobalArgs struct {
	/* Repository path */
	RepositoryPath string
	/* Correlation ID */
	CorrelationId string
	/* Identity to use */
	Identity string
	/* Force the operation if possible */
	Force bool
	/* Run operation without connecting to server */
	Offline bool
	/* Use only local data */
	Local bool
	/* Use only remote data */
	Remote bool
	/* Dry run mode, only report what would have been changed and perform no changes to local file system */
	DryRun bool
	/* Avoid recording last access timestamps in the data stores */
	NoAtime bool
	/* Maximum number of parallel connections for bulk data transfer */
	MaxConnections uint32
	/* Search limit when iterating revisions */
	SearchLimit uint32
	/* Allow matching to the nearest matching revision when a perfect match is not available */
	SearchNearest bool
	/* Prevent the automatic incremental/step GC for this operation; it otherwise runs in the background on write operations. `repository gc` always runs a full pass regardless */
	NoGc bool
	/* Use in-memory stores instead of file-backed stores. No store data is
	read from or written to the .urc/immutable/ and .urc/mutable/ directories. */
	InMemory bool
	/* Maximum number of files being processed in parallel */
	FileCountLimit uint64
	/* Maximum total size of all files being processed in parallel */
	FileSizeLimit uint64
	/* Maximum number of parallel compression tasks */
	CompressTaskLimit uint64
	/* Keep store references alive after a repository call completes to avoid
	repeated store open/close cycles for consecutive API calls in the same process. */
	StoreKeepAlive bool
	/* Duration in seconds to keep store references alive. Only used when
	`store_keep_alive` is set. 0 means use the default (10 seconds). */
	StoreKeepAliveSeconds uint64
	/* Force sync data to storage media during store flush */
	SyncData bool
	/* Cache fragment payloads fetched from remote in the local store. Without
	this only state fragments and fragments flagged for local cache priority
	are retained */
	Cache bool
}

type LoreGlobalArgsFFI struct {
	/* Repository path */
	RepositoryPath LoreString
	/* Correlation ID */
	CorrelationId LoreString
	/* Identity to use */
	Identity LoreString
	/* Force the operation if possible */
	Force uint8
	/* Run operation without connecting to server */
	Offline uint8
	/* Use only local data */
	Local uint8
	/* Use only remote data */
	Remote uint8
	/* Dry run mode, only report what would have been changed and perform no changes to local file system */
	DryRun uint8
	/* Avoid recording last access timestamps in the data stores */
	NoAtime uint8
	/* Maximum number of parallel connections for bulk data transfer */
	MaxConnections uint32
	/* Search limit when iterating revisions */
	SearchLimit uint32
	/* Allow matching to the nearest matching revision when a perfect match is not available */
	SearchNearest uint8
	/* Prevent the automatic incremental/step GC for this operation; it otherwise runs in the background on write operations. `repository gc` always runs a full pass regardless */
	NoGc uint8
	/* Use in-memory stores instead of file-backed stores. No store data is
	read from or written to the .urc/immutable/ and .urc/mutable/ directories. */
	InMemory uint8
	/* Maximum number of files being processed in parallel */
	FileCountLimit uint64
	/* Maximum total size of all files being processed in parallel */
	FileSizeLimit uint64
	/* Maximum number of parallel compression tasks */
	CompressTaskLimit uint64
	/* Keep store references alive after a repository call completes to avoid
	repeated store open/close cycles for consecutive API calls in the same process. */
	StoreKeepAlive uint8
	/* Duration in seconds to keep store references alive. Only used when
	`store_keep_alive` is set. 0 means use the default (10 seconds). */
	StoreKeepAliveSeconds uint64
	/* Force sync data to storage media during store flush */
	SyncData uint8
	/* Cache fragment payloads fetched from remote in the local store. Without
	this only state fragments and fragments flagged for local cache priority
	are retained */
	Cache uint8
}

func NewLoreGlobalArgs(opts LoreGlobalArgs) (LoreGlobalArgsFFI, func()) {
	valRepositoryPath, cleanupRepositoryPath := NewLoreString(opts.RepositoryPath)
	valCorrelationId, cleanupCorrelationId := NewLoreString(opts.CorrelationId)
	valIdentity, cleanupIdentity := NewLoreString(opts.Identity)
	valForce, cleanupForce := Newuint8(opts.Force)
	valOffline, cleanupOffline := Newuint8(opts.Offline)
	valLocal, cleanupLocal := Newuint8(opts.Local)
	valRemote, cleanupRemote := Newuint8(opts.Remote)
	valDryRun, cleanupDryRun := Newuint8(opts.DryRun)
	valNoAtime, cleanupNoAtime := Newuint8(opts.NoAtime)
	valSearchNearest, cleanupSearchNearest := Newuint8(opts.SearchNearest)
	valNoGc, cleanupNoGc := Newuint8(opts.NoGc)
	valInMemory, cleanupInMemory := Newuint8(opts.InMemory)
	valStoreKeepAlive, cleanupStoreKeepAlive := Newuint8(opts.StoreKeepAlive)
	valSyncData, cleanupSyncData := Newuint8(opts.SyncData)
	valCache, cleanupCache := Newuint8(opts.Cache)

	cleanup := func() {
		cleanupRepositoryPath()
		cleanupCorrelationId()
		cleanupIdentity()
		cleanupForce()
		cleanupOffline()
		cleanupLocal()
		cleanupRemote()
		cleanupDryRun()
		cleanupNoAtime()
		cleanupSearchNearest()
		cleanupNoGc()
		cleanupInMemory()
		cleanupStoreKeepAlive()
		cleanupSyncData()
		cleanupCache()
	}

	return LoreGlobalArgsFFI{
		RepositoryPath:        valRepositoryPath,
		CorrelationId:         valCorrelationId,
		Identity:              valIdentity,
		Force:                 valForce,
		Offline:               valOffline,
		Local:                 valLocal,
		Remote:                valRemote,
		DryRun:                valDryRun,
		NoAtime:               valNoAtime,
		MaxConnections:        opts.MaxConnections,
		SearchLimit:           opts.SearchLimit,
		SearchNearest:         valSearchNearest,
		NoGc:                  valNoGc,
		InMemory:              valInMemory,
		FileCountLimit:        opts.FileCountLimit,
		FileSizeLimit:         opts.FileSizeLimit,
		CompressTaskLimit:     opts.CompressTaskLimit,
		StoreKeepAlive:        valStoreKeepAlive,
		StoreKeepAliveSeconds: opts.StoreKeepAliveSeconds,
		SyncData:              valSyncData,
		Cache:                 valCache,
	}, cleanup
}

type LoreAuthUserInfoArgs struct {
	/* User IDs to resolve; empty resolves the current user locally */
	UserIds []string
}

type LoreAuthUserInfoArgsFFI struct {
	/* User IDs to resolve; empty resolves the current user locally */
	UserIds LoreStringArrayFFI
}

func NewLoreAuthUserInfoArgs(opts LoreAuthUserInfoArgs) (LoreAuthUserInfoArgsFFI, func()) {
	valUserIds, cleanupUserIds := NewLoreStringArray(opts.UserIds)

	cleanup := func() {
		cleanupUserIds()
	}

	return LoreAuthUserInfoArgsFFI{
		UserIds: valUserIds,
	}, cleanup
}

type LoreAuthLoginWithTokenArgs struct {
	/* Remote URL; empty resolves from the repository config */
	RemoteUrl string
	/* Authentication token */
	Token string
	/* Token type */
	TokenType string
	/* Auth service URL with scheme (e.g. `ucs-auth://auth.example.com`); used
	directly when non-empty, required when no remote URL is available */
	AuthUrl string
}

type LoreAuthLoginWithTokenArgsFFI struct {
	/* Remote URL; empty resolves from the repository config */
	RemoteUrl LoreString
	/* Authentication token */
	Token LoreString
	/* Token type */
	TokenType LoreString
	/* Auth service URL with scheme (e.g. `ucs-auth://auth.example.com`); used
	directly when non-empty, required when no remote URL is available */
	AuthUrl LoreString
}

func NewLoreAuthLoginWithTokenArgs(opts LoreAuthLoginWithTokenArgs) (LoreAuthLoginWithTokenArgsFFI, func()) {
	valRemoteUrl, cleanupRemoteUrl := NewLoreString(opts.RemoteUrl)
	valToken, cleanupToken := NewLoreString(opts.Token)
	valTokenType, cleanupTokenType := NewLoreString(opts.TokenType)
	valAuthUrl, cleanupAuthUrl := NewLoreString(opts.AuthUrl)

	cleanup := func() {
		cleanupRemoteUrl()
		cleanupToken()
		cleanupTokenType()
		cleanupAuthUrl()
	}

	return LoreAuthLoginWithTokenArgsFFI{
		RemoteUrl: valRemoteUrl,
		Token:     valToken,
		TokenType: valTokenType,
		AuthUrl:   valAuthUrl,
	}, cleanup
}

type LoreAuthListArgs struct {
	/* Include the decrypted cached token in each identity */
	WithToken bool
}

type LoreAuthListArgsFFI struct {
	/* Include the decrypted cached token in each identity */
	WithToken uint8
}

func NewLoreAuthListArgs(opts LoreAuthListArgs) (LoreAuthListArgsFFI, func()) {
	valWithToken, cleanupWithToken := Newuint8(opts.WithToken)

	cleanup := func() {
		cleanupWithToken()
	}

	return LoreAuthListArgsFFI{
		WithToken: valWithToken,
	}, cleanup
}

type LoreAuthLogoutArgs struct {
	/* Auth service URL; empty resolves from the repository */
	AuthUrl string
	/* Resource ID (e.g. `urc-{id}`); empty removes all tokens for the auth URL */
	Resource string
	/* User identity to remove; empty removes all identities */
	UserId string
}

type LoreAuthLogoutArgsFFI struct {
	/* Auth service URL; empty resolves from the repository */
	AuthUrl LoreString
	/* Resource ID (e.g. `urc-{id}`); empty removes all tokens for the auth URL */
	Resource LoreString
	/* User identity to remove; empty removes all identities */
	UserId LoreString
}

func NewLoreAuthLogoutArgs(opts LoreAuthLogoutArgs) (LoreAuthLogoutArgsFFI, func()) {
	valAuthUrl, cleanupAuthUrl := NewLoreString(opts.AuthUrl)
	valResource, cleanupResource := NewLoreString(opts.Resource)
	valUserId, cleanupUserId := NewLoreString(opts.UserId)

	cleanup := func() {
		cleanupAuthUrl()
		cleanupResource()
		cleanupUserId()
	}

	return LoreAuthLogoutArgsFFI{
		AuthUrl:  valAuthUrl,
		Resource: valResource,
		UserId:   valUserId,
	}, cleanup
}

type LoreAuthClearArgs struct {
	Unused bool
}

type LoreAuthClearArgsFFI struct {
	Unused uint8
}

func NewLoreAuthClearArgs(opts LoreAuthClearArgs) (LoreAuthClearArgsFFI, func()) {
	valUnused, cleanupUnused := Newuint8(opts.Unused)

	cleanup := func() {
		cleanupUnused()
	}

	return LoreAuthClearArgsFFI{
		Unused: valUnused,
	}, cleanup
}

type LoreAuthLocalUserInfoArgs struct {
	/* Auth service remote URL; empty resolves from the repository's remote environment */
	AuthEndpoint string
	/* User identities to resolve; empty resolves the current user */
	UserIds []string
	/* Emit cached token details for identities with a local token */
	WithToken bool
}

type LoreAuthLocalUserInfoArgsFFI struct {
	/* Auth service remote URL; empty resolves from the repository's remote environment */
	AuthEndpoint LoreString
	/* User identities to resolve; empty resolves the current user */
	UserIds LoreStringArrayFFI
	/* Emit cached token details for identities with a local token */
	WithToken uint8
}

func NewLoreAuthLocalUserInfoArgs(opts LoreAuthLocalUserInfoArgs) (LoreAuthLocalUserInfoArgsFFI, func()) {
	valAuthEndpoint, cleanupAuthEndpoint := NewLoreString(opts.AuthEndpoint)
	valUserIds, cleanupUserIds := NewLoreStringArray(opts.UserIds)
	valWithToken, cleanupWithToken := Newuint8(opts.WithToken)

	cleanup := func() {
		cleanupAuthEndpoint()
		cleanupUserIds()
		cleanupWithToken()
	}

	return LoreAuthLocalUserInfoArgsFFI{
		AuthEndpoint: valAuthEndpoint,
		UserIds:      valUserIds,
		WithToken:    valWithToken,
	}, cleanup
}

type LoreAuthLoginInteractiveArgs struct {
	/* Remote URL; empty resolves from the repository config */
	RemoteUrl string
	/* Emit the login URL instead of opening a browser */
	NoBrowser bool
}

type LoreAuthLoginInteractiveArgsFFI struct {
	/* Remote URL; empty resolves from the repository config */
	RemoteUrl LoreString
	/* Emit the login URL instead of opening a browser */
	NoBrowser uint8
}

func NewLoreAuthLoginInteractiveArgs(opts LoreAuthLoginInteractiveArgs) (LoreAuthLoginInteractiveArgsFFI, func()) {
	valRemoteUrl, cleanupRemoteUrl := NewLoreString(opts.RemoteUrl)
	valNoBrowser, cleanupNoBrowser := Newuint8(opts.NoBrowser)

	cleanup := func() {
		cleanupRemoteUrl()
		cleanupNoBrowser()
	}

	return LoreAuthLoginInteractiveArgsFFI{
		RemoteUrl: valRemoteUrl,
		NoBrowser: valNoBrowser,
	}, cleanup
}

type LoreBranchCreateArgs struct {
	/* Name of the branch */
	Branch string
	/* Category of the branch */
	Category string
	/* Optional explicit branch ID (hex-encoded 16-byte context) */
	Id string
}

type LoreBranchCreateArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
	/* Category of the branch */
	Category LoreString
	/* Optional explicit branch ID (hex-encoded 16-byte context) */
	Id LoreString
}

func NewLoreBranchCreateArgs(opts LoreBranchCreateArgs) (LoreBranchCreateArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valCategory, cleanupCategory := NewLoreString(opts.Category)
	valId, cleanupId := NewLoreString(opts.Id)

	cleanup := func() {
		cleanupBranch()
		cleanupCategory()
		cleanupId()
	}

	return LoreBranchCreateArgsFFI{
		Branch:   valBranch,
		Category: valCategory,
		Id:       valId,
	}, cleanup
}

type LoreBranchInfoArgs struct {
	/* Name of the branch */
	Branch string
}

type LoreBranchInfoArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
}

func NewLoreBranchInfoArgs(opts LoreBranchInfoArgs) (LoreBranchInfoArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupBranch()
	}

	return LoreBranchInfoArgsFFI{
		Branch: valBranch,
	}, cleanup
}

type LoreBranchDiffArgs struct {
	/* Source branch name */
	Source string
	/* Target branch name */
	Target string
	/* Optional path in the repository to limit the diff to */
	Path string
	/* Attempt to auto resolve conflicts */
	AutoResolve bool
}

type LoreBranchDiffArgsFFI struct {
	/* Source branch name */
	Source LoreString
	/* Target branch name */
	Target LoreString
	/* Optional path in the repository to limit the diff to */
	Path LoreString
	/* Attempt to auto resolve conflicts */
	AutoResolve uint8
}

func NewLoreBranchDiffArgs(opts LoreBranchDiffArgs) (LoreBranchDiffArgsFFI, func()) {
	valSource, cleanupSource := NewLoreString(opts.Source)
	valTarget, cleanupTarget := NewLoreString(opts.Target)
	valPath, cleanupPath := NewLoreString(opts.Path)
	valAutoResolve, cleanupAutoResolve := Newuint8(opts.AutoResolve)

	cleanup := func() {
		cleanupSource()
		cleanupTarget()
		cleanupPath()
		cleanupAutoResolve()
	}

	return LoreBranchDiffArgsFFI{
		Source:      valSource,
		Target:      valTarget,
		Path:        valPath,
		AutoResolve: valAutoResolve,
	}, cleanup
}

type LoreBranchProtectArgs struct {
	/* Name of the branch */
	Branch string
}

type LoreBranchProtectArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
}

func NewLoreBranchProtectArgs(opts LoreBranchProtectArgs) (LoreBranchProtectArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupBranch()
	}

	return LoreBranchProtectArgsFFI{
		Branch: valBranch,
	}, cleanup
}

type LoreBranchUnprotectArgs struct {
	/* Name of the branch */
	Branch string
}

type LoreBranchUnprotectArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
}

func NewLoreBranchUnprotectArgs(opts LoreBranchUnprotectArgs) (LoreBranchUnprotectArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupBranch()
	}

	return LoreBranchUnprotectArgsFFI{
		Branch: valBranch,
	}, cleanup
}

type LoreBranchArchiveArgs struct {
	/* Name of the branch */
	Branch string
}

type LoreBranchArchiveArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
}

func NewLoreBranchArchiveArgs(opts LoreBranchArchiveArgs) (LoreBranchArchiveArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupBranch()
	}

	return LoreBranchArchiveArgsFFI{
		Branch: valBranch,
	}, cleanup
}

type LoreBranchListArgs struct {
	/* Include archived local branches in listing */
	Archived bool
}

type LoreBranchListArgsFFI struct {
	/* Include archived local branches in listing */
	Archived uint8
}

func NewLoreBranchListArgs(opts LoreBranchListArgs) (LoreBranchListArgsFFI, func()) {
	valArchived, cleanupArchived := Newuint8(opts.Archived)

	cleanup := func() {
		cleanupArchived()
	}

	return LoreBranchListArgsFFI{
		Archived: valArchived,
	}, cleanup
}

type LoreBranchMergeAbortArgs struct {
	/* Optional link path for link-scoped abort */
	Link string
	/* Abort only the main repository merge, keeping link pin updates */
	IgnoreLinks bool
}

type LoreBranchMergeAbortArgsFFI struct {
	/* Optional link path for link-scoped abort */
	Link LoreString
	/* Abort only the main repository merge, keeping link pin updates */
	IgnoreLinks uint8
}

func NewLoreBranchMergeAbortArgs(opts LoreBranchMergeAbortArgs) (LoreBranchMergeAbortArgsFFI, func()) {
	valLink, cleanupLink := NewLoreString(opts.Link)
	valIgnoreLinks, cleanupIgnoreLinks := Newuint8(opts.IgnoreLinks)

	cleanup := func() {
		cleanupLink()
		cleanupIgnoreLinks()
	}

	return LoreBranchMergeAbortArgsFFI{
		Link:        valLink,
		IgnoreLinks: valIgnoreLinks,
	}, cleanup
}

type LoreBranchMergeUnresolveArgs struct {
	/* Paths to mark unresolved */
	Paths []string
}

type LoreBranchMergeUnresolveArgsFFI struct {
	/* Paths to mark unresolved */
	Paths LoreStringArrayFFI
}

func NewLoreBranchMergeUnresolveArgs(opts LoreBranchMergeUnresolveArgs) (LoreBranchMergeUnresolveArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreBranchMergeUnresolveArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreBranchMergeIntoArgs struct {
	/* Name of the target branch to merge into */
	Branch string
	/* ID of the target branch to merge into */
	BranchId LoreContext
	/* Commit message for the auto-commit */
	Message string
	/* Optional link path for link-scoped merge into */
	Link string
	/* Merge only the main repository, skipping all linked repositories */
	IgnoreLinks bool
}

type LoreBranchMergeIntoArgsFFI struct {
	/* Name of the target branch to merge into */
	Branch LoreString
	/* ID of the target branch to merge into */
	BranchId LoreContext
	/* Commit message for the auto-commit */
	Message LoreString
	/* Optional link path for link-scoped merge into */
	Link LoreString
	/* Merge only the main repository, skipping all linked repositories */
	IgnoreLinks uint8
}

func NewLoreBranchMergeIntoArgs(opts LoreBranchMergeIntoArgs) (LoreBranchMergeIntoArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valMessage, cleanupMessage := NewLoreString(opts.Message)
	valLink, cleanupLink := NewLoreString(opts.Link)
	valIgnoreLinks, cleanupIgnoreLinks := Newuint8(opts.IgnoreLinks)

	cleanup := func() {
		cleanupBranch()
		cleanupMessage()
		cleanupLink()
		cleanupIgnoreLinks()
	}

	return LoreBranchMergeIntoArgsFFI{
		Branch:      valBranch,
		BranchId:    opts.BranchId,
		Message:     valMessage,
		Link:        valLink,
		IgnoreLinks: valIgnoreLinks,
	}, cleanup
}

type LoreBranchMergeResolveArgs struct {
	/* Paths to mark resolved */
	Paths []string
}

type LoreBranchMergeResolveArgsFFI struct {
	/* Paths to mark resolved */
	Paths LoreStringArrayFFI
}

func NewLoreBranchMergeResolveArgs(opts LoreBranchMergeResolveArgs) (LoreBranchMergeResolveArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreBranchMergeResolveArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreBranchMergeResolveMineArgs struct {
	/* Paths to resolve as "mine" */
	Paths []string
}

type LoreBranchMergeResolveMineArgsFFI struct {
	/* Paths to resolve as "mine" */
	Paths LoreStringArrayFFI
}

func NewLoreBranchMergeResolveMineArgs(opts LoreBranchMergeResolveMineArgs) (LoreBranchMergeResolveMineArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreBranchMergeResolveMineArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreBranchMergeResolveTheirsArgs struct {
	/* Paths to resolve as "theirs" */
	Paths []string
}

type LoreBranchMergeResolveTheirsArgsFFI struct {
	/* Paths to resolve as "theirs" */
	Paths LoreStringArrayFFI
}

func NewLoreBranchMergeResolveTheirsArgs(opts LoreBranchMergeResolveTheirsArgs) (LoreBranchMergeResolveTheirsArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreBranchMergeResolveTheirsArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreBranchMergeRestartArgs struct {
	/* Paths to re-materialize */
	Paths []string
}

type LoreBranchMergeRestartArgsFFI struct {
	/* Paths to re-materialize */
	Paths LoreStringArrayFFI
}

func NewLoreBranchMergeRestartArgs(opts LoreBranchMergeRestartArgs) (LoreBranchMergeRestartArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreBranchMergeRestartArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreBranchMergeStartArgs struct {
	/* Name of the source branch to merge into the current branch */
	Branch string
	/* Message to use for an auto commit if no conflicts arise */
	Message string
	/* Disable auto commit even if no conflicts arise */
	NoCommit bool
	/* Optional link path for link-scoped merge */
	Link string
	/* Merge only the main repository, skipping all linked repositories */
	IgnoreLinks bool
}

type LoreBranchMergeStartArgsFFI struct {
	/* Name of the source branch to merge into the current branch */
	Branch LoreString
	/* Message to use for an auto commit if no conflicts arise */
	Message LoreString
	/* Disable auto commit even if no conflicts arise */
	NoCommit uint8
	/* Optional link path for link-scoped merge */
	Link LoreString
	/* Merge only the main repository, skipping all linked repositories */
	IgnoreLinks uint8
}

func NewLoreBranchMergeStartArgs(opts LoreBranchMergeStartArgs) (LoreBranchMergeStartArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valMessage, cleanupMessage := NewLoreString(opts.Message)
	valNoCommit, cleanupNoCommit := Newuint8(opts.NoCommit)
	valLink, cleanupLink := NewLoreString(opts.Link)
	valIgnoreLinks, cleanupIgnoreLinks := Newuint8(opts.IgnoreLinks)

	cleanup := func() {
		cleanupBranch()
		cleanupMessage()
		cleanupNoCommit()
		cleanupLink()
		cleanupIgnoreLinks()
	}

	return LoreBranchMergeStartArgsFFI{
		Branch:      valBranch,
		Message:     valMessage,
		NoCommit:    valNoCommit,
		Link:        valLink,
		IgnoreLinks: valIgnoreLinks,
	}, cleanup
}

type LoreBranchSwitchArgs struct {
	/* Name of the branch */
	Branch string
	/* Hash of the revision */
	Revision string
	/* Reset local modified files to match the incoming revision */
	Reset bool
	/* Only update anchor tracking without modifying or verifying files */
	Bare bool
}

type LoreBranchSwitchArgsFFI struct {
	/* Name of the branch */
	Branch LoreString
	/* Hash of the revision */
	Revision LoreString
	/* Reset local modified files to match the incoming revision */
	Reset uint8
	/* Only update anchor tracking without modifying or verifying files */
	Bare uint8
}

func NewLoreBranchSwitchArgs(opts LoreBranchSwitchArgs) (LoreBranchSwitchArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valReset, cleanupReset := Newuint8(opts.Reset)
	valBare, cleanupBare := Newuint8(opts.Bare)

	cleanup := func() {
		cleanupBranch()
		cleanupRevision()
		cleanupReset()
		cleanupBare()
	}

	return LoreBranchSwitchArgsFFI{
		Branch:   valBranch,
		Revision: valRevision,
		Reset:    valReset,
		Bare:     valBare,
	}, cleanup
}

type LoreBranchResetArgs struct {
	/* Revision to reset the local LATEST pointer to */
	Revision string
	/* Branch to reset, current branch if empty */
	Branch string
}

type LoreBranchResetArgsFFI struct {
	/* Revision to reset the local LATEST pointer to */
	Revision LoreString
	/* Branch to reset, current branch if empty */
	Branch LoreString
}

func NewLoreBranchResetArgs(opts LoreBranchResetArgs) (LoreBranchResetArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupRevision()
		cleanupBranch()
	}

	return LoreBranchResetArgsFFI{
		Revision: valRevision,
		Branch:   valBranch,
	}, cleanup
}

type LoreBranchPushArgs struct {
	/* Optional branch to push, current branch if not given */
	Branch string
	/* Allow the server to fast-forward merge if the target branch head has moved */
	FastForwardMerge bool
}

type LoreBranchPushArgsFFI struct {
	/* Optional branch to push, current branch if not given */
	Branch LoreString
	/* Allow the server to fast-forward merge if the target branch head has moved */
	FastForwardMerge uint8
}

func NewLoreBranchPushArgs(opts LoreBranchPushArgs) (LoreBranchPushArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valFastForwardMerge, cleanupFastForwardMerge := Newuint8(opts.FastForwardMerge)

	cleanup := func() {
		cleanupBranch()
		cleanupFastForwardMerge()
	}

	return LoreBranchPushArgsFFI{
		Branch:           valBranch,
		FastForwardMerge: valFastForwardMerge,
	}, cleanup
}

type LoreBranchMetadataGetArgs struct {
	/* Branch name or identifier */
	Branch string
	/* Metadata key (empty string lists all) */
	Key string
}

type LoreBranchMetadataGetArgsFFI struct {
	/* Branch name or identifier */
	Branch LoreString
	/* Metadata key (empty string lists all) */
	Key LoreString
}

func NewLoreBranchMetadataGetArgs(opts LoreBranchMetadataGetArgs) (LoreBranchMetadataGetArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupBranch()
		cleanupKey()
	}

	return LoreBranchMetadataGetArgsFFI{
		Branch: valBranch,
		Key:    valKey,
	}, cleanup
}

type LoreBranchMetadataSetArgs struct {
	/* Branch name or identifier */
	Branch string
	/* Metadata keys to set (parallel with `values`/`formats`) */
	Keys []string
	/* Values to set, one per key (decoded per the matching `formats` entry) */
	Values []string
	/* Value type for each key, one per key */
	Formats LoreMetadataTypeArray
}

type LoreBranchMetadataSetArgsFFI struct {
	/* Branch name or identifier */
	Branch LoreString
	/* Metadata keys to set (parallel with `values`/`formats`) */
	Keys LoreStringArrayFFI
	/* Values to set, one per key (decoded per the matching `formats` entry) */
	Values LoreStringArrayFFI
	/* Value type for each key, one per key */
	Formats LoreMetadataTypeArrayFFI
}

func NewLoreBranchMetadataSetArgs(opts LoreBranchMetadataSetArgs) (LoreBranchMetadataSetArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)
	valValues, cleanupValues := NewLoreStringArray(opts.Values)
	valFormats, cleanupFormats := NewLoreMetadataTypeArray(opts.Formats)

	cleanup := func() {
		cleanupBranch()
		cleanupKeys()
		cleanupValues()
		cleanupFormats()
	}

	return LoreBranchMetadataSetArgsFFI{
		Branch:  valBranch,
		Keys:    valKeys,
		Values:  valValues,
		Formats: valFormats,
	}, cleanup
}

type LoreBranchMetadataClearArgs struct {
	/* Branch name or identifier */
	Branch string
	/* Keys to clear (empty array clears all user-defined keys) */
	Keys []string
}

type LoreBranchMetadataClearArgsFFI struct {
	/* Branch name or identifier */
	Branch LoreString
	/* Keys to clear (empty array clears all user-defined keys) */
	Keys LoreStringArrayFFI
}

func NewLoreBranchMetadataClearArgs(opts LoreBranchMetadataClearArgs) (LoreBranchMetadataClearArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)

	cleanup := func() {
		cleanupBranch()
		cleanupKeys()
	}

	return LoreBranchMetadataClearArgsFFI{
		Branch: valBranch,
		Keys:   valKeys,
	}, cleanup
}

type LoreFileInfoArgs struct {
	/* Array of paths */
	Paths []string
	/* Revision to get info for */
	Revision string
	/* Calculate the filtered local filesystem hash and size */
	Local bool
	/* Calculate the filtered repository size */
	Filtered bool
}

type LoreFileInfoArgsFFI struct {
	/* Array of paths */
	Paths LoreStringArrayFFI
	/* Revision to get info for */
	Revision LoreString
	/* Calculate the filtered local filesystem hash and size */
	Local uint8
	/* Calculate the filtered repository size */
	Filtered uint8
}

func NewLoreFileInfoArgs(opts LoreFileInfoArgs) (LoreFileInfoArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valLocal, cleanupLocal := Newuint8(opts.Local)
	valFiltered, cleanupFiltered := Newuint8(opts.Filtered)

	cleanup := func() {
		cleanupPaths()
		cleanupRevision()
		cleanupLocal()
		cleanupFiltered()
	}

	return LoreFileInfoArgsFFI{
		Paths:    valPaths,
		Revision: valRevision,
		Local:    valLocal,
		Filtered: valFiltered,
	}, cleanup
}

type LoreFileDiffArgs struct {
	/* An array of paths */
	Paths []string
	/* Source revision */
	SourceRevision string
	/* Target revision */
	TargetRevision string
	/* Produce three-way merge output with conflict markers */
	Diff3 bool
	/* Number of unchanged context lines per unified-diff hunk */
	ContextLines uint32
	/* Treat lines that differ only in trailing whitespace as equal */
	IgnoreWhitespaceEol bool
	/* Collapse runs of internal whitespace to a single space for comparison */
	IgnoreWhitespaceInline bool
}

type LoreFileDiffArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
	/* Source revision */
	SourceRevision LoreString
	/* Target revision */
	TargetRevision LoreString
	/* Produce three-way merge output with conflict markers */
	Diff3 uint8
	/* Number of unchanged context lines per unified-diff hunk */
	ContextLines uint32
	/* Treat lines that differ only in trailing whitespace as equal */
	IgnoreWhitespaceEol uint8
	/* Collapse runs of internal whitespace to a single space for comparison */
	IgnoreWhitespaceInline uint8
}

func NewLoreFileDiffArgs(opts LoreFileDiffArgs) (LoreFileDiffArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valSourceRevision, cleanupSourceRevision := NewLoreString(opts.SourceRevision)
	valTargetRevision, cleanupTargetRevision := NewLoreString(opts.TargetRevision)
	valDiff3, cleanupDiff3 := Newuint8(opts.Diff3)
	valIgnoreWhitespaceEol, cleanupIgnoreWhitespaceEol := Newuint8(opts.IgnoreWhitespaceEol)
	valIgnoreWhitespaceInline, cleanupIgnoreWhitespaceInline := Newuint8(opts.IgnoreWhitespaceInline)

	cleanup := func() {
		cleanupPaths()
		cleanupSourceRevision()
		cleanupTargetRevision()
		cleanupDiff3()
		cleanupIgnoreWhitespaceEol()
		cleanupIgnoreWhitespaceInline()
	}

	return LoreFileDiffArgsFFI{
		Paths:                  valPaths,
		SourceRevision:         valSourceRevision,
		TargetRevision:         valTargetRevision,
		Diff3:                  valDiff3,
		ContextLines:           opts.ContextLines,
		IgnoreWhitespaceEol:    valIgnoreWhitespaceEol,
		IgnoreWhitespaceInline: valIgnoreWhitespaceInline,
	}, cleanup
}

type LoreFileHashArgs struct {
	/* An array of paths */
	Paths []string
}

type LoreFileHashArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
}

func NewLoreFileHashArgs(opts LoreFileHashArgs) (LoreFileHashArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreFileHashArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreFileHistoryArgs struct {
	/* A path to a file */
	Path string
	/* Optional revision specifier to start from */
	Revision string
	/* Restrict history to revisions on this branch */
	Branch string
	/* Number of revisions to list */
	Length uint32
	/* Number of revisions to search initially */
	Depth uint32
}

type LoreFileHistoryArgsFFI struct {
	/* A path to a file */
	Path LoreString
	/* Optional revision specifier to start from */
	Revision LoreString
	/* Restrict history to revisions on this branch */
	Branch LoreString
	/* Number of revisions to list */
	Length uint32
	/* Number of revisions to search initially */
	Depth uint32
}

func NewLoreFileHistoryArgs(opts LoreFileHistoryArgs) (LoreFileHistoryArgsFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupPath()
		cleanupRevision()
		cleanupBranch()
	}

	return LoreFileHistoryArgsFFI{
		Path:     valPath,
		Revision: valRevision,
		Branch:   valBranch,
		Length:   opts.Length,
		Depth:    opts.Depth,
	}, cleanup
}

type LoreFileMetadataClearArgs struct {
	/* Which file to clear metadata for */
	Path string
}

type LoreFileMetadataClearArgsFFI struct {
	/* Which file to clear metadata for */
	Path LoreString
}

func NewLoreFileMetadataClearArgs(opts LoreFileMetadataClearArgs) (LoreFileMetadataClearArgsFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupPath()
	}

	return LoreFileMetadataClearArgsFFI{
		Path: valPath,
	}, cleanup
}

type LoreFileMetadataGetArgs struct {
	/* Revision to get metadata for */
	Revision string
	/* Where to get metadata for */
	Path string
	/* Metadata key */
	Key string
}

type LoreFileMetadataGetArgsFFI struct {
	/* Revision to get metadata for */
	Revision LoreString
	/* Where to get metadata for */
	Path LoreString
	/* Metadata key */
	Key LoreString
}

func NewLoreFileMetadataGetArgs(opts LoreFileMetadataGetArgs) (LoreFileMetadataGetArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valPath, cleanupPath := NewLoreString(opts.Path)
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupRevision()
		cleanupPath()
		cleanupKey()
	}

	return LoreFileMetadataGetArgsFFI{
		Revision: valRevision,
		Path:     valPath,
		Key:      valKey,
	}, cleanup
}

type LoreFileMetadataListArgs struct {
	/* What to list metadata for */
	Path string
	/* Revision to list metadata for */
	Revision string
}

type LoreFileMetadataListArgsFFI struct {
	/* What to list metadata for */
	Path LoreString
	/* Revision to list metadata for */
	Revision LoreString
}

func NewLoreFileMetadataListArgs(opts LoreFileMetadataListArgs) (LoreFileMetadataListArgsFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)

	cleanup := func() {
		cleanupPath()
		cleanupRevision()
	}

	return LoreFileMetadataListArgsFFI{
		Path:     valPath,
		Revision: valRevision,
	}, cleanup
}

type LoreFileMetadataSetArgs struct {
	/* An array of paths */
	Paths []string
	/* An array of keys */
	Keys []string
	/* An array of values */
	Values []string
	/* Pointer to an array of formats */
	Formats LoreMetadataTypeArray
	/* Pointer to an array of entry counts per path */
	Entries LoreUint32Array
}

type LoreFileMetadataSetArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
	/* An array of keys */
	Keys LoreStringArrayFFI
	/* An array of values */
	Values LoreStringArrayFFI
	/* Pointer to an array of formats */
	Formats LoreMetadataTypeArrayFFI
	/* Pointer to an array of entry counts per path */
	Entries LoreUint32ArrayFFI
}

func NewLoreFileMetadataSetArgs(opts LoreFileMetadataSetArgs) (LoreFileMetadataSetArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)
	valValues, cleanupValues := NewLoreStringArray(opts.Values)
	valFormats, cleanupFormats := NewLoreMetadataTypeArray(opts.Formats)
	valEntries, cleanupEntries := NewLoreUint32Array(opts.Entries)

	cleanup := func() {
		cleanupPaths()
		cleanupKeys()
		cleanupValues()
		cleanupFormats()
		cleanupEntries()
	}

	return LoreFileMetadataSetArgsFFI{
		Paths:   valPaths,
		Keys:    valKeys,
		Values:  valValues,
		Formats: valFormats,
		Entries: valEntries,
	}, cleanup
}

type LoreFileResetArgs struct {
	/* Pointer to an array of paths */
	Paths []string
	/* Revision to reset files into */
	Revision string
	/* Purge untracked files */
	Purge bool
}

type LoreFileResetArgsFFI struct {
	/* Pointer to an array of paths */
	Paths LoreStringArrayFFI
	/* Revision to reset files into */
	Revision LoreString
	/* Purge untracked files */
	Purge uint8
}

func NewLoreFileResetArgs(opts LoreFileResetArgs) (LoreFileResetArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valPurge, cleanupPurge := Newuint8(opts.Purge)

	cleanup := func() {
		cleanupPaths()
		cleanupRevision()
		cleanupPurge()
	}

	return LoreFileResetArgsFFI{
		Paths:    valPaths,
		Revision: valRevision,
		Purge:    valPurge,
	}, cleanup
}

type LoreFileResetToLastMergedArgs struct {
	/* Pointer to an array of paths */
	Paths []string
	/* Branch whose last merged revision to reset to */
	Branch string
	/* Purge untracked files */
	Purge bool
}

type LoreFileResetToLastMergedArgsFFI struct {
	/* Pointer to an array of paths */
	Paths LoreStringArrayFFI
	/* Branch whose last merged revision to reset to */
	Branch LoreString
	/* Purge untracked files */
	Purge uint8
}

func NewLoreFileResetToLastMergedArgs(opts LoreFileResetToLastMergedArgs) (LoreFileResetToLastMergedArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valPurge, cleanupPurge := Newuint8(opts.Purge)

	cleanup := func() {
		cleanupPaths()
		cleanupBranch()
		cleanupPurge()
	}

	return LoreFileResetToLastMergedArgsFFI{
		Paths:  valPaths,
		Branch: valBranch,
		Purge:  valPurge,
	}, cleanup
}

type LoreFileStageArgs struct {
	/* An array of paths */
	Paths []string
	/* Case change handling, 0 = error, 1 = update filesystem (keep), 2 = update repository (rename) */
	CaseChange uint32
	/* Force a recursive filesystem scan of directory paths (no effect on file paths) */
	Scan bool
}

type LoreFileStageArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
	/* Case change handling, 0 = error, 1 = update filesystem (keep), 2 = update repository (rename) */
	CaseChange uint32
	/* Force a recursive filesystem scan of directory paths (no effect on file paths) */
	Scan uint8
}

func NewLoreFileStageArgs(opts LoreFileStageArgs) (LoreFileStageArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valScan, cleanupScan := Newuint8(opts.Scan)

	cleanup := func() {
		cleanupPaths()
		cleanupScan()
	}

	return LoreFileStageArgsFFI{
		Paths:      valPaths,
		CaseChange: opts.CaseChange,
		Scan:       valScan,
	}, cleanup
}

type LoreFileStageMergeArgs struct {
	/* Paths to files to stage as merge */
	Paths []string
}

type LoreFileStageMergeArgsFFI struct {
	/* Paths to files to stage as merge */
	Paths LoreStringArrayFFI
}

func NewLoreFileStageMergeArgs(opts LoreFileStageMergeArgs) (LoreFileStageMergeArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreFileStageMergeArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreFileStageMoveArgs struct {
	/* Original path of file */
	FromPath string
	/* New path of file */
	ToPath string
}

type LoreFileStageMoveArgsFFI struct {
	/* Original path of file */
	FromPath LoreString
	/* New path of file */
	ToPath LoreString
}

func NewLoreFileStageMoveArgs(opts LoreFileStageMoveArgs) (LoreFileStageMoveArgsFFI, func()) {
	valFromPath, cleanupFromPath := NewLoreString(opts.FromPath)
	valToPath, cleanupToPath := NewLoreString(opts.ToPath)

	cleanup := func() {
		cleanupFromPath()
		cleanupToPath()
	}

	return LoreFileStageMoveArgsFFI{
		FromPath: valFromPath,
		ToPath:   valToPath,
	}, cleanup
}

type LoreFileDirtyArgs struct {
	/* An array of paths */
	Paths []string
}

type LoreFileDirtyArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
}

func NewLoreFileDirtyArgs(opts LoreFileDirtyArgs) (LoreFileDirtyArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreFileDirtyArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreFileDirtyMoveArgs struct {
	/* Original path of file */
	FromPath string
	/* New path of file */
	ToPath string
}

type LoreFileDirtyMoveArgsFFI struct {
	/* Original path of file */
	FromPath LoreString
	/* New path of file */
	ToPath LoreString
}

func NewLoreFileDirtyMoveArgs(opts LoreFileDirtyMoveArgs) (LoreFileDirtyMoveArgsFFI, func()) {
	valFromPath, cleanupFromPath := NewLoreString(opts.FromPath)
	valToPath, cleanupToPath := NewLoreString(opts.ToPath)

	cleanup := func() {
		cleanupFromPath()
		cleanupToPath()
	}

	return LoreFileDirtyMoveArgsFFI{
		FromPath: valFromPath,
		ToPath:   valToPath,
	}, cleanup
}

type LoreFileDirtyCopyArgs struct {
	/* Source path of file */
	FromPath string
	/* Destination path of copy */
	ToPath string
}

type LoreFileDirtyCopyArgsFFI struct {
	/* Source path of file */
	FromPath LoreString
	/* Destination path of copy */
	ToPath LoreString
}

func NewLoreFileDirtyCopyArgs(opts LoreFileDirtyCopyArgs) (LoreFileDirtyCopyArgsFFI, func()) {
	valFromPath, cleanupFromPath := NewLoreString(opts.FromPath)
	valToPath, cleanupToPath := NewLoreString(opts.ToPath)

	cleanup := func() {
		cleanupFromPath()
		cleanupToPath()
	}

	return LoreFileDirtyCopyArgsFFI{
		FromPath: valFromPath,
		ToPath:   valToPath,
	}, cleanup
}

type LoreFileUnstageArgs struct {
	/* An array of paths */
	Paths []string
}

type LoreFileUnstageArgsFFI struct {
	/* An array of paths */
	Paths LoreStringArrayFFI
}

func NewLoreFileUnstageArgs(opts LoreFileUnstageArgs) (LoreFileUnstageArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreFileUnstageArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreFileWriteArgs struct {
	/* Address of data to write; takes precedence over `path` when non-empty */
	Address string
	/* Repository path to the file; used when `address` is empty */
	Path string
	/* Revision of the file to write (used with `path`) */
	Revision string
	/* Destination filesystem path to write to */
	Output string
}

type LoreFileWriteArgsFFI struct {
	/* Address of data to write; takes precedence over `path` when non-empty */
	Address LoreString
	/* Repository path to the file; used when `address` is empty */
	Path LoreString
	/* Revision of the file to write (used with `path`) */
	Revision LoreString
	/* Destination filesystem path to write to */
	Output LoreString
}

func NewLoreFileWriteArgs(opts LoreFileWriteArgs) (LoreFileWriteArgsFFI, func()) {
	valAddress, cleanupAddress := NewLoreString(opts.Address)
	valPath, cleanupPath := NewLoreString(opts.Path)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valOutput, cleanupOutput := NewLoreString(opts.Output)

	cleanup := func() {
		cleanupAddress()
		cleanupPath()
		cleanupRevision()
		cleanupOutput()
	}

	return LoreFileWriteArgsFFI{
		Address:  valAddress,
		Path:     valPath,
		Revision: valRevision,
		Output:   valOutput,
	}, cleanup
}

type LoreFileObliterateArgs struct {
	/* Address of data to obliterate; takes precedence over `path` when non-empty */
	Address string
	/* Repository path to obliterate; used when `address` is empty */
	Path string
}

type LoreFileObliterateArgsFFI struct {
	/* Address of data to obliterate; takes precedence over `path` when non-empty */
	Address LoreString
	/* Repository path to obliterate; used when `address` is empty */
	Path LoreString
}

func NewLoreFileObliterateArgs(opts LoreFileObliterateArgs) (LoreFileObliterateArgsFFI, func()) {
	valAddress, cleanupAddress := NewLoreString(opts.Address)
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupAddress()
		cleanupPath()
	}

	return LoreFileObliterateArgsFFI{
		Address: valAddress,
		Path:    valPath,
	}, cleanup
}

type LoreFileDumpArgs struct {
	/* Address of data to dump; takes precedence over `path` when non-empty */
	Address string
	/* Repository path to dump; used when `address` is empty */
	Path string
}

type LoreFileDumpArgsFFI struct {
	/* Address of data to dump; takes precedence over `path` when non-empty */
	Address LoreString
	/* Repository path to dump; used when `address` is empty */
	Path LoreString
}

func NewLoreFileDumpArgs(opts LoreFileDumpArgs) (LoreFileDumpArgsFFI, func()) {
	valAddress, cleanupAddress := NewLoreString(opts.Address)
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupAddress()
		cleanupPath()
	}

	return LoreFileDumpArgsFFI{
		Address: valAddress,
		Path:    valPath,
	}, cleanup
}

type LoreFileDependencyAddArgs struct {
	/* Source file paths that will have dependencies added. */
	Paths []string
	/* Dependency target file paths (flat array). */
	Dependencies []string
	/* Tags to apply to the added dependencies (flat array). */
	Tags []string
	/* Number of dependencies per source file path. */
	DepCounts LoreUint32Array
	/* Number of tags per dependency entry. */
	TagCounts LoreUint32Array
	/* Skip cycle detection. */
	Force bool
}

type LoreFileDependencyAddArgsFFI struct {
	/* Source file paths that will have dependencies added. */
	Paths LoreStringArrayFFI
	/* Dependency target file paths (flat array). */
	Dependencies LoreStringArrayFFI
	/* Tags to apply to the added dependencies (flat array). */
	Tags LoreStringArrayFFI
	/* Number of dependencies per source file path. */
	DepCounts LoreUint32ArrayFFI
	/* Number of tags per dependency entry. */
	TagCounts LoreUint32ArrayFFI
	/* Skip cycle detection. */
	Force uint8
}

func NewLoreFileDependencyAddArgs(opts LoreFileDependencyAddArgs) (LoreFileDependencyAddArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valDependencies, cleanupDependencies := NewLoreStringArray(opts.Dependencies)
	valTags, cleanupTags := NewLoreStringArray(opts.Tags)
	valDepCounts, cleanupDepCounts := NewLoreUint32Array(opts.DepCounts)
	valTagCounts, cleanupTagCounts := NewLoreUint32Array(opts.TagCounts)
	valForce, cleanupForce := Newuint8(opts.Force)

	cleanup := func() {
		cleanupPaths()
		cleanupDependencies()
		cleanupTags()
		cleanupDepCounts()
		cleanupTagCounts()
		cleanupForce()
	}

	return LoreFileDependencyAddArgsFFI{
		Paths:        valPaths,
		Dependencies: valDependencies,
		Tags:         valTags,
		DepCounts:    valDepCounts,
		TagCounts:    valTagCounts,
		Force:        valForce,
	}, cleanup
}

type LoreFileDependencyRemoveArgs struct {
	/* Source file paths to remove dependencies from. */
	Paths []string
	/* Dependency target paths to remove (flat array). */
	Dependencies []string
	/* Tags to remove. */
	Tags []string
	/* Number of dependencies per source file. */
	DepCounts LoreUint32Array
	/* Number of tags per dependency entry. */
	TagCounts LoreUint32Array
}

type LoreFileDependencyRemoveArgsFFI struct {
	/* Source file paths to remove dependencies from. */
	Paths LoreStringArrayFFI
	/* Dependency target paths to remove (flat array). */
	Dependencies LoreStringArrayFFI
	/* Tags to remove. */
	Tags LoreStringArrayFFI
	/* Number of dependencies per source file. */
	DepCounts LoreUint32ArrayFFI
	/* Number of tags per dependency entry. */
	TagCounts LoreUint32ArrayFFI
}

func NewLoreFileDependencyRemoveArgs(opts LoreFileDependencyRemoveArgs) (LoreFileDependencyRemoveArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valDependencies, cleanupDependencies := NewLoreStringArray(opts.Dependencies)
	valTags, cleanupTags := NewLoreStringArray(opts.Tags)
	valDepCounts, cleanupDepCounts := NewLoreUint32Array(opts.DepCounts)
	valTagCounts, cleanupTagCounts := NewLoreUint32Array(opts.TagCounts)

	cleanup := func() {
		cleanupPaths()
		cleanupDependencies()
		cleanupTags()
		cleanupDepCounts()
		cleanupTagCounts()
	}

	return LoreFileDependencyRemoveArgsFFI{
		Paths:        valPaths,
		Dependencies: valDependencies,
		Tags:         valTags,
		DepCounts:    valDepCounts,
		TagCounts:    valTagCounts,
	}, cleanup
}

type LoreFileDependencyListArgs struct {
	/* Files to query. */
	Paths []string
	/* Revision to query at. */
	Revision string
	/* Follow transitive dependencies recursively. */
	Recursive bool
	/* Return dependents instead of dependencies. */
	Reverse bool
	/* Filter results by tags. */
	Tags []string
	/* Maximum recursion depth (0 = unlimited). */
	DepthLimit uint32
}

type LoreFileDependencyListArgsFFI struct {
	/* Files to query. */
	Paths LoreStringArrayFFI
	/* Revision to query at. */
	Revision LoreString
	/* Follow transitive dependencies recursively. */
	Recursive uint8
	/* Return dependents instead of dependencies. */
	Reverse uint8
	/* Filter results by tags. */
	Tags LoreStringArrayFFI
	/* Maximum recursion depth (0 = unlimited). */
	DepthLimit uint32
}

func NewLoreFileDependencyListArgs(opts LoreFileDependencyListArgs) (LoreFileDependencyListArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valRecursive, cleanupRecursive := Newuint8(opts.Recursive)
	valReverse, cleanupReverse := Newuint8(opts.Reverse)
	valTags, cleanupTags := NewLoreStringArray(opts.Tags)

	cleanup := func() {
		cleanupPaths()
		cleanupRevision()
		cleanupRecursive()
		cleanupReverse()
		cleanupTags()
	}

	return LoreFileDependencyListArgsFFI{
		Paths:      valPaths,
		Revision:   valRevision,
		Recursive:  valRecursive,
		Reverse:    valReverse,
		Tags:       valTags,
		DepthLimit: opts.DepthLimit,
	}, cleanup
}

type LoreLockFileAcquireArgs struct {
	/* Paths to acquire locks on */
	Paths []string
	/* Branch the locks are acquired on */
	Branch string
}

type LoreLockFileAcquireArgsFFI struct {
	/* Paths to acquire locks on */
	Paths LoreStringArrayFFI
	/* Branch the locks are acquired on */
	Branch LoreString
}

func NewLoreLockFileAcquireArgs(opts LoreLockFileAcquireArgs) (LoreLockFileAcquireArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupPaths()
		cleanupBranch()
	}

	return LoreLockFileAcquireArgsFFI{
		Paths:  valPaths,
		Branch: valBranch,
	}, cleanup
}

type LoreLockFileStatusArgs struct {
	/* Paths to get the lock status of */
	Paths []string
	/* Branch the locks were acquired on */
	Branch string
}

type LoreLockFileStatusArgsFFI struct {
	/* Paths to get the lock status of */
	Paths LoreStringArrayFFI
	/* Branch the locks were acquired on */
	Branch LoreString
}

func NewLoreLockFileStatusArgs(opts LoreLockFileStatusArgs) (LoreLockFileStatusArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)

	cleanup := func() {
		cleanupPaths()
		cleanupBranch()
	}

	return LoreLockFileStatusArgsFFI{
		Paths:  valPaths,
		Branch: valBranch,
	}, cleanup
}

type LoreLockFileQueryArgs struct {
	/* Branch to query locks on */
	Branch string
	/* Owner filter; empty matches any owner */
	Owner string
	/* Path filter; empty matches any path */
	Path string
}

type LoreLockFileQueryArgsFFI struct {
	/* Branch to query locks on */
	Branch LoreString
	/* Owner filter; empty matches any owner */
	Owner LoreString
	/* Path filter; empty matches any path */
	Path LoreString
}

func NewLoreLockFileQueryArgs(opts LoreLockFileQueryArgs) (LoreLockFileQueryArgsFFI, func()) {
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valOwner, cleanupOwner := NewLoreString(opts.Owner)
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupBranch()
		cleanupOwner()
		cleanupPath()
	}

	return LoreLockFileQueryArgsFFI{
		Branch: valBranch,
		Owner:  valOwner,
		Path:   valPath,
	}, cleanup
}

type LoreLockFileReleaseArgs struct {
	/* Paths to release locks on */
	Paths []string
	/* Branch the locks were acquired on */
	Branch string
	/* Owner of the lock */
	Owner string
	/* Owner id of the lock */
	OwnerId string
}

type LoreLockFileReleaseArgsFFI struct {
	/* Paths to release locks on */
	Paths LoreStringArrayFFI
	/* Branch the locks were acquired on */
	Branch LoreString
	/* Owner of the lock */
	Owner LoreString
	/* Owner id of the lock */
	OwnerId LoreString
}

func NewLoreLockFileReleaseArgs(opts LoreLockFileReleaseArgs) (LoreLockFileReleaseArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valOwner, cleanupOwner := NewLoreString(opts.Owner)
	valOwnerId, cleanupOwnerId := NewLoreString(opts.OwnerId)

	cleanup := func() {
		cleanupPaths()
		cleanupBranch()
		cleanupOwner()
		cleanupOwnerId()
	}

	return LoreLockFileReleaseArgsFFI{
		Paths:   valPaths,
		Branch:  valBranch,
		Owner:   valOwner,
		OwnerId: valOwnerId,
	}, cleanup
}

type LoreLinkAddArgs struct {
	/* Link repository URL */
	Link string
	/* Path within this repository where the link is added */
	LinkPath string
	/* Source path within the linked repository; `/` or `\` means the root */
	SourcePath string
	/* Branch or revision to set the link pin at */
	Pin string
	/* Disable automatic branch creation in the linked repository */
	DisableBranching bool
}

type LoreLinkAddArgsFFI struct {
	/* Link repository URL */
	Link LoreString
	/* Path within this repository where the link is added */
	LinkPath LoreString
	/* Source path within the linked repository; `/` or `\` means the root */
	SourcePath LoreString
	/* Branch or revision to set the link pin at */
	Pin LoreString
	/* Disable automatic branch creation in the linked repository */
	DisableBranching uint8
}

func NewLoreLinkAddArgs(opts LoreLinkAddArgs) (LoreLinkAddArgsFFI, func()) {
	valLink, cleanupLink := NewLoreString(opts.Link)
	valLinkPath, cleanupLinkPath := NewLoreString(opts.LinkPath)
	valSourcePath, cleanupSourcePath := NewLoreString(opts.SourcePath)
	valPin, cleanupPin := NewLoreString(opts.Pin)
	valDisableBranching, cleanupDisableBranching := Newuint8(opts.DisableBranching)

	cleanup := func() {
		cleanupLink()
		cleanupLinkPath()
		cleanupSourcePath()
		cleanupPin()
		cleanupDisableBranching()
	}

	return LoreLinkAddArgsFFI{
		Link:             valLink,
		LinkPath:         valLinkPath,
		SourcePath:       valSourcePath,
		Pin:              valPin,
		DisableBranching: valDisableBranching,
	}, cleanup
}

type LoreLinkRemoveArgs struct {
	/* Path within this repository where the link is removed */
	LinkPath string
}

type LoreLinkRemoveArgsFFI struct {
	/* Path within this repository where the link is removed */
	LinkPath LoreString
}

func NewLoreLinkRemoveArgs(opts LoreLinkRemoveArgs) (LoreLinkRemoveArgsFFI, func()) {
	valLinkPath, cleanupLinkPath := NewLoreString(opts.LinkPath)

	cleanup := func() {
		cleanupLinkPath()
	}

	return LoreLinkRemoveArgsFFI{
		LinkPath: valLinkPath,
	}, cleanup
}

type LoreLinkListArgs struct {
	Unused int
}

type LoreLinkListArgsFFI struct {
	Unused int
}

func NewLoreLinkListArgs(opts LoreLinkListArgs) (LoreLinkListArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreLinkListArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreLinkUpdateArgs struct {
	/* Path within this repository of the link to update */
	LinkPath string
	/* Branch or specific revision to pin the link to */
	Pin string
}

type LoreLinkUpdateArgsFFI struct {
	/* Path within this repository of the link to update */
	LinkPath LoreString
	/* Branch or specific revision to pin the link to */
	Pin LoreString
}

func NewLoreLinkUpdateArgs(opts LoreLinkUpdateArgs) (LoreLinkUpdateArgsFFI, func()) {
	valLinkPath, cleanupLinkPath := NewLoreString(opts.LinkPath)
	valPin, cleanupPin := NewLoreString(opts.Pin)

	cleanup := func() {
		cleanupLinkPath()
		cleanupPin()
	}

	return LoreLinkUpdateArgsFFI{
		LinkPath: valLinkPath,
		Pin:      valPin,
	}, cleanup
}

type LoreRepositoryCloneArgs struct {
	/* URL to the repository */
	RepositoryUrl string
	/* [Optional] Revision to clone */
	Revision string
	/* [Optional] Client side view filter to use */
	View string
	/* Clone without any files */
	Bare bool
	/* Clone virtually using split-write filesystem */
	Virtually bool
	/* Use direct file write */
	DirectFileWrite bool
	/* Use direct file I/O instead of memory mapping files */
	DirectFileIo bool
	/* (Optional) Layer module */
	Layer string
	/* (Optional) Layer metadata key to link revisions with */
	LayerMetadata string
	/* (Optional) File containing list of files to prefetch */
	Prefetch string
	/* Use the shared store instead of a local immutable store */
	UseSharedStore bool
	/* [Optional] Path to use for the shared store, an empty string means to use the default */
	SharedStorePath string
	/* Clone without local repository tracking (memory-only stores) */
	NoTracking bool
	/* Root files for dependency-based selective clone */
	RootFiles []string
	/* Tags to filter dependencies by during resolution */
	DependencyTags []string
	/* Follow transitive dependencies recursively */
	DependencyRecursive bool
	/* Maximum dependency traversal depth. 0 means unlimited. */
	DependencyDepthLimit uint32
}

type LoreRepositoryCloneArgsFFI struct {
	/* URL to the repository */
	RepositoryUrl LoreString
	/* [Optional] Revision to clone */
	Revision LoreString
	/* [Optional] Client side view filter to use */
	View LoreString
	/* Clone without any files */
	Bare uint8
	/* Clone virtually using split-write filesystem */
	Virtually uint8
	/* Use direct file write */
	DirectFileWrite uint8
	/* Use direct file I/O instead of memory mapping files */
	DirectFileIo uint8
	/* (Optional) Layer module */
	Layer LoreString
	/* (Optional) Layer metadata key to link revisions with */
	LayerMetadata LoreString
	/* (Optional) File containing list of files to prefetch */
	Prefetch LoreString
	/* Use the shared store instead of a local immutable store */
	UseSharedStore uint8
	/* [Optional] Path to use for the shared store, an empty string means to use the default */
	SharedStorePath LoreString
	/* Clone without local repository tracking (memory-only stores) */
	NoTracking uint8
	/* Root files for dependency-based selective clone */
	RootFiles LoreStringArrayFFI
	/* Tags to filter dependencies by during resolution */
	DependencyTags LoreStringArrayFFI
	/* Follow transitive dependencies recursively */
	DependencyRecursive uint8
	/* Maximum dependency traversal depth. 0 means unlimited. */
	DependencyDepthLimit uint32
}

func NewLoreRepositoryCloneArgs(opts LoreRepositoryCloneArgs) (LoreRepositoryCloneArgsFFI, func()) {
	valRepositoryUrl, cleanupRepositoryUrl := NewLoreString(opts.RepositoryUrl)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valView, cleanupView := NewLoreString(opts.View)
	valBare, cleanupBare := Newuint8(opts.Bare)
	valVirtually, cleanupVirtually := Newuint8(opts.Virtually)
	valDirectFileWrite, cleanupDirectFileWrite := Newuint8(opts.DirectFileWrite)
	valDirectFileIo, cleanupDirectFileIo := Newuint8(opts.DirectFileIo)
	valLayer, cleanupLayer := NewLoreString(opts.Layer)
	valLayerMetadata, cleanupLayerMetadata := NewLoreString(opts.LayerMetadata)
	valPrefetch, cleanupPrefetch := NewLoreString(opts.Prefetch)
	valUseSharedStore, cleanupUseSharedStore := Newuint8(opts.UseSharedStore)
	valSharedStorePath, cleanupSharedStorePath := NewLoreString(opts.SharedStorePath)
	valNoTracking, cleanupNoTracking := Newuint8(opts.NoTracking)
	valRootFiles, cleanupRootFiles := NewLoreStringArray(opts.RootFiles)
	valDependencyTags, cleanupDependencyTags := NewLoreStringArray(opts.DependencyTags)
	valDependencyRecursive, cleanupDependencyRecursive := Newuint8(opts.DependencyRecursive)

	cleanup := func() {
		cleanupRepositoryUrl()
		cleanupRevision()
		cleanupView()
		cleanupBare()
		cleanupVirtually()
		cleanupDirectFileWrite()
		cleanupDirectFileIo()
		cleanupLayer()
		cleanupLayerMetadata()
		cleanupPrefetch()
		cleanupUseSharedStore()
		cleanupSharedStorePath()
		cleanupNoTracking()
		cleanupRootFiles()
		cleanupDependencyTags()
		cleanupDependencyRecursive()
	}

	return LoreRepositoryCloneArgsFFI{
		RepositoryUrl:        valRepositoryUrl,
		Revision:             valRevision,
		View:                 valView,
		Bare:                 valBare,
		Virtually:            valVirtually,
		DirectFileWrite:      valDirectFileWrite,
		DirectFileIo:         valDirectFileIo,
		Layer:                valLayer,
		LayerMetadata:        valLayerMetadata,
		Prefetch:             valPrefetch,
		UseSharedStore:       valUseSharedStore,
		SharedStorePath:      valSharedStorePath,
		NoTracking:           valNoTracking,
		RootFiles:            valRootFiles,
		DependencyTags:       valDependencyTags,
		DependencyRecursive:  valDependencyRecursive,
		DependencyDepthLimit: opts.DependencyDepthLimit,
	}, cleanup
}

type LoreRepositoryInfoArgs struct {
	/* URL of the remote repository to query */
	RepositoryUrl string
}

type LoreRepositoryInfoArgsFFI struct {
	/* URL of the remote repository to query */
	RepositoryUrl LoreString
}

func NewLoreRepositoryInfoArgs(opts LoreRepositoryInfoArgs) (LoreRepositoryInfoArgsFFI, func()) {
	valRepositoryUrl, cleanupRepositoryUrl := NewLoreString(opts.RepositoryUrl)

	cleanup := func() {
		cleanupRepositoryUrl()
	}

	return LoreRepositoryInfoArgsFFI{
		RepositoryUrl: valRepositoryUrl,
	}, cleanup
}

type LoreRepositoryDumpArgs struct {
	/* Revision to dump; empty string uses the current revision */
	Revision string
	/* Repository-relative path to start dumping from; empty dumps the root */
	Path string
	/* Maximum tree traversal depth */
	MaxDepth uintptr
}

type LoreRepositoryDumpArgsFFI struct {
	/* Revision to dump; empty string uses the current revision */
	Revision LoreString
	/* Repository-relative path to start dumping from; empty dumps the root */
	Path LoreString
	/* Maximum tree traversal depth */
	MaxDepth uintptr
}

func NewLoreRepositoryDumpArgs(opts LoreRepositoryDumpArgs) (LoreRepositoryDumpArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupRevision()
		cleanupPath()
	}

	return LoreRepositoryDumpArgsFFI{
		Revision: valRevision,
		Path:     valPath,
		MaxDepth: opts.MaxDepth,
	}, cleanup
}

type LoreRepositoryCreateArgs struct {
	/* URL to the repository */
	RepositoryUrl string
	/* Optional repository description */
	Description string
	/* Optional repository ID, set to empty string to generate a new ID */
	Id string
	/* Use the shared store instead of a local immutable store */
	UseSharedStore bool
	/* [Optional] Path to use for the shared store, an empty string means to use the default */
	SharedStorePath string
}

type LoreRepositoryCreateArgsFFI struct {
	/* URL to the repository */
	RepositoryUrl LoreString
	/* Optional repository description */
	Description LoreString
	/* Optional repository ID, set to empty string to generate a new ID */
	Id LoreString
	/* Use the shared store instead of a local immutable store */
	UseSharedStore uint8
	/* [Optional] Path to use for the shared store, an empty string means to use the default */
	SharedStorePath LoreString
}

func NewLoreRepositoryCreateArgs(opts LoreRepositoryCreateArgs) (LoreRepositoryCreateArgsFFI, func()) {
	valRepositoryUrl, cleanupRepositoryUrl := NewLoreString(opts.RepositoryUrl)
	valDescription, cleanupDescription := NewLoreString(opts.Description)
	valId, cleanupId := NewLoreString(opts.Id)
	valUseSharedStore, cleanupUseSharedStore := Newuint8(opts.UseSharedStore)
	valSharedStorePath, cleanupSharedStorePath := NewLoreString(opts.SharedStorePath)

	cleanup := func() {
		cleanupRepositoryUrl()
		cleanupDescription()
		cleanupId()
		cleanupUseSharedStore()
		cleanupSharedStorePath()
	}

	return LoreRepositoryCreateArgsFFI{
		RepositoryUrl:   valRepositoryUrl,
		Description:     valDescription,
		Id:              valId,
		UseSharedStore:  valUseSharedStore,
		SharedStorePath: valSharedStorePath,
	}, cleanup
}

type LoreRepositoryFlushArgs struct {
	Unused int
}

type LoreRepositoryFlushArgsFFI struct {
	Unused int
}

func NewLoreRepositoryFlushArgs(opts LoreRepositoryFlushArgs) (LoreRepositoryFlushArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryFlushArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryGcArgs struct {
	Unused int
}

type LoreRepositoryGcArgsFFI struct {
	Unused int
}

func NewLoreRepositoryGcArgs(opts LoreRepositoryGcArgs) (LoreRepositoryGcArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryGcArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryReleaseArgs struct {
	Unused int
}

type LoreRepositoryReleaseArgsFFI struct {
	Unused int
}

func NewLoreRepositoryReleaseArgs(opts LoreRepositoryReleaseArgs) (LoreRepositoryReleaseArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryReleaseArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreLayerAddArgs struct {
	/* Path in the current repository where the layer should be placed */
	TargetPath string
	/* Repository to add as a layer */
	SourceRepository string
	/* Path in the layer repository where the layer should start */
	SourcePath string
	/* Metadata key to use to match revisions */
	Metadata string
}

type LoreLayerAddArgsFFI struct {
	/* Path in the current repository where the layer should be placed */
	TargetPath LoreString
	/* Repository to add as a layer */
	SourceRepository LoreString
	/* Path in the layer repository where the layer should start */
	SourcePath LoreString
	/* Metadata key to use to match revisions */
	Metadata LoreString
}

func NewLoreLayerAddArgs(opts LoreLayerAddArgs) (LoreLayerAddArgsFFI, func()) {
	valTargetPath, cleanupTargetPath := NewLoreString(opts.TargetPath)
	valSourceRepository, cleanupSourceRepository := NewLoreString(opts.SourceRepository)
	valSourcePath, cleanupSourcePath := NewLoreString(opts.SourcePath)
	valMetadata, cleanupMetadata := NewLoreString(opts.Metadata)

	cleanup := func() {
		cleanupTargetPath()
		cleanupSourceRepository()
		cleanupSourcePath()
		cleanupMetadata()
	}

	return LoreLayerAddArgsFFI{
		TargetPath:       valTargetPath,
		SourceRepository: valSourceRepository,
		SourcePath:       valSourcePath,
		Metadata:         valMetadata,
	}, cleanup
}

type LoreLayerRemoveArgs struct {
	/* Path in the current repository where the layer is placed */
	TargetPath string
	/* Repository added as a layer at the given path */
	SourceRepository string
	/* Remove all untracked files and directories inside the layer mount */
	Purge bool
}

type LoreLayerRemoveArgsFFI struct {
	/* Path in the current repository where the layer is placed */
	TargetPath LoreString
	/* Repository added as a layer at the given path */
	SourceRepository LoreString
	/* Remove all untracked files and directories inside the layer mount */
	Purge uint8
}

func NewLoreLayerRemoveArgs(opts LoreLayerRemoveArgs) (LoreLayerRemoveArgsFFI, func()) {
	valTargetPath, cleanupTargetPath := NewLoreString(opts.TargetPath)
	valSourceRepository, cleanupSourceRepository := NewLoreString(opts.SourceRepository)
	valPurge, cleanupPurge := Newuint8(opts.Purge)

	cleanup := func() {
		cleanupTargetPath()
		cleanupSourceRepository()
		cleanupPurge()
	}

	return LoreLayerRemoveArgsFFI{
		TargetPath:       valTargetPath,
		SourceRepository: valSourceRepository,
		Purge:            valPurge,
	}, cleanup
}

type LoreLayerListArgs struct {
	Unused int
}

type LoreLayerListArgsFFI struct {
	Unused int
}

func NewLoreLayerListArgs(opts LoreLayerListArgs) (LoreLayerListArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreLayerListArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryListArgs struct {
	/* Remote URL to list repositories from */
	Url string
}

type LoreRepositoryListArgsFFI struct {
	/* Remote URL to list repositories from */
	Url LoreString
}

func NewLoreRepositoryListArgs(opts LoreRepositoryListArgs) (LoreRepositoryListArgsFFI, func()) {
	valUrl, cleanupUrl := NewLoreString(opts.Url)

	cleanup := func() {
		cleanupUrl()
	}

	return LoreRepositoryListArgsFFI{
		Url: valUrl,
	}, cleanup
}

type LoreRepositoryStatusArgs struct {
	/* Include staged state in the report */
	Staged bool
	/* Reconcile against the filesystem and refresh dirty tracking.

	By default, status reports the currently tracked state: the
	staged revision (if any) plus any files and directories already
	marked dirty. No filesystem reads are performed beyond the existing
	dirty flags — clean or unmarked files on disk are not inspected even
	if they differ from the current revision.

	When enabled, the filesystem is walked under each requested path, every
	file is reconciled against the current revision, and dirty flags are
	set or cleared accordingly. The refreshed flags are persisted in the
	staged state so subsequent operations (commit, stage, status) see an
	accurate picture without rescanning. */
	Scan bool
	/* Verify dirty flags against the filesystem without a full scan.

	When enabled, files already marked dirty are re-examined individually: a
	dirty file whose on-disk content matches its tracked node (same size,
	and same content when the modification time differs) has its dirty flag
	cleared and is omitted from the report, unless it is also staged.
	Structural dirty actions (add/move/copy/delete) are always reported.
	The refreshed flags are persisted in the staged state. */
	CheckDirty bool
	/* Reset the tracked state before computing status */
	Reset bool
	/* Include the sync point in the report */
	SyncPoint bool
	/* Only emit revision info, skipping all diffs */
	RevisionOnly bool
	/* Count directories and files (view-filtered) in the staged state if
	present, otherwise the current revision */
	Count bool
	/* Repository-relative paths to limit the status check to; empty checks all */
	Paths []string
}

type LoreRepositoryStatusArgsFFI struct {
	/* Include staged state in the report */
	Staged uint8
	/* Reconcile against the filesystem and refresh dirty tracking.

	By default, status reports the currently tracked state: the
	staged revision (if any) plus any files and directories already
	marked dirty. No filesystem reads are performed beyond the existing
	dirty flags — clean or unmarked files on disk are not inspected even
	if they differ from the current revision.

	When enabled, the filesystem is walked under each requested path, every
	file is reconciled against the current revision, and dirty flags are
	set or cleared accordingly. The refreshed flags are persisted in the
	staged state so subsequent operations (commit, stage, status) see an
	accurate picture without rescanning. */
	Scan uint8
	/* Verify dirty flags against the filesystem without a full scan.

	When enabled, files already marked dirty are re-examined individually: a
	dirty file whose on-disk content matches its tracked node (same size,
	and same content when the modification time differs) has its dirty flag
	cleared and is omitted from the report, unless it is also staged.
	Structural dirty actions (add/move/copy/delete) are always reported.
	The refreshed flags are persisted in the staged state. */
	CheckDirty uint8
	/* Reset the tracked state before computing status */
	Reset uint8
	/* Include the sync point in the report */
	SyncPoint uint8
	/* Only emit revision info, skipping all diffs */
	RevisionOnly uint8
	/* Count directories and files (view-filtered) in the staged state if
	present, otherwise the current revision */
	Count uint8
	/* Repository-relative paths to limit the status check to; empty checks all */
	Paths LoreStringArrayFFI
}

func NewLoreRepositoryStatusArgs(opts LoreRepositoryStatusArgs) (LoreRepositoryStatusArgsFFI, func()) {
	valStaged, cleanupStaged := Newuint8(opts.Staged)
	valScan, cleanupScan := Newuint8(opts.Scan)
	valCheckDirty, cleanupCheckDirty := Newuint8(opts.CheckDirty)
	valReset, cleanupReset := Newuint8(opts.Reset)
	valSyncPoint, cleanupSyncPoint := Newuint8(opts.SyncPoint)
	valRevisionOnly, cleanupRevisionOnly := Newuint8(opts.RevisionOnly)
	valCount, cleanupCount := Newuint8(opts.Count)
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupStaged()
		cleanupScan()
		cleanupCheckDirty()
		cleanupReset()
		cleanupSyncPoint()
		cleanupRevisionOnly()
		cleanupCount()
		cleanupPaths()
	}

	return LoreRepositoryStatusArgsFFI{
		Staged:       valStaged,
		Scan:         valScan,
		CheckDirty:   valCheckDirty,
		Reset:        valReset,
		SyncPoint:    valSyncPoint,
		RevisionOnly: valRevisionOnly,
		Count:        valCount,
		Paths:        valPaths,
	}, cleanup
}

type LoreRepositoryStoreImmutableQueryArgs struct {
	/* Fragment address to query */
	Address string
	/* Recurse into and query subfragments */
	Recurse bool
}

type LoreRepositoryStoreImmutableQueryArgsFFI struct {
	/* Fragment address to query */
	Address LoreString
	/* Recurse into and query subfragments */
	Recurse uint8
}

func NewLoreRepositoryStoreImmutableQueryArgs(opts LoreRepositoryStoreImmutableQueryArgs) (LoreRepositoryStoreImmutableQueryArgsFFI, func()) {
	valAddress, cleanupAddress := NewLoreString(opts.Address)
	valRecurse, cleanupRecurse := Newuint8(opts.Recurse)

	cleanup := func() {
		cleanupAddress()
		cleanupRecurse()
	}

	return LoreRepositoryStoreImmutableQueryArgsFFI{
		Address: valAddress,
		Recurse: valRecurse,
	}, cleanup
}

type LoreRepositoryVerifyStateArgs struct {
	/* Repository-relative path to verify; empty verifies the whole repository */
	Path string
	/* Heal detected inconsistencies */
	Heal bool
}

type LoreRepositoryVerifyStateArgsFFI struct {
	/* Repository-relative path to verify; empty verifies the whole repository */
	Path LoreString
	/* Heal detected inconsistencies */
	Heal uint8
}

func NewLoreRepositoryVerifyStateArgs(opts LoreRepositoryVerifyStateArgs) (LoreRepositoryVerifyStateArgsFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)
	valHeal, cleanupHeal := Newuint8(opts.Heal)

	cleanup := func() {
		cleanupPath()
		cleanupHeal()
	}

	return LoreRepositoryVerifyStateArgsFFI{
		Path: valPath,
		Heal: valHeal,
	}, cleanup
}

type LoreRevisionCommitArgs struct {
	/* Commit message */
	Message string
	/* If set, commit only this linked repository (mount path relative to repo root) */
	Link string
	/* Array of link relative paths that have specific messages */
	LinkPaths []string
	/* Array of messages corresponding to each link path (parallel array with `link_paths`) */
	LinkMessages []string
	/* If set, commit only this layer (mount path relative to repo root) */
	Layer string
	/* Array of layer mount paths that have specific messages */
	LayerPaths []string
	/* Array of messages corresponding to each layer path (parallel array with `layer_paths`) */
	LayerMessages []string
	/* Emit per-fragment write stats during the commit */
	Stats bool
}

type LoreRevisionCommitArgsFFI struct {
	/* Commit message */
	Message LoreString
	/* If set, commit only this linked repository (mount path relative to repo root) */
	Link LoreString
	/* Array of link relative paths that have specific messages */
	LinkPaths LoreStringArrayFFI
	/* Array of messages corresponding to each link path (parallel array with `link_paths`) */
	LinkMessages LoreStringArrayFFI
	/* If set, commit only this layer (mount path relative to repo root) */
	Layer LoreString
	/* Array of layer mount paths that have specific messages */
	LayerPaths LoreStringArrayFFI
	/* Array of messages corresponding to each layer path (parallel array with `layer_paths`) */
	LayerMessages LoreStringArrayFFI
	/* Emit per-fragment write stats during the commit */
	Stats uint8
}

func NewLoreRevisionCommitArgs(opts LoreRevisionCommitArgs) (LoreRevisionCommitArgsFFI, func()) {
	valMessage, cleanupMessage := NewLoreString(opts.Message)
	valLink, cleanupLink := NewLoreString(opts.Link)
	valLinkPaths, cleanupLinkPaths := NewLoreStringArray(opts.LinkPaths)
	valLinkMessages, cleanupLinkMessages := NewLoreStringArray(opts.LinkMessages)
	valLayer, cleanupLayer := NewLoreString(opts.Layer)
	valLayerPaths, cleanupLayerPaths := NewLoreStringArray(opts.LayerPaths)
	valLayerMessages, cleanupLayerMessages := NewLoreStringArray(opts.LayerMessages)
	valStats, cleanupStats := Newuint8(opts.Stats)

	cleanup := func() {
		cleanupMessage()
		cleanupLink()
		cleanupLinkPaths()
		cleanupLinkMessages()
		cleanupLayer()
		cleanupLayerPaths()
		cleanupLayerMessages()
		cleanupStats()
	}

	return LoreRevisionCommitArgsFFI{
		Message:       valMessage,
		Link:          valLink,
		LinkPaths:     valLinkPaths,
		LinkMessages:  valLinkMessages,
		Layer:         valLayer,
		LayerPaths:    valLayerPaths,
		LayerMessages: valLayerMessages,
		Stats:         valStats,
	}, cleanup
}

type LoreRevisionAmendArgs struct {
	/* New commit message */
	Message string
}

type LoreRevisionAmendArgsFFI struct {
	/* New commit message */
	Message LoreString
}

func NewLoreRevisionAmendArgs(opts LoreRevisionAmendArgs) (LoreRevisionAmendArgsFFI, func()) {
	valMessage, cleanupMessage := NewLoreString(opts.Message)

	cleanup := func() {
		cleanupMessage()
	}

	return LoreRevisionAmendArgsFFI{
		Message: valMessage,
	}, cleanup
}

type LoreRevisionInfoArgs struct {
	/* Revision to get info for; empty for current */
	Revision string
	/* Include delta against parent */
	Delta bool
	/* Include file metadata entries */
	Metadata bool
}

type LoreRevisionInfoArgsFFI struct {
	/* Revision to get info for; empty for current */
	Revision LoreString
	/* Include delta against parent */
	Delta uint8
	/* Include file metadata entries */
	Metadata uint8
}

func NewLoreRevisionInfoArgs(opts LoreRevisionInfoArgs) (LoreRevisionInfoArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valDelta, cleanupDelta := Newuint8(opts.Delta)
	valMetadata, cleanupMetadata := Newuint8(opts.Metadata)

	cleanup := func() {
		cleanupRevision()
		cleanupDelta()
		cleanupMetadata()
	}

	return LoreRevisionInfoArgsFFI{
		Revision: valRevision,
		Delta:    valDelta,
		Metadata: valMetadata,
	}, cleanup
}

type LoreRevisionDiffArgs struct {
	/* Source revision to diff from */
	RevisionSource string
	/* Target revision to diff to; empty for current */
	RevisionTarget string
	/* Repository-relative paths to restrict the diff to; empty for all */
	Paths []string
}

type LoreRevisionDiffArgsFFI struct {
	/* Source revision to diff from */
	RevisionSource LoreString
	/* Target revision to diff to; empty for current */
	RevisionTarget LoreString
	/* Repository-relative paths to restrict the diff to; empty for all */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionDiffArgs(opts LoreRevisionDiffArgs) (LoreRevisionDiffArgsFFI, func()) {
	valRevisionSource, cleanupRevisionSource := NewLoreString(opts.RevisionSource)
	valRevisionTarget, cleanupRevisionTarget := NewLoreString(opts.RevisionTarget)
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupRevisionSource()
		cleanupRevisionTarget()
		cleanupPaths()
	}

	return LoreRevisionDiffArgsFFI{
		RevisionSource: valRevisionSource,
		RevisionTarget: valRevisionTarget,
		Paths:          valPaths,
	}, cleanup
}

type LoreRevisionFindArgs struct {
	/* Metadata key to search for; non-empty selects key/value search */
	Key string
	/* Metadata value to match against `key` */
	Value string
	/* Revision number to search for when `key` is empty; 0 disables */
	Number uint64
}

type LoreRevisionFindArgsFFI struct {
	/* Metadata key to search for; non-empty selects key/value search */
	Key LoreString
	/* Metadata value to match against `key` */
	Value LoreString
	/* Revision number to search for when `key` is empty; 0 disables */
	Number uint64
}

func NewLoreRevisionFindArgs(opts LoreRevisionFindArgs) (LoreRevisionFindArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)
	valValue, cleanupValue := NewLoreString(opts.Value)

	cleanup := func() {
		cleanupKey()
		cleanupValue()
	}

	return LoreRevisionFindArgsFFI{
		Key:    valKey,
		Value:  valValue,
		Number: opts.Number,
	}, cleanup
}

type LoreRevisionHistoryArgs struct {
	/* Start from this revision; empty for current */
	Revision string
	/* Restrict to this branch; empty for current */
	Branch string
	/* Stop at revisions created before this date (Unix timestamp; 0 disables) */
	Date uint64
	/* Maximum number of revisions to return; 0 for unlimited */
	Length uint32
	/* Stop when reaching a different branch */
	OnlyBranch bool
}

type LoreRevisionHistoryArgsFFI struct {
	/* Start from this revision; empty for current */
	Revision LoreString
	/* Restrict to this branch; empty for current */
	Branch LoreString
	/* Stop at revisions created before this date (Unix timestamp; 0 disables) */
	Date uint64
	/* Maximum number of revisions to return; 0 for unlimited */
	Length uint32
	/* Stop when reaching a different branch */
	OnlyBranch uint8
}

func NewLoreRevisionHistoryArgs(opts LoreRevisionHistoryArgs) (LoreRevisionHistoryArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valBranch, cleanupBranch := NewLoreString(opts.Branch)
	valOnlyBranch, cleanupOnlyBranch := Newuint8(opts.OnlyBranch)

	cleanup := func() {
		cleanupRevision()
		cleanupBranch()
		cleanupOnlyBranch()
	}

	return LoreRevisionHistoryArgsFFI{
		Revision:   valRevision,
		Branch:     valBranch,
		Date:       opts.Date,
		Length:     opts.Length,
		OnlyBranch: valOnlyBranch,
	}, cleanup
}

type LoreRevisionRestoreArgs struct {
	/* Commit message for the restored revision */
	Message string
}

type LoreRevisionRestoreArgsFFI struct {
	/* Commit message for the restored revision */
	Message LoreString
}

func NewLoreRevisionRestoreArgs(opts LoreRevisionRestoreArgs) (LoreRevisionRestoreArgsFFI, func()) {
	valMessage, cleanupMessage := NewLoreString(opts.Message)

	cleanup := func() {
		cleanupMessage()
	}

	return LoreRevisionRestoreArgsFFI{
		Message: valMessage,
	}, cleanup
}

type LoreRevisionMetadataClearArgs struct {
	Unused int
}

type LoreRevisionMetadataClearArgsFFI struct {
	Unused int
}

func NewLoreRevisionMetadataClearArgs(opts LoreRevisionMetadataClearArgs) (LoreRevisionMetadataClearArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionMetadataClearArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRevisionMetadataGetArgs struct {
	/* Metadata key to look up */
	Key string
	/* Revision to get metadata for; empty for current */
	Revision string
}

type LoreRevisionMetadataGetArgsFFI struct {
	/* Metadata key to look up */
	Key LoreString
	/* Revision to get metadata for; empty for current */
	Revision LoreString
}

func NewLoreRevisionMetadataGetArgs(opts LoreRevisionMetadataGetArgs) (LoreRevisionMetadataGetArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)
	valRevision, cleanupRevision := NewLoreString(opts.Revision)

	cleanup := func() {
		cleanupKey()
		cleanupRevision()
	}

	return LoreRevisionMetadataGetArgsFFI{
		Key:      valKey,
		Revision: valRevision,
	}, cleanup
}

type LoreRevisionMetadataListArgs struct {
	/* Revision to list metadata for; empty for current */
	Revision string
}

type LoreRevisionMetadataListArgsFFI struct {
	/* Revision to list metadata for; empty for current */
	Revision LoreString
}

func NewLoreRevisionMetadataListArgs(opts LoreRevisionMetadataListArgs) (LoreRevisionMetadataListArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)

	cleanup := func() {
		cleanupRevision()
	}

	return LoreRevisionMetadataListArgsFFI{
		Revision: valRevision,
	}, cleanup
}

type LoreRevisionMetadataSetArgs struct {
	/* Metadata keys (parallel with `values` and `formats`) */
	Keys []string
	/* Metadata values, decoded per the matching format */
	Values []string
	/* Value type for each entry */
	Formats LoreMetadataTypeArray
}

type LoreRevisionMetadataSetArgsFFI struct {
	/* Metadata keys (parallel with `values` and `formats`) */
	Keys LoreStringArrayFFI
	/* Metadata values, decoded per the matching format */
	Values LoreStringArrayFFI
	/* Value type for each entry */
	Formats LoreMetadataTypeArrayFFI
}

func NewLoreRevisionMetadataSetArgs(opts LoreRevisionMetadataSetArgs) (LoreRevisionMetadataSetArgsFFI, func()) {
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)
	valValues, cleanupValues := NewLoreStringArray(opts.Values)
	valFormats, cleanupFormats := NewLoreMetadataTypeArray(opts.Formats)

	cleanup := func() {
		cleanupKeys()
		cleanupValues()
		cleanupFormats()
	}

	return LoreRevisionMetadataSetArgsFFI{
		Keys:    valKeys,
		Values:  valValues,
		Formats: valFormats,
	}, cleanup
}

type LoreRevisionSyncArgs struct {
	/* Revision to synchronize to; empty for branch tip */
	Revision string
	/* Fast forward and keep local changes when syncing to a local revision */
	ForwardChanges bool
	/* Reset local modified files to match the incoming revision */
	Reset bool
	/* Root files for dependency-based selective sync */
	RootFiles []string
	/* Tags to filter dependencies by during resolution */
	DependencyTags []string
	/* Follow transitive dependencies recursively */
	DependencyRecursive bool
	/* Maximum dependency traversal depth; 0 means unlimited */
	DependencyDepthLimit uint32
}

type LoreRevisionSyncArgsFFI struct {
	/* Revision to synchronize to; empty for branch tip */
	Revision LoreString
	/* Fast forward and keep local changes when syncing to a local revision */
	ForwardChanges uint8
	/* Reset local modified files to match the incoming revision */
	Reset uint8
	/* Root files for dependency-based selective sync */
	RootFiles LoreStringArrayFFI
	/* Tags to filter dependencies by during resolution */
	DependencyTags LoreStringArrayFFI
	/* Follow transitive dependencies recursively */
	DependencyRecursive uint8
	/* Maximum dependency traversal depth; 0 means unlimited */
	DependencyDepthLimit uint32
}

func NewLoreRevisionSyncArgs(opts LoreRevisionSyncArgs) (LoreRevisionSyncArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valForwardChanges, cleanupForwardChanges := Newuint8(opts.ForwardChanges)
	valReset, cleanupReset := Newuint8(opts.Reset)
	valRootFiles, cleanupRootFiles := NewLoreStringArray(opts.RootFiles)
	valDependencyTags, cleanupDependencyTags := NewLoreStringArray(opts.DependencyTags)
	valDependencyRecursive, cleanupDependencyRecursive := Newuint8(opts.DependencyRecursive)

	cleanup := func() {
		cleanupRevision()
		cleanupForwardChanges()
		cleanupReset()
		cleanupRootFiles()
		cleanupDependencyTags()
		cleanupDependencyRecursive()
	}

	return LoreRevisionSyncArgsFFI{
		Revision:             valRevision,
		ForwardChanges:       valForwardChanges,
		Reset:                valReset,
		RootFiles:            valRootFiles,
		DependencyTags:       valDependencyTags,
		DependencyRecursive:  valDependencyRecursive,
		DependencyDepthLimit: opts.DependencyDepthLimit,
	}, cleanup
}

type LoreRevisionRevertArgs struct {
	/* Revision to revert */
	Revision string
	/* Message to use for an auto-commit if no conflicts arise */
	Message string
	/* Disable auto-commit even if no conflicts arise */
	NoCommit bool
}

type LoreRevisionRevertArgsFFI struct {
	/* Revision to revert */
	Revision LoreString
	/* Message to use for an auto-commit if no conflicts arise */
	Message LoreString
	/* Disable auto-commit even if no conflicts arise */
	NoCommit uint8
}

func NewLoreRevisionRevertArgs(opts LoreRevisionRevertArgs) (LoreRevisionRevertArgsFFI, func()) {
	valRevision, cleanupRevision := NewLoreString(opts.Revision)
	valMessage, cleanupMessage := NewLoreString(opts.Message)
	valNoCommit, cleanupNoCommit := Newuint8(opts.NoCommit)

	cleanup := func() {
		cleanupRevision()
		cleanupMessage()
		cleanupNoCommit()
	}

	return LoreRevisionRevertArgsFFI{
		Revision: valRevision,
		Message:  valMessage,
		NoCommit: valNoCommit,
	}, cleanup
}

type LoreRevisionRevertAbortArgs struct {
	Unused int
}

type LoreRevisionRevertAbortArgsFFI struct {
	Unused int
}

func NewLoreRevisionRevertAbortArgs(opts LoreRevisionRevertAbortArgs) (LoreRevisionRevertAbortArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionRevertAbortArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRevisionRevertUnresolveArgs struct {
	/* Repository-relative paths to mark unresolved */
	Paths []string
}

type LoreRevisionRevertUnresolveArgsFFI struct {
	/* Repository-relative paths to mark unresolved */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionRevertUnresolveArgs(opts LoreRevisionRevertUnresolveArgs) (LoreRevisionRevertUnresolveArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreRevisionRevertUnresolveArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreRevisionRevertRestartArgs struct {
	/* Repository-relative paths to re-materialize for resolution */
	Paths []string
}

type LoreRevisionRevertRestartArgsFFI struct {
	/* Repository-relative paths to re-materialize for resolution */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionRevertRestartArgs(opts LoreRevisionRevertRestartArgs) (LoreRevisionRevertRestartArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreRevisionRevertRestartArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreRevisionRevertResolveArgs struct {
	/* Repository-relative paths to mark resolved */
	Paths []string
}

type LoreRevisionRevertResolveArgsFFI struct {
	/* Repository-relative paths to mark resolved */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionRevertResolveArgs(opts LoreRevisionRevertResolveArgs) (LoreRevisionRevertResolveArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreRevisionRevertResolveArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreRevisionRevertResolveMineArgs struct {
	/* Repository-relative paths to resolve in favor of "mine" */
	Paths []string
}

type LoreRevisionRevertResolveMineArgsFFI struct {
	/* Repository-relative paths to resolve in favor of "mine" */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionRevertResolveMineArgs(opts LoreRevisionRevertResolveMineArgs) (LoreRevisionRevertResolveMineArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreRevisionRevertResolveMineArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreRevisionRevertResolveTheirsArgs struct {
	/* Repository-relative paths to resolve in favor of "theirs" */
	Paths []string
}

type LoreRevisionRevertResolveTheirsArgsFFI struct {
	/* Repository-relative paths to resolve in favor of "theirs" */
	Paths LoreStringArrayFFI
}

func NewLoreRevisionRevertResolveTheirsArgs(opts LoreRevisionRevertResolveTheirsArgs) (LoreRevisionRevertResolveTheirsArgsFFI, func()) {
	valPaths, cleanupPaths := NewLoreStringArray(opts.Paths)

	cleanup := func() {
		cleanupPaths()
	}

	return LoreRevisionRevertResolveTheirsArgsFFI{
		Paths: valPaths,
	}, cleanup
}

type LoreSharedStoreCreateArgs struct {
	/* Remote URL backing the store */
	RemoteUrl string
	/* Path where the store will be created; empty string uses the default location */
	Path string
	/* Set this as the default shared store in the global config */
	MakeDefault bool
}

type LoreSharedStoreCreateArgsFFI struct {
	/* Remote URL backing the store */
	RemoteUrl LoreString
	/* Path where the store will be created; empty string uses the default location */
	Path LoreString
	/* Set this as the default shared store in the global config */
	MakeDefault uint8
}

func NewLoreSharedStoreCreateArgs(opts LoreSharedStoreCreateArgs) (LoreSharedStoreCreateArgsFFI, func()) {
	valRemoteUrl, cleanupRemoteUrl := NewLoreString(opts.RemoteUrl)
	valPath, cleanupPath := NewLoreString(opts.Path)
	valMakeDefault, cleanupMakeDefault := Newuint8(opts.MakeDefault)

	cleanup := func() {
		cleanupRemoteUrl()
		cleanupPath()
		cleanupMakeDefault()
	}

	return LoreSharedStoreCreateArgsFFI{
		RemoteUrl:   valRemoteUrl,
		Path:        valPath,
		MakeDefault: valMakeDefault,
	}, cleanup
}

type LoreSharedStoreInfoArgs struct {
	Unused int
}

type LoreSharedStoreInfoArgsFFI struct {
	Unused int
}

func NewLoreSharedStoreInfoArgs(opts LoreSharedStoreInfoArgs) (LoreSharedStoreInfoArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreSharedStoreInfoArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreSharedStoreSetUseAutomaticallyArgs struct {
	/* Automatically use the shared store */
	Enabled bool
}

type LoreSharedStoreSetUseAutomaticallyArgsFFI struct {
	/* Automatically use the shared store */
	Enabled uint8
}

func NewLoreSharedStoreSetUseAutomaticallyArgs(opts LoreSharedStoreSetUseAutomaticallyArgs) (LoreSharedStoreSetUseAutomaticallyArgsFFI, func()) {
	valEnabled, cleanupEnabled := Newuint8(opts.Enabled)

	cleanup := func() {
		cleanupEnabled()
	}

	return LoreSharedStoreSetUseAutomaticallyArgsFFI{
		Enabled: valEnabled,
	}, cleanup
}

type LoreStorageOpenArgs struct {
	/* Path to an existing lore repository; must be empty when `in_memory` is set */
	RepositoryPath string
	/* Open a fresh in-memory store; `repository_path` must then be empty */
	InMemory bool
	/* Remote endpoint binding for ops that consult a peer; honored only when `has_remote_config` is set */
	RemoteConfig LoreStorageRemoteConfig
	/* Activate `remote_config`; otherwise the handle has no remote */
	HasRemoteConfig bool
	/* Soft cap on total immutable-store bytes (compactor target). A non-zero cache target enables
	incremental background GC for the handle; `0` then selects the default. Shared disk backends
	inherit the first opener's value */
	CacheTargetBytes uint64
	/* Soft cap on immutable-store fragment count (evictor target). A non-zero cache target enables
	incremental background GC for the handle; `0` then selects the default */
	CacheTargetFragments uint64
}

type LoreStorageOpenArgsFFI struct {
	/* Path to an existing lore repository; must be empty when `in_memory` is set */
	RepositoryPath LoreString
	/* Open a fresh in-memory store; `repository_path` must then be empty */
	InMemory uint8
	/* Remote endpoint binding for ops that consult a peer; honored only when `has_remote_config` is set */
	RemoteConfig LoreStorageRemoteConfig
	/* Activate `remote_config`; otherwise the handle has no remote */
	HasRemoteConfig uint8
	/* Soft cap on total immutable-store bytes (compactor target). A non-zero cache target enables
	incremental background GC for the handle; `0` then selects the default. Shared disk backends
	inherit the first opener's value */
	CacheTargetBytes uint64
	/* Soft cap on immutable-store fragment count (evictor target). A non-zero cache target enables
	incremental background GC for the handle; `0` then selects the default */
	CacheTargetFragments uint64
}

func NewLoreStorageOpenArgs(opts LoreStorageOpenArgs) (LoreStorageOpenArgsFFI, func()) {
	valRepositoryPath, cleanupRepositoryPath := NewLoreString(opts.RepositoryPath)
	valInMemory, cleanupInMemory := Newuint8(opts.InMemory)
	valHasRemoteConfig, cleanupHasRemoteConfig := Newuint8(opts.HasRemoteConfig)

	cleanup := func() {
		cleanupRepositoryPath()
		cleanupInMemory()
		cleanupHasRemoteConfig()
	}

	return LoreStorageOpenArgsFFI{
		RepositoryPath:       valRepositoryPath,
		InMemory:             valInMemory,
		RemoteConfig:         opts.RemoteConfig,
		HasRemoteConfig:      valHasRemoteConfig,
		CacheTargetBytes:     opts.CacheTargetBytes,
		CacheTargetFragments: opts.CacheTargetFragments,
	}, cleanup
}

type LoreStoragePutArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Buffers to store; each runs independently and emits its own `PUT_ITEM_COMPLETE` */
	Items LoreStoragePutItemArray
}

type LoreStoragePutArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Buffers to store; each runs independently and emits its own `PUT_ITEM_COMPLETE` */
	Items LoreStoragePutItemArrayFFI
}

func NewLoreStoragePutArgs(opts LoreStoragePutArgs) (LoreStoragePutArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStoragePutItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStoragePutArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageGetArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to read; each runs independently and emits its own event sequence */
	Items LoreStorageGetItemArray
}

type LoreStorageGetArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to read; each runs independently and emits its own event sequence */
	Items LoreStorageGetItemArrayFFI
}

func NewLoreStorageGetArgs(opts LoreStorageGetArgs) (LoreStorageGetArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageGetItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageGetArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageCloseArgs struct {
	/* Handle to release; from `LORE_EVENT_STORAGE_OPENED` */
	Handle LoreStore
}

type LoreStorageCloseArgsFFI struct {
	/* Handle to release; from `LORE_EVENT_STORAGE_OPENED` */
	Handle LoreStore
}

func NewLoreStorageCloseArgs(opts LoreStorageCloseArgs) (LoreStorageCloseArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageCloseArgsFFI{
		Handle: opts.Handle,
	}, cleanup
}

type LoreStorageFlushArgs struct {
	/* Open handle whose pending writes to flush */
	Handle LoreStore
}

type LoreStorageFlushArgsFFI struct {
	/* Open handle whose pending writes to flush */
	Handle LoreStore
}

func NewLoreStorageFlushArgs(opts LoreStorageFlushArgs) (LoreStorageFlushArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreStorageFlushArgsFFI{
		Handle: opts.Handle,
	}, cleanup
}

type LoreStorageGetMetadataArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to look up; each runs independently and emits its own `GET_METADATA_ITEM_COMPLETE` */
	Items LoreStorageGetMetadataItemArray
}

type LoreStorageGetMetadataArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to look up; each runs independently and emits its own `GET_METADATA_ITEM_COMPLETE` */
	Items LoreStorageGetMetadataItemArrayFFI
}

func NewLoreStorageGetMetadataArgs(opts LoreStorageGetMetadataArgs) (LoreStorageGetMetadataArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageGetMetadataItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageGetMetadataArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageObliterateArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to delete; each runs independently and emits its own `OBLITERATE_ITEM_COMPLETE` */
	Items LoreStorageObliterateItemArray
}

type LoreStorageObliterateArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses to delete; each runs independently and emits its own `OBLITERATE_ITEM_COMPLETE` */
	Items LoreStorageObliterateItemArrayFFI
}

func NewLoreStorageObliterateArgs(opts LoreStorageObliterateArgs) (LoreStorageObliterateArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageObliterateItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageObliterateArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageMutableLoadArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Keys to read; each runs independently and emits its own `MUTABLE_LOAD_ITEM_COMPLETE` */
	Items LoreStorageMutableLoadItemArray
}

type LoreStorageMutableLoadArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Keys to read; each runs independently and emits its own `MUTABLE_LOAD_ITEM_COMPLETE` */
	Items LoreStorageMutableLoadItemArrayFFI
}

func NewLoreStorageMutableLoadArgs(opts LoreStorageMutableLoadArgs) (LoreStorageMutableLoadArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageMutableLoadItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageMutableLoadArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageMutableStoreArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Key-value pairs to write; each runs independently and emits its own `MUTABLE_STORE_ITEM_COMPLETE` */
	Items LoreStorageMutableStoreItemArray
}

type LoreStorageMutableStoreArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Key-value pairs to write; each runs independently and emits its own `MUTABLE_STORE_ITEM_COMPLETE` */
	Items LoreStorageMutableStoreItemArrayFFI
}

func NewLoreStorageMutableStoreArgs(opts LoreStorageMutableStoreArgs) (LoreStorageMutableStoreArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageMutableStoreItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageMutableStoreArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageMutableCompareAndSwapArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Swaps to perform; each runs independently and emits its own `MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE` */
	Items LoreStorageMutableCompareAndSwapItemArray
}

type LoreStorageMutableCompareAndSwapArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Swaps to perform; each runs independently and emits its own `MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE` */
	Items LoreStorageMutableCompareAndSwapItemArrayFFI
}

func NewLoreStorageMutableCompareAndSwapArgs(opts LoreStorageMutableCompareAndSwapArgs) (LoreStorageMutableCompareAndSwapArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageMutableCompareAndSwapItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageMutableCompareAndSwapArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageMutableListArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Listings to perform; each runs independently and emits its own entries and terminal event */
	Items LoreStorageMutableListItemArray
}

type LoreStorageMutableListArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Listings to perform; each runs independently and emits its own entries and terminal event */
	Items LoreStorageMutableListItemArrayFFI
}

func NewLoreStorageMutableListArgs(opts LoreStorageMutableListArgs) (LoreStorageMutableListArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageMutableListItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageMutableListArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageCopyArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Copy requests; each runs independently and emits its own `COPY_ITEM_COMPLETE` */
	Items LoreStorageCopyItemArray
}

type LoreStorageCopyArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Copy requests; each runs independently and emits its own `COPY_ITEM_COMPLETE` */
	Items LoreStorageCopyItemArrayFFI
}

func NewLoreStorageCopyArgs(opts LoreStorageCopyArgs) (LoreStorageCopyArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageCopyItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageCopyArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStoragePutFileArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Files to store; each runs independently and emits its own `PUT_ITEM_COMPLETE` */
	Items LoreStoragePutFileItemArray
}

type LoreStoragePutFileArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Files to store; each runs independently and emits its own `PUT_ITEM_COMPLETE` */
	Items LoreStoragePutFileItemArrayFFI
}

func NewLoreStoragePutFileArgs(opts LoreStoragePutFileArgs) (LoreStoragePutFileArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStoragePutFileItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStoragePutFileArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageGetFileArgs struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses and destination paths; each runs independently */
	Items LoreStorageGetFileItemArray
}

type LoreStorageGetFileArgsFFI struct {
	/* Open storage handle */
	Handle LoreStore
	/* Addresses and destination paths; each runs independently */
	Items LoreStorageGetFileItemArrayFFI
}

func NewLoreStorageGetFileArgs(opts LoreStorageGetFileArgs) (LoreStorageGetFileArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageGetFileItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageGetFileArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreStorageUploadArgs struct {
	/* Open storage handle; must have been opened with `remote_config` */
	Handle LoreStore
	/* Addresses to push to remote; each runs independently and emits its own `UPLOAD_ITEM_COMPLETE` */
	Items LoreStorageUploadItemArray
}

type LoreStorageUploadArgsFFI struct {
	/* Open storage handle; must have been opened with `remote_config` */
	Handle LoreStore
	/* Addresses to push to remote; each runs independently and emits its own `UPLOAD_ITEM_COMPLETE` */
	Items LoreStorageUploadItemArrayFFI
}

func NewLoreStorageUploadArgs(opts LoreStorageUploadArgs) (LoreStorageUploadArgsFFI, func()) {
	valItems, cleanupItems := NewLoreStorageUploadItemArray(opts.Items)

	cleanup := func() {
		cleanupItems()
	}

	return LoreStorageUploadArgsFFI{
		Handle: opts.Handle,
		Items:  valItems,
	}, cleanup
}

type LoreServiceStartArgs struct {
	Unused int
}

type LoreServiceStartArgsFFI struct {
	Unused int
}

func NewLoreServiceStartArgs(opts LoreServiceStartArgs) (LoreServiceStartArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreServiceStartArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreServiceStopArgs struct {
	/* Stop all repositories rather than just the current one */
	All bool
}

type LoreServiceStopArgsFFI struct {
	/* Stop all repositories rather than just the current one */
	All uint8
}

func NewLoreServiceStopArgs(opts LoreServiceStopArgs) (LoreServiceStopArgsFFI, func()) {
	valAll, cleanupAll := Newuint8(opts.All)

	cleanup := func() {
		cleanupAll()
	}

	return LoreServiceStopArgsFFI{
		All: valAll,
	}, cleanup
}

type LoreNotificationSubscribeArgs struct {
	Unused int
}

type LoreNotificationSubscribeArgsFFI struct {
	Unused int
}

func NewLoreNotificationSubscribeArgs(opts LoreNotificationSubscribeArgs) (LoreNotificationSubscribeArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreNotificationSubscribeArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreNotificationUnsubscribeArgs struct {
	Unused int
}

type LoreNotificationUnsubscribeArgsFFI struct {
	Unused int
}

func NewLoreNotificationUnsubscribeArgs(opts LoreNotificationUnsubscribeArgs) (LoreNotificationUnsubscribeArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreNotificationUnsubscribeArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryMetadataGetArgs struct {
	/* Metadata key to fetch; empty string lists all entries */
	Key string
}

type LoreRepositoryMetadataGetArgsFFI struct {
	/* Metadata key to fetch; empty string lists all entries */
	Key LoreString
}

func NewLoreRepositoryMetadataGetArgs(opts LoreRepositoryMetadataGetArgs) (LoreRepositoryMetadataGetArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupKey()
	}

	return LoreRepositoryMetadataGetArgsFFI{
		Key: valKey,
	}, cleanup
}

type LoreRepositoryMetadataSetArgs struct {
	/* Metadata keys to set, positionally aligned with `values` and `formats` */
	Keys []string
	/* Values to set, one per key, encoded per the matching `formats` entry */
	Values []string
	/* Value format/type for each key-value pair */
	Formats LoreMetadataTypeArray
}

type LoreRepositoryMetadataSetArgsFFI struct {
	/* Metadata keys to set, positionally aligned with `values` and `formats` */
	Keys LoreStringArrayFFI
	/* Values to set, one per key, encoded per the matching `formats` entry */
	Values LoreStringArrayFFI
	/* Value format/type for each key-value pair */
	Formats LoreMetadataTypeArrayFFI
}

func NewLoreRepositoryMetadataSetArgs(opts LoreRepositoryMetadataSetArgs) (LoreRepositoryMetadataSetArgsFFI, func()) {
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)
	valValues, cleanupValues := NewLoreStringArray(opts.Values)
	valFormats, cleanupFormats := NewLoreMetadataTypeArray(opts.Formats)

	cleanup := func() {
		cleanupKeys()
		cleanupValues()
		cleanupFormats()
	}

	return LoreRepositoryMetadataSetArgsFFI{
		Keys:    valKeys,
		Values:  valValues,
		Formats: valFormats,
	}, cleanup
}

type LoreRepositoryMetadataClearArgs struct {
	/* Keys to clear; empty array clears all user-defined keys */
	Keys []string
}

type LoreRepositoryMetadataClearArgsFFI struct {
	/* Keys to clear; empty array clears all user-defined keys */
	Keys LoreStringArrayFFI
}

func NewLoreRepositoryMetadataClearArgs(opts LoreRepositoryMetadataClearArgs) (LoreRepositoryMetadataClearArgsFFI, func()) {
	valKeys, cleanupKeys := NewLoreStringArray(opts.Keys)

	cleanup := func() {
		cleanupKeys()
	}

	return LoreRepositoryMetadataClearArgsFFI{
		Keys: valKeys,
	}, cleanup
}

type LoreRepositoryInstanceListArgs struct {
	Unused int
}

type LoreRepositoryInstanceListArgsFFI struct {
	Unused int
}

func NewLoreRepositoryInstanceListArgs(opts LoreRepositoryInstanceListArgs) (LoreRepositoryInstanceListArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryInstanceListArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryInstancePruneArgs struct {
	Unused int
}

type LoreRepositoryInstancePruneArgsFFI struct {
	Unused int
}

func NewLoreRepositoryInstancePruneArgs(opts LoreRepositoryInstancePruneArgs) (LoreRepositoryInstancePruneArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryInstancePruneArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryUpdatePathArgs struct {
	Unused int
}

type LoreRepositoryUpdatePathArgsFFI struct {
	Unused int
}

func NewLoreRepositoryUpdatePathArgs(opts LoreRepositoryUpdatePathArgs) (LoreRepositoryUpdatePathArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRepositoryUpdatePathArgsFFI{
		Unused: opts.Unused,
	}, cleanup
}

type LoreRepositoryConfigGetArgs struct {
	/* Config key to read (`remote_url` or `identity`) */
	Key string
}

type LoreRepositoryConfigGetArgsFFI struct {
	/* Config key to read (`remote_url` or `identity`) */
	Key LoreString
}

func NewLoreRepositoryConfigGetArgs(opts LoreRepositoryConfigGetArgs) (LoreRepositoryConfigGetArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupKey()
	}

	return LoreRepositoryConfigGetArgsFFI{
		Key: valKey,
	}, cleanup
}

type LoreRevisionTreeLoadArgs struct {
	/* Open storage handle the revision tree is loaded against */
	Store LoreStore
	/* Repository partition the loaded revision belongs to */
	Repository LorePartition
	/* Revision to open; `0` opens an empty tree for an initial commit */
	RevisionHash LoreHash
}

type LoreRevisionTreeLoadArgsFFI struct {
	/* Open storage handle the revision tree is loaded against */
	Store LoreStore
	/* Repository partition the loaded revision belongs to */
	Repository LorePartition
	/* Revision to open; `0` opens an empty tree for an initial commit */
	RevisionHash LoreHash
}

func NewLoreRevisionTreeLoadArgs(opts LoreRevisionTreeLoadArgs) (LoreRevisionTreeLoadArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeLoadArgsFFI{
		Store:        opts.Store,
		Repository:   opts.Repository,
		RevisionHash: opts.RevisionHash,
	}, cleanup
}

type LoreRevisionTreeCloseArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Revision-tree handle to release */
	Handle LoreRevisionTree
}

type LoreRevisionTreeCloseArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Revision-tree handle to release */
	Handle LoreRevisionTree
}

func NewLoreRevisionTreeCloseArgs(opts LoreRevisionTreeCloseArgs) (LoreRevisionTreeCloseArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeCloseArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
	}, cleanup
}

type LoreRevisionTreeResolvePathArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to resolve against */
	Handle LoreRevisionTree
	/* UTF-8 path relative to the tree root; empty resolves to the root node */
	Path string
}

type LoreRevisionTreeResolvePathArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to resolve against */
	Handle LoreRevisionTree
	/* UTF-8 path relative to the tree root; empty resolves to the root node */
	Path LoreString
}

func NewLoreRevisionTreeResolvePathArgs(opts LoreRevisionTreeResolvePathArgs) (LoreRevisionTreeResolvePathArgsFFI, func()) {
	valPath, cleanupPath := NewLoreString(opts.Path)

	cleanup := func() {
		cleanupPath()
	}

	return LoreRevisionTreeResolvePathArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		Path:   valPath,
	}, cleanup
}

type LoreRevisionTreeListChildrenArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Directory node whose children are streamed */
	ParentNodeId uint32
}

type LoreRevisionTreeListChildrenArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Directory node whose children are streamed */
	ParentNodeId uint32
}

func NewLoreRevisionTreeListChildrenArgs(opts LoreRevisionTreeListChildrenArgs) (LoreRevisionTreeListChildrenArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeListChildrenArgsFFI{
		Id:           opts.Id,
		Handle:       opts.Handle,
		ParentNodeId: opts.ParentNodeId,
	}, cleanup
}

type LoreRevisionTreeNodeInfoArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Node whose record is fetched; the root id also yields `root_info` */
	NodeId uint32
}

type LoreRevisionTreeNodeInfoArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Node whose record is fetched; the root id also yields `root_info` */
	NodeId uint32
}

func NewLoreRevisionTreeNodeInfoArgs(opts LoreRevisionTreeNodeInfoArgs) (LoreRevisionTreeNodeInfoArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeNodeInfoArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		NodeId: opts.NodeId,
	}, cleanup
}

type LoreRevisionTreeNodePathArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Node whose full UTF-8 path is reconstructed by walking parents */
	NodeId uint32
}

type LoreRevisionTreeNodePathArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Node whose full UTF-8 path is reconstructed by walking parents */
	NodeId uint32
}

func NewLoreRevisionTreeNodePathArgs(opts LoreRevisionTreeNodePathArgs) (LoreRevisionTreeNodePathArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeNodePathArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		NodeId: opts.NodeId,
	}, cleanup
}

type LoreRevisionTreeAddArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Parent node the new child is added under */
	ParentNodeId uint32
	/* UTF-8 name of the new child within its parent */
	Name string
	/* `NodeKind` encoding: FILE=1, DIRECTORY=2, LINK=3 */
	Kind uint32
	/* POSIX permission bits for the new node */
	Mode uint16
	/* Content size in bytes (leaf nodes) */
	Size uint64
	/* Content address `(hash, file_id context)` of the new node */
	Address LoreAddress
}

type LoreRevisionTreeAddArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Parent node the new child is added under */
	ParentNodeId uint32
	/* UTF-8 name of the new child within its parent */
	Name LoreString
	/* `NodeKind` encoding: FILE=1, DIRECTORY=2, LINK=3 */
	Kind uint32
	/* POSIX permission bits for the new node */
	Mode uint16
	/* Content size in bytes (leaf nodes) */
	Size uint64
	/* Content address `(hash, file_id context)` of the new node */
	Address LoreAddress
}

func NewLoreRevisionTreeAddArgs(opts LoreRevisionTreeAddArgs) (LoreRevisionTreeAddArgsFFI, func()) {
	valName, cleanupName := NewLoreString(opts.Name)

	cleanup := func() {
		cleanupName()
	}

	return LoreRevisionTreeAddArgsFFI{
		Id:           opts.Id,
		Handle:       opts.Handle,
		ParentNodeId: opts.ParentNodeId,
		Name:         valName,
		Kind:         opts.Kind,
		Mode:         opts.Mode,
		Size:         opts.Size,
		Address:      opts.Address,
	}, cleanup
}

type LoreRevisionTreeDeleteArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Subtree root to mark deleted, including its transitive children */
	NodeId uint32
}

type LoreRevisionTreeDeleteArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Subtree root to mark deleted, including its transitive children */
	NodeId uint32
}

func NewLoreRevisionTreeDeleteArgs(opts LoreRevisionTreeDeleteArgs) (LoreRevisionTreeDeleteArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeDeleteArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		NodeId: opts.NodeId,
	}, cleanup
}

type LoreRevisionTreeModifyArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Leaf node to update; non-leaf targets are rejected */
	NodeId uint32
	/* New POSIX permission bits */
	Mode uint16
	/* New content size in bytes */
	Size uint64
	/* New content address; the existing `file_id` context is preserved */
	Address LoreAddress
}

type LoreRevisionTreeModifyArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Leaf node to update; non-leaf targets are rejected */
	NodeId uint32
	/* New POSIX permission bits */
	Mode uint16
	/* New content size in bytes */
	Size uint64
	/* New content address; the existing `file_id` context is preserved */
	Address LoreAddress
}

func NewLoreRevisionTreeModifyArgs(opts LoreRevisionTreeModifyArgs) (LoreRevisionTreeModifyArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeModifyArgsFFI{
		Id:      opts.Id,
		Handle:  opts.Handle,
		NodeId:  opts.NodeId,
		Mode:    opts.Mode,
		Size:    opts.Size,
		Address: opts.Address,
	}, cleanup
}

type LoreRevisionTreeMoveArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Node to move; its `file_id` is preserved across the move */
	NodeId uint32
	/* Parent node the moved node is reparented under */
	DestinationParentId uint32
	/* UTF-8 name the moved node takes at the destination */
	DstName string
}

type LoreRevisionTreeMoveArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Node to move; its `file_id` is preserved across the move */
	NodeId uint32
	/* Parent node the moved node is reparented under */
	DestinationParentId uint32
	/* UTF-8 name the moved node takes at the destination */
	DstName LoreString
}

func NewLoreRevisionTreeMoveArgs(opts LoreRevisionTreeMoveArgs) (LoreRevisionTreeMoveArgsFFI, func()) {
	valDstName, cleanupDstName := NewLoreString(opts.DstName)

	cleanup := func() {
		cleanupDstName()
	}

	return LoreRevisionTreeMoveArgsFFI{
		Id:                  opts.Id,
		Handle:              opts.Handle,
		NodeId:              opts.NodeId,
		DestinationParentId: opts.DestinationParentId,
		DstName:             valDstName,
	}, cleanup
}

type LoreRevisionTreeMetadataSetArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Metadata key; re-setting it overwrites the pending value */
	Key string
	/* Value stored under the key */
	Value string
	/* Value encoding, matching `LoreRevisionMetadataSetArgs::formats` */
	Format uint32
}

type LoreRevisionTreeMetadataSetArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to mutate */
	Handle LoreRevisionTree
	/* Metadata key; re-setting it overwrites the pending value */
	Key LoreString
	/* Value stored under the key */
	Value LoreString
	/* Value encoding, matching `LoreRevisionMetadataSetArgs::formats` */
	Format uint32
}

func NewLoreRevisionTreeMetadataSetArgs(opts LoreRevisionTreeMetadataSetArgs) (LoreRevisionTreeMetadataSetArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)
	valValue, cleanupValue := NewLoreString(opts.Value)

	cleanup := func() {
		cleanupKey()
		cleanupValue()
	}

	return LoreRevisionTreeMetadataSetArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		Key:    valKey,
		Value:  valValue,
		Format: opts.Format,
	}, cleanup
}

type LoreRevisionTreeMetadataGetArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Metadata key to read; pending edits take precedence over the revision */
	Key string
}

type LoreRevisionTreeMetadataGetArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to read from */
	Handle LoreRevisionTree
	/* Metadata key to read; pending edits take precedence over the revision */
	Key LoreString
}

func NewLoreRevisionTreeMetadataGetArgs(opts LoreRevisionTreeMetadataGetArgs) (LoreRevisionTreeMetadataGetArgsFFI, func()) {
	valKey, cleanupKey := NewLoreString(opts.Key)

	cleanup := func() {
		cleanupKey()
	}

	return LoreRevisionTreeMetadataGetArgsFFI{
		Id:     opts.Id,
		Handle: opts.Handle,
		Key:    valKey,
	}, cleanup
}

type LoreRevisionTreeCommitArgs struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to freeze and commit */
	Handle LoreRevisionTree
	/* Branch whose tip is atomically advanced to the new revision */
	Branch LoreBranchId
	/* Commit tuneables (local-only vs remote-uploading) */
	Options LoreRevisionTreeCommitOptions
}

type LoreRevisionTreeCommitArgsFFI struct {
	/* Per-call correlation id echoed back in events */
	Id uint64
	/* Loaded revision-tree handle to freeze and commit */
	Handle LoreRevisionTree
	/* Branch whose tip is atomically advanced to the new revision */
	Branch LoreBranchId
	/* Commit tuneables (local-only vs remote-uploading) */
	Options LoreRevisionTreeCommitOptions
}

func NewLoreRevisionTreeCommitArgs(opts LoreRevisionTreeCommitArgs) (LoreRevisionTreeCommitArgsFFI, func()) {

	cleanup := func() {
	}

	return LoreRevisionTreeCommitArgsFFI{
		Id:      opts.Id,
		Handle:  opts.Handle,
		Branch:  opts.Branch,
		Options: opts.Options,
	}, cleanup
}
