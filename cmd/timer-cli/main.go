package main

import (
	"fmt"
	"os"
	"timer-cli/internal/config"
	"timer-cli/internal/task"
	"timer-cli/internal/ui"
	"timer-cli/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	utils.InitLogger(cfg.LogLevel)

	// Initialize task manager
	taskManager := task.NewTaskManager()

	// Start the TUI
	if err := ui.Start(taskManager); err != nil {
		utils.Logger.Fatalf("Failed to start UI: %v", err)
	}
}