from flask import Flask, render_template
from flask_socketio import SocketIO
from core.core import ScreenMonitor, TaskManager
import threading
import time

app = Flask(__name__)
socketio = SocketIO(app)


class SnitchApp:
    def __init__(self):
        self.screen_monitor = ScreenMonitor()
        self.task_manager = TaskManager()
        self.monitoring_thread = None
        self.is_monitoring = False

    def start_monitoring(self):
        """Start the monitoring thread."""
        self.is_monitoring = True
        while self.is_monitoring:
            self.screen_monitor.capture_screen()
            # TODO: Add ML-based screen analysis
            time.sleep(5)  # Check every 5 seconds

    def stop_monitoring(self):
        """Stop the monitoring thread."""
        self.is_monitoring = False
        if self.monitoring_thread:
            self.monitoring_thread.join()


# Initialize the app
snitch = SnitchApp()


@app.route("/")
def index():
    return render_template("tasks.html")


@app.route("/analytics")
def analytics():
    return render_template("analytics.html")


@app.route("/focus")
def focus():
    return render_template("focus.html")


@app.route("/settings")
def settings():
    return render_template("settings.html")


@socketio.on("connect")
def handle_connect():
    socketio.emit("status_update", {"message": "Connected to Snitch"})


@socketio.on("set_task")
def handle_set_task(data):
    task = data.get("task")
    if task:
        snitch.task_manager.set_current_task(task)
        socketio.emit("status_update", {"message": f"🎯 Currently working on: {task}"})
        socketio.emit("activity_update", {"message": f"Started task: {task}"})

        # Start monitoring if not already running
        if not snitch.monitoring_thread or not snitch.monitoring_thread.is_alive():
            snitch.monitoring_thread = threading.Thread(target=snitch.start_monitoring)
            snitch.monitoring_thread.daemon = True
            snitch.monitoring_thread.start()


def run_app():
    """Run the Flask application."""
    socketio.run(app, debug=True)


if __name__ == "__main__":
    run_app()
