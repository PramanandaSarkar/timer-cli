package task

import (
	"sync"
	"time"
)

type Task struct {
	ID         int
	Name       string
	Duration   time.Duration
	Remaining  time.Duration
	State      string // "running", "paused", "completed"
	Loop       bool
	StartTime  time.Time
	PauseTime  time.Time
}

type TaskManager struct {
	Tasks       []*Task
	CurrentTask *Task
	mu          sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		Tasks: make([]*Task, 0),
	}
}

func (tm *TaskManager) AddTask(name string, duration time.Duration) *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task := &Task{
		ID:        len(tm.Tasks) + 1,
		Name:      name,
		Duration:  duration,
		Remaining: duration,
		State:     "paused",
	}
	tm.Tasks = append(tm.Tasks, task)
	return task
}

func (tm *TaskManager) StartTask(id int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.Tasks {
		if task.ID == id {
			if tm.CurrentTask != nil && tm.CurrentTask.State == "running" {
				tm.CurrentTask.State = "paused"
			}
			tm.CurrentTask = task
			task.State = "running"
			task.StartTime = time.Now()
		}
	}
}