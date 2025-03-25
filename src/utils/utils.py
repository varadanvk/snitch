import os
import logging
import json
from typing import Dict, Any, Optional
from datetime import datetime
import mss
import numpy as np
from PIL import Image
import dotenv
from twilio.rest import Client


# Set up logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler()],
)
logger = logging.getLogger("snitch")


def capture_screenshot() -> np.ndarray:
    """Capture a screenshot of the primary monitor."""
    with mss.mss() as sct:
        screenshot = sct.grab(sct.monitors[1])  # Primary monitor
        img = Image.frombytes("RGB", screenshot.size, screenshot.rgb)
        return np.array(img)


def save_screenshot(img: np.ndarray, path: str) -> str:
    """Save a screenshot to disk."""
    if not os.path.exists(os.path.dirname(path)):
        os.makedirs(os.path.dirname(path), exist_ok=True)
    
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"{path}/screenshot_{timestamp}.png"
    
    Image.fromarray(img).save(filename)
    return filename


def load_config(config_path: str) -> Dict[str, Any]:
    """Load configuration from a JSON file."""
    try:
        if os.path.exists(config_path):
            with open(config_path, "r") as f:
                return json.load(f)
        return {}
    except Exception as e:
        logger.error(f"Error loading config: {e}")
        return {}


def save_config(config: Dict[str, Any], config_path: str) -> bool:
    """Save configuration to a JSON file."""
    try:
        os.makedirs(os.path.dirname(config_path), exist_ok=True)
        with open(config_path, "w") as f:
            json.dump(config, f, indent=2)
        return True
    except Exception as e:
        logger.error(f"Error saving config: {e}")
        return False


class TwilioHelper:
    """Helper class for Twilio integration."""
    
    def __init__(self):
        """Initialize Twilio client if environment variables are set."""
        dotenv.load_dotenv()
        self.account_sid = os.getenv("TWILIO_ACCOUNT_SID")
        self.auth_token = os.getenv("TWILIO_AUTH_TOKEN")
        self.from_number = os.getenv("TWILIO_PHONE_NUMBER")
        
        self.client = None
        if self.account_sid and self.auth_token and self.from_number:
            try:
                self.client = Client(self.account_sid, self.auth_token)
                logger.info("Twilio client initialized successfully")
            except Exception as e:
                logger.error(f"Failed to initialize Twilio client: {e}")
    
    def send_message(self, to_number: str, message: str) -> bool:
        """Send a text message using Twilio."""
        if not self.client or not self.from_number:
            logger.error("Twilio client not initialized")
            return False
        
        try:
            self.client.messages.create(
                body=message,
                from_=self.from_number,
                to=to_number
            )
            logger.info(f"Message sent to {to_number}")
            return True
        except Exception as e:
            logger.error(f"Failed to send message: {e}")
            return False


class ActivityTracker:
    """Helper class for tracking and analyzing user activity."""
    
    def __init__(self, history_file: str):
        """Initialize the activity tracker."""
        self.history_file = history_file
        self.history = self._load_history()
    
    def _load_history(self) -> Dict[str, Any]:
        """Load activity history from file."""
        return load_config(self.history_file)
    
    def save_activity(self, activity_type: str, details: Dict[str, Any]) -> None:
        """Save an activity to the history."""
        if 'activities' not in self.history:
            self.history['activities'] = []
        
        activity = {
            'timestamp': datetime.now().isoformat(),
            'type': activity_type,
            **details
        }
        
        self.history['activities'].append(activity)
        save_config(self.history, self.history_file)
    
    def get_daily_summary(self) -> Dict[str, Any]:
        """Get a summary of today's activities."""
        today = datetime.now().date().isoformat()
        
        today_activities = [
            activity for activity in self.history.get('activities', [])
            if activity.get('timestamp', '').startswith(today)
        ]
        
        productive_time = 0
        distracting_time = 0
        
        for activity in today_activities:
            if activity.get('type') == 'productivity':
                if activity.get('productive', False):
                    productive_time += activity.get('duration', 0)
                else:
                    distracting_time += activity.get('duration', 0)
        
        return {
            'date': today,
            'productive_time': productive_time,
            'distracting_time': distracting_time,
            'activities': len(today_activities)
        }


class CustomException(Exception):
    """Base class for custom exceptions in the Snitch app."""
    pass


class ConfigError(CustomException):
    """Exception raised for errors in the configuration."""
    pass


class MLError(CustomException):
    """Exception raised for errors in ML processing."""
    pass


class NotificationError(CustomException):
    """Exception raised for errors in sending notifications."""
    pass