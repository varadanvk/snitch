import os
from datetime import datetime, timedelta
import time
import random
import logging
import subprocess
import numpy as np
from PIL import Image
from typing import Dict, List, Tuple, Optional, Any
import requests
import json

# Configure logger
logger = logging.getLogger("snitch.ml")

# Sample activity classes for testing without a real ML model
ACTIVITY_CLASSES = {
    "productive": [
        "coding",
        "writing",
        "reading documentation",
        "email",
        "spreadsheet",
    ],
    "distracting": [
        "social media",
        "video streaming",
        "gaming",
        "shopping",
        "news browsing",
    ],
}

# Sample sassy messages for notifications
SASSY_MESSAGES = {
    "distracted": [
        "Busted! That doesn't look like work to me...",
        "I see you slacking. Your future self is disappointed.",
        "Eyes on the prize! That cat video can wait.",
        "Nice try. Get back to work, please.",
        "That's definitely not what you're supposed to be doing.",
    ],
    "productive": [
        "Look at you being all productive!",
        "Great focus! Keep it up!",
        "You're crushing it right now!",
        "This is what productivity looks like. Nice!",
        "Your future self thanks you for this focus!",
    ],
    "reminder": [
        "Hey, remember that task you were supposed to be working on?",
        "Just a gentle nudge to get back on track.",
        "Focus check: Are you still working on your task?",
        "Quick reminder of what you're supposed to be doing.",
        "Your goals are waiting... time to refocus!",
    ],
}


class OllamaInterface:
    """Interface for communicating with Ollama for local LLM inference."""

    def __init__(self, model: str = "llava:latest"):
        """Initialize the Ollama interface with the specified model."""
        self.model = model
        self.base_url = "http://localhost:11434/api"
        self.is_available = self._check_availability()

        if not self.is_available:
            logger.warning("Ollama is not available. Some features will be limited.")
        else:
            logger.info(f"Ollama is available, using model: {self.model}")

    def _check_availability(self) -> bool:
        """Check if Ollama is running and available."""
        try:
            response = requests.get(f"{self.base_url}/tags")
            if response.status_code == 200:
                # Check if the model exists
                models = response.json().get("models", [])
                model_exists = any(model.get("name", "").startswith(self.model.split(":")[0]) for model in models)
                
                if not model_exists:
                    logger.warning(f"Model {self.model} not found. Please run 'ollama pull {self.model}'")
                    logger.warning("Using a vision-capable model like 'llava' is recommended")
                    # We still return True since Ollama is available, even if the model isn't
                
                return True
            return False
        except requests.RequestException:
            return False

    def is_ollama_installed(self) -> bool:
        """Check if Ollama is installed."""
        try:
            result = subprocess.run(
                ["which", "ollama"], stdout=subprocess.PIPE, stderr=subprocess.PIPE
            )
            return result.returncode == 0
        except Exception:
            return False

    def start_ollama_if_needed(self) -> bool:
        """Start Ollama if it's installed but not running."""
        if self.is_available:
            return True

        if not self.is_ollama_installed():
            logger.warning(
                "Ollama is not installed. Please install it for better features."
            )
            return False

        try:
            # Try to start Ollama as a background process
            subprocess.Popen(
                ["ollama", "serve"],
                stdout=subprocess.DEVNULL,
                stderr=subprocess.DEVNULL,
            )

            # Wait for Ollama to start (with timeout)
            for _ in range(5):  # Try for 5 seconds
                time.sleep(1)
                if self._check_availability():
                    self.is_available = True
                    logger.info("Successfully started Ollama")
                    return True

            logger.warning("Failed to start Ollama within timeout period")
            return False
        except Exception as e:
            logger.error(f"Error starting Ollama: {e}")
            return False

    def generate_analysis(self, prompt: str, image_data: Optional[np.ndarray] = None) -> str:
        """
        Generate text using Ollama, optionally with an image.
        
        Args:
            prompt: Text prompt for the model
            image_data: Optional numpy array containing image data
        
        Returns:
            The model's response as a string
        """
        if not self.is_available and not self.start_ollama_if_needed():
            # Fallback to simple response if Ollama is not available
            logger.warning("Using fallback responses since Ollama is not available")
            return "Ollama is not available for advanced analysis."
            
        try:
            payload = {
                "model": self.model,
                "prompt": prompt,
                "stream": False,
            }
            
            # If image is provided, encode it as base64 and add to the payload
            if image_data is not None:
                import base64
                from io import BytesIO
                
                # Convert numpy array to PIL Image
                image = Image.fromarray(image_data)
                
                # Convert PIL image to base64 string
                buffered = BytesIO()
                image.save(buffered, format="JPEG")
                img_str = base64.b64encode(buffered.getvalue()).decode()
                
                # Add image to payload with proper format for Ollama
                payload["images"] = [img_str]
                
                logger.info("Sending image to Ollama for analysis")
            
            response = requests.post(
                f"{self.base_url}/generate",
                json=payload,
            )

            if response.status_code == 200:
                result = response.json()
                return result.get("response", "")
            else:
                logger.error(f"Error from Ollama API: {response.text}")
                return "Error generating analysis"
        except Exception as e:
            logger.error(f"Failed to generate analysis: {e}")
            return "Failed to communicate with Ollama"


