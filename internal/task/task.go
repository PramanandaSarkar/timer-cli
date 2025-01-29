package task

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Task struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	Remaining time.Duration `json:"remaining"`
	State     string        `json:"state"` // "running", "paused", "completed"
	Loop      bool          `json:"loop"`
	StartTime time.Time     `json:"start_time"`
	PauseTime time.Time     `json:"pause_time"`
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

func (tm *TaskManager) StopTask(id int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.Tasks {
		if task.ID == id && task.State == "running" {
			task.State = "paused"
			task.PauseTime = time.Now()
		}
	}
}

func (tm *TaskManager) ModifyTaskName(id int, newName string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.Tasks {
		if task.ID == id {
			task.Name = newName
			break
		}
	}
}

func (tm *TaskManager) ModifyTaskDuration(id int, newDuration time.Duration) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for _, task := range tm.Tasks {
		if task.ID == id {
			task.Duration = newDuration
			task.Remaining = newDuration
			break
		}
	}
}

func (tm *TaskManager) SaveTasksToFile(filename string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error saving tasks:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tm.Tasks)
	if err != nil {
		fmt.Println("Error encoding tasks:", err)
	}
}

func (tm *TaskManager) LoadTasksFromFile(filename string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return // No saved tasks
		}
		fmt.Println("Error loading tasks:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tm.Tasks)
	if err != nil {
		fmt.Println("Error decoding tasks:", err)
	}
}

func (tm *TaskManager) Quit() {
	tm.SaveTasksToFile("tasks.json")
}
