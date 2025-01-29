package ui

import (
	"testing"
	"timer-cli/internal/task"
)

func TestStartUI(t *testing.T) {
	taskManager := task.NewTaskManager()
	err := Start(taskManager)
	if err != nil {
		t.Errorf("Failed to start UI: %v", err)
	}
}