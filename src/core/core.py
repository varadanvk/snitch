import time
import os
import logging
from datetime import datetime
from typing import Optional, List, Dict, Any, Tuple
import mss
import numpy as np
from PIL import Image
import threading

from src.utils.utils import TwilioHelper, ActivityTracker, save_config, load_config
from src.ml.ml import (
    ScreenAnalyzer,
    ActivityClassifier,
    MessageGenerator,
    ActivityHistory,
    OllamaInterface,
)

# Configure logger
logger = logging.getLogger("snitch.core")

# Default paths
DEFAULT_CONFIG_DIR = os.path.expanduser("~/.snitch")
DEFAULT_CONFIG_PATH = os.path.join(DEFAULT_CONFIG_DIR, "config.json")
DEFAULT_HISTORY_PATH = os.path.join(DEFAULT_CONFIG_DIR, "history.json")
DEFAULT_SCREENSHOTS_DIR = os.path.join(DEFAULT_CONFIG_DIR, "screenshots")


class ScreenMonitor:
    def __init__(self):
        """Initialize the screen monitor."""
        self.sct = mss.mss()
        self._last_capture: Optional[np.ndarray] = None
        self._last_capture_time: Optional[datetime] = None

        # Create screenshot directory if it doesn't exist
        os.makedirs(DEFAULT_SCREENSHOTS_DIR, exist_ok=True)

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

    def save_screenshot(self) -> Optional[str]:
        """Save the last captured screenshot to disk."""
        if self._last_capture is None:
            return None

        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = os.path.join(DEFAULT_SCREENSHOTS_DIR, f"screenshot_{timestamp}.png")

        try:
            Image.fromarray(self._last_capture).save(filename)
            return filename
        except Exception as e:
            logger.error(f"Error saving screenshot: {e}")
            return None


class TaskManager:
    def __init__(self):
        """Initialize the task manager."""
        self.current_task: Optional[str] = None
        self.task_history: List[tuple[str, datetime, Optional[datetime]]] = []
        self.productive_apps: List[str] = []
        self.distracting_apps: List[str] = []
        self._current_task_start: Optional[datetime] = None

        # Load configuration if available
        self._load_config()

    def set_current_task(self, task: str):
        """Set the current task and log the previous one."""
        if self.current_task:
            self.task_history.append(
                (self.current_task, self._current_task_start, datetime.now())
            )
        self.current_task = task
        self._current_task_start = datetime.now()

        # Save the updated configuration
        self._save_config()

    def add_productive_app(self, app_name: str):
        """Add an application to the productive list."""
        if app_name not in self.productive_apps:
            self.productive_apps.append(app_name)
            self._save_config()

    def add_distracting_app(self, app_name: str):
        """Add an application to the distracting list."""
        if app_name not in self.distracting_apps:
            self.distracting_apps.append(app_name)
            self._save_config()

    def get_task_history(self) -> List[tuple[str, datetime, Optional[datetime]]]:
        """Get the history of completed tasks."""
        return self.task_history.copy()

    def is_app_productive(self, app_name: str) -> bool:
        """Check if an application is marked as productive."""
        return app_name in self.productive_apps

    def is_app_distracting(self, app_name: str) -> bool:
        """Check if an application is marked as distracting."""
        return app_name in self.distracting_apps

    def _load_config(self) -> None:
        """Load task manager configuration from disk."""
        config = load_config(DEFAULT_CONFIG_PATH)
        if config:
            self.productive_apps = config.get("productive_apps", [])
            self.distracting_apps = config.get("distracting_apps", [])

    def _save_config(self) -> None:
        """Save task manager configuration to disk."""
        config = {
            "productive_apps": self.productive_apps,
            "distracting_apps": self.distracting_apps,
        }
        save_config(config, DEFAULT_CONFIG_PATH)


class NotificationManager:
    """Manages local notifications and app alerts."""

    def __init__(self, on_notification=None):
        """
        Initialize the notification manager.

        Args:
            on_notification: Callback function to call when a notification is sent
        """
        self.message_generator = MessageGenerator()
        self.on_notification = on_notification
        self.notification_history = []
        self.last_notification_time = None
        self.min_notification_interval = 6  # Minimum seconds between notifications

    def send_notification(
        self, message_type: str, context: Dict[str, Any] = None
    ) -> bool:
        """
        Send a notification.

        Args:
            message_type: Type of message ("distracted", "productive", "reminder")
            context: Additional context for the notification

        Returns:
            True if notification was sent, False otherwise
        """
        now = datetime.now()

        # Check if we should rate limit
        if (
            self.last_notification_time
            and (now - self.last_notification_time).total_seconds()
            < self.min_notification_interval
        ):
            return False

        # Generate message
        message = self.message_generator.generate_message(message_type, context)

        # Record this notification
        notification = {
            "timestamp": now,
            "message_type": message_type,
            "message": message,
            "context": context,
        }
        self.notification_history.append(notification)
        self.last_notification_time = now

        # Call the notification callback if provided
        if self.on_notification:
            self.on_notification(message, message_type)

        return True

    def get_recent_notifications(self, count: int = 10) -> List[Dict[str, Any]]:
        """Get the most recent notifications."""
        return self.notification_history[-count:]


