import time
from datetime import datetime
from typing import Optional, List
import mss
import numpy as np
from PIL import Image


class ScreenMonitor:
    def __init__(self):
        self.sct = mss.mss()
        self._last_capture: Optional[np.ndarray] = None
        self._last_capture_time: Optional[datetime] = None

    def capture_screen(self) -> np.ndarray:
        """Capture the current screen state."""
        screenshot = self.sct.grab(self.sct.monitors[1])  # Primary monitor
        img = Image.frombytes("RGB", screenshot.size, screenshot.rgb)
        self._last_capture = np.array(img)
        self._last_capture_time = datetime.now()
        return self._last_capture

    def get_last_capture(self) -> Optional[np.ndarray]:
        """Get the most recent screen capture."""
        return self._last_capture

    def get_last_capture_time(self) -> Optional[datetime]:
        """Get the timestamp of the most recent capture."""
        return self._last_capture_time


class TaskManager:
    def __init__(self):
        self.current_task: Optional[str] = None
        self.task_history: List[tuple[str, datetime, Optional[datetime]]] = []
        self.productive_apps: List[str] = []
        self.distracting_apps: List[str] = []

    def set_current_task(self, task: str):
        """Set the current task and log the previous one."""
        if self.current_task:
            self.task_history.append(
                (self.current_task, self._current_task_start, datetime.now())
            )
        self.current_task = task
        self._current_task_start = datetime.now()

    def add_productive_app(self, app_name: str):
        """Add an application to the productive list."""
        if app_name not in self.productive_apps:
            self.productive_apps.append(app_name)

    def add_distracting_app(self, app_name: str):
        """Add an application to the distracting list."""
        if app_name not in self.distracting_apps:
            self.distracting_apps.append(app_name)

    def get_task_history(self) -> List[tuple[str, datetime, Optional[datetime]]]:
        """Get the history of completed tasks."""
        return self.task_history.copy()

    def is_app_productive(self, app_name: str) -> bool:
        """Check if an application is marked as productive."""
        return app_name in self.productive_apps

    def is_app_distracting(self, app_name: str) -> bool:
        """Check if an application is marked as distracting."""
        return app_name in self.distracting_apps
