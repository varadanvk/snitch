import os
import sys
import unittest
from unittest.mock import MagicMock, patch
import numpy as np
from datetime import datetime

# Add the src directory to the path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from src.core.core import ScreenMonitor, TaskManager, NotificationManager
from src.core.core import ConfigManager, SnitchCore


class TestScreenMonitor(unittest.TestCase):
    """Tests for the ScreenMonitor class."""
    
    @patch("src.core.core.mss.mss")
    def test_capture_screen(self, mock_mss):
        """Test that the screen capture works correctly."""
        # Set up mocks
        mock_sct = MagicMock()
        mock_mss.return_value = mock_sct
        
        # Create a fake screenshot
        mock_screenshot = MagicMock()
        mock_screenshot.size = (800, 600)
        mock_screenshot.rgb = b"fake_rgb_data" * 800 * 600 * 3
        mock_sct.grab.return_value = mock_screenshot
        
        # Create the screen monitor
        monitor = ScreenMonitor()
        
        # Capture the screen
        result = monitor.capture_screen()
        
        # Check that the correct methods were called
        mock_sct.grab.assert_called_once()
        
        # Check that we got a result
        self.assertIsNotNone(result)
        self.assertIsInstance(result, np.ndarray)


class TestTaskManager(unittest.TestCase):
    """Tests for the TaskManager class."""
    
    def test_set_current_task(self):
        """Test setting the current task."""
        manager = TaskManager()
        
        # Set an initial task
        task_name = "Test Task"
        manager.set_current_task(task_name)
        
        # Check that the task was set
        self.assertEqual(manager.current_task, task_name)
        
        # Set a new task and check that the history is updated
        new_task = "New Task"
        manager.set_current_task(new_task)
        
        # Check that the new task is set
        self.assertEqual(manager.current_task, new_task)
        
        # Check that the history was updated
        self.assertEqual(len(manager.task_history), 1)
        self.assertEqual(manager.task_history[0][0], task_name)
    
    def test_add_productive_app(self):
        """Test adding a productive app."""
        manager = TaskManager()
        
        # Add a productive app
        app_name = "productive_app"
        manager.add_productive_app(app_name)
        
        # Check that the app was added
        self.assertIn(app_name, manager.productive_apps)
        
        # Add the same app again and check that it doesn't duplicate
        manager.add_productive_app(app_name)
        self.assertEqual(manager.productive_apps.count(app_name), 1)
    
    def test_add_distracting_app(self):
        """Test adding a distracting app."""
        manager = TaskManager()
        
        # Add a distracting app
        app_name = "distracting_app"
        manager.add_distracting_app(app_name)
        
        # Check that the app was added
        self.assertIn(app_name, manager.distracting_apps)
        
        # Add the same app again and check that it doesn't duplicate
        manager.add_distracting_app(app_name)
        self.assertEqual(manager.distracting_apps.count(app_name), 1)


class TestNotificationManager(unittest.TestCase):
    """Tests for the NotificationManager class."""
    
    def test_send_notification(self):
        """Test sending a notification."""
        # Create a mock callback
        callback = MagicMock()
        
        # Create the notification manager
        manager = NotificationManager(callback)
        
        # Send a notification
        result = manager.send_notification("distracted", {"task": "test"})
        
        # Check that the notification was sent
        self.assertTrue(result)
        
        # Check that the callback was called
        callback.assert_called_once()
        
        # Check that the notification was added to the history
        self.assertEqual(len(manager.notification_history), 1)
        self.assertEqual(manager.notification_history[0]["message_type"], "distracted")


class TestConfigManager(unittest.TestCase):
    """Tests for the ConfigManager class."""
    
    @patch("src.core.core.load_config")
    @patch("src.core.core.save_config")
    def test_update_config(self, mock_save_config, mock_load_config):
        """Test updating the configuration."""
        # Set up mocks
        mock_load_config.return_value = {}
        
        # Create the config manager
        manager = ConfigManager()
        
        # Update the configuration
        updates = {"monitoring_interval": 10, "sensitivity": "high"}
        manager.update_config(updates)
        
        # Check that the config was updated
        self.assertEqual(manager.config["monitoring_interval"], 10)
        self.assertEqual(manager.config["sensitivity"], "high")
        
        # Check that save_config was called
        mock_save_config.assert_called()


class TestSnitchCore(unittest.TestCase):
    """Tests for the SnitchCore class."""
    
    @patch("src.core.core.ScreenMonitor")
    @patch("src.core.core.TaskManager")
    @patch("src.core.core.NotificationManager")
    @patch("src.core.core.ConfigManager")
    @patch("src.core.core.ScreenAnalyzer")
    @patch("src.core.core.ActivityClassifier")
    @patch("src.core.core.AccountabilityManager")
    @patch("src.core.core.ActivityHistory")
    @patch("src.core.core.ActivityTracker")
    def test_initialization(self, mock_tracker, mock_history, mock_accountability,
                           mock_classifier, mock_analyzer, mock_config,
                           mock_notification, mock_task, mock_screen):
        """Test that the core initializes its components correctly."""
        # Set up specific mocks
        mock_task_instance = MagicMock()
        mock_task.return_value = mock_task_instance
        mock_task_instance.productive_apps = ["app1"]
        mock_task_instance.distracting_apps = ["app2"]
        
        # Create the core
        core = SnitchCore()
        
        # Check that all components were initialized
        self.assertIsNotNone(core.config_manager)
        self.assertIsNotNone(core.screen_monitor)
        self.assertIsNotNone(core.task_manager)
        self.assertIsNotNone(core.screen_analyzer)
        self.assertIsNotNone(core.activity_classifier)
        self.assertIsNotNone(core.notification_manager)
        self.assertIsNotNone(core.accountability_manager)
        self.assertIsNotNone(core.activity_history)
        
    @patch("src.core.core.threading.Thread")
    def test_start_monitoring(self, mock_thread):
        """Test starting the monitoring thread."""
        # Create mocks
        mock_thread_instance = MagicMock()
        mock_thread.return_value = mock_thread_instance
        
        # Create a SnitchCore with mocked components
        core = SnitchCore()
        core.config_manager = MagicMock()
        core.screen_monitor = MagicMock()
        core.task_manager = MagicMock()
        core.notification_manager = MagicMock()
        core.accountability_manager = MagicMock()
        core.screen_analyzer = MagicMock()
        core.activity_history = MagicMock()
        core.activity_tracker = MagicMock()
        core.activity_classifier = MagicMock()
        
        # Start monitoring
        result = core.start_monitoring()
        
        # Check that monitoring was started
        self.assertTrue(result)
        self.assertTrue(core.is_monitoring)
        
        # Check that a thread was created and started
        mock_thread.assert_called_once()
        mock_thread_instance.start.assert_called_once()


if __name__ == "__main__":
    unittest.main()