class ScreenAnalyzer:
    """Analyzes screenshots to determine user activity."""

    def __init__(self):
        """Initialize the screen analyzer."""
        self.ollama = OllamaInterface()

    def analyze_screenshot(self, screenshot: np.ndarray) -> Dict[str, Any]:
        """
        Analyze a screenshot to determine what the user is doing.
        Uses Ollama to analyze the image if available.
        """
        # Try to use Ollama with the image for analysis
        if self.ollama.is_available or self.ollama.start_ollama_if_needed():
            try:
                # Prepare prompt for vision model analysis
                prompt = """
                Analyze this screenshot of my computer screen. What activity am I doing?
                Is it productive work or a distraction? Please be specific about what you see.
                Format your response as JSON with these fields:
                {  
                    "activity_type": "productive" or "distracting",
                    "activity": "specific activity description",
                    "confidence": 0.0 to 1.0,
                    "reasoning": "brief explanation"
                }
                """
                
                # Send the image to Ollama
                response = self.ollama.generate_analysis(prompt, screenshot)
                
                try:
                    # Parse the response as JSON
                    result = json.loads(response)
                    logger.info(f"Ollama screenshot analysis: {result}")
                    return result
                except json.JSONDecodeError:
                    # If response isn't valid JSON, try to extract JSON pattern
                    import re
                    json_pattern = re.search(r'\{[\s\S]*\}', response)
                    if json_pattern:
                        try:
                            result = json.loads(json_pattern.group(0))
                            return result
                        except:
                            pass
                    
                    logger.warning(f"Could not parse Ollama response as JSON: {response[:100]}...")
            except Exception as e:
                logger.error(f"Error using Ollama for screenshot analysis: {e}")
        
        # Fallback to simulated analysis if Ollama failed or isn't available
        logger.info("Using fallback simulated analysis")
        activity_type = random.choice(["productive", "distracting"])
        activity = random.choice(ACTIVITY_CLASSES[activity_type])
        confidence = random.uniform(0.7, 0.98)

        return {
            "activity_type": activity_type,
            "activity": activity,
            "confidence": confidence,
            "reasoning": "Fallback analysis (Ollama not available)"
        }

    def get_detailed_analysis(self, screenshot: np.ndarray) -> Dict[str, Any]:
        """
        Get a more detailed analysis of the screenshot using Ollama.

        This would be used for deeper insights, not regular monitoring.
        """
        if not self.ollama.is_available and not self.ollama.start_ollama_if_needed():
            # Return simulated analysis if Ollama is not available
            return self.analyze_screenshot(screenshot)

        prompt = """
        Provide a detailed analysis of this screenshot of my computer screen:
        1. What applications or websites are visible?
        2. What specific activity am I engaged in?
        3. Is this activity productive or distracting?
        4. What specific elements on screen suggest productive or distracting behavior?
        5. How confident are you in this assessment?
        
        Format your response as JSON with these fields:
        {
            "activity_type": "productive" or "distracting",
            "activity": "specific activity description",
            "applications": ["list", "of", "detected", "applications"],
            "confidence": 0.0 to 1.0,
            "reasoning": "detailed explanation",
            "suggestions": "optional improvement suggestions if distracting"
        }
        """

        try:
            # Send the image directly to Ollama
            response = self.ollama.generate_analysis(prompt, screenshot)
            
            try:
                # Parse the response as JSON
                result = json.loads(response)
                logger.info(f"Detailed analysis: {result}")
                return result
            except json.JSONDecodeError:
                # If response isn't valid JSON, try to extract JSON pattern
                import re
                json_pattern = re.search(r'\{[\s\S]*\}', response)
                if json_pattern:
                    try:
                        result = json.loads(json_pattern.group(0))
                        return result
                    except:
                        pass
                
                logger.warning(f"Could not parse detailed analysis as JSON: {response[:100]}...")
                # Create a structured response from the text
                return {
                    "activity_type": "unknown",
                    "activity": "unclassified activity",
                    "confidence": 0.5,
                    "reasoning": f"Parsing error. Raw response: {response[:200]}..."
                }
        except Exception as e:
            logger.error(f"Error in detailed analysis: {e}")
            # Fall back to simpler analysis
            return self.analyze_screenshot(screenshot)


class ActivityClassifier:
    """Determines if user activity is productive or a distraction."""

    def __init__(
        self, productive_apps: List[str] = None, distracting_apps: List[str] = None
    ):
        """Initialize with lists of productive and distracting apps."""
        self.productive_apps = productive_apps or []
        self.distracting_apps = distracting_apps or []

    def classify_activity(self, app_name: str, window_title: str) -> Tuple[bool, float]:
        """
        Classify an activity as productive or distracting based on app and window title.

        Returns:
        - is_productive: Boolean indicating if the activity is productive
        - confidence: Confidence score between 0 and 1
        """
        # Check if the app is in our predefined lists
        if any(app.lower() in app_name.lower() for app in self.productive_apps):
            return True, 0.9

        if any(app.lower() in app_name.lower() for app in self.distracting_apps):
            return False, 0.9

        # For demo purposes, simulate a classification
        # In a real app, you would use a trained model here
        is_productive = random.choice(
            [True, False, True]
        )  # Slight bias toward productive
        confidence = random.uniform(0.6, 0.85)

        return is_productive, confidence

    def add_productive_app(self, app_name: str) -> None:
        """Add an app to the productive list."""
        if app_name not in self.productive_apps:
            self.productive_apps.append(app_name)

    def add_distracting_app(self, app_name: str) -> None:
        """Add an app to the distracting list."""
        if app_name not in self.distracting_apps:
            self.distracting_apps.append(app_name)


