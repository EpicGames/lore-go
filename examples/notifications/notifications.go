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
	COMMIT_MESSAGE = "Test commit for notifications"
	LOG_FILE_PATH  = "./LoreRepositories"
)

// generateID generates a random hex string for unique repository names
func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

var (
	REPOSITORY_NAME = "NotificationRepo" + generateID()
	REPOSITORY_PATH = filepath.Join("./LoreRepositories", REPOSITORY_NAME)
)

// Track notification events
var notificationReceived bool

// eventHandler handles general callback events
func eventHandler(event *types.LoreEventFFI, userContext uint64) {
	if event.Tag == types.LoreEventTag_LOG {
		if logEvent, ok := event.GetData().(*types.LoreLogEventDataFFI); ok {
			if logEvent.Level > types.LoreLogLevel_DEBUG {
				fmt.Println(logEvent.Message)
			}
		}
	}
}

// notificationHandler handles notification-specific events
func notificationHandler(event *types.LoreEventFFI, userContext uint64) {
	switch event.Tag {
	case types.LoreEventTag_NOTIFICATION_SUBSCRIBED:
		fmt.Println("🔔 NOTIFICATION: Notification subscription confirmed")

	case types.LoreEventTag_NOTIFICATION_UNSUBSCRIBED:
		fmt.Println("🔔 NOTIFICATION: Notification unsubscription confirmed")

	case types.LoreEventTag_NOTIFICATION_BRANCH_PUSHED:
		notificationReceived = true
		if pushEvent, ok := event.GetData().(*types.LoreNotificationBranchPushedEventDataFFI); ok {
			branchContext := pushEvent.Branch.String()
			userId := pushEvent.UserId.String()
			revisionHash := pushEvent.Revision.String()
			fmt.Printf("🔔 NOTIFICATION: Branch pushed!\n")
			fmt.Printf("   Branch Context: %s\n", branchContext)
			fmt.Printf("   Revision: %s\n", revisionHash[:16])
			fmt.Printf("   User ID: %s\n", userId)
		}

	case types.LoreEventTag_NOTIFICATION_BRANCH_CREATED:
		if createEvent, ok := event.GetData().(*types.LoreNotificationBranchCreatedEventDataFFI); ok {
			branchContext := createEvent.Branch.String()
			fmt.Printf("🔔 NOTIFICATION: Branch created (context: %s)\n", branchContext)
		}

	case types.LoreEventTag_NOTIFICATION_BRANCH_DELETED:
		if deleteEvent, ok := event.GetData().(*types.LoreNotificationBranchDeletedEventDataFFI); ok {
			branchContext := deleteEvent.Branch.String()
			fmt.Printf("🔔 NOTIFICATION: Branch deleted (context: %s)\n", branchContext)
		}

	case types.LoreEventTag_NOTIFICATION_RESOURCE_LOCKED:
		fmt.Println("🔔 NOTIFICATION: Resource locked")

	case types.LoreEventTag_NOTIFICATION_RESOURCE_UNLOCKED:
		fmt.Println("🔔 NOTIFICATION: Resource unlocked")

	case types.LoreEventTag_LOG:
		if logEvent, ok := event.GetData().(*types.LoreLogEventDataFFI); ok {
			if logEvent.Level > types.LoreLogLevel_DEBUG {
				fmt.Println(logEvent.Message)
			}
		}
	}
}

// createFile generates a test file to commit
func createFile() error {
	filePath := filepath.Join(REPOSITORY_PATH, "notification-test.txt")
	content := "This file was created to test notification events" + generateID()

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
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
	fmt.Printf("✓ Lore %s success\n", operationName)
}

