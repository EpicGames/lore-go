// Copyright Epic Games, Inc. All Rights Reserved.

package types

type LoreLogLevel uint32

const (
	LoreLogLevel_NONE  LoreLogLevel = 0
	LoreLogLevel_TRACE LoreLogLevel = 1
	LoreLogLevel_DEBUG LoreLogLevel = 2
	LoreLogLevel_INFO  LoreLogLevel = 3
	LoreLogLevel_WARN  LoreLogLevel = 4
	LoreLogLevel_ERROR LoreLogLevel = 5
)

type LoreBranchLocation uint32

const (
	LoreBranchLocation_LOCAL  LoreBranchLocation = 0
	LoreBranchLocation_REMOTE LoreBranchLocation = 1
)

type LoreFileAction uint32

const (
	LoreFileAction_KEEP   LoreFileAction = 0
	LoreFileAction_ADD    LoreFileAction = 1
	LoreFileAction_DELETE LoreFileAction = 2
	LoreFileAction_MOVE   LoreFileAction = 3
	LoreFileAction_COPY   LoreFileAction = 4
)

type LoreNodeType uint32

const (
	LoreNodeType_DIRECTORY LoreNodeType = 0
	LoreNodeType_FILE      LoreNodeType = 1
	LoreNodeType_LINK      LoreNodeType = 2
)

type LoreErrorCode uint32

const (
	LoreErrorCode_NONE              LoreErrorCode = 0
	LoreErrorCode_INVALID_ARGUMENTS LoreErrorCode = 1
	LoreErrorCode_ADDRESS_NOT_FOUND LoreErrorCode = 2
	LoreErrorCode_INTERNAL          LoreErrorCode = 3
	LoreErrorCode_SLOW_DOWN         LoreErrorCode = 4
)

type LoreMetadataType uint32

const (
	LoreMetadataType_BINARY  LoreMetadataType = 0
	LoreMetadataType_NUMERIC LoreMetadataType = 1
	LoreMetadataType_STRING  LoreMetadataType = 2
)

type LoreKeyType uint32

const (
	LoreKeyType_UNTYPED               LoreKeyType = 0
	LoreKeyType_BRANCH_METADATA       LoreKeyType = 1
	LoreKeyType_BRANCH_ID             LoreKeyType = 2
	LoreKeyType_BRANCH_LATEST_POINTER LoreKeyType = 3
	LoreKeyType_REPOSITORY_METADATA   LoreKeyType = 4
	LoreKeyType_REPOSITORY_ID         LoreKeyType = 5
	LoreKeyType_INSTANCE              LoreKeyType = 6
)

type LoreMetadataTag uint32

const (
	LoreMetadataTag_ADDRESS LoreMetadataTag = 0
	LoreMetadataTag_BOOLEAN LoreMetadataTag = 1
	LoreMetadataTag_BINARY  LoreMetadataTag = 2
	LoreMetadataTag_CONTEXT LoreMetadataTag = 3
	LoreMetadataTag_HASH    LoreMetadataTag = 4
	LoreMetadataTag_NUMERIC LoreMetadataTag = 5
	LoreMetadataTag_STRING  LoreMetadataTag = 6
)

type LoreEventTag uint32

