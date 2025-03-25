import os
import sys
import unittest
from unittest.mock import MagicMock, patch
import numpy as np
import json

# Add the src directory to the path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from src.ml.ml import (
    OllamaInterface,
    ScreenAnalyzer,
    ActivityClassifier,
    MessageGenerator,
    SASSY_MESSAGES
)


class TestOllamaInterface(unittest.TestCase):
    """Tests for the OllamaInterface class."""
    
    @patch("src.ml.ml.requests.get")
    def test_check_availability(self, mock_get):
        """Test checking if Ollama is available."""
        # Set up mock response for available Ollama
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_get.return_value = mock_response
        
        # Create the interface and check availability
        interface = OllamaInterface()
        
        # Verify the result
        self.assertTrue(interface.is_available)
        mock_get.assert_called_once()
        
        # Set up mock response for unavailable Ollama
        mock_get.reset_mock()
        mock_response.status_code = 404
        
        # Create a new interface
        interface = OllamaInterface()
        interface.is_available = False  # Reset for testing
        
        # Check availability again
        result = interface._check_availability()
        
        # Verify the result
        self.assertFalse(result)
        mock_get.assert_called_once()
    
    @patch("src.ml.ml.requests.post")
    def test_generate_analysis(self, mock_post):
        """Test generating analysis with Ollama."""
        # Set up mock response
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"response": "Test analysis"}
        mock_post.return_value = mock_response
        
        # Create the interface with availability mocked
        interface = OllamaInterface()
        interface.is_available = True
        
        # Generate analysis
        result = interface.generate_analysis("Test prompt")
        
        # Verify the result
        self.assertEqual(result, "Test analysis")
        mock_post.assert_called_once()
        
        # Test error handling
        mock_post.reset_mock()
        mock_response.status_code = 500
        
        # Generate analysis again
        result = interface.generate_analysis("Test prompt")
        
        # Verify error handling
        self.assertIn("Error", result)
        mock_post.assert_called_once()


class TestScreenAnalyzer(unittest.TestCase):
    """Tests for the ScreenAnalyzer class."""
    
    def test_analyze_screenshot(self):
        """Test analyzing a screenshot."""
        # Create the analyzer
        analyzer = ScreenAnalyzer()
        
        # Create a dummy screenshot
        screenshot = np.zeros((800, 600, 3), dtype=np.uint8)
        
        # Analyze the screenshot
        result = analyzer.analyze_screenshot(screenshot)
        
        # Verify the result structure
        self.assertIn("activity_type", result)
        self.assertIn("activity", result)
        self.assertIn("confidence", result)
        
        # Check that the activity type is valid
        self.assertIn(result["activity_type"], ["productive", "distracting"])
        
        # Check that the confidence is between 0 and 1
        self.assertGreaterEqual(result["confidence"], 0.0)
        self.assertLessEqual(result["confidence"], 1.0)
    
    @patch("src.ml.ml.OllamaInterface.generate_analysis")
    def test_get_detailed_analysis(self, mock_generate):
        """Test getting detailed analysis with Ollama."""
        # Set up mock response
        mock_response = json.dumps({
            "activity_type": "productive",
            "activity": "coding",
            "confidence": 0.9,
            "reasoning": "The user is writing code in an IDE"
        })
        mock_generate.return_value = mock_response
        
        # Create the analyzer with a mocked Ollama interface
        analyzer = ScreenAnalyzer()
        analyzer.ollama.is_available = True
        
        # Create a dummy screenshot
        screenshot = np.zeros((800, 600, 3), dtype=np.uint8)
        
        # Get detailed analysis
        with patch("src.ml.ml.os.path.exists", return_value=True), \
             patch("src.ml.ml.os.remove"):
            result = analyzer.get_detailed_analysis(screenshot)
        
        # Verify the result
        self.assertEqual(result["activity_type"], "productive")
        self.assertEqual(result["activity"], "coding")
        self.assertEqual(result["confidence"], 0.9)
        self.assertEqual(result["reasoning"], "The user is writing code in an IDE")
        
        # Test fallback when Ollama fails
        mock_generate.side_effect = Exception("Test error")
        
        # Get detailed analysis again
        with patch("src.ml.ml.os.path.exists", return_value=True), \
             patch("src.ml.ml.os.remove"):
            result = analyzer.get_detailed_analysis(screenshot)
        
        # Verify fallback to simpler analysis
        self.assertIn("activity_type", result)
        self.assertIn("activity", result)
        self.assertIn("confidence", result)


class TestActivityClassifier(unittest.TestCase):
    """Tests for the ActivityClassifier class."""
    
    def test_classify_activity(self):
        """Test classifying an activity."""
        # Create a classifier with known apps
        productive_apps = ["vscode", "zoom", "terminal"]
        distracting_apps = ["twitter", "youtube", "facebook"]
        classifier = ActivityClassifier(productive_apps, distracting_apps)
        
        # Test a known productive app
        is_productive, confidence = classifier.classify_activity("vscode", "Working on project")
        self.assertTrue(is_productive)
        self.assertGreaterEqual(confidence, 0.9)
        
        # Test a known distracting app
        is_productive, confidence = classifier.classify_activity("youtube", "Watching videos")
        self.assertFalse(is_productive)
        self.assertGreaterEqual(confidence, 0.9)
        
        # Test an unknown app
        is_productive, confidence = classifier.classify_activity("unknown", "Random activity")
        self.assertIsInstance(is_productive, bool)
        self.assertGreaterEqual(confidence, 0.0)
        self.assertLessEqual(confidence, 1.0)
    
    def test_add_productive_app(self):
        """Test adding a productive app."""
        classifier = ActivityClassifier()
        
        # Add a productive app
        app_name = "productive_app"
        classifier.add_productive_app(app_name)
        
        # Verify the app was added
        self.assertIn(app_name, classifier.productive_apps)
    
    def test_add_distracting_app(self):
        """Test adding a distracting app."""
        classifier = ActivityClassifier()
        
        # Add a distracting app
        app_name = "distracting_app"
        classifier.add_distracting_app(app_name)
        
        # Verify the app was added
        self.assertIn(app_name, classifier.distracting_apps)


class TestMessageGenerator(unittest.TestCase):
    """Tests for the MessageGenerator class."""
    
    def test_generate_message(self):
        """Test generating a message."""
        # Create the generator with a mocked Ollama interface
        ollama = MagicMock()
        ollama.is_available = False
        generator = MessageGenerator(ollama)
        
        # Generate a message for a distracted user
        message = generator.generate_message("distracted")
        
        # Verify the message
        self.assertIsInstance(message, str)
        self.assertGreater(len(message), 0)
        
        # Check that the message is from the distracted category
        self.assertIn(message, SASSY_MESSAGES["distracted"])
        
        # Test with Ollama available
        ollama.reset_mock()
        ollama.is_available = True
        ollama.generate_analysis.return_value = "Custom generated message"
        
        # Generate a message with the mock
        with patch("src.ml.ml.random.random", return_value=0.5):
            message = generator.generate_message("distracted", {"current_task": "coding"})
        
        # Verify Ollama was used
        ollama.generate_analysis.assert_called_once()
        self.assertEqual(message, "Custom generated message")
        
        # Test fallback when Ollama fails
        ollama.reset_mock()
        ollama.generate_analysis.side_effect = Exception("Test error")
        
        # Generate a message with the mock
        with patch("src.ml.ml.random.random", return_value=0.5):
            message = generator.generate_message("distracted")
        
        # Verify fallback to predefined messages
        self.assertIn(message, SASSY_MESSAGES["distracted"])


if __name__ == "__main__":
    unittest.main()