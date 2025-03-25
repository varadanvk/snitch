#!/usr/bin/env python3
"""
Run script for Snitch.
This is a convenience script to run the Snitch application.
"""
import os
import sys
from pathlib import Path
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)

# Add the project directory to the path
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

# Import the main module and run it
from src.main import main

if __name__ == "__main__":
    print("Starting Snitch - AI Productivity Watchdog")
    print("Make sure Ollama is installed and running (with a vision-capable model)")
    print("Use 'ollama pull llava' and 'ollama serve' if not already set up")
    sys.exit(main())