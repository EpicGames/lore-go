// Copyright Epic Games, Inc. All Rights Reserved.

//go:generate go run github.com/EpicGames/lore-go/cmd/fetch-lore-lib

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	lore "github.com/EpicGames/lore-go/native"
	"github.com/EpicGames/lore-go/types"
)

// Configuration
const (
	COMMIT_MESSAGE = "Initial commit"
	LOG_FILE_PATH  = "./LoreRepositories"
)

// generateID generates a random hex string for unique repository names
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var (
	REPOSITORY_NAME = "EpicRepo" + generateID()
	REPOSITORY_PATH = filepath.Join("./LoreRepositories", REPOSITORY_NAME)
)

// eventHandler handles callback events
func eventHandler(event *types.LoreEventFFI, userContext uint64) {
	if event.Tag == types.LoreEventTag_LOG {
		if logEvent, ok := event.GetData().(*types.LoreLogEventDataFFI); ok {
			if logEvent.Level > types.LoreLogLevel_DEBUG {
				fmt.Println(logEvent.Message)
			}
		}
	}
}

// createFiles generates files to commit to repository
func createFiles() error {
	files := []string{
		filepath.Join(REPOSITORY_PATH, "file.txt"),
		filepath.Join(REPOSITORY_PATH, "log.txt"),
	}

	content := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et"

	for _, file := range files {
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// verifyResult checks the result and exits on failure
func verifyResult(operationName string, result int32, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if result != 0 {
		fmt.Printf("Lore %s failed.\n", operationName)
		os.Exit(1)
	}
	fmt.Printf("Lore %s success.\n", operationName)
}

func main() {
	fmt.Println("Lore Vanilla Example (purego)")
	fmt.Println("============================")

	// If a remote URL is provided as the first CLI arg, run in online mode
	// (push the revision and clone the repository back). Otherwise run a
	// fully offline example that only creates a local repository and commits
	// a file. Authentication is not handled by this example; if the remote
	// requires it, run `lore auth` before invoking this program.
	online := len(os.Args) > 1
	remoteUrl := ""
	if online {
		remoteUrl = os.Args[1]
		fmt.Printf("Running in online mode against: %s\n", remoteUrl)
	} else {
		fmt.Println("Running in offline mode (pass a remote URL as the first arg to enable push/clone)")
	}

	repositoryUrl := REPOSITORY_NAME
	if online {
		repositoryUrl = remoteUrl + "/" + REPOSITORY_NAME
	}

	// Set up callback config
	callback := types.LoreEventCallbackConfig{
		Callback:    eventHandler,
		UserContext: 0,
	}

	// Configure logging
	logConfig, cleanupLogConfig := types.NewLoreLogConfig(types.LoreLogConfig{
		File:     true,
		FilePath: LOG_FILE_PATH,
		Level:    types.LoreLogLevel_DEBUG,
	})
	defer cleanupLogConfig()
	result, err := lore.LogConfigure(&logConfig)
	verifyResult("Setup", result, err)

	// Set up global args
	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: REPOSITORY_PATH,
		Offline:        !online,
	})
	defer cleanupGlobals()

	// Create repository
	{
		repoArgs, cleanupRepo := types.NewLoreRepositoryCreateArgs(types.LoreRepositoryCreateArgs{
			RepositoryUrl: repositoryUrl,
		})
		defer cleanupRepo()

		result, err := lore.RepositoryCreate(&globals, &repoArgs, &callback)
		verifyResult("Repo Create", result, err)
	}

	// Create files to commit to the new repository
	if err := createFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create files: %v\n", err)
		os.Exit(1)
	}

	// Stage files
	{
		paths := []string{
			filepath.Join(REPOSITORY_PATH, "file.txt"),
			filepath.Join(REPOSITORY_PATH, "log.txt"),
		}
		stageArgs, cleanupStage := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
			Paths: paths,
		})
		defer cleanupStage()

		result, err := lore.FileStage(&globals, &stageArgs, &callback)
		verifyResult("File Stage", result, err)
	}

	// Revision commit
	{
		commitArgs, cleanupCommit := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
			Message: COMMIT_MESSAGE,
		})
		defer cleanupCommit()

		result, err := lore.RevisionCommit(&globals, &commitArgs, &callback)
		verifyResult("Revision Commit", result, err)
	}

	if online {
		// Branch push
		{
			pushArgs, cleanupPush := types.NewLoreBranchPushArgs(types.LoreBranchPushArgs{})
			defer cleanupPush()

			result, err := lore.BranchPush(&globals, &pushArgs, &callback)
			verifyResult("Branch Push", result, err)
		}

		// Clone repository
		{
			clonePath := REPOSITORY_PATH + "_clone"
			globalsClone, cleanupGlobalsClone := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
				RepositoryPath: clonePath,
			})
			defer cleanupGlobalsClone()

			cloneArgs, cleanupClone := types.NewLoreRepositoryCloneArgs(types.LoreRepositoryCloneArgs{
				RepositoryUrl: repositoryUrl,
			})
			defer cleanupClone()

			result, err := lore.RepositoryClone(&globalsClone, &cloneArgs, &callback)
			verifyResult("Repository Clone", result, err)
		}
	}

	// Shut down the library
	result, err = lore.Shutdown()
	verifyResult("Shutdown", result, err)
}
