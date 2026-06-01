// Copyright Epic Games, Inc. All Rights Reserved.

//go:generate go run github.com/EpicGames/lore-go/cmd/fetch-lore-lib

package main

import (
	"fmt"
	"os"

	lore "github.com/EpicGames/lore-go/native"
	"github.com/EpicGames/lore-go/types"
)

func main() {
	fmt.Println("Lore Repository Status Example (purego)")
	fmt.Println("========================================")

	// Get repository path from command line or use current directory
	repoPath := "."
	if len(os.Args) > 1 {
		repoPath = os.Args[1]
	}

	fmt.Printf("Repository path: %s\n\n", repoPath)

	// Set up global args
	globals, cleanupGlobals := types.NewLoreGlobalArgs(types.LoreGlobalArgs{
		RepositoryPath: repoPath,
		Offline:        true,
	})
	defer cleanupGlobals()

	// Set up repository status args
	args, cleanupArgs := types.NewLoreRepositoryStatusArgs(types.LoreRepositoryStatusArgs{
		Staged:    true,
		Scan:      true,
		SyncPoint: true,
	})
	defer cleanupArgs()

	// Event counter
	eventCount := 0

	// array for cloned events
	var clonedEvents []types.LoreEvent

	// Define callback that will receive events
	callback := func(event *types.LoreEventFFI, userContext uint64) {
		eventCount++
		eventType := event.Tag
		fmt.Printf("Event #%d: Type=%d, UserContext=%d\n", eventCount, eventType, userContext)

		// Access event data based on type
		if eventType == types.LoreEventTag_LOG {
			// Get a zero-copy view of the log event, narrow the type with type assertion
			if logEvent, ok := event.GetData().(*types.LoreLogEventDataFFI); ok {
				fmt.Printf("  [%s] %s\n", logEvent.Location, logEvent.Message)
			}
		} else if eventType == types.LoreEventTag_REPOSITORY_STATUS_REVISION {
			if revisionEvent, ok := event.GetData().(*types.LoreRepositoryStatusRevisionEventDataFFI); ok {
				fmt.Printf(" Branch: %s\n Revision: %s\n", revisionEvent.BranchName, revisionEvent.Revision)
			}
		} else if eventType == types.LoreEventTag_REPOSITORY_STATUS_FILE {
			// Clone the event for later use (FFI memory is only valid during callback)
			clonedEvent := event.Clone()
			clonedEvents = append(clonedEvents, clonedEvent)
			fmt.Printf("  Cloned REPOSITORY_STATUS_FILE event (total: %d)\n", len(clonedEvents))
		}

		if eventType == types.LoreEventTag_END {
			fmt.Println("\nReceived END event - callback will be cleaned up")
		}
	}

	// Configure the callback with a user context value
	callbackConfig := types.LoreEventCallbackConfig{
		Callback:    callback,
		UserContext: 12345, // Example user context - can be any uint64 value
	}

	// Call the repository status function
	fmt.Println("Calling lore_repository_status...")
	result, err := lore.RepositoryStatus(&globals, &args, &callbackConfig)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nResult code: %d\n", result)
	fmt.Printf("Total events received: %d\n", eventCount)
	fmt.Printf("Cloned REPOSITORY_STATUS_FILE events: %d\n", len(clonedEvents))

	// Demonstrate that cloned events are still accessible after the callback
	if len(clonedEvents) > 0 {
		fmt.Println("\nCloned events are still accessible:")
		for i, clonedEvent := range clonedEvents {
			fmt.Printf("  Event #%d: Tag=%d\n", i+1, clonedEvent.Tag)
			if clonedEvent, ok := clonedEvent.Data.(types.LoreRepositoryStatusFileEventData); ok {
				fmt.Printf("    %s: %d\n", clonedEvent.Path, clonedEvent.Action)
			}
		}
	}

	// Print Lore version
	version, err := lore.Version()
	fmt.Printf("Lore version: %s\n", version)

	// Shut down
	lore.Shutdown()
}
