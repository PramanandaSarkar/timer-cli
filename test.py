import time
import threading
from playsound import playsound
from rich.console import Console
from rich.table import Table

# Task Manager Class
class TimerCLI:
    def __init__(self):
        self.tasks = {}
        self.current_task = None
        self.lock = threading.Lock()
        self.console = Console()

    def add_task(self, name):
        if name in self.tasks:
            self.console.print(f"Task '{name}' already exists!", style="bold red")
            return
        self.tasks[name] = {"time": 0, "status": "paused", "loop": False}
        self.console.print(f"Task '{name}' added!", style="bold green")

    def start_task(self, name):
        if self.current_task:
            self.console.print(f"Task '{self.current_task}' is already running. Pause it first.", style="bold red")
            return
        if name not in self.tasks:
            self.console.print(f"Task '{name}' does not exist!", style="bold red")
            return

        self.current_task = name
        self.tasks[name]["status"] = "running"
        threading.Thread(target=self.run_task, args=(name,), daemon=True).start()
        self.console.print(f"Task '{name}' started!", style="bold green")

    def pause_task(self, name):
        if name not in self.tasks or self.tasks[name]["status"] != "running":
            self.console.print(f"Task '{name}' is not running!", style="bold red")
            return

        self.tasks[name]["status"] = "paused"
        self.current_task = None
        self.console.print(f"Task '{name}' paused.", style="bold yellow")

    def loop_task(self, name):
        if name not in self.tasks:
            self.console.print(f"Task '{name}' does not exist!", style="bold red")
            return
        self.tasks[name]["loop"] = not self.tasks[name]["loop"]
        status = "looping" if self.tasks[name]["loop"] else "not looping"
        self.console.print(f"Task '{name}' is now {status}.", style="bold green")

    def run_task(self, name):
        while self.current_task == name:
            with self.lock:
                self.tasks[name]["time"] += 1
            time.sleep(1)

            # Notify if looping is enabled and task reaches a threshold (example: 10 sec)
            if self.tasks[name]["loop"] and self.tasks[name]["time"] % 10 == 0:
                playsound("notification.mp3")
                self.console.print(f"Loop alert for task '{name}'!", style="bold magenta")

    def show_tasks(self):
        table = Table(title="Task Manager")
        table.add_column("Task", justify="left")
        table.add_column("Time (sec)", justify="center")
        table.add_column("Status", justify="center")

        for task, info in self.tasks.items():
            table.add_row(task, str(info["time"]), info["status"] + (", loop" if info["loop"] else ""))
        self.console.print(table)

# Main Program
def main():
    timer_cli = TimerCLI()
    timer_cli.add_task("relaxing")
    timer_cli.start_task("relaxing")
    time.sleep(5)
    timer_cli.pause_task("relaxing")
    timer_cli.show_tasks()

if __name__ == "__main__":
    main()
