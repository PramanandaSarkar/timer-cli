package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"timer-cli/internal/task"
)

const saveFile = "tasks.json"

func Start(taskManager *task.TaskManager) error {
	// Load tasks from file on startup
	taskManager.LoadTasksFromFile(saveFile)

	app := tview.NewApplication()

	// Create UI panels
	tasksPanel := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	cmdInput := tview.NewInputField().SetLabel("$ ")

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tasksPanel, 0, 1, false).
		AddItem(cmdInput, 1, 1, true)

	// Update tasks display
	updateTasksPanel := func() {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Current Time: [::b]%s[::-]\n\n", time.Now().Format("15:04:05")))

		for _, t := range taskManager.Tasks {
			// Update remaining time if task is running
			if t.State == "running" {
				t.Remaining = t.Duration - time.Since(t.StartTime)
				if t.Remaining <= 0 {
					t.State = "completed"
					t.Remaining = 0
				}
			}

			status := fmt.Sprintf("%d %s - %s (%s", t.ID, t.Name, t.Remaining.Round(time.Second), t.State)
			if t.Loop {
				status += ", loop"
			}
			status += ")"
			sb.WriteString(status + "\n")
		}

		tasksPanel.SetText(sb.String())
	}

	// Command handling
	cmdInput.SetDoneFunc(func(key tcell.Key) {
		cmd := strings.TrimSpace(cmdInput.GetText())
		cmdInput.SetText("")

		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			return
		}

		switch parts[0] {
		case "add", "+", "a":
			if len(parts) < 3 {
				return
			}
			name := parts[1]
			duration, err := time.ParseDuration(parts[2])
			if err != nil {
				return
			}
			taskManager.AddTask(name, duration)

		case "start", "s", "continue", "c":
			if len(parts) < 2 {
				return
			}
			id := parseID(parts[1])
			taskManager.StartTask(id)

		case "stop", "pause", "p":
			if len(parts) < 2 {
				return
			}
			id := parseID(parts[1])
			taskManager.StopTask(id)

		case "modify", "m":
			if len(parts) < 4 {
				return
			}
			id := parseID(parts[1])
			if parts[2] == "-n" {
				taskManager.ModifyTaskName(id, parts[3])
			} else if parts[2] == "-d" {
				duration, err := time.ParseDuration(parts[3])
				if err != nil {
					return
				}
				taskManager.ModifyTaskDuration(id, duration)
			}

		case "quit", "exit", "q":
			taskManager.SaveTasksToFile(saveFile)
			app.Stop()

		default:
			tasksPanel.SetText("[red]Unknown command![-]")
		}

		updateTasksPanel()
	})

	// Auto-refresh tasks panel every second
	go func() {
		for {
			time.Sleep(1 * time.Second)
			app.QueueUpdateDraw(updateTasksPanel)
		}
	}()

	return app.SetRoot(flex, true).EnableMouse(true).Run()
}

func parseID(s string) int {
	var id int
	fmt.Sscanf(s, "%d", &id)
	return id
}