class AccountabilityManager:
    """Manages accountability buddies and notifications."""

    def __init__(self):
        """Initialize the accountability manager."""
        self.twilio_helper = TwilioHelper()
        self.buddies = []
        self.snitch_mode_enabled = False
        self.last_snitch_time = None
        self.min_snitch_interval = 3600  # Minimum seconds between snitches (1 hour)

        # Load configuration
        self._load_config()

    def add_buddy(self, name: str, phone_number: str) -> bool:
        """Add an accountability buddy."""
        for buddy in self.buddies:
            if buddy["phone_number"] == phone_number:
                return False

        self.buddies.append(
            {"name": name, "phone_number": phone_number, "enabled": True}
        )

        self._save_config()
        return True

    def remove_buddy(self, phone_number: str) -> bool:
        """Remove an accountability buddy."""
        initial_count = len(self.buddies)
        self.buddies = [b for b in self.buddies if b["phone_number"] != phone_number]

        if len(self.buddies) < initial_count:
            self._save_config()
            return True
        return False

    def toggle_buddy(self, phone_number: str, enabled: bool) -> bool:
        """Enable or disable a specific buddy."""
        for buddy in self.buddies:
            if buddy["phone_number"] == phone_number:
                buddy["enabled"] = enabled
                self._save_config()
                return True
        return False

    def toggle_snitch_mode(self, enabled: bool) -> None:
        """Enable or disable snitch mode."""
        self.snitch_mode_enabled = enabled
        self._save_config()

    def snitch(self, context: Dict[str, Any]) -> bool:
        """
        Send a snitch message to all enabled buddies.

        Returns True if at least one message was sent.
        """
        if not self.snitch_mode_enabled or not self.buddies:
            return False

        now = datetime.now()

        # Check if we should rate limit
        if (
            self.last_snitch_time
            and (now - self.last_snitch_time).total_seconds() < self.min_snitch_interval
        ):
            return False

        # Generate the snitch message
        task = context.get("current_task", "their work")
        activity = context.get("activity", "getting distracted")

        message = f"🚨 SNITCH ALERT: Your buddy should be working on {task} but they're {activity} instead!"

        # Send to all enabled buddies
        messages_sent = 0
        for buddy in self.buddies:
            if buddy.get("enabled", True):
                success = self.twilio_helper.send_message(
                    buddy["phone_number"], message
                )
                if success:
                    messages_sent += 1

        # Update last snitch time if any messages were sent
        if messages_sent > 0:
            self.last_snitch_time = now
            return True

        return False

    def _load_config(self) -> None:
        """Load accountability manager configuration."""
        config = load_config(DEFAULT_CONFIG_PATH) or {}
        accountability_config = config.get("accountability", {})

        self.buddies = accountability_config.get("buddies", [])
        self.snitch_mode_enabled = accountability_config.get(
            "snitch_mode_enabled", False
        )

    def _save_config(self) -> None:
        """Save accountability manager configuration."""
        config = load_config(DEFAULT_CONFIG_PATH) or {}

        config["accountability"] = {
            "buddies": self.buddies,
            "snitch_mode_enabled": self.snitch_mode_enabled,
        }

        save_config(config, DEFAULT_CONFIG_PATH)


class ConfigManager:
    """Manages app settings and preferences."""

    def __init__(self):
        """Initialize the configuration manager."""
        # Create config directory if it doesn't exist
        os.makedirs(DEFAULT_CONFIG_DIR, exist_ok=True)

        # Initialize with default values
        self.config = {
            "monitoring_interval": 5,  # Seconds between screen captures
            "notification_interval": 10,  # Seconds between notifications
            "sensitivity": "medium",  # Productivity detection sensitivity
            "focused_hours": {"start": 9, "end": 17},  # Default work hours
            "theme": "system",  # UI theme
            "save_screenshots": False,  # Whether to save screenshots
        }

        # Load existing configuration if available
        self._load_config()

    def get_config(self) -> Dict[str, Any]:
        """Get the current configuration."""
        return self.config.copy()

    def update_config(self, updates: Dict[str, Any]) -> None:
        """Update configuration with new values."""
        self.config.update(updates)
        self._save_config()

    def reset_to_defaults(self) -> None:
        """Reset configuration to default values."""
        # Initialize with default values again
        self.config = {
            "monitoring_interval": 5,
            "notification_interval": 10,
            "sensitivity": "medium",
            "focused_hours": {"start": 9, "end": 17},
            "theme": "system",
            "save_screenshots": False,
        }
        self._save_config()

    def _load_config(self) -> None:
        """Load configuration from file."""
        config = load_config(DEFAULT_CONFIG_PATH)
        if config and "app_settings" in config:
            self.config.update(config["app_settings"])

    def _save_config(self) -> None:
        """Save configuration to file."""
        config = load_config(DEFAULT_CONFIG_PATH) or {}
        config["app_settings"] = self.config
        save_config(config, DEFAULT_CONFIG_PATH)


