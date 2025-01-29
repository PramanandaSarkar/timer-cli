package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"timer-cli/internal/task"
	"strings"
	"fmt"
	"time"
)

func Start(taskManager *task.TaskManager) error {
	app := tview.NewApplication()

	// Create panels
	tasksPanel := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	cmdInput := tview.NewInputField().SetLabel("Command: ")

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tasksPanel, 0, 1, false).
		AddItem(cmdInput, 1, 1, true)

	// Update tasks display
	updateTasksPanel := func() {
		var sb strings.Builder
		sb.WriteString("[::b]Current Tasks:[::-]\n")
		sb.WriteString(fmt.Sprintf("Current Time: %s\n\n", time.Now().Format("15:04:05")))

		for _, task := range taskManager.Tasks {
			status := fmt.Sprintf("%s - %s (%s", task.Name, task.Remaining.Round(time.Second), task.State)
			if task.Loop {
				status += ", loop"
			}
			status += ")"
			sb.WriteString(status + "\n")
		}

		tasksPanel.SetText(sb.String())
	}

	// Command handling
	cmdInput.SetDoneFunc(func(key tcell.Key) {
		cmd := cmdInput.GetText()
		cmdInput.SetText("")

		parts := strings.Split(cmd, " ")
		if len(parts) < 1 {
			return
		}

		switch parts[0] {
		case "add":
			if len(parts) < 3 {
				return
			}
			name := parts[1]
			duration, err := time.ParseDuration(parts[2])
			if err != nil {
				return
			}
			taskManager.AddTask(name, duration)

		case "start":
			if len(parts) < 2 {
				return
			}
			id := parseID(parts[1])
			taskManager.StartTask(id)
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