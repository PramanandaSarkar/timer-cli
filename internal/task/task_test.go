package task

import (
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	tm := NewTaskManager()
	task := tm.AddTask("test", 30*time.Minute)

	if task.Name != "test" || task.Duration != 30*time.Minute {
		t.Errorf("AddTask failed: unexpected task details")
	}
}