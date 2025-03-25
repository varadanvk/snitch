#!/usr/bin/env python3
"""
Run script for Snitch.
This is a convenience script to run the Snitch application.
"""
import os
import sys
from pathlib import Path

# Add the project directory to the path
current_dir = Path(__file__).parent
sys.path.insert(0, str(current_dir))

# Import the main module and run it
from src.main import main

if __name__ == "__main__":
    sys.exit(main())