func main() {
	fmt.Println("Lore Notifications Example")
	fmt.Println("=========================")
	fmt.Println()

	// If a remote URL is provided as the first CLI arg, run in online mode
	// (push the repository and subscribe to notifications). Otherwise run a
	// fully offline example that only creates a local repository and commits
	// a file — notifications require a remote server. Authentication is not
	// handled by this example; if the remote requires it, run `lore auth`
	// before invoking this program.
	online := len(os.Args) > 1
	remoteUrl := ""
	if online {
		remoteUrl = os.Args[1]
		fmt.Printf("Running in online mode against: %s\n", remoteUrl)
	} else {
		fmt.Println("Running in offline mode (pass a remote URL as the first arg to enable notifications)")
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
	fmt.Println("\nCreating repository...")
	{
		repoArgs, cleanupRepo := types.NewLoreRepositoryCreateArgs(types.LoreRepositoryCreateArgs{
			RepositoryUrl: repositoryUrl,
		})
		defer cleanupRepo()

		result, err := lore.RepositoryCreate(&globals, &repoArgs, &callback)
		verifyResult("Repo Create", result, err)
	}

	// Create a test file
	fmt.Println("\nCreating test file...")
	if err := createFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Test file created")

	// Stage the file
	fmt.Println("\nStaging file...")
	{
		paths := []string{
			filepath.Join(REPOSITORY_PATH, "notification-test.txt"),
		}
		stageArgs, cleanupStage := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
			Paths: paths,
		})
		defer cleanupStage()

		result, err := lore.FileStage(&globals, &stageArgs, &callback)
		verifyResult("File Stage", result, err)
	}

	// Commit the changes
	fmt.Println("\nCommitting changes...")
	{
		commitArgs, cleanupCommit := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
			Message: COMMIT_MESSAGE,
		})
		defer cleanupCommit()

		result, err := lore.RevisionCommit(&globals, &commitArgs, &callback)
		verifyResult("Revision Commit", result, err)
	}

	if online {
		// Push the repository. The server needs to know the repository before notification
		// subscription works
		fmt.Println("\nPushing the repository...")
		{
			pushArgs, cleanupPush := types.NewLoreBranchPushArgs(types.LoreBranchPushArgs{})
			defer cleanupPush()

			result, err := lore.BranchPush(&globals, &pushArgs, &callback)
			verifyResult("Branch Push", result, err)
		}

		// Subscribe to notifications
		fmt.Println("\nSubscribing to notifications...")
		{
			notificationCallback := types.LoreEventCallbackConfig{
				Callback:    notificationHandler,
				UserContext: 0,
			}

			subscribeArgs, cleanupSubscribe := types.NewLoreNotificationSubscribeArgs(types.LoreNotificationSubscribeArgs{})
			defer cleanupSubscribe()

			result, err := lore.NotificationSubscribe(&globals, &subscribeArgs, &notificationCallback)
			verifyResult("Notification Subscribe", result, err)
		}

		// Create a test file
		fmt.Println("\nCreating test file...")
		if err := createFile(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Test file created")

		// Stage the file
		fmt.Println("\nStaging file...")
		{
			paths := []string{
				filepath.Join(REPOSITORY_PATH, "notification-test.txt"),
			}
			stageArgs, cleanupStage := types.NewLoreFileStageArgs(types.LoreFileStageArgs{
				Paths: paths,
			})
			defer cleanupStage()

			result, err := lore.FileStage(&globals, &stageArgs, &callback)
			verifyResult("File Stage", result, err)
		}

		// Commit the changes
		fmt.Println("\nCommitting changes...")
		{
			commitArgs, cleanupCommit := types.NewLoreRevisionCommitArgs(types.LoreRevisionCommitArgs{
				Message: COMMIT_MESSAGE,
			})
			defer cleanupCommit()

			result, err := lore.RevisionCommit(&globals, &commitArgs, &callback)
			verifyResult("Revision Commit", result, err)
		}

		// Push to trigger notification
		fmt.Println("\nPushing changes (this should trigger a notification)...")
		{
			pushArgs, cleanupPush := types.NewLoreBranchPushArgs(types.LoreBranchPushArgs{})
			defer cleanupPush()

			result, err := lore.BranchPush(&globals, &pushArgs, &callback)
			verifyResult("Branch Push", result, err)
		}

		// Verify notification was received
		if notificationReceived {
			fmt.Println("\n✅ SUCCESS: Notification event was received!")
		} else {
			fmt.Println("\n⚠️  WARNING: No notification event was received (this may be expected depending on server configuration)")
		}

		// Unsubscribe from notifications
		fmt.Println("\nUnsubscribing from notifications...")
		{
			unsubscribeArgs, cleanupUnsubscribe := types.NewLoreNotificationUnsubscribeArgs(types.LoreNotificationUnsubscribeArgs{})
			defer cleanupUnsubscribe()

			result, err := lore.NotificationUnsubscribe(&globals, &unsubscribeArgs, &callback)
			verifyResult("Notification Unsubscribe", result, err)
		}
	} else {
		fmt.Println("\nSkipping push and notification subscription (offline mode).")
	}

	// Shut down the library
	fmt.Println("\nShutting down...")
	result, err = lore.Shutdown()
	verifyResult("Shutdown", result, err)

	fmt.Println("\n=========================")
	fmt.Println("Example completed successfully!")
}