const (
	LoreEventTag_PROGRESS                                       LoreEventTag = 0
	LoreEventTag_ERROR                                          LoreEventTag = 1
	LoreEventTag_COMPLETE                                       LoreEventTag = 2
	LoreEventTag_METADATA                                       LoreEventTag = 3
	LoreEventTag_LOG                                            LoreEventTag = 4
	LoreEventTag_END                                            LoreEventTag = 5
	LoreEventTag_MAINTENANCE                                    LoreEventTag = 6
	LoreEventTag_AUTH_URL                                       LoreEventTag = 7
	LoreEventTag_AUTH_USER_INFO                                 LoreEventTag = 8
	LoreEventTag_AUTH_USER_TOKEN                                LoreEventTag = 9
	LoreEventTag_AUTH_IDENTITY                                  LoreEventTag = 10
	LoreEventTag_BRANCH_CREATE                                  LoreEventTag = 11
	LoreEventTag_BRANCH_MULTIPLE_INSTANCE                       LoreEventTag = 12
	LoreEventTag_BRANCH_ARCHIVE                                 LoreEventTag = 13
	LoreEventTag_BRANCH_LIST_BEGIN                              LoreEventTag = 14
	LoreEventTag_BRANCH_LIST_ENTRY                              LoreEventTag = 15
	LoreEventTag_BRANCH_LIST_END                                LoreEventTag = 16
	LoreEventTag_BRANCH_MERGE_ABORT_BEGIN                       LoreEventTag = 17
	LoreEventTag_BRANCH_MERGE_ABORT_END                         LoreEventTag = 18
	LoreEventTag_BRANCH_INFO                                    LoreEventTag = 19
	LoreEventTag_BRANCH_DIFF_BEGIN                              LoreEventTag = 20
	LoreEventTag_BRANCH_DIFF_CHANGE_BEGIN                       LoreEventTag = 21
	LoreEventTag_BRANCH_DIFF_CHANGE                             LoreEventTag = 22
	LoreEventTag_BRANCH_DIFF_CHANGE_END                         LoreEventTag = 23
	LoreEventTag_BRANCH_DIFF_CONFLICT_BEGIN                     LoreEventTag = 24
	LoreEventTag_BRANCH_DIFF_CONFLICT                           LoreEventTag = 25
	LoreEventTag_BRANCH_DIFF_CONFLICT_END                       LoreEventTag = 26
	LoreEventTag_BRANCH_DIFF_END                                LoreEventTag = 27
	LoreEventTag_BRANCH_LATEST_LIST_ENTRY                       LoreEventTag = 28
	LoreEventTag_BRANCH_MERGE_CONFLICT_FILE                     LoreEventTag = 29
	LoreEventTag_BRANCH_MERGE_LINK_SKIPPED                      LoreEventTag = 30
	LoreEventTag_BRANCH_MERGE_UNRESOLVE_FILE                    LoreEventTag = 31
	LoreEventTag_BRANCH_MERGE_UNRESOLVE_REVISION                LoreEventTag = 32
	LoreEventTag_BRANCH_MERGE_INTO_FILE_BEGIN                   LoreEventTag = 33
	LoreEventTag_BRANCH_MERGE_INTO_FILE                         LoreEventTag = 34
	LoreEventTag_BRANCH_MERGE_INTO_FILE_END                     LoreEventTag = 35
	LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_BEGIN               LoreEventTag = 36
	LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_PROGRESS            LoreEventTag = 37
	LoreEventTag_BRANCH_MERGE_INTO_FRAGMENT_END                 LoreEventTag = 38
	LoreEventTag_BRANCH_MERGE_INTO_REVISION                     LoreEventTag = 39
	LoreEventTag_BRANCH_MERGE_INTO_SYNC_BEGIN                   LoreEventTag = 40
	LoreEventTag_BRANCH_MERGE_INTO_SYNC_END                     LoreEventTag = 41
	LoreEventTag_BRANCH_MERGE_RESOLVE_FILE                      LoreEventTag = 42
	LoreEventTag_BRANCH_MERGE_RESOLVE_REVISION                  LoreEventTag = 43
	LoreEventTag_BRANCH_MERGE_START_BEGIN                       LoreEventTag = 44
	LoreEventTag_BRANCH_MERGE_START_END                         LoreEventTag = 45
	LoreEventTag_CHERRY_PICK_START_BEGIN                        LoreEventTag = 46
	LoreEventTag_CHERRY_PICK_START_END                          LoreEventTag = 47
	LoreEventTag_CHERRY_PICK_ABORT_BEGIN                        LoreEventTag = 48
	LoreEventTag_CHERRY_PICK_ABORT_END                          LoreEventTag = 49
	LoreEventTag_CHERRY_PICK_CONFLICT_FILE                      LoreEventTag = 50
	LoreEventTag_CHERRY_PICK_UNRESOLVE_FILE                     LoreEventTag = 51
	LoreEventTag_CHERRY_PICK_UNRESOLVE_REVISION                 LoreEventTag = 52
	LoreEventTag_CHERRY_PICK_RESOLVE_FILE                       LoreEventTag = 53
	LoreEventTag_CHERRY_PICK_RESOLVE_REVISION                   LoreEventTag = 54
	LoreEventTag_REVERT_START_BEGIN                             LoreEventTag = 55
	LoreEventTag_REVERT_START_END                               LoreEventTag = 56
	LoreEventTag_REVERT_ABORT_BEGIN                             LoreEventTag = 57
	LoreEventTag_REVERT_ABORT_END                               LoreEventTag = 58
	LoreEventTag_REVERT_RESOLVE_FILE                            LoreEventTag = 59
	LoreEventTag_REVERT_RESOLVE_REVISION                        LoreEventTag = 60
	LoreEventTag_REVERT_CONFLICT_FILE                           LoreEventTag = 61
	LoreEventTag_REVERT_UNRESOLVE_FILE                          LoreEventTag = 62
	LoreEventTag_REVERT_UNRESOLVE_REVISION                      LoreEventTag = 63
	LoreEventTag_BRANCH_PROTECT                                 LoreEventTag = 64
	LoreEventTag_BRANCH_PUSH                                    LoreEventTag = 65
	LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_BEGIN              LoreEventTag = 66
	LoreEventTag_BRANCH_PUSH_REVISION_UPDATE_END                LoreEventTag = 67
	LoreEventTag_BRANCH_PUSH_FRAGMENT_BEGIN                     LoreEventTag = 68
	LoreEventTag_BRANCH_PUSH_FRAGMENT_PROGRESS                  LoreEventTag = 69
	LoreEventTag_BRANCH_PUSH_FRAGMENT_END                       LoreEventTag = 70
	LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_BEGIN                LoreEventTag = 71
	LoreEventTag_BRANCH_PUSH_BRANCH_CREATE_END                  LoreEventTag = 72
	LoreEventTag_BRANCH_PUSH_REVISION_PUSH_BEGIN                LoreEventTag = 73
	LoreEventTag_BRANCH_PUSH_REVISION_PUSH_UPDATE               LoreEventTag = 74
	LoreEventTag_BRANCH_PUSH_REVISION_PUSH_END                  LoreEventTag = 75
	LoreEventTag_BRANCH_RESET                                   LoreEventTag = 76
	LoreEventTag_BRANCH_SWITCH_BEGIN                            LoreEventTag = 77
	LoreEventTag_BRANCH_SWITCH_END                              LoreEventTag = 78
	LoreEventTag_BRANCH_UNPROTECT                               LoreEventTag = 79
	LoreEventTag_FILE_INFO                                      LoreEventTag = 80
	LoreEventTag_FILE_DIFF                                      LoreEventTag = 81
	LoreEventTag_FILE_HASH                                      LoreEventTag = 82
	LoreEventTag_FILE_HISTORY                                   LoreEventTag = 83
	LoreEventTag_FILE_WRITE                                     LoreEventTag = 84
	LoreEventTag_FILE_OBLITERATE                                LoreEventTag = 85
	LoreEventTag_FILE_DUMP                                      LoreEventTag = 86
	LoreEventTag_FILE_DEPENDENCY_ADD_BEGIN                      LoreEventTag = 87
	LoreEventTag_FILE_DEPENDENCY_ADD_ENTRY                      LoreEventTag = 88
	LoreEventTag_FILE_DEPENDENCY_ADD_END                        LoreEventTag = 89
	LoreEventTag_FILE_DEPENDENCY_REMOVE_BEGIN                   LoreEventTag = 90
	LoreEventTag_FILE_DEPENDENCY_REMOVE_ENTRY                   LoreEventTag = 91
	LoreEventTag_FILE_DEPENDENCY_REMOVE_END                     LoreEventTag = 92
	LoreEventTag_FILE_DEPENDENCY_LIST_BEGIN                     LoreEventTag = 93
	LoreEventTag_FILE_DEPENDENCY_LIST_FILE                      LoreEventTag = 94
	LoreEventTag_FILE_DEPENDENCY_LIST_ENTRY                     LoreEventTag = 95
	LoreEventTag_FILE_DEPENDENCY_LIST_FILE_END                  LoreEventTag = 96
	LoreEventTag_FILE_DEPENDENCY_LIST_END                       LoreEventTag = 97
	LoreEventTag_FILE_RESET_BEGIN                               LoreEventTag = 98
	LoreEventTag_FILE_RESET_PROGRESS                            LoreEventTag = 99
	LoreEventTag_FILE_RESET_END                                 LoreEventTag = 100
	LoreEventTag_FILE_RESET_FILE                                LoreEventTag = 101
	LoreEventTag_FILTER_EXCLUDE                                 LoreEventTag = 102
	LoreEventTag_FILE_STAGE_BEGIN                               LoreEventTag = 103
	LoreEventTag_FILE_STAGE_PROGRESS                            LoreEventTag = 104
	LoreEventTag_FILE_STAGE_END                                 LoreEventTag = 105
	LoreEventTag_FILE_STAGE_REVISION                            LoreEventTag = 106
	LoreEventTag_FILE_STAGE_FILE                                LoreEventTag = 107
	LoreEventTag_FILE_UNSTAGE_BEGIN                             LoreEventTag = 108
	LoreEventTag_FILE_UNSTAGE_PROGRESS                          LoreEventTag = 109
	LoreEventTag_FILE_UNSTAGE_END                               LoreEventTag = 110
	LoreEventTag_FILE_UNSTAGE_REVISION                          LoreEventTag = 111
	LoreEventTag_FILE_UNSTAGE_FILE                              LoreEventTag = 112
	LoreEventTag_FRAGMENT_WRITE                                 LoreEventTag = 113
	LoreEventTag_LAYER_ADD                                      LoreEventTag = 114
	LoreEventTag_LAYER_ENTRY                                    LoreEventTag = 115
	LoreEventTag_LAYER_REMOVE                                   LoreEventTag = 116
	LoreEventTag_LAYER_STAGED_ENTRY                             LoreEventTag = 117
	LoreEventTag_LINK_CHANGE                                    LoreEventTag = 118
	LoreEventTag_LINK_ENTRY                                     LoreEventTag = 119
	LoreEventTag_LOCK_FILE_ACQUIRE_BEGIN                        LoreEventTag = 120
	LoreEventTag_LOCK_FILE_ACQUIRE                              LoreEventTag = 121
	LoreEventTag_LOCK_FILE_STATUS_BEGIN                         LoreEventTag = 122
	LoreEventTag_LOCK_FILE_STATUS                               LoreEventTag = 123
	LoreEventTag_LOCK_FILE_QUERY_BEGIN                          LoreEventTag = 124
	LoreEventTag_LOCK_FILE_QUERY                                LoreEventTag = 125
	LoreEventTag_LOCK_FILE_RELEASE_BEGIN                        LoreEventTag = 126
	LoreEventTag_LOCK_FILE_RELEASE                              LoreEventTag = 127
	LoreEventTag_METADATA_CLEAR_FILE                            LoreEventTag = 128
	LoreEventTag_METADATA_CLEAR_REVISION                        LoreEventTag = 129
	LoreEventTag_PATH_IGNORE                                    LoreEventTag = 130
	LoreEventTag_REPOSITORY_CREATE                              LoreEventTag = 131
	LoreEventTag_REPOSITORY_CLONE_BEGIN                         LoreEventTag = 132
	LoreEventTag_REPOSITORY_CLONE_PROGRESS                      LoreEventTag = 133
	LoreEventTag_REPOSITORY_CLONE_END                           LoreEventTag = 134
	LoreEventTag_DEPENDENCY_RESOLVE_BEGIN                       LoreEventTag = 135
	LoreEventTag_DEPENDENCY_RESOLVE_ITEM                        LoreEventTag = 136
	LoreEventTag_DEPENDENCY_RESOLVE_END                         LoreEventTag = 137
	LoreEventTag_REPOSITORY_DATA                                LoreEventTag = 138
	LoreEventTag_REPOSITORY_CONFIG_GET                          LoreEventTag = 139
	LoreEventTag_REPOSITORY_DUMP_BEGIN                          LoreEventTag = 140
	LoreEventTag_REPOSITORY_DUMP_END                            LoreEventTag = 141
	LoreEventTag_REPOSITORY_LIST_ENTRY                          LoreEventTag = 142
	LoreEventTag_REPOSITORY_INSTANCE                            LoreEventTag = 143
	LoreEventTag_REPOSITORY_VERIFY_STATE_BEGIN                  LoreEventTag = 144
	LoreEventTag_REPOSITORY_VERIFY_STATE_END                    LoreEventTag = 145
	LoreEventTag_REPOSITORY_VERIFY_FRAGMENT                     LoreEventTag = 146
	LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_MATCH               LoreEventTag = 147
	LoreEventTag_REPOSITORY_VERIFY_FRAGMENT_REMOTE              LoreEventTag = 148
	LoreEventTag_REPOSITORY_STATE_DUMP                          LoreEventTag = 149
	LoreEventTag_REPOSITORY_STATE_DUMP_NODE                     LoreEventTag = 150
	LoreEventTag_REPOSITORY_STATUS_REVISION                     LoreEventTag = 151
	LoreEventTag_REPOSITORY_STATUS_FILE                         LoreEventTag = 152
	LoreEventTag_REPOSITORY_STATUS_COUNT                        LoreEventTag = 153
	LoreEventTag_REPOSITORY_STATUS_SUMMARY                      LoreEventTag = 154
	LoreEventTag_REPOSITORY_STORE_IMMUTABLE_QUERY               LoreEventTag = 155
	LoreEventTag_REVISION_COMMIT_BEGIN                          LoreEventTag = 156
	LoreEventTag_REVISION_COMMIT_PROGRESS                       LoreEventTag = 157
	LoreEventTag_REVISION_COMMIT_END                            LoreEventTag = 158
	LoreEventTag_REVISION_COMMIT_REVISION                       LoreEventTag = 159
	LoreEventTag_REVISION_INFO                                  LoreEventTag = 160
	LoreEventTag_REVISION_INFO_DELTA                            LoreEventTag = 161
	LoreEventTag_REVISION_DIFF_FILE                             LoreEventTag = 162
	LoreEventTag_REVISION_FIND                                  LoreEventTag = 163
	LoreEventTag_REVISION_HISTORY                               LoreEventTag = 164
	LoreEventTag_REVISION_HISTORY_ENTRY                         LoreEventTag = 165
	LoreEventTag_REVISION_RESTORE_FILE_BEGIN                    LoreEventTag = 166
	LoreEventTag_REVISION_RESTORE_FILE                          LoreEventTag = 167
	LoreEventTag_REVISION_RESTORE_FILE_END                      LoreEventTag = 168
	LoreEventTag_REVISION_RESTORE_FRAGMENT_BEGIN                LoreEventTag = 169
	LoreEventTag_REVISION_RESTORE_FRAGMENT_PROGRESS             LoreEventTag = 170
	LoreEventTag_REVISION_RESTORE_FRAGMENT_END                  LoreEventTag = 171
	LoreEventTag_REVISION_RESTORE_REVISION                      LoreEventTag = 172
	LoreEventTag_REVISION_RESTORE_SYNC_BEGIN                    LoreEventTag = 173
	LoreEventTag_REVISION_RESTORE_SYNC_END                      LoreEventTag = 174
	LoreEventTag_REVISION_RESOLVE                               LoreEventTag = 175
	LoreEventTag_REVISION_SYNC_TARGET                           LoreEventTag = 176
	LoreEventTag_REVISION_SYNC_FILE                             LoreEventTag = 177
	LoreEventTag_REVISION_SYNC_PROGRESS                         LoreEventTag = 178
	LoreEventTag_REVISION_SYNC_REVISION                         LoreEventTag = 179
	LoreEventTag_REVISION_BISECT                                LoreEventTag = 180
	LoreEventTag_NOTIFICATION_BRANCH_CREATED                    LoreEventTag = 181
	LoreEventTag_NOTIFICATION_BRANCH_DELETED                    LoreEventTag = 182
	LoreEventTag_NOTIFICATION_BRANCH_PUSHED                     LoreEventTag = 183
	LoreEventTag_NOTIFICATION_RESOURCE_LOCKED                   LoreEventTag = 184
	LoreEventTag_NOTIFICATION_RESOURCE_UNLOCKED                 LoreEventTag = 185
	LoreEventTag_NOTIFICATION_SUBSCRIBED                        LoreEventTag = 186
	LoreEventTag_NOTIFICATION_UNSUBSCRIBED                      LoreEventTag = 187
	LoreEventTag_SHARED_STORE_CREATE                            LoreEventTag = 188
	LoreEventTag_SHARED_STORE_INFO                              LoreEventTag = 189
	LoreEventTag_LINK_STAGED_ENTRY                              LoreEventTag = 190
	LoreEventTag_STORAGE_OPENED                                 LoreEventTag = 191
	LoreEventTag_STORAGE_PUT_ITEM_COMPLETE                      LoreEventTag = 192
	LoreEventTag_STORAGE_GET_HEADER                             LoreEventTag = 193
	LoreEventTag_STORAGE_GET_DATA                               LoreEventTag = 194
	LoreEventTag_STORAGE_GET_ITEM_COMPLETE                      LoreEventTag = 195
	LoreEventTag_STORAGE_GET_METADATA_ITEM_COMPLETE             LoreEventTag = 196
	LoreEventTag_STORAGE_COPY_ITEM_COMPLETE                     LoreEventTag = 197
	LoreEventTag_STORAGE_OBLITERATE_ITEM_COMPLETE               LoreEventTag = 198
	LoreEventTag_STORAGE_UPLOAD_ITEM_COMPLETE                   LoreEventTag = 199
	LoreEventTag_REVISION_TREE_LOADED                           LoreEventTag = 200
	LoreEventTag_REVISION_TREE_RESOLVE_PATH_COMPLETE            LoreEventTag = 201
	LoreEventTag_REVISION_TREE_CHILD                            LoreEventTag = 202
	LoreEventTag_REVISION_TREE_NODE_INFO                        LoreEventTag = 203
	LoreEventTag_REVISION_TREE_NODE_PATH                        LoreEventTag = 204
	LoreEventTag_REVISION_TREE_ADD_COMPLETE                     LoreEventTag = 205
	LoreEventTag_REVISION_TREE_DELETE_COMPLETE                  LoreEventTag = 206
	LoreEventTag_REVISION_TREE_MODIFY_COMPLETE                  LoreEventTag = 207
	LoreEventTag_REVISION_TREE_MOVE_COMPLETE                    LoreEventTag = 208
	LoreEventTag_REVISION_TREE_METADATA_SET_COMPLETE            LoreEventTag = 209
	LoreEventTag_REVISION_TREE_METADATA_GET_COMPLETE            LoreEventTag = 210
	LoreEventTag_REVISION_TREE_COMMIT_COMPLETE                  LoreEventTag = 211
	LoreEventTag_REVISION_TREE_CLOSE_COMPLETE                   LoreEventTag = 212
	LoreEventTag_STORAGE_MUTABLE_LOAD_ITEM_COMPLETE             LoreEventTag = 213
	LoreEventTag_STORAGE_MUTABLE_STORE_ITEM_COMPLETE            LoreEventTag = 214
	LoreEventTag_STORAGE_MUTABLE_COMPARE_AND_SWAP_ITEM_COMPLETE LoreEventTag = 215
	LoreEventTag_STORAGE_MUTABLE_LIST_ENTRY                     LoreEventTag = 216
	LoreEventTag_STORAGE_MUTABLE_LIST_ITEM_COMPLETE             LoreEventTag = 217
	LoreEventTag_EVICTION_BEGIN                                 LoreEventTag = 218
	LoreEventTag_EVICTION_PROGRESS                              LoreEventTag = 219
	LoreEventTag_EVICTION_END                                   LoreEventTag = 220
	LoreEventTag_COMPACTION_BEGIN                               LoreEventTag = 221
	LoreEventTag_COMPACTION_PROGRESS                            LoreEventTag = 222
	LoreEventTag_COMPACTION_END                                 LoreEventTag = 223
)
