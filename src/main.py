#!/usr/bin/env python3
import os
import sys
import logging
import customtkinter as ctk
from src.core.core import SnitchCore
from src.ui.ui import ModernUI

# Ensure log directory exists
log_dir = os.path.expanduser("~/.snitch")
os.makedirs(log_dir, exist_ok=True)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    handlers=[logging.StreamHandler(), logging.FileHandler(os.path.join(log_dir, "snitch.log"))],
)
logger = logging.getLogger("snitch.main")


class SnitchApp:
    """Main Snitch application class."""
    
    def __init__(self):
        """Initialize the Snitch application."""
        # Initialize core functionality
        self.core = SnitchCore(notification_callback=self.handle_notification)
        
        # Set up UI
        self.root = ctk.CTk()
        self.ui = ModernUI(self.root)
        
        # Set up callbacks
        self.setup_callbacks()
        
        logger.info("Snitch application initialized")
    
    def setup_callbacks(self):
        """Set up the callbacks between UI and core."""
        self.ui.initialize(self.core.set_task)
    
    def handle_notification(self, message, message_type):
        """Handle notifications from the core."""
        self.ui.update_status(message)
        self.ui.add_activity(f"{message_type.capitalize()}: {message}")
    
    def run(self):
        """Run the application."""
        try:
            # Start the GUI
            self.root.mainloop()
        except KeyboardInterrupt:
            logger.info("Application terminated by user")
        except Exception as e:
            logger.error(f"Unhandled exception: {e}")
        finally:
            # Cleanup
            if self.core.is_monitoring:
                self.core.stop_monitoring()
            logger.info("Application shutdown")


def main():
    """Main entry point for the application."""
    try:
        # Set environment variable for high DPI scaling
        os.environ["CUSTOMTKINTER_DPI_SCALE"] = "1.0"
        
        # Create and run the app
        app = SnitchApp()
        app.run()
        
        return 0
    except Exception as e:
        logger.error(f"Fatal error: {e}")
        return 1


if __name__ == "__main__":
    sys.exit(main())