class SnitchCore:
    """Core manager that coordinates all components of the Snitch app."""

    def __init__(self, notification_callback=None):
        """Initialize the Snitch core."""
        self.config_manager = ConfigManager()
        self.screen_monitor = ScreenMonitor()
        self.task_manager = TaskManager()
        self.screen_analyzer = ScreenAnalyzer()
        self.activity_classifier = ActivityClassifier(
            self.task_manager.productive_apps, self.task_manager.distracting_apps
        )
        self.notification_manager = NotificationManager(notification_callback)
        self.accountability_manager = AccountabilityManager()
        self.activity_history = ActivityHistory()

        self.monitoring_thread = None
        self.is_monitoring = False
        self.activity_tracker = ActivityTracker(DEFAULT_HISTORY_PATH)

        logger.info("SnitchCore initialized successfully")

    def start_monitoring(self) -> bool:
        """
        Start the monitoring thread.

        Returns True if monitoring was started, False if already running.
        """
        if self.is_monitoring:
            return False

        self.is_monitoring = True
        self.monitoring_thread = threading.Thread(target=self._monitoring_loop)
        self.monitoring_thread.daemon = True
        self.monitoring_thread.start()

        logger.info("Monitoring started")
        return True

    def stop_monitoring(self) -> bool:
        """
        Stop the monitoring thread.

        Returns True if monitoring was stopped, False if not running.
        """
        if not self.is_monitoring:
            return False

        self.is_monitoring = False
        if self.monitoring_thread:
            self.monitoring_thread.join(timeout=1.0)

        logger.info("Monitoring stopped")
        return True

    def set_task(self, task: str) -> None:
        """Set the current task."""
        self.task_manager.set_current_task(task)
        logger.info(f"Task set: {task}")

        # Start monitoring if not already running
        if not self.is_monitoring:
            self.start_monitoring()

    def add_productive_app(self, app_name: str) -> None:
        """Add an app to the productive list."""
        self.task_manager.add_productive_app(app_name)
        self.activity_classifier.add_productive_app(app_name)

    def add_distracting_app(self, app_name: str) -> None:
        """Add an app to the distracting list."""
        self.task_manager.add_distracting_app(app_name)
        self.activity_classifier.add_distracting_app(app_name)

    def add_accountability_buddy(self, name: str, phone_number: str) -> bool:
        """Add an accountability buddy."""
        return self.accountability_manager.add_buddy(name, phone_number)

    def toggle_snitch_mode(self, enabled: bool) -> None:
        """Enable or disable snitch mode."""
        self.accountability_manager.toggle_snitch_mode(enabled)

    def get_productivity_summary(self) -> Dict[str, Any]:
        """Get a summary of productivity."""
        return self.activity_tracker.get_daily_summary()

    def _monitoring_loop(self) -> None:
        """Main monitoring loop that runs in a separate thread."""
        last_notification_time = datetime.now()
        config = self.config_manager.get_config()
        monitoring_interval = config["monitoring_interval"]
        notification_interval = config["notification_interval"]

        while self.is_monitoring:
            try:
                # Capture and analyze the screen
                screenshot = self.screen_monitor.capture_screen()
                analysis = self.screen_analyzer.analyze_screenshot(screenshot)

                # Get the current classification
                activity_type = analysis["activity_type"]
                activity = analysis["activity"]
                is_productive = activity_type == "productive"

                # Save to history
                self.activity_history.add_activity(
                    datetime.now(), activity_type, is_productive, {"activity": activity}
                )

                # Track in activity log
                self.activity_tracker.save_activity(
                    "productivity",
                    {
                        "productive": is_productive,
                        "activity": activity,
                        "duration": monitoring_interval,
                    },
                )

                # Check if we should send a notification
                now = datetime.now()
                notification_due = (
                    now - last_notification_time
                ).total_seconds() >= notification_interval

                if notification_due and not is_productive:
                    context = {
                        "current_task": self.task_manager.current_task,
                        "activity": activity,
                    }

                    # Send notification
                    notification_sent = self.notification_manager.send_notification(
                        "distracted", context
                    )

                    if notification_sent:
                        last_notification_time = now

                        # Check if we should snitch
                        self.accountability_manager.snitch(context)

                # Save screenshots if enabled
                if config["save_screenshots"]:
                    self.screen_monitor.save_screenshot()

                # Sleep for the monitoring interval
                time.sleep(monitoring_interval)

            except Exception as e:
                logger.error(f"Error in monitoring loop: {e}")
                time.sleep(monitoring_interval)

        logger.info("Monitoring loop terminated")