class MessageGenerator:
    """Generates personalized notification messages."""

    def __init__(self, ollama_interface: Optional[OllamaInterface] = None):
        """Initialize with optional Ollama interface for advanced messaging."""
        self.ollama = ollama_interface or OllamaInterface()

    def generate_message(
        self, message_type: str, context: Dict[str, Any] = None
    ) -> str:
        """
        Generate a message based on the type and context.

        Args:
            message_type: Type of message ("distracted", "productive", "reminder")
            context: Additional context for personalization

        Returns:
            A personalized message
        """
        context = context or {}

        # Try to generate using Ollama if available
        if (
            self.ollama.is_available and random.random() < 0.7
        ):  # 70% chance to use Ollama
            try:
                task = context.get("current_task", "your task")
                activity = context.get("activity", "something")

                prompt = f"""
                Generate a short, friendly but sassy notification to the user who is supposed 
                to be working on "{task}" but is actually {activity}. 
                Keep it under 100 characters, be motivational but with a touch of humor.
                Don't use hashtags or emojis.
                """

                response = self.ollama.generate_analysis(prompt, None)

                # If we got a reasonable response, use it
                if response and 5 < len(response) < 150:
                    return response.strip()
            except Exception as e:
                logger.error(f"Error generating personalized message: {e}")

        # Fall back to predefined messages
        if message_type in SASSY_MESSAGES:
            return random.choice(SASSY_MESSAGES[message_type])

        # Default fallback
        return "Hey, focus on your work!"


class ActivityHistory:
    """Tracks and analyzes patterns of user behavior."""

    def __init__(self, max_history: int = 100):
        """Initialize with a maximum history size."""
        self.activities = []
        self.max_history = max_history

    def add_activity(
        self,
        timestamp: datetime,
        activity_type: str,
        is_productive: bool,
        details: Dict[str, Any] = None,
    ) -> None:
        """Add an activity to the history."""
        details = details or {}

        activity = {
            "timestamp": timestamp,
            "activity_type": activity_type,
            "is_productive": is_productive,
            **details,
        }

        self.activities.append(activity)

        # Trim history if needed
        if len(self.activities) > self.max_history:
            self.activities = self.activities[-self.max_history :]

    def get_productivity_ratio(self, timeframe_hours: int = 24) -> float:
        """
        Calculate the productivity ratio over a given timeframe.

        Returns a value between 0 and 1 where 1 is 100% productive.
        """
        if not self.activities:
            return 0.5  # Default when no data

        # Filter activities within the timeframe
        cutoff_time = datetime.now() - timedelta(hours=timeframe_hours)
        recent_activities = [a for a in self.activities if a["timestamp"] > cutoff_time]

        if not recent_activities:
            return 0.5

        # Count productive vs. total activities
        productive_count = sum(1 for a in recent_activities if a["is_productive"])

        return productive_count / len(recent_activities)

    def identify_patterns(self) -> Dict[str, Any]:
        """
        Identify patterns in user behavior.

        Returns insights about productivity patterns.
        """
        if not self.activities:
            return {"message": "Not enough data to identify patterns"}

        # Simple analysis for demo
        # In a real app, you would use more sophisticated analysis
        productive_count = sum(1 for a in self.activities if a["is_productive"])
        productivity_ratio = productive_count / len(self.activities)

        # Analyze time of day patterns (morning, afternoon, evening)
        time_of_day = {
            "morning": {"productive": 0, "distracting": 0},
            "afternoon": {"productive": 0, "distracting": 0},
            "evening": {"productive": 0, "distracting": 0},
        }

        for activity in self.activities:
            hour = activity["timestamp"].hour

            if 5 <= hour < 12:
                period = "morning"
            elif 12 <= hour < 18:
                period = "afternoon"
            else:
                period = "evening"

            if activity["is_productive"]:
                time_of_day[period]["productive"] += 1
            else:
                time_of_day[period]["distracting"] += 1

        # Determine most and least productive periods
        most_productive = max(
            time_of_day.items(),
            key=lambda x: x[1]["productive"] / (sum(x[1].values()) or 1),
        )
        least_productive = min(
            time_of_day.items(),
            key=lambda x: x[1]["productive"] / (sum(x[1].values()) or 1),
        )

        return {
            "overall_productivity": productivity_ratio,
            "most_productive_period": most_productive[0],
            "least_productive_period": least_productive[0],
            "time_of_day": time_of_day,
